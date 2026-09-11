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

	"github.com/zeromicro/go-zero/core/logx"
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

func validateEffectivePlan(organization *model.Organization, plan *model.PricingPlan) error {
	if plan == nil || plan.Status != constants.CommonStatusActive {
		return errorx.New("pricing plan is disabled")
	}
	if organization == nil || organization.Type != plan.Type {
		return errorx.New("organization type not match pricing plan type")
	}
	return nil
}

// Create creates a new organization subscription.
func (l *OrganizationSubscriptionLogic) Create(ctx context.Context, input dto.OrganizationSubscriptionCreateDto) (result *dto.OrganizationSubscriptionEntity, err error) {
	plan, err := l.repos.PricingPlanRepo.GetOne(ctx, input.PricingPlanId, nil)
	if err != nil {
		return nil, errorx.Wrap(err, "fail to find pricing plan")
	}
	if plan.Type != input.Type {
		return nil, errorx.New("pricing plan type not match")
	}
	if plan.Status != constants.CommonStatusActive {
		return nil, errorx.New("pricing plan not found or disabled")
	}

	policy, err := l.repos.PricingPlanPolicyRepo.GetOne(ctx, input.PricingPlanPolicyId)
	if err != nil {
		return nil, errorx.New("fail to find pricing plan policy")
	}

	if policy.PricingPlanId != input.PricingPlanId {
		err = errorx.New("pricing plan policy not match")
		return
	}

	code, err := l.generateCode(ctx, nil)
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

	err = l.repos.OrganizationSubscriptionRepo.Transaction(ctx, func(tx *gorm.DB) error {
		var innerErr error
		lockedOrganization, innerErr := l.repos.OrganizationRepo.LockByIdWithTx(ctx, tx, input.OrganizationId)
		if innerErr != nil {
			return innerErr
		}
		plan, innerErr = l.repos.PricingPlanRepo.GetOneWithTx(ctx, tx, input.PricingPlanId, []gormx.Preload{{Query: "Resources"}})
		if innerErr != nil {
			return errorx.Wrap(innerErr, "fail to find pricing plan")
		}
		if plan.Type != input.Type {
			return errorx.New("pricing plan type not match")
		}
		if innerErr = validateEffectivePlan(lockedOrganization, plan); innerErr != nil {
			return innerErr
		}
		// fetch current non-expired subscriptions for the org
		oldSubs, innerErr := l.repos.OrganizationSubscriptionRepo.FindActiveByOrgIdWithTx(ctx, tx, input.OrganizationId)
		if innerErr != nil {
			return innerErr
		}
		for _, oldSub := range oldSubs {
			if oldSub.PricingPlanId == plan.Id {
				continue
			}
			oldPlan, findErr := l.repos.PricingPlanRepo.GetOneWithTx(ctx, tx, oldSub.PricingPlanId, []gormx.Preload{{Query: "Resources"}})
			if findErr != nil {
				return findErr
			}
			if innerErr = reconcileRoleResourcesForPlanChange(ctx, l.repos, tx, oldSub.OrganizationId, oldPlan, plan); innerErr != nil {
				return innerErr
			}
		}
		// set non-expired subscriptions to expired; Suspended subs preserved.
		updateColumns := map[string]interface{}{
			"status": constants.OrganizationSubscriptionStatusExpired,
		}
		equalFilter := map[string]interface{}{
			"organization_id": input.OrganizationId,
		}
		notFilter := map[string]interface{}{
			"status": []constants.OrganizationSubscriptionStatus{constants.OrganizationSubscriptionStatusExpired, constants.OrganizationSubscriptionStatusSuspended},
		}
		innerErr = l.repos.OrganizationSubscriptionRepo.UpdateColumnsWithTx(ctx, tx, updateColumns, equalFilter, notFilter)
		if innerErr != nil {
			return innerErr
		}
		// create new active subscription. M5: becoming effective does not auto-bind
		// plan resources — the org admin configures role-resource bindings manually.
		subscription, innerErr = l.repos.OrganizationSubscriptionRepo.CreateWithTx(ctx, tx, subscription, false)
		if innerErr != nil {
			return innerErr
		}
		if len(oldSubs) == 0 {
			return reconcileRoleResourcesForPlanChange(ctx, l.repos, tx, subscription.OrganizationId, nil, plan)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	// rebuild this org's policies so the plan gate recomputes its resource scope.
	l.publishOrgPolicyRebuild(ctx, []int64{subscription.OrganizationId}, plan.Type)
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

	var plan *model.PricingPlan
	err = l.repos.OrganizationSubscriptionRepo.Transaction(ctx, func(tx *gorm.DB) error {
		var innerErr error
		lockedOrganization, innerErr := l.repos.OrganizationRepo.LockByIdWithTx(ctx, tx, subscription.OrganizationId)
		if innerErr != nil {
			return innerErr
		}
		plan, innerErr = l.repos.PricingPlanRepo.GetOneWithTx(ctx, tx, subscription.PricingPlanId, []gormx.Preload{{Query: "Resources"}})
		if innerErr != nil {
			return errorx.Wrap(innerErr, "fail to find pricing plan")
		}
		if newEffective {
			if innerErr = validateEffectivePlan(lockedOrganization, plan); innerErr != nil {
				return innerErr
			}
		}
		if newEffective {
			if innerErr = l.expireOtherEffectiveSubscriptions(ctx, tx, subscription.OrganizationId, id); innerErr != nil {
				return innerErr
			}
		}
		subscription, innerErr = l.repos.OrganizationSubscriptionRepo.UpdateOneByIdWithTx(ctx, tx, id, subscription)
		if innerErr != nil {
			return innerErr
		}
		// M5: only becoming ineffective unbinds (remove-side); the rebuild recomputes the gate.
		return syncSubscriptionEffect(ctx, l.repos, tx, subscription.OrganizationId, plan, oldEffective, newEffective)
	})
	if err != nil {
		return nil, err
	}

	l.publishOrgPolicyRebuild(ctx, []int64{subscription.OrganizationId}, plan.Type)
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

	var plan *model.PricingPlan
	err = l.repos.OrganizationSubscriptionRepo.Transaction(ctx, func(tx *gorm.DB) error {
		var innerErr error
		lockedOrganization, innerErr := l.repos.OrganizationRepo.LockByIdWithTx(ctx, tx, subscription.OrganizationId)
		if innerErr != nil {
			return innerErr
		}
		plan, innerErr = l.repos.PricingPlanRepo.GetOneWithTx(ctx, tx, subscription.PricingPlanId, []gormx.Preload{{Query: "Resources"}})
		if innerErr != nil {
			return errorx.Wrap(innerErr, "fail to find pricing plan")
		}
		if newEffective {
			if innerErr = validateEffectivePlan(lockedOrganization, plan); innerErr != nil {
				return innerErr
			}
		}
		if newEffective {
			if innerErr = l.expireOtherEffectiveSubscriptions(ctx, tx, subscription.OrganizationId, input.Id); innerErr != nil {
				return innerErr
			}
		}
		_, innerErr = l.repos.OrganizationSubscriptionRepo.UpdateOneByIdWithTx(ctx, tx, input.Id, subscription)
		if innerErr != nil {
			return innerErr
		}
		// M5: only becoming ineffective unbinds (remove-side); the rebuild recomputes the gate.
		return syncSubscriptionEffect(ctx, l.repos, tx, subscription.OrganizationId, plan, oldEffective, newEffective)
	})
	if err != nil {
		return err
	}

	l.publishOrgPolicyRebuild(ctx, []int64{subscription.OrganizationId}, plan.Type)
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
	if status == constants.OrganizationSubscriptionStatusExpired {
		subscription.ExpireDate = &sql.NullTime{Time: time.Now(), Valid: true}
	}
	newEffective := subscription.Status.IsEffective()

	var plan *model.PricingPlan
	err = l.repos.OrganizationSubscriptionRepo.Transaction(ctx, func(tx *gorm.DB) error {
		var innerErr error
		lockedOrganization, innerErr := l.repos.OrganizationRepo.LockByIdWithTx(ctx, tx, subscription.OrganizationId)
		if innerErr != nil {
			return innerErr
		}
		plan, innerErr = l.repos.PricingPlanRepo.GetOneWithTx(ctx, tx, subscription.PricingPlanId, []gormx.Preload{{Query: "Resources"}})
		if innerErr != nil {
			return errorx.Wrap(innerErr, "fail to find pricing plan")
		}
		if newEffective {
			if innerErr = validateEffectivePlan(lockedOrganization, plan); innerErr != nil {
				return innerErr
			}
		}
		if newEffective {
			if innerErr = l.expireOtherEffectiveSubscriptions(ctx, tx, subscription.OrganizationId, id); innerErr != nil {
				return innerErr
			}
		}
		_, innerErr = l.repos.OrganizationSubscriptionRepo.UpdateOneByIdWithTx(ctx, tx, id, subscription)
		if innerErr != nil {
			return innerErr
		}
		// M5: only becoming ineffective unbinds (remove-side); the rebuild recomputes the gate.
		return syncSubscriptionEffect(ctx, l.repos, tx, subscription.OrganizationId, plan, oldEffective, newEffective)
	})
	if err != nil {
		return err
	}

	l.publishOrgPolicyRebuild(ctx, []int64{subscription.OrganizationId}, plan.Type)
	return nil
}

// SuspendEffectiveByOrganizationId suspends every effective subscription of an
// org and unbinds its plan resources (remove-side). Used when an org is disabled.
// Subscriptions are suspended, not deleted (no delete op exists).
func (l *OrganizationSubscriptionLogic) SuspendEffectiveByOrganizationId(ctx context.Context, organizationId int64, userType constants.UserType) error {
	err := l.repos.OrganizationSubscriptionRepo.Transaction(ctx, func(tx *gorm.DB) error {
		if _, innerErr := l.repos.OrganizationRepo.LockByIdWithTx(ctx, tx, organizationId); innerErr != nil {
			return innerErr
		}
		subs, innerErr := l.repos.OrganizationSubscriptionRepo.FindActiveByOrgIdWithTx(ctx, tx, organizationId)
		if innerErr != nil {
			return innerErr
		}
		for _, sub := range subs {
			if sub == nil || !sub.Status.IsEffective() {
				continue
			}
			plan, planErr := l.repos.PricingPlanRepo.GetOneWithTx(ctx, tx, sub.PricingPlanId, []gormx.Preload{{Query: "Resources"}})
			if planErr != nil {
				return errorx.Wrap(planErr, "fail to find pricing plan")
			}
			if reconcileErr := reconcileRoleResourcesForPlanChange(ctx, l.repos, tx, sub.OrganizationId, plan, nil); reconcileErr != nil {
				return reconcileErr
			}
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

	l.publishOrgPolicyRebuild(ctx, []int64{organizationId}, userType)
	return nil
}

func (l *OrganizationSubscriptionLogic) expireOtherEffectiveSubscriptions(ctx context.Context, tx *gorm.DB, organizationId, keepId int64) error {
	subs, err := l.repos.OrganizationSubscriptionRepo.FindActiveByOrgIdWithTx(ctx, tx, organizationId)
	if err != nil {
		return err
	}
	for _, sub := range subs {
		if sub == nil || sub.Id == keepId || !sub.Status.IsEffective() {
			continue
		}
		plan, err := l.repos.PricingPlanRepo.GetOneWithTx(ctx, tx, sub.PricingPlanId, []gormx.Preload{{Query: "Resources"}})
		if err != nil {
			return errorx.Wrap(err, "fail to find pricing plan")
		}
		if err = reconcileRoleResourcesForPlanChange(ctx, l.repos, tx, organizationId, plan, nil); err != nil {
			return err
		}
		if err = l.repos.OrganizationSubscriptionRepo.UpdateColumnsByIdWithTx(ctx, tx, sub.Id, map[string]interface{}{
			"status": constants.OrganizationSubscriptionStatusExpired,
		}); err != nil {
			return err
		}
	}
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

	code, err := l.generateCode(ctx, nil)
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
		organization, innerErr := l.repos.OrganizationRepo.LockByIdWithTx(ctx, tx, organizationId)
		if innerErr != nil {
			return innerErr
		}
		if organization.Type != userType {
			return errorx.New("organization type not match pricing plan type")
		}
		if innerErr = l.expireOtherEffectiveSubscriptions(ctx, tx, organizationId, 0); innerErr != nil {
			return innerErr
		}
		if innerErr = reconcileRoleResourcesForPlanChange(ctx, l.repos, tx, organizationId, nil, plan); innerErr != nil {
			return innerErr
		}

		_, createErr := l.repos.OrganizationSubscriptionRepo.CreateWithTx(ctx, tx, subscription, false)
		return createErr
	})
	if err != nil {
		return err
	}

	// recompute this org's policies: the trial expanded its plan-resource scope.
	l.publishOrgPolicyRebuild(ctx, []int64{organizationId}, plan.Type)

	return nil
}

func (l *OrganizationSubscriptionLogic) UpdateStatus(ctx context.Context, organizationType constants.UserType) (err error) {
	// default plan + its free (ActualAmount == 0) policy, used to auto-recover
	// expired subscriptions so an org never fully loses access on expiry.
	planFilter := []gormx.WhereRaw{
		{
			Raw: "is_default = ? and type = ? and status = ?",
			Args: []any{
				true,
				organizationType,
				constants.CommonStatusActive,
			},
		},
	}
	planPreloads := []gormx.Preload{
		{Query: "Policies"},
		{Query: "Resources"},
	}
	recovery := expiredSubscriptionRecovery{Now: time.Now()}
	recovery.DefaultPlan, _ = l.repos.PricingPlanRepo.FindOneByFilter(ctx, planFilter, planPreloads...)
	if recovery.DefaultPlan != nil {
		recovery.DefaultPolicy, _ = slice.FindBy(recovery.DefaultPlan.Policies, func(index int, item *model.PricingPlanPolicy) bool {
			return item.ActualAmount == 0
		})
	}

	// subscriptions still in an effective status for this user-type
	activeSubs, err := l.repos.OrganizationSubscriptionRepo.FindActiveSubscriptions(ctx, organizationType)
	if err != nil {
		return
	}

	now := recovery.Now
	expiringDate := now.AddDate(0, 0, 7)
	for _, sub := range activeSubs {
		if sub == nil {
			continue
		}
		newStatus := sub.Status
		expireDate := now
		if sub.ExpireDate != nil && sub.ExpireDate.Valid {
			expireDate = sub.ExpireDate.Time
		}
		switch sub.Status {
		case constants.OrganizationSubscriptionStatusActive:
			if now.After(expireDate) {
				newStatus = constants.OrganizationSubscriptionStatusExpired
			} else if expireDate.Before(expiringDate) {
				newStatus = constants.OrganizationSubscriptionStatusExpiring
			}
		case constants.OrganizationSubscriptionStatusExpiring:
			if now.After(expireDate) {
				newStatus = constants.OrganizationSubscriptionStatusExpired
			}
		case constants.OrganizationSubscriptionStatusTrialing:
			if now.After(expireDate) {
				newStatus = constants.OrganizationSubscriptionStatusExpired
			}
		}

		if newStatus == sub.Status {
			continue
		}
		// non-Expired transition (e.g. Active -> Expiring): simple status update
		if newStatus != constants.OrganizationSubscriptionStatusExpired {
			_ = l.repos.OrganizationSubscriptionRepo.UpdateColumnsById(ctx, sub.Id, map[string]any{
				"status": newStatus,
			})
			continue
		}

		// Expired transition: recover via default plan when possible
		txErr := l.repos.OrganizationSubscriptionRepo.Transaction(ctx, func(tx *gorm.DB) error {
			_, innerErr := l.handleExpiredSubscription(ctx, tx, sub, recovery)
			return innerErr
		})
		if txErr != nil {
			logx.WithContext(ctx).WithFields(
				logx.Field("type", "update_organization_subscription_status"),
				logx.Field("subscription_id", sub.Id),
				logx.Field("organization_id", sub.OrganizationId),
				logx.Field("organization_type", organizationType),
			).
				Errorf("fail to handle expired subscription %d: %v", sub.Id, txErr)
		} else {
			l.publishOrgPolicyRebuild(ctx, []int64{sub.OrganizationId}, organizationType)
		}
	}

	return
}

// expiredSubscriptionRecovery bundles the default-plan recovery context shared
// across every expired subscription in a single UpdateStatus pass.
type expiredSubscriptionRecovery struct {
	DefaultPlan   *model.PricingPlan
	DefaultPolicy *model.PricingPlanPolicy
	Now           time.Time
}

// handleExpiredSubscription marks sub Expired, then recovers access via the
// default plan + free policy if both exist: same plan renews in place; different
// plan creates a new sub and unbinds only old-plan-exclusive resources. Returns
// per-org removed policies for the caller to publish after commit.
func (l *OrganizationSubscriptionLogic) handleExpiredSubscription(
	ctx context.Context,
	tx *gorm.DB,
	sub *model.OrganizationSubscription,
	recovery expiredSubscriptionRecovery,
) (map[int64]*organizationPolicyChanges, error) {
	var err error
	if _, err = l.repos.OrganizationRepo.LockByIdWithTx(ctx, tx, sub.OrganizationId); err != nil {
		return nil, err
	}
	if err = l.repos.OrganizationSubscriptionRepo.UpdateColumnsByIdWithTx(ctx, tx, sub.Id, map[string]interface{}{
		"status": constants.OrganizationSubscriptionStatusExpired,
	}); err != nil {
		return nil, err
	}

	if recovery.DefaultPlan == nil || recovery.DefaultPolicy == nil {
		oldPlan, planErr := l.repos.PricingPlanRepo.GetOneWithTx(ctx, tx, sub.PricingPlanId, []gormx.Preload{{Query: "Resources"}})
		if planErr != nil {
			return nil, errorx.Wrap(planErr, "fail to find expired subscription pricing plan")
		}
		return nil, reconcileRoleResourcesForPlanChange(ctx, l.repos, tx, sub.OrganizationId, oldPlan, nil)
	}

	// same plan -> renew in place
	if sub.PricingPlanId == recovery.DefaultPlan.Id {
		if err = l.repos.OrganizationSubscriptionRepo.UpdateColumnsByIdWithTx(ctx, tx, sub.Id, map[string]interface{}{
			"status":                 constants.OrganizationSubscriptionStatusActive,
			"expire_date":            calculateExpireDate(recovery.DefaultPolicy.Type, recovery.Now),
			"pricing_plan_policy_id": recovery.DefaultPolicy.Id,
			"original_amount":        recovery.DefaultPolicy.OriginalAmount,
			"actual_amount":          recovery.DefaultPolicy.ActualAmount,
		}); err != nil {
			return nil, err
		}
		return nil, nil
	}

	// different plan -> create new subscription under the default plan
	oldPlan, err := l.repos.PricingPlanRepo.GetOneWithTx(ctx, tx, sub.PricingPlanId, []gormx.Preload{{Query: "Resources"}})
	if err != nil {
		return nil, errorx.Wrap(err, "fail to find expired subscription pricing plan")
	}

	code, err := l.generateCode(ctx, tx)
	if err != nil {
		return nil, err
	}
	newSub := &model.OrganizationSubscription{
		Code:                code,
		OrganizationId:      sub.OrganizationId,
		PricingPlanId:       recovery.DefaultPlan.Id,
		PricingPlanPolicyId: recovery.DefaultPolicy.Id,
		Status:              constants.OrganizationSubscriptionStatusActive,
		ExpireDate:          calculateExpireDate(recovery.DefaultPolicy.Type, recovery.Now),
		OriginalAmount:      recovery.DefaultPolicy.OriginalAmount,
		ActualAmount:        recovery.DefaultPolicy.ActualAmount,
	}
	if _, err = l.repos.OrganizationSubscriptionRepo.CreateWithTx(ctx, tx, newSub, false); err != nil {
		return nil, err
	}

	return nil, reconcileRoleResourcesForPlanChange(ctx, l.repos, tx, sub.OrganizationId, oldPlan, recovery.DefaultPlan)
}

// publishOrgPolicyRebuild emits a per-org ActionPolicyRebuild for each org whose
// subscription state changed; each receiving replica recomputes that org's
// policies from post-commit DB state.
func (l *OrganizationSubscriptionLogic) publishOrgPolicyRebuild(ctx context.Context, affectedOrgIds []int64, userType constants.UserType) {
	seen := map[int64]bool{}
	for _, orgId := range affectedOrgIds {
		if orgId <= 0 || seen[orgId] {
			continue
		}
		seen[orgId] = true
		orgId := orgId
		threading.GoSafe(func() {
			subCtx := context.WithoutCancel(ctx)
			msg := auth.Msg{
				Action:         auth.ActionPolicyRebuild,
				Scope:          auth.ScopeOrganization,
				Type:           userType,
				OrganizationId: orgId,
			}
			_ = l.publishAuthEvent(subCtx, l.authTopicMap[userType], msg)
		})
	}
}

// maxCodeRetries is the number of attempts generateCode makes to produce a
// unique subscription code before giving up.
const maxCodeRetries = 5

// generateCode produces a unique subscription code. It retries on collision with
// an existing code (checked via ExistsByCode) up to maxCodeRetries times.
func (l *OrganizationSubscriptionLogic) generateCode(ctx context.Context, tx *gorm.DB) (string, error) {
	for i := 0; i < maxCodeRetries; i++ {
		randN := rand2.N(99999)
		randStr := fmt.Sprintf("%05d", randN)
		code := fmt.Sprintf("SU%s%s", time.Now().Format("20060102150405"), randStr)
		var exists bool
		var err error
		if tx == nil {
			exists, err = l.repos.OrganizationSubscriptionRepo.IsExistedByColumns(ctx, map[string]interface{}{
				"code": code,
			})
		} else {
			exists, err = l.repos.OrganizationSubscriptionRepo.IsExistedByColumnsWithTx(ctx, tx, map[string]interface{}{
				"code": code,
			})
		}
		// exists, err := l.repos.OrganizationSubscriptionRepo.ExistsByCode(ctx, code)
		if err != nil {
			return "", err
		}
		if !exists {
			return code, nil
		}
	}
	return "", errorx.New("fail to generate unique subscription code")
}

// syncSubscriptionEffect propagates a subscription's effective-boundary crossing
// to role-resource bindings. M5: only becoming ineffective unbinds (remove-side);
// becoming effective never auto-binds (admin configures, gate recomputes on rebuild).
func syncSubscriptionEffect(ctx context.Context, repos *repo.AllRepos, tx *gorm.DB, organizationId int64, plan *model.PricingPlan, oldEffective, newEffective bool) error {
	if oldEffective == newEffective {
		return nil
	}
	if newEffective {
		return reconcileRoleResourcesForPlanChange(ctx, repos, tx, organizationId, nil, plan)
	}
	return reconcileRoleResourcesForPlanChange(ctx, repos, tx, organizationId, plan, nil)
}

func reconcileRoleResourcesForPlanChange(ctx context.Context, repos *repo.AllRepos, tx *gorm.DB, organizationId int64, sourcePlan, targetPlan *model.PricingPlan) error {
	plan := sourcePlan
	if targetPlan != nil {
		plan = targetPlan
	}
	if plan == nil {
		return nil
	}
	roles, err := repos.RoleRepo.FindByOrganizationIdsWithTx(ctx, tx, []int64{organizationId}, plan.Type)
	if err != nil {
		return err
	}
	sourceIDs := map[int64]bool{}
	if sourcePlan != nil {
		for _, resource := range sourcePlan.Resources {
			if resource != nil && !resource.IsPublic {
				sourceIDs[resource.Id] = true
			}
		}
	}
	roleIds := make([]int64, 0, len(roles))
	for _, role := range roles {
		roleIds = append(roleIds, role.Id)
	}
	sourceResourceIds := make([]int64, 0, len(sourceIDs))
	for id := range sourceIDs {
		sourceResourceIds = append(sourceResourceIds, id)
	}
	roleResourceScopes := map[int64]map[int64]constants.DataScope{}
	if sourcePlan != nil {
		relations, relationErr := repos.RoleAndResourceRepo.FindByRoleIdsAndResourceIdsWithTx(ctx, tx, roleIds, sourceResourceIds)
		if relationErr != nil {
			return relationErr
		}
		for _, relation := range relations {
			if relation == nil {
				continue
			}
			if roleResourceScopes[relation.RoleId] == nil {
				roleResourceScopes[relation.RoleId] = map[int64]constants.DataScope{}
			}
			roleResourceScopes[relation.RoleId][relation.ResourceId] = relation.DataScope
		}
	}
	targetIDs := map[int64]bool{}
	if targetPlan != nil {
		for _, resource := range targetPlan.Resources {
			if resource != nil && !resource.IsPublic && resource.Status == constants.CommonStatusActive {
				targetIDs[resource.Id] = true
			}
		}
	}
	currentScopes := map[int64]map[int64]constants.DataScope{}
	if len(targetIDs) > 0 {
		targetResourceIds := make([]int64, 0, len(targetIDs))
		for id := range targetIDs {
			targetResourceIds = append(targetResourceIds, id)
		}
		relations, relationErr := repos.RoleAndResourceRepo.FindByRoleIdsAndResourceIdsWithTx(ctx, tx, roleIds, targetResourceIds)
		if relationErr != nil {
			return relationErr
		}
		for _, relation := range relations {
			if currentScopes[relation.RoleId] == nil {
				currentScopes[relation.RoleId] = map[int64]constants.DataScope{}
			}
			currentScopes[relation.RoleId][relation.ResourceId] = relation.DataScope
		}
	}
	history := map[int64]map[int64]constants.DataScope{}
	if targetPlan != nil {
		snapshots, findErr := repos.RoleResourceSnapshotRepo.FindByOrganizationAndPlanWithTx(ctx, tx, organizationId, targetPlan.Id)
		if findErr != nil {
			return findErr
		}
		for _, snapshot := range snapshots {
			history[snapshot.RoleId] = map[int64]constants.DataScope{}
			for _, resource := range snapshot.Resources {
				scope := resource.DataScope
				if scope == "" {
					scope = constants.DataScopeAll
				}
				history[snapshot.RoleId][resource.ResourceId] = scope
			}
		}
	}
	for _, role := range roles {
		if role.Name == constants.AdminRoleName {
			continue
		}
		if sourcePlan != nil {
			resources := make([]model.RoleResourceSnapshotResource, 0)
			for _, resource := range role.Resources {
				if resource != nil && !resource.IsPublic && sourceIDs[resource.Id] {
					scope := constants.DataScopeAll
					if roleResourceScopes[role.Id] != nil && roleResourceScopes[role.Id][resource.Id] != "" {
						scope = roleResourceScopes[role.Id][resource.Id]
					}
					resources = append(resources, model.RoleResourceSnapshotResource{ResourceId: resource.Id, DataScope: scope})
				}
			}
			if err = repos.RoleResourceSnapshotRepo.SaveWithTx(ctx, tx, &model.RoleResourceSnapshot{OrganizationId: organizationId, PricingPlanId: sourcePlan.Id, RoleId: role.Id, Resources: resources}); err != nil {
				return err
			}
		}
		wanted := map[int64]bool{}
		for _, resource := range role.Resources {
			if resource != nil && !resource.IsPublic && targetIDs[resource.Id] {
				wanted[resource.Id] = true
			}
		}
		for id := range history[role.Id] {
			if targetIDs[id] {
				wanted[id] = true
			}
		}
		var add, remove []*model.Resource
		addScopes := map[int64]constants.DataScope{}
		for _, resource := range role.Resources {
			if resource != nil && !resource.IsPublic && !wanted[resource.Id] {
				remove = append(remove, &model.Resource{Base: model.Base{Id: resource.Id}})
			}
		}
		for id := range wanted {
			if _, found := currentScopes[role.Id][id]; !found {
				add = append(add, &model.Resource{Base: model.Base{Id: id}})
				addScopes[id] = history[role.Id][id]
			}
		}
		if err = repos.RoleRepo.UpdateResourcesWithTx(ctx, tx, role, nil, remove); err != nil {
			return err
		}
		for resourceId, snapshotScope := range history[role.Id] {
			if !wanted[resourceId] || snapshotScope != constants.DataScopeAll {
				continue
			}
			if currentScopes[role.Id][resourceId] != constants.DataScopeAll {
				if err = repos.RoleAndResourceRepo.UpdateDataScopeWithTx(ctx, tx, role.Id, resourceId, constants.DataScopeAll); err != nil {
					return err
				}
			}
		}
		newRelations := make([]*model.RoleAndResource, 0, len(add))
		for _, resource := range add {
			scope := addScopes[resource.Id]
			if scope == "" {
				scope = constants.DataScopeAll
			}
			newRelations = append(newRelations, &model.RoleAndResource{RoleId: role.Id, ResourceId: resource.Id, DataScope: scope})
		}
		if len(newRelations) > 0 {
			if _, err = repos.RoleAndResourceRepo.BatchCreateWithTx(ctx, tx, newRelations, false); err != nil {
				return err
			}
		}
	}
	return nil
}
