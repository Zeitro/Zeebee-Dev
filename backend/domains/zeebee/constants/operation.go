package constants

type OperationEntityType string

const (
	OperationEntityTypeBorrowerPropertyTransaction OperationEntityType = "BorrowerPropertyTransaction"
	OperationEntityTypeLoan                        OperationEntityType = "Loan"
)

type OperationAction string

const (
	// Borrower property transaction actions
	OperationActionBorrowerPropertyTransactionCreate          OperationAction = "Transaction created"
	OperationActionBorrowerPropertyTransactionUpdateBaseInfo  OperationAction = "Basic information updated"
	OperationActionBorrowerPropertyTransactionUpdateStatus    OperationAction = "Status changed"
	OperationActionBorrowerPropertyTransactionUpdateChecklist OperationAction = "Checklist changed"
	OperationActionBorrowerPropertyTransactionAddDocument     OperationAction = "Document created"
	OperationActionBorrowerPropertyTransactionDelete          OperationAction = "Transaction deleted"
	OperationActionBorrowerPropertyTransactionUploadDocument  OperationAction = "Documents uploaded"
	OperationActionBorrowerPropertyTransactionDeleteDocument  OperationAction = "Document deleted"

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
