package logic

import (
	"context"
	"database/sql"
	"fmt"
	"time"

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

func (l *OrganizationSubscriptionLogic) Create(ctx context.Context, input dto.OrganizationSubscriptionCreateDto) (*dto.OrganizationSubscriptionEntity, error) {
	plan, err := l.repos.PricingPlanRepo.GetOne(ctx, input.PricingPlanId)
	if err != nil {
		return nil, errorx.Wrap(err, "pricing plan not found")
	}
	if plan.Status != constants.CommonStatusActive {
		return nil, errorx.New("pricing plan not found or disabled")
	}

	policy, err := l.repos.PricingPlanPolicyRepo.GetOne(ctx, input.PricingPlanPolicyId)
	if err != nil {
		return nil, errorx.New("pricing plan policy not found")
	}

	code := fmt.Sprintf("SU%d", time.Now().UnixNano())
	expireDate := l.calculateExpireDate(policy.Type)

	status := constants.OrganizationSubscriptionStatusActive
	if plan.TrialDays > 0 {
		status = constants.OrganizationSubscriptionStatusTrialing
	}

	subscription := &model.OrganizationSubscription{
		Code:                code,
		OrganizationId:      input.OrganizationId,
		PricingPlanId:       input.PricingPlanId,
		PricingPlanPolicyId: input.PricingPlanPolicyId,
		Status:              status,
		ExpireDate:          expireDate,
		OriginalAmount:      policy.OriginalAmount,
		ActualAmount:        policy.ActualAmount,
	}

	err = l.repos.OrganizationSubscriptionRepo.Transaction(ctx, func(tx *gorm.DB) error {
		activeSubs, innerErr := l.repos.OrganizationSubscriptionRepo.FindActiveByOrgIdWithTx(ctx, tx, input.OrganizationId)
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

		var createErr error
		subscription, createErr = l.repos.OrganizationSubscriptionRepo.CreateWithTx(ctx, tx, subscription, false)
		return createErr
	})
	if err != nil {
		return nil, err
	}

	l.publishRoleReloadEvent(ctx, input.OrganizationId, plan.Type)
	return l.transToDetailDto(ctx, subscription)
}

func (l *OrganizationSubscriptionLogic) Renew(ctx context.Context, id int64, input dto.OrganizationSubscriptionRenewDto) (*dto.OrganizationSubscriptionEntity, error) {
	subscription, err := l.repos.OrganizationSubscriptionRepo.GetOne(ctx, id)
	if err != nil {
		return nil, errorx.Wrap(err, "organization subscription not found")
	}

	parsedDate, err := time.Parse(time.RFC3339, input.ExpireDate)
	if err != nil {
		return nil, errorx.New("invalid expire_date format, expected RFC3339")
	}

	subscription.ExpireDate = &sql.NullTime{Time: parsedDate, Valid: true}
	subscription, err = l.repos.OrganizationSubscriptionRepo.UpdateOneById(ctx, id, subscription)
	if err != nil {
		return nil, err
	}

	plan, _ := l.repos.PricingPlanRepo.GetOne(ctx, subscription.PricingPlanId)
	if plan != nil {
		l.publishRoleReloadEvent(ctx, subscription.OrganizationId, plan.Type)
	}

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

	subscription.ExpireDate = &sql.NullTime{Time: expireDate, Valid: true}
	_, err = l.repos.OrganizationSubscriptionRepo.UpdateOneById(ctx, id, subscription)
	if err != nil {
		return err
	}

	plan, _ := l.repos.PricingPlanRepo.GetOne(ctx, subscription.PricingPlanId)
	if plan != nil {
		l.publishRoleReloadEvent(ctx, subscription.OrganizationId, plan.Type)
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

	plan, _ := l.repos.PricingPlanRepo.GetOne(ctx, subscription.PricingPlanId)
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

func (l *OrganizationSubscriptionLogic) calculateExpireDate(policyType constants.PricingPolicyType) *sql.NullTime {
	var days int
	switch policyType {
	case constants.PricingPolicyTypeMonthly:
		days = 30
	case constants.PricingPolicyTypeQuarterly:
		days = 90
	case constants.PricingPolicyTypeYearly:
		days = 365
	default:
		days = 30
	}
	return &sql.NullTime{Time: time.Now().AddDate(0, 0, days), Valid: true}
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

	plan, innerErr := l.repos.PricingPlanRepo.GetOne(ctx, sub.PricingPlanId)
	if innerErr == nil && plan != nil {
		result.PricingPlanName = plan.Name
	}

	return result, nil
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
