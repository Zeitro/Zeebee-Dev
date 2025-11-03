package repo

import (
	"gorm.io/gorm"

	"zeebee-backend/app/borrower/rpc/internal/repo/model"
	"zeebee-backend/common/gormx"
)

type BorrowerApplicantRepo struct {
	*gormx.BaseRepo[*model.BorrowerApplicant]
	m *model.BorrowerApplicant
}

func NewBorrowerApplicantRepo(db *gorm.DB) *BorrowerApplicantRepo {
	return &BorrowerApplicantRepo{
		BaseRepo: gormx.NewBaseRepo[*model.BorrowerApplicant](db),
		m:        new(model.BorrowerApplicant),
	}
}
