package repo

import (
	"context"

	"gorm.io/gorm"

	"zeebee-backend/common/copierx"
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

func (repo *PricingPlanAndResourceRepo) FindResourceIdsByPlanIdWithTx(ctx context.Context, tx *gorm.DB, pricingPlanId int64) (ids []int64, err error) {
	err = tx.WithContext(ctx).Model(&model.PricingPlanAndResource{}).Where("pricing_plan_id = ?", pricingPlanId).Pluck("resource_id", &ids).Error
	return ids, errorx.WithStack(err)
}

// ReplaceResourcesByPlanId replaces all resource bindings for a plan using GORM many2many association
func (repo *PricingPlanAndResourceRepo) ReplaceResourcesByPlanId(ctx context.Context, planId int64, resourceIds []int64) error {
	resources := make([]*model.Resource, 0, len(resourceIds))
	for _, rid := range resourceIds {
		resources = append(resources, &model.Resource{Base: model.Base{Id: rid}})
	}
	plan := &model.PricingPlan{Base: model.Base{Id: planId}}
	err := repo.GetDB().WithContext(ctx).
		Model(plan).
		Association("Resources").
		Replace(resources)
	return errorx.WithStack(err)
}

// UpdateResourcesWithTx updates a plan's resource associations within a transaction.
// New resources are appended, removed resources are deleted, unchanged ones are left alone.
func (repo *PricingPlanAndResourceRepo) UpdateResourcesWithTx(ctx context.Context, tx *gorm.DB, plan *model.PricingPlan, newResources []*model.Resource, removedResources []*model.Resource) error {
	copyPlan, err := copierx.Trans[model.PricingPlan, model.PricingPlan](plan)
	if err != nil {
		return err
	}
	copyPlan.RemoveRelationships()
	query := tx.WithContext(ctx).Model(copyPlan)
	if len(newResources) > 0 {
		err = query.Association("Resources").Append(newResources)
		if err != nil {
			return errorx.WithStack(err)
		}
	}
	copyPlan.RemoveRelationships()
	if len(removedResources) > 0 {
		err = query.Association("Resources").Delete(removedResources)
		if err != nil {
			return errorx.WithStack(err)
		}
	}
	return nil
}
