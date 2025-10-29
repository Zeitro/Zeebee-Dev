package model

import "zeebee-backend/common/constants"

type File struct {
	Base

	// FileType   constants.FileType    `gorm:"column:file_type;size:10;not null;default:'file'" json:"file_type"`
	EntityType constants.FileEntityType `gorm:"column:entity_type;size:50;not null;default:'LoanDocument'" json:"entity_type"`
	EntityId   int64                    `gorm:"column:entity_id;not null;default:0" json:"entity_id"`
	FileUrl    string                   `gorm:"column:file_url;size:255;not null;default:''" json:"file_url"`
}

func (m *File) TableName() string {
	return "files"
}
