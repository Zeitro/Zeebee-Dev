package repo

import (
	"context"
	"fmt"

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

func (repo *ApplicantRepo) GetBorrowerNames(ctx context.Context, borrowerId int64) (result []string, err error) {
	err = repo.BaseRepo.Client(ctx).
		Select("full_name").
		Where("borrower_id = ?", borrowerId).
		Distinct("full_name").
		Pluck("full_name", &result).
		Error
	return result, errors.WithStack(err)
}

func (repo *ApplicantRepo) GetApplicantsFullNamesByLoanId(ctx context.Context, loanId int64) (result []string, err error) {
	err = repo.BaseRepo.Client(ctx).
		Select("full_name").
		Where("loan_id = ?", loanId).
		Distinct("full_name").
		Pluck("full_name", &result).
		Error
	return result, errors.WithStack(err)
}

// GetByPropertyBaseId get applicants by property base id
func (repo *ApplicantRepo) GetByPropertyBaseId(ctx context.Context, propertyBaseId int64) (result []*model.Applicant, err error) {
	loanPropertyTable := model.LoanProperty{}.TableName()
	applicantTable := model.Applicant{}.TableName()
	propertyBaseTable := model.PropertyBase{}.TableName()
	err = repo.BaseRepo.Client(ctx).
		Joins(fmt.Sprintf("left join %s on %s.loan_property_id = %s.id", loanPropertyTable, applicantTable, loanPropertyTable)).
		Joins(fmt.Sprintf("left join %s on %s.property_base_id = %s.id", propertyBaseTable, loanPropertyTable, propertyBaseTable)).
		Where(fmt.Sprintf("%s.id = ?", propertyBaseTable), propertyBaseId).
		Find(&result).Error
	return result, errors.WithStack(err)
}
