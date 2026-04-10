package dto

import (
	"database/sql"

	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data"
)

type LoanEntityDto struct {
	Id                     int64                        `json:"id"`
	Code                   string                       `json:"code"`
	Stage                  constants.LoanStage          `json:"stage"`
	PendingStatus          constants.LoanPendingStatus  `json:"pending_status"`
	LoanPurpose            constants.LoanPurpose        `json:"loan_purpose"`
	LoanType               constants.LoanType           `json:"loan_type"`
	TotalAmount            float64                      `json:"total_amount"`
	InitialAmount          float64                      `json:"initial_amount"`
	LoanProperty           *LoanPropertyEntityDto       `json:"loan_property"`
	LenderCompany          string                       `json:"lender_company"`
	ClosingDate            *sql.NullTime                `json:"closing_date"`
	HoldbackAmount         float64                      `json:"holdback_amount"`
	HoldbackReleasedAmount float64                      `json:"holdback_released_amount"`
	Rate                   float32                      `json:"rate"`
	ClosingFund            float64                      `json:"closing_fund"`
	PrincipalAmount        float64                      `json:"principal_amount"`
	MonthlyPayment         float64                      `json:"monthly_payment"`
	LienPosition           string                       `json:"lien_position"`
	MaturityDate           *sql.NullTime                `json:"maturity_date"`
	PrepaymentPenalty      int32                        `json:"prepayment_penalty"`
	ServiceCompany         string                       `json:"service_company"`
	PaidOffDate            *sql.NullTime                `json:"paid_off_date"`
	CreatedAt              *sql.NullTime                `json:"created_at"`
	UpdatedAt              *sql.NullTime                `json:"updated_at"`
	OrganizationId         int64                        `json:"organization_id"`
	BorrowerId             int64                        `json:"borrower_id"`
	BusinessStrategy       []constants.BusinessStrategy `json:"business_strategy"`
	Applicants             []*LoanApplicantEntityDto    `json:"applicants"`
	Source                 constants.LoanSource         `json:"source"`
	LoanTerm               int32                        `json:"loan_term"`
	Ltv                    float64                      `json:"ltv"`
	EstCreditScore         int32                        `json:"est_credit_score"`
	ReservesAmount         float64                      `json:"reserves_amount"`
	TransactionExperiences []*LoanExperienceEntityDto   `json:"transaction_experiences"`
	RealEstateExperiences  []*LoanExperienceEntityDto   `json:"real_estate_experiences"`
	LoanBorrower           *LoanBorrowerEntityDto       `json:"loan_borrower"`
	LoanDemographic        *LoanDemographicEntityDto    `json:"loan_demographic"`
	LoanThirdServices      []*LoanThirdServiceEntityDto `json:"loan_third_services"`
	OriginalContractDate   *sql.NullTime                `json:"original_contract_date"`
	ExtensionDate          *sql.NullTime                `json:"extension_date"`
	LoanProcessor          string                       `json:"loan_processor"`
	Extend                 *data.LoanExtendValueEntity  `json:"extend"`
	LoanDocuments          []*LoanDocumentEntityDto     `json:"loan_documents"`
}

type LoanBorrowerPageRequestDto struct {
	BasePageRequestDto
	Filter LoanBorrowerPageRequestFilter `json:"filter"`
	Sorter BasePageSorter                `json:"sorter"`
}

type LoanBorrowerPageRequestFilter struct {
	UserId         int64 `json:"user_id,omitempty"`
	OrganizationId int64 `json:"organization_id"`
	// BorrowerOrganizationId int64    `json:"borrower_organization_id"`
	Stage         string   `json:"stage,omitempty"`
	Keyword       string   `json:"keyword,omitempty"`
	KeywordType   string   `json:"keyword_type,omitempty"`
	LoanPurpose   []string `json:"loan_purpose,omitempty"`
	LoanType      []string `json:"loan_type,omitempty"`
	Lender        []string `json:"lender,omitempty"`
	PendingStatus []string `json:"pending_status,omitempty"`
}

type LoanLenderPageRequestDto struct {
	BasePageRequestDto
	Filter LoanLenderPageRequestFilter `json:"filter,omitempty"`
	Sorter BasePageSorter              `json:"sorter"`
}

type LoanLenderPageRequestFilter struct {
	OrganizationId int64    `json:"organization_id"`
	Stage          string   `json:"stage,omitempty"`
	Keyword        string   `json:"keyword,omitempty"`
	KeywordType    string   `json:"keyword_type,omitempty"`
	LoanPurpose    []string `json:"loan_purpose,omitempty"`
	LoanType       []string `json:"loan_type,omitempty"`
	PendingStatus  []string `json:"pending_status,omitempty"`
}

type LoanAdminPageRequestDto struct {
	BasePageRequestDto
	Filter LoanAdminPageRequestFilter `json:"filter,omitempty"`
	Sorter BasePageSorter             `json:"sorter"`
}

type LoanAdminPageRequestFilter struct {
	// OrganizationId int64    `json:"organization_id"`
	Stage         string   `json:"stage,omitempty"`
	Keyword       string   `json:"keyword,omitempty"`
	KeywordType   string   `json:"keyword_type,omitempty"`
	LoanPurpose   []string `json:"loan_purpose,omitempty"`
	LoanType      []string `json:"loan_type,omitempty"`
	PendingStatus []string `json:"pending_status,omitempty"`
	Lender        []string `json:"lender,omitempty"`
}
