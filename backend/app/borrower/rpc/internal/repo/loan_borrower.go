package repo

import (
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
