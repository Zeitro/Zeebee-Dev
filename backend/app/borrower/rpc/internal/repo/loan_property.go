package repo

import (
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
