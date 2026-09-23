package dto

import (
	"database/sql"
	"time"

	"zeebee-backend/domains/zeebee/constants"
)

type ProjectManagementTaskDto struct {
	Id             int64
	CreatedAt      time.Time
	UpdatedAt      time.Time
	OrganizationId int64
	EntityType     constants.ProjectManagementOwnerType
	EntityId       int64
	WorkflowId     int64
	ParentTaskId   int64
	StageId        int64
	Phase          constants.ProjectManagementStagePhase
	CreatorId      int64
	CreatorName    string
	Title          string
	Description    string
	Priority       constants.ProjectManagementPriority
	WorkHours      float64
	StartAt        *sql.NullTime
	EndAt          *sql.NullTime
	TaskType       *ProjectManagementTaskTypeSnapshotDto
	Assignees      []*ProjectManagementTaskAssigneeDto
	Attachments    []*ProjectManagementTaskAttachmentDto
	Children       []*ProjectManagementTaskDto
	Dependencies   []*ProjectManagementTaskDependencyDto
	ProjectName    string
	StageName      string
	WorkflowName   string
	OutOfPlan      bool
}

type ProjectManagementTaskOverdueDto struct {
	Overdue        bool
	OverdueSeconds int64
}

type ProjectManagementTaskCreateDto struct {
	ProjectManagementTaskDto
	TaskTypeId                  int64
	AssigneeOrganizationUserIds []int64
}

type ProjectManagementTaskUpdateDto struct {
	StageId                     int64
	WorkflowId                  int64
	Title                       string
	Description                 string
	Priority                    constants.ProjectManagementPriority
	WorkHours                   float64
	StartAt                     *sql.NullTime
	EndAt                       *sql.NullTime
	TaskTypeId                  int64
	AssigneeOrganizationUserIds []int64
	Dependencies                []*ProjectManagementTaskDependencyDto
	Attachments                 []*ProjectManagementTaskAttachmentDto
}

type ProjectManagementTaskAssigneeDto struct {
	Id                       int64
	TaskId                   int64
	StageId                  int64
	OrganizationUserId       int64
	OrganizationUserFullName string
}

type ProjectManagementTaskFilterDto struct {
	ProjectId  int64
	StageId    int64
	Phase      []constants.ProjectManagementStagePhase
	Priority   constants.ProjectManagementPriority
	TaskTypeId int64
	AssigneeId int64
	Keyword    string
	Overdue    *bool
}

type ProjectManagementTaskCandidateDto struct {
	EntityType constants.ProjectManagementOwnerType
	EntityId   int64
	TaskId     int64
	Keyword    string
}

type ProjectManagementTaskCandidateEntityDto struct {
	Id    int64
	Title string
	Phase constants.ProjectManagementStagePhase
}

type ProjectManagementTaskPageRequestDto struct {
	Page   int64
	Limit  int64
	Filter ProjectManagementTaskFilterDto
}

type ProjectManagementTaskViewDto struct {
	ProjectId       int64
	ProjectStartAt  *sql.NullTime
	ProjectEndAt    *sql.NullTime
	Tasks           []*ProjectManagementTaskDto
	TimelineStartAt *sql.NullTime
	TimelineEndAt   *sql.NullTime
}

type ProjectManagementMyTaskPageRequestDto struct {
	Page   int64
	Limit  int64
	Filter ProjectManagementTaskFilterDto
}
