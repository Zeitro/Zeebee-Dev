package dto

import (
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
