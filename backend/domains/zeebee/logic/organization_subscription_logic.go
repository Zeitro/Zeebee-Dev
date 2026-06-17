package logic

import (
	"context"
	"database/sql"
	"fmt"
	rand2 "math/rand/v2"
	"time"

	"github.com/duke-git/lancet/v2/slice"
	"gorm.io/gorm"

	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/gormx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/repo"

	"github.com/zeromicro/go-zero/core/threading"

	"zeebee-backend/domains/zeebee/auth"
)

type OrganizationSubscriptionLogic struct {
	*BaseLogic
	authTopicMap map[constants.UserType]string
}

func NewOrganizationSubscriptionLogic(baseLogic *BaseLogic, authTopicMap map[constants.UserType]string) *OrganizationSubscriptionLogic {
	return &OrganizationSubscriptionLogic{
		BaseLogic:    baseLogic,
		authTopicMap: authTopicMap,
	}
}

// Create creates a new organization subscription.
func (l *OrganizationSubscriptionLogic) Create(ctx context.Context, input dto.OrganizationSubscriptionCreateDto) (result *dto.OrganizationSubscriptionEntity, err error) {
	plan, err := l.repos.PricingPlanRepo.GetOne(ctx, input.PricingPlanId, []gormx.Preload{{Query: "Resources"}})
	if err != nil {
		return nil, errorx.Wrap(err, "fail to find pricing plan")
	}
	if plan.Type != input.Type {
		err = errorx.New("pricing plan type not match")
		return
	}

	if plan.Status != constants.CommonStatusActive {
		return nil, errorx.New("pricing plan not found or disabled")
	}

	// H4: the organization's type must match the plan's type, otherwise the
	// subscription would bind plan resources to the wrong user-type roles.
	organization, err := l.repos.OrganizationRepo.GetOrganizationById(ctx, input.OrganizationId, false)
	if err != nil {
		return nil, errorx.Wrap(err, "fail to find organization")
	}
	if organization == nil || organization.Type != plan.Type {
		return nil, errorx.New("organization type not match pricing plan type")
	}

	policy, err := l.repos.PricingPlanPolicyRepo.GetOne(ctx, input.PricingPlanPolicyId)
	if err != nil {
		return nil, errorx.New("fail to find pricing plan policy")
	}

	if policy.PricingPlanId != input.PricingPlanId {
		err = errorx.New("pricing plan policy not match")
		return
	}

	code, err := l.generateCode(ctx)
	if err != nil {
		return nil, err
	}
	expireDate := calculateExpireDate(policy.Type, time.Now())

	subscription := &model.OrganizationSubscription{
		Code:                code,
		OrganizationId:      input.OrganizationId,
		PricingPlanId:       input.PricingPlanId,
		PricingPlanPolicyId: input.PricingPlanPolicyId,
		Status:              constants.OrganizationSubscriptionStatusActive,
		ExpireDate:          expireDate,
		OriginalAmount:      policy.OriginalAmount,
		ActualAmount:        policy.ActualAmount,
	}

	var organizationPolicies map[int64]*organizationPolicyChanges
	err = l.repos.OrganizationSubscriptionRepo.Transaction(ctx, func(tx *gorm.DB) error {
		var innerErr error
		// fetch current non-expired subscriptions for the org
		oldSubs, innerErr := l.repos.OrganizationSubscriptionRepo.FindActiveByOrgIdWithTx(ctx, tx, input.OrganizationId)
		if innerErr != nil {
			return innerErr
		}
		// unbind resources for currently effective old subscriptions
		removedPolicies, innerErr := unbindEffectiveSubscriptions(ctx, l.repos, tx, oldSubs)
		if innerErr != nil {
			return innerErr
		}
		// set all non-expired subscriptions to expired
		updateColumns := map[string]interface{}{
			"status": constants.OrganizationSubscriptionStatusExpired,
		}
		equalFilter := map[string]interface{}{
			"organization_id": input.OrganizationId,
		}
		notFilter := map[string]interface{}{
			"status": constants.OrganizationSubscriptionStatusExpired,
		}
		innerErr = l.repos.OrganizationSubscriptionRepo.UpdateColumnsWithTx(ctx, tx, updateColumns, equalFilter, notFilter)
		if innerErr != nil {
			return innerErr
		}
		// create new active subscription
		subscription, innerErr = l.repos.OrganizationSubscriptionRepo.CreateWithTx(ctx, tx, subscription, false)
		if innerErr != nil {
			return innerErr
		}
		// bind new plan resources to org roles
		addedPolicies, innerErr := syncSubscriptionEffect(ctx, l.repos, tx, subscriptionEffectSyncReq{
			OrganizationId: subscription.OrganizationId,
			Plan:           plan,
			OldEffective:   false,
			NewEffective:   true,
		})
		if innerErr != nil {
			return innerErr
		}
		organizationPolicies = make(map[int64]*organizationPolicyChanges)
		mergeOrganizationPolicies(organizationPolicies, removedPolicies)
		mergeOrganizationPolicies(organizationPolicies, addedPolicies)
		return nil
	})
	if err != nil {
		return nil, err
	}

	l.publishOrgPolicyReloadEvents(ctx, organizationPolicies, plan.Type)
	return l.transToDetailDto(ctx, subscription)
}

