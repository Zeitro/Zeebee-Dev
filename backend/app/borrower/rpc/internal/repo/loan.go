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

// UpdateOneWithTx update with tx
func (repo *LoanRepo) UpdateOneWithTx(ctx context.Context, tx *gorm.DB, id int64, data *model.Loan) (result *model.Loan, err error) {
	err = tx.WithContext(ctx).Select("*").Omit("id", "source", "created_at").Where("id = ?", id).Updates(data).Error
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return data, nil
}

// GetLoanCountByPropertyId
func (repo *LoanRepo) GetLoanCountByPropertyId(ctx context.Context, propertyId int64) (count int64, err error) {
	err = repo.Client(ctx).Where("property_id = ?", propertyId).Count(&count).Error
	return count, errors.WithStack(err)
}
