package dto

import (
	"database/sql"
	"time"

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
	CreatedAt              time.Time                    `json:"created_at"`
	UpdatedAt              time.Time                    `json:"updated_at"`
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

	// HoldbackReleasedAmount   float64                 `json:"release_holdback_amount"`
	HoldbackRemainingAmount float64                 `json:"holdback_remaining_amount"`
	DueStatus               constants.LoanDueStatus `json:"due_status"`
	DueDays                 int32                   `json:"due_days"`
}

type LoanBorrowerPageRequestDto struct {
	BasePageRequestDto
	Filter LoanBorrowerPageRequestFilter `json:"filter"`
	Sorter BasePageSorter                `json:"sorter"`
}

type LoanBorrowerPageRequestFilter struct {
	UserId                        int64 `json:"user_id,omitempty"`
	OrganizationId                int64 `json:"organization_id"`
	BorrowerPropertyTransactionId int64 `json:"borrower_property_transaction_id"`
	// BorrowerOrganizationId int64    `json:"borrower_organization_id"`
	Stage         constants.LoanStage    `json:"stage,omitempty"`
	Keyword       string                 `json:"keyword,omitempty"`
	KeywordType   string                 `json:"keyword_type,omitempty"`
	LoanPurpose   []string               `json:"loan_purpose,omitempty"`
	LoanType      []string               `json:"loan_type,omitempty"`
	Lender        []string               `json:"lender,omitempty"`
	PendingStatus []string               `json:"pending_status,omitempty"`
	Source        []constants.LoanSource `json:"source,omitempty"`
	DueInDays     int32                  `json:"due_in_days,omitempty"`
	StartDate     *sql.NullTime          `json:"start_date,omitempty"`
	EndDate       *sql.NullTime          `json:"end_date,omitempty"`
	DateType      string                 `json:"date_type"`
}

type LoanLenderPageRequestDto struct {
	BasePageRequestDto
	Filter LoanLenderPageRequestFilter `json:"filter,omitempty"`
	Sorter BasePageSorter              `json:"sorter"`
}

type LoanLenderPageRequestFilter struct {
	OrganizationId int64               `json:"organization_id"`
	Stage          constants.LoanStage `json:"stage,omitempty"`
	Keyword        string              `json:"keyword,omitempty"`
	KeywordType    string              `json:"keyword_type,omitempty"`
	LoanPurpose    []string            `json:"loan_purpose,omitempty"`
	LoanType       []string            `json:"loan_type,omitempty"`
	PendingStatus  []string            `json:"pending_status,omitempty"`
}

type LoanAdminPageRequestDto struct {
	BasePageRequestDto
	Filter LoanAdminPageRequestFilter `json:"filter,omitempty"`
	Sorter BasePageSorter             `json:"sorter"`
}

type LoanAdminPageRequestFilter struct {
	// OrganizationId int64    `json:"organization_id"`
	Stage         constants.LoanStage `json:"stage,omitempty"`
	Keyword       string              `json:"keyword,omitempty"`
	KeywordType   string              `json:"keyword_type,omitempty"`
	LoanPurpose   []string            `json:"loan_purpose,omitempty"`
	LoanType      []string            `json:"loan_type,omitempty"`
	PendingStatus []string            `json:"pending_status,omitempty"`
	Lender        []string            `json:"lender,omitempty"`
}

type LoanSummaryRequestDto struct {
	OrganizationId     int64                  `json:"organization_id"`
	OrganizationUserId int64                  `json:"organization_user_id"`
	StartDate          *sql.NullTime          `json:"start_date"`
	EndDate            *sql.NullTime          `json:"end_date"`
	DateType           string                 `json:"date_type"`
	Source             []constants.LoanSource `json:"source"`
	DueInDays          int32                  `json:"due_in_days"`
	Stage              constants.LoanStage    `json:"stage"`
	Sorter             BasePageSorter         `json:"sorter"`
}

type LoanSummaryResponseDto struct {
	TotalLoanAmount              float64          `json:"total_loan_amount"`
	TotalPrincipalAmount         float64          `json:"total_principal_amount"`
	TotalHoldbackAmount          float64          `json:"total_holdback_amount"`
	TotalClosingFundAmount       float64          `json:"total_closing_fund_amount"`
	TotalMonthlyPaymentAmount    float64          `json:"total_monthly_payment_amount"`
	TotalHoldbackReleasedAmount  float64          `json:"total_holdback_released_amount"`
	TotalHoldbackRemainingAmount float64          `json:"total_remaining_holdback_amount"`
	TotalReservesAmount          float64          `json:"total_reserves_amount"`
	TotalLoanCount               int64            `json:"total_loan_count"`
	AverageRate                  float64          `json:"average_rate"`
	List                         []*LoanEntityDto `json:"list" gorm:"-"`
}
