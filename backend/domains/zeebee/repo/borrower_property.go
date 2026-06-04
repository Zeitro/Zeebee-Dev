package repo

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jinzhu/copier"
	"github.com/lib/pq"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/data/dto"

	"zeebee-backend/common/copierx"
	"zeebee-backend/common/gormx"
	"zeebee-backend/common/timex"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/model"
)

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

func (repo *BorrowerPropertyRepo) GetOrCreateBorrowerPropertyWithTxForBorrower(ctx context.Context, tx *gorm.DB, organizationId int64, creatorId int64, creatorName string, address model.Address) (result *model.BorrowerProperty, err error) {
	if len(address.City) == 0 || len(address.Address1) == 0 || len(address.State) == 0 {
		return nil, fmt.Errorf("invalid borrower address")
	}
	// if not existed,create one
	err = tx.WithContext(ctx).
		Where("organization_id = ? ", organizationId).
		Where("address->>'state' = ?", address.State).
		Where("address->>'city' = ?", address.City).
		Where("address->>'address1' = ?", address.Address1).
		Find(&result).
		Error
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	if result == nil || result.Id == 0 {
		data := &model.BorrowerProperty{
			CreatorId:      creatorId,
			CreatorName:    creatorName,
			Address:        address,
			OrganizationId: organizationId,
		}
		err = tx.WithContext(ctx).Omit(clause.Associations).Create(data).Error
		if err != nil {
			return nil, errorx.WithStack(err)
		}
		result = data
	}
	return
}

// GetOneByIdForBorrower return one property
func (repo *BorrowerPropertyRepo) GetOneByIdForBorrower(ctx context.Context, id int64, creatorId int64) (result *model.BorrowerProperty, err error) {
	err = repo.Client(ctx).
		Where("id = ? and creator_id = ?", id, creatorId).
		First(&result).
		Error
	return result, errorx.WithStack(err)
}

// GetOneByIdForOrganization return one property by id and organization id
func (repo *BorrowerPropertyRepo) GetOneByIdForOrganization(ctx context.Context, id int64, organizationId int64) (result *model.BorrowerProperty, err error) {
	err = repo.Client(ctx).
		Where("id = ? and organization_id = ?", id, organizationId).
		First(&result).
		Error
	return result, errorx.WithStack(err)
}

// SyncTotalAmountAndLoanCountWithTx sync total amount and loan count
func (repo *BorrowerPropertyRepo) SyncTotalAmountAndLoanCountWithTx(ctx context.Context, tx *gorm.DB, id int64) (err error) {
	var result *model.BorrowerProperty
	err = tx.WithContext(ctx).Where("id = ?", id).
		Preload("LoanProperties").
		Preload("LoanProperties.Loan").
		Find(&result).Error
	if err != nil {
		return errorx.WithStack(err)
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
	return errorx.WithStack(err)
}

// Page page
func (repo *BorrowerPropertyRepo) PageForBorrower(ctx context.Context, pageReq dto.BorrowerPropertyPageRequestDto) (page gormx.PageResponse[*model.BorrowerProperty], err error) {
	borrowerPropertyTable := model.BorrowerProperty{}.TableName()
	loanPropertyTable := model.LoanProperty{}.TableName()
	applicantTable := model.LoanApplicant{}.TableName()
	loanRepo := NewLoanRepo(repo.GetDB())
	transRepo := NewBorrowerPropertyTransactionRepo(repo.GetDB())

	filter := pageReq.Filter
	sorter := pageReq.Sorter
	query := repo.Client(ctx).
		Where(fmt.Sprintf("%s.organization_id = ?", borrowerPropertyTable), filter.OrganizationId)
	// keyword
	if len(filter.Keyword) > 0 {
		query = query.Where(
			fmt.Sprintf("%s.address->>'state' like @keyword or %s.address->>'city' like @keyword or %s.address->>'address1' like @keyword", borrowerPropertyTable, borrowerPropertyTable, borrowerPropertyTable), sql.Named("keyword", "%"+filter.Keyword+"%"))
	}
	// applicant
	if len(filter.Applicant) > 0 {
		query = query.Distinct(fmt.Sprintf("%s.*", borrowerPropertyTable)).
			Joins(fmt.Sprintf("left join %s on %s.id = %s.borrower_property_id", loanPropertyTable, borrowerPropertyTable, loanPropertyTable)).
			Joins(fmt.Sprintf("left join %s on %s.id = %s.loan_property_id", applicantTable, loanPropertyTable, applicantTable)).
			Where(fmt.Sprintf("%s.full_name in ?", applicantTable), filter.Applicant)
	}
	// status
	if len(filter.Status) > 0 {
		query = query.Where(fmt.Sprintf("%s.status @> ?", borrowerPropertyTable), pq.Array(filter.Status))
	}

	var sortColumns []string

	if len(sorter.Column) > 0 {
		switch sorter.Column {
		case "total_amount":
			sortColumns = append(sortColumns,
				fmt.Sprintf("%s.total_amount %s", borrowerPropertyTable, sorter.Sort),
			)
		case "address1", "state", "city":
			sortColumns = append(sortColumns,
				fmt.Sprintf("%s.address->>'%s' %s", borrowerPropertyTable, sorter.Column, sorter.Sort),
			)
		default:
			sortColumns = append(sortColumns,
				fmt.Sprintf("%s.%s %s", borrowerPropertyTable, sorter.Column, sorter.Sort),
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
			return page, errorx.WithStack(err)
		}
		page.List[i].Loans = loans

		// get transactions
		transactions, err := transRepo.FindAll(ctx, map[string]interface{}{
			"borrower_property_id": borrowerProperty.Id,
		})
		if err != nil {
			return page, errorx.WithStack(err)
		}
		page.List[i].Transactions = transactions

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
		return nil, errorx.WithStack(err)
	}
	return data, nil
}

// AllForOrganization return all property for organization
func (repo *BorrowerPropertyRepo) AllForOrganization(ctx context.Context, organizationId int64) (result []*model.BorrowerProperty, err error) {
	err = repo.Client(ctx).
		Where("organization_id = ?", organizationId).
		Order("id desc").
		Find(&result).Error
	return result, errorx.WithStack(err)
}

// IsExisted check given property base is existed
func (repo *BorrowerPropertyRepo) IsExisted(ctx context.Context, id int64, organizationId int64, address *model.Address) (isExisted bool, err error) {
	if address == nil {
		err = errorx.New("Address Is Empty")
		return false, errorx.WithStack(err)
	}
	query := repo.Client(ctx).
		Where("organization_id = ?", organizationId).
		Where("address->>'state' = ?", address.State).
		Where("address->>'city' = ?", address.City).
		Where("address->>'address1' = ?", address.Address1)
	if id > 0 {
		query = query.Where("id != ?", id)
	}
	var count int64
	err = query.Count(&count).Error
	if err != nil {
		return false, errorx.WithStack(err)
	}

	return count > 0, nil
}

// TranToLoanPropertyWithId get BorrowerProperty from id and trans it to LoanProperty
func (repo *BorrowerPropertyRepo) TranToLoanPropertyWithId(ctx context.Context, id int64, organizationId int64) (result *model.LoanProperty, err error) {
	var property *model.BorrowerProperty
	property, err = repo.GetOneByIdForOrganization(ctx, id, organizationId)
	if err != nil {
		return nil, errorx.WithStack(err)
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
		return nil, errorx.WithStack(err)
	}
	result.Id = 0
	result.BorrowerPropertyId = property.Id

	return result, nil
}
