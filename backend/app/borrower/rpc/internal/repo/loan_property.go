package repo

import (
	"context"

	"github.com/pkg/errors"
	"gorm.io/gorm"

	"zeebee-backend/app/borrower/rpc/internal/repo/model"
	"zeebee-backend/common/gormx"
)

type LoanPropertyRepo struct {
	*gormx.BaseRepo[*model.LoanProperty]
	m *model.LoanProperty
}

func NewLoanPropertyRepo(db *gorm.DB) *LoanPropertyRepo {
	return &LoanPropertyRepo{
		BaseRepo: gormx.NewBaseRepo[*model.LoanProperty](db),
		m:        new(model.LoanProperty),
	}
}

// FindOneByLoanId find one by loan id
func (repo *LoanPropertyRepo) FindOneByLoanId(ctx context.Context, loanId int64) (result *model.LoanProperty, err error) {
	err = repo.Client(ctx).
		Where("loan_id = ?", loanId).
		Preload("PropertyBase").
		First(&result).Error
	return result, errors.WithStack(err)
}
