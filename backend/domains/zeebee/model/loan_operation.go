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

type LoanOperationStatus struct {
	Stage         constants.LoanStage         `json:"stage"`
	PendingStatus constants.LoanPendingStatus `json:"pending_status"`
}

type OperationLoanDocument struct {
	Category constants.LoanDocumentCategory `json:"category"`
	FileType constants.LoanDocumentFileType `json:"file_type"`
	Status   constants.LoanDocumentStatus   `json:"status"`
	Files    []*OperationFile               `json:"files,omitempty"`
}

type OperationFile struct {
	EntityId   int64                    `json:"entity_id"`
	EntityType constants.FileEntityType `json:"entity_type"`
	FileUrl    string                   `json:"file_url"`
}

type LoanOperationMetadata struct {
	Comment       string                   `json:"comment,omitempty"`
	LoanDocuments []*OperationLoanDocument `json:"loan_documents,omitempty"`
	// Files         []*OperationFile `json:"files,omitempty"`
}
