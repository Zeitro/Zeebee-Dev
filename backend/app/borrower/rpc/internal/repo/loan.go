package repo

import (
	"context"

	"github.com/pkg/errors"
	"gorm.io/gorm"

	"zeebee-backend/app/borrower/rpc/internal/repo/model"
	"zeebee-backend/common/gormx"
)

type LoanRepo struct {
	*gormx.BaseRepo[*model.Loan]
	m *model.Loan
}

func NewLoanRepo(db *gorm.DB) *LoanRepo {
	return &LoanRepo{
		BaseRepo: gormx.NewBaseRepo[*model.Loan](db),
		m:        new(model.Loan),
	}
}

func (repo *LoanRepo) GetBorrowerLoan(ctx context.Context, id int64, borrowerId int64) (result *model.Loan, err error) {
	err = repo.BaseRepo.Client(ctx).
		Where("id = ? and borrower_id = ?", id, borrowerId).
		Preload("Applicants").
		Preload("Property").
		First(&result).Error
	return result, errors.WithStack(err)
}
