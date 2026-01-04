package types

import "zeebee-backend/domains/zeebee/constants"

type LoanOperationStatus struct {
	Stage         constants.LoanStage         `json:"stage"`
	PendingStatus constants.LoanPendingStatus `json:"pending_status"`
}
type LoanDocument struct {
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
	Comment       string          `json:"comment,omitempty"`
	LoanDocuments []*LoanDocument `json:"loan_documents,omitempty"`
	// Files         []*OperationFile `json:"files,omitempty"`
}
