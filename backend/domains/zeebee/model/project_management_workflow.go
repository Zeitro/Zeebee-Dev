package model

import (
	"zeebee-backend/domains/zeebee/constants"
)

type ProjectManagementWorkflow struct {
	Base
	Name               string                                      `gorm:"column:name;type:varchar(100);not null;default:''" json:"name"`
	Description        string                                      `gorm:"column:description;type:varchar(500);not null;default:''" json:"description"`
	IsDefault          bool                                        `gorm:"column:is_default;not null;default:false" json:"is_default"`
	EntityType         constants.ProjectManagementOwnerType        `gorm:"column:entity_type;type:varchar(32);not null;default:'Project'" json:"entity_type"`
	EntityId           int64                                       `gorm:"column:entity_id;not null;default:0" json:"entity_id"`
	SourceWorkflowId   int64                                       `gorm:"column:source_workflow_id;not null;default:0" json:"source_workflow_id"`
	SourceWorkflowName string                                      `gorm:"column:source_workflow_name;type:varchar(100);not null;default:''" json:"source_workflow_name"`
	// relationship
	Stages             []*ProjectManagementStage                   `gorm:"many2many:project_management_workflow_stages;joinForeignKey:WorkflowId;joinReferences:StageId" json:"stages"`
	Transitions        []*ProjectManagementWorkflowTransition `gorm:"foreignKey:WorkflowId;references:Id" json:"transitions"`
}

func (ProjectManagementWorkflow) TableName() string { return "project_management_workflows" }

// RemoveRelationships removes relationship data to avoid circular references.
func (m *ProjectManagementWorkflow) RemoveRelationships() {
	m.Stages = nil
	m.Transitions = nil
}
