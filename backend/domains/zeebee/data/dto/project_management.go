package dto

import (
	"time"

	"zeebee-backend/domains/zeebee/constants"
)

type ProjectManagementActor struct {
	OrganizationId     int64
	OrganizationUserId int64
	Name               string
}

type ProjectManagementWorkflowStageDto struct {
	Id        int64
	Name      string
	Phase     constants.ProjectManagementStagePhase
	SortOrder int
}

type ProjectManagementWorkflowTransitionDto struct {
	Id                int64
	FromStageId       int64
	ToStageId         int64
	AllowedActorTypes []constants.ProjectManagementMemberRole
}

type ProjectManagementWorkflowDto struct {
	Id                 int64
	Name               string
	Description        string
	Enabled            bool
	IsDefault          bool
	SourceTemplateId   int64
	SourceTemplateName string
	Stages             []ProjectManagementWorkflowStageDto
	Transitions        []ProjectManagementWorkflowTransitionDto
}

type ProjectManagementWorkflowFilterDto struct {
	Keyword string
	Enabled *bool
}

type ProjectManagementWorkflowPageRequestDto struct {
	Page   int64
	Limit  int64
	Filter ProjectManagementWorkflowFilterDto
}

type ProjectManagementTaskTypeSnapshotDto struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Color       string `json:"color"`
}

type ProjectManagementTaskTypeDto struct {
	Id          int64
	Name        string
	Description string
	Color       string
	Enabled     bool
}

type ProjectManagementTaskTypeFilterDto struct {
	Keyword string
	Enabled *bool
}

type ProjectManagementTaskTypePageRequestDto struct {
	Page   int64
	Limit  int64
	Filter ProjectManagementTaskTypeFilterDto
}

type ProjectManagementProjectDto struct {
	Id             int64
	OrganizationId int64
	EntityId       int64
	EntityType     string
	Name           string
	Description    string
	CoverUrl       string
	OwnerId        int64
	OwnerName      string
}

type ProjectManagementProjectCreateDto struct {
	ProjectManagementProjectDto
	WorkflowTemplateId int64
}

type ProjectManagementProjectUpdateDto struct {
	Name        string
	Description string
	CoverUrl    string
}

type ProjectManagementProjectFilterDto struct {
	Keyword    string
	EntityId   int64
	EntityType string
}

type ProjectManagementProjectPageRequestDto struct {
	Page   int64
	Limit  int64
	Filter ProjectManagementProjectFilterDto
}

type ProjectManagementMemberDto struct {
	ProjectId          int64
	OrganizationUserId int64
	Role               constants.ProjectManagementMemberRole
}

type ProjectManagementMemberDetailDto struct {
	Id                 int64
	ProjectId          int64
	OrganizationUserId int64
	Role               constants.ProjectManagementMemberRole
	FullName           string
	Email              string
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

type ProjectManagementMemberFilterDto struct {
	Keyword string
	Role    constants.ProjectManagementMemberRole
}

type ProjectManagementMemberPageRequestDto struct {
	Page   int64
	Limit  int64
	Filter ProjectManagementMemberFilterDto
}

type ProjectManagementTaskDto struct {
	Id             int64
	OrganizationId int64
	ProjectId      int64
	ParentTaskId   int64
	StageId        int64
	CreatorId      int64
	CreatorName    string
	Title          string
	Description    string
	Priority       string
	StartAt        *time.Time
	EndAt          *time.Time
	TaskType       *ProjectManagementTaskTypeSnapshotDto
}

type ProjectManagementTaskOverdueDto struct {
	Overdue        bool
	OverdueSeconds int64
}

type ProjectManagementTaskCreateDto struct {
	ProjectManagementTaskDto
	TaskTypeId  int64
	AssigneeIds []int64
}

type ProjectManagementTaskUpdateDto struct {
	Title       string
	Description string
	Priority    string
	StartAt     *time.Time
	EndAt       *time.Time
	TaskTypeId  int64
}

type ProjectManagementTaskAssigneeDto struct {
	TaskId               int64
	OrganizationUserId   int64
	OrganizationUserName string
}

type ProjectManagementTaskDependencyDto struct {
	Id               int64
	TaskId           int64
	DependencyTaskId int64
	DependencyType   constants.ProjectManagementDependencyType
	DependencyRule   constants.ProjectManagementDependencyRule
}

type ProjectManagementTaskDependencyPageRequestDto struct {
	Page        int64
	Limit       int64
	TaskId      int64
	Predecessor bool
}

type ProjectManagementAttachmentDto struct {
	FileId      int64
	TaskId      int64
	FileName    string
	FileUrl     string
	StorageType string
}

type ProjectManagementAttachmentCreateDto struct {
	TaskId      int64
	FileName    string
	FileUrl     string
	StorageType constants.StorageType
}

type ProjectManagementAttachmentPageRequestDto struct {
	Page  int64
	Limit int64
}

type ProjectManagementTaskFilterDto struct {
	ProjectId  int64
	StageId    int64
	TaskTypeId int64
	AssigneeId int64
	Keyword    string
	Overdue    *bool
}

type ProjectManagementTaskPageRequestDto struct {
	Page   int64
	Limit  int64
	Filter ProjectManagementTaskFilterDto
}

type ProjectManagementMyTaskPageRequestDto struct {
	Page   int64
	Limit  int64
	Filter ProjectManagementTaskFilterDto
}

type ProjectManagementOperationFilterDto struct {
	ProjectId     int64
	TaskId        int64
	OperationType string
}

type ProjectManagementOperationDto struct {
	ProjectId    int64
	TaskId       int64
	OperatorId   int64
	OperatorName string
	Operation    string
	Metadata     []byte
}

type ProjectManagementOperationPageRequestDto struct {
	Page   int64
	Limit  int64
	Filter ProjectManagementOperationFilterDto
}
