package repo

import (
	"context"
	"fmt"

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

// CheckOrganizationUserExistedByEmail check organization user exist
func (repo *OrganizationUserRepo) CheckOrganizationUserExistedByEmail(ctx context.Context, organizationId int64, email string, userType constants.UserType) (existed bool, err error) {
	userTable := model.User{}.TableName()
	organizationUserTable := model.OrganizationUser{}.TableName()
	var count int64
	err = repo.Client(ctx).
		Joins(fmt.Sprintf("%s.id = %s.user_id", userTable, organizationUserTable)).
		Where(fmt.Sprintf("%s.email = ? and %s.type = ?", userTable, userTable), email, userType).
		Where(fmt.Sprintf("%s.organization_id = ?", organizationUserTable), organizationId).
		Count(&count).Error
	return count > 0, errorx.WithStack(err)
}

// AllOrganizationUsers return all organization users by organizationId
func (repo *OrganizationUserRepo) AllOrganizationUsers(ctx context.Context, organizationId int64) (results []*model.OrganizationUser, err error) {
	err = repo.Client(ctx).
		Where("organization_id = ? and status = ?", organizationId, constants.OrganizationStatusActive).
		Preload("User").
		Find(&results).Error
	return results, errorx.WithStack(err)
}
