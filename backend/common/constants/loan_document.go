package constants

import (
	"github.com/duke-git/lancet/v2/slice"

	"zeebee-backend/common/enums"
)

type LoanDocumentType string

const LoanDocumentTypeName = "loan_document_type"

const (
	LoanDocumentTypeApplyForm        LoanDocumentType = "ApplyForm"
	LoanDocumentTypeIdentification   LoanDocumentType = "Identification"
	LoanDocumentTypeIncome           LoanDocumentType = "Income"
	LoanDocumentTypeProperty         LoanDocumentType = "Property"
	LoanDocumentTypeAssets           LoanDocumentType = "Assets"
	LoanDocumentTypeCredit           LoanDocumentType = "Credit"
	LoanDocumentTypeMiscellaneous    LoanDocumentType = "Miscellaneous"
	LoanDocumentTypeInternalDocument LoanDocumentType = "InternalDocument"
)

func (c LoanDocumentType) IsValid() bool {
	arr := []LoanDocumentType{
		LoanDocumentTypeApplyForm,
		LoanDocumentTypeIdentification,
		LoanDocumentTypeIncome,
		LoanDocumentTypeProperty,
		LoanDocumentTypeAssets,
		LoanDocumentTypeCredit,
		LoanDocumentTypeMiscellaneous,
		LoanDocumentTypeInternalDocument,
	}
	_, ok := slice.FindBy(arr, func(index int, item LoanDocumentType) bool {
		return c == item
	})
	return ok
}

// uploader type
type LoanDocumentUploaderType string

const LoanDocumentUploaderTypeName = "loan_document_uploader_type"
const (
	LoanDocumentUploaderTypeBorrower LoanDocumentUploaderType = "Borrower"
	LoanDocumentUploaderTypeLender   LoanDocumentUploaderType = "Lender"
)

// loan  document status: Pending Uploand, Requested, Ready For review
type LoanDocumentStatus string

const LoanDocumentStatusName = "loan_document_status"
const (
	LoanDocumentStatusPendingUpload  LoanDocumentStatus = "PendingUpload"
	LoanDocumentStatusRequested      LoanDocumentStatus = "Requested"
	LoanDocumentStatusReadyForReview LoanDocumentStatus = "ReadyForReview"
)

func init() {
	// register loan document type
	enums.Register(LoanDocumentTypeName, []enums.EnumOpt{
		{
			Label: string(LoanDocumentTypeApplyForm),
			Value: string(LoanDocumentTypeApplyForm),
		},
		{
			Label: string(LoanDocumentTypeIdentification),
			Value: string(LoanDocumentTypeIdentification),
		},
		{
			Label: string(LoanDocumentTypeIncome),
			Value: string(LoanDocumentTypeIncome),
		},
		{
			Label: string(LoanDocumentTypeProperty),
			Value: string(LoanDocumentTypeProperty),
		},
		{
			Label: string(LoanDocumentTypeAssets),
			Value: string(LoanDocumentTypeAssets),
		},
		{
			Label: string(LoanDocumentTypeCredit),
			Value: string(LoanDocumentTypeCredit),
		},
		{
			Label: string(LoanDocumentTypeMiscellaneous),
			Value: string(LoanDocumentTypeMiscellaneous),
		},
		{
			Label: string(LoanDocumentTypeInternalDocument),
			Value: string(LoanDocumentTypeInternalDocument),
		},
	})
	// register loan document uploader type
	enums.Register(LoanDocumentUploaderTypeName, []enums.EnumOpt{
		{
			Label: string(LoanDocumentUploaderTypeBorrower),
			Value: string(LoanDocumentUploaderTypeBorrower),
		},
		{
			Label: string(LoanDocumentUploaderTypeLender),
			Value: string(LoanDocumentUploaderTypeLender),
		},
	})

	// register loan document status
	enums.Register(LoanDocumentStatusName, []enums.EnumOpt{
		{
			Label: string(LoanDocumentStatusPendingUpload),
			Value: string(LoanDocumentStatusPendingUpload),
		},
		{
			Label: string(LoanDocumentStatusRequested),
			Value: string(LoanDocumentStatusRequested),
		},
		{
			Label: string(LoanDocumentStatusReadyForReview),
			Value: string(LoanDocumentStatusReadyForReview),
		},
	})
}
