package repo

import (
	"gorm.io/gorm"

	"zeebee-backend/app/lender/rpc/internal/repo/model"
	"zeebee-backend/common/gormx"
)

type LoanDocumentRepo struct {
	*gormx.BaseRepo[*model.LoanDocument]
	m *model.LoanDocument
}

func NewLoanDocumentRepo(db *gorm.DB) *LoanDocumentRepo {
	return &LoanDocumentRepo{
		BaseRepo: gormx.NewBaseRepo[*model.LoanDocument](db),
		m:        new(model.LoanDocument),
	}
}
