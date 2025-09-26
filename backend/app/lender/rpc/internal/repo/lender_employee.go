package repo

import (
	"gorm.io/gorm"

	"zeebee-backend/app/lender/rpc/internal/repo/model"
	"zeebee-backend/common/gormx"
)

type LenderEmployeeRepo struct {
	*gormx.BaseRepo[*model.LenderEmployee]
	m *model.LenderEmployee
}

func NewLenderEmployeeRepo(db *gorm.DB) *LenderEmployeeRepo {
	return &LenderEmployeeRepo{
		BaseRepo: gormx.NewBaseRepo[*model.LenderEmployee](db),
		m:        new(model.LenderEmployee),
	}
}
