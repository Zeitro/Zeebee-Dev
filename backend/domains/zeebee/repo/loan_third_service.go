package repo

import (
	"context"

	"gorm.io/gorm"

	"zeebee-backend/common/errorx"
	"zeebee-backend/common/gormx"
	"zeebee-backend/domains/zeebee/model"
)

type LoanThirdServiceRepo struct {
	*gormx.BaseRepo[*model.LoanThirdService]
	m *model.LoanThirdService
}

func NewLoanThirdServiceRepo(db *gorm.DB) *LoanThirdServiceRepo {
	return &LoanThirdServiceRepo{
		BaseRepo: gormx.NewBaseRepo[*model.LoanThirdService](db),
		m:        new(model.LoanThirdService),
	}
}

// UpdateOneWithTx
func (repo *LoanThirdServiceRepo) UpdateOneWithTx(ctx context.Context, tx *gorm.DB, id int64, data *model.LoanThirdService) (result *model.LoanThirdService, err error) {
	err = tx.WithContext(ctx).Select("*").Omit("id", "created_at").Where("id = ?", id).Updates(data).Error
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	return data, nil
}

func (repo *LoanThirdServiceRepo) DeleteByLoanIdWithTx(ctx context.Context, tx *gorm.DB, loanId int64) (err error) {
	err = tx.WithContext(ctx).
		Model(&model.LoanThirdService{}).
		Where("loan_id = ?", loanId).
		Delete(&model.LoanThirdService{}).
		Error
	return errorx.WithStack(err)
}
