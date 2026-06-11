package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"zeebee-backend/domains/zeebee/constants"
)

// DriveSyncFileItem is a single file's sync outcome.
type DriveSyncFileItem struct {
	FileName string `json:"file_name"`
	DriveUrl string `json:"drive_url"`
	Status   string `json:"status"`
}

// DriveSyncFiles is stored as a JSONB array. It implements driver.Valuer /
// sql.Scanner so it serializes uniformly on every write path — including
// map-based GORM Updates, where the `serializer:json` field tag is NOT applied
// and would otherwise hand the raw slice to the pgx driver.
type DriveSyncFiles []DriveSyncFileItem

func (f DriveSyncFiles) Value() (driver.Value, error) {
	if len(f) == 0 {
		return "[]", nil
	}
	b, err := json.Marshal([]DriveSyncFileItem(f))
	if err != nil {
		return nil, err
	}
	return string(b), nil
}

func (f *DriveSyncFiles) Scan(src interface{}) error {
	if src == nil {
		*f = nil
		return nil
	}
	var b []byte
	switch v := src.(type) {
	case []byte:
		b = v
	case string:
		b = []byte(v)
	default:
		return fmt.Errorf("unsupported source type for DriveSyncFiles: %T", src)
	}
	if len(b) == 0 {
		*f = nil
		return nil
	}
	return json.Unmarshal(b, (*[]DriveSyncFileItem)(f))
}

// DriveSyncJob records one asynchronous "sync transaction to Google Drive" run.
// The HTTP request only creates the row (status=Running) and returns immediately;
// a detached goroutine performs the sync and updates this row when finished.
type DriveSyncJob struct {
	Base
	OrganizationId   int64              `gorm:"column:organization_id;not null;default:0" json:"organization_id"`
	OrganizationType constants.UserType `gorm:"column:organization_type;not null;default:''" json:"organization_type"`
	TransactionId    int64              `gorm:"column:transaction_id;not null;default:0" json:"transaction_id"`
	Status           string             `gorm:"column:status;not null;default:'Running'" json:"status"`
	Total            int64              `gorm:"column:total;not null;default:0" json:"total"`
	Created          int64              `gorm:"column:created;not null;default:0" json:"created"`
	Replaced         int64              `gorm:"column:replaced;not null;default:0" json:"replaced"`
	Failed           int64              `gorm:"column:failed;not null;default:0" json:"failed"`
	Removed          int64              `gorm:"column:removed;not null;default:0" json:"removed"`
	FolderWebLink    string             `gorm:"column:folder_web_link;not null;default:''" json:"folder_web_link"`
	Files            DriveSyncFiles     `gorm:"column:files;type:jsonb;default:'[]'" json:"files"`
	Error            string             `gorm:"column:error;not null;default:''" json:"error"`
}

func (m DriveSyncJob) TableName() string {
	return "drive_sync_jobs"
}
