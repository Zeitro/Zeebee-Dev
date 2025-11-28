package model

import (
	"database/sql"

	"github.com/lib/pq"

	"zeebee-backend/common/constants"
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
	PrincipalAmount        float64                     `gorm:"column:principal_amount;type:numeric(10,2);not null;default:0" json:"principle_amount"`
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
	Extend                 LoanExtendValues            `json:"extend" gorm:"column:extend;type:jsonb;serializer:json"`

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
