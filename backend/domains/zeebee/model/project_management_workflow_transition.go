package model

import "zeebee-backend/domains/zeebee/constants"

// ProjectManagementWorkflowStageTransition stores owner-stage IDs directly.
type ProjectManagementWorkflowStageTransition struct {
	Base
	WorkflowId        int64                                   `gorm:"column:workflow_id;not null;default:0" json:"workflow_id"`
	Name              string                                  `gorm:"column:name;type:varchar(100);not null;default:''" json:"name"`
	FromStageId       int64                                   `gorm:"column:from_stage_id;not null;default:0" json:"from_stage_id"`
	ToStageId         int64                                   `gorm:"column:to_stage_id;not null;default:0" json:"to_stage_id"`
	AllowedActorTypes []constants.ProjectManagementMemberRole `gorm:"column:allowed_actor_types;type:jsonb;not null;default:'[]';serializer:json" json:"allowed_actor_types"`
}

func (ProjectManagementWorkflowStageTransition) TableName() string {
	return "project_management_workflow_transitions"
}
