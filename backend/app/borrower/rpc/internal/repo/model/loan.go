package model

import (
	"database/sql"

	"github.com/lib/pq"

	"zeebee-backend/common/constants"
)

type Loan struct {
	Base
	BorrowerId int64 `gorm:"column:borrower_id;not null;default:0" json:"borrower_id"`
	// PropertyId             int64                       `gorm:"column:property_id;not null;default:0" json:"property_id"`
	LenderId               int64                       `gorm:"column:lender_id;not null;default:0" json:"lender_id"`
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

	Borrower     *Borrower     `json:"borrower" gorm:"foreignKey:BorrowerId;references:Id"`
	LoanProperty *LoanProperty `json:"loan_property" gorm:"foreignKey:LoanId;references:Id"`
	Lender       *Lender       `json:"lender" gorm:"foreignKey:LenderId;references:Id"`
	Applicants   []*Applicant  `json:"applicants" gorm:"foreignKey:LoanId;references:Id"`
}

func (m Loan) TableName() string {
	return "loans"
}
