package dto

import "time"

type PricingPlanPolicyEntityDto struct {
	Id             int64     `json:"id"`
	PricingPlanId  int64     `json:"pricing_plan_id"`
	OriginalAmount float64   `json:"original_amount"`
	ActualAmount   float64   `json:"actual_amount"`
	Type           string    `json:"type"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// PricingPlanResource DTOs

type PricingPlanResourceEntityDto struct {
	Id            int64 `json:"id"`
	PricingPlanId int64 `json:"pricing_plan_id"`
	ResourceId    int64 `json:"resource_id"`
}
