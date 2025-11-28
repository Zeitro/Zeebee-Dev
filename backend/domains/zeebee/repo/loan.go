package repo

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"zeebee-backend/common/errorx"

	"zeebee-backend/common/constants"
	"zeebee-backend/common/gormx"
	"zeebee-backend/domains/zeebee/model"
)

type LoanBorrowerPageReq struct {
	Page   int64                      `json:"page"`
	Limit  int64                      `json:"limit"`
	Filter *LoanBorrowerPageReqFilter `json:"filter"`
	Sorter *LoanBorrowerPageReqSorter `json:"sorter"`
}

type LoanBorrowerPageReqFilter struct {
	UserId         int64    `json:"user_id,omitempty"`
	OrganizationId int64    `json:"organization_id,omitempty"`
	Stage          string   `json:"stage,omitempty"`
	Keyword        string   `json:"keyword,omitempty"`
	KeywordType    string   `json:"keyword_type,omitempty"`
	LoanPurpose    []string `json:"loan_purpose,omitempty"`
	LoanType       []string `json:"loan_type,omitempty"`
	Lender         []string `json:"lender,omitempty"`
	PendingStatus  []string `json:"pending_status,omitempty"`
}

type LoanBorrowerPageReqSorter struct {
	Column string `json:"column,omitempty"`
	Sort   string `json:"sort,omitempty"`
}

type LoanLenderCompanyPageReq struct {
	Page   int64                           `json:"page,omitempty"`
	Limit  int64                           `son:"limit,omitempty"`
	Filter *LoanLenderCompanyPageReqFilter `json:"filter,omitempty"`
	Sorter *LoanLenderCompanyPageReqSorter `json:"sorter,omitempty"`
}

type LoanLenderCompanyPageReqFilter struct {
	LenderCompanyId int64    `json:"lender_company_id,omitempty"`
	Stage           string   `json:"stage,omitempty"`
	Keyword         string   `json:"keyword,omitempty"`
	KeywordType     string   `json:"keyword_type,omitempty"`
	LoanPurpose     []string `json:"loan_purpose,omitempty"`
	LoanType        []string `json:"loan_type,omitempty"`
	PendingStatus   []string `json:"pending_status,omitempty"`
}
type LoanLenderCompanyPageReqSorter struct {
	Column string `json:"column,omitempty"`
	Sort   string `json:"sort,omitempty"`
}

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

func (repo *LoanRepo) GetBorrowerLoan(ctx context.Context, id int64, organizationId int64) (result *model.Loan, err error) {
	err = repo.BaseRepo.Client(ctx).
		Where("id = ? and borrower_organization_id = ?", id, organizationId).
		Preload("Applicants", func(db *gorm.DB) *gorm.DB {
			return db.Order("loan_applicants.is_primary desc")
		}).
		Preload("LoanProperty").
		Preload("LoanBorrower").
		Preload("LoanExperiences", func(db *gorm.DB) *gorm.DB {
			return db.Order("loan_experiences.id asc")
		}).
		Preload("LoanThirdServices", func(db *gorm.DB) *gorm.DB {
			return db.Order("loan_third_services.id asc")
		}).
		Preload("LoanDemographic").
		Preload("LoanDocuments", func(db *gorm.DB) *gorm.DB {
			return db.Order("loan_documents.id asc")
		}).
		Preload("LoanDocuments.Files", func(db *gorm.DB) *gorm.DB {
			return db.Order("files.id asc")
		}).
		First(&result).Error
	return result, errorx.WithStack(err)
}

func (repo *LoanRepo) GetBorrowerLoanWithTx(ctx context.Context, tx *gorm.DB, id int64, organizationId int64) (result *model.Loan, err error) {
	err = tx.WithContext(ctx).
		Where("id = ? and borrower_organization_id = ?", id, organizationId).
		Preload("Applicants", func(db *gorm.DB) *gorm.DB {
			return db.Order("loan_applicants.is_primary desc")
		}).
		Preload("LoanProperty").
		Preload("LoanBorrower").
		Preload("LoanExperiences", func(db *gorm.DB) *gorm.DB {
			return db.Order("loan_experiences.id asc")
		}).
		Preload("LoanThirdServices", func(db *gorm.DB) *gorm.DB {
			return db.Order("loan_third_services.id asc")
		}).
		Preload("LoanDemographic", func(db *gorm.DB) *gorm.DB {
			return db.Order("loan_demographics.id asc")
		}).
		Preload("LoanDocuments", func(db *gorm.DB) *gorm.DB {
			return db.Order("loan_documents.id asc")
		}).
		Preload("LoanDocuments.Files", func(db *gorm.DB) *gorm.DB {
			return db.Order("files.id asc")
		}).
		First(&result).Error
	return result, errorx.WithStack(err)
}

func (repo *LoanRepo) GetLenderLoan(ctx context.Context, id int64, lenderOrganizationId int64, withRelationships bool) (result *model.Loan, err error) {

	query := repo.BaseRepo.Client(ctx).
		Where("id = ? and lender_organization_id = ?", id, lenderOrganizationId)
	if withRelationships {
		query = query.Preload("Applicants", func(db *gorm.DB) *gorm.DB {
			return db.Order("loan_applicants.is_primary desc")
		}).
			Preload("LoanProperty").
			Preload("LoanBorrower").
			Preload("LoanExperiences", func(db *gorm.DB) *gorm.DB {
				return db.Order("loan_experiences.id asc")
			}).
			Preload("LoanThirdServices", func(db *gorm.DB) *gorm.DB {
				return db.Order("loan_third_services.id asc")
			}).
			Preload("LoanDemographic").
			Preload("LoanDocuments", func(db *gorm.DB) *gorm.DB {
				return db.Order("loan_documents.id asc")
			}).
			Preload("LoanDocuments.Files", func(db *gorm.DB) *gorm.DB {
				return db.Order("files.id asc")
			})
	}
	err = query.First(&result).Error
	return result, errorx.WithStack(err)
}

