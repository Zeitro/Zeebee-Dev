package repo

import (
	"context"

	"github.com/pkg/errors"
	"gorm.io/gorm"

	"zeebee-backend/app/borrower/rpc/internal/repo/model"
	"zeebee-backend/common/gormx"
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
	err = tx.WithContext(ctx).Select("*").Omit("id", "created_at").Where("id = ?", id).Updates(data).Error
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return data, nil
}
