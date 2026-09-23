package model

import "zeebee-backend/domains/zeebee/constants"

type ProjectManagementStage struct {
	Base
	EntityType  constants.ProjectManagementOwnerType  `gorm:"column:entity_type;type:varchar(32);not null" json:"entity_type"`
	EntityId    int64                                 `gorm:"column:entity_id;not null;default:0" json:"entity_id"`
	Name        string                                `gorm:"column:name;type:varchar(100);not null;default:''" json:"name"`
	Description string                                `gorm:"column:description;type:varchar(500);not null;default:''" json:"description"`
	Phase       constants.ProjectManagementStagePhase `gorm:"column:phase;type:varchar(32);not null" json:"phase"`
	Color       string                                `gorm:"column:color;type:varchar(32);not null;default:''" json:"color"`
	SortOrder   int                                   `gorm:"column:sort_order;not null;default:0" json:"sort_order"`
}

func (ProjectManagementStage) TableName() string { return "project_management_stages" }
