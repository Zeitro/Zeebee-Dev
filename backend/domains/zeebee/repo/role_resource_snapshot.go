package repo

import (
	"context"

	"gorm.io/gorm"

	"zeebee-backend/common/errorx"
	"zeebee-backend/common/gormx"
	"zeebee-backend/domains/zeebee/model"
)

type RoleResourceSnapshotRepo struct {
	*gormx.BaseRepo[*model.RoleResourceSnapshot]
}

func NewRoleResourceSnapshotRepo(db *gorm.DB) *RoleResourceSnapshotRepo {
	return &RoleResourceSnapshotRepo{BaseRepo: gormx.NewBaseRepo[*model.RoleResourceSnapshot](db)}
}

func (r *RoleResourceSnapshotRepo) SaveWithTx(ctx context.Context, tx *gorm.DB, snapshot *model.RoleResourceSnapshot) error {
	var existing model.RoleResourceSnapshot
	err := tx.WithContext(ctx).Where("organization_id = ? AND pricing_plan_id = ? AND role_id = ?", snapshot.OrganizationId, snapshot.PricingPlanId, snapshot.RoleId).First(&existing).Error
	if err == nil {
		existing.Resources = snapshot.Resources
		return errorx.WithStack(tx.WithContext(ctx).Save(&existing).Error)
	}
	if err != gorm.ErrRecordNotFound {
		return errorx.WithStack(err)
	}
	return errorx.WithStack(tx.WithContext(ctx).Create(snapshot).Error)
}

func (r *RoleResourceSnapshotRepo) FindByOrganizationAndPlanWithTx(ctx context.Context, tx *gorm.DB, organizationId, pricingPlanId int64) (snapshots []*model.RoleResourceSnapshot, err error) {
	err = tx.WithContext(ctx).Where("organization_id = ? AND pricing_plan_id = ?", organizationId, pricingPlanId).Find(&snapshots).Error
	return snapshots, errorx.WithStack(err)
}

func (r *RoleResourceSnapshotRepo) DeleteByRoleIdWithTx(ctx context.Context, tx *gorm.DB, roleId int64) error {
	return errorx.WithStack(tx.WithContext(ctx).Where("role_id = ?", roleId).Delete(&model.RoleResourceSnapshot{}).Error)
}

func (r *RoleResourceSnapshotRepo) DeleteByPricingPlanIdWithTx(ctx context.Context, tx *gorm.DB, pricingPlanId int64) error {
	return errorx.WithStack(tx.WithContext(ctx).Where("pricing_plan_id = ?", pricingPlanId).Delete(&model.RoleResourceSnapshot{}).Error)
}
