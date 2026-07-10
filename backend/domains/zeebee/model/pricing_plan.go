package model

import (
	"zeebee-backend/domains/zeebee/constants"
)

type PricingPlan struct {
	Base
	Name        string                    `gorm:"column:name;type:varchar(100);not null;default:''" json:"name"`
	Subtitle    string                    `gorm:"column:subtitle;type:varchar(255);not null;default:''" json:"subtitle"`
	Description string                    `gorm:"column:description;type:varchar(255);not null;default:''" json:"description"`
	Status      constants.CommonStatus    `gorm:"column:status;type:varchar(50);not null;default:'Active'" json:"status"`
	TrialDays   int16                     `gorm:"column:trial_days;type:int2;not null;default:0" json:"trial_days"`
	SortOrder   int16                     `gorm:"column:sort_order;type:int2;not null;default:6" json:"sort_order"`
	IsPublic    bool                      `gorm:"column:is_public;type:bool;not null;default:false" json:"is_public"`
	IsDefault   bool                      `gorm:"column:is_default;type:bool;not null;default:false" json:"is_default"`
	Type        constants.UserType        `gorm:"column:type;type:varchar(50);not null;default:''" json:"type"`
	PlanType    constants.PricingPlanType `gorm:"column:plan_type;type:varchar(50);not null;default:'Standard'" json:"plan_type"`

	Policies      []*PricingPlanPolicy        `json:"policies,omitempty" gorm:"foreignKey:PricingPlanId;references:Id"`
	Subscriptions []*OrganizationSubscription `json:"subscriptions,omitempty" gorm:"foreignKey:PricingPlanId;references:Id"`
	Resources     []*Resource                 `json:"resources,omitempty" gorm:"many2many:pricing_plan_and_resources;joinForeignKey:pricing_plan_id;joinReferences:resource_id" `
}

func (PricingPlan) TableName() string {
	return "pricing_plans"
}

func (m *PricingPlan) RemoveRelationships() {
	m.Policies = nil
	// m.PlanAndResources = nil
	m.Resources = nil
	m.Subscriptions = nil
}
