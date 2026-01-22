package repo

import (
	"gorm.io/gorm"

	"zeebee-backend/common/gormx"
	"zeebee-backend/domains/zeebee/model"
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
