package logic

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
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

func (l *OrganizationSubscriptionLogic) Create(ctx context.Context, input dto.OrganizationSubscriptionCreateDto) (result *dto.OrganizationSubscriptionEntity, err error) {
	plan, err := l.repos.PricingPlanRepo.GetOne(ctx, input.PricingPlanId, nil)
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

	policy, err := l.repos.PricingPlanPolicyRepo.GetOne(ctx, input.PricingPlanPolicyId)
	if err != nil {
		return nil, errorx.New("fail to find pricing plan policy")
	}

	if policy.PricingPlanId != input.PricingPlanId {
		err = errorx.New("pricing plan policy not match")
		return
	}

	code := fmt.Sprintf("SU%s%s", time.Now().Format("20060102150405"), randomHex(4))
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
		// set active subscription to inactive
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
		subscription, innerErr = l.repos.OrganizationSubscriptionRepo.CreateWithTx(ctx, tx, subscription, false)
		return innerErr
	})
	if err != nil {
		return nil, err
	}

	// l.publishRoleReloadEvent(ctx, input.OrganizationId, plan.Type)
	return l.transToDetailDto(ctx, subscription)
}

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

	oldEffective := subscriptionEffective(subscription.Status)
	// renew base: start from now if expired, else extend from current expire date
	base := time.Now()
	if subscription.ExpireDate != nil && subscription.ExpireDate.Valid && subscription.ExpireDate.Time.After(base) {
		base = subscription.ExpireDate.Time
	}
	subscription.ExpireDate = calculateExpireDate(policy.Type, base)
	subscription.Status = constants.OrganizationSubscriptionStatusActive
	newEffective := subscriptionEffective(subscription.Status)

	var organizationPolicies map[int64]*organizationPolicyChanges
	err = l.repos.OrganizationSubscriptionRepo.Transaction(ctx, func(tx *gorm.DB) error {
		var innerErr error
		subscription, innerErr = l.repos.OrganizationSubscriptionRepo.UpdateOneByIdWithTx(ctx, tx, id, subscription)
		if innerErr != nil {
			return innerErr
		}
		if oldEffective == newEffective {
			return nil
		}
		resourceIds := planResourceIds(plan)
		if len(resourceIds) == 0 {
			return nil
		}
		organizationPolicies, innerErr = syncPlanResourcesToOrganizationRoles(ctx, l.repos, tx, planResourceSyncReq{
			OrganizationIds: []int64{subscription.OrganizationId},
			ResourceIds:     resourceIds,
			UserType:        plan.Type,
			IsBind:          newEffective,
		})
		return innerErr
	})
	if err != nil {
		return nil, err
	}

	l.publishOrgPolicyReloadEvents(ctx, organizationPolicies, plan.Type)
	return l.transToDetailDto(ctx, subscription)
}

func (l *OrganizationSubscriptionLogic) UpdateExpireDate(ctx context.Context, id int64, expireDate time.Time) error {
	subscription, err := l.repos.OrganizationSubscriptionRepo.GetOne(ctx, id)
	if err != nil {
		return errorx.Wrap(err, "organization subscription not found")
	}

	if subscription.ExpireDate != nil && subscription.ExpireDate.Valid && subscription.ExpireDate.Time.Equal(expireDate) {
		return nil
	}

	oldEffective := subscriptionEffective(subscription.Status)
	subscription.ExpireDate = &sql.NullTime{Time: expireDate, Valid: true}
	if expireDate.After(time.Now()) {
		subscription.Status = constants.OrganizationSubscriptionStatusActive
	} else {
		subscription.Status = constants.OrganizationSubscriptionStatusExpired
	}
	newEffective := subscriptionEffective(subscription.Status)

	plan, _ := l.repos.PricingPlanRepo.GetOne(ctx, subscription.PricingPlanId, []gormx.Preload{{Query: "Resources"}})

	var organizationPolicies map[int64]*organizationPolicyChanges
	err = l.repos.OrganizationSubscriptionRepo.Transaction(ctx, func(tx *gorm.DB) error {
		var innerErr error
		_, innerErr = l.repos.OrganizationSubscriptionRepo.UpdateOneByIdWithTx(ctx, tx, id, subscription)
		if innerErr != nil {
			return innerErr
		}
		if plan == nil || oldEffective == newEffective {
			return nil
		}
		resourceIds := planResourceIds(plan)
		if len(resourceIds) == 0 {
			return nil
		}
		organizationPolicies, innerErr = syncPlanResourcesToOrganizationRoles(ctx, l.repos, tx, planResourceSyncReq{
			OrganizationIds: []int64{subscription.OrganizationId},
			ResourceIds:     resourceIds,
			UserType:        plan.Type,
			IsBind:          newEffective,
		})
		return innerErr
	})
	if err != nil {
		return err
	}

	if plan != nil {
		l.publishOrgPolicyReloadEvents(ctx, organizationPolicies, plan.Type)
	}

	return nil
}

