package model

import (
	"zeebee-backend/domains/zeebee/constants"
)

type PricingPlanPolicy struct {
	Base
	PricingPlanId  int64                       `gorm:"column:pricing_plan_id;type:int4;not null;default:0" json:"pricing_plan_id"`
	OriginalAmount float64                     `gorm:"column:original_amount;type:decimal(10,2);not null;default:0" json:"original_amount"`
	ActualAmount   float64                     `gorm:"column:actual_amount;type:decimal(10,2);not null;default:0" json:"actual_amount"`
	Type           constants.PricingPolicyType `gorm:"column:type;type:varchar(20);not null;default:'Monthly'" json:"type"`

	PricingPlan *PricingPlan `json:"pricing_plan,omitempty" gorm:"foreignKey:PricingPlanId;references:Id"`
}

func (PricingPlanPolicy) TableName() string {
	return "pricing_plan_policies"
}

func (m *PricingPlanPolicy) RemoveRelationships() {
	m.PricingPlan = nil
}
