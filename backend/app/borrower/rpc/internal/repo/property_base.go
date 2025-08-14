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
		err = tx.WithContext(ctx).Create(data).Error
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
		Omit("id", "created_at").
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
		totalAmount += loanProperty.Loan.TotalAmount
	}
	err = tx.WithContext(ctx).
		Model(new(model.PropertyBase)).
		Where("id = ?", id).
		Update("total_amount", totalAmount).Error
	return errors.WithStack(err)
}

// Page page
func (repo *PropertyBaseRepo) Page(ctx context.Context, pageReq *borrower.PropertyPageRequestData) (page gormx.PageResponse[*model.PropertyBase], err error) {
	propertyBaseTable := model.PropertyBase{}.TableName()
	loanPropertyTable := model.LoanProperty{}.TableName()
	applicantTable := model.Applicant{}.TableName()
	loanRepo := NewLoanRepo(repo.BaseRepo.GetDB())
	applicantRepo := NewApplicantRepo(repo.BaseRepo.GetDB())

	filter := pageReq.Filter
	sorter := pageReq.Sorter
	query := repo.Client(ctx).
		Where(fmt.Sprintf("%s.borrower_id = ?", propertyBaseTable), filter.BorrowerId)
	// keyword
	if len(filter.Keyword) > 0 {
		query = query.Where(
			fmt.Sprintf("%s.address->>'state' like @keyword or %s.address->>'city' like @keyword or %s.address->>'address1' like @keyword", propertyBaseTable, propertyBaseTable, propertyBaseTable), sql.Named("keyword", "%"+filter.Keyword+"%"))
	}
	// applicant
	if len(filter.Applicant) > 0 {
		query = query.Joins(fmt.Sprintf("left join %s on %s.id = %s.property_base_id", loanPropertyTable, propertyBaseTable, loanPropertyTable)).
			Joins(fmt.Sprintf("left join %s on %s.id = %s.loan_property_id", applicantTable, loanPropertyTable, applicantTable)).
			Where(fmt.Sprintf("%s.name in ?", applicantTable), filter.Applicant)
	}

	var sortColumns []string

	if sorter != nil {
		switch sorter.Column {
		case "total_amount":
			sortColumns = append(sortColumns,
				fmt.Sprintf("%s.total_amount %s", propertyBaseTable, sorter.Sort),
			)
		}
	}
	if sorter == nil || len(sortColumns) == 0 {
		sortColumns = append(sortColumns, fmt.Sprintf("%s.total_amount desc", propertyBaseTable))
	}

	for _, column := range sortColumns {
		query = query.Order(column)
	}

	page, err = gormx.Page[*model.PropertyBase](query, gormx.PageRequest{
		Page:        pageReq.Page,
		Limit:       pageReq.Limit,
		CountColumn: fmt.Sprintf("%s.id", propertyBaseTable),
	})

	if err != nil {
		return page, errors.WithStack(err)
	}

	for i, propertyBase := range page.List {
		loans, err := loanRepo.GetByPropertyBaseId(ctx, propertyBase.Id)

		if err != nil {
			return page, errors.WithStack(err)
		}
		page.List[i].Loans = loans

		applicants, err := applicantRepo.GetByPropertyBaseId(ctx, propertyBase.Id)
		if err != nil {
			return page, errors.WithStack(err)
		}
		page.List[i].Applicants = applicants
	}

	return page, nil
}
