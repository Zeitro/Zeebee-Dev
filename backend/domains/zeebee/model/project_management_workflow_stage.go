package model

import "zeebee-backend/domains/zeebee/constants"

type ProjectManagementWorkflowStage struct {
	Base
	WorkflowId int64                                 `gorm:"column:workflow_id;not null;default:0" json:"workflow_id"`
	Name       string                                `gorm:"column:name;type:varchar(100);not null;default:''" json:"name"`
	Phase      constants.ProjectManagementStagePhase `gorm:"column:phase;type:varchar(32);not null;default:'Pending'" json:"phase"`
	SortOrder  int                                   `gorm:"column:sort_order;not null;default:0" json:"sort_order"`
}

func (ProjectManagementWorkflowStage) TableName() string { return "project_management_workflow_stages" }
