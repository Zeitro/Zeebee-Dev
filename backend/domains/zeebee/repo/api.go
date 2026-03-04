package repo

import (
	"context"

	"gorm.io/gorm"

	"zeebee-backend/common/errorx"
	"zeebee-backend/common/gormx"
	"zeebee-backend/domains/zeebee/constants"
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
func (repo *ApiRepo) GetApi(ctx context.Context, id int64, apiType constants.UserType, withPreload bool) (result *model.Api, err error) {
	query := repo.Client(ctx).
		Where("id = ?", id).
		Where("type = ?", apiType)
	if withPreload {
		query = query.
			Preload("Resources")
	}
	err = query.First(&result).Error
	return result, errorx.WithStack(err)
}

// DeleteApi
func (repo *ApiRepo) DeleteApi(ctx context.Context, api *model.Api) (err error) {
	if api == nil {
		return nil
	}
	err = repo.Transaction(ctx, func(tx *gorm.DB) error {
		var innerErr error
		resources := api.Resources
		if len(resources) > 0 {
			innerErr = tx.WithContext(ctx).
				Model(api).
				Association("Resources").
				Clear()
			if innerErr != nil {
				return innerErr
			}
		}
		innerErr = tx.WithContext(ctx).
			Delete(api).Error
		return innerErr
	})

	return errorx.WithStack(err)
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

// GetAuthenticationWhitelist return authentication whitelist
func (repo *ApiRepo) GetAuthenticationWhitelist(ctx context.Context, userType constants.UserType) (whitelist []*model.Api, err error) {
	err = repo.Client(ctx).
		Select("id", "type", "path", "method", "status", "auth_required").
		Where("type = ? and auth_required = ?", userType, false).
		Find(&whitelist).Error
	return whitelist, errorx.WithStack(err)
}

// GetAuthorizationWhitelist return authorized whitelist
func (repo *ApiRepo) GetAuthorizationWhitelist(ctx context.Context, userType constants.UserType) (whitelist []*model.Api, err error) {
	err = repo.Client(ctx).
		Select("id", "type", "path", "method", "status", "permission_required").
		Where("type = ? and permission_required = ?", userType, false).
		Find(&whitelist).Error
	return whitelist, errorx.WithStack(err)
}

// GetAllApis
func (repo *ApiRepo) GetAllApis(ctx context.Context, apiType constants.UserType, withPreload bool) (apis []*model.Api, err error) {
	query := repo.Client(ctx).Where("type = ? and status = ?", apiType, constants.CommonStatusActive)
	if withPreload {
		query = query.Preload("Resources")
	}
	err = query.Find(&apis).Error
	return apis, errorx.WithStack(err)
}
