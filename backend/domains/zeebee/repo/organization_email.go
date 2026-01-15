package repo

import (
	"context"

	"gorm.io/gorm"

	"zeebee-backend/common/errorx"
	"zeebee-backend/common/gormx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/model"
)

type OrganizationEmailRepo struct {
	*gormx.BaseRepo[*model.OrganizationEmail]
	m *model.OrganizationEmail
}

func NewOrganizationEmailRepo(db *gorm.DB) *OrganizationEmailRepo {
	return &OrganizationEmailRepo{
		BaseRepo: gormx.NewBaseRepo[*model.OrganizationEmail](db),
		m:        new(model.OrganizationEmail),
	}
}

func (repo *OrganizationEmailRepo) GetOrganizationEmailByOrganizationIdWithTx(ctx context.Context, tx *gorm.DB, organizationId int64, purpose constants.EmailPurpose) (result *model.OrganizationEmail, err error) {
	err = tx.WithContext(ctx).
		Where("organization_id = ? and purpose = ?", organizationId, purpose).
		First(&result).Error
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	if result == nil || result.Id == 0 {
		return nil, errorx.New("Organization Email Not Found")
	}
	return result, nil
}

func (repo *OrganizationEmailRepo) GetOrganizationEmailByOrganizationId(ctx context.Context, organizationId int64, purpose constants.EmailPurpose) (result *model.OrganizationEmail, err error) {
	err = repo.Client(ctx).
		Where("organization_id = ? and purpose = ?", organizationId, purpose).
		First(&result).Error
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	if result == nil || result.Id == 0 {
		return nil, errorx.New("Organization Email Not Found")
	}
	return result, nil
}

// GetOrCreateOrganizationEmailWithTx retrieves an existing email authorization or creates a new one within a transaction.
// It first checks if an authorization exists for the given organization and email scope.
// If not found, it creates a new authorization record.
func (repo *OrganizationEmailRepo) GetOrCreateOrganizationEmailWithTx(ctx context.Context, tx *gorm.DB, data *model.OrganizationEmail) (result *model.OrganizationEmail, isNew bool, err error) {
	if data == nil || (data.OrganizationId == 0 && len(data.Email) == 0 && len(data.Purpose) == 0) {
		return nil, false, errorx.New("Invalid Data")
	}
	result = &model.OrganizationEmail{}
	query := tx.WithContext(ctx).
		Where("organization_id = ?", data.OrganizationId).
		Where("purpose = ?", data.Purpose)
	err = query.Find(result).Error
	if err != nil {
		return nil, false, errorx.WithStack(err)
	}
	if result == nil || result.Id == 0 {
		err = tx.WithContext(ctx).Create(data).Error
		if err != nil {
			return nil, false, errorx.WithStack(err)
		}
		result = data
		isNew = true
	}
	return result, isNew, nil
}

// DeleteByIdWithTx deletes an email authorization record by its ID within a transaction.
func (repo *OrganizationEmailRepo) DeleteByIdWithTx(ctx context.Context, tx *gorm.DB, id int64) (err error) {
	var data *model.OrganizationEmail
	err = tx.WithContext(ctx).
		Model(&model.OrganizationEmail{}).
		Preload("EmailAuthorizations").
		Where("id = ?", id).
		First(&data).Error
	if err != nil {
		return errorx.WithStack(err)
	}
	err = tx.WithContext(ctx).Select("EmailAuthorizations").Delete(data).Error
	return err
}
