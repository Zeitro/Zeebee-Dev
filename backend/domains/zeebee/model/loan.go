package model

import (
	"database/sql"

	"github.com/lib/pq"

	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data"
)

type Loan struct {
	Base
	Code                   string                      `gorm:"column:code;size:100;not null;default:'';<-:create" json:"code"`
	CreatorId              int64                       `gorm:"column:creator_id;not null;default:0" json:"creator_id"`
	CreatorName            string                      `gorm:"column:creator_name;type:varchar(100);default:'';not null" json:"creator_name"`
	LenderOrganizationId   int64                       `gorm:"column:lender_organization_id;not null;default:0" json:"lender_organization_id"`
	BorrowerOrganizationId int64                       `gorm:"column:borrower_organization_id;not null;default:0" json:"borrower_organization_id"`
	LenderCompany          string                      `gorm:"column:lender_company;size:100;not null;default:''" json:"lender_company"`
	ServiceCompany         string                      `gorm:"column:service_company;size:100;not null;default:''" json:"service_company"`
	Stage                  constants.LoanStage         `gorm:"column:stage;size:10;not null;default:''" json:"stage"`
	PendingStatus          constants.LoanPendingStatus `gorm:"column:pending_status;size:10;not null;default:''" json:"pending_status"`
	LoanPurpose            constants.LoanPurpose       `gorm:"column:loan_purpose;size:20;not null;default:''" json:"loan_purpose"`
	LoanType               constants.LoanType          `gorm:"column:loan_type;size:30;not null;default:''" json:"loan_type"`
	TotalAmount            float64                     `gorm:"column:total_amount;type:numeric(10,2);not null;default:0" json:"total_amount"`
	PrincipalAmount        float64                     `gorm:"column:principal_amount;type:numeric(10,2);not null;default:0" json:"principal_amount"`
	InitialAmount          float64                     `gorm:"column:initial_amount;type:numeric(10,2);not null;default:0" json:"initial_amount"`
	HoldbackAmount         float64                     `gorm:"column:holdback_amount;type:numeric(10,2);not null;default:0" json:"holdback_amount"`
	HoldbackReleasedAmount float64                     `gorm:"column:holdback_released_amount;type:numeric(10,2);not null;default:0" json:"holdback_released_amount"`
	PrepaymentPenalty      int32                       `gorm:"column:prepayment_penalty;type:int2;not null;default:0" json:"prepayment_penalty"`
	LienPosition           string                      `gorm:"column:lien_position;size:10;not null;default:''" json:"lien_position"`
	Rate                   float32                     `gorm:"column:rate;type:numeric(4,2);not null;default:0" json:"rate"`
	ClosingFund            float64                     `gorm:"column:closing_fund;type:numeric(10,2);not null;default:0" json:"closing_fund"`
	ClosingDate            *sql.NullTime               `gorm:"column:closing_date" json:"closing_date"`
	MaturityDate           *sql.NullTime               `gorm:"column:maturity_date" json:"maturity_date"`
	PaidOffDate            *sql.NullTime               `gorm:"column:paid_off_date" json:"paid_off_date"`
	MonthlyPayment         float64                     `gorm:"column:monthly_payment;type:numeric(10,2);not null;default:0" json:"monthly_payment"`
	BusinessStrategy       pq.StringArray              `gorm:"column:business_strategy;type:text[]" json:"business_strategy"`
	Source                 constants.LoanSource        `gorm:"column:source;size:30;not null;default:''" json:"source"`
	LoanTerm               int32                       `gorm:"column:loan_term;type:int2;not null;default:0" json:"loan_term"`
	Ltv                    float64                     `gorm:"column:ltv;type:numeric(10,2);not null;default:0" json:"ltv"`
	EstCreditScore         int32                       `gorm:"column:est_credit_score;type:int2;not null;default:0" json:"est_credit_score"`
	ReservesAmount         float64                     `gorm:"column:reserves_amount;type:numeric(10,2);not null;default:0" json:"reserves_amount"`
	OriginalContractDate   *sql.NullTime               `json:"original_contract_date" gorm:"column:original_contract_date"`
	ExtensionDate          *sql.NullTime               `json:"extension_date" gorm:"column:extension_date"`
	LoanProcessor          string                      `json:"loan_processor" gorm:"column:loan_processor;size:100;not null;default:''"`
	Extend                         data.LoanExtendValueEntity  `json:"extend" gorm:"column:extend;type:jsonb;serializer:json"`
	BorrowerPropertyTransactionId  int64                       `gorm:"column:borrower_property_transaction_id;not null;default:0" json:"borrower_property_transaction_id"`

	// relationships
	Creator              *User               `json:"creator" gorm:"foreignKey:CreatorId;references:Id"`
	BorrowerOrganization *Organization       `json:"borrower_organization" gorm:"foreignKey:BorrowerOrganizationId;references:Id"`
	LenderOrganization   *Organization       `json:"lender_organization" gorm:"foreignKey:LenderOrganizationId;references:Id"`
	LoanProperty         *LoanProperty       `json:"loan_property" gorm:"foreignKey:LoanId;references:Id"`
	Lender               *Organization       `json:"lender" gorm:"foreignKey:LenderOrganizationId;references:Id"`
	Applicants           []*LoanApplicant    `json:"applicants" gorm:"foreignKey:LoanId;references:Id"`
	LoanBorrower         *LoanBorrower       `json:"loan_borrower" gorm:"foreignKey:LoanId;references:Id"`
	LoanExperiences      []*LoanExperience   `json:"loan_experiences" gorm:"foreignKey:LoanId;references:Id"`
	LoanThirdServices    []*LoanThirdService `json:"loan_third_services" gorm:"foreignKey:LoanId;references:Id"`
	LoanDemographic      *LoanDemographic    `json:"loan_demographic" gorm:"foreignKey:LoanId;references:Id"`
	LoanDocuments        []*LoanDocument     `json:"loan_documents" gorm:"foreignKey:LoanId;references:Id"`
}

