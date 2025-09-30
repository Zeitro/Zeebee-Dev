package types

import "zeebee-backend/common/constants"

type LoanOperationStatus struct {
	Stage         constants.LoanStage         `json:"stage"`
	PendingStatus constants.LoanPendingStatus `json:"pending_status"`
}
type LoanDocument struct {
	Category constants.LoanDocumentCategory `json:"category"`
	FileType constants.LoanDocumentType     `json:"file_type"`
	Status   constants.LoanDocumentStatus   `json:"status"`
	Files    []*OperationFile               `json:"files,omitempty"`
}

type OperationFile struct {
	EntityId   int64
	EntityType constants.FileEntityType
	FileUrl    string
}

type LoanOperationMetadata struct {
	Comment       string          `json:"comment,omitempty"`
	LoanDocuments []*LoanDocument `json:"loan_documents,omitempty"`
	// Files         []*OperationFile `json:"files,omitempty"`
}
