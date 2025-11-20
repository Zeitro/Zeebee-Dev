package repo

import (
	"context"

	"gorm.io/gorm"

	"zeebee-backend/common/errorx"

	"zeebee-backend/common/gormx"
	"zeebee-backend/domains/zeebee/model"
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

// GetLenderEmployee return one lender employee
func (repo *LenderEmployeeRepo) GetLenderEmployee(ctx context.Context, id int64) (res *model.LenderEmployee, err error) {
	err = repo.Client(ctx).
		Where("id = ?", id).
		Preload("LenderCompanies").
		First(&res).
		Error
	return res, errorx.WithStack(err)
}

func (repo *LenderEmployeeRepo) CreateEmployeeWithTx(ctx context.Context, tx *gorm.DB, data *model.LenderEmployee) (res *model.LenderEmployee, err error) {
	err = tx.WithContext(ctx).Create(data).Error
	if err != nil {
		return res, errorx.WithStack(err)
	}
	res = data

	return res, errorx.WithStack(err)
}
