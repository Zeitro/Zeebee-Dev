package repo

import (
	"gorm.io/gorm"

	"zeebee-backend/app/borrower/rpc/internal/repo/model"
	"zeebee-backend/common/gormx"
)

type UsStateRepo struct {
	*gormx.BaseRepo[*model.UsState]
	m *model.UsState
}

func NewUsStateRepo(db *gorm.DB) *UsStateRepo {
	return &UsStateRepo{
		BaseRepo: gormx.NewBaseRepo[*model.UsState](db),
		m:        new(model.UsState),
	}
}
