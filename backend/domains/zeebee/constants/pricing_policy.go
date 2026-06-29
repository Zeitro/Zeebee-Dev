package constants

import "zeebee-backend/common/enums"

// PricingPolicyType pricing policy type
type PricingPolicyType string

const PricingPolicyTypeName = "pricing_policy_type"

const (
	PricingPolicyTypeMonthly   PricingPolicyType = "Monthly"
	PricingPolicyTypeQuarterly PricingPolicyType = "Quarterly"
	PricingPolicyTypeYearly    PricingPolicyType = "Yearly"
)

func (c PricingPolicyType) Validate() error {
	ranges := []PricingPolicyType{
		PricingPolicyTypeMonthly,
		PricingPolicyTypeQuarterly,
		PricingPolicyTypeYearly,
	}
	return validateConstant(c, ranges, PricingPolicyTypeName)
}

func init() {
	enums.Register(PricingPolicyTypeName, []enums.EnumOpt{
		{Label: string(PricingPolicyTypeMonthly), Value: string(PricingPolicyTypeMonthly)},
		{Label: string(PricingPolicyTypeQuarterly), Value: string(PricingPolicyTypeQuarterly)},
		{Label: string(PricingPolicyTypeYearly), Value: string(PricingPolicyTypeYearly)},
	})
}
