package dto

import (
	"zeebee-backend/domains/zeebee/constants"
)

// StripeCheckoutRequestDto carries a self-serve checkout request. UserId is the
// operator taken from the JWT claim; only the organization owner may check out.
type StripeCheckoutRequestDto struct {
	OrganizationId      int64              `json:"organization_id"`
	UserId              int64              `json:"user_id"`
	PricingPlanId       int64              `json:"pricing_plan_id"`
	PricingPlanPolicyId int64              `json:"pricing_plan_policy_id"`
	UserType            constants.UserType `json:"user_type"`
}

// StripeAutoRenewRequestDto toggles auto-renew (Stripe cancel_at_period_end).
type StripeAutoRenewRequestDto struct {
	OrganizationId int64 `json:"organization_id"`
	UserId         int64 `json:"user_id"`
	Enabled        bool  `json:"enabled"`
}
