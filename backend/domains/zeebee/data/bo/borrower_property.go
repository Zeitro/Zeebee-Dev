package bo

import (
	"time"

	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data"
)

type BorrowerPropertyBo struct {
	Id             int64                  `json:"id"`
	PropertyType   constants.PropertyType `json:"property_type"`
	TotalAmount    float64                `json:"total_amount"`
	Unit           int32                  `json:"unit"`
	LivingAreaSqft float32                `json:"living_area_sqft"`
	Bedrooms       int32                  `json:"bedrooms"`
	Bathrooms      int32                  `json:"bathrooms"`
	Address        data.Address           `json:"address"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
	LoanCount      int32                  `json:"loan_count"`
}
