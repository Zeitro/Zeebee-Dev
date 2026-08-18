package model

import (
	"github.com/lib/pq"

	"zeebee-backend/domains/zeebee/constants"
)

type ProjectManagementWorkflow struct {
	Base
	Kind               constants.ProjectManagementWorkflowKind `gorm:"column:kind;type:varchar(32);not null" json:"kind"`
	OrganizationId     int64                                   `gorm:"column:organization_id;not null;default:0" json:"organization_id"`
	ProjectId          int64                                   `gorm:"column:project_id;not null;default:0" json:"project_id"`
	Name               string                                  `gorm:"column:name;type:varchar(100);not null;default:''" json:"name"`
	Description        string                                  `gorm:"column:description;type:varchar(500);not null;default:''" json:"description"`
	Enabled            bool                                    `gorm:"column:enabled;not null;default:true" json:"enabled"`
	IsDefault          bool                                    `gorm:"column:is_default;not null;default:false" json:"is_default"`
	SourceTemplateId   int64                                   `gorm:"column:source_template_id;not null;default:0" json:"source_template_id"`
	SourceTemplateName string                                  `gorm:"column:source_template_name;type:varchar(100);not null;default:''" json:"source_template_name"`
	Stages             []*ProjectManagementWorkflowStage       `gorm:"foreignKey:WorkflowId" json:"stages"`
	Transitions        []*ProjectManagementWorkflowTransition  `gorm:"foreignKey:WorkflowId" json:"transitions"`
}

func (ProjectManagementWorkflow) TableName() string { return "project_management_workflows" }

type ProjectManagementWorkflowStage struct {
	Base
	WorkflowId int64                                 `gorm:"column:workflow_id;not null" json:"workflow_id"`
	Name       string                                `gorm:"column:name;type:varchar(100);not null;default:''" json:"name"`
	Phase      constants.ProjectManagementStagePhase `gorm:"column:phase;type:varchar(32);not null" json:"phase"`
	SortOrder  int                                   `gorm:"column:sort_order;not null" json:"sort_order"`
}

func (ProjectManagementWorkflowStage) TableName() string { return "project_management_workflow_stages" }

type ProjectManagementWorkflowTransition struct {
	Base
	WorkflowId        int64          `gorm:"column:workflow_id;not null" json:"workflow_id"`
	FromStageId       int64          `gorm:"column:from_stage_id;not null" json:"from_stage_id"`
	ToStageId         int64          `gorm:"column:to_stage_id;not null" json:"to_stage_id"`
	AllowedActorTypes pq.StringArray `gorm:"column:allowed_actor_types;type:text[];not null;default:'{}'" json:"allowed_actor_types"`
}

func (ProjectManagementWorkflowTransition) TableName() string {
	return "project_management_workflow_transitions"
}
