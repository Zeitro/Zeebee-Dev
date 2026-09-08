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

	var replacedStripeSubIds []string
	err = l.repos.OrganizationSubscriptionRepo.Transaction(ctx, func(tx *gorm.DB) error {
		var innerErr error
		// fetch current non-expired subscriptions for the org
		oldSubs, innerErr := l.repos.OrganizationSubscriptionRepo.FindActiveByOrgIdWithTx(ctx, tx, input.OrganizationId)
		if innerErr != nil {
			return innerErr
		}
		replacedStripeSubIds = effectiveStripeSubscriptionIds(oldSubs)
		// unbind resources for currently effective old subscriptions (remove-side)
		if _, innerErr = unbindEffectiveSubscriptions(ctx, l.repos, tx, oldSubs); innerErr != nil {
			return innerErr
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
		return nil
	})
	if err != nil {
		return nil, err
	}

	// a replaced Stripe subscription would keep charging: cancel it remotely
	// (best-effort; failures are logged for manual follow-up)
	for _, stripeSubId := range replacedStripeSubIds {
		l.cancelStripeSubscriptionQuietly(ctx, stripeSubId)
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

	plan, err := l.repos.PricingPlanRepo.GetOne(ctx, subscription.PricingPlanId, []gormx.Preload{{Query: "Resources"}})
	if err != nil || plan == nil {
		err = errorx.New("fail to find pricing plan")
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

	err = l.repos.OrganizationSubscriptionRepo.Transaction(ctx, func(tx *gorm.DB) error {
		var innerErr error
		subscription, innerErr = l.repos.OrganizationSubscriptionRepo.UpdateOneByIdWithTx(ctx, tx, id, subscription)
		if innerErr != nil {
			return innerErr
		}
		// M5: only becoming ineffective unbinds (remove-side); the rebuild recomputes the gate.
		_, innerErr = syncSubscriptionEffect(ctx, l.repos, tx, subscriptionEffectSyncReq{
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

	plan, err := l.repos.PricingPlanRepo.GetOne(ctx, subscription.PricingPlanId, []gormx.Preload{{Query: "Resources"}})
	if err != nil {
		return errorx.Wrap(err, "fail to find pricing plan")
	}

	err = l.repos.OrganizationSubscriptionRepo.Transaction(ctx, func(tx *gorm.DB) error {
		var innerErr error
		_, innerErr = l.repos.OrganizationSubscriptionRepo.UpdateOneByIdWithTx(ctx, tx, input.Id, subscription)
		if innerErr != nil {
			return innerErr
		}
		// M5: only becoming ineffective unbinds (remove-side); the rebuild recomputes the gate.
		_, innerErr = syncSubscriptionEffect(ctx, l.repos, tx, subscriptionEffectSyncReq{
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
	newEffective := subscription.Status.IsEffective()

	plan, err := l.repos.PricingPlanRepo.GetOne(ctx, subscription.PricingPlanId, []gormx.Preload{{Query: "Resources"}})
	if err != nil {
		return errorx.Wrap(err, "fail to find pricing plan")
	}

	err = l.repos.OrganizationSubscriptionRepo.Transaction(ctx, func(tx *gorm.DB) error {
		var innerErr error
		_, innerErr = l.repos.OrganizationSubscriptionRepo.UpdateOneByIdWithTx(ctx, tx, id, subscription)
		if innerErr != nil {
			return innerErr
		}
		// M5: only becoming ineffective unbinds (remove-side); the rebuild recomputes the gate.
		_, innerErr = syncSubscriptionEffect(ctx, l.repos, tx, subscriptionEffectSyncReq{
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

	l.publishOrgPolicyRebuild(ctx, []int64{subscription.OrganizationId}, plan.Type)
	return nil
}

// SuspendEffectiveByOrganizationId suspends every effective subscription of an
// org and unbinds its plan resources (remove-side). Used when an org is disabled.
// Subscriptions are suspended, not deleted (no delete op exists).
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

	l.publishOrgPolicyRebuild(ctx, orgPolicyKeys(organizationPolicies), userType)
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

// SubscribeDefaultPlan subscribes an organization to the Active default plan:
// with TrialDays > 0 it creates a Trialing subscription; otherwise it activates
// the plan's free (ActualAmount == 0) policy directly so a new org still gets
// access without a trial. No-op when neither option is available.
func (l *OrganizationSubscriptionLogic) SubscribeDefaultPlan(ctx context.Context, organizationId int64, userType constants.UserType) error {
	plan, err := l.repos.PricingPlanRepo.FindDefaultByType(ctx, userType)
	if err != nil || plan == nil {
		return nil
	}

	if plan.Status != constants.CommonStatusActive {
		return nil
	}

	plan, err = l.repos.PricingPlanRepo.GetOne(ctx, plan.Id, []gormx.Preload{{Query: "Resources"}, {Query: "Policies"}})
	if err != nil || plan == nil {
		return nil
	}

	code, err := l.generateCode(ctx, nil)
	if err != nil {
		return err
	}

	subscription := &model.OrganizationSubscription{
		Code:           code,
		OrganizationId: organizationId,
		PricingPlanId:  plan.Id,
	}
	if plan.TrialDays > 0 {
		subscription.PricingPlanPolicyId = 0
		subscription.Status = constants.OrganizationSubscriptionStatusTrialing
		subscription.ExpireDate = &sql.NullTime{Time: time.Now().AddDate(0, 0, int(plan.TrialDays)), Valid: true}
	} else {
		freePolicy, _ := slice.FindBy(plan.Policies, func(index int, item *model.PricingPlanPolicy) bool {
			return item.ActualAmount == 0
		})
		if freePolicy == nil {
			return nil
		}
		subscription.PricingPlanPolicyId = freePolicy.Id
		subscription.Status = constants.OrganizationSubscriptionStatusActive
		subscription.ExpireDate = calculateExpireDate(freePolicy.Type, time.Now())
		subscription.OriginalAmount = freePolicy.OriginalAmount
	}

	var replacedStripeSubIds []string
	err = l.repos.OrganizationSubscriptionRepo.Transaction(ctx, func(tx *gorm.DB) error {
		activeSubs, innerErr := l.repos.OrganizationSubscriptionRepo.FindActiveByOrgIdWithTx(ctx, tx, organizationId)
		if innerErr != nil {
			return innerErr
		}
		replacedStripeSubIds = effectiveStripeSubscriptionIds(activeSubs)
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

	// stop remote billing for replaced Stripe subscriptions (best-effort)
	for _, stripeSubId := range replacedStripeSubIds {
		l.cancelStripeSubscriptionQuietly(ctx, stripeSubId)
	}

	// recompute this org's policies: the new subscription expanded its plan-resource scope.
	l.publishOrgPolicyRebuild(ctx, []int64{organizationId}, plan.Type)

	return nil
}

func (l *OrganizationSubscriptionLogic) UpdateStatus(ctx context.Context, organizationType constants.UserType) (err error) {
	// default plan + its free (ActualAmount == 0) policy, used to auto-recover
	// expired subscriptions so an org never fully loses access on expiry.
	planFilter := []gormx.WhereRaw{
		{
			Raw: "is_default = ? and type = ?",
			Args: []any{
				true,
				organizationType,
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
		// Stripe auto-renewing subscriptions are normally extended by the
		// invoice.paid webhook before they expire. The cron only acts as a
		// fallback: expiry gets a grace window for webhook delivery/retry,
		// and the Expiring reminder is skipped (renewal is automatic).
		stripeAutoRenew := sub.Source == constants.OrganizationSubscriptionSourceStripe && sub.IsAutoRenew
		deadline := expireDate
		if stripeAutoRenew {
			deadline = expireDate.AddDate(0, 0, stripeAutoRenewGraceDays)
		}
		switch sub.Status {
		case constants.OrganizationSubscriptionStatusActive:
			if now.After(deadline) {
				newStatus = constants.OrganizationSubscriptionStatusExpired
			} else if expireDate.Before(expiringDate) && !stripeAutoRenew {
				newStatus = constants.OrganizationSubscriptionStatusExpiring
			}
		case constants.OrganizationSubscriptionStatusExpiring:
			if now.After(deadline) {
				newStatus = constants.OrganizationSubscriptionStatusExpired
			}
		case constants.OrganizationSubscriptionStatusTrialing:
			if now.After(deadline) {
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

		// force-expired Stripe subscription (webhook never arrived within the
		// grace window): also cancel remotely so Stripe stops charging
		if stripeAutoRenew && sub.StripeSubscriptionId != "" {
			l.cancelStripeSubscriptionQuietly(ctx, sub.StripeSubscriptionId)
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
	if err = l.repos.OrganizationSubscriptionRepo.UpdateColumnsByIdWithTx(ctx, tx, sub.Id, map[string]interface{}{
		"status": constants.OrganizationSubscriptionStatusExpired,
	}); err != nil {
		return nil, err
	}

	if recovery.DefaultPlan == nil || recovery.DefaultPolicy == nil {
		return nil, nil
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

	// unbind old-only resources; overlap kept, new-only left to admin
	_, oldOnly := diffResourceIds(planResourceIds(oldPlan), planResourceIds(recovery.DefaultPlan))
	if len(oldOnly) == 0 {
		return nil, nil
	}
	removed, err := syncPlanResourcesToOrganizationRoles(ctx, l.repos, tx, planResourceSyncReq{
		OrganizationIds: []int64{sub.OrganizationId},
		ResourceIds:     oldOnly,
		UserType:        oldPlan.Type,
		IsBind:          false,
	})
	if err != nil {
		return nil, err
	}
	return removed, nil
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

// orgPolicyKeys returns the org ids tracked in an organizationPolicyChanges map.
func orgPolicyKeys(m map[int64]*organizationPolicyChanges) []int64 {
	out := make([]int64, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
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

// stripeAutoRenewGraceDays is the window the cron leaves for the invoice.paid
// webhook (including Stripe retries) to renew a subscription before it is
// force-expired locally.
const stripeAutoRenewGraceDays = 3

// effectiveStripeSubscriptionIds collects the Stripe subscription ids of
// currently effective Stripe-sourced subscriptions. Callers cancel these
// remotely after replacing them locally, otherwise Stripe keeps charging.
func effectiveStripeSubscriptionIds(subs []*model.OrganizationSubscription) []string {
	var ids []string
	for _, sub := range subs {
		if sub == nil || !sub.Status.IsEffective() {
			continue
		}
		if sub.Source == constants.OrganizationSubscriptionSourceStripe && sub.StripeSubscriptionId != "" {
			ids = append(ids, sub.StripeSubscriptionId)
		}
	}
	return ids
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

// syncSubscriptionEffect propagates a subscription's effective-boundary crossing
// to role-resource bindings. M5: only becoming ineffective unbinds (remove-side);
// becoming effective never auto-binds (admin configures, gate recomputes on rebuild).
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
