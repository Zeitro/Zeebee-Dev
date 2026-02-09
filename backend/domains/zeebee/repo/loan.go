package repo

import (
	"context"
	"fmt"
	"strconv"

	"gorm.io/gorm"

	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/data/dto"

	"zeebee-backend/common/gormx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/model"
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

// Detail retrieves a loan record by its ID.
func (repo *LoanRepo) Detail(ctx context.Context, id int64, withPreload bool) (result *model.Loan, err error) {
	query := repo.Client(ctx).
		Where("id = ?", id)
	return repo.getDetailByQuery(query, withPreload)
}

func (repo *LoanRepo) getDetailByQuery(query *gorm.DB, withPreload bool) (result *model.Loan, err error) {
	if query == nil {
		err = errorx.New("query is nil")
		return nil, errorx.WithStack(err)
	}
	if withPreload {
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

func (repo *LoanRepo) GetBorrowerLoan(ctx context.Context, id int64, organizationId int64) (result *model.Loan, err error) {
	query := repo.BaseRepo.Client(ctx).
		Where("id = ? and borrower_organization_id = ?", id, organizationId)
	return repo.getDetailByQuery(query, true)
}

func (repo *LoanRepo) GetBorrowerLoanWithTx(ctx context.Context, tx *gorm.DB, id int64, organizationId int64) (result *model.Loan, err error) {
	query := tx.WithContext(ctx).
		Where("id = ? and borrower_organization_id = ?", id, organizationId)
	return repo.getDetailByQuery(query, true)
}

func (repo *LoanRepo) GetLenderLoan(ctx context.Context, id int64, lenderOrganizationId int64, withPreload bool) (result *model.Loan, err error) {
	query := repo.BaseRepo.Client(ctx).
		Where("id = ? and lender_organization_id = ?", id, lenderOrganizationId)
	return repo.getDetailByQuery(query, withPreload)
}

// UpdateOneWithTx update with tx
func (repo *LoanRepo) UpdateOneWithTx(ctx context.Context, tx *gorm.DB, id int64, data *model.Loan) (result *model.Loan, err error) {
	omitFields := []string{"id", "source", "created_at", "code", "creator_id", "creator_name"}
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

// UpdateStageWithTx update stage with tx
func (repo *LoanRepo) UpdateStageWithTx(ctx context.Context, tx *gorm.DB, id int64, stage constants.LoanStage) (err error) {
	err = tx.WithContext(ctx).Model(&model.Loan{}).
		Where("id = ?", id).
		Update("stage", stage).Error
	return errorx.WithStack(err)
}

// UpdatePendingStatusWithTx update pending status with tx
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

// PageForBorrower page for borrower
func (repo *LoanRepo) PageForBorrower(ctx context.Context, pageReq dto.LoanBorrowerPageRequestDto) (page gormx.PageResponse[*model.Loan], err error) {
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
	if len(sorter.Column) > 0 {
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
	if len(sortColumns) == 0 {
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

func (repo *LoanRepo) PageForLender(ctx context.Context, pageReq dto.LoanLenderPageRequestDto) (page gormx.PageResponse[*model.Loan], err error) {
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
		Where(fmt.Sprintf("%s.lender_organization_id = ?", loanTable), filter.OrganizationId).
		Preload("Property").
		Preload("Applicants")
	if len(filter.Stage) > 0 {
		query = query.Where(fmt.Sprintf("%s.stage = ?", loanTable), filter.Stage)
		if filter.Stage == string(constants.LoanStageDraft) {
			query = query.Where(fmt.Sprintf("%s.source in ?", loanTable), []constants.LoanSource{constants.LoanSourcePreapproval, constants.LoanSourceSystem})
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
	if len(sorter.Column) > 0 {
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
	if len(sortColumns) == 0 {
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

func (repo *LoanRepo) PageForAdmin(ctx context.Context, pageReq dto.LoanAdminPageRequestDto) (page gormx.PageResponse[*model.Loan], err error) {
	loanTable := model.Loan{}.TableName()
	loanPropertyTable := model.LoanProperty{}.TableName()
	applicantTable := model.LoanApplicant{}.TableName()
	// organizationTable := model.Organization{}.TableName()
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
		Preload("Property").
		Preload("Applicants")
	sources := []constants.LoanSource{constants.LoanSourceSystem}
	if len(filter.Stage) > 0 {
		query = query.Where(fmt.Sprintf("%s.stage = ?", loanTable), filter.Stage)
		if filter.Stage == string(constants.LoanStageDraft) {
			sources = append(sources, constants.LoanSourcePreapproval)
		}
	}
	query = query.Where(fmt.Sprintf("%s.source in ?", loanTable), sources)

	hasJoinProperty := false

	// keyword
	if len(filter.Keyword) > 0 {
		switch filter.KeywordType {
		case "id":
			id, _ := strconv.ParseInt(filter.Keyword, 10, 64)
			query = query.Where(fmt.Sprintf("%s.id = ?", loanTable), id)
		case "code":
			query = query.Where(fmt.Sprintf("%s.code like ?", loanTable), "%"+filter.Keyword+"%")
		case "lender_organization_name":
			// query = query.Joins(fmt.Sprintf("left join %s on %s.lender_organization_id = %s.id", organizationTable, loanTable, organizationTable)).
			// 	Where(fmt.Sprintf("%s.name like ?", organizationTable), "%"+filter.Keyword+"%").
			// 	Where(fmt.Sprintf("%s.type = ?", organizationTable), constants.UserTypeLender)
			query = query.Where(fmt.Sprintf("%s.lender_company like ?", loanTable), "%"+filter.Keyword+"%")
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
	// lender
	if len(filter.Lender) > 0 {
		query = query.Where(fmt.Sprintf("%s.lender_company in ?", loanTable), filter.Lender)
	}

	sortColumns := []string{}

	// sort
	if len(sorter.Column) > 0 {
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
	if len(sortColumns) == 0 {
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
		Distinct("lender_company").
		Order("lender_company asc")

	if creatorId > 0 {
		query = query.Where("creator_id = ?", creatorId)
	}
	err = query.Pluck("lender_company", &result).
		Error
	return result, errorx.WithStack(err)
}

// GetLenderCompanyForAdmin get lender company for admin
func (repo *LoanRepo) GetLenderCompanyForAdmin(ctx context.Context) (result []string, err error) {
	query := repo.BaseRepo.Client(ctx).
		Select("lender_company").
		Where("source in ?", []constants.LoanSource{constants.LoanSourceSystem, constants.LoanSourcePreapproval}).
		Distinct("lender_company").
		Order("lender_company asc")
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
