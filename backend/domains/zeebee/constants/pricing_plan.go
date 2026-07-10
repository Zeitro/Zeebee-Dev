package constants

import (
	"zeebee-backend/common/enums"
)

// PricingPlanType pricing plan type
type PricingPlanType string

const PricingPlanTypeName = "pricing_plan_type"

const (
	PricingPlanTypeStandard PricingPlanType = "Standard"
	PricingPlanTypeCustom   PricingPlanType = "Custom"
)

func (c PricingPlanType) Validate() error {
	ranges := []PricingPlanType{
		PricingPlanTypeStandard,
		PricingPlanTypeCustom,
	}
	return validateConstant(c, ranges, PricingPlanTypeName)
}

func init() {
	enums.Register(PricingPlanTypeName, []enums.EnumOpt{
		{Label: string(PricingPlanTypeStandard), Value: string(PricingPlanTypeStandard)},
		{Label: string(PricingPlanTypeCustom), Value: string(PricingPlanTypeCustom)},
	})

}
