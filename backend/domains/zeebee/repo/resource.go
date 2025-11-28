package repo

import (
	"context"

	"gorm.io/gorm"

	"zeebee-backend/common/constants"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/gormx"
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

// IsExisted check resource is existed
func (repo *ResourceRepo) IsExisted(ctx context.Context, data *model.Resource) (isExisted bool, err error) {
	query := repo.Client(ctx).
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
		First(&result).Error
	return result, errorx.WithStack(err)
}