// Detail gets the details of an organization subscription.
func (l *OrganizationSubscriptionLogic) Detail(ctx context.Context, id int64) (*dto.OrganizationSubscriptionEntity, error) {
	subscription, err := l.repos.OrganizationSubscriptionRepo.GetOne(ctx, id)
	if err != nil {
		return nil, errorx.Wrap(err, "fail to find organization subscription")
	}

	return l.transToDetailDto(ctx, subscription)
}

// Renew renews an organization subscription.
func (l *OrganizationSubscriptionLogic) Renew(ctx context.Context, id int64) (result *dto.OrganizationSubscriptionEntity, err error) {
	subscription, err := l.repos.OrganizationSubscriptionRepo.GetOne(ctx, id)
	if err != nil {
		err = errorx.New("organization subscription not found")
		return
	}

	plan, err := l.repos.PricingPlanRepo.GetOne(ctx, subscription.PricingPlanId, []gormx.Preload{{Query: "Resources"}})
	if err != nil || plan == nil {
		err = errorx.New("fail to find pricing plan not found or disabled")
		return
	}

	if plan.Status != constants.CommonStatusActive {
		return nil, errorx.New("pricing plan is disabled")
	}

	policy, err := l.repos.PricingPlanPolicyRepo.GetOne(ctx, subscription.PricingPlanPolicyId)
	if err != nil {
		err = errorx.New("fail to find pricing plan")
		return
	}

	if policy.PricingPlanId != subscription.PricingPlanId {
		err = errorx.New("pricing plan policy not match")
		return
	}

	oldEffective := subscription.Status.IsEffective()
	// renew base: start from now if expired, else extend from current expire date
	base := time.Now()
	if subscription.ExpireDate != nil && subscription.ExpireDate.Valid && subscription.ExpireDate.Time.After(base) {
		base = subscription.ExpireDate.Time
	}
	subscription.ExpireDate = calculateExpireDate(policy.Type, base)
	subscription.Status = constants.OrganizationSubscriptionStatusActive
	newEffective := subscription.Status.IsEffective()

	var organizationPolicies map[int64]*organizationPolicyChanges
	err = l.repos.OrganizationSubscriptionRepo.Transaction(ctx, func(tx *gorm.DB) error {
		var innerErr error
		subscription, innerErr = l.repos.OrganizationSubscriptionRepo.UpdateOneByIdWithTx(ctx, tx, id, subscription)
		if innerErr != nil {
			return innerErr
		}
		organizationPolicies, innerErr = syncSubscriptionEffect(ctx, l.repos, tx, subscriptionEffectSyncReq{
			OrganizationId: subscription.OrganizationId,
			Plan:           plan,
			OldEffective:   oldEffective,
			NewEffective:   newEffective,
		})
		return innerErr
	})
	if err != nil {
		return nil, err
	}

	l.publishOrgPolicyReloadEvents(ctx, organizationPolicies, plan.Type)
	return l.transToDetailDto(ctx, subscription)
}

