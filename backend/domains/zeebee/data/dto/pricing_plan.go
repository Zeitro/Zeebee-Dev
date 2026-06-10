package dto

import (
	"time"

	"zeebee-backend/domains/zeebee/constants"
)

type PricingPlanEntity struct {
	Id          int64     `json:"id"`
	Name        string    `json:"name"`
	Subtitle    string    `json:"subtitle"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	TrialDays   int16     `json:"trial_days"`
	SortOrder   int16     `json:"sort_order"`
	IsPublic    bool      `json:"is_public"`
	IsDefault   bool      `json:"is_default"`
	Type        string    `json:"type"`
	PlanType    string    `json:"plan_type"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// PricingPlan create/update DTOs

type PricingPlanCreateDto struct {
	Name        string                    `json:"name"`
	Subtitle    string                    `json:"subtitle,omitempty"`
	Description string                    `json:"description,omitempty"`
	Status      constants.CommonStatus    `json:"status,omitempty"`
	TrialDays   int16                     `json:"trial_days,omitempty"`
	SortOrder   int16                     `json:"sort_order,omitempty"`
	IsPublic    bool                      `json:"is_public,omitempty"`
	IsDefault   bool                      `json:"is_default"`
	Type        constants.UserType        `json:"type"`
	PlanType    constants.PricingPlanType `json:"plan_type,omitempty"`
}

type PricingPlanUpdateDto struct {
	Id          int64                     `json:"id"`
	Subtitle    string                    `json:"subtitle,omitempty"`
	Description string                    `json:"description,omitempty"`
	Status      constants.CommonStatus    `json:"status,omitempty"`
	TrialDays   int16                     `json:"trial_days,omitempty"`
	SortOrder   int16                     `json:"sort_order,omitempty"`
	IsPublic    bool                      `json:"is_public,omitempty"`
	IsDefault   bool                      `json:"is_default"`
	PlanType    constants.PricingPlanType `json:"plan_type,omitempty"`
}

type PricingPlanBindResourcesDto struct {
	Id          int64   `json:"id"`
	ResourceIds []int64 `json:"resource_ids"`
}

// Page request/response

type PricingPlanPageRequestDto struct {
	BasePageRequestDto
	Filter PricingPlanPageRequestFilter `json:"filter"`
	Sorter BasePageSorter               `json:"sorter"`
}

type PricingPlanPageRequestFilter struct {
	Keyword  string `json:"keyword,omitempty"`
	Status   string `json:"status,omitempty"`
	Type     string `json:"type,omitempty"`
	PlanType string `json:"plan_type,omitempty"`
}

type PricingPlanPageEntity struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	Subtitle string `json:"subtitle"`
	// Description string    `json:"description"`
	Status    string    `json:"status"`
	TrialDays int16     `json:"trial_days"`
	SortOrder int16     `json:"sort_order"`
	IsPublic  bool      `json:"is_public"`
	IsDefault bool      `json:"is_default"`
	Type      string    `json:"type"`
	PlanType  string    `json:"plan_type"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// PricingPlanEntityDto detail entity (includes policies)
type PricingPlanEntityDto struct {
	Id          int64                        `json:"id"`
	Name        string                       `json:"name"`
	Subtitle    string                       `json:"subtitle"`
	Description string                       `json:"description"`
	Status      string                       `json:"status"`
	TrialDays   int16                        `json:"trial_days"`
	SortOrder   int16                        `json:"sort_order"`
	IsPublic    bool                         `json:"is_public"`
	IsDefault   bool                         `json:"is_default"`
	Type        string                       `json:"type"`
	PlanType    string                       `json:"plan_type"`
	Policies    []PricingPlanPolicyEntityDto `json:"policies"`
	CreatedAt   time.Time                    `json:"created_at"`
	UpdatedAt   time.Time                    `json:"updated_at"`
}

// PricingPlanPolicy DTOs

type PricingPlanPolicyCreateDto struct {
	PricingPlanId  int64                       `json:"pricing_plan_id"`
	OriginalAmount float64                     `json:"original_amount"`
	ActualAmount   float64                     `json:"actual_amount"`
	Type           constants.PricingPolicyType `json:"type"`
}

type PricingPlanPolicyUpdateDto struct {
	Id             int64                       `json:"id"`
	OriginalAmount float64                     `json:"original_amount,omitempty"`
	ActualAmount   float64                     `json:"actual_amount,omitempty"`
	Type           constants.PricingPolicyType `json:"type,omitempty"`
}
