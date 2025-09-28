package types

import "zeebee-backend/common/constants"

type LoanOperationStatus struct {
	Stage         constants.LoanStage         `json:"stage"`
	PendingStatus constants.LoanPendingStatus `json:"pending_status"`
}

type LoanOperationMetadata struct {
	Comment string `json:"comment"`
}
