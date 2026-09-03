package dto

import "zeebee-backend/domains/zeebee/constants"

type ProjectManagementTaskDependencyDto struct {
	Id               int64
	TaskId           int64
	DependencyTaskId int64
	DependencyType   constants.ProjectManagementDependencyType
	DependencyRule   constants.ProjectManagementDependencyRule
}

type ProjectManagementTaskDependencyFilterDto struct {
	TaskId      int64
	Predecessor bool
}

type ProjectManagementTaskDependencyPageRequestDto struct {
	Page   int64
	Limit  int64
	Filter ProjectManagementTaskDependencyFilterDto
}
