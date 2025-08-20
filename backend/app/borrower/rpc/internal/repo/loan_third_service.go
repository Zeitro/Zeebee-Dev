package repo

import (
	"gorm.io/gorm"

	"zeebee-backend/app/borrower/rpc/internal/repo/model"
	"zeebee-backend/common/gormx"
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
