package dto

import (
	"database/sql"
	"time"

	"zeebee-backend/domains/zeebee/constants"
)

type OrganizationSubscriptionEntity struct {
	Id                  int64                  `json:"id"`
	Code                string                 `json:"code"`
	OrganizationId      int64                  `json:"organization_id"`
	Organization        *OrganizationEntityDto `json:"organization"`
	PricingPlanId       int64                  `json:"pricing_plan_id"`
	PricingPlan         *PricingPlanEntityDto  `json:"pricing_plan_name"`
	PricingPlanPolicyId int64                  `json:"pricing_plan_policy_id"`
	Status              string                 `json:"status"`
	ExpireDate          *sql.NullTime          `json:"expire_date"`
	OriginalAmount      float64                `json:"original_amount"`
	ActualAmount        float64                `json:"actual_amount"`
	IsAutoRenew         bool                   `json:"is_auto_renew"`
	CreatedAt           time.Time              `json:"created_at"`
	UpdatedAt           time.Time              `json:"updated_at"`
}

type OrganizationSubscriptionCreateDto struct {
	OrganizationId      int64              `json:"organization_id"`
	PricingPlanId       int64              `json:"pricing_plan_id"`
	PricingPlanPolicyId int64              `json:"pricing_plan_policy_id"`
	Type                constants.UserType `json:"type"`
}

type OrganizationSubscriptionRenewDto struct {
	Id         int64         `json:"id"`
	ExpireDate *sql.NullTime `json:"expire_date"`
}

type OrganizationSubscriptionSwitchStatusDto struct {
	Id     int64                                    `json:"id"`
	Status constants.OrganizationSubscriptionStatus `json:"status"`
}

type OrganizationSubscriptionPageRequestDto struct {
	BasePageRequestDto
	Filter OrganizationSubscriptionPageRequestFilter `json:"filter"`
	Sorter BasePageSorter                            `json:"sorter"`
}

type OrganizationSubscriptionPageRequestFilter struct {
	OrganizationId int64  `json:"organization_id,omitempty"`
	PricingPlanId  int64  `json:"pricing_plan_id,omitempty"`
	Status         string `json:"status,omitempty"`
}

type OrganizationSubscriptionPageEntity struct {
	Id   int64  `json:"id"`
	Code string `json:"code"`
	// OrganizationId      int64         `json:"organization_id"`
	OrganizationName string `json:"organization_name"`
	// PricingPlanId       int64         `json:"pricing_plan_id"`
	PricingPlanName string `json:"pricing_plan_name"`
	// PricingPlanPolicyId int64         `json:"pricing_plan_policy_id"`
	Status         string        `json:"status"`
	ExpireDate     *sql.NullTime `json:"expire_date"`
	OriginalAmount float64       `json:"original_amount"`
	ActualAmount   float64       `json:"actual_amount"`
	IsAutoRenew    bool          `json:"is_auto_renew"`
	CreatedAt      time.Time     `json:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at"`
}
