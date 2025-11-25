package repo

import (
	"context"
	"fmt"

	"github.com/jinzhu/copier"
	"gorm.io/gorm"

	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/gormx"
	"zeebee-backend/common/timex"
	"zeebee-backend/domains/zeebee/model"
)

type BorrowerApplicantBorrowerPageReq struct {
	Page   int64                                   `json:"page"`
	Limit  int64                                   `json:"limit"`
	Filter *BorrowerApplicantBorrowerPageReqFilter `json:"filter"`
	Sorter *BorrowerApplicantBorrowerPageReqSorter `json:"sorter"`
}

type BorrowerApplicantBorrowerPageReqFilter struct {
	BorrowerId int64  `json:"borrower_id"`
	Keyword    string `json:"keyword"`
}
type BorrowerApplicantBorrowerPageReqSorter struct {
	Column string `json:"column"`
	Sort   string `json:"sort"`
}

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
func (repo *BorrowerApplicantRepo) GetBorrowerApplicantById(ctx context.Context, id int64, organizationId int64) (result *model.BorrowerApplicant, err error) {
	err = repo.Client(ctx).
		Where("id = ? and organization_id = ?", id, organizationId).
		First(&result).
		Error
	return result, errorx.WithStack(err)
}

func (repo *BorrowerApplicantRepo) PageForBorrower(ctx context.Context, pageReq *BorrowerApplicantBorrowerPageReq) (page gormx.PageResponse[*model.BorrowerApplicant], err error) {

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

	return page, errorx.WithStack(err)

}

// AllForBorrower return all of the borrower's entities
func (repo *BorrowerApplicantRepo) AllForBorrower(ctx context.Context, borrowerId int64) (result []*model.BorrowerApplicant, err error) {
	err = repo.Client(ctx).
		Where("borrower_id = ?", borrowerId).
		// Omit("address").
		Order("id desc").
		Find(&result).Error
	return result, errorx.WithStack(err)
}

// TranToLoanApplicantWithId get BorrowerApplicant from id and trans it to LoanApplicant
func (repo *BorrowerApplicantRepo) TranToLoanApplicantWithId(ctx context.Context, id int64, organizationId int64) (result *model.LoanApplicant, err error) {
	var applicant *model.BorrowerApplicant
	applicant, err = repo.GetBorrowerApplicantById(ctx, id, organizationId)
	if err != nil {
		return nil, errorx.WithStack(err)
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
		return nil, errorx.WithStack(err)
	}
	result.Id = 0
	result.BorrowerApplicantId = applicant.Id

	return result, nil
}
