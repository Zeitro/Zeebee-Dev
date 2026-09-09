package model

import "zeebee-backend/domains/zeebee/constants"

type ProjectManagementWorkflowTransition struct {
	Base
	WorkflowId        int64                                   `gorm:"column:workflow_id;not null;default:0" json:"workflow_id"`
	Name              string                                  `gorm:"column:name;not null;default:''" json:"name"`
	FromStageId       int64                                   `gorm:"column:from_stage_id;not null;default:0" json:"from_stage_id"`
	ToStageId         int64                                   `gorm:"column:to_stage_id;not null;default:0" json:"to_stage_id"`
	AllowedActorTypes []constants.ProjectManagementMemberRole `gorm:"column:allowed_actor_types;type:jsonb;not null;default:'[]';serializer:json" json:"allowed_actor_types"`
	// relationships
	Workflow  *ProjectManagementWorkflow      `gorm:"foreignKey:WorkflowId;references:Id" json:"workflow,omitempty"`
	FromStage *ProjectManagementWorkflowStage `gorm:"foreignKey:FromStageId;references:Id" json:"from_stage,omitempty"`
	ToStage   *ProjectManagementWorkflowStage `gorm:"foreignKey:ToStageId;references:Id" json:"to_stage,omitempty"`
}

func (ProjectManagementWorkflowTransition) TableName() string {
	return "project_management_workflow_transitions"
}

// RemoveRelationships removes relationship data to avoid circular references.
func (m *ProjectManagementWorkflowTransition) RemoveRelationships() {
	m.Workflow = nil
	m.FromStage = nil
	m.ToStage = nil
}
