package model

import (
	"gorm.io/datatypes"

	"zeebee-backend/domains/zeebee/constants"
)

type Operation struct {
	Base
	EntityId     int64                      `gorm:"column:entity_id;not null;default:0" json:"entity_id"`
	EntityType   constants.OperationEntityType `gorm:"column:entity_type;type:varchar(100);not null;default:''" json:"entity_type"`
	OperatorId   int64                      `gorm:"column:operator_id;not null;default:0" json:"operator_id"`
	OperatorName string                     `gorm:"column:operator_name;type:varchar(50);not null;default:''" json:"operator_name"`
	OperatorType string                     `gorm:"column:operator_type;type:varchar(50);not null;default:''" json:"operator_type"`
	Operation    constants.OperationAction  `gorm:"column:operation;type:varchar(50);not null;default:''" json:"operation"`
	Note         string                     `gorm:"column:note;type:varchar(255);not null;default:''" json:"note"`
	Metadata     datatypes.JSON             `gorm:"column:metadata;type:jsonb;not null;default:'{}'" json:"metadata"`
	OrganizationId int64                    `gorm:"column:organization_id;not null;default:0" json:"organization_id"`
}

func (m Operation) TableName() string {
	return "operations"
}
