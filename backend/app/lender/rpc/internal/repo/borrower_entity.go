package repo

import (
	"gorm.io/gorm"

	"zeebee-backend/app/lender/rpc/internal/repo/model"
	"zeebee-backend/common/gormx"
)

type BorrowerEntityRepo struct {
	*gormx.BaseRepo[*model.BorrowerEntity]
	m *model.BorrowerEntity
}

func NewBorrowerEntityRepo(db *gorm.DB) *BorrowerEntityRepo {
	return &BorrowerEntityRepo{
		BaseRepo: gormx.NewBaseRepo[*model.BorrowerEntity](db),
		m:        new(model.BorrowerEntity),
	}
}
