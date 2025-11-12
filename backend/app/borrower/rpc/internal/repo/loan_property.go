package repo

import (
	"context"

	"github.com/pkg/errors"
	"gorm.io/gorm"

	"zeebee-backend/app/borrower/rpc/internal/repo/model"
	"zeebee-backend/common/gormx"
)

type LoanPropertyRepo struct {
	*gormx.BaseRepo[*model.LoanProperty]
	m *model.LoanProperty
}

func NewLoanPropertyRepo(db *gorm.DB) *LoanPropertyRepo {
	return &LoanPropertyRepo{
		BaseRepo: gormx.NewBaseRepo[*model.LoanProperty](db),
		m:        new(model.LoanProperty),
	}
}

// GetOneById get one by id
func (repo *LoanPropertyRepo) GetOneById(ctx context.Context, id int64) (result *model.LoanProperty, err error) {
	err = repo.Client(ctx).Where("id = ?", id).
		Preload("Applicants").
		Preload("Loan").
		First(&result).Error
	return result, errors.WithStack(err)
}

// FindOneByLoanId find one by loan id
func (repo *LoanPropertyRepo) FindOneByLoanId(ctx context.Context, loanId int64) (result *model.LoanProperty, err error) {
	err = repo.Client(ctx).
		Where("loan_id = ?", loanId).
		First(&result).Error
	return result, errors.WithStack(err)
}

// GetLoanPropertyCountByPropertyBaseId get loan property count by property_base_id
func (repo *LoanPropertyRepo) GetLoanPropertyCountByBorrowerPropertyId(ctx context.Context, borrowerPropertyId int64) (count int64, err error) {
	err = repo.Client(ctx).
		Where("borrower_property_id = ?", borrowerPropertyId).
		Count(&count).Error
	return count, errors.WithStack(err)
}

// UpdateOneWithTx
func (repo *LoanPropertyRepo) UpdateOneWithTx(ctx context.Context, tx *gorm.DB, id int64, data *model.LoanProperty) (result *model.LoanProperty, err error) {
	err = tx.WithContext(ctx).Select("*").Omit("id", "created_at", "loan_id").Where("id = ?", id).Updates(data).Error
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return data, nil
}

func (repo *LoanPropertyRepo) DeleteByLoanIdWithTx(ctx context.Context, tx *gorm.DB, loanId int64) (err error) {
	err = tx.WithContext(ctx).
		Model(&model.LoanProperty{}).
		Where("loan_id = ?", loanId).
		Delete(&model.LoanProperty{}).
		Error
	return errors.WithStack(err)
}