// UpdateOneWithTx update with tx
func (repo *LoanRepo) UpdateOneWithTx(ctx context.Context, tx *gorm.DB, id int64, data *model.Loan) (result *model.Loan, err error) {
	omitFields := []string{"id", "source", "created_at", "code", "creator_id", "creator_naem"}
	var applyFields []string
	omitFields = append(omitFields, applyFields...)
	err = tx.WithContext(ctx).
		Select("*").
		Omit(omitFields...).
		Where("id = ?", id).
		Updates(data).Error
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	return data, nil
}

// UpdateStageWithTx
func (repo *LoanRepo) UpdateStageWithTx(ctx context.Context, tx *gorm.DB, id int64, stage constants.LoanStage) (err error) {
	err = tx.WithContext(ctx).Model(&model.Loan{}).
		Where("id = ?", id).
		Update("stage", stage).Error
	return errorx.WithStack(err)
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
	return errorx.WithStack(err)
}

// page
func (repo *LoanRepo) PageForBorrower(ctx context.Context, pageReq *LoanBorrowerPageReq) (page gormx.PageResponse[*model.Loan], err error) {
	loanTable := model.Loan{}.TableName()
	loanPropertyTable := model.LoanProperty{}.TableName()
	applicantTable := model.LoanApplicant{}.TableName()
	loanPropertyRepo := NewLoanPropertyRepo(repo.BaseRepo.GetDB())
	applicantRepo := NewLoanApplicantRepo(repo.BaseRepo.GetDB())

	filter := pageReq.Filter
	sorter := pageReq.Sorter
	query := repo.Client(ctx).
		Where(fmt.Sprintf("%s.borrower_organization_id = ?", loanTable), filter.OrganizationId).
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
				Where(fmt.Sprintf("%s.address->>'address1' like ? or %s.address->>'address2' like ?", loanPropertyTable, loanPropertyTable), "%"+filter.Keyword+"%", "%"+filter.Keyword+"%")
		case "city":
			hasJoinProperty = true
			query = query.Joins(fmt.Sprintf("left join %s on %s.id = %s.loan_id", loanPropertyTable, loanTable, loanPropertyTable)).
				Where(fmt.Sprintf("%s.address->>'city' like ?", loanPropertyTable), "%"+filter.Keyword+"%")
		case "state":
			hasJoinProperty = true
			query = query.Joins(fmt.Sprintf("left join %s on %s.id = %s.loan_id", loanPropertyTable, loanTable, loanPropertyTable)).
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
				query = query.Joins(fmt.Sprintf("left join %s on %s.id = %s.loan_id", loanPropertyTable, loanTable, loanPropertyTable))
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
		return page, errorx.WithStack(err)
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

func (repo *LoanRepo) PageForLenderCompany(ctx context.Context, pageReq *LoanLenderCompanyPageReq) (page gormx.PageResponse[*model.Loan], err error) {
	loanTable := model.Loan{}.TableName()
	loanPropertyTable := model.LoanProperty{}.TableName()
	applicantTable := model.LoanApplicant{}.TableName()
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
				Where(fmt.Sprintf("%s.address->>'address1' like ? or %s.address->>'address2' like ?", loanPropertyTable, loanPropertyTable), "%"+filter.Keyword+"%", "%"+filter.Keyword+"%")
		case "city":
			hasJoinProperty = true
			query = query.Joins(fmt.Sprintf("left join %s on %s.id = %s.loan_id", loanPropertyTable, loanTable, loanPropertyTable)).
				Where(fmt.Sprintf("%s.address->>'city' like ?", loanPropertyTable), "%"+filter.Keyword+"%")
		case "state":
			hasJoinProperty = true
			query = query.Joins(fmt.Sprintf("left join %s on %s.id = %s.loan_id", loanPropertyTable, loanTable, loanPropertyTable)).
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
		return page, errorx.WithStack(err)
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

// GetLenderCompanyForBorrower get lender company
func (repo *LoanRepo) GetLenderCompanyForBorrower(ctx context.Context, creatorId int64, organizationId int64) (result []string, err error) {
	query := repo.BaseRepo.Client(ctx).
		Select("lender_company").
		Where("borrower_organization_id = ?", organizationId).
		Distinct("lender_company")

	if creatorId > 0 {
		query = query.Where("creator_id = ?", creatorId)
	}
	err = query.Pluck("lender_company", &result).
		Error
	return result, errorx.WithStack(err)
}

// GetByBorrowerPropertyId get loans by borrower property id
func (repo *LoanRepo) GetByBorrowerPropertyId(ctx context.Context, borrowerPropertyId int64) (result []*model.Loan, err error) {
	loanPropertyTable := model.LoanProperty{}.TableName()
	loanTable := model.Loan{}.TableName()
	borrowerPropertyTable := model.BorrowerProperty{}.TableName()
	err = repo.BaseRepo.Client(ctx).
		Joins(fmt.Sprintf("left join %s on %s.id = %s.loan_id", loanPropertyTable, loanTable, loanPropertyTable)).
		Joins(fmt.Sprintf("left join %s on %s.borrower_property_id = %s.id", borrowerPropertyTable, loanPropertyTable, borrowerPropertyTable)).
		Where(fmt.Sprintf("%s.id = ?", borrowerPropertyTable), borrowerPropertyId).
		Find(&result).Error
	return result, errorx.WithStack(err)
}
