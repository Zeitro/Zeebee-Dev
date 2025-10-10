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
		Preload("LoanProperty").
		Preload("LoanProperty.PropertyBase").
		Preload("LoanBorrower").
		Preload("LoanExperiences").
		Preload("LoanThirdServices").
		Preload("LoanDemographic").
		// Preload("ApplyFormFile").
		Preload("LoanDocuments.Files").
		First(&result).Error
	return result, errors.WithStack(err)
}

func (repo *LoanRepo) GetBorrowerLoanWithTx(ctx context.Context, tx *gorm.DB, id int64, borrowerId int64) (result *model.Loan, err error) {
	err = tx.WithContext(ctx).
		Where("id = ? and borrower_id = ?", id, borrowerId).
		Preload("Applicants").
		Preload("LoanProperty").
		Preload("LoanProperty.PropertyBase").
		Preload("LoanBorrower").
		Preload("LoanExperiences").
		Preload("LoanThirdServices").
		Preload("LoanDemographic").
		// Preload("ApplyFormFile").
		Preload("LoanDocuments.Files").
		First(&result).Error
	return result, errors.WithStack(err)
}

// UpdateOneWithTx update with tx
func (repo *LoanRepo) UpdateOneWithTx(ctx context.Context, tx *gorm.DB, id int64, data *model.Loan, isApply bool) (result *model.Loan, err error) {
	omitFields := []string{"id", "source", "created_at", "code", "borrower_id"}
	var applyFields []string
	if !isApply {
		applyFields = []string{"existing_lender_company", "existing_rate", "existing_maturity_date", "has_assignment_fee", "assignment_fee_amount", "has_subordinate_debt", "subordinate_debt_amount", "scheduled_close_escrow_date", "loan_term", "arv", "ltv", "liquidity_amount", "budget", "mid_construction", "est_credit_score", "est_yearly_rental_income", "dscr", "reserves_amount", "has_delayed_financing", "has_collateralized_blanket", "has_properties_vacant", "has_short_term_rental"}
	} else {
		applyFields = []string{"original_contract_date", "extension_date", "loan_processor", "extend"}
	}
	omitFields = append(omitFields, applyFields...)
	err = tx.WithContext(ctx).
		Select("*").
		Omit(omitFields...).
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
	propertyBaseTable := model.PropertyBase{}.TableName()
	loanPropertyTable := model.LoanProperty{}.TableName()
	applicantTable := model.Applicant{}.TableName()
	// propertyBaseRepo := NewPropertyBaseRepo(repo.BaseRepo.GetDB())
	loanPropertyRepo := NewLoanPropertyRepo(repo.BaseRepo.GetDB())
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
			query = query.Joins(fmt.Sprintf("left join %s on %s.id = %s.loan_id", loanPropertyTable, loanTable, loanPropertyTable)).
				Joins(fmt.Sprintf("left join %s on %s.property_base_id = %s.id", propertyBaseTable, loanPropertyTable, propertyBaseTable)).
				Where(fmt.Sprintf("%s.address->>'address1' like ? or %s.address->>'address2' like ?", propertyBaseTable, propertyBaseTable), "%"+filter.Keyword+"%", "%"+filter.Keyword+"%")
		case "city":
			hasJoinProperty = true
			query = query.Joins(fmt.Sprintf("left join %s on %s.id = %s.loan_id", loanPropertyTable, loanTable, loanPropertyTable)).
				Joins(fmt.Sprintf("left join %s on %s.property_base_id = %s.id", propertyBaseTable, loanPropertyTable, propertyBaseTable)).
				Where(fmt.Sprintf("%s.address->>'city' like ?", propertyBaseTable), "%"+filter.Keyword+"%")
		case "state":
			hasJoinProperty = true
			query = query.Joins(fmt.Sprintf("left join %s on %s.id = %s.loan_id", loanPropertyTable, loanTable, loanPropertyTable)).
				Joins(fmt.Sprintf("left join %s on %s.property_base_id = %s.id", propertyBaseTable, loanPropertyTable, propertyBaseTable)).
				Where(fmt.Sprintf("%s.address->>'state' like ?", propertyBaseTable), "%"+filter.Keyword+"%")
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
				query = query.Joins(fmt.Sprintf("left join %s on %s.id = %s.loan_id", loanPropertyTable, loanTable, loanPropertyTable)).
					Joins(fmt.Sprintf("left join %s on %s.property_base_id = %s.id", propertyBaseTable, loanPropertyTable, propertyBaseTable))
			}
			sortColumns = append(sortColumns,
				fmt.Sprintf("%s.address->>'address1' %s", propertyBaseTable, sorter.Sort),
				fmt.Sprintf("%s.address->>'city' %s", propertyBaseTable, sorter.Sort),
				fmt.Sprintf("%s.address->>'state' %s", propertyBaseTable, sorter.Sort),
			)
		default:
			sortColumns = append(sortColumns, fmt.Sprintf("%s.%s %s", loanTable, sorter.Column, sorter.Sort))
		}
	}
	if sorter == nil || len(sortColumns) == 0 {
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
		loanProperty, _ := loanPropertyRepo.FindOneByLoanId(ctx, loan.Id)
		page.List[i].LoanProperty = loanProperty
		applicants, _ := applicantRepo.FindAll(ctx, map[string]interface{}{
			"loan_id": loan.Id,
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

// GetByPropertyBaseId get loans by property base id
func (repo *LoanRepo) GetByPropertyBaseId(ctx context.Context, propertyBaseId int64) (result []*model.Loan, err error) {
	loanPropertyTable := model.LoanProperty{}.TableName()
	loanTable := model.Loan{}.TableName()
	propertyBaseTable := model.PropertyBase{}.TableName()
	err = repo.BaseRepo.Client(ctx).
		Joins(fmt.Sprintf("left join %s on %s.id = %s.loan_id", loanPropertyTable, loanTable, loanPropertyTable)).
		Joins(fmt.Sprintf("left join %s on %s.property_base_id = %s.id", propertyBaseTable, loanPropertyTable, propertyBaseTable)).
		Where(fmt.Sprintf("%s.id = ?", propertyBaseTable), propertyBaseId).
		Find(&result).Error
	return result, errors.WithStack(err)
}
