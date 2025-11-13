package repo

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"zeebee-backend/common/constants"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/gormx"
	"zeebee-backend/common/timex"
	"zeebee-backend/common/types"
	"zeebee-backend/domains/zeebee/model"
)

type PropertyBorrowerPageReq struct {
	Page   int64                          `json:"page"`
	Limit  int64                          `json:"limit"`
	Filter *PropertyBorrowerPageReqFilter `json:"filter"`
	Sorter *PropertyBorrowerPageReqSorter `json:"sorter"`
}

type PropertyBorrowerPageReqFilter struct {
	BorrowerId int64    `json:"borrower_id"`
	Keyword    string   `json:"keyword,omitempty"`
	Applicant  []string `json:"applicant,omitempty"`
}

type PropertyBorrowerPageReqSorter struct {
	Column string `json:"column,omitempty"`
	Sort   string `json:"sort,omitempty"`
}

type BorrowerPropertyRepo struct {
	*gormx.BaseRepo[*model.BorrowerProperty]
	m *model.BorrowerProperty
}

func NewBorrowerPropertyRepo(db *gorm.DB) *BorrowerPropertyRepo {
	return &BorrowerPropertyRepo{
		BaseRepo: gormx.NewBaseRepo[*model.BorrowerProperty](db),
		m:        new(model.BorrowerProperty),
	}
}

// GetOneWithTx return one property
func (repo *BorrowerPropertyRepo) GetOneWithTx(ctx context.Context, tx *gorm.DB, id int64, preloads []string) (result *model.BorrowerProperty, err error) {
	query := tx.WithContext(ctx).Where("id = ?", id)
	for _, preload := range preloads {
		query = query.Preload(preload)
	}
	err = query.First(&result).Error
	return
}

// GetOne return one property
func (repo *BorrowerPropertyRepo) GetOne(ctx context.Context, id int64) (result *model.BorrowerProperty, err error) {
	query := repo.BaseRepo.Client(ctx).Where("id = ?", id).
		Preload("LoanProperties").
		Preload("LoanProperties.Applicants").
		Preload("LoanProperties.Loan", func(db *gorm.DB) *gorm.DB {
			return db.Order("loans.updated_at desc")
		})

	err = query.First(&result).Error
	return
}

// GetOrCreateWithTxForBorrower query or create property
func (repo *BorrowerPropertyRepo) GetOrCreateWithTxForBorrower(ctx context.Context, tx *gorm.DB, data *model.BorrowerProperty) (result *model.BorrowerProperty, err error) {
	// if not existed,create one
	err = tx.WithContext(ctx).
		Where("borrower_id = ? and unit = ?", data.BorrowerId, data.Unit).
		Where("address->>'state' = ?", data.Address.State).
		Where("address->>'city' = ?", data.Address.City).
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

func (repo *BorrowerPropertyRepo) GetOrCreateBorrowerPropertyWithTxForBorrower(ctx context.Context, tx *gorm.DB, borrowerId int64, address model.Address) (result *model.BorrowerProperty, err error) {
	if len(address.City) == 0 || len(address.Address1) == 0 || len(address.State) == 0 {
		return nil, fmt.Errorf("Invalid Borrower Address")
	}
	// if not existed,create one
	err = tx.WithContext(ctx).
		Where("borrower_id = ? ", borrowerId).
		Where("address->>'state' = ?", address.State).
		Where("address->>'city' = ?", address.City).
		Where("address->>'address1' = ?", address.Address1).
		Find(&result).
		Error
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if result == nil || result.Id == 0 {
		data := &model.BorrowerProperty{
			BorrowerId: borrowerId,
			Address:    address,
		}
		err = tx.WithContext(ctx).Omit(clause.Associations).Create(data).Error
		if err != nil {
			return nil, errors.WithStack(err)
		}
		result = data
	}
	return
}

// GetOneWithTxByPropertyForBorrower return one property
func (repo *BorrowerPropertyRepo) GetOneWithTxByPropertyForBorrower(ctx context.Context, tx *gorm.DB, data *model.BorrowerProperty) (result *model.BorrowerProperty, err error) {
	err = tx.WithContext(ctx).
		Where("borrower_id = ? and unit = ?", data.BorrowerId, data.Unit).
		Where("address->>'state' = ?", data.Address.State).
		Where("address->>'city' = ?", data.Address.City).
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

// GetOneByIdForBorrower return one property
func (repo *BorrowerPropertyRepo) GetOneByIdForBorrower(ctx context.Context, id int64, borrowerId int64) (result *model.BorrowerProperty, err error) {
	err = repo.Client(ctx).
		Where("id = ? and borrower_id = ?", id, borrowerId).
		First(&result).
		Error
	return result, errors.WithStack(err)
}

// UpdateOneWithTx update one property
func (repo *BorrowerPropertyRepo) UpdateOneWithTx(ctx context.Context, tx *gorm.DB, id int64, data *model.BorrowerProperty) (result *model.BorrowerProperty, err error) {
	err = tx.WithContext(ctx).Select("*").
		Omit("id", "created_at", "applicants").
		Where("id = ?", id).
		Updates(data).Error
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return data, nil
}

// SyncTotalAmountAndLoanCountWithTx sync total amount and loan count
func (repo *BorrowerPropertyRepo) SyncTotalAmountAndLoanCountWithTx(ctx context.Context, tx *gorm.DB, id int64) (err error) {
	var result *model.BorrowerProperty
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
	var loanCount int
	for _, loanProperty := range result.LoanProperties {
		if loanProperty.Loan == nil {
			continue
		}
		loanCount++
		if loanProperty.Loan.Stage != constants.LoanStageActive {
			totalAmount += loanProperty.Loan.TotalAmount
		}
	}

	err = tx.WithContext(ctx).
		Model(new(model.BorrowerProperty)).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"total_amount": totalAmount,
			"loan_count":   loanCount,
		}).Error
	// Updates("total_amount", totalAmount).Error
	return errors.WithStack(err)
}

