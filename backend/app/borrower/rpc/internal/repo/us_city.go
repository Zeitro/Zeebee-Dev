package repo

import (
	"gorm.io/gorm"

	"zeebee-backend/app/borrower/rpc/internal/repo/model"
	"zeebee-backend/common/gormx"
)

type UsCityRepo struct {
	*gormx.BaseRepo[*model.UsCity]
	m *model.UsCity
}

func NewUsCityRepo(db *gorm.DB) *UsCityRepo {
	return &UsCityRepo{
		BaseRepo: gormx.NewBaseRepo[*model.UsCity](db),
		m:        new(model.UsCity),
	}
}
