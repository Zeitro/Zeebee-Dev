package dto

import (
	"zeebee-backend/domains/zeebee/constants"
)

// Checklist structures (cached in transaction.checklist JSONB)
type ChecklistRoot struct {
	AllChecked                         bool             `json:"all_checked"`
	TotalDocumentCount                 int32            `json:"total_document_count"`
	TotalRequiredDocumentCount         int32            `json:"total_required_document_count"`
	TotalRequiredUploadedDocumentCount int32            `json:"total_required_uploaded_document_count"`
	TotalOptionalDocumentCount         int32            `json:"total_optional_document_count"`
	TotalOptionalUploadedDocumentCount int32            `json:"total_optional_uploaded_document_count"`
	RequiredDocuments                  []*ChecklistNode `json:"required_documents"`
	OptionalDocuments                  []*ChecklistNode `json:"optional_documents"`
}

type ChecklistNode struct {
	Name                               string           `json:"name"`
	DocumentId                         int64            `json:"document_id"`
	IsRequired                         bool             `json:"is_required"`
	SubDocumentCount                   int32            `json:"sub_document_count"`
	UploadedDocumentCount              int32            `json:"uploaded_document_count"`
	RequiredDocumentCount              int32            `json:"required_document_count"`
	OptionalDocumentCount              int32            `json:"optional_document_count"`
	TotalRequiredUploadedDocumentCount int32            `json:"total_required_uploaded_document_count"`
	TotalOptionalUploadedDocumentCount int32            `json:"total_optional_uploaded_document_count"`
	Children                           []*ChecklistNode `json:"children,omitempty"`
	FileCount                          int32            `json:"file_count,omitempty"`
	Files                              []*FileEntityDto `json:"files,omitempty"`
}

type BpTransactionSwitchStatusDto struct {
	Id             int64                                       `json:"id"`
	Status         constants.BorrowerPropertyTransactionStatus `json:"status"`
	OrganizationId int64                                       `json:"organization_id"`
	OperatorId     int64                                       `json:"operator_id"`
}

// Page request/filter/response

type BorrowerPropertyTransactionPageRequestDto struct {
	BasePageRequestDto
	Filter BorrowerPropertyTransactionPageRequestFilter `json:"filter"`
	Sorter BasePageSorter                               `json:"sorter"`
}

type BorrowerPropertyTransactionPageRequestFilter struct {
	OrganizationId int64  `json:"organization_id"`
	Keyword        string `json:"keyword,omitempty"`
	Type           string `json:"type,omitempty"`
	Status         string `json:"status,omitempty"`
}

type BorrowerPropertyTransactionPageEntity struct {
	Id        int64      `json:"id"`
	Code      string     `json:"code"`
	Type      string     `json:"type"`
	Status    string     `json:"status"`
	CoeDate   string     `json:"coe_date"`
	CreatedAt string     `json:"created_at"`
	UpdatedAt string     `json:"updated_at"`
	Address   AddressDto `json:"address"`
}
