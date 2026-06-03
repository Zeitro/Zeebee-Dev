package dto

import (
	"golang.org/x/oauth2"

	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/drive"
)

// DriveIntegrationDto is the per-organization Google Drive connection state.
type DriveIntegrationDto struct {
	IsAuthorized       bool   `json:"is_authorized"`
	AuthEmail          string `json:"auth_email"`
	FolderName         string `json:"folder_name"`
	DriveId            string `json:"drive_id"`
	DriveName          string `json:"drive_name"`
	RootFolderUrl      string `json:"root_folder_url"`
	AuthorizedByUserId int64  `json:"authorized_by_user_id"`
}

// DriveAuthUrlRequestDto carries the inputs needed to build a consent URL.
type DriveAuthUrlRequestDto struct {
	OrganizationId    int64
	UserId            int64
	OrganizationType  constants.UserType
	AuthCodeGenerator *drive.DriveAuthCodeUrlGenerator
}

// DriveAuthorizeRequestDto carries the OAuth callback code and oauth config.
type DriveAuthorizeRequestDto struct {
	Code             string
	State            string
	OrganizationId   int64
	UserId           int64
	OrganizationType constants.UserType
	Oauth2Config     *oauth2.Config
}

// DriveListDrivesRequestDto carries inputs for listing shared drives.
type DriveListDrivesRequestDto struct {
	OrganizationId   int64
	OrganizationType constants.UserType
	Oauth2Config     *oauth2.Config
}

// DriveUpdateSelectionRequestDto switches the target drive (shared drive or My Drive).
type DriveUpdateSelectionRequestDto struct {
	OrganizationId   int64
	OrganizationType constants.UserType
	DriveId          string
	DriveName        string
}

// DriveUpdateFolderRequestDto renames the root folder.
type DriveUpdateFolderRequestDto struct {
	OrganizationId   int64
	OrganizationType constants.UserType
	FolderName       string
}

// DriveSyncTransactionRequestDto carries inputs for syncing one transaction's documents.
type DriveSyncTransactionRequestDto struct {
	TransactionId    int64
	OrganizationId   int64
	OrganizationType constants.UserType
	Oauth2Config     *oauth2.Config
}

// DriveSyncedFileDto is the per-file outcome of a sync.
type DriveSyncedFileDto struct {
	FileName string `json:"file_name"`
	DriveUrl string `json:"drive_url"`
	Status   string `json:"status"` // created | replaced | failed | removed
}

// DriveSyncTransactionResultDto aggregates a sync run.
type DriveSyncTransactionResultDto struct {
	Total         int64                `json:"total"`
	Created       int64                `json:"created"`
	Replaced      int64                `json:"replaced"`
	Failed        int64                `json:"failed"`
	Removed       int64                `json:"removed"`
	FolderWebLink string               `json:"folder_web_link"`
	Files         []DriveSyncedFileDto `json:"files"`
}
