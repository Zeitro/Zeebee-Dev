package constants

import "zeebee-backend/common/enums"

type LoanOperationAction string

const LoanOperationActionName = "loan_operation_action"

const (
	LoanOperationActionCreate  LoanOperationAction = "Create"
	LoanOperationActionUpdate  LoanOperationAction = "Update"
	LoanOperationActionReject  LoanOperationAction = "Reject"
	LoanOperationActionApprove LoanOperationAction = "Approve"
	LoanOperationActionComment LoanOperationAction = "Comment"
	LoanOperationActionDelete  LoanOperationAction = "Delete"
)

type LoanOperationOperatorType string

const LoanOperationOperatorTypeName = "loan_operation_operator_type"
const (
	LoanOperationOperatorTypeBorrower LoanOperationOperatorType = "Borrower"
	LoanOperationOperatorTypeLender   LoanOperationOperatorType = "Lender"
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
	})
}
