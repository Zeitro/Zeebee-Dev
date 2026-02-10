package dto

import (
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data"
)

type LoanOperationCreateDto struct {
	LoanId        int64                               `json:"loan_id"`
	OperatorId    int64                               `json:"operator_id"`
	From          data.LoanOperationStatus            `json:"from"`
	To            data.LoanOperationStatus            `json:"to"`
	Action        constants.LoanOperationAction       `json:"action"`
	OperationType constants.LoanOperationOperatorType `json:"operation_type"`
	Metadata      *data.LoanOperationMetadata         `json:"metadata"`
}
