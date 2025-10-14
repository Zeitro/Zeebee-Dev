package repo

import (
	"gorm.io/gorm"

	"zeebee-backend/app/borrower/rpc/internal/repo/model"
	"zeebee-backend/common/gormx"
)

type LoanOperationRepo struct {
	*gormx.BaseRepo[*model.LoanOperation]
	m *model.LoanOperation
}

func NewLoanOperationRepo(db *gorm.DB) *LoanOperationRepo {
	return &LoanOperationRepo{
		BaseRepo: gormx.NewBaseRepo[*model.LoanOperation](db),
		m:        new(model.LoanOperation),
	}
}
