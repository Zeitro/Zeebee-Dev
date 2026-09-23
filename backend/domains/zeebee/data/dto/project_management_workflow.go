package dto

import "zeebee-backend/domains/zeebee/constants"

type ProjectManagementWorkflowStageDto struct {
	Id          int64
	Name        string
	Description string
	Phase       constants.ProjectManagementStagePhase
	Color       string
	SortOrder   int
}

type ProjectManagementWorkflowTransitionDto struct {
	Id                int64
	Name              string
	FromStageId       int64
	ToStageId         int64
	AllowedActorTypes []constants.ProjectManagementMemberRole
}

type ProjectManagementWorkflowUpdateTransitionDto struct {
	Id                int64
	Name              string
	FromStageId       int64
	ToStageId         int64
	AllowedActorTypes []constants.ProjectManagementMemberRole
}

type ProjectManagementWorkflowUpdateStageDto struct {
	Id          int64
	Name        string
	Description string
	Phase       constants.ProjectManagementStagePhase
	Color       string
	SortOrder   int
	Transitions []ProjectManagementWorkflowUpdateTransitionDto
}

type ProjectManagementWorkflowUpdateDto struct {
	Name        string
	Description string
	Stages      []ProjectManagementWorkflowUpdateStageDto
}

type ProjectManagementWorkflowDto struct {
	Id                 int64
	Name               string
	Description        string
	IsDefault          bool
	SourceTemplateId   int64
	SourceTemplateName string
	EntityType         constants.ProjectManagementOwnerType
	EntityId           int64
	SourceWorkflowId   int64
	SourceWorkflowName string
	Stages             []ProjectManagementWorkflowStageDto
	Transitions        []ProjectManagementWorkflowTransitionDto
}
