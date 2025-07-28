package repo

import (
	"gorm.io/gorm"

	"zeebee-backend/app/borrower/rpc/internal/repo/model"
	"zeebee-backend/common/gormx"
)

type LoanRepo struct {
	*gormx.BaseRepo[*model.Loan]
	m *model.Loan
}

func NewLoanRepo(db *gorm.DB) *LoanRepo {
	return &LoanRepo{
		BaseRepo: gormx.NewBaseRepo[*model.Loan](db),
		m:        new(model.Loan),
	}
}