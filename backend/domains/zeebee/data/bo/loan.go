package bo

import (
	"database/sql"

	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data"
)

type LoanEntityBo struct {
	Id                     int64                        `json:"id"`
	Code                   string                       `json:"code"`
	Stage                  constants.LoanStage          `json:"stage"`
	PendingStatus          constants.LoanPendingStatus  `json:"pending_status"`
	LoanPurpose            constants.LoanPurpose        `json:"loan_purpose"`
	LoanType               constants.LoanType           `json:"loan_type"`
	TotalAmount            float64                      `json:"total_amount"`
	InitialAmount          float64                      `json:"initial_amount"`
	LoanProperty           *LoanPropertyEntityBo        `json:"loan_property"`
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
	Applicants             []*LoanApplicantEntityBo     `json:"applicants"`
	Source                 constants.LoanSource         `json:"source"`
	LoanTerm               int32                        `json:"loan_term"`
	Ltv                    float64                      `json:"ltv"`
	EstCreditScore         int32                        `json:"est_credit_score"`
	ReservesAmount         float64                      `json:"reserves_amount"`
	TransactionExperiences []*LoanExperienceEntityBo    `json:"transaction_experiences"`
	RealEstateExperiences  []*LoanExperienceEntityBo    `json:"real_estate_experiences"`
	LoanBorrower           *LoanBorrowerEntityBo        `json:"loan_borrower"`
	LoanDemographic        *LoanDemographicEntityBo     `json:"loan_demographic"`
	LoanThirdServices      []*LoanThirdServiceEntityBo  `json:"loan_third_services"`
	OriginalContractDate   *sql.NullTime                `json:"original_contract_date"`
	ExtensionDate          *sql.NullTime                `json:"extension_date"`
	LoanProcessor          string                       `json:"loan_processor"`
	Extend                 *data.LoanExtendValueEntity  `json:"extend"`
	LoanDocuments          []*LoanDocumentEntityBo      `json:"loan_documents"`
}
