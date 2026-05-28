package dto

import (
	"time"

	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data"
)

type BorrowerPropertyDto struct {
	Id             int64                  `json:"id"`
	PropertyType   constants.PropertyType `json:"property_type"`
	TotalAmount    float64                `json:"total_amount"`
	Unit           int32                  `json:"unit"`
	LivingAreaSqft float32                `json:"living_area_sqft"`
	LotAreaSqft    float32                `json:"lot_area_sqft"`
	Bedrooms       int32                  `json:"bedrooms"`
	Bathrooms      int32                  `json:"bathrooms"`
	Address        data.Address           `json:"address"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
	LoanCount      int32                  `json:"loan_count"`
}

type BorrowerPropertyPageRequestDto struct {
	BasePageRequestDto
	Filter BorrowerPropertyPageRequestFilter `json:"filter"`
	Sorter BasePageSorter                    `json:"sorter"`
}

type BorrowerPropertyPageRequestFilter struct {
	UserId         int64    `json:"user_id"`
	OrganizationId int64    `json:"organization_id"`
	Keyword        string   `json:"keyword,omitempty"`
	Applicant      []string `json:"applicant,omitempty"`
}
