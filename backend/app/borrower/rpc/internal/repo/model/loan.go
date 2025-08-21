package model

import (
	"database/sql"

	"github.com/lib/pq"

	"zeebee-backend/app/borrower/rpc/internal/types"
	"zeebee-backend/common/constants"
)

// CREATE TABLE "public"."loans" (
//
//	"id" int4 NOT NULL DEFAULT nextval('loans_id_seq'::regclass),
//	"borrower_id" int4 NOT NULL DEFAULT 0,
//	"lender_id" int4 NOT NULL DEFAULT 0,
//	"stage" varchar(10) COLLATE "pg_catalog"."default" NOT NULL DEFAULT 'Pending'::character varying,
//	"pending_status" varchar(10) COLLATE "pg_catalog"."default" NOT NULL DEFAULT 'Quoting'::character varying,
//	"loan_purpose" varchar(20) COLLATE "pg_catalog"."default" NOT NULL DEFAULT 'Purchase'::character varying,
//	"total_amount" numeric(10,2) NOT NULL DEFAULT 0,
//	"initial_amount" numeric(10,2) NOT NULL DEFAULT 0,
//	"closing_date" date,
//	"holdback_amount" numeric(10,2) NOT NULL DEFAULT 0,
//	"rate" numeric(4,2) NOT NULL DEFAULT 0,
//	"closing_fund" numeric(10,2) NOT NULL DEFAULT 0,
//	"created_at" timestamptz(6) NOT NULL DEFAULT now(),
//	"updated_at" timestamptz(6) NOT NULL DEFAULT now(),
//	"is_del" int2 NOT NULL DEFAULT 0,
//	"loan_type" varchar(255) COLLATE "pg_catalog"."default" NOT NULL DEFAULT 'Bridge'::character varying,
//	"paid_off_date" date,
//	"principal_amount" numeric(10,2) NOT NULL DEFAULT 0,
//	"monthly_payment" numeric(10,2) NOT NULL DEFAULT 0,
//	"maturity_date" date,
//	"holdback_released_amount" numeric(10,2) NOT NULL DEFAULT 0,
//	"lender_company" varchar(255) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
//	"service_company" varchar(255) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
//	"business_strategy" text[] COLLATE "pg_catalog"."default" DEFAULT '{}'::text[],
//	"source" varchar(10) COLLATE "pg_catalog"."default" NOT NULL DEFAULT 'system'::character varying,
//	"lien_position" varchar(255) COLLATE "pg_catalog"."default" DEFAULT ''::character varying,
//	"prepayment_penalty" int2 NOT NULL DEFAULT 0,
//	"existing_lender_company" varchar(255) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
//	"existing_rate" numeric(4,2) NOT NULL DEFAULT 0,
//	"existing_maturity_date" date,
//	"purchase_amount" numeric(10,2) NOT NULL DEFAULT 0,
//	"purchase_date" date,
//	"has_assignment_fee" bool NOT NULL DEFAULT false,
//	"assignment_fee_amount" numeric(10,2) NOT NULL DEFAULT 0,
//	"has_subordinate_debt" bool NOT NULL DEFAULT false,
//	"subordinate_debt_amount" numeric(10,2) NOT NULL DEFAULT 0,
//	"scheduled_close_escrow_date" date,
//	"loan_term" int2 NOT NULL DEFAULT 0,
//	"arv" numeric(10,2) NOT NULL DEFAULT 0,
//	"ltv" numeric(10,2) NOT NULL DEFAULT 0,
//	"liquidity_amount" numeric(10,2) NOT NULL DEFAULT 0,
//	"budget" jsonb NOT NULL DEFAULT '{}'::jsonb,
//	"mid_construction" bool NOT NULL DEFAULT false,
//	"est_credit_score" int2 NOT NULL DEFAULT 0,
//	"est_yearly_rental_income" numeric(10,2) NOT NULL DEFAULT 0,
//	"yearly_liability_expense" numeric(10,2) NOT NULL DEFAULT 0,
//	"dscr" numeric(8,2) NOT NULL DEFAULT 0,
//	"reserves" numeric(10,2) NOT NULL DEFAULT 0,
//	"has_delayed_financing" bool NOT NULL DEFAULT false,
//	"has_collateralized_blanket" bool NOT NULL DEFAULT false,
//	"has_properties_vacant" bool NOT NULL DEFAULT false,
//	"has_short_term_rental" bool NOT NULL DEFAULT false,
//	"additional_info" varchar(3000) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
//	CONSTRAINT "loans_pkey" PRIMARY KEY ("id")
//
// )
type Loan struct {
	Base
	BorrowerId int64 `gorm:"column:borrower_id;not null;default:0" json:"borrower_id"`
	// PropertyId             int64                       `gorm:"column:property_id;not null;default:0" json:"property_id"`
	LenderId                 int64                       `gorm:"column:lender_id;not null;default:0" json:"lender_id"`
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
	Budget                   types.LoanBudget            `gorm:"column:budget;type:jsonb;not null;default:'{}'" json:"budget"`
	MidConstruction          bool                        `gorm:"column:mid_construction;not null;default:false" json:"mid_construction"`
	EstCreditScore           int32                       `gorm:"column:est_credit_score;type:int2;not null;default:0" json:"est_credit_score"`
	EstYearlyRentalIncome    float64                     `gorm:"column:est_yearly_rental_income;type:numeric(10,2);not null;default:0" json:"est_yearly_rental_income"`
	Dscr                     float64                     `gorm:"column:dscr;type:numeric(10,2);not null;default:0" json:"dscr"`
	Reserves                 float64                     `gorm:"column:reserve;type:numeric(10,2);not null;default:0" json:"reserves"`
	HasDelayedFinancing      bool                        `gorm:"column:has_delayed_financing;not null;default:false" json:"has_delayed_financing"`
	HasCollateralizedBlanket bool                        `gorm:"column:has_collateralized_blanket;not null;default:false" json:"has_collateralized_blanket"`
	HasPropertiesVacant      bool                        `gorm:"column:has_properties_vacant;not null;default:false" json:"has_properties_vacant"`
	HasShortTermRental       bool                        `gorm:"column:has_short_term_rental;not null;default:false" json:"has_short_term_rental"`
	AdditionalInfo           string                      `gorm:"column:additional_info;type:varchar(3000)" json:"additional_info"`

	// relationships
	Borrower          *Borrower           `json:"borrower" gorm:"foreignKey:BorrowerId;references:Id"`
	LoanProperty      *LoanProperty       `json:"loan_property" gorm:"foreignKey:LoanId;references:Id"`
	Lender            *Lender             `json:"lender" gorm:"foreignKey:LenderId;references:Id"`
	Applicants        []*Applicant        `json:"applicants" gorm:"foreignKey:LoanId;references:Id"`
	LoanBorrower      *LoanBorrower       `json:"loan_borrower" gorm:"foreignKey:LoanId;references:Id"`
	LoanExperiences   []*LoanExperience   `json:"loan_experiences" gorm:"foreignKey:LoanId;references:Id"`
	LoanThirdServices []*LoanThirdService `json:"loan_third_services" gorm:"foreignKey:LoanId;references:Id"`
	LoanDemographic   *LoanDemographic    `json:"loan_demographic" gorm:"foreignKey:LoanId;references:Id"`
}

func (m Loan) TableName() string {
	return "loans"
}
