package repo

import (
	"context"

	"gorm.io/gorm"

	"zeebee-backend/common/constants"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/gormx"
	"zeebee-backend/domains/zeebee/model"
)

type OrganizationUserRepo struct {
	*gormx.BaseRepo[*model.OrganizationUser]
	m *model.OrganizationUser
}

func NewOrganizationUserRepo(db *gorm.DB) *OrganizationUserRepo {
	return &OrganizationUserRepo{
		BaseRepo: gormx.NewBaseRepo[*model.OrganizationUser](db),
		m:        new(model.OrganizationUser),
	}
}

// GetOrganizationUser get organization user by organizationId and userId
func (repo *OrganizationUserRepo) GetOrganizationUser(ctx context.Context, organizationId int64, userId int64) (res *model.OrganizationUser, err error) {
	err = repo.Client(ctx).
		Where("organization_id = ? and user_id = ?", organizationId, userId).
		First(&res).Error
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	return res, nil
}

// DeleteOrganizationUserWithTx delete relation with transaction ctx
func (repo *OrganizationUserRepo) DeleteOrganizationUseWithTx(ctx context.Context, tx *gorm.DB, organizationId int64, userId int64) (err error) {
	err = tx.WithContext(ctx).
		Where("organization_id = ? and user_id = ?", organizationId, userId).
		Delete(&model.OrganizationUser{}).Error
	return errorx.WithStack(err)
}

// UpdateOrganizationUserStatus update relation status
func (repo *OrganizationUserRepo) UpdateOrganizationUserStatus(ctx context.Context, organizationId int64, userId int64, status constants.OrganizationUserStatus) (err error) {
	err = repo.Client(ctx).Model(&model.OrganizationUser{}).
		Where("organization_id = ? and user_id = ?", organizationId, userId).
		Update("status", status).Error
	return errorx.WithStack(err)
}
