package model

import (
	"zeebee-backend/domains/zeebee/constants"
)

type ProjectManagementWorkflow struct {
	Base
	Kind               constants.ProjectManagementWorkflowKind `gorm:"column:kind;type:varchar(32);not null;default:'OrganizationTemplate'" json:"kind"`
	OrganizationId     int64                                   `gorm:"column:organization_id;not null;default:0" json:"organization_id"`
	ProjectId          int64                                   `gorm:"column:project_id;not null;default:0" json:"project_id"`
	Name               string                                  `gorm:"column:name;type:varchar(100);not null;default:''" json:"name"`
	Description        string                                  `gorm:"column:description;type:varchar(500);not null;default:''" json:"description"`
	IsDefault          bool                                    `gorm:"column:is_default;not null;default:false" json:"is_default"`
	SourceTemplateId   int64                                   `gorm:"column:source_template_id;not null;default:0" json:"source_template_id"`
	SourceTemplateName string                                  `gorm:"column:source_template_name;type:varchar(100);not null;default:''" json:"source_template_name"`
	// relationships
	Stages      []*ProjectManagementWorkflowStage      `gorm:"foreignKey:WorkflowId" json:"stages"`
	Transitions []*ProjectManagementWorkflowTransition `gorm:"foreignKey:WorkflowId" json:"transitions"`
}

func (ProjectManagementWorkflow) TableName() string { return "project_management_workflows" }

// RemoveRelationships removes relationship data to avoid circular references.
func (m *ProjectManagementWorkflow) RemoveRelationships() {
	m.Stages = nil
	m.Transitions = nil
}
