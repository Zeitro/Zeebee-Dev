package repo

import (
	"context"

	"gorm.io/gorm"

	"zeebee-backend/common/errorx"
	"zeebee-backend/common/gormx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/model"
)

type EmailAuthorizationRepo struct {
	*gormx.BaseRepo[*model.EmailAuthorization]
	m *model.EmailAuthorization
}

func NewEmailAuthorizationRepo(db *gorm.DB) *EmailAuthorizationRepo {
	return &EmailAuthorizationRepo{
		BaseRepo: gormx.NewBaseRepo[*model.EmailAuthorization](db),
		m:        new(model.EmailAuthorization),
	}
}

// DeleteEmailAuthorizationByOrganizationEmailIdWithTx deletes email authorization records
// from the database within a transaction for a given organization email ID and purpose.
// Note: If organizationEmailId is 0, the function returns without performing any operation.
func (repo *EmailAuthorizationRepo) DeleteEmailAuthorizationByOrganizationEmailIdWithTx(ctx context.Context, tx *gorm.DB, organizationEmailId int64, purpose constants.EmailPurpose) (err error) {
	if organizationEmailId == 0 {
		return err
	}
	err = tx.WithContext(ctx).
		Where("organization_email_id = ? and purpose = ?", organizationEmailId, purpose).
		Delete(&model.EmailAuthorization{}).Error
	return errorx.WithStack(err)
}

// GetEmailAuthorizationByOrganizationEmailIdWithTx retrieves an email authorization record by organization ID and email scope within a transaction.
// Note: Returns an error if no record is found or if any database error occurs
func (repo *EmailAuthorizationRepo) GetEmailAuthorizationByOrganizationEmailIdWithTx(ctx context.Context, tx *gorm.DB, organizationEmailId int64, purpose constants.EmailPurpose) (result *model.EmailAuthorization, err error) {
	err = tx.WithContext(ctx).
		Where("organization_email_id = ? and purpose = ?", organizationEmailId, purpose).
		First(result).Error
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	if result == nil || result.Id == 0 {
		return nil, errorx.New("Email Authorization Not Found")
	}
	return result, nil
}

// GetEmailAuthorizationByOrganizationEmailId retrieves an email authorization record from the database
// based on the provided organization email ID and email purpose.
func (repo *EmailAuthorizationRepo) GetEmailAuthorizationByOrganizationEmailId(ctx context.Context, organizationEmailId int64, userId int64, purpose constants.EmailPurpose) (result *model.EmailAuthorization, err error) {
	err = repo.Client(ctx).
		Where("organization_email_id = ? and user_id = ?", organizationEmailId, userId).
		Where("purpose = ?", purpose).
		First(&result).Error
	return result, errorx.WithStack(err)
}

// GetOrCreateEmailAuthorizationWithTx retrieves an existing email authorization or creates a new one within a transaction.
// It first checks if an authorization exists for the given organization and email scope.
// If not found, it creates a new authorization record.
func (repo *EmailAuthorizationRepo) GetOrCreateEmailAuthorizationWithTx(ctx context.Context, tx *gorm.DB, data *model.EmailAuthorization) (result *model.EmailAuthorization, isNew bool, err error) {
	if data == nil || (data.OrganizationEmailId == 0 && data.UserId == 0) {
		return nil, false, errorx.New("Invalid Data")
	}
	query := tx.WithContext(ctx).
		Where("organization_email_id = ?", data.OrganizationEmailId).
		Where("user_id = ?", data.UserId).
		Where("purpose = ?", data.Purpose)
	result = &model.EmailAuthorization{}
	err = query.Find(result).Error
	if err != nil {
		return nil, false, errorx.WithStack(err)
	}
	if result == nil || result.Id == 0 {
		data.IsAuthorized = false
		err = tx.WithContext(ctx).Create(data).Error
		if err != nil {
			return nil, false, errorx.WithStack(err)
		}
		result = data
		isNew = true
	}
	return result, isNew, nil
}
