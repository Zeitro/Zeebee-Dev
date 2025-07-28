package repo

import (
	"gorm.io/gorm"

	"zeebee-backend/app/borrower/rpc/internal/repo/model"
	"zeebee-backend/common/gormx"
)

type PropertyRepo struct {
	*gormx.BaseRepo[*model.Property]
	m *model.Property
}

func NewPropertyRepo(db *gorm.DB) *PropertyRepo {
	return &PropertyRepo{
		BaseRepo: gormx.NewBaseRepo[*model.Property](db),
		m:        new(model.Property),
	}
}