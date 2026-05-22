package constants

type OperationEntityType string

const (
	OperationEntityTypeBorrowerPropertyTransaction OperationEntityType = "BorrowerPropertyTransaction"
	OperationEntityTypeLoan                        OperationEntityType = "Loan"
)

type OperationAction string

const (
	// Borrower propperty transaction actions
	OperationActionBorrowerPropertyTransactionCreate          OperationAction = "Transaction Created"
	OperationActionBorrowerPropertyTransactionUpdateBaseInfo  OperationAction = "Transaction base information updated"
	OperationActionBorrowerPropertyTransactionUpdateStatus    OperationAction = "Transaction status changed"
	OperationActionBorrowerPropertyTransactionUpdateChecklist OperationAction = "Transaction checklist changed"
	OperationActionBorrowerPropertyTransactionDelete          OperationAction = "Borrower Property Transaction Delete"
	OperationActionBorrowerPropertyTransactionUploadDocument  OperationAction = "Borrower Property Transaction Upload Document"
	OperationActionBorrowerPropertyTransactionDeleteDocument  OperationAction = "Borrower Property Transaction Delete Document"

	// Loan actions
	OperationActionLoanCreate       OperationAction = "Loan Create"
	OperationActionLoanUpdate       OperationAction = "Loan Update"
	OperationActionLoanStatusChange OperationAction = "Loan Status Change"
)

type StorageType string

const (
	StorageTypeAwsS3       StorageType = "Aws S3"
	StorageTypeGoogleDrive StorageType = "Google Drive"
)

type UploadType string

const (
	UploadTypeManual      UploadType = "Manual"
	UploadTypeGoogleDrive UploadType = "Google Drive"
)