// UpdateExpireDate updates the expire date of an organization subscription.
func (l *OrganizationSubscriptionLogic) UpdateExpireDate(ctx context.Context, input dto.OrganizationSubscriptionUpdateExpireDateRequestDto) error {

	if input.ExpireDate == nil {
		return errorx.New("invalid expire date")
	}

	expireDate := input.ExpireDate.Time

	subscription, err := l.repos.OrganizationSubscriptionRepo.GetOne(ctx, input.Id)
	if err != nil {
		return errorx.Wrap(err, "organization subscription not found")
	}

	if subscription.ExpireDate != nil && subscription.ExpireDate.Valid && subscription.ExpireDate.Time.Equal(expireDate) {
		return nil
	}

	oldEffective := subscription.Status.IsEffective()
	subscription.ExpireDate = &sql.NullTime{Time: expireDate, Valid: true}
	// M2: a Suspended subscription stays Suspended when its expire date changes;
	// re-activating it requires an explicit SwitchStatus call. Only non-Suspended
	// statuses are derived from the new expire date.
	if subscription.Status != constants.OrganizationSubscriptionStatusSuspended {
		if expireDate.After(time.Now()) {
			subscription.Status = constants.OrganizationSubscriptionStatusActive
		} else {
			subscription.Status = constants.OrganizationSubscriptionStatusExpired
		}
	}
	newEffective := subscription.Status.IsEffective()

	plan, err := l.repos.PricingPlanRepo.GetOne(ctx, subscription.PricingPlanId, []gormx.Preload{{Query: "Resources"}})
	if err != nil {
		return errorx.Wrap(err, "fail to find pricing plan")
	}

	var organizationPolicies map[int64]*organizationPolicyChanges
	err = l.repos.OrganizationSubscriptionRepo.Transaction(ctx, func(tx *gorm.DB) error {
		var innerErr error
		_, innerErr = l.repos.OrganizationSubscriptionRepo.UpdateOneByIdWithTx(ctx, tx, input.Id, subscription)
		if innerErr != nil {
			return innerErr
		}
		organizationPolicies, innerErr = syncSubscriptionEffect(ctx, l.repos, tx, subscriptionEffectSyncReq{
			OrganizationId: subscription.OrganizationId,
			Plan:           plan,
			OldEffective:   oldEffective,
			NewEffective:   newEffective,
		})
		return innerErr
	})
	if err != nil {
		return err
	}

	l.publishOrgPolicyReloadEvents(ctx, organizationPolicies, plan.Type)
	return nil
}

// SwitchStatus switches the status of an organization subscription.
func (l *OrganizationSubscriptionLogic) SwitchStatus(ctx context.Context, id int64, status constants.OrganizationSubscriptionStatus) error {
	subscription, err := l.repos.OrganizationSubscriptionRepo.GetOne(ctx, id)
	if err != nil {
		return errorx.Wrap(err, "organization subscription not found")
	}

	if subscription.Status == status {
		return nil
	}

	oldEffective := subscription.Status.IsEffective()
	subscription.Status = status
	newEffective := subscription.Status.IsEffective()

	plan, err := l.repos.PricingPlanRepo.GetOne(ctx, subscription.PricingPlanId, []gormx.Preload{{Query: "Resources"}})
	if err != nil {
		return errorx.Wrap(err, "fail to find pricing plan")
	}

	var organizationPolicies map[int64]*organizationPolicyChanges
	err = l.repos.OrganizationSubscriptionRepo.Transaction(ctx, func(tx *gorm.DB) error {
		var innerErr error
		_, innerErr = l.repos.OrganizationSubscriptionRepo.UpdateOneByIdWithTx(ctx, tx, id, subscription)
		if innerErr != nil {
			return innerErr
		}
		organizationPolicies, innerErr = syncSubscriptionEffect(ctx, l.repos, tx, subscriptionEffectSyncReq{
			OrganizationId: subscription.OrganizationId,
			Plan:           plan,
			OldEffective:   oldEffective,
			NewEffective:   newEffective,
		})
		return innerErr
	})
	if err != nil {
		return err
	}

	l.publishOrgPolicyReloadEvents(ctx, organizationPolicies, plan.Type)
	return nil
}

