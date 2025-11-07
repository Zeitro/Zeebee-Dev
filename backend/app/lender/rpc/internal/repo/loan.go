package repo

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"gorm.io/gorm"

	"zeebee-backend/app/lender/rpc/internal/repo/model"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/constants"

	// "zeebee-backend/app/lender/rpc/pb/lender"
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

func (repo *LoanRepo) GetLenderLoan(ctx context.Context, id int64, lenderCompanyId int64, withRelationships bool) (result *model.Loan, err error) {

	query := repo.BaseRepo.Client(ctx).
		Where("id = ? and lender_company_id = ?", id, lenderCompanyId)
	if withRelationships {
		query = query.Preload("Applicants", func(db *gorm.DB) *gorm.DB {
			return db.Order("loan_applicants.is_primary desc")
		}).
			Preload("LoanProperty").
			// Preload("LoanProperty.BorrowerProperty").
			Preload("LoanBorrower").
			Preload("LoanExperiences", func(db *gorm.DB) *gorm.DB {
				return db.Order("loan_experiences.id asc")
			}).
			Preload("LoanThirdServices", func(db *gorm.DB) *gorm.DB {
				return db.Order("loan_third_services.id asc")
			}).
			Preload("LoanDemographic").
			// Preload("ApplyFormFile").
			Preload("LoanDocuments", func(db *gorm.DB) *gorm.DB {
				return db.Order("loan_documents.id asc")
			}).
			Preload("LoanDocuments.Files", func(db *gorm.DB) *gorm.DB {
				return db.Order("files.id asc")
			})
	}
	err = query.First(&result).Error
	return result, errors.WithStack(err)
}

func (repo *LoanRepo) GetLenderLoanWithTx(ctx context.Context, tx *gorm.DB, id int64, lenderCompanyId int64) (result *model.Loan, err error) {
	err = tx.WithContext(ctx).
		Where("id = ? and lender_company_id = ?", id, lenderCompanyId).
		Preload("Applicants", func(db *gorm.DB) *gorm.DB {
			return db.Order("loan_applicants.is_primary desc")
		}).
		Preload("LoanProperty").
		// Preload("LoanProperty.BorrowerProperty").
		Preload("LoanBorrower").
		Preload("LoanExperiences", func(db *gorm.DB) *gorm.DB {
			return db.Order("loan_experiences.id asc")
		}).
		Preload("LoanThirdServices", func(db *gorm.DB) *gorm.DB {
			return db.Order("loan_third_services.id asc")
		}).
		Preload("LoanDemographic", func(db *gorm.DB) *gorm.DB {
			return db.Order("loan_demographic.id asc")
		}).
		// Preload("ApplyFormFile").
		Preload("LoanDocuments", func(db *gorm.DB) *gorm.DB {
			return db.Order("loan_documents.id asc")
		}).
		Preload("LoanDocuments.Files", func(db *gorm.DB) *gorm.DB {
			return db.Order("files.id asc")
		}).
		First(&result).Error
	return result, errors.WithStack(err)
}

