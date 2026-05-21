package dto

import (
	"database/sql"

	"zeebee-backend/domains/zeebee/constants"
)

type BpTransactionCreateDto struct {
	BorrowerPropertyId int64                       `json:"borrower_property_id"`
	Type               constants.BorrowerPropertyTransactionType `json:"type"`
	CoeDate            *sql.NullTime               `json:"coe_date,omitempty"`
	PurchaseAmount     float64                     `json:"purchase_amount,omitempty"`
	EmdAmount          float64                     `json:"emd_amount,omitempty"`
	PropertyType       constants.PropertyType      `json:"property_type,omitempty"`
	Pm                 string                      `json:"pm,omitempty"`
	Strategy           string                      `json:"strategy,omitempty"`
	Note               string                      `json:"note,omitempty"`
	OrganizationId     int64                       `json:"organization_id"`
	OperatorId         int64                       `json:"operator_id"`
}

type BpTransactionUpdateDto struct {
	Id             int64                  `json:"id"`
	CoeDate        *sql.NullTime          `json:"coe_date,omitempty"`
	PurchaseAmount *float64               `json:"purchase_amount,omitempty"`
	EmdAmount      *float64               `json:"emd_amount,omitempty"`
	PropertyType   constants.PropertyType `json:"property_type,omitempty"`
	Pm             string                 `json:"pm,omitempty"`
	Strategy       string                 `json:"strategy,omitempty"`
	Note           string                 `json:"note,omitempty"`
	OrganizationId int64                  `json:"organization_id"`
	OperatorId     int64                  `json:"operator_id"`
}

type BpTransactionChecklistDto struct {
	Id             int64  `json:"id"`
	AllChecked     bool   `json:"all_checked"`
	OrganizationId int64  `json:"organization_id"`
	OperatorId     int64  `json:"operator_id"`
}

// Checklist structures (cached in transaction.checklist JSONB)

type ChecklistRoot struct {
	AllChecked                          bool              `json:"all_checked"`
	TotalDocumentCount                  int32             `json:"total_document_count"`
	TotalRequiredDocumentCount          int32             `json:"total_required_document_count"`
	TotalRequiredUpdatedDocumentCount   int32             `json:"total_required_updated_document_count"`
	TotalOptionalDocumentCount          int32             `json:"total_optional_document_count"`
	TotalOptionalUpdatedDocumentCount   int32             `json:"total_optional_updated_document_count"`
	RequiredDocuments                   []ChecklistNode   `json:"required_documents"`
	OptionalDocuments                   []ChecklistNode   `json:"optional_documents"`
}

type ChecklistNode struct {
	Name                  string          `json:"name"`
	DocumentId            int64           `json:"document_id"`
	IsRequired            bool            `json:"is_required"`
	SubDocumentCount      int32           `json:"sub_document_count"`
	UploadedDocumentCount int32           `json:"uploaded_document_count"`
	Children              []ChecklistNode `json:"children,omitempty"`
	FileCount             int32           `json:"file_count,omitempty"`
}

type BpTransactionSwitchStatusDto struct {
	Id             int64                         `json:"id"`
	Status         constants.BorrowerPropertyTransactionStatus `json:"status"`
	OrganizationId int64                         `json:"organization_id"`
	OperatorId     int64                         `json:"operator_id"`
}
