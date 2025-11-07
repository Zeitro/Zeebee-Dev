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

type BorrowerApplicantRepo struct {
	*gormx.BaseRepo[*model.BorrowerApplicant]
	m *model.BorrowerApplicant
}

func NewBorrowerApplicantRepo(db *gorm.DB) *BorrowerApplicantRepo {
	return &BorrowerApplicantRepo{
		BaseRepo: gormx.NewBaseRepo[*model.BorrowerApplicant](db),
		m:        new(model.BorrowerApplicant),
	}
}

// GetBorrowerApplicantById
func (repo *BorrowerApplicantRepo) GetBorrowerApplicantById(ctx context.Context, id int64, borrowerId int64) (result *model.BorrowerApplicant, err error) {
	err = repo.Client(ctx).
		Where("id = ? and borrower_id = ?", id, borrowerId).
		First(&result).
		Error
	return result, errors.WithStack(err)
}

func (repo *BorrowerApplicantRepo) Page(ctx context.Context, pageReq *borrower.BorrowerApplicantPageRequestData) (page gormx.PageResponse[*model.BorrowerApplicant], err error) {

	filter := pageReq.Filter
	sorter := pageReq.Sorter
	query := repo.Client(ctx).
		Where("borrower_id = ?", filter.BorrowerId)
	if len(filter.Keyword) > 0 {
		query = query.Where("full_name like ?", "%"+filter.Keyword+"%")
	}
	if len(sorter.Column) == 0 {
		sorter.Column = "updated_at"
	}
	if len(sorter.Sort) == 0 {
		sorter.Sort = "desc"
	}
	query = query.Order(fmt.Sprintf("%s %s", sorter.Column, sorter.Sort))
	page, err = gormx.Page[*model.BorrowerApplicant](query, gormx.PageRequest{
		Page:        pageReq.Page,
		Limit:       pageReq.Limit,
		CountColumn: "id",
	})

	return page, errors.WithStack(err)

}

// All return all of the borrower's entities
func (repo *BorrowerApplicantRepo) All(ctx context.Context, borrowerId int64) (result []*model.BorrowerApplicant, err error) {
	err = repo.Client(ctx).
		Where("borrower_id = ?", borrowerId).
		// Omit("address").
		Order("id desc").
		Find(&result).Error
	return result, errors.WithStack(err)
}

// TranToLoanApplicantWithId get BorrowerApplicant from id and trans it to LoanApplicant
func (repo *BorrowerApplicantRepo) TranToLoanApplicantWithId(ctx context.Context, id int64, borrowerId int64) (result *model.LoanApplicant, err error) {
	var applicant *model.BorrowerApplicant
	applicant, err = repo.GetBorrowerApplicantById(ctx, id, borrowerId)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	result = &model.LoanApplicant{}
	err = copier.CopyWithOption(result, applicant, copier.Option{
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
	result.BorrowerApplicantId = applicant.Id

	return result, nil
}
