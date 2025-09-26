package repo

import (
	"gorm.io/gorm"

	"zeebee-backend/app/lender/rpc/internal/repo/model"
	"zeebee-backend/common/gormx"
)

type BorrowerRepo struct {
	*gormx.BaseRepo[*model.Borrower]
	m *model.Borrower
}

func NewBorrowerRepo(db *gorm.DB) *BorrowerRepo {
	return &BorrowerRepo{
		BaseRepo: gormx.NewBaseRepo[*model.Borrower](db),
		m:        new(model.Borrower),
	}
}
