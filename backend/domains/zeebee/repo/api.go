package repo

import (
	"context"

	"gorm.io/gorm"

	"zeebee-backend/common/constants"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/gormx"
	"zeebee-backend/domains/zeebee/model"
)

type ApiRepo struct {
	*gormx.BaseRepo[*model.Api]
	m *model.Api
}

func NewApiRepo(db *gorm.DB) *ApiRepo {
	return &ApiRepo{
		BaseRepo: gormx.NewBaseRepo[*model.Api](db),
		m:        new(model.Api),
	}
}

// GetApi return api by id
func (repo *ApiRepo) All(ctx context.Context, apiType constants.UserType) (result []*model.Api, err error) {
	err = repo.Client(ctx).
		Where("type = ?", apiType).
		Order("parent_id asc").
		Order("sort_order asc").
		Find(&result).Error
	return result, errorx.WithStack(err)
}

// IsExisted check api is existed
func (repo *ApiRepo) IsExisted(ctx context.Context, data *model.Api) (isExisted bool, err error) {
	query := repo.Client(ctx).
		Where("name = ?", data.Name).
		Where("type = ?", data.Type).
		Where("path = ?", data.Path).
		Where("method = ?", data.Method)
	if data.Id > 0 {
		query = query.Where("id != ?", data.Id)
	}
	var count int64
	err = query.Count(&count).Error
	return count > 0, errorx.WithStack(err)
}

// GetApi return api by id
func (repo *ApiRepo) GetApi(ctx context.Context, id int64, apiType constants.UserType) (result *model.Api, err error) {
	err = repo.Client(ctx).
		Where("id = ?", id).
		Where("type = ?", apiType).
		First(&result).Error
	return result, errorx.WithStack(err)
}

// GetApisByParentId return apis by parent id
func (repo *ApiRepo) GetApisByParentId(ctx context.Context, parentId int64, apiType constants.UserType) (result []*model.Api, err error) {
	err = repo.Client(ctx).
		Where("parent_id = ? and type = ?", parentId, apiType).
		Order("sort_order asc").
		Find(&result).Error
	return result, errorx.WithStack(err)
}

// GetApis return apis by ids
func (repo *ApiRepo) GetApisWithTx(ctx context.Context, tx *gorm.DB, ids []int64, apiType constants.UserType) (result []*model.Api, err error) {
	err = tx.WithContext(ctx).
		Model(&model.Api{}).
		Where("id in (?) and type = ?", ids, apiType).
		Order("sort_order asc").
		Find(&result).Error
	return result, errorx.WithStack(err)
}
