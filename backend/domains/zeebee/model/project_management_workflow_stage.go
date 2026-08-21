package model

import "zeebee-backend/domains/zeebee/constants"

type ProjectManagementWorkflowStage struct {
	Base
	WorkflowId int64                                 `gorm:"column:workflow_id;not null;default:0" json:"workflow_id"`
	Name       string                                `gorm:"column:name;type:varchar(100);not null;default:''" json:"name"`
	Phase      constants.ProjectManagementStagePhase `gorm:"column:phase;type:varchar(32);not null;default:'Pending'" json:"phase"`
	Color      string                                `gorm:"column:color;type:varchar(32);not null;default:''" json:"color"`
	SortOrder  int                                   `gorm:"column:sort_order;not null;default:0" json:"sort_order"`

	// Relationships
	Transitions []ProjectManagementWorkflowTransition `gorm:"foreignKey:FromStageId" json:"transitions"`
}

func (ProjectManagementWorkflowStage) TableName() string { return "project_management_workflow_stages" }


// RemoveRelationships removes relationship data to avoid circular references.
func (m *ProjectManagementWorkflowStage) RemoveRelationships() {
	m.Transitions = nil
}