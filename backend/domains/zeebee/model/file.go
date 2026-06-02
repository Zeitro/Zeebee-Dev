package model

import "zeebee-backend/domains/zeebee/constants"

type File struct {
	Base

	EntityType constants.FileEntityType `gorm:"column:entity_type;size:50;not null;default:'LoanDocument'" json:"entity_type"`
	EntityId   int64                    `gorm:"column:entity_id;not null;default:0" json:"entity_id"`
	FileUrl    string                   `gorm:"column:file_url;size:255;not null;default:''" json:"file_url"`
	FileName   string                   `gorm:"column:file_name;type:varchar(255);not null;default:''" json:"file_name"`
	StorageType constants.StorageType    `gorm:"column:storage_type;type:varchar(20);not null;default:'Aws S3'" json:"storage_type"`
	ExternalId string                   `gorm:"column:external_id;type:varchar(255);not null;default:''" json:"external_id"`
	UploadType constants.UploadType     `gorm:"column:upload_type;type:varchar(20);not null;default:'Manual'" json:"upload_type"`
}

func (m *File) TableName() string {
	return "files"
}
