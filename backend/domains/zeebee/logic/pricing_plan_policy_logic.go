package logic

import (
	"context"

	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/gormx"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
)

type PricingPlanPolicyLogic struct {
	*BaseLogic
}

func NewPricingPlanPolicyLogic(baseLogic *BaseLogic) *PricingPlanPolicyLogic {
	return &PricingPlanPolicyLogic{
		BaseLogic: baseLogic,
	}
}

func (l *PricingPlanPolicyLogic) Create(ctx context.Context, input dto.PricingPlanPolicyCreateRequestDto) (*dto.PricingPlanPolicyEntityDto, error) {
	_, err := l.repos.PricingPlanRepo.GetOne(ctx, input.PricingPlanId, nil)
	if err != nil {
		return nil, errorx.Wrap(err, "pricing plan not found")
	}
	filters := []gormx.WhereRaw{
		{
			Raw:  "pricing_plan_id = ? and type = ?",
			Args: []interface{}{input.PricingPlanId, input.Type},
		},
	}

	exists, err := l.repos.PricingPlanPolicyRepo.IsExistedWithFilter(ctx, filters)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errorx.New("pricing plan policy with same type already exists")
	}

	policyModel, err := copierx.Trans[dto.PricingPlanPolicyCreateRequestDto, model.PricingPlanPolicy](&input)
	if err != nil {
		return nil, err
	}

	result, err := l.repos.PricingPlanPolicyRepo.Create(ctx, policyModel)
	if err != nil {
		return nil, err
	}

	return copierx.Trans[model.PricingPlanPolicy, dto.PricingPlanPolicyEntityDto](result)
}

func (l *PricingPlanPolicyLogic) Update(ctx context.Context, input dto.PricingPlanPolicyUpdateRequestDto) (*dto.PricingPlanPolicyEntityDto, error) {
	policy, err := l.repos.PricingPlanPolicyRepo.GetOne(ctx, input.Id)
	if err != nil {
		return nil, errorx.Wrap(err, "fail to find pricing plan policy")
	}

	err = copierx.Assign[dto.PricingPlanPolicyUpdateRequestDto, model.PricingPlanPolicy](&input, policy)
	if err != nil {
		return nil, err
	}

	// check policy type
	filters := []gormx.WhereRaw{
		{
			Raw:  "pricing_plan_id = ? and type = ? and id != ?",
			Args: []interface{}{policy.PricingPlanId, input.Type, input.Id},
		},
	}
	exists, err := l.repos.PricingPlanPolicyRepo.IsExistedWithFilter(ctx, filters)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errorx.New("pricing plan policy with same type already exists")
	}

	policy, err = l.repos.PricingPlanPolicyRepo.UpdateOneById(ctx, input.Id, policy)
	if err != nil {
		return nil, err
	}

	return copierx.Trans[model.PricingPlanPolicy, dto.PricingPlanPolicyEntityDto](policy)
}

func (l *PricingPlanPolicyLogic) Delete(ctx context.Context, id int64) error {
	_, err := l.repos.PricingPlanPolicyRepo.GetOne(ctx, id)
	if err != nil {
		return errorx.Wrap(err, "fail to find pricing plan policy")
	}

	_, err = l.repos.PricingPlanPolicyRepo.DeleteByID(ctx, id)
	return err
}
