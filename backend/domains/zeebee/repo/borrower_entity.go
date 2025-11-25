package repo

import (
	"context"
	"fmt"

	"github.com/jinzhu/copier"
	"gorm.io/gorm"

	"zeebee-backend/common/errorx"

	"zeebee-backend/common/copierx"
	"zeebee-backend/common/gormx"
	"zeebee-backend/common/timex"
	"zeebee-backend/domains/zeebee/model"
)

type BorrowerEntityBorrowerPageReq struct {
	Page   int64                                `json:"page"`
	Limit  int64                                `json:"limit"`
	Filter *BorrowerEntityBorrowerPageReqFilter `json:"filter"`
	Sorter *BorrowerEntityBorrowerPageReqSorter `json:"sorter"`
}

type BorrowerEntityBorrowerPageReqFilter struct {
	BorrowerId int64  `json:"borrower_id"`
	Keyword    string `json:"keyword"`
}
type BorrowerEntityBorrowerPageReqSorter struct {
	Column string `json:"column"`
	Sort   string `json:"sort"`
}

type BorrowerEntityRepo struct {
	*gormx.BaseRepo[*model.BorrowerEntity]
	m *model.BorrowerEntity
}

func NewBorrowerEntityRepo(db *gorm.DB) *BorrowerEntityRepo {
	return &BorrowerEntityRepo{
		BaseRepo: gormx.NewBaseRepo[*model.BorrowerEntity](db),
		m:        new(model.BorrowerEntity),
	}
}

// GetBorrowerEntityById return BorrowerEntity by id and borrowerId
func (repo *BorrowerEntityRepo) GetBorrowerEntityById(ctx context.Context, id int64, organizationId int64) (result *model.BorrowerEntity, err error) {
	err = repo.Client(ctx).
		Where("id = ? and organization_id = ?", id, organizationId).
		First(&result).
		Error
	return result, errorx.WithStack(err)
}

// PageForBorrower return page
func (repo *BorrowerEntityRepo) PageForBorrower(ctx context.Context, pageReq *BorrowerEntityBorrowerPageReq) (page gormx.PageResponse[*model.BorrowerEntity], err error) {

	filter := pageReq.Filter
	sorter := pageReq.Sorter
	query := repo.Client(ctx).
		Where("borrower_id = ?", filter.BorrowerId)
	if len(filter.Keyword) > 0 {
		query = query.Omit("mailing_address").Where("name like ?", "%"+filter.Keyword+"%")
	}
	if len(sorter.Column) == 0 {
		sorter.Column = "updated_at"
	}
	if len(sorter.Sort) == 0 {
		sorter.Sort = "desc"
	}
	query = query.Order(fmt.Sprintf("%s %s", sorter.Column, sorter.Sort))
	page, err = gormx.Page[*model.BorrowerEntity](query, gormx.PageRequest{
		Page:        pageReq.Page,
		Limit:       pageReq.Limit,
		CountColumn: "id",
	})

	return page, errorx.WithStack(err)
}

// AllForBorrower return all of the borrower's entities
func (repo *BorrowerEntityRepo) AllForBorrower(ctx context.Context, borrowerId int64) (result []*model.BorrowerEntity, err error) {
	err = repo.Client(ctx).
		Where("borrower_id = ?", borrowerId).
		Order("id desc").
		Find(&result).Error
	return result, errorx.WithStack(err)
}

// TranToLoanBorrowerWithId get BorrowerEntity from id and trans it to LoanBorrower
func (repo *BorrowerEntityRepo) TranToLoanBorrowerWithId(ctx context.Context, id int64, organizationId int64) (result *model.LoanBorrower, err error) {
	var entity *model.BorrowerEntity
	entity, err = repo.GetBorrowerEntityById(ctx, id, organizationId)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	result = &model.LoanBorrower{}
	err = copier.CopyWithOption(result, entity, copier.Option{
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
	result.BorrowerEntityId = entity.Id

	return result, nil
}
