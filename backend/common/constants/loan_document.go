package constants

import (
	"zeebee-backend/common/enums"
)

type LoanDocumentType string

const LoanDocumentTypeName = "loan_document_type"

const (
	// borrower document
	LoanDocumentTypeValidID                = "Valid ID of Borrower"
	LoanDocumentTypePurchaseAgreement      = "Purchase Agreement"
	LoanDocumentTypeEvidenceOfInsurance    = "Evidence of Insurance"
	LoanDocumentTypePreliminaryTitleReport = "Preliminary Title Report/Escrow Contact"
	LoanDocumentTypeREExperience           = "RE Experience"
	LoanDocumentTypeBankStatements         = "Bank Statements for Most Recent Two Months"

	// lender document
	LoanDocumentTypeBorrowerCreditScore                         = "Borrower's Credit Score"
	LoanDocumentTypeBorrowerDefaultFraudRiskViewCriminalRecords = "Borrower's Default/Fraud/Risk/View/Criminal Records And RE Experience"
	LoanDocumentTypePropertyData                                = "Property Data/Historical Property Ownership"
	LoanDocumentTypeTitleReport                                 = "Title Report"
)

// uploader type
type LoanDocumentCategory string

const LoanDocumentCategoryName = "loan_document_category"
const (
	LoanDocumentUploaderTypeBorrower = "Borrower"
	LoanDocumentUploaderTypeLender   = "Lender"
	LoanDocumentUploaderTypeExtend   = "Extend"
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
			Label: string(LoanDocumentTypeValidID),
			Value: string(LoanDocumentTypeValidID),
		},
		{
			Label: string(LoanDocumentTypePurchaseAgreement),
			Value: string(LoanDocumentTypePurchaseAgreement),
		},
		{
			Label: string(LoanDocumentTypeEvidenceOfInsurance),
			Value: string(LoanDocumentTypeEvidenceOfInsurance),
		},
		{
			Label: string(LoanDocumentTypePreliminaryTitleReport),
			Value: string(LoanDocumentTypePreliminaryTitleReport),
		},
		{
			Label: string(LoanDocumentTypeREExperience),
			Value: string(LoanDocumentTypeREExperience),
		},
		{
			Label: string(LoanDocumentTypeBankStatements),
			Value: string(LoanDocumentTypeBankStatements),
		},
		{
			Label: string(LoanDocumentTypeBorrowerCreditScore),
			Value: string(LoanDocumentTypeBorrowerCreditScore),
		},
		{
			Label: string(LoanDocumentTypeBorrowerDefaultFraudRiskViewCriminalRecords),
			Value: string(LoanDocumentTypeBorrowerDefaultFraudRiskViewCriminalRecords),
		},
		{
			Label: string(LoanDocumentTypePropertyData),
			Value: string(LoanDocumentTypePropertyData),
		},
		{
			Label: string(LoanDocumentTypeTitleReport),
			Value: string(LoanDocumentTypeTitleReport),
		},
	})
	// register loan document category
	enums.Register(LoanDocumentCategoryName, []enums.EnumOpt{
		{
			Label: string(LoanDocumentUploaderTypeBorrower),
			Value: string(LoanDocumentUploaderTypeBorrower),
		},
		{
			Label: string(LoanDocumentUploaderTypeLender),
			Value: string(LoanDocumentUploaderTypeLender),
		},
		{
			Label: string(LoanDocumentUploaderTypeExtend),
			Value: string(LoanDocumentUploaderTypeExtend),
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

func GetLoanDocumentTypes() map[LoanDocumentCategory][]LoanDocumentType {
	return map[LoanDocumentCategory][]LoanDocumentType{
		LoanDocumentUploaderTypeBorrower: []LoanDocumentType{
			LoanDocumentTypeValidID,
			LoanDocumentTypePurchaseAgreement,
			LoanDocumentTypeEvidenceOfInsurance,
			LoanDocumentTypePreliminaryTitleReport,
			LoanDocumentTypeREExperience,
			LoanDocumentTypeBankStatements,
		},
		LoanDocumentUploaderTypeLender: []LoanDocumentType{
			LoanDocumentTypeBorrowerCreditScore,
			LoanDocumentTypeBorrowerDefaultFraudRiskViewCriminalRecords,
			LoanDocumentTypePropertyData,
			LoanDocumentTypeTitleReport,
		},
	}

}
