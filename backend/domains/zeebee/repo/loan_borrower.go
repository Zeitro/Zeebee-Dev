package repo

import (
	"context"

	"github.com/pkg/errors"
	"gorm.io/gorm"

	"zeebee-backend/common/gormx"
	"zeebee-backend/domains/zeebee/model"
)

type LoanBorrowerRepo struct {
	*gormx.BaseRepo[*model.LoanBorrower]
	m *model.LoanBorrower
}

func NewLoanBorrowerRepo(db *gorm.DB) *LoanBorrowerRepo {
	return &LoanBorrowerRepo{
		BaseRepo: gormx.NewBaseRepo[*model.LoanBorrower](db),
		m:        new(model.LoanBorrower),
	}
}

// UpdateOneWithTx
func (repo *LoanBorrowerRepo) UpdateOneWithTx(ctx context.Context, tx *gorm.DB, id int64, data *model.LoanBorrower) (result *model.LoanBorrower, err error) {
	err = tx.WithContext(ctx).Select("*").
		Omit("id", "created_at", "loan_id").
		Where("id = ?", id).
		Updates(data).
		Error
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return data, nil
}
func (repo *LoanBorrowerRepo) DeleteByLoanIdWithTx(ctx context.Context, tx *gorm.DB, loanId int64) (err error) {
	err = tx.WithContext(ctx).
		Model(&model.LoanBorrower{}).
		Where("loan_id = ?", loanId).
		Delete(&model.LoanBorrower{}).
		Error
	return errors.WithStack(err)
}
