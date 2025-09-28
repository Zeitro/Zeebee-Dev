package model

import (
	"database/sql"

	"github.com/lib/pq"

	"zeebee-backend/app/lender/rpc/internal/types"

	"zeebee-backend/common/constants"
)

type Loan struct {
	Base
	Code                     string                      `gorm:"column:code;size:100;not null;default:'';<-:create" json:"code"`
	BorrowerId               int64                       `gorm:"column:borrower_id;not null;default:0" json:"borrower_id"`
	LenderCompanyId          int64                       `gorm:"column:lender_company_id;not null;default:0" json:"lender_company_id"`
	LenderCompany            string                      `gorm:"column:lender_company;size:100;not null;default:''" json:"lender_company"`
	ServiceCompany           string                      `gorm:"column:service_company;size:100;not null;default:''" json:"service_company"`
	Stage                    constants.LoanStage         `gorm:"column:stage;size:10;not null;default:''" json:"stage"`
	PendingStatus            constants.LoanPendingStatus `gorm:"column:pending_status;size:10;not null;default:''" json:"pending_status"`
	LoanPurpose              constants.LoanPurpose       `gorm:"column:loan_purpose;size:20;not null;default:''" json:"loan_purpose"`
	LoanType                 constants.LoanType          `gorm:"column:loan_type;size:30;not null;default:''" json:"loan_type"`
	TotalAmount              float64                     `gorm:"column:total_amount;type:numeric(10,2);not null;default:0" json:"total_amount"`
	PrincipalAmount          float64                     `gorm:"column:principal_amount;type:numeric(10,2);not null;default:0" json:"principle_amount"`
	InitialAmount            float64                     `gorm:"column:initial_amount;type:numeric(10,2);not null;default:0" json:"initial_amount"`
	HoldbackAmount           float64                     `gorm:"column:holdback_amount;type:numeric(10,2);not null;default:0" json:"holdback_amount"`
	HoldbackReleasedAmount   float64                     `gorm:"column:holdback_released_amount;type:numeric(10,2);not null;default:0" json:"holdback_released_amount"`
	PrepaymentPenalty        int32                       `gorm:"column:prepayment_penalty;type:int2;not null;default:0" json:"prepayment_penalty"`
	LienPosition             string                      `gorm:"column:lien_position;size:10;not null;default:''" json:"lien_position"`
	Rate                     float32                     `gorm:"column:rate;type:numeric(4,2);not null;default:0" json:"rate"`
	ClosingFund              float64                     `gorm:"column:closing_fund;type:numeric(10,2);not null;default:0" json:"closing_fund"`
	ClosingDate              *sql.NullTime               `gorm:"column:closing_date" json:"closing_date"`
	MaturityDate             *sql.NullTime               `gorm:"column:maturity_date" json:"maturity_date"`
	PaidOffDate              *sql.NullTime               `gorm:"column:paid_off_date" json:"paid_off_date"`
	MonthlyPayment           float64                     `gorm:"column:monthly_payment;type:numeric(10,2);not null;default:0" json:"monthly_payment"`
	BusinessStrategy         pq.StringArray              `gorm:"column:business_strategy;type:text[]" json:"business_strategy"`
	Source                   constants.LoanSource        `gorm:"column:source;size:30;not null;default:''" json:"source"`
	ExistingLenderCompany    string                      `gorm:"column:existing_lender_company;size:100;not null;default:''" json:"existing_lender_company"`
	ExistingRate             float32                     `gorm:"column:existing_rate;type:numeric(4,2);not null;default:0" json:"existing_rate"`
	ExistingMaturityDate     *sql.NullTime               `gorm:"column:existing_maturity_date" json:"existing_maturity_date"`
	PurchaseAmount           float64                     `gorm:"column:purchase_amount;type:numeric(10,2);not null;default:0" json:"purchase_amount"`
	PurchaseDate             *sql.NullTime               `gorm:"column:purchase_date" json:"purchase_date"`
	HasAssignmentFee         bool                        `gorm:"column:has_assignment_fee;not null;default:false" json:"has_assignment_fee"`
	AssignmentFeeAmount      float64                     `gorm:"column:assignment_fee_amount;type:numeric(10,2);not null;default:0" json:"assignment_fee_amount"`
	HasSubordinateDebt       bool                        `gorm:"column:has_subordinate_debt;not null;default:false" json:"has_subordinate_debt"`
	SubordinateDebtAmount    float64                     `gorm:"column:subordinate_debt_amount;type:numeric(10,2);not null;default:0" json:"subordinate_debt_amount"`
	ScheduledCloseEscrowDate *sql.NullTime               `gorm:"column:scheduled_close_escrow_date" json:"scheduled_close_escrow_date"`
	LoanTerm                 int32                       `gorm:"column:loan_term;type:int2;not null;default:0" json:"loan_term"`
	Arv                      float64                     `gorm:"column:arv;type:numeric(10,2);not null;default:0" json:"arv"`
	Ltv                      float64                     `gorm:"column:ltv;type:numeric(10,2);not null;default:0" json:"ltv"`
	LiquidityAmount          float64                     `gorm:"column:liquidity_amount;type:numeric(10,2);not null;default:0" json:"liquidity_amount"`
	Budget                   types.LoanBudget            `gorm:"column:budget;type:jsonb;not null;default:'{}';serializer:json" json:"budget"`
	MidConstruction          bool                        `gorm:"column:mid_construction;not null;default:false" json:"mid_construction"`
	EstCreditScore           int32                       `gorm:"column:est_credit_score;type:int2;not null;default:0" json:"est_credit_score"`
	EstYearlyRentalIncome    float64                     `gorm:"column:est_yearly_rental_income;type:numeric(10,2);not null;default:0" json:"est_yearly_rental_income"`
	Dscr                     float64                     `gorm:"column:dscr;type:numeric(10,2);not null;default:0" json:"dscr"`
	ReservesAmount           float64                     `gorm:"column:reserves_amount;type:numeric(10,2);not null;default:0" json:"reserves_amount"`
	HasDelayedFinancing      bool                        `gorm:"column:has_delayed_financing;not null;default:false" json:"has_delayed_financing"`
	HasCollateralizedBlanket bool                        `gorm:"column:has_collateralized_blanket;not null;default:false" json:"has_collateralized_blanket"`
	HasPropertiesVacant      bool                        `gorm:"column:has_properties_vacant;not null;default:false" json:"has_properties_vacant"`
	HasShortTermRental       bool                        `gorm:"column:has_short_term_rental;not null;default:false" json:"has_short_term_rental"`
	OriginalContractDate     *sql.NullTime               `json:"original_contract_date,optional" gorm:"column:original_contract_date"`
	ExtensionDate            *sql.NullTime               `json:"extension_date,optional" gorm:"column:extension_date"`
	LoanProcessor            string                      `json:"loan_processor,optional" gorm:"column:loan_processor;size:100;not null;default:''"`
	Extend                   types.LoanExtendValues      `json:"extend" gorm:"column:extend;type:jsonb;serializer:json"`

	// relationships
	Borrower          *Borrower           `json:"borrower" gorm:"foreignKey:BorrowerId;references:Id"`
	LoanProperty      *LoanProperty       `json:"loan_property" gorm:"foreignKey:LoanId;references:Id"`
	Lender            *LenderCompany      `json:"lender" gorm:"foreignKey:LenderCompanyId;references:Id"`
	Applicants        []*Applicant        `json:"applicants" gorm:"foreignKey:LoanId;references:Id"`
	LoanBorrower      *LoanBorrower       `json:"loan_borrower" gorm:"foreignKey:LoanId;references:Id"`
	LoanExperiences   []*LoanExperience   `json:"loan_experiences" gorm:"foreignKey:LoanId;references:Id"`
	LoanThirdServices []*LoanThirdService `json:"loan_third_services" gorm:"foreignKey:LoanId;references:Id"`
	LoanDemographic   *LoanDemographic    `json:"loan_demographic" gorm:"foreignKey:LoanId;references:Id"`
	ApplyFormFile     *File               `json:"apply_form_file" gorm:"polymorphicType:FileType;polymorphicId:RefId;polymorphicValue:ApplyForm"`
	// IdentificationFile   *File               `json:"identification_file" gorm:"polymorphicType:FileType;polymorphicId:RefId;polymorphicValue:Identification"`
	// IncomeFile           *File               `json:"income_file" gorm:"polymorphicType:FileType;polymorphicId:RefId;polymorphicValue:Income"`
	// PropertyFile         *File               `json:"property_file" gorm:"polymorphicType:FileType;polymorphicId:RefId;polymorphicValue:Property"`
	// AssetsFile           *File               `json:"assets_file" gorm:"polymorphicType:FileType;polymorphicId:RefId;polymorphicValue:Assets"`
	// CreditFile           *File               `json:"credit_file" gorm:"polymorphicType:FileType;polymorphicId:RefId;polymorphicValue:Credit"`
	// MiscellaneousFile    *File               `json:"miscellaneous_file" gorm:"polymorphicType:FileType;polymorphicId:RefId;polymorphicValue:Miscellaneous"`
	// InternalDocumentFile *File               `json:"internal_document_file" gorm:"polymorphicType:FileType;polymorphicId:RefId;polymorphicValue:InternalDocument"`
	Documents []*File `json:"documents" gorm:"polymorphicType:BizType;polymorphicId:RefId;polymorphicValue:Loan"`
}

func (m Loan) TableName() string {
	return "loans"
}

func (m *Loan) RemoveRelationships() {
	m.Borrower = nil
	m.LoanProperty = nil
	m.Lender = nil
	m.Applicants = nil
	m.LoanBorrower = nil
	m.LoanExperiences = nil
	m.LoanThirdServices = nil
	m.LoanDemographic = nil
	m.ApplyFormFile = nil
	m.Documents = nil
}
