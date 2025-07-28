package repo

import (
	"gorm.io/gorm"

	"zeebee-backend/app/borrower/rpc/internal/repo/model"
	"zeebee-backend/common/gormx"
)

type LenderRepo struct {
	*gormx.BaseRepo[*model.Lender]
	m *model.Lender
}

func NewLenderRepo(db *gorm.DB) *LenderRepo {
	return &LenderRepo{
		BaseRepo: gormx.NewBaseRepo[*model.Lender](db),
		m:        new(model.Lender),
	}
}