package repo

import (
	"gorm.io/gorm"

	"zeebee-backend/app/borrower/rpc/internal/repo/model"
	"zeebee-backend/common/gormx"
)

type LenderCompanyRepo struct {
	*gormx.BaseRepo[*model.LenderCompany]
	m *model.LenderCompany
}

func NewLenderCompanyRepo(db *gorm.DB) *LenderCompanyRepo {
	return &LenderCompanyRepo{
		BaseRepo: gormx.NewBaseRepo[*model.LenderCompany](db),
		m:        new(model.LenderCompany),
	}
}
