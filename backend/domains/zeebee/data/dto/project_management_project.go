package dto

import "zeebee-backend/domains/zeebee/constants"

type ProjectManagementProjectDto struct {
	Id             int64
	OrganizationId int64
	EntityId       int64
	EntityType     constants.ProjectManagementProjectEntityType
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
	EntityType constants.ProjectManagementProjectEntityType
}

type ProjectManagementProjectPageRequestDto struct {
	Page   int64
	Limit  int64
	Filter ProjectManagementProjectFilterDto
}
