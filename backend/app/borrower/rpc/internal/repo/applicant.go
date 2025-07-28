package repo

import (
	"gorm.io/gorm"

	"zeebee-backend/app/borrower/rpc/internal/repo/model"
	"zeebee-backend/common/gormx"
)

type ApplicantRepo struct {
	*gormx.BaseRepo[*model.Applicant]
	m *model.Applicant
}

func NewApplicantRepo(db *gorm.DB) *ApplicantRepo {
	return &ApplicantRepo{
		BaseRepo: gormx.NewBaseRepo[*model.Applicant](db),
		m:        new(model.Applicant),
	}
}