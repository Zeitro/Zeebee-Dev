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

// GetApiWithTx return api by id with transaction
func (repo *ApiRepo) GetApiWithTx(ctx context.Context, tx *gorm.DB, id int64, apiType constants.UserType, withPreload bool) (result *model.Api, err error) {
	query := tx.WithContext(ctx).Model(&model.Api{}).
		Where("id = ?", id).
		Where("type = ?", apiType)
	if withPreload {
		query = query.
			Preload("Resources")
	}
	err = query.First(&result).Error
	return result, errorx.WithStack(err)
}

// DeleteApiWithTx delete api with transaction
func (repo *ApiRepo) DeleteApiWithTx(ctx context.Context, tx *gorm.DB, api *model.Api) (err error) {
	if api == nil {
		return nil
	}
	resources := api.Resources
	if len(resources) > 0 {
		err = tx.WithContext(ctx).
			Model(api).
			Association("Resources").
			Clear()
		if err != nil {
			return errorx.WithStack(err)
		}
	}
	err = tx.WithContext(ctx).
		Delete(api).Error
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

// GetOrCreateApiWithTx
func (repo *ApiRepo) GetOrCreateApiWithTx(ctx context.Context, tx *gorm.DB, api *model.Api) (result *model.Api, err error) {
	if api == nil {
		return nil, errorx.New("Api is nil")
	}
	query := tx.WithContext(ctx).Model(&model.Api{})

	err = query.Where("name = ? and path = ? and method = ? and type = ?", api.Name, api.Path, api.Method, api.Type).Find(&result).Error
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	if result != nil && result.Id > 0 {
		return result, nil
	}
	//  create
	err = query.Create(&api).Error
	return api, errorx.WithStack(err)
}

// FindApiCategoryByNameAndTypeWithTx finds a category by name and type (path=” AND method=”)
func (repo *ApiRepo) FindApiCategoryByNameAndTypeWithTx(ctx context.Context, tx *gorm.DB, name string, apiType constants.UserType) (result *model.Api, err error) {
	err = tx.WithContext(ctx).Model(&model.Api{}).
		Where("name = ? AND type = ? AND path = '' AND method = ''", name, apiType).
		First(&result).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return result, errorx.WithStack(err)
}

// FindApiByPathMethodAndTypeWithTx finds an API by path, method and type
func (repo *ApiRepo) FindApiByPathMethodAndTypeWithTx(ctx context.Context, tx *gorm.DB, path, method string, apiType constants.UserType) (result *model.Api, err error) {
	err = tx.WithContext(ctx).Model(&model.Api{}).
		Where("path = ? AND method = ? AND type = ?", path, method, apiType).
		First(&result).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return result, errorx.WithStack(err)
}

// GetMaxSortOrderWithTx returns the max sort_order for a given parent_id and type
func (repo *ApiRepo) GetMaxSortOrderWithTx(ctx context.Context, tx *gorm.DB, parentId int64, apiType constants.UserType) (result int32, err error) {
	err = tx.WithContext(ctx).Model(&model.Api{}).
		Where("parent_id = ? AND type = ?", parentId, apiType).
		Select("COALESCE(MAX(sort_order), 0)").
		Scan(&result).Error
	return result, errorx.WithStack(err)
}
