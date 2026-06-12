package dto

import (
	"time"

	"zeebee-backend/domains/zeebee/constants"
)

type FileEntityDto struct {
	Id          int64                    `json:"id"`
	EntityType  constants.FileEntityType `json:"entity_type"`
	EntityId    int64                    `json:"entity_id"`
	FileUrl     string                   `json:"file_url"`
	FileName    string                   `json:"file_name"`
	StorageType constants.StorageType    `json:"storage_type"`
	ExternalId  string                   `json:"external_id"`
	UploadType  constants.UploadType     `json:"upload_type"`
	CreatedAt   time.Time                `json:"created_at"`
	UpdatedAt   time.Time                `json:"updated_at"`
	// Google Drive 同步标记：drive_synced_at 非空即表示当前在 Drive 对应文件夹中。
	DriveUrl      string `json:"drive_url"`
	DriveSyncedAt string `json:"drive_synced_at"`
}
