package constants

import "zeebee-backend/common/enums"

type LoanOperationAction string

const LoanOperationActionName = "loan_operation_action"

const (
	LoanOperationActionCreate                LoanOperationAction = "Create"
	LoanOperationActionUpdate                LoanOperationAction = "Update"
	LoanOperationActionReject                LoanOperationAction = "Reject"
	LoanOperationActionApprove               LoanOperationAction = "Approve"
	LoanOperationActionComment               LoanOperationAction = "Comment"
	LoanOperationActionDelete                LoanOperationAction = "Delete"
	LoanOperationActionCreateDocumentType    LoanOperationAction = "CreateDocumentType"
	LoanOperationActionSendDocumentRequest   LoanOperationAction = "SendDocumentRequest"
	LoanOperationActionDeleteDocumentRequest LoanOperationAction = "DeleteDocumentRequest"
	LoanOperationActionUploadDocumentFile    LoanOperationAction = "UpdateDocumentFile"
	LoanOperationActionDeleteDocumentFile    LoanOperationAction = "DeleteDocumentFile"
)

type LoanOperationOperatorType string

const LoanOperationOperatorTypeName = "loan_operation_operator_type"
const (
	LoanOperationOperatorTypeBorrower LoanOperationOperatorType = "Borrower"
	LoanOperationOperatorTypeLender   LoanOperationOperatorType = "Lender"
	LoanOperationOperatorTypeSystem   LoanOperationOperatorType = "System"
)

func init() {
	// register loan operation action
	enums.Register(LoanOperationActionName, []enums.EnumOpt{
		{
			Label: string(LoanOperationActionCreate),
			Value: string(LoanOperationActionCreate),
		},
		{
			Label: string(LoanOperationActionUpdate),
			Value: string(LoanOperationActionUpdate),
		},
		{
			Label: string(LoanOperationActionReject),
			Value: string(LoanOperationActionReject),
		},
		{
			Label: string(LoanOperationActionApprove),
			Value: string(LoanOperationActionApprove),
		},
		{
			Label: string(LoanOperationActionComment),
			Value: string(LoanOperationActionComment),
		},
		{
			Label: string(LoanOperationActionDelete),
			Value: string(LoanOperationActionDelete),
		},
		{
			Label: string(LoanOperationActionCreateDocumentType),
			Value: string(LoanOperationActionCreateDocumentType),
		},
		{
			Label: string(LoanOperationActionSendDocumentRequest),
			Value: string(LoanOperationActionSendDocumentRequest),
		},
		{
			Label: string(LoanOperationActionDeleteDocumentRequest),
			Value: string(LoanOperationActionDeleteDocumentRequest),
		},
		{
			Label: string(LoanOperationActionUploadDocumentFile),
			Value: string(LoanOperationActionUploadDocumentFile),
		},
		{
			Label: string(LoanOperationActionDeleteDocumentFile),
			Value: string(LoanOperationActionDeleteDocumentFile),
		},
	})

	// register loan operation operator type
	enums.Register(LoanOperationOperatorTypeName, []enums.EnumOpt{
		{
			Label: string(LoanOperationOperatorTypeBorrower),
			Value: string(LoanOperationOperatorTypeBorrower),
		},
		{
			Label: string(LoanOperationOperatorTypeLender),
			Value: string(LoanOperationOperatorTypeLender),
		},
		{
			Label: string(LoanOperationOperatorTypeSystem),
			Value: string(LoanOperationOperatorTypeSystem),
		},
	})
}
