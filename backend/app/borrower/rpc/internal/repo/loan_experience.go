package repo

import (
	"gorm.io/gorm"

	"zeebee-backend/app/borrower/rpc/internal/repo/model"
	"zeebee-backend/common/gormx"
)

type LoanExperienceRepo struct {
	*gormx.BaseRepo[*model.LoanExperience]
	m *model.LoanExperience
}

func NewLoanExperienceRepo(db *gorm.DB) *LoanExperienceRepo {
	return &LoanExperienceRepo{
		BaseRepo: gormx.NewBaseRepo[*model.LoanExperience](db),
		m:        new(model.LoanExperience),
	}
}
