package dto

import (
	"database/sql"

	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data"
)

type LoanExperienceEntityDto struct {
	Id               int64                        `json:"id"`
	OwnershipOnTitle string                       `json:"ownership_on_title"`
	PropertyType     constants.PropertyType       `json:"property_type"`
	PurchaseAmount   float64                      `json:"purchase_amount"`
	PurchaseDate     *sql.NullTime                `json:"purchase_date"`
	SoldAmount       float64                      `json:"sold_amount"`
	SoldDate         *sql.NullTime                `json:"sold_date"`
	BudgetOrCost     float64                      `json:"budget_or_cost"`
	Profit           float64                      `json:"profit"`
	IsTenantOccupied bool                         `json:"is_tenant_occupied"`
	Type             constants.LoanExperienceType `json:"type"`
	Address          data.Address                 `json:"address"`
}
