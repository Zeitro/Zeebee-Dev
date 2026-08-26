package dto

import (
	"database/sql"

	"zeebee-backend/domains/zeebee/constants"
)

type ProjectManagementTaskDto struct {
	Id             int64
	OrganizationId int64
	ProjectId      int64
	ParentTaskId   int64
	StageId        int64
	Phase          constants.ProjectManagementStagePhase
	CreatorId      int64
	CreatorName    string
	Title          string
	Description    string
	Priority       constants.ProjectManagementPriority
	StartAt        *sql.NullTime
	EndAt          *sql.NullTime
	TaskType       *ProjectManagementTaskTypeSnapshotDto
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
	Title                       string
	Description                 string
	Priority                    constants.ProjectManagementPriority
	StartAt                     *sql.NullTime
	EndAt                       *sql.NullTime
	TaskTypeId                  int64
	AssigneeOrganizationUserIds []int64
	Dependencies                []*ProjectManagementTaskDependencyDto
}

type ProjectManagementTaskAssigneeDto struct {
	TaskId               int64
	StageId              int64
	OrganizationUserId   int64
	OrganizationUserName string
}

type ProjectManagementTaskFilterDto struct {
	ProjectId  int64
	StageId    int64
	Phase      []constants.ProjectManagementStagePhase
	TaskTypeId int64
	AssigneeId int64
	Keyword    string
	Overdue    *bool
}

type ProjectManagementTaskCandidateDto struct {
	ProjectId int64
	TaskId    int64
	Keyword   string
}

type ProjectManagementTaskCandidateEntityDto struct {
	Id    int64
	Title string
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
