package dto

import (
	"database/sql"

	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data"
)

type LoanPropertyEntityDto struct {
	Id                           int64                  `json:"id"`
	PropertyType                 constants.PropertyType `json:"property_type"`
	BorrowerProperty             BorrowerPropertyDto    `json:"borrower_property"`
	LivingAreaSqft               float32                `json:"living_area_sqft"`
	LotAreaSqft                  float32                `json:"lot_area_sqft"`
	Bedrooms                     int32                  `json:"bedrooms"`
	Bathrooms                    int32                  `json:"bathrooms"`
	FuturePropertyType           constants.PropertyType `json:"future_property_type"`
	FutureUnit                   int32                  `json:"future_unit"`
	FutureLivingAreaSqft         float32                `json:"future_living_area_sqft"`
	FutureBedrooms               int32                  `json:"future_bedrooms"`
	FutureBathrooms              int32                  `json:"future_bathrooms"`
	Owner                        string                 `json:"owner"`
	Unit                         int32                  `json:"unit"`
	ExistingLenderCompany        string                 `json:"existing_lender_company"`
	ExistingRate                 float32                `json:"existing_rate"`
	ExistingMaturityDate         *sql.NullTime          `json:"existing_maturity_date"`
	PurchaseAmount               float64                `json:"purchase_amount"`
	PurchaseDate                 *sql.NullTime          `json:"purchase_date"`
	HasAssignmentFee             bool                   `json:"has_assignment_fee"`
	AssignmentFeeAmount          float64                `json:"assignment_fee_amount"`
	HasSubordinateDebt           bool                   `json:"has_subordinate_debt"`
	SubordinateDebtAmount        float64                `json:"subordinate_debt_amount"`
	ScheduledCloseEscrowDate     *sql.NullTime          `json:"scheduled_close_escrow_date"`
	LiquidityAmount              float64                `json:"liquidity_amount"`
	Budget                       data.LoanBudgetEntity  `json:"budget"`
	MidConstruction              bool                   `json:"mid_construction"`
	EstYearlyRentalIncomeAmount  float64                `json:"est_yearly_rental_income_amount"`
	YearlyLiabilityExpenseAmount float64                `json:"yearly_liability_expense_amount"`
	Dscr                         float64                `json:"dscr"`
	HasDelayedFinancing          bool                   `json:"has_delayed_financing"`
	HasCollateralizedBlanket     bool                   `json:"has_collateralized_blanket"`
	HasPropertiesVacant          bool                   `json:"has_properties_vacant"`
	HasShortTermRental           bool                   `json:"has_short_term_rental"`
	EstPropertyValueAmount       float64                `json:"est_property_value_amount"`
	Address                      data.Address           `json:"address"`
	Arv                          float64                `json:"arv"`
}
