package repo

import (
	"context"

	"gorm.io/gorm"

	"zeebee-backend/common/errorx"
	"zeebee-backend/common/gormx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/model"
)

type RoleAndResourceRepo struct {
	*gormx.BaseRepo[*model.RoleAndResource]
	m *model.RoleAndResource
}

func NewRoleAndResourceRepo(db *gorm.DB) *RoleAndResourceRepo {
	return &RoleAndResourceRepo{
		BaseRepo: gormx.NewBaseRepo[*model.RoleAndResource](db),
		m:        new(model.RoleAndResource),
	}
}

func (repo *RoleAndResourceRepo) GetRelation(ctx context.Context, roleId int64, resourceId int64) (res *model.RoleAndResource, err error) {
	err = repo.Client(ctx).
		Where("role_id = ? and resource_id = ?", roleId, resourceId).
		Find(&res).Error
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	return res, nil
}

// GetRoleRelationsWithTx get role relations
func (repo *RoleAndResourceRepo) GetRoleRelationsWithTx(ctx context.Context, tx *gorm.DB, roleId int64) (results []*model.RoleAndResource, err error) {
	err = tx.WithContext(ctx).Model(&model.RoleAndResource{}).
		Where("role_id = ?", roleId).
		Find(&results).Error
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	return results, nil
}

// FindByRoleIdsAndResourceIdsWithTx finds role-resource bindings matching the given role IDs and resource IDs.
func (repo *RoleAndResourceRepo) FindByRoleIdsAndResourceIdsWithTx(ctx context.Context, tx *gorm.DB, roleIds []int64, resourceIds []int64) (results []*model.RoleAndResource, err error) {
	if len(roleIds) == 0 || len(resourceIds) == 0 {
		return []*model.RoleAndResource{}, nil
	}
	err = tx.WithContext(ctx).Model(&model.RoleAndResource{}).
		Where("role_id IN ? AND resource_id IN ?", roleIds, resourceIds).
		Find(&results).Error
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	return results, nil
}

func (repo *RoleAndResourceRepo) UpdateDataScopeWithTx(ctx context.Context, tx *gorm.DB, roleId, resourceId int64, dataScope constants.DataScope) error {
	err := tx.WithContext(ctx).Model(&model.RoleAndResource{}).
		Where("role_id = ? AND resource_id = ?", roleId, resourceId).
		Update("data_scope", dataScope).Error
	return errorx.WithStack(err)
}
