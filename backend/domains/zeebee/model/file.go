package model

import (
	"database/sql"

	"zeebee-backend/domains/zeebee/constants"
)

type File struct {
	Base

	EntityType constants.FileEntityType `gorm:"column:entity_type;size:50;not null;default:'LoanDocument'" json:"entity_type"`
	EntityId   int64                    `gorm:"column:entity_id;not null;default:0" json:"entity_id"`
	FileUrl    string                   `gorm:"column:file_url;size:255;not null;default:''" json:"file_url"`
	FileName   string                   `gorm:"column:file_name;type:varchar(255);not null;default:''" json:"file_name"`
	StorageType constants.StorageType    `gorm:"column:storage_type;type:varchar(20);not null;default:'Aws S3'" json:"storage_type"`
	ExternalId string                   `gorm:"column:external_id;type:varchar(255);not null;default:''" json:"external_id"`
	UploadType constants.UploadType     `gorm:"column:upload_type;type:varchar(20);not null;default:'Manual'" json:"upload_type"`
	// Google Drive 同步标记：drive_synced_at 非空即表示当前已同步到 Drive。
	DriveUrl      string        `gorm:"column:drive_url;type:varchar(512);not null;default:''" json:"drive_url"`
	DriveSyncedAt *sql.NullTime `gorm:"column:drive_synced_at;type:timestamptz" json:"drive_synced_at"`
}

func (m *File) TableName() string {
	return "files"
}
