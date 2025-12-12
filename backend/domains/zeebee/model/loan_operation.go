package model

import (
	"zeebee-backend/domains/zeebee/constants"
)

type LoanOperation struct {
	Base
	LoanId       int64                               `gorm:"column:loan_id;not null;default:0" json:"loan_id"`
	OperatorId   int64                               `gorm:"column:operator_id;not null;default:0" json:"operator_id"`
	OperatorName string                              `gorm:"column:operator_name;size:50;not null;default:''" json:"operator_name"`
	OperatorType constants.LoanOperationOperatorType `gorm:"column:operator_type;size:50;not null;default:''" json:"operator_type"`
	Action       constants.LoanOperationAction       `gorm:"column:action;size:50;not null;default:''" json:"action"`
	From         LoanOperationStatus                 `gorm:"column:from;type:jsonb;not null;default:'{}';serializer:json" json:"from"`
	To           LoanOperationStatus                 `gorm:"column:to;type:jsonb;not null;default:'{}';serializer:json" json:"to"`
	Metadata     *LoanOperationMetadata              `gorm:"column:metadata;type:jsonb;not null;default:'{}';serializer:json" json:"metadata"`
}

func (m LoanOperation) TableName() string {
	return "loan_operations"
}
