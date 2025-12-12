package repo

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"zeebee-backend/common/errorx"
	"zeebee-backend/common/gormx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/model"
)

type ResourceRepo struct {
	*gormx.BaseRepo[*model.Resource]
	m *model.Resource
}

func NewResourceRepo(db *gorm.DB) *ResourceRepo {
	return &ResourceRepo{
		BaseRepo: gormx.NewBaseRepo[*model.Resource](db),
		m:        new(model.Resource),
	}
}

// GetResource return api by id
func (repo *ResourceRepo) All(ctx context.Context, apiType constants.UserType) (result []*model.Resource, err error) {
	err = repo.Client(ctx).
		Where("type = ?", apiType).
		Order("parent_id asc").
		Order("sort_order asc").
		Find(&result).Error
	return result, errorx.WithStack(err)
}

// IsExistedWithTx check resource is existed
func (repo *ResourceRepo) IsExistedWithTx(ctx context.Context, tx *gorm.DB, data *model.Resource) (isExisted bool, err error) {
	query := repo.Client(ctx).WithContext(ctx).Model(&model.Resource{}).
		Where("name = ?", data.Name).
		Where("type = ?", data.Type).
		Where("path = ?", data.Path).
		Where("resource_type = ?", data.ResourceType)
	if data.Id > 0 {
		query = query.Where("id != ?", data.Id)
	}
	var count int64
	err = query.Count(&count).Error
	return count > 0, errorx.WithStack(err)
}

// GetResource return api by id
func (repo *ResourceRepo) GetResource(ctx context.Context, id int64, apiType constants.UserType) (result *model.Resource, err error) {
	err = repo.Client(ctx).
		Where("id = ?", id).
		Where("type = ?", apiType).
		Preload("Apis", func(db *gorm.DB) *gorm.DB {
			return db.Order("apis.sort_order asc")
		}).
		First(&result).Error
	return result, errorx.WithStack(err)
}

func (repo *ResourceRepo) GetResourceWithTx(ctx context.Context, tx *gorm.DB, id int64, apiType constants.UserType) (result *model.Resource, err error) {
	err = tx.WithContext(ctx).Model(&model.Resource{}).
		Where("id = ?", id).
		Where("type = ?", apiType).
		Preload("Apis", func(db *gorm.DB) *gorm.DB {
			return db.Order("apis.sort_order asc")
		}).
		First(&result).Error
	return result, errorx.WithStack(err)
}

// GetResourceByParentId get resources by parent id and api type
func (repo *ResourceRepo) GetResourceByParentId(ctx context.Context, parentId int64, apiType constants.UserType) (result []*model.Resource, err error) {
	err = repo.Client(ctx).
		Where("parent_id = ? and type = ?", parentId, apiType).
		Order("sort_order asc").
		Find(&result).Error
	return result, errorx.WithStack(err)
}

// DelByIdWithTx delete resource by id
func (repo *ResourceRepo) DelByIdWithTx(ctx context.Context, tx *gorm.DB, id int64, typ constants.UserType) (err error) {
	var resource *model.Resource
	resource, err = repo.GetResourceWithTx(ctx, tx, id, typ)
	if err != nil {
		return err
	}
	err = tx.WithContext(ctx).
		Model(&model.Resource{}).
		Where("id = ?", id).
		Select(clause.Associations).
		Delete(resource).Error
	return errorx.WithStack(err)
}

// RemoveApis remove apis with tx and return error if occurred
func (repo *ResourceRepo) RemoveApisWithTx(ctx context.Context, tx *gorm.DB, resource *model.Resource, apis []*model.Api) (err error) {
	if resource == nil || len(apis) == 0 {
		return
	}
	err = tx.WithContext(ctx).Model(resource).Association("Apis").Delete(apis)
	return errorx.WithStack(err)
}

// AddApisWithTx add apis with tx and return error if occurred
func (repo *ResourceRepo) AddApisWithTx(ctx context.Context, tx *gorm.DB, resource *model.Resource, apis []*model.Api) (err error) {
	if resource == nil || len(apis) == 0 {
		return
	}
	err = tx.WithContext(ctx).Model(resource).Association("Apis").Append(apis)
	return errorx.WithStack(err)
}

func (repo *ResourceRepo) UpdateApisWithTx(ctx context.Context, tx *gorm.DB, resource *model.Resource, apis []*model.Api) (err error) {
	if resource == nil || len(apis) == 0 {
		return
	}
	err = tx.WithContext(ctx).Model(resource).Association("Apis").Replace(apis)
	return errorx.WithStack(err)
}

// GetResourcesWithTx get resources with tx
func (repo *ResourceRepo) GetResourcesWithTx(ctx context.Context, tx *gorm.DB, ids []int64, apiType constants.UserType) (result []*model.Resource, err error) {
	err = tx.WithContext(ctx).
		Model(&model.Resource{}).
		Where("id in (?) and type = ?", ids, apiType).
		Order("sort_order asc").
		Find(&result).Error
	return result, errorx.WithStack(err)
}
