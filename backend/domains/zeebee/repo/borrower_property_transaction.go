package repo

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"gorm.io/gorm"

	"zeebee-backend/common/errorx"
	"zeebee-backend/common/gormx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
)

// BorrowerPropertyTransactionPageRow maps the JOIN query result (transaction + address)
type BorrowerPropertyTransactionPageRow struct {
	Id                 int64                                       `gorm:"column:id"`
	BorrowerPropertyId int64                                       `gorm:"column:borrower_property_id"`
	OrganizationId     int64                                       `gorm:"column:organization_id"`
	Code               string                                      `gorm:"column:code"`
	Type               constants.BorrowerPropertyTransactionType   `gorm:"column:type"`
	Status             constants.BorrowerPropertyTransactionStatus `gorm:"column:status"`
	CoeDate            *sql.NullTime                               `gorm:"column:coe_date"`
	PurchaseAmount     float64                                     `gorm:"column:purchase_amount"`
	EmdAmount          float64                                     `gorm:"column:emd_amount"`
	PropertyType       constants.PropertyType                      `gorm:"column:property_type"`
	Pm                 string                                      `gorm:"column:pm"`
	Strategy           string                                      `gorm:"column:strategy"`
	Note               string                                      `gorm:"column:note"`
	CreatedAt          time.Time                                   `gorm:"column:created_at"`
	UpdatedAt          time.Time                                   `gorm:"column:updated_at"`
	Address            dto.AddressDto                              `gorm:"column:address;;serializer:json"`
}

type BorrowerPropertyTransactionRepo struct {
	*gormx.BaseRepo[*model.BorrowerPropertyTransaction]
}

func NewBorrowerPropertyTransactionRepo(db *gorm.DB) *BorrowerPropertyTransactionRepo {
	return &BorrowerPropertyTransactionRepo{
		BaseRepo: gormx.NewBaseRepo[*model.BorrowerPropertyTransaction](db),
	}
}

func (repo *BorrowerPropertyTransactionRepo) GetOne(ctx context.Context, id int64) (result *model.BorrowerPropertyTransaction, err error) {
	err = repo.Client(ctx).
		Where("id = ?", id).
		Preload("BorrowerProperty").
		First(&result).Error
	return result, errorx.WithStack(err)
}

func (repo *BorrowerPropertyTransactionRepo) GetOrganizationTransaction(ctx context.Context, id int64, organizationId int64) (result *model.BorrowerPropertyTransaction, err error) {
	err = repo.Client(ctx).
		Where("id = ? and organization_id = ?", id, organizationId).
		Preload("BorrowerProperty").
		First(&result).Error
	return result, errorx.WithStack(err)
}

func (repo *BorrowerPropertyTransactionRepo) GetByCode(ctx context.Context, code string) (result *model.BorrowerPropertyTransaction, err error) {
	err = repo.Client(ctx).
		Where("code = ?", code).
		First(&result).Error
	return result, errorx.WithStack(err)
}

func (repo *BorrowerPropertyTransactionRepo) AllByOrganization(ctx context.Context, organizationId int64) (result []*model.BorrowerPropertyTransaction, err error) {
	err = repo.Client(ctx).
		Where("organization_id = ? and status in ?", organizationId, []constants.BorrowerPropertyTransactionStatus{constants.BorrowerPropertyTransactionStatusPending, constants.BorrowerPropertyTransactionStatusClosed}).
		Preload("BorrowerProperty").
		Order("id desc").
		Find(&result).Error
	return result, errorx.WithStack(err)
}

func (repo *BorrowerPropertyTransactionRepo) ListByBorrowerPropertyId(ctx context.Context, borrowerPropertyId int64, organizationId int64) (result []*model.BorrowerPropertyTransaction, err error) {
	err = repo.Client(ctx).
		Where("borrower_property_id = ? and organization_id = ?", borrowerPropertyId, organizationId).
		Preload("BorrowerProperty").
		Order("created_at desc").
		Find(&result).Error
	return result, errorx.WithStack(err)
}

