package model

import (
	"database/sql"

	"zeebee-backend/domains/zeebee/constants"
)

type LoanDocument struct {
	Base
	LoanId         int64                          `gorm:"column:loan_id;not null;default:0" json:"loan_id"`
	Category       constants.LoanDocumentCategory `gorm:"column:category;size:100;not null" json:"category"`
	FileType       constants.LoanDocumentFileType `gorm:"column:file_type;size:30;not null;default:'file'" json:"file_type"`
	Status         constants.LoanDocumentStatus   `gorm:"column:status;size:255;not null;default:''" json:"status"`
	Comment        string                         `gorm:"column:comment;size:255;not null;default:''" json:"comment"`
	LastUploadedAt *sql.NullTime                  `gorm:"column:last_uploaded_at" json:"last_uploaded_at"`

	// relationships
	Files []*File `gorm:"polymorphicType:EntityType;polymorphicId:EntityId;polymorphicValue:LoanDocument" json:"files"`
}

func (m *LoanDocument) TableName() string {
	return "loan_documents"
}
