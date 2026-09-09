package dto

import (
	"zeebee-backend/domains/zeebee/constants"
)

type LoanOperationCreateDto struct {
	LoanId        int64                               `json:"loan_id"`
	OperatorId    int64                               `json:"operator_id"`
	From          LoanOperationStatus                 `json:"from"`
	To            LoanOperationStatus                 `json:"to"`
	Action        constants.LoanOperationAction       `json:"action"`
	OperationType constants.LoanOperationOperatorType `json:"operation_type"`
	Metadata      *LoanOperationMetadata              `json:"metadata"`
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
	Comment                      string                   `json:"comment,omitempty"`
	LoanDocuments                []*OperationLoanDocument `json:"loan_documents,omitempty"`
	OriginalLenderOrganizationId int64                    `json:"original_lender_organization_id,omitempty"`
	NewLenderOrganizationId      int64                    `json:"new_lender_organization_id,omitempty"`
	// Files         []*OperationFile `json:"files,omitempty"`
}
