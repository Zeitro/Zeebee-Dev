package constants

import (
	"zeebee-backend/common/enums"
)

type LoanDocumentFileType string

const LoanDocumentFileTypeName = "loan_document_file_type"

const (
	// borrower document
	LoanDocumentFileTypeValidID                = "Valid ID of Borrower"
	LoanDocumentFileTypePurchaseAgreement      = "Purchase Agreement"
	LoanDocumentFileTypeEvidenceOfInsurance    = "Evidence of Insurance"
	LoanDocumentFileTypePreliminaryTitleReport = "Preliminary Title Report/Escrow Contact"
	LoanDocumentFileTypeREExperience           = "RE Experience"
	LoanDocumentFileTypeBankStatements         = "Bank Statements for Most Recent Two Months"

	// lender document
	LoanDocumentFileTypeBorrowerCreditScore                         = "Borrower's Credit Score"
	LoanDocumentFileTypeBorrowerDefaultFraudRiskViewCriminalRecords = "Borrower's Default/Fraud/Risk/View/Criminal Records And RE Experience"
	LoanDocumentFileTypePropertyData                                = "Property Data/Historical Property Ownership"
	LoanDocumentFileTypeTitleReport                                 = "Title Report"

	// lender closing
	LoanDocumentFileTypeClosingDocs = "Closing Docs"
	LoanDocumentFileTypeTitleDocs   = "Title Docs"
	LoanDocumentFileTypeFormalLOI   = "Formal LOI"

	// init
	LoanDocumentFileTypeApplyForm = "Apply Form"
)

// uploader type
type LoanDocumentCategory string

const LoanDocumentCategoryName = "loan_document_category"
const (
	LoanDocumentCategoryBorrower      = "Borrower"
	LoanDocumentCategoryLenderClosing = "LenderClosing"
	LoanDocumentCategoryLender        = "Lender"
	LoanDocumentCategoryExtend        = "Extend"
	LoanDocumentCategoryInit          = "Init"
)

// loan  document status: Pending Uploand, Requested, Ready For review
type LoanDocumentStatus string

const LoanDocumentStatusName = "loan_document_status"
const (
	LoanDocumentStatusPendingUpload    LoanDocumentStatus = "PendingUpload"
	LoanDocumentStatusPendingRequested LoanDocumentStatus = "PendingRequested"
	LoanDocumentStatusRequested        LoanDocumentStatus = "Requested"
	LoanDocumentStatusReadyForReview   LoanDocumentStatus = "ReadyForReview"
)

func init() {
	// register loan document type
	enums.Register(LoanDocumentFileTypeName, []enums.EnumOpt{
		{
			Label: string(LoanDocumentFileTypeValidID),
			Value: string(LoanDocumentFileTypeValidID),
		},
		{
			Label: string(LoanDocumentFileTypePurchaseAgreement),
			Value: string(LoanDocumentFileTypePurchaseAgreement),
		},
		{
			Label: string(LoanDocumentFileTypeEvidenceOfInsurance),
			Value: string(LoanDocumentFileTypeEvidenceOfInsurance),
		},
		{
			Label: string(LoanDocumentFileTypePreliminaryTitleReport),
			Value: string(LoanDocumentFileTypePreliminaryTitleReport),
		},
		{
			Label: string(LoanDocumentFileTypeREExperience),
			Value: string(LoanDocumentFileTypeREExperience),
		},
		{
			Label: string(LoanDocumentFileTypeBankStatements),
			Value: string(LoanDocumentFileTypeBankStatements),
		},
		{
			Label: string(LoanDocumentFileTypeBorrowerCreditScore),
			Value: string(LoanDocumentFileTypeBorrowerCreditScore),
		},
		{
			Label: string(LoanDocumentFileTypeBorrowerDefaultFraudRiskViewCriminalRecords),
			Value: string(LoanDocumentFileTypeBorrowerDefaultFraudRiskViewCriminalRecords),
		},
		{
			Label: string(LoanDocumentFileTypePropertyData),
			Value: string(LoanDocumentFileTypePropertyData),
		},
		{
			Label: string(LoanDocumentFileTypeTitleReport),
			Value: string(LoanDocumentFileTypeTitleReport),
		},
	})
	// register loan document category
	enums.Register(LoanDocumentCategoryName, []enums.EnumOpt{
		{
			Label: string(LoanDocumentCategoryBorrower),
			Value: string(LoanDocumentCategoryBorrower),
		},
		{
			Label: string(LoanDocumentCategoryLenderClosing),
			Value: string(LoanDocumentCategoryLenderClosing),
		},
		{
			Label: string(LoanDocumentCategoryLender),
			Value: string(LoanDocumentCategoryLender),
		},
		{
			Label: string(LoanDocumentCategoryExtend),
			Value: string(LoanDocumentCategoryExtend),
		},
		{
			Label: string(LoanDocumentCategoryInit),
			Value: string(LoanDocumentCategoryInit),
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

func GetLoanDocumentTypes() map[LoanDocumentCategory][]LoanDocumentFileType {
	return map[LoanDocumentCategory][]LoanDocumentFileType{
		LoanDocumentCategoryBorrower: []LoanDocumentFileType{
			LoanDocumentFileTypeValidID,
			LoanDocumentFileTypePurchaseAgreement,
			LoanDocumentFileTypeEvidenceOfInsurance,
			LoanDocumentFileTypePreliminaryTitleReport,
			LoanDocumentFileTypeREExperience,
			LoanDocumentFileTypeBankStatements,
		},
		LoanDocumentCategoryLender: []LoanDocumentFileType{
			LoanDocumentFileTypeBorrowerCreditScore,
			LoanDocumentFileTypeBorrowerDefaultFraudRiskViewCriminalRecords,
			LoanDocumentFileTypePropertyData,
			LoanDocumentFileTypeTitleReport,
		},
		LoanDocumentCategoryLenderClosing: []LoanDocumentFileType{
			LoanDocumentFileTypeClosingDocs,
			LoanDocumentFileTypeTitleDocs,
			LoanDocumentFileTypeFormalLOI,
		},
	}

}