func (m Loan) TableName() string {
	return "loans"
}

func (m *Loan) RemoveRelationships() {
	m.Creator = nil
	m.BorrowerOrganization = nil
	m.LenderOrganization = nil
	m.LoanProperty = nil
	m.Lender = nil
	m.Applicants = nil
	m.LoanBorrower = nil
	m.LoanExperiences = nil
	m.LoanThirdServices = nil
	m.LoanDemographic = nil
	m.LoanDocuments = nil
}

// LoanExtendValues represents extended information related to a loan.
// type LoanExtendValues struct {
// 	Pricing             string               `json:"pricing,omitempty"`
// 	LoanStatusUpdates   string               `json:"loan_status_updates,omitempty"`
// 	MeetingNotes        string               `json:"meeting_notes,omitempty"`
// 	ActionList          string               `json:"action_list,omitempty"`
// 	BusinessLine        string               `json:"business_line,omitempty"`
// 	LoanPreapprovalInfo *LoanPreapprovalInfo `json:"loan_preapproval_info,omitempty"`
// }
//
// type LoanPreapprovalInfo struct {
// 	Input  *LoanPreapprovalInfoInput  `json:"input,omitempty"`
// 	Result *LoanPreapprovalInfoResult `json:"result,omitempty"`
// }
//
// type LoanPreapprovalInfoInput struct {
// 	Borrower *LoanPreapprovalInfoInputBorrower `json:"borrower,omitempty"`
// 	Property *LoanPreapprovalInfoInputProperty `json:"property,omitempty"`
// 	Loan     *LoanPreapprovalInfoInputLoan     `json:"loan,omitempty"`
// }
//
// type LoanPreapprovalInfoResult struct {
// 	LenderOrganizationName string  `json:"lender_organization_name,omitempty"`
// 	LenderOrganizationId   int64   `json:"lender_organization_id,omitempty"`
// 	IsApproved             bool    `json:"is_approved,omitempty"`
// 	RejectReason           string  `json:"reject_reason,omitempty"`
// 	TotalAmount            float64 `json:"total_amount,omitempty"`
// 	InterestReservedAmount float64 `json:"interest_reserved_amount,omitempty"`
// 	Rate                   float32 `json:"rate,omitempty"`
// 	MonthlyPayment         float64 `json:"monthly_payment,omitempty"`
// 	CashToClose            float64 `json:"cash_to_close,omitempty"`
// 	InitialAmount          float64 `json:"initial_amount,omitempty"`
// 	HoldbackAmount         float64 `json:"holdback_amount,omitempty"`
// 	LoanTerm               int32   `json:"loan_term,omitempty"`
// 	Ltv                    float32 `json:"ltv,omitempty"`
// 	Ltc                    float32 `json:"ltc,omitempty"`
// 	Ltarv                  float32 `json:"ltarv,omitempty"`
// 	ExperienceLevel        string  `json:"experience_level,omitempty"`
// }
//
// type LoanPreapprovalInfoInputBorrower struct {
// 	Name                string `json:"name,omitempty"`
// 	Email               string `json:"email,omitempty"`
// 	Phone               string `json:"phone,omitempty"`
// 	EstCreditScore      int32  `json:"est_credit_score,omitempty"`
// 	FlipPropertyCount   int32  `json:"flip_property_count,omitempty"`
// 	PurchaseRentalCount int32  `json:"purchase_rental_property_count,omitempty"`
// }
//
// type LoanPreapprovalInfoInputProperty struct {
// 	Address       *Address `json:"address,omitempty"`
// 	PropertyType  string   `json:"property_type,omitempty"`
// 	IsPermitReady bool     `json:"is_permit_ready,omitempty"`
// 	Unit          int32    `json:"unit,omitempty"`
// }
//
// type LoanPreapprovalInfoInputLoan struct {
// 	LoanTerm               int32         `json:"loan_term,omitempty"`
// 	EstPropertyValueAmount float64       `json:"est_property_value_amount,omitempty"`
// 	Arv                    float64       `json:"arv,omitempty"`
// 	TotalAmount            float64       `json:"total_amount,omitempty"`
// 	PrincipalAmount        float64       `json:"principal_amount,omitempty"`
// 	PurchaseDate           *sql.NullTime `json:"purchase_date,omitempty"`
// 	PurchaseAmount         float64       `json:"purchase_amount,omitempty"`
// 	CostSpentToDateAmount  float64       `json:"cost_spent_to_date_amount,omitempty"`
// 	BudgetAmount           float64       `json:"budget_amount,omitempty"`
// 	InitialAmount          float64       `json:"initial_amount,omitempty"`
// 	HoldbackAmount         float64       `json:"holdback_amount,omitempty"`
// 	InterestReserveMonths  int32         `json:"interest_reserve_months,omitempty"`
// 	LoanPurpose            string        `json:"loan_purpose,omitempty"`
// 	LoanType               string        `json:"loan_type,omitempty"`
// 	HasAssignmentFee       bool          `json:"has_assignment_fee,omitempty"`
// 	AssignmentFeeAmount    float64       `json:"assignment_fee_amount,omitempty"`
// }
