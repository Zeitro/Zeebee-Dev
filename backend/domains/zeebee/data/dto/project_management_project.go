package dto

import (
	"database/sql"

	"zeebee-backend/domains/zeebee/constants"
)

type ProjectManagementProjectDto struct {
	Id             int64
	OrganizationId int64
	EntityId       int64
	EntityType     constants.ProjectManagementProjectEntityType
	Name           string
	Description    string
	Priority       constants.ProjectManagementPriority
	CoverUrl       string
	OwnerId        int64
	OwnerName      string
	ProjectType    constants.ProjectManagementProjectType
	StartAt        *sql.NullTime
	EndAt          *sql.NullTime
	PeriodValue    int
	PeriodUnit     constants.ProjectManagementPeriodUnit
}

type ProjectManagementProjectCreateDto struct {
	EntityId      int64
	EntityType    constants.ProjectManagementProjectEntityType
	Name          string
	Description   string
	Priority      constants.ProjectManagementPriority
	CoverUrl      string
	ProjectType   constants.ProjectManagementProjectType
	StartAt       *sql.NullTime
	PeriodValue   int
	PeriodUnit    constants.ProjectManagementPeriodUnit
	TemplateId    int64
	CopyAssignees bool
}

type ProjectManagementTemplateAssigneeResultDto struct {
	OrganizationUserId int64
	Reason             string
}

type ProjectManagementProjectCreateAssigneeResultDto struct {
	CopiedOrganizationUserIds []int64
	SkippedAssignees          []*ProjectManagementTemplateAssigneeResultDto
}

type ProjectManagementProjectUpdateDto struct {
	Name        string
	Description string
	Priority    constants.ProjectManagementPriority
	CoverUrl    string
	StartAt     *sql.NullTime
	PeriodValue *int
	PeriodUnit  *constants.ProjectManagementPeriodUnit
}

type ProjectManagementProjectFilterDto struct {
	Keyword            string
	EntityId           int64
	EntityType         constants.ProjectManagementProjectEntityType
	OrganizationUserId int64
}

type ProjectManagementProjectPageRequestDto struct {
	Page   int64
	Limit  int64
	Filter ProjectManagementProjectFilterDto
}
