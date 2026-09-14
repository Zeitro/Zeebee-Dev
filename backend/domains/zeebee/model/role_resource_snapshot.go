package model

import "zeebee-backend/domains/zeebee/constants"

type RoleResourceSnapshot struct {
	Base
	OrganizationId int64                          `gorm:"column:organization_id;type:int4;not null" json:"organization_id"`
	PricingPlanId  int64                          `gorm:"column:pricing_plan_id;type:int4;not null" json:"pricing_plan_id"`
	RoleId         int64                          `gorm:"column:role_id;type:int4;not null" json:"role_id"`
	Resources      []RoleResourceSnapshotResource `gorm:"column:resources;type:jsonb;not null;default:'[]';serializer:json" json:"resources"`
}

type RoleResourceSnapshotResource struct {
	ResourceId int64               `json:"resource_id"`
	DataScope  constants.DataScope `json:"data_scope"`
}

func (RoleResourceSnapshot) TableName() string {
	return "role_resource_snapshots"
}
