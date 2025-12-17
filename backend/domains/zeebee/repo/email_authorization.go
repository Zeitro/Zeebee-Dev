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

// GetOrCreateEmailAuthorizationWithTx retrieves an existing email authorization or creates a new one within a transaction.
// It first checks if an authorization exists for the given organization and email scope.
// If not found, it creates a new authorization record.
// Notes:
//   - Returns an error if input data is nil
//   - Uses errorx package for error wrapping and stack traces
//   - The operation is performed within the provided transaction
func (repo *EmailAuthorizationRepo) GetOrCreateEmailAuthorizationWithTx(ctx context.Context, tx *gorm.DB, data *model.EmailAuthorization) (result *model.EmailAuthorization, isNew bool, err error) {
	if data == nil || (data.OrganizationId == 0 && data.UserId == 0) {
		return nil, false, errorx.New("Invalid Data")
	}
	query := tx.WithContext(ctx).Where("email_scope = ?", data.EmailScope)
	if data.OrganizationId > 0 {
		query = query.Where("organization_id = ?", data.OrganizationId)
	}
	if data.UserId > 0 {
		query = query.Where("user_id = ?", data.UserId)
	}
	err = query.Find(result).
		Error
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
