package repo

import (
	"context"

	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"zeebee-backend/app/lender/rpc/internal/repo/model"
	// "zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/constants"
	"zeebee-backend/common/gormx"
)

type PropertyBaseRepo struct {
	*gormx.BaseRepo[*model.PropertyBase]
	m *model.PropertyBase
}

func NewPropertyBaseRepo(db *gorm.DB) *PropertyBaseRepo {
	return &PropertyBaseRepo{
		BaseRepo: gormx.NewBaseRepo[*model.PropertyBase](db),
		m:        new(model.PropertyBase),
	}
}

func (repo *PropertyBaseRepo) GetOneWithTx(ctx context.Context, tx *gorm.DB, id int64, preloads []string) (result *model.PropertyBase, err error) {
	query := tx.WithContext(ctx).Where("id = ?", id)
	for _, preload := range preloads {
		query = query.Preload(preload)
	}
	err = query.First(&result).Error
	return
}

func (repo *PropertyBaseRepo) GetOne(ctx context.Context, id int64) (result *model.PropertyBase, err error) {
	query := repo.BaseRepo.Client(ctx).Where("id = ?", id).
		Preload("LoanProperties").
		Preload("LoanProperties.Applicants").
		Preload("LoanProperties.Loan")

	err = query.First(&result).Error
	return
}

// GetOrCreateWithTx query or create property
func (repo *PropertyBaseRepo) GetOrCreateWithTx(ctx context.Context, tx *gorm.DB, data *model.PropertyBase) (result *model.PropertyBase, err error) {
	// if not existed,create one
	err = tx.WithContext(ctx).
		Where("borrower_id = ? and unit = ?", data.BorrowerId, data.Unit).
		Where("address->>'state' = ?", data.Address.State).
		Where("address->>'city' = ?", data.Address.City).
		// Where("address->>'zip_code' = ?", data.Address.ZipCode).
		Where("address->>'address1' = ?", data.Address.Address1).
		Where("address->>'address2' = ?", data.Address.Address2).
		Find(&result).
		Error
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if result == nil || result.Id == 0 {
		err = tx.WithContext(ctx).Omit(clause.Associations).Create(data).Error
		if err != nil {
			return nil, errors.WithStack(err)
		}
		result = data
	}
	return
}

func (repo *PropertyBaseRepo) GetOneWithTxByProperty(ctx context.Context, tx *gorm.DB, data *model.PropertyBase) (result *model.PropertyBase, err error) {
	err = tx.WithContext(ctx).
		Where("borrower_id = ? and unit = ?", data.BorrowerId, data.Unit).
		Where("address->>'state' = ?", data.Address.State).
		Where("address->>'city' = ?", data.Address.City).
		// Where("address->>'zip_code' = ?", data.Address.ZipCode).
		Where("address->>'address1' = ?", data.Address.Address1).
		Where("address->>'address2' = ?", data.Address.Address2).
		Find(&result).
		Error
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if result == nil || result.Id == 0 {
		return nil, nil
	}
	return
}

func (repo *PropertyBaseRepo) UpdateOneWithTx(ctx context.Context, tx *gorm.DB, id int64, data *model.PropertyBase) (result *model.PropertyBase, err error) {
	err = tx.WithContext(ctx).Select("*").
		Omit("id", "created_at", "applicants").
		Where("id = ?", id).
		Updates(data).Error
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return data, nil
}

// SyncTotalAmountWithTx sync total amount
func (repo *PropertyBaseRepo) SyncTotalAmountWithTx(ctx context.Context, tx *gorm.DB, id int64) (err error) {
	var result *model.PropertyBase
	err = tx.WithContext(ctx).Where("id = ?", id).
		Preload("LoanProperties").
		Preload("LoanProperties.Loan").
		Find(&result).Error
	if err != nil {
		return errors.WithStack(err)
	}
	if result == nil || result.Id == 0 {
		return nil
	}
	var totalAmount float64
	for _, loanProperty := range result.LoanProperties {
		if loanProperty.Loan == nil {
			continue
		}
		if loanProperty.Loan.Stage != constants.LoanStageActive {
			totalAmount += loanProperty.Loan.TotalAmount
		}
	}
	err = tx.WithContext(ctx).
		Model(new(model.PropertyBase)).
		Where("id = ?", id).
		Update("total_amount", totalAmount).Error
	return errors.WithStack(err)
}
