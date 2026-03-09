package repo

import (
	"context"

	"gorm.io/gorm"

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

// GetResource return all resources
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
	query := tx.WithContext(ctx).Model(&model.Resource{}).
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

// GetResource get resource by id and api type
func (repo *ResourceRepo) GetResource(ctx context.Context, id int64, apiType constants.UserType, withPreload bool) (result *model.Resource, err error) {
	query := repo.Client(ctx).
		Where("id = ?", id).
		Where("type = ?", apiType)
	if withPreload {
		query = query.Preload("Apis", func(db *gorm.DB) *gorm.DB {
			return db.Order("apis.sort_order asc")
		})
	}
	err = query.First(&result).Error
	return result, errorx.WithStack(err)
}

// GetResourceWithTx get resource with tx
func (repo *ResourceRepo) GetResourceWithTx(ctx context.Context, tx *gorm.DB, id int64, apiType constants.UserType, withPreload bool) (result *model.Resource, err error) {
	query := tx.WithContext(ctx).Model(&model.Resource{}).
		Where("id = ?", id).
		Where("type = ?", apiType)
	if withPreload {
		query = query.Preload("Apis", func(db *gorm.DB) *gorm.DB {
			return db.Order("apis.sort_order asc")
		}).
			Preload("Roles")
	}
	err = query.First(&result).Error
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

// GetResourceByParentIdWithTx get resources with tx
func (repo *ResourceRepo) GetResourceByParentIdWithTx(ctx context.Context, tx *gorm.DB, parentId int64, apiType constants.UserType) (result []*model.Resource, err error) {
	err = tx.WithContext(ctx).Model(&model.Resource{}).
		Where("parent_id = ? and type = ?", parentId, apiType).
		Order("sort_order asc").
		Find(&result).Error
	return result, errorx.WithStack(err)
}

// DeleteResourceWithTx delete resource
func (repo *ResourceRepo) DeleteResourceWithTx(ctx context.Context, tx *gorm.DB, resource *model.Resource) (err error) {
	if resource == nil {
		return
	}
	apis := resource.Apis
	if len(apis) > 0 {
		err = repo.RemoveApisWithTx(ctx, tx, resource, apis)
		if err != nil {
			return err
		}
	}

	// delete role relations
	roleAndResourceTable := model.RoleAndResource{}.TableName()
	err = tx.WithContext(ctx).
		Table(roleAndResourceTable).
		Where("resource_id = ?", resource.Id).
		Delete(&model.RoleAndResource{}).Error
	if err != nil {
		return errorx.WithStack(err)
	}

	// delete resource
	err = tx.WithContext(ctx).
		Delete(resource).Error
	if err != nil {
		return errorx.WithStack(err)
	}

	return nil
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

// AllResourcesForAuth return all resources for auth
func (repo *ResourceRepo) AllResourcesForAuth(ctx context.Context, organizationId int64, apiType constants.UserType) (result []*model.Resource, err error) {
	err = repo.Client(ctx).Where("type = ? and status = ?", apiType, constants.CommonStatusActive).
		Preload("Roles", func(db *gorm.DB) *gorm.DB {
			query := db.Where("status = ?", constants.CommonStatusActive)
			if organizationId > 0 {
				query = query.Where("organization_id = ?", organizationId)
			}

			return query
		}).Find(&result).Error
	return result, errorx.WithStack(err)
}

// AllResources return all resources
func (repo *ResourceRepo) AllResources(ctx context.Context, organizationId int64, resourceType constants.UserType, withPreload bool) (result []*model.Resource, err error) {
	query := repo.Client(ctx).Where("type = ? and status = ?", resourceType, constants.CommonStatusActive)
	if withPreload {
		query = query.
			// Preload("Apis").
			Preload("Roles", func(db *gorm.DB) *gorm.DB {
				subQuery := db
				if organizationId > 0 {
					subQuery = subQuery.Where("organization_id = ?", organizationId)
				}
				return subQuery
			})

	}
	err = query.Find(&result).Error
	return result, errorx.WithStack(err)
}