// SuspendEffectiveByOrganizationId suspends every effective subscription of an
// organization in one transaction and unbinds its plan resources (remove-side).
// Used when an organization is disabled so its plan-granted access is revoked.
// Subscriptions are NOT deleted — the system has no delete-subscription
// operation; re-enabling the organization requires reactivating subscriptions.
func (l *OrganizationSubscriptionLogic) SuspendEffectiveByOrganizationId(ctx context.Context, organizationId int64, userType constants.UserType) error {
	organizationPolicies := make(map[int64]*organizationPolicyChanges)
	err := l.repos.OrganizationSubscriptionRepo.Transaction(ctx, func(tx *gorm.DB) error {
		subs, innerErr := l.repos.OrganizationSubscriptionRepo.FindActiveByOrgIdWithTx(ctx, tx, organizationId)
		if innerErr != nil {
			return innerErr
		}
		for _, sub := range subs {
			if sub == nil || !sub.Status.IsEffective() {
				continue
			}
			plan, planErr := l.repos.PricingPlanRepo.GetOne(ctx, sub.PricingPlanId, []gormx.Preload{{Query: "Resources"}})
			if planErr != nil || plan == nil {
				continue
			}
			// remove-side: unbind this subscription's plan resources from the org's roles
			removed, unbindErr := syncPlanResourcesToOrganizationRoles(ctx, l.repos, tx, planResourceSyncReq{
				OrganizationIds: []int64{sub.OrganizationId},
				ResourceIds:     planResourceIds(plan),
				UserType:        plan.Type,
				IsBind:          false,
			})
			if unbindErr != nil {
				return unbindErr
			}
			mergeOrganizationPolicies(organizationPolicies, removed)
			sub.Status = constants.OrganizationSubscriptionStatusSuspended
			if _, updateErr := l.repos.OrganizationSubscriptionRepo.UpdateOneByIdWithTx(ctx, tx, sub.Id, sub); updateErr != nil {
				return updateErr
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	l.publishOrgPolicyReloadEvents(ctx, organizationPolicies, userType)
	return nil
}

// Page returns a page of organization subscriptions.
func (l *OrganizationSubscriptionLogic) Page(ctx context.Context, req dto.OrganizationSubscriptionPageRequestDto) (gormx.PageResponse[dto.OrganizationSubscriptionPageEntity], error) {
	pageResult, err := l.repos.OrganizationSubscriptionRepo.Page(ctx, req)
	if err != nil {
		return gormx.PageResponse[dto.OrganizationSubscriptionPageEntity]{}, err
	}

	resultPtr, err := copierx.Trans[gormx.PageResponse[repo.OrganizationSubscriptionPageRow], gormx.PageResponse[dto.OrganizationSubscriptionPageEntity]](&pageResult)
	if err != nil {
		return gormx.PageResponse[dto.OrganizationSubscriptionPageEntity]{}, err
	}

	return *resultPtr, nil
}

func calculateExpireDate(policyType constants.PricingPolicyType, base time.Time) *sql.NullTime {
	months := 1
	switch policyType {
	case constants.PricingPolicyTypeQuarterly:
		months = 3
	case constants.PricingPolicyTypeYearly:
		months = 12
	default:
		months = 1
	}
	// Billing term is N full calendar months inclusive of the start day:
	// expire = (start - 1 day) + N months, day clamped to the target month's
	// last day. e.g. monthly Jan 1 -> Jan 31, Jan 31 -> Feb 28/29.
	expire := addMonthsClamped(base.AddDate(0, 0, -1), months)
	return &sql.NullTime{Time: expire, Valid: true}
}

// addMonthsClamped adds months to t, clamping the day to the target month's
// last day (Jan 31 + 1 month -> Feb 28/29). Unlike time.AddDate it never
// overflows the day into the following month.
func addMonthsClamped(t time.Time, months int) time.Time {
	year, month, day := t.Date()
	totalMonths := int(month) - 1 + months
	year += totalMonths / 12
	month = time.Month(totalMonths%12 + 1)
	if last := lastDayOfMonth(year, month, t.Location()); day > last {
		day = last
	}
	return time.Date(year, month, day, t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location())
}

// lastDayOfMonth returns the last calendar day of the given year/month in loc.
func lastDayOfMonth(year int, month time.Month, loc *time.Location) int {
	return time.Date(year, month+1, 0, 0, 0, 0, 0, loc).Day()
}

func (l *OrganizationSubscriptionLogic) transToDetailDto(ctx context.Context, sub *model.OrganizationSubscription) (*dto.OrganizationSubscriptionEntity, error) {
	result, err := copierx.Trans[model.OrganizationSubscription, dto.OrganizationSubscriptionEntity](sub)
	if err != nil {
		return nil, err
	}

	org, innerErr := l.repos.OrganizationRepo.GetOrganizationById(ctx, sub.OrganizationId, false)
	if innerErr == nil && org != nil {
		orgDto, err := copierx.Trans[model.Organization, dto.OrganizationEntityDto](org)
		if err != nil {
			return nil, err
		}
		result.Organization = orgDto
	}

	plan, innerErr := l.repos.PricingPlanRepo.GetOne(ctx, sub.PricingPlanId, nil)
	if innerErr == nil && plan != nil {
		planDto, err := copierx.Trans[model.PricingPlan, dto.PricingPlanEntityDto](plan)
		if err != nil {
			return nil, err
		}
		result.PricingPlan = planDto
	}

	return result, nil
}

// SubscribeDefaultPlan finds the default plan and creates a trial subscription.
// Only applies to plans with TrialDays > 0 and Active status.
// Creates subscription with Trialing status and expire_date based on plan.TrialDays.
func (l *OrganizationSubscriptionLogic) SubscribeDefaultPlan(ctx context.Context, organizationId int64, userType constants.UserType) error {
	plan, err := l.repos.PricingPlanRepo.FindDefaultByType(ctx, userType)
	if err != nil || plan == nil {
		return nil
	}

	if plan.Status != constants.CommonStatusActive || plan.TrialDays <= 0 {
		return nil
	}

	plan, err = l.repos.PricingPlanRepo.GetOne(ctx, plan.Id, []gormx.Preload{{Query: "Resources"}})
	if err != nil || plan == nil {
		return nil
	}

	code, err := l.generateCode(ctx)
	if err != nil {
		return err
	}
	expireDate := &sql.NullTime{Time: time.Now().AddDate(0, 0, int(plan.TrialDays)), Valid: true}

	subscription := &model.OrganizationSubscription{
		Code:                code,
		OrganizationId:      organizationId,
		PricingPlanId:       plan.Id,
		PricingPlanPolicyId: 0,
		Status:              constants.OrganizationSubscriptionStatusTrialing,
		ExpireDate:          expireDate,
		OriginalAmount:      0,
		ActualAmount:        0,
	}

	err = l.repos.OrganizationSubscriptionRepo.Transaction(ctx, func(tx *gorm.DB) error {
		activeSubs, innerErr := l.repos.OrganizationSubscriptionRepo.FindActiveByOrgIdWithTx(ctx, tx, organizationId)
		if innerErr != nil {
			return innerErr
		}
		for _, sub := range activeSubs {
			sub.Status = constants.OrganizationSubscriptionStatusExpired
			_, innerErr = l.repos.OrganizationSubscriptionRepo.UpdateOneByIdWithTx(ctx, tx, sub.Id, sub)
			if innerErr != nil {
				return innerErr
			}
		}

		_, createErr := l.repos.OrganizationSubscriptionRepo.CreateWithTx(ctx, tx, subscription, false)
		return createErr
	})
	if err != nil {
		return err
	}

	// M5: trial plan resources are not auto-bound to roles — the org admin
	// configures role-resource bindings manually. No auth event is needed: the
	// enforcer plan-resource gate is recomputed from the DB on the next
	// RebuildPolicy.
	return nil
}

// maxCodeRetries is the number of attempts generateCode makes to produce a
// unique subscription code before giving up.
const maxCodeRetries = 5

// generateCode produces a unique subscription code. It retries on collision with
// an existing code (checked via ExistsByCode) up to maxCodeRetries times.
func (l *OrganizationSubscriptionLogic) generateCode(ctx context.Context) (string, error) {
	for i := 0; i < maxCodeRetries; i++ {
		randN := rand2.N(99999)
		randStr := fmt.Sprintf("%05d", randN)
		code := fmt.Sprintf("SU%s%s", time.Now().Format("20060102150405"), randStr)
		exists, err := l.repos.OrganizationSubscriptionRepo.ExistsByCode(ctx, code)
		if err != nil {
			return "", err
		}
		if !exists {
			return code, nil
		}
	}
	return "", errorx.New("fail to generate unique subscription code")
}

// planResourceIds extracts resource IDs from a pricing plan.
func planResourceIds(plan *model.PricingPlan) []int64 {
	if plan == nil {
		return nil
	}
	return slice.Map(plan.Resources, func(index int, resource *model.Resource) int64 {
		return resource.Id
	})
}

// subscriptionEffectSyncReq bundles inputs for syncing a single organization's
// plan resources when a subscription crosses the effective boundary.
type subscriptionEffectSyncReq struct {
	OrganizationId int64
	Plan           *model.PricingPlan
	OldEffective   bool
	NewEffective   bool
}

// syncSubscriptionEffect binds/unbinds a single organization's plan resources
// when a subscription crosses the effective boundary (OldEffective != NewEffective).
// No-op when effective status is unchanged or the plan has no resources.
// Returns per-org policy changes for downstream auth reload events.
// syncSubscriptionEffect propagates a single organization's subscription
// effective-boundary crossing to role-resource bindings.
//
// M5: becoming effective (NewEffective=true) does NOT auto-bind plan resources
// to roles — the org admin configures role-resource bindings manually. The
// subscription only defines the available plan-resource scope; the enforcer
// gate (organizationPlanResourceMap) is recomputed from the DB on the next
// RebuildPolicy, so no event is needed on the bind side.
// Becoming ineffective (NewEffective=false) unbinds the plan's resources
// (remove-side) and collects the removed policies for immediate revocation.
func syncSubscriptionEffect(ctx context.Context, repos *repo.AllRepos, tx *gorm.DB, req subscriptionEffectSyncReq) (map[int64]*organizationPolicyChanges, error) {
	if req.OldEffective == req.NewEffective {
		return nil, nil
	}
	if req.NewEffective {
		return nil, nil
	}
	resourceIds := planResourceIds(req.Plan)
	if len(resourceIds) == 0 {
		return nil, nil
	}
	return syncPlanResourcesToOrganizationRoles(ctx, repos, tx, planResourceSyncReq{
		OrganizationIds: []int64{req.OrganizationId},
		ResourceIds:     resourceIds,
		UserType:        req.Plan.Type,
		IsBind:          false,
	})
}

// unbindEffectiveSubscriptions unbinds plan resources for subscriptions that
// are currently effective (Active/Trialing) and groups their policy changes.
// Subscriptions sharing a plan are unbound once. Non-effective statuses are skipped.
// orgIdByUserType maps each org's affected policies, keyed by organization id.
func unbindEffectiveSubscriptions(ctx context.Context, repos *repo.AllRepos, tx *gorm.DB, subs []*model.OrganizationSubscription) (map[int64]*organizationPolicyChanges, error) {
	merged := make(map[int64]*organizationPolicyChanges)
	// group effective subs by plan id to avoid reloading the same plan
	planOrgs := make(map[int64][]int64)
	for _, sub := range subs {
		if sub == nil || !sub.Status.IsEffective() {
			continue
		}
		planOrgs[sub.PricingPlanId] = append(planOrgs[sub.PricingPlanId], sub.OrganizationId)
	}
	for planId, orgIds := range planOrgs {
		plan, innerErr := repos.PricingPlanRepo.GetOne(ctx, planId, []gormx.Preload{{Query: "Resources"}})
		if innerErr != nil || plan == nil {
			continue
		}
		resourceIds := planResourceIds(plan)
		if len(resourceIds) == 0 {
			continue
		}
		orgPolicies, innerErr := syncPlanResourcesToOrganizationRoles(ctx, repos, tx, planResourceSyncReq{
			OrganizationIds: slice.Unique(orgIds),
			ResourceIds:     resourceIds,
			UserType:        plan.Type,
			IsBind:          false,
		})
		if innerErr != nil {
			return nil, innerErr
		}
		mergeOrganizationPolicies(merged, orgPolicies)
	}
	return merged, nil
}

// publishOrgPolicyReloadEvents sends per-org RoleReload events with explicit policy changes.
func (l *OrganizationSubscriptionLogic) publishOrgPolicyReloadEvents(ctx context.Context, organizationPolicies map[int64]*organizationPolicyChanges, userType constants.UserType) {
	for orgId, changes := range organizationPolicies {
		orgId := orgId
		changes := changes
		threading.GoSafe(func() {
			subCtx := context.WithoutCancel(ctx)
			msg := auth.Msg{
				Action:          auth.ActionRoleReload,
				Scope:           auth.ScopeOrganization,
				Type:            userType,
				OrganizationId:  orgId,
				AddedPolicies:   changes.added,
				RemovedPolicies: changes.removed,
			}
			_ = l.publishAuthEvent(subCtx, l.authTopicMap[userType], msg)
		})
	}
}
