package dto

import (
	"database/sql"

	"gorm.io/datatypes"

	"zeebee-backend/domains/zeebee/constants"
)

type OperationCreateDto struct {
	EntityId       int64                         `json:"entity_id"`
	EntityType     constants.OperationEntityType `json:"entity_type"`
	OperatorId     int64                         `json:"operator_id"`
	OperatorName   string                        `json:"operator_name"`
	OperatorType   constants.UserType            `json:"operator_type"`
	Operation      constants.OperationAction     `json:"operation"`
	Note           string                        `json:"note,omitempty"`
	Metadata       datatypes.JSON                `json:"metadata,omitempty"`
	OrganizationId int64                         `json:"organization_id"`
}

type OperationStatsRequestDto struct {
	EntityId   int64                         `json:"entity_id"`
	EntityType constants.OperationEntityType `json:"entity_type"`
	StartDate  *sql.NullTime                 `json:"start_date"`
	EndDate    *sql.NullTime                 `json:"end_date"`
}

type OperationListRequestDto struct {
	EntityId   int64                         `json:"entity_id"`
	EntityType constants.OperationEntityType `json:"entity_type"`
	StartDate  *sql.NullTime                 `json:"start_date"`
	EndDate    *sql.NullTime                 `json:"end_date"`
	Operation  constants.OperationAction     `json:"operation"`
}

type OperationPageFilterDto struct {
	TaskId        int64
	OperationType constants.OperationAction
	EntityType    constants.OperationEntityType
	EntityId      int64
}

type OperationPageRequestDto struct {
	Page   int64
	Limit  int64
	Filter OperationPageFilterDto
}