func (l *OrganizationSubscriptionLogic) SwitchStatus(ctx context.Context, id int64, status constants.OrganizationSubscriptionStatus) error {
	subscription, err := l.repos.OrganizationSubscriptionRepo.GetOne(ctx, id)
	if err != nil {
		return errorx.Wrap(err, "organization subscription not found")
	}

	if subscription.Status == status {
		return nil
	}

	subscription.Status = status
	_, err = l.repos.OrganizationSubscriptionRepo.UpdateOneById(ctx, id, subscription)
	if err != nil {
		return err
	}

	plan, _ := l.repos.PricingPlanRepo.GetOne(ctx, subscription.PricingPlanId, nil)
	if plan != nil {
		l.publishRoleReloadEvent(ctx, subscription.OrganizationId, plan.Type)
	}

	return nil
}

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
		result.OrganizationName = org.Name
	}

	plan, innerErr := l.repos.PricingPlanRepo.GetOne(ctx, sub.PricingPlanId, nil)
	if innerErr == nil && plan != nil {
		result.PricingPlanName = plan.Name
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

	code := fmt.Sprintf("SU%s%s", time.Now().Format("20060102150405"), randomHex(4))
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

	resourceIds := planResourceIds(plan)

	var organizationPolicies map[int64]*organizationPolicyChanges
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
		if createErr != nil {
			return createErr
		}

		if len(resourceIds) == 0 {
			return nil
		}
		organizationPolicies, innerErr = syncPlanResourcesToOrganizationRoles(ctx, l.repos, tx, planResourceSyncReq{
			OrganizationIds: []int64{organizationId},
			ResourceIds:     resourceIds,
			UserType:        plan.Type,
			IsBind:          true,
		})
		return innerErr
	})
	if err != nil {
		return err
	}

	l.publishOrgPolicyReloadEvents(ctx, organizationPolicies, plan.Type)
	return nil
}

// randomHex returns n random bytes as a hex-encoded string.
func randomHex(n int) string {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "0000"
	}
	return hex.EncodeToString(b)
}

// subscriptionEffective reports whether a subscription status grants plan resource access.
// Mirrors FindEffectiveByPlanIdWithTx: every status except Expired is effective.
func subscriptionEffective(status constants.OrganizationSubscriptionStatus) bool {
	return status != constants.OrganizationSubscriptionStatusExpired
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

func (l *OrganizationSubscriptionLogic) publishRoleReloadEvent(ctx context.Context, orgId int64, userType constants.UserType) {
	threading.GoSafe(func() {
		subCtx := context.WithoutCancel(ctx)
		msg := auth.Msg{
			Action:         auth.ActionRoleReload,
			Scope:          auth.ScopeOrganization,
			Type:           userType,
			OrganizationId: orgId,
		}
		topic := l.authTopicMap[userType]
		_ = l.publishAuthEvent(subCtx, topic, msg)
	})
}
