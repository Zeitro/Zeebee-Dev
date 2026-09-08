package dto

import (
	"time"

	"zeebee-backend/domains/zeebee/constants"
)

type ProjectManagementAttachmentDto struct {
	FileId      int64
	TaskId      int64
	FileName    string
	FileUrl     string
	StorageType constants.StorageType
	UploadType  constants.UploadType
}

type ProjectManagementTaskAttachmentDto struct {
	Id          int64
	EntityId    int64
	EntityType  constants.FileEntityType
	CreatedAt   time.Time
	UpdatedAt   time.Time
	FileName    string
	FileUrl     string
	StorageType constants.StorageType
	UploadType  constants.UploadType
}

type ProjectManagementAttachmentCreateDto struct {
	TaskId      int64
	FileName    string
	FileUrl     string
	StorageType constants.StorageType
	UploadType  constants.UploadType
}
