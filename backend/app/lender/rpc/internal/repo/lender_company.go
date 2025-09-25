package repo

import (
	"gorm.io/gorm"

	"zeebee-backend/app/lender/rpc/internal/repo/model"
	"zeebee-backend/common/gormx"
)

type LenderCompanyRepo struct {
	*gormx.BaseRepo[*model.LenderCompany]
	m *model.LenderCompany
}

func NewLenderCompanyRepo(db *gorm.DB) *LenderRepo {
	return &LenderRepo{
		BaseRepo: gormx.NewBaseRepo[*model.LenderCompany](db),
		m:        new(model.LenderCompany),
	}
}
