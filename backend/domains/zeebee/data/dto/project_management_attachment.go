package dto

import "zeebee-backend/domains/zeebee/constants"

type ProjectManagementAttachmentDto struct {
	FileId      int64
	TaskId      int64
	FileName    string
	FileUrl     string
	StorageType constants.StorageType
}

type ProjectManagementAttachmentCreateDto struct {
	TaskId      int64
	FileName    string
	FileUrl     string
	StorageType constants.StorageType
}
