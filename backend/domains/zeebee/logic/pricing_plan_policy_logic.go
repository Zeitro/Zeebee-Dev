package logic

import (
	"context"

	"gorm.io/gorm"

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

	policyModel, err := copierx.Trans[dto.PricingPlanPolicyCreateRequestDto, model.PricingPlanPolicy](&input)
	if err != nil {
		return nil, err
	}

	// duplicate-check and create run in one transaction to close the TOCTOU gap
	err = l.repos.PricingPlanPolicyRepo.Transaction(ctx, func(tx *gorm.DB) error {
		exists, innerErr := l.repos.PricingPlanPolicyRepo.IsExistedWithFilterWithTx(ctx, tx, filters)
		if innerErr != nil {
			return innerErr
		}
		if exists {
			return errorx.New("pricing plan policy with same type already exists")
		}
		policyModel, innerErr = l.repos.PricingPlanPolicyRepo.CreateWithTx(ctx, tx, policyModel, false)
		return innerErr
	})
	if err != nil {
		return nil, err
	}

	return copierx.Trans[model.PricingPlanPolicy, dto.PricingPlanPolicyEntityDto](policyModel)
}

func (l *PricingPlanPolicyLogic) Update(ctx context.Context, input dto.PricingPlanPolicyUpdateRequestDto) (*dto.PricingPlanPolicyEntityDto, error) {
	var result *model.PricingPlanPolicy
	// load, duplicate-check, and update run in one transaction to close the TOCTOU gap
	err := l.repos.PricingPlanPolicyRepo.Transaction(ctx, func(tx *gorm.DB) error {
		policy, innerErr := l.repos.PricingPlanPolicyRepo.FindOneByIdWithTx(ctx, tx, input.Id)
		if innerErr != nil {
			return errorx.Wrap(innerErr, "fail to find pricing plan policy")
		}

		innerErr = copierx.Assign[dto.PricingPlanPolicyUpdateRequestDto, model.PricingPlanPolicy](&input, policy)
		if innerErr != nil {
			return innerErr
		}

		// check policy type uniqueness (excluding self)
		filters := []gormx.WhereRaw{
			{
				Raw:  "pricing_plan_id = ? and type = ? and id != ?",
				Args: []interface{}{policy.PricingPlanId, input.Type, input.Id},
			},
		}
		exists, innerErr := l.repos.PricingPlanPolicyRepo.IsExistedWithFilterWithTx(ctx, tx, filters)
		if innerErr != nil {
			return innerErr
		}
		if exists {
			return errorx.New("pricing plan policy with same type already exists")
		}

		result, innerErr = l.repos.PricingPlanPolicyRepo.UpdateOneByIdWithTx(ctx, tx, input.Id, policy)
		return innerErr
	})
	if err != nil {
		return nil, err
	}

	return copierx.Trans[model.PricingPlanPolicy, dto.PricingPlanPolicyEntityDto](result)
}

func (l *PricingPlanPolicyLogic) Delete(ctx context.Context, id int64) error {
	_, err := l.repos.PricingPlanPolicyRepo.GetOne(ctx, id)
	if err != nil {
		return errorx.Wrap(err, "fail to find pricing plan policy")
	}

	_, err = l.repos.PricingPlanPolicyRepo.DeleteByID(ctx, id)
	return err
}
