package repo

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"zeebee-backend/common/errorx"
	"zeebee-backend/common/gormx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/model"
)

type DriveAuthorizationRepo struct {
	*gormx.BaseRepo[*model.DriveAuthorization]
}

func NewDriveAuthorizationRepo(db *gorm.DB) *DriveAuthorizationRepo {
	return &DriveAuthorizationRepo{
		BaseRepo: gormx.NewBaseRepo[*model.DriveAuthorization](db),
	}
}

// GetByOrganization returns the authorization row for an organization, or (nil, nil)
// when none exists (a missing row simply means "not connected").
func (repo *DriveAuthorizationRepo) GetByOrganization(ctx context.Context, organizationId int64, organizationType constants.UserType) (result *model.DriveAuthorization, err error) {
	err = repo.Client(ctx).
		Where("organization_id = ? and organization_type = ?", organizationId, organizationType).
		First(&result).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, errorx.WithStack(err)
	}
	return result, nil
}

// Upsert creates the row when absent, otherwise partially updates it. Preference
// fields (folder_name/root_folder_id/root_folder_url) are only overwritten when the
// incoming value is non-empty, so a token refresh never clobbers admin settings.
func (repo *DriveAuthorizationRepo) Upsert(ctx context.Context, data *model.DriveAuthorization) (err error) {
	existing, err := repo.GetByOrganization(ctx, data.OrganizationId, data.OrganizationType)
	if err != nil {
		return err
	}
	if existing == nil {
		if data.FolderName == "" {
			data.FolderName = "Zeebee"
		}
		if data.DriveName == "" {
			data.DriveName = "My Drive"
		}
		return errorx.WithStack(repo.Client(ctx).Create(data).Error)
	}

	updates := map[string]interface{}{
		"auth_email":            data.AuthEmail,
		"is_authorized":         data.IsAuthorized,
		"token":                 data.Token,
		"authorized_by_user_id": data.AuthorizedByUserId,
	}
	if data.FolderName != "" {
		updates["folder_name"] = data.FolderName
	}
	if data.RootFolderId != "" {
		updates["root_folder_id"] = data.RootFolderId
	}
	if data.RootFolderUrl != "" {
		updates["root_folder_url"] = data.RootFolderUrl
	}
	return repo.updateByOrganization(ctx, data.OrganizationId, data.OrganizationType, updates)
}

// UpdateToken persists a refreshed token only.
func (repo *DriveAuthorizationRepo) UpdateToken(ctx context.Context, organizationId int64, organizationType constants.UserType, token *model.DriveToken) (err error) {
	return repo.updateByOrganization(ctx, organizationId, organizationType, map[string]interface{}{
		"token": token,
	})
}

// UpdateRootFolder caches the ensured root folder id and web link.
func (repo *DriveAuthorizationRepo) UpdateRootFolder(ctx context.Context, organizationId int64, organizationType constants.UserType, rootFolderId, rootFolderUrl string) (err error) {
	return repo.updateByOrganization(ctx, organizationId, organizationType, map[string]interface{}{
		"root_folder_id":  rootFolderId,
		"root_folder_url": rootFolderUrl,
	})
}

// UpdateDriveSelection sets the target drive and clears the cached root folder so
// the next sync rebuilds folders on the newly selected drive.
func (repo *DriveAuthorizationRepo) UpdateDriveSelection(ctx context.Context, organizationId int64, organizationType constants.UserType, driveId, driveName string) (err error) {
	return repo.updateByOrganization(ctx, organizationId, organizationType, map[string]interface{}{
		"drive_id":        driveId,
		"drive_name":      driveName,
		"root_folder_id":  "",
		"root_folder_url": "",
	})
}

// UpdateFolderName sets the root folder name and clears the cached root folder.
func (repo *DriveAuthorizationRepo) UpdateFolderName(ctx context.Context, organizationId int64, organizationType constants.UserType, folderName string) (err error) {
	return repo.updateByOrganization(ctx, organizationId, organizationType, map[string]interface{}{
		"folder_name":     folderName,
		"root_folder_id":  "",
		"root_folder_url": "",
	})
}

// Disconnect clears token/email/root folder but keeps folder_name/drive selection.
func (repo *DriveAuthorizationRepo) Disconnect(ctx context.Context, organizationId int64, organizationType constants.UserType) (err error) {
	return repo.updateByOrganization(ctx, organizationId, organizationType, map[string]interface{}{
		"is_authorized":         false,
		"auth_email":            "",
		"token":                 &model.DriveToken{},
		"authorized_by_user_id": 0,
		"root_folder_id":        "",
		"root_folder_url":       "",
	})
}

func (repo *DriveAuthorizationRepo) updateByOrganization(ctx context.Context, organizationId int64, organizationType constants.UserType, updates map[string]interface{}) (err error) {
	err = repo.Client(ctx).
		Model(new(model.DriveAuthorization)).
		Where("organization_id = ? and organization_type = ?", organizationId, organizationType).
		Updates(updates).Error
	return errorx.WithStack(err)
}
