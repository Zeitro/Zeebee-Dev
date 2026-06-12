package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"golang.org/x/oauth2"

	"zeebee-backend/domains/zeebee/constants"
)

// DriveToken wraps an OAuth2 token for JSONB storage.
//
// It implements driver.Valuer / sql.Scanner so it serializes to JSON on every
// write path uniformly — including map-based GORM Updates (e.g. Disconnect,
// UpdateToken), where GORM does NOT apply the `serializer:json` field tag and
// would otherwise hand the raw struct to the pgx driver ("cannot find encode
// plan"). Hence the field tag intentionally omits `serializer:json`.
type DriveToken struct {
	*oauth2.Token
}

// Value serializes the wrapped token to a JSON string for the jsonb column.
// A nil inner token is stored as an empty JSON object rather than SQL NULL,
// matching the column default ('{}').
func (t DriveToken) Value() (driver.Value, error) {
	if t.Token == nil {
		return "{}", nil
	}
	b, err := json.Marshal(t.Token)
	if err != nil {
		return nil, err
	}
	return string(b), nil
}

// Scan parses the jsonb value back into the wrapped token. Empty/null/{} are
// treated as "no token" so a disconnected row scans to a nil inner token.
func (t *DriveToken) Scan(src interface{}) error {
	if src == nil {
		t.Token = nil
		return nil
	}
	var b []byte
	switch v := src.(type) {
	case []byte:
		b = v
	case string:
		b = []byte(v)
	default:
		return fmt.Errorf("unsupported source type for DriveToken: %T", src)
	}
	s := string(b)
	if len(b) == 0 || s == "{}" || s == "null" {
		t.Token = nil
		return nil
	}
	var tok oauth2.Token
	if err := json.Unmarshal(b, &tok); err != nil {
		return err
	}
	t.Token = &tok
	return nil
}

// DriveAuthorization stores a single organization's Google Drive OAuth state.
// There is one row per (organization_id, organization_type); uniqueness is
// enforced at the application layer via Upsert (mirrors email_authorizations).
type DriveAuthorization struct {
	Base
	OrganizationId     int64              `gorm:"column:organization_id;not null;default:0" json:"organization_id"`
	OrganizationType   constants.UserType `gorm:"column:organization_type;not null;default:''" json:"organization_type"`
	AuthEmail          string             `gorm:"column:auth_email;not null;default:''" json:"auth_email"`
	IsAuthorized       bool               `gorm:"column:is_authorized;not null;default:false" json:"is_authorized"`
	DriveId            string             `gorm:"column:drive_id;not null;default:''" json:"drive_id"`
	DriveName          string             `gorm:"column:drive_name;not null;default:'My Drive'" json:"drive_name"`
	FolderName         string             `gorm:"column:folder_name;not null;default:'Zeebee'" json:"folder_name"`
	RootFolderId       string             `gorm:"column:root_folder_id;not null;default:''" json:"root_folder_id"`
	RootFolderUrl      string             `gorm:"column:root_folder_url;not null;default:''" json:"root_folder_url"`
	Token              *DriveToken        `gorm:"column:token;type:jsonb;default:'{}'" json:"token"`
	AuthorizedByUserId int64              `gorm:"column:authorized_by_user_id;not null;default:0" json:"authorized_by_user_id"`
}

func (m DriveAuthorization) TableName() string {
	return "drive_authorizations"
}
