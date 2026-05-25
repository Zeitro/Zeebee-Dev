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
	OperationActionBorrowerPropertyTransactionUpdateBaseInfo  OperationAction = "Transaction base info updated"
	OperationActionBorrowerPropertyTransactionUpdateStatus    OperationAction = "Transaction status changed"
	OperationActionBorrowerPropertyTransactionUpdateChecklist OperationAction = "Transaction checklist changed"
	OperationActionBorrowerPropertyTransactionAddDocument     OperationAction = "Transaction document added"
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

func (s StorageType) IsValid() bool {
	switch s {
	case StorageTypeAwsS3, StorageTypeGoogleDrive:
		return true
	}
	return false
}

type UploadType string

const (
	UploadTypeManual      UploadType = "Manual"
	UploadTypeGoogleDrive UploadType = "Google Drive"
)

func (s UploadType) IsValid() bool {
	switch s {
	case UploadTypeManual, UploadTypeGoogleDrive:
		return true
	}
	return false
}
