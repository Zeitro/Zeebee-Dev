package dto

import "zeebee-backend/domains/zeebee/constants"

type ProjectManagementOperationFilterDto struct {
	ProjectId     int64
	TaskId        int64
	OperationType constants.OperationAction
}

type ProjectManagementOperationDto struct {
	ProjectId    int64
	TaskId       int64
	OperatorId   int64
	OperatorName string
	Operation    constants.OperationAction
	Metadata     []byte
}

type ProjectManagementOperationPageRequestDto struct {
	Page   int64
	Limit  int64
	Filter ProjectManagementOperationFilterDto
}
