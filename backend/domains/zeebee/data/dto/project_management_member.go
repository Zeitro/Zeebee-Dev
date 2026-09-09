package dto

import (
	"time"

	"zeebee-backend/domains/zeebee/constants"
)

type ProjectManagementMemberDto struct {
	ProjectId          int64
	MemberId           int64
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
	ProjectId int64
	Keyword   string
	Role      constants.ProjectManagementMemberRole
}

type ProjectManagementMemberCandidateDto struct {
	ProjectId int64
	Keyword   string
	Purpose   constants.ProjectManagementMemberCandidatePurpose
}

type ProjectManagementMemberCandidateEntityDto struct {
	OrganizationUserId int64
	FullName           string
	Email              string
}

type ProjectManagementMemberPageRequestDto struct {
	Page   int64
	Limit  int64
	Filter ProjectManagementMemberFilterDto
}
