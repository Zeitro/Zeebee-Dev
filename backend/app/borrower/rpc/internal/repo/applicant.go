package repo

import (
	"context"

	"github.com/pkg/errors"
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

// UpdateOneWithTx
func (repo *ApplicantRepo) UpdateOneWithTx(ctx context.Context, tx *gorm.DB, id int64, data *model.Applicant) (result *model.Applicant, err error) {
	err = tx.WithContext(ctx).Select("*").Omit("id", "created_at").Where("id = ?", id).Updates(data).Error
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return data, nil
}
