package model

import "zeebee-backend/domains/zeebee/constants"

type ProjectManagementTemplate struct {
	Base
	TemplateType    constants.ProjectManagementTemplateType `gorm:"column:template_type;type:varchar(32);not null" json:"template_type"`
	OrganizationId  int64                                   `gorm:"column:organization_id;not null;default:0" json:"organization_id"`
	Name            string                                  `gorm:"column:name;type:varchar(100);not null;default:''" json:"name"`
	Description     string                                  `gorm:"column:description;type:varchar(500);not null;default:''" json:"description"`
	Status          constants.CommonStatus                  `gorm:"column:status;type:varchar(32);not null;default:'Active'" json:"status"`
	IsDefault       bool                                    `gorm:"column:is_default;not null;default:false" json:"is_default"`
	SortOrder       int                                     `gorm:"column:sort_order;not null;default:0" json:"sort_order"`
	SourceProjectId int64                                   `gorm:"column:source_project_id;not null;default:0" json:"source_project_id"`
	// relationships
	Stages []*ProjectManagementStage `gorm:"polymorphic:Entity;polymorphicType:EntityType;polymorphicId:EntityId;polymorphicValue:Template" json:"stages"`
}

func (ProjectManagementTemplate) TableName() string { return "project_management_templates" }

func (m *ProjectManagementTemplate) RemoveRelationships() {
	m.Stages = nil
}
