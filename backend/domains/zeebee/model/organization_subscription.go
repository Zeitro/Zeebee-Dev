package model

import (
	"database/sql"

	"zeebee-backend/domains/zeebee/constants"
)

type OrganizationSubscription struct {
	Base

	Code                string                                   `gorm:"column:code;type:varchar(50);not null;default:'';uniqueIndex;<-:create" json:"code"`
	OrganizationId      int64                                    `gorm:"column:organization_id;type:int4;not null;default:0" json:"organization_id"`
	PricingPlanId       int64                                    `gorm:"column:pricing_plan_id;type:int4;not null;default:0" json:"pricing_plan_id"`
	PricingPlanPolicyId int64                                    `gorm:"column:pricing_plan_policy_id;type:int4;not null;default:0" json:"pricing_plan_policy_id"`
	Status              constants.OrganizationSubscriptionStatus `gorm:"column:status;type:varchar(50);not null;default:'Active'" json:"status"`
	ExpireDate          *sql.NullTime                            `gorm:"column:expire_date;type:date;not null" json:"expire_date"`
	OriginalAmount      float64                                  `gorm:"column:original_amount;type:decimal(10,2);not null;default:0" json:"original_amount"`
	ActualAmount        float64                                  `gorm:"column:actual_amount;type:decimal(10,2);not null;default:0" json:"actual_amount"`
	IsAutoRenew         bool                                     `gorm:"column:is_auto_renew;type:bool;not null;default:false" json:"is_auto_renew"`
}

func (OrganizationSubscription) TableName() string {
	return "organization_subscriptions"
}
