package model

import "zeebee-backend/common/constants"

type LoanDocument struct {
	Base
	LoanId       int64                              `gorm:"column:loan_id;not null;default:0" json:"loan_id"`
	UploaderType constants.LoanDocumentUploaderType `gorm:"column:uploader_type;size:20;not null" json:"uploader_type"`
	FileType     constants.LoanDocumentType         `gorm:"column:file_type;size:255;not null;default:'file'" json:"file_type"`
	Status       constants.LoanDocumentStatus       `gorm:"column:status;size:255;not null;default:''" json:"status"`
	Comment      string                             `gorm:"column:comment;size:255;not null;default:''" json:"comment"`

	// relationships
	Files []*File `gorm:"polymorphicType:EntityType;polymorphicId:EntityId;polymorphicValue:LoanDocument" json:"files"`
}

func (m *LoanDocument) TableName() string {
	return "loan_documents"
}
