package repo

import (
	"context"
	"fmt"

	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"zeebee-backend/app/borrower/rpc/internal/repo/model"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/gormx"
	"zeebee-backend/common/timex"
)

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

func (repo *BorrowerEntityRepo) GetBorrowerEntityById(ctx context.Context, id int64, borrowerId int64) (result *model.BorrowerEntity, err error) {
	err = repo.Client(ctx).
		Where("id = ? and borrower_id = ?", id, borrowerId).
		First(&result).
		Error
	return result, errors.WithStack(err)
}

// Page return page
func (repo *BorrowerEntityRepo) Page(ctx context.Context, pageReq *borrower.BorrowerEntityPageRequestData) (page gormx.PageResponse[*model.BorrowerEntity], err error) {

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

	return page, errors.WithStack(err)

}

// All return all of the borrower's entities
func (repo *BorrowerEntityRepo) All(ctx context.Context, borrowerId int64) (result []*model.BorrowerEntity, err error) {
	err = repo.Client(ctx).
		Where("borrower_id = ?", borrowerId).
		Omit("mailing_address").
		Order("id desc").
		Find(&result).Error
	return result, errors.WithStack(err)
}

// TranToLoanBorrowerWithId get BorrowerEntity from id and trans it to LoanBorrower
func (repo *BorrowerEntityRepo) TranToLoanBorrowerWithId(ctx context.Context, id int64, borrowerId int64) (result *model.LoanBorrower, err error) {
	var entity *model.BorrowerEntity
	entity, err = repo.GetBorrowerEntityById(ctx, id, borrowerId)
	if err != nil {
		return nil, errors.WithStack(err)
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
		return nil, errors.WithStack(err)
	}
	result.Id = 0
	result.BorrowerEntityId = entity.Id

	return result, nil
}
