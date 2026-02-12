package model

import (
	"database/sql"

	"zeebee-backend/domains/zeebee/constants"
)

type LoanExperience struct {
	Base
	LoanId           int64                        `gorm:"column:loan_id;not null;default:0" json:"loan_id"`
	OwnershipOnTitle string                       `gorm:"column:ownership_on_title;size:100;not null;default:''" json:"ownership_on_title"`
	PropertyType     constants.PropertyType       `gorm:"column:property_type;size:50;not null;default:''" json:"property_type"`
	PurchaseAmount   float64                      `gorm:"column:purchase_amount;type:numeric(10,2);not null;default:0" json:"purchase_amount"`
	PurchaseDate     *sql.NullTime                `gorm:"column:purchase_date" json:"purchase_date"`
	SoldAmount       float64                      `gorm:"column:sold_amount;type:numeric(10,2);not null;default:0" json:"sold_amount"`
	SoldDate         *sql.NullTime                `gorm:"column:sold_date" json:"sold_date"`
	BudgetOrCost     float64                      `gorm:"column:budget_or_cost;type:numeric(10,2);not null;default:0" json:"budget_or_cost"`
	Profit           float64                      `gorm:"column:profit;type:numeric(10,2);not null;default:0" json:"profit"`
	IsTenantOccupied bool                         `gorm:"column:is_tenant_occupied;not null;default:false" json:"is_tenant_occupied"`
	Type             constants.LoanExperienceType `gorm:"column:type;size:30;not null;default:''" json:"type"`
	Address          Address                      `gorm:"column:address;type:jsonb;not null;default:'{}';serializer:json" json:"address"`
	// relationships
	Loan *Loan `json:"loan" gorm:"foreignKey:LoanId;references:Id"`
}

func (m LoanExperience) TableName() string {
	return "loan_experiences"
}

func (m *LoanExperience) RemoveRelationships() {
	m.Loan = nil
}

// LoanExtendValues represents extended information related to a loan.
type LoanExtendValues struct {
	Pricing           string `json:"pricing"`
	LoanStatusUpdates string `json:"loan_status_updates"`
	MeetingNotes      string `json:"meeting_notes"`
	ActionList        string `json:"action_list"`
	BusinessLine      string `json:"business_line"`
}
