package repo

import (
	"context"
	"fmt"

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

// GetBorrowerNames return borrower names
func (repo *LoanApplicantRepo) GetBorrowerNames(ctx context.Context, borrowerId int64) (result []string, err error) {
	err = repo.BaseRepo.Client(ctx).
		Select("full_name").
		Where("borrower_id = ?", borrowerId).
		Distinct("full_name").
		Pluck("full_name", &result).
		Error
	return result, errors.WithStack(err)
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

// GetByPropertyByBorrowerPropertyId get applicants by borrower property id
func (repo *LoanApplicantRepo) GetByPropertyByBorrowerPropertyId(ctx context.Context, borrowerPropertyId int64) (result []*model.LoanApplicant, err error) {
	loanPropertyTable := model.LoanProperty{}.TableName()
	applicantTable := model.LoanApplicant{}.TableName()
	borrowePropertyTable := model.BorrowerProperty{}.TableName()
	err = repo.BaseRepo.Client(ctx).
		Joins(fmt.Sprintf("left join %s on %s.loan_property_id = %s.id", loanPropertyTable, applicantTable, loanPropertyTable)).
		Joins(fmt.Sprintf("left join %s on %s.borrower_property_id = %s.id", borrowePropertyTable, loanPropertyTable, borrowePropertyTable)).
		Where(fmt.Sprintf("%s.id = ?", borrowePropertyTable), borrowerPropertyId).
		Find(&result).Error
	return result, errors.WithStack(err)
}
