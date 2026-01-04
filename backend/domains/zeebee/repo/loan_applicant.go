package repo

import (
	"context"

	"github.com/pkg/errors"
	"gorm.io/gorm"

	"zeebee-backend/common/gormx"
	"zeebee-backend/domains/zeebee/model"
)

type LoanApplicantRepo struct {
	*gormx.BaseRepo[*model.LoanApplicant]
	m *model.LoanApplicant
}

func NewLoanApplicantRepo(db *gorm.DB) *LoanApplicantRepo {
	return &LoanApplicantRepo{
		BaseRepo: gormx.NewBaseRepo[*model.LoanApplicant](db),
		m:        new(model.LoanApplicant),
	}
}

// UpdateOneWithTx
func (repo *LoanApplicantRepo) UpdateOneWithTx(ctx context.Context, tx *gorm.DB, id int64, data *model.LoanApplicant) (result *model.LoanApplicant, err error) {
	err = tx.WithContext(ctx).Select("*").Omit("id", "created_at").Where("id = ?", id).Updates(data).Error
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return data, nil
}

// GetApplicantsFullNamesByLoanId return applicants full names by loan id
func (repo *LoanApplicantRepo) GetApplicantsFullNamesByLoanId(ctx context.Context, loanId int64) (result []string, err error) {
	err = repo.BaseRepo.Client(ctx).
		Select("full_name").
		Where("loan_id = ?", loanId).
		Distinct("full_name").
		Pluck("full_name", &result).
		Error
	return result, errors.WithStack(err)
}
