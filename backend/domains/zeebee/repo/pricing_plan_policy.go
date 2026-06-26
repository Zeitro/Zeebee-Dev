package repo

import (
	"context"

	"gorm.io/gorm"

	"zeebee-backend/common/errorx"
	"zeebee-backend/common/gormx"
	"zeebee-backend/domains/zeebee/model"
)

type PricingPlanPolicyRepo struct {
	*gormx.BaseRepo[*model.PricingPlanPolicy]
}

func NewPricingPlanPolicyRepo(db *gorm.DB) *PricingPlanPolicyRepo {
	return &PricingPlanPolicyRepo{
		BaseRepo: gormx.NewBaseRepo[*model.PricingPlanPolicy](db),
	}
}

func (repo *PricingPlanPolicyRepo) GetOne(ctx context.Context, id int64) (result *model.PricingPlanPolicy, err error) {
	err = repo.Client(ctx).Where("id = ?", id).First(&result).Error
	return result, errorx.WithStack(err)
}

func (repo *PricingPlanPolicyRepo) ExistsByPlanIdAndType(ctx context.Context, pricingPlanId int64, policyType string) (bool, error) {
	return repo.IsExistedByColumns(ctx, map[string]interface{}{
		"pricing_plan_id": pricingPlanId,
		"type":            policyType,
	})
}

func (repo *PricingPlanPolicyRepo) ListByPlanId(ctx context.Context, pricingPlanId int64) (result []*model.PricingPlanPolicy, err error) {
	err = repo.Client(ctx).
		Where("pricing_plan_id = ?", pricingPlanId).
		Order("id asc").
		Find(&result).Error
	return result, errorx.WithStack(err)
}
