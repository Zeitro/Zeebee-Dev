package model

import (
	"database/sql"

	"zeebee-backend/domains/zeebee/constants"
)

type LoanProperty struct {
	Base
	LoanId                       int64                  `gorm:"column:loan_id;not null;default:0" json:"loan_id"`
	BorrowerPropertyId           int64                  `gorm:"column:borrower_property_id;not null;default:0" json:"borrower_property_id"`
	Unit                         int32                  `gorm:"column:unit;not null;default:1" json:"unit"`
	PropertyType                 constants.PropertyType `gorm:"column:property_type;size:30;not null;default:''" json:"property_type"`
	LivingAreaSqft               float32                `gorm:"column:living_area_sqft;type:numeric(8,2);not null;default:0" json:"living_area_sqft"`
	Bedrooms                     int32                  `gorm:"column:bedrooms;not null;default:0" json:"bedrooms"`
	Bathrooms                    int32                  `gorm:"column:bathrooms;not null;default:0" json:"bathrooms"`
	FuturePropertyType           constants.PropertyType `gorm:"column:future_property_type;size:30;not null;default:''" json:"future_property_type"`
	FutureUnit                   int32                  `gorm:"column:future_unit;not null;default:0" json:"future_unit"`
	FutureLivingAreaSqft         float32                `gorm:"column:future_living_area_sqft;type:numeric(8,2);not null;default:0" json:"future_living_area_sqft"`
	FutureBedrooms               int32                  `gorm:"column:future_bedrooms;not null;default:0" json:"future_bedrooms"`
	FutureBathrooms              int32                  `gorm:"column:future_bathrooms;not null;default:0" json:"future_bathrooms"`
	Owner                        string                 `gorm:"column:owner;size:255;not null;default:''" json:"owner"`
	Address                      Address                `gorm:"column:address;type:jsonb;serializer:json" json:"address"`
	ExistingLenderCompany        string                 `gorm:"column:existing_lender_company;size:100;not null;default:''" json:"existing_lender_company"`
	ExistingRate                 float32                `gorm:"column:existing_rate;type:numeric(4,2);not null;default:0" json:"existing_rate"`
	ExistingMaturityDate         *sql.NullTime          `gorm:"column:existing_maturity_date" json:"existing_maturity_date"`
	PurchaseAmount               float64                `gorm:"column:purchase_amount;type:decimal(10,2);not null;default:0" json:"purchase_amount"`
	PurchaseDate                 *sql.NullTime          `gorm:"column:purchase_date" json:"purchase_date"`
	HasAssignmentFee             bool                   `gorm:"column:has_assignment_fee;not null;default:false" json:"has_assignment_fee"`
	AssignmentFeeAmount          float64                `gorm:"column:assignment_fee_amount;type:decimal(10,2);not null;default:0" json:"assignment_fee_amount"`
	HasSubordinateDebt           bool                   `gorm:"column:has_subordinate_debt;not null;default:false" json:"has_subordinate_debt"`
	SubordinateDebtAmount        float64                `gorm:"column:subordinate_debt_amount;type:decimal(10,2);not null;default:0" json:"subordinate_debt_amount"`
	ScheduledCloseEscrowDate     *sql.NullTime          `gorm:"column:scheduled_close_escrow_date" json:"scheduled_close_escrow_date"`
	LiquidityAmount              float64                `gorm:"column:liquidity_amount;type:decimal(10,2);not null;default:0" json:"liquidity_amount"`
	Budget                       LoanBudget             `gorm:"column:budget;type:jsonb;not null;default:'{}';serializer:json" json:"budget"`
	MidConstruction              bool                   `gorm:"column:mid_construction;not null;default:false" json:"mid_construction"`
	EstYearlyRentalIncomeAmount  float64                `gorm:"column:est_yearly_rental_income_amount;type:decimal(10,2);not null;default:0" json:"est_yearly_rental_income_amount"`
	YearlyLiabilityExpenseAmount float64                `gorm:"column:yearly_liability_expense_amount;type:decimal(10,2);not null;default:0" json:"yearly_liability_expense_amount"`
	Dscr                         float64                `gorm:"column:dscr;type:numeric(10,2);not null;default:0" json:"dscr"`
	HasDelayedFinancing          bool                   `gorm:"column:has_delayed_financing;not null;default:false" json:"has_delayed_financing"`
	HasCollateralizedBlanket     bool                   `gorm:"column:has_collateralized_blanket;not null;default:false" json:"has_collateralized_blanket"`
	HasPropertiesVacant          bool                   `gorm:"column:has_properties_vacant;not null;default:false" json:"has_properties_vacant"`
	HasShortTermRental           bool                   `gorm:"column:has_short_term_rental;not null;default:false" json:"has_short_term_rental"`
	EstPropertyValueAmount       float64                `gorm:"column:est_property_value_amount;type:decimal(10,2);not null;default:0" json:"est_property_value_amount"`
	Arv                          float64                `gorm:"column:arv;type:decimal(10,2);not null;default:0" json:"arv"`

	// relationships
	BorrowerProperty *BorrowerProperty `json:"borrower_property" gorm:"foreignKey:BorrowerPropertyId;references:Id"`
	Loan             *Loan             `json:"loan" gorm:"foreignKey:LoanId;references:Id"`
	Applicants       []*LoanApplicant  `json:"applicants" gorm:"foreignKey:LoanPropertyId;references:Id"`
}

func (m LoanProperty) TableName() string {
	return "loan_properties"
}

func (m *LoanProperty) RemoveRelationships() {
	m.BorrowerProperty = nil
	m.Loan = nil
	m.Applicants = nil
}

type LoanBudget struct {
	OriginalBudget        float64 `json:"original_budget"`
	RemainingBudget       float64 `json:"remaining_budget"`
	CostSpentToDateAmount float64 `json:"cost_spent_to_date_amount"`
}