// Page page
func (repo *BorrowerPropertyRepo) PageForBorrower(ctx context.Context, pageReq *PropertyBorrowerPageReq) (page gormx.PageResponse[*model.BorrowerProperty], err error) {
	borrowerPropertyTable := model.BorrowerProperty{}.TableName()
	loanPropertyTable := model.LoanProperty{}.TableName()
	applicantTable := model.LoanApplicant{}.TableName()
	loanRepo := NewLoanRepo(repo.BaseRepo.GetDB())

	filter := pageReq.Filter
	sorter := pageReq.Sorter
	query := repo.Client(ctx).
		Where(fmt.Sprintf("%s.borrower_id = ?", borrowerPropertyTable), filter.BorrowerId)
	// keyword
	if len(filter.Keyword) > 0 {
		query = query.Where(
			fmt.Sprintf("%s.address->>'state' like @keyword or %s.address->>'city' like @keyword or %s.address->>'address1' like @keyword", borrowerPropertyTable, borrowerPropertyTable, borrowerPropertyTable), sql.Named("keyword", "%"+filter.Keyword+"%"))
	}
	// applicant
	if len(filter.Applicant) > 0 {
		query = query.Joins(fmt.Sprintf("left join %s on %s.id = %s.borrower_property_id", loanPropertyTable, borrowerPropertyTable, loanPropertyTable)).
			Joins(fmt.Sprintf("left join %s on %s.id = %s.loan_property_id", applicantTable, loanPropertyTable, applicantTable)).
			Where(fmt.Sprintf("%s.full_name in ?", applicantTable), filter.Applicant)
	}

	var sortColumns []string

	if sorter != nil {
		switch sorter.Column {
		case "total_amount":
			sortColumns = append(sortColumns,
				fmt.Sprintf("%s.total_amount %s", borrowerPropertyTable, sorter.Sort),
			)
		}
	}
	sortColumns = append(sortColumns, fmt.Sprintf("%s.updated_at desc", borrowerPropertyTable))

	for _, column := range sortColumns {
		query = query.Order(column)
	}

	page, err = gormx.Page[*model.BorrowerProperty](query, gormx.PageRequest{
		Page:        pageReq.Page,
		Limit:       pageReq.Limit,
		CountColumn: fmt.Sprintf("%s.id", borrowerPropertyTable),
	})

	if err != nil {
		return page, errorx.WithStack(err)
	}

	for i, borrowerProperty := range page.List {
		loans, err := loanRepo.GetByBorrowerPropertyId(ctx, borrowerProperty.Id)

		if err != nil {
			return page, errors.WithStack(err)
		}
		page.List[i].Loans = loans
	}

	return page, nil
}

// UpdateOne update one
func (repo *BorrowerPropertyRepo) UpdateOne(ctx context.Context, id int64, data *model.BorrowerProperty) (result *model.BorrowerProperty, err error) {
	err = repo.Client(ctx).
		Model(new(model.BorrowerProperty)).
		Where("id = ?", id).
		Updates(data).Error
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return data, nil
}

// AllForBorrower return all of the borrower's properties
func (repo *BorrowerPropertyRepo) AllForBorrower(ctx context.Context, borrowerId int64) (result []*model.BorrowerProperty, err error) {
	err = repo.Client(ctx).
		Where("borrower_id = ?", borrowerId).
		Order("id desc").
		Find(&result).Error
	return result, errors.WithStack(err)
}

// IsExisted check given property base is existed
func (repo *BorrowerPropertyRepo) IsExisted(ctx context.Context, id int64, borrowerId int64, address *types.Address) (isExisted bool, err error) {
	if address == nil {
		err = errors.New("Address Is Empty")
		return false, errors.WithStack(err)
	}
	query := repo.Client(ctx).
		Where("borrower_id = ?", borrowerId).
		Where("address->>'state' = ?", address.State).
		Where("address->>'city' = ?", address.City).
		Where("address->>'address1' = ?", address.Address1)
	if id > 0 {
		query = query.Where("id != ?", id)
	}
	var count int64
	err = query.Count(&count).Error
	if err != nil {
		return false, errors.WithStack(err)
	}

	return count > 0, nil
}

// TranToLoanPropertyWithId get BorrowerProperty from id and trans it to LoanProperty
func (repo *BorrowerPropertyRepo) TranToLoanPropertyWithId(ctx context.Context, id int64, borrowerId int64) (result *model.LoanProperty, err error) {
	var property *model.BorrowerProperty
	property, err = repo.GetOneByIdForBorrower(ctx, id, borrowerId)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	result = &model.LoanProperty{}
	err = copier.CopyWithOption(result, property, copier.Option{
		CaseSensitive: true,
		DeepCopy:      true,
		Converters: []copier.TypeConverter{
			copierx.StrToSqlTimePtr(timex.LayoutUsStyle),
		},
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}
	result.Id = 0
	result.BorrowerPropertyId = property.Id

	return result, nil
}
