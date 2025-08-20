package repo

import (
	"gorm.io/gorm"

	"zeebee-backend/app/borrower/rpc/internal/repo/model"
	"zeebee-backend/common/gormx"
)

type LoanDemographicRepo struct {
	*gormx.BaseRepo[*model.LoanDemographic]
	m *model.LoanDemographic
}

func NewLoanDemographicRepo(db *gorm.DB) *LoanDemographicRepo {
	return &LoanDemographicRepo{
		BaseRepo: gormx.NewBaseRepo[*model.LoanDemographic](db),
		m:        new(model.LoanDemographic),
	}
}
