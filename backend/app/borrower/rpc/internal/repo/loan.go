package repo

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"gorm.io/gorm"

	"zeebee-backend/app/borrower/rpc/internal/repo/model"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/gormx"
)

type LoanRepo struct {
	*gormx.BaseRepo[*model.Loan]
	m *model.Loan
}

func NewLoanRepo(db *gorm.DB) *LoanRepo {
	return &LoanRepo{
		BaseRepo: gormx.NewBaseRepo[*model.Loan](db),
		m:        new(model.Loan),
	}
}

func (repo *LoanRepo) GetBorrowerLoan(ctx context.Context, id int64, borrowerId int64) (result *model.Loan, err error) {
	err = repo.BaseRepo.Client(ctx).
		Where("id = ? and borrower_id = ?", id, borrowerId).
		Preload("Applicants").
		Preload("Property").
		First(&result).Error
	return result, errors.WithStack(err)
}

// UpdateOneWithTx update with tx
func (repo *LoanRepo) UpdateOneWithTx(ctx context.Context, tx *gorm.DB, id int64, data *model.Loan) (result *model.Loan, err error) {
	err = tx.WithContext(ctx).
		Select("*").
		Omit("id", "source", "created_at").
		Where("id = ?", id).
		Updates(data).Error
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return data, nil
}

// GetLoanCountByPropertyId
func (repo *LoanRepo) GetLoanCountByPropertyId(ctx context.Context, propertyId int64) (count int64, err error) {
	err = repo.Client(ctx).
		Where("property_id = ?", propertyId).
		Count(&count).Error
	return count, errors.WithStack(err)
}

// page
func (repo *LoanRepo) Page(ctx context.Context, pageReq *borrower.LoanPageRequestData) (page gormx.PageResponse[*model.Loan], err error) {
	loanTable := model.Loan{}.TableName()
	propertyTable := model.Property{}.TableName()
	applicantTable := model.Applicant{}.TableName()
	propertyRepo := NewPropertyRepo(repo.BaseRepo.GetDB())
	applicantRepo := NewApplicantRepo(repo.BaseRepo.GetDB())

	filter := pageReq.Filter
	sorter := pageReq.Sorter
	query := repo.Client(ctx).
		Where(fmt.Sprintf("%s.borrower_id = ?", loanTable), filter.BorrowerId).
		Preload("Property").
		Preload("Applicants")
	if len(filter.Stage) > 0 {
		query = query.Where(fmt.Sprintf("%s.stage = ?", loanTable), filter.Stage)
	}
	hasJoinProperty := false

	// keyword
	if len(filter.Keyword) > 0 {
		switch filter.KeywordType {
		case "address":
			hasJoinProperty = true
			query = query.Joins(fmt.Sprintf("left join %s on %s.property_id = %s.id", propertyTable, loanTable, propertyTable)).
				Where(fmt.Sprintf("%s.address->>'address1' like ? or %s.address->>'address2' like ?", propertyTable, propertyTable), "%"+filter.Keyword+"%", "%"+filter.Keyword+"%")
		case "city":
			hasJoinProperty = true
			query = query.Joins(fmt.Sprintf("left join %s on %s.property_id = %s.id", propertyTable, loanTable, propertyTable)).
				Where(fmt.Sprintf("%s.address->>'city' like ?", propertyTable), "%"+filter.Keyword+"%")
		case "state":
			hasJoinProperty = true
			query = query.Joins(fmt.Sprintf("left join %s on %s.property_id = %s.id", propertyTable, loanTable, propertyTable)).
				Where(fmt.Sprintf("%s.address->>'state' like ?", propertyTable), "%"+filter.Keyword+"%")
		case "borrower":
			query = query.Joins(fmt.Sprintf("left join %s on %s.id = %s.loan_id", applicantTable, loanTable, applicantTable)).
				Where(fmt.Sprintf("%s.full_name like ?", applicantTable), "%"+filter.Keyword+"%")
		}
	}
	// loan purpose
	if len(filter.LoanPurpose) > 0 {
		query = query.Where(fmt.Sprintf("%s.loan_purpose in ?", loanTable), filter.LoanPurpose)
	}
	// loan type
	if len(filter.LoanType) > 0 {
		query = query.Where(fmt.Sprintf("%s.loan_type in ?", loanTable), filter.LoanType)
	}
	// lender
	if len(filter.Lender) > 0 {
		query = query.Where(fmt.Sprintf("%s.lender_company in ?", loanTable), filter.Lender)
	}
	// pending status
	if len(filter.PendingStatus) > 0 {
		query = query.Where(fmt.Sprintf("%s.pending_status in ?", loanTable), filter.PendingStatus)
	}
	sortColumns := []string{}

	// sort
	if sorter != nil {
		switch sorter.Column {
		case "address":
			if !hasJoinProperty {
				query = query.Joins(fmt.Sprintf("left join %s on %s.property_id = %s.id", propertyTable, loanTable, propertyTable))
			}
			sortColumns = append(sortColumns,
				fmt.Sprintf("%s.address->>'address1' %s", propertyTable, sorter.Sort),
				fmt.Sprintf("%s.address->>'city' %s", propertyTable, sorter.Sort),
				fmt.Sprintf("%s.address->>'state' %s", propertyTable, sorter.Sort),
			)
		default:
			sortColumns = append(sortColumns, fmt.Sprintf("%s.%s %s", loanTable, sorter.Column, sorter.Sort))
		}
	}
	if sorter == nil {
		sortColumns = append(sortColumns, fmt.Sprintf("%s.created_at desc", loanTable))
	}
	for _, column := range sortColumns {
		query = query.Order(column)
	}
	page, err = gormx.Page[*model.Loan](query, gormx.PageRequest{
		Page:        pageReq.Page,
		Limit:       pageReq.Limit,
		CountColumn: fmt.Sprintf("%s.id", loanTable),
	})

	if err != nil {
		return page, errors.WithStack(err)
	}

	for i, loan := range page.List {
		property, _ := propertyRepo.FindOneById(ctx, loan.PropertyId)
		page.List[i].Property = property
		applicants, _ := applicantRepo.FindAll(ctx, map[string]interface{}{
			applicantRepo.m.FieldLoanId(): loan.Id,
		})
		page.List[i].Applicants = applicants
	}

	return
}

// GetLenderCompany get lender company
func (repo *LoanRepo) GetLenderCompany(ctx context.Context, borrowerId int64) (result []string, err error) {
	err = repo.BaseRepo.Client(ctx).
		Select("lender_company").
		Where("borrower_id = ?", borrowerId).
		Distinct("lender_company").
		Pluck("lender_company", &result).
		Error
	return result, errors.WithStack(err)
}
