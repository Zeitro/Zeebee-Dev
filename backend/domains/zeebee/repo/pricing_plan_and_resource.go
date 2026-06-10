package repo

import (
	"context"

	"gorm.io/gorm"

	"zeebee-backend/common/errorx"
	"zeebee-backend/common/gormx"
	"zeebee-backend/domains/zeebee/model"
)

type PricingPlanAndResourceRepo struct {
	*gormx.BaseRepo[*model.PricingPlanAndResource]
}

func NewPricingPlanAndResourceRepo(db *gorm.DB) *PricingPlanAndResourceRepo {
	return &PricingPlanAndResourceRepo{
		BaseRepo: gormx.NewBaseRepo[*model.PricingPlanAndResource](db),
	}
}

func (repo *PricingPlanAndResourceRepo) DeleteByPlanId(ctx context.Context, pricingPlanId int64) error {
	err := repo.Client(ctx).
		Where("pricing_plan_id = ?", pricingPlanId).
		Delete(&model.PricingPlanAndResource{}).Error
	return errorx.WithStack(err)
}

func (repo *PricingPlanAndResourceRepo) FindByPlanId(ctx context.Context, pricingPlanId int64) (result []*model.PricingPlanAndResource, err error) {
	err = repo.Client(ctx).
		Where("pricing_plan_id = ?", pricingPlanId).
		Find(&result).Error
	return result, errorx.WithStack(err)
}

func (repo *PricingPlanAndResourceRepo) BatchCreate(ctx context.Context, items []*model.PricingPlanAndResource) error {
	if len(items) == 0 {
		return nil
	}
	err := repo.GetDB().WithContext(ctx).Create(&items).Error
	return errorx.WithStack(err)
}

func (repo *PricingPlanAndResourceRepo) FindResourceIdsByPlanId(ctx context.Context, pricingPlanId int64) (ids []int64, err error) {
	err = repo.Client(ctx).
		Where("pricing_plan_id = ?", pricingPlanId).
		Pluck("resource_id", &ids).Error
	return ids, errorx.WithStack(err)
}