func (repo *BorrowerPropertyTransactionRepo) FindByBorrowerPropertyIdWithTx(ctx context.Context, tx *gorm.DB, borrowerPropertyId int64) (result []*model.BorrowerPropertyTransaction, err error) {
	err = tx.WithContext(ctx).
		Where("borrower_property_id = ? and is_del = 0", borrowerPropertyId).
		Find(&result).Error
	return result, errorx.WithStack(err)
}

func (repo *BorrowerPropertyTransactionRepo) UpdateOne(ctx context.Context, id int64, data *model.BorrowerPropertyTransaction) (result *model.BorrowerPropertyTransaction, err error) {
	err = repo.Client(ctx).
		Model(new(model.BorrowerPropertyTransaction)).
		Where("id = ?", id).
		Updates(data).Error
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	return data, nil
}

// IsBoundToLoan check if the transaction is bound to a loan
func (repo *BorrowerPropertyTransactionRepo) IsBoundToLoan(ctx context.Context, id int64) (bool, error) {
	var count int64
	err := repo.Client(ctx).
		Model(model.Loan{}).
		Where("borrower_property_transaction_id = ?", id).
		Count(&count).Error
	if err != nil {
		return false, errorx.WithStack(err)
	}

	return count > 0, nil
}

// Page get borrower property transactions
func (repo *BorrowerPropertyTransactionRepo) Page(ctx context.Context, pageReq dto.BorrowerPropertyTransactionPageRequestDto) (page gormx.PageResponse[BorrowerPropertyTransactionPageRow], err error) {
	transTable := model.BorrowerPropertyTransaction{}.TableName()
	propTable := model.BorrowerProperty{}.TableName()

	filter := pageReq.Filter
	query := repo.Client(ctx).
		Select(fmt.Sprintf("%s.*, %s.address", transTable, propTable)).
		Joins(fmt.Sprintf("LEFT JOIN %s ON %s.borrower_property_id = %s.id", propTable, transTable, propTable)).
		Where(fmt.Sprintf("%s.organization_id = ?", transTable), filter.OrganizationId).
		Where(fmt.Sprintf("%s.is_del = 0 and %s.is_del = 0", transTable, propTable))

	// keyword: fuzzy match on property address (address1, city, state)
	if len(filter.Keyword) > 0 {
		query = query.Where(fmt.Sprintf(
			"%s.code LIKE @keyword OR %s.address->>'address1' LIKE @keyword OR %s.address->>'city' LIKE @keyword OR %s.address->>'state' LIKE @keyword",
			transTable, propTable, propTable, propTable,
		), sql.Named("keyword", "%"+filter.Keyword+"%"))
	}

	// type: exact match
	if len(filter.Type) > 0 {
		query = query.Where(fmt.Sprintf("%s.type = ?", transTable), filter.Type)
	}

	// status: exact match
	if len(filter.Status) > 0 {
		query = query.Where(fmt.Sprintf("%s.status = ?", transTable), filter.Status)
	}

	// sort
	var sortColumns []string
	if len(pageReq.Sorter.Column) > 0 {
		switch pageReq.Sorter.Column {
		case "created_at", "updated_at":
			sortColumns = append(sortColumns, fmt.Sprintf("%s.%s %s", transTable, pageReq.Sorter.Column, pageReq.Sorter.Sort))
		}
	}
	sortColumns = append(sortColumns, fmt.Sprintf("%s.created_at DESC", transTable))
	for _, col := range sortColumns {
		query = query.Order(col)
	}

	page, err = gormx.Page[BorrowerPropertyTransactionPageRow](query, gormx.PageRequest{
		Page:        pageReq.Page,
		Limit:       pageReq.Limit,
		CountColumn: fmt.Sprintf("%s.id", transTable),
	})
	if err != nil {
		return page, errorx.WithStack(err)
	}

	return page, nil
}
