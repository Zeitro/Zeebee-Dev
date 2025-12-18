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

// GetEmailAuthorizationByOrganizationIdWithTx retrieves an email authorization record by organization ID and email scope within a transaction.
// Note: Returns an error if no record is found or if any database error occurs
func (repo *EmailAuthorizationRepo) GetEmailAuthorizationByOrganizationIdWithTx(ctx context.Context, tx *gorm.DB, organizationId int64, emailScope constants.EmailScope) (result *model.EmailAuthorization, err error) {
	err = tx.WithContext(ctx).
		Where("organization_id = ? and email_scope = ?", organizationId, emailScope).
		First(result).Error
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	if result == nil || result.Id == 0 {
		return nil, errorx.New("Email Authorization Not Found")
	}
	return result, nil
}

// GetEmailAuthorization retrieves an email authorization record from the database
// based on the provided organization ID, user ID, and email scope.
func (repo *EmailAuthorizationRepo) GetEmailAuthorization(ctx context.Context, organizationId int64, userId int64, emailScope constants.EmailScope) (result *model.EmailAuthorization, err error) {
	query := repo.Client(ctx).
		Where("organization_id = ?", organizationId).
		Where("user_id = ?", userId)
	err = query.Where("email_scope = ?", emailScope).First(&result).Error
	return result, errorx.WithStack(err)
}

// GetOrCreateEmailAuthorizationWithTx retrieves an existing email authorization or creates a new one within a transaction.
// It first checks if an authorization exists for the given organization and email scope.
// If not found, it creates a new authorization record.
func (repo *EmailAuthorizationRepo) GetOrCreateEmailAuthorizationWithTx(ctx context.Context, tx *gorm.DB, data *model.EmailAuthorization) (result *model.EmailAuthorization, isNew bool, err error) {
	if data == nil || (data.OrganizationId == 0 && data.UserId == 0) {
		return nil, false, errorx.New("Invalid Data")
	}
	query := tx.WithContext(ctx).
		Where("organization_id = ?", data.OrganizationId).
		Where("user_id = ?", data.UserId).
		Where("email_scope = ?", data.EmailScope)
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
