package logic

import (
	"context"

	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/gormx"
	"zeebee-backend/domains/zeebee/auth"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"

	"github.com/zeromicro/go-zero/core/threading"
)

type PricingPlanLogic struct {
	*BaseLogic
	authTopicMap map[constants.UserType]string
}

func NewPricingPlanLogic(baseLogic *BaseLogic, authTopicMap map[constants.UserType]string) *PricingPlanLogic {
	return &PricingPlanLogic{
		BaseLogic:    baseLogic,
		authTopicMap: authTopicMap,
	}
}

func (l *PricingPlanLogic) Create(ctx context.Context, input dto.PricingPlanCreateDto) (*dto.PricingPlanEntityDto, error) {
	exists, err := l.repos.PricingPlanRepo.ExistsByNameAndType(ctx, input.Name, input.Type)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errorx.New("pricing plan name already exists under this type")
	}

	planModel, err := copierx.Trans[dto.PricingPlanCreateDto, model.PricingPlan](&input)
	if err != nil {
		return nil, err
	}

	result, err := l.repos.PricingPlanRepo.Create(ctx, planModel)
	if err != nil {
		return nil, err
	}

	l.publishResourceReloadEvent(ctx, result.Type)
	return l.transToDetailDto(ctx, result)
}

func (l *PricingPlanLogic) Update(ctx context.Context, id int64, input dto.PricingPlanUpdateDto) (*dto.PricingPlanEntityDto, error) {
	plan, err := l.repos.PricingPlanRepo.GetOne(ctx, id)
	if err != nil {
		return nil, errorx.Wrap(err, "pricing plan not found")
	}

	err = copierx.Assign[dto.PricingPlanUpdateDto, model.PricingPlan](&input, plan)
	if err != nil {
		return nil, err
	}

	plan, err = l.repos.PricingPlanRepo.UpdateOneById(ctx, id, plan)
	if err != nil {
		return nil, err
	}

	l.publishResourceReloadEvent(ctx, plan.Type)
	return l.transToDetailDto(ctx, plan)
}

func (l *PricingPlanLogic) Delete(ctx context.Context, id int64) error {
	plan, err := l.repos.PricingPlanRepo.GetOne(ctx, id)
	if err != nil {
		return errorx.Wrap(err, "pricing plan not found")
	}

	subscriptions, err := l.repos.OrganizationSubscriptionRepo.FindByPricingPlanId(ctx, id)
	if err != nil {
		return err
	}
	for _, sub := range subscriptions {
		if sub.Status != constants.OrganizationSubscriptionStatusExpired {
			return errorx.New("cannot delete pricing plan with active subscriptions")
		}
	}

	_, err = l.repos.PricingPlanRepo.DeleteByIDWithTx(ctx, l.repos.PricingPlanRepo.GetDB(), id)
	if err != nil {
		return err
	}

	l.publishResourceReloadEvent(ctx, plan.Type)
	return nil
}

func (l *PricingPlanLogic) Detail(ctx context.Context, id int64) (*dto.PricingPlanEntityDto, error) {
	plan, err := l.repos.PricingPlanRepo.FindOneWithPreloadById(ctx, id, []string{"Policies"})
	if err != nil {
		return nil, errorx.Wrap(err, "pricing plan not found")
	}
	return l.transToDetailDto(ctx, plan)
}

func (l *PricingPlanLogic) Page(ctx context.Context, req dto.PricingPlanPageRequestDto) (gormx.PageResponse[dto.PricingPlanPageEntity], error) {
	pageResult, err := l.repos.PricingPlanRepo.Page(ctx, req)
	if err != nil {
		return gormx.PageResponse[dto.PricingPlanPageEntity]{}, err
	}

	resultPtr, err := copierx.Trans[gormx.PageResponse[model.PricingPlan], gormx.PageResponse[dto.PricingPlanPageEntity]](&pageResult)
	if err != nil {
		return gormx.PageResponse[dto.PricingPlanPageEntity]{}, err
	}

	return *resultPtr, nil
}

func (l *PricingPlanLogic) SwitchStatus(ctx context.Context, id int64, status constants.CommonStatus) error {
	plan, err := l.repos.PricingPlanRepo.GetOne(ctx, id)
	if err != nil {
		return errorx.Wrap(err, "pricing plan not found")
	}

	if plan.Status == status {
		return nil
	}

	plan.Status = status
	_, err = l.repos.PricingPlanRepo.UpdateOneById(ctx, id, plan)
	if err != nil {
		return err
	}

	l.publishResourceReloadEvent(ctx, plan.Type)
	return nil
}

func (l *PricingPlanLogic) BindResources(ctx context.Context, planId int64, resourceIds []int64) ([]*dto.PricingPlanResourceEntityDto, error) {
	plan, err := l.repos.PricingPlanRepo.GetOne(ctx, planId)
	if err != nil {
		return nil, errorx.Wrap(err, "pricing plan not found")
	}

	err = l.repos.PricingPlanAndResourceRepo.ReplaceResourcesByPlanId(ctx, planId, resourceIds)
	if err != nil {
		return nil, err
	}

	l.publishResourceReloadEvent(ctx, plan.Type)

	bindings, err := l.repos.PricingPlanAndResourceRepo.FindByPlanId(ctx, planId)
	if err != nil {
		return nil, err
	}

	listPtr, err := copierx.Trans[[]*model.PricingPlanAndResource, []*dto.PricingPlanResourceEntityDto](&bindings)
	if err != nil {
		return nil, err
	}
	return *listPtr, nil
}

func (l *PricingPlanLogic) ResourceList(ctx context.Context, planId int64) ([]*dto.PricingPlanResourceEntityDto, error) {
	_, err := l.repos.PricingPlanRepo.GetOne(ctx, planId)
	if err != nil {
		return nil, errorx.Wrap(err, "pricing plan not found")
	}

	bindings, err := l.repos.PricingPlanAndResourceRepo.FindByPlanId(ctx, planId)
	if err != nil {
		return nil, err
	}

	listPtr, err := copierx.Trans[[]*model.PricingPlanAndResource, []*dto.PricingPlanResourceEntityDto](&bindings)
	if err != nil {
		return nil, err
	}
	return *listPtr, nil
}

func (l *PricingPlanLogic) transToDetailDto(ctx context.Context, plan *model.PricingPlan) (*dto.PricingPlanEntityDto, error) {
	result, err := copierx.Trans[model.PricingPlan, dto.PricingPlanEntityDto](plan)
	if err != nil {
		return nil, err
	}

	// load policies if not preloaded
	if plan.Policies == nil {
		policies, innerErr := l.repos.PricingPlanPolicyRepo.ListByPlanId(ctx, plan.Id)
		if innerErr != nil {
			return nil, innerErr
		}
		plan.Policies = policies
	}

	policyDtos, err := copierx.Trans[[]*model.PricingPlanPolicy, []dto.PricingPlanPolicyEntityDto](&plan.Policies)
	if err != nil {
		return nil, err
	}
	result.Policies = *policyDtos

	return result, nil
}

func (l *PricingPlanLogic) publishResourceReloadEvent(ctx context.Context, userType constants.UserType) {
	threading.GoSafe(func() {
		subCtx := context.WithoutCancel(ctx)
		msg := auth.Msg{
			Action: auth.ActionResourceReload,
			Scope:  auth.ScopeGlobal,
			Type:   userType,
		}
		topic := l.authTopicMap[userType]
		_ = l.publishAuthEvent(subCtx, topic, msg)
	})
}
