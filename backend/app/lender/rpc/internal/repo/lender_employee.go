package repo

import (
	"context"

	"github.com/pkg/errors"
	"gorm.io/gorm"

	"zeebee-backend/app/lender/rpc/internal/repo/model"
	"zeebee-backend/common/gormx"
)

type LenderEmployeeRepo struct {
	*gormx.BaseRepo[*model.LenderEmployee]
	m *model.LenderEmployee
}

func NewLenderEmployeeRepo(db *gorm.DB) *LenderEmployeeRepo {
	return &LenderEmployeeRepo{
		BaseRepo: gormx.NewBaseRepo[*model.LenderEmployee](db),
		m:        new(model.LenderEmployee),
	}
}

func (repo *LenderEmployeeRepo) GetLenderEmployee(ctx context.Context, id int64) (res *model.LenderEmployee, err error) {
	err = repo.Client(ctx).
		Where("id = ?", id).
		Preload("LenderCompanies").
		First(&res).
		Error
	return res, errors.WithStack(err)
}

func (repo *LenderEmployeeRepo) CreateEmployeeWithTx(ctx context.Context, tx *gorm.DB, data *model.LenderEmployee) (res *model.LenderEmployee, err error) {
	err = tx.WithContext(ctx).Create(data).Error
	if err != nil {
		return res, errors.WithStack(err)
	}
	res = data

	return res, errors.WithStack(err)
}
