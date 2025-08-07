package repo

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pkg/errors"
	"gorm.io/gorm"

	"zeebee-backend/app/borrower/rpc/internal/repo/model"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/gormx"
)

type PropertyRepo struct {
	*gormx.BaseRepo[*model.Property]
	m *model.Property
}

func NewPropertyRepo(db *gorm.DB) *PropertyRepo {
	return &PropertyRepo{
		BaseRepo: gormx.NewBaseRepo[*model.Property](db),
		m:        new(model.Property),
	}
}

func (repo *PropertyRepo) GetOneWithTx(ctx context.Context, tx *gorm.DB, id int64, preloads []string) (result *model.Property, err error) {
	query := tx.WithContext(ctx).Where("id = ?", id)
	for _, preload := range preloads {
		query = query.Preload(preload)
	}
	err = query.First(&result).Error
	return
}

func (repo *PropertyRepo) GetOne(ctx context.Context, id int64, preloads []string) (result *model.Property, err error) {
	query := repo.BaseRepo.Client(ctx).Where("id = ?", id)
	for _, preload := range preloads {
		query = query.Preload(preload)
	}
	err = query.First(&result).Error
	return
}

// GetOrCreateWithTx query or create property
func (repo *PropertyRepo) GetOrCreateWithTx(ctx context.Context, tx *gorm.DB, data *model.Property) (result *model.Property, err error) {
	// if not existed,create one
	err = tx.WithContext(ctx).
		Where("borrower_id = ? and property_type = ? and unit = ?", data.BorrowerId, data.PropertyType, data.Unit).
		Where("address->>'state' = ?", data.Address.State).
		Where("address->>'city' = ?", data.Address.City).
		Where("address->>'zip_code' = ?", data.Address.ZipCode).
		Where("address->>'address1' = ?", data.Address.Address1).
		Where("address->>'address2' = ?", data.Address.Address2).
		Find(&result).
		Error
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if result == nil || result.Id == 0 {
		err = tx.WithContext(ctx).Create(data).Error
		if err != nil {
			return nil, errors.WithStack(err)
		}
		result = data
	}
	return
}

func (repo *PropertyRepo) GetOneWithTxByProperty(ctx context.Context, tx *gorm.DB, data *model.Property) (result *model.Property, err error) {
	err = tx.WithContext(ctx).
		Where("borrower_id = ? and property_type = ? and unit = ?", data.BorrowerId, data.PropertyType, data.Unit).
		Where("address->>'state' = ?", data.Address.State).
		Where("address->>'city' = ?", data.Address.City).
		Where("address->>'zip_code' = ?", data.Address.ZipCode).
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

func (repo *PropertyRepo) UpdateOneWithTx(ctx context.Context, tx *gorm.DB, id int64, data *model.Property) (result *model.Property, err error) {
	err = tx.WithContext(ctx).Select("*").
		Omit("id", "created_at").
		Where("id = ?", id).
		Updates(data).Error
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return data, nil
}

// SyncTotalAmountWithTx sync total amount
func (repo *PropertyRepo) SyncTotalAmountWithTx(ctx context.Context, tx *gorm.DB, id int64) (err error) {
	var result *model.Property
	err = tx.WithContext(ctx).Where("id = ?", id).Preload("Loans").Find(&result).Error
	if err != nil {
		return errors.WithStack(err)
	}
	if result == nil || result.Id == 0 {
		return nil
	}
	var totalAmount float64
	for _, loan := range result.Loans {
		totalAmount += loan.TotalAmount
	}
	err = tx.WithContext(ctx).
		Model(new(model.Property)).
		Where("id = ?", id).
		Update("total_amount", totalAmount).Error
	return errors.WithStack(err)
}

// Page
func (repo *PropertyRepo) Page(ctx context.Context, pageReq *borrower.PropertyPageRequestData) (page gormx.PageResponse[*model.Property], err error) {
	propertyTable := model.Property{}.TableName()
	applicantTable := model.Applicant{}.TableName()
	loanRepo := NewLoanRepo(repo.BaseRepo.GetDB())
	applicantRepo := NewApplicantRepo(repo.BaseRepo.GetDB())

	filter := pageReq.Filter
	sorter := pageReq.Sorter
	query := repo.Client(ctx).
		Where(fmt.Sprintf("%s.borrower_id = ?", propertyTable), filter.BorrowerId).
		Preload("Loans").
		Preload("Applicants")
	// keyword
	if len(filter.Keyword) > 0 {
		query = query.Where(
			fmt.Sprintf("%s.address->>'state' like @keyword or %s.address->>'city' like @keyword or %s.address->>'address1' like @keyword", propertyTable, propertyTable, propertyTable), sql.Named("keyword", "%"+filter.Keyword+"%"))
	}
	// applicant
	if len(filter.Applicant) > 0 {
		query = query.Joins(fmt.Sprintf("left join %s on %s.id = %s.loan_id", applicantTable, propertyTable, applicantTable)).
			Where(fmt.Sprintf("%s.name in ?", applicantTable), filter.Applicant)
	}

	sortColumns := []string{}

	if sorter != nil {
		switch sorter.Column {
		case "total_amount":
			sortColumns = append(sortColumns,
				fmt.Sprintf("%s.total_amount %s", propertyTable, sorter.Sort),
			)
		}
	}
	if sorter == nil || len(sortColumns) == 0 {
		sortColumns = append(sortColumns, fmt.Sprintf("%s.total_amount desc", propertyTable))
	}

	for _, column := range sortColumns {
		query = query.Order(column)
	}

	page, err = gormx.Page[*model.Property](query, gormx.PageRequest{
		Page:        pageReq.Page,
		Limit:       pageReq.Limit,
		CountColumn: fmt.Sprintf("%s.id", propertyTable),
	})

	if err != nil {
		return page, errors.WithStack(err)
	}

	for i, property := range page.List {
		loans, err := loanRepo.FindAll(ctx, map[string]interface{}{
			"property_id": property.Id,
		})

		if err != nil {
			return page, errors.WithStack(err)
		}
		page.List[i].Loans = loans

		applicants, err := applicantRepo.FindAll(ctx, map[string]interface{}{
			"property_id": property.Id,
		})
		if err != nil {
			return page, errors.WithStack(err)
		}
		page.List[i].Applicants = applicants
	}

	return page, nil
}
