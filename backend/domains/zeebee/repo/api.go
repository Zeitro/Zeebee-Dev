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