// UpdateOneWithTx update with tx
func (repo *LoanRepo) UpdateOneWithTx(ctx context.Context, tx *gorm.DB, id int64, data *model.Loan) (result *model.Loan, err error) {
	omitFields := []string{"id", "source", "created_at", "code", "borrower_id", "lender_company_id"}
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

// UpdateStageWithTx
func (repo *LoanRepo) UpdateStageWithTx(ctx context.Context, tx *gorm.DB, id int64, stage constants.LoanStage) (err error) {
	err = tx.WithContext(ctx).Model(&model.Loan{}).
		Where("id = ?", id).
		Update("stage", stage).Error
	return errors.WithStack(err)
}

// UpdatePendingStatusWithTx
func (repo *LoanRepo) UpdatePendingStatusWithTx(ctx context.Context, tx *gorm.DB, id int64, pendingStatus constants.LoanPendingStatus, stage constants.LoanStage) (err error) {
	err = tx.WithContext(ctx).Model(&model.Loan{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"stage":          stage,
			"pending_status": pendingStatus,
		}).
		Error
	return errors.WithStack(err)
}

// GetLoanCountByPropertyId
func (repo *LoanRepo) GetLoanCountByPropertyId(ctx context.Context, propertyId int64) (count int64, err error) {
	err = repo.Client(ctx).
		Where("property_id = ?", propertyId).
		Count(&count).Error
	return count, errors.WithStack(err)
}

// page
func (repo *LoanRepo) Page(ctx context.Context, pageReq *lender.LoanPageRequestData) (page gormx.PageResponse[*model.Loan], err error) {
	loanTable := model.Loan{}.TableName()
	// borrowerPropertyTable := model.BorrowerProperty{}.TableName()
	loanPropertyTable := model.LoanProperty{}.TableName()
	applicantTable := model.LoanApplicant{}.TableName()
	// propertyBaseRepo := NewBorrowerPropertyRepo(repo.BaseRepo.GetDB())
	loanPropertyRepo := NewLoanPropertyRepo(repo.BaseRepo.GetDB())
	applicantRepo := NewLoanApplicantRepo(repo.BaseRepo.GetDB())

	filter := pageReq.Filter
	sorter := pageReq.Sorter
	omitFields := []string{
		fmt.Sprintf("%s.extend", loanTable),
		fmt.Sprintf("%s.loan_processor", loanTable),
		fmt.Sprintf("%s.extension_date", loanTable),
		fmt.Sprintf("%s.original_contract_date", loanTable),
		fmt.Sprintf("%s.loan_processor", loanTable),
		fmt.Sprintf("%s.bugdet", loanTable),
	}
	query := repo.Client(ctx).
		Omit(omitFields...).
		Where(fmt.Sprintf("%s.lender_company_id = ?", loanTable), filter.LenderCompanyId).
		Preload("Property").
		Preload("Applicants")
	if len(filter.Stage) > 0 {
		query = query.Where(fmt.Sprintf("%s.stage = ?", loanTable), filter.Stage)
		if filter.Stage == string(constants.LoanStageDraft) {
			query = query.Where(fmt.Sprintf("%s.source in ?", loanTable), []constants.LoanSource{constants.LoanSourcePreApproval, constants.LoanSourceSystem})
		}
	}
	hasJoinProperty := false

	// keyword
	if len(filter.Keyword) > 0 {
		switch filter.KeywordType {
		case "address":
			hasJoinProperty = true
			query = query.Joins(fmt.Sprintf("left join %s on %s.id = %s.loan_id", loanPropertyTable, loanTable, loanPropertyTable)).
				// Joins(fmt.Sprintf("left join %s on %s.property_base_id = %s.id", borrowerPropertyTable, loanPropertyTable, borrowerPropertyTable)).
				Where(fmt.Sprintf("%s.address->>'address1' like ? or %s.address->>'address2' like ?", loanPropertyTable, loanPropertyTable), "%"+filter.Keyword+"%", "%"+filter.Keyword+"%")
		case "city":
			hasJoinProperty = true
			query = query.Joins(fmt.Sprintf("left join %s on %s.id = %s.loan_id", loanPropertyTable, loanTable, loanPropertyTable)).
				// Joins(fmt.Sprintf("left join %s on %s.property_base_id = %s.id", borrowerPropertyTable, loanPropertyTable, borrowerPropertyTable)).
				Where(fmt.Sprintf("%s.address->>'city' like ?", loanPropertyTable), "%"+filter.Keyword+"%")
		case "state":
			hasJoinProperty = true
			query = query.Joins(fmt.Sprintf("left join %s on %s.id = %s.loan_id", loanPropertyTable, loanTable, loanPropertyTable)).
				// Joins(fmt.Sprintf("left join %s on %s.property_base_id = %s.id", borrowerPropertyTable, loanPropertyTable, borrowerPropertyTable)).
				Where(fmt.Sprintf("%s.address->>'state' like ?", loanPropertyTable), "%"+filter.Keyword+"%")
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
				query = query.Joins(fmt.Sprintf("left join %s on %s.id = %s.loan_id", loanPropertyTable, loanTable, loanPropertyTable))
				// Joins(fmt.Sprintf("left join %s on %s.property_base_id = %s.id", borrowerPropertyTable, loanPropertyTable, borrowerPropertyTable))
			}
			sortColumns = append(sortColumns,
				fmt.Sprintf("%s.address->>'address1' %s", loanPropertyTable, sorter.Sort),
				fmt.Sprintf("%s.address->>'city' %s", loanPropertyTable, sorter.Sort),
				fmt.Sprintf("%s.address->>'state' %s", loanPropertyTable, sorter.Sort),
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

// GetByBorrowerPropertyId get loans by property base id
// func (repo *LoanRepo) GetByBorrowerPropertyId(ctx context.Context, propertyBaseId int64) (result []*model.Loan, err error) {
// 	loanPropertyTable := model.LoanProperty{}.TableName()
// 	loanTable := model.Loan{}.TableName()
// 	borrowerPropertyTable := model.BorrowerProperty{}.TableName()
// 	err = repo.BaseRepo.Client(ctx).
// 		Joins(fmt.Sprintf("left join %s on %s.id = %s.loan_id", loanPropertyTable, loanTable, loanPropertyTable)).
// 		Joins(fmt.Sprintf("left join %s on %s.property_base_id = %s.id", borrowerPropertyTable, loanPropertyTable, borrowerPropertyTable)).
// 		Where(fmt.Sprintf("%s.id = ?", borrowerPropertyTable), propertyBaseId).
// 		Find(&result).Error
// 	return result, errors.WithStack(err)
// }
