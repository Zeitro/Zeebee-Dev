package model

import (
	"database/sql"

	"zeebee-backend/domains/zeebee/constants"
)

type ProjectManagementTemplate struct {
	Base
	TemplateType    constants.ProjectManagementTemplateType `gorm:"column:template_type;type:varchar(32);not null" json:"template_type"`
	OrganizationId  int64                                   `gorm:"column:organization_id;not null;default:0" json:"organization_id"`
	Name            string                                  `gorm:"column:name;type:varchar(100);not null;default:''" json:"name"`
	Description     string                                  `gorm:"column:description;type:varchar(500);not null;default:''" json:"description"`
	Status          string                                  `gorm:"column:status;type:varchar(32);not null;default:'Active'" json:"status"`
	IsDefault       bool                                    `gorm:"column:is_default;not null;default:false" json:"is_default"`
	SortOrder       int                                     `gorm:"column:sort_order;not null;default:0" json:"sort_order"`
	ProjectType     constants.ProjectManagementProjectType  `gorm:"column:project_type;type:varchar(32)" json:"project_type"`
	StartAt         *sql.NullTime                           `gorm:"column:start_at;type:date" json:"start_at"`
	EndAt           *sql.NullTime                           `gorm:"column:end_at;type:date" json:"end_at"`
	PeriodValue     int                                     `gorm:"column:period_value" json:"period_value"`
	PeriodUnit      constants.ProjectManagementPeriodUnit   `gorm:"column:period_unit;type:varchar(16)" json:"period_unit"`
	SourceProjectId int64                                   `gorm:"column:source_project_id;not null;default:0" json:"source_project_id"`
	Stages          []*ProjectManagementStage               `gorm:"polymorphic:Entity;polymorphicType:EntityType;polymorphicId:EntityId;polymorphicValue:Template" json:"stages"`
}

func (ProjectManagementTemplate) TableName() string { return "project_management_templates" }
