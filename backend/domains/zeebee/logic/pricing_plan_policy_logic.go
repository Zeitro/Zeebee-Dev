package logic

import (
	"context"

	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
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

	exists, err := l.repos.PricingPlanPolicyRepo.ExistsByPlanIdAndType(ctx, input.PricingPlanId, string(input.Type))
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

func (l *PricingPlanPolicyLogic) Update(ctx context.Context, id int64, input dto.PricingPlanPolicyUpdateRequestDto) (*dto.PricingPlanPolicyEntityDto, error) {
	policy, err := l.repos.PricingPlanPolicyRepo.GetOne(ctx, id)
	if err != nil {
		return nil, errorx.Wrap(err, "pricing plan policy not found")
	}

	err = copierx.Assign[dto.PricingPlanPolicyUpdateRequestDto, model.PricingPlanPolicy](&input, policy)
	if err != nil {
		return nil, err
	}

	policy, err = l.repos.PricingPlanPolicyRepo.UpdateOneById(ctx, id, policy)
	if err != nil {
		return nil, err
	}

	return copierx.Trans[model.PricingPlanPolicy, dto.PricingPlanPolicyEntityDto](policy)
}

func (l *PricingPlanPolicyLogic) Delete(ctx context.Context, id int64) error {
	_, err := l.repos.PricingPlanPolicyRepo.GetOne(ctx, id)
	if err != nil {
		return errorx.Wrap(err, "pricing plan policy not found")
	}

	_, err = l.repos.PricingPlanPolicyRepo.DeleteByIDWithTx(ctx, l.repos.PricingPlanPolicyRepo.GetDB(), id)
	return err
}
