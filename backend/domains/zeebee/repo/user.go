package repo

import (
	"context"

	"gorm.io/gorm"

	"zeebee-backend/common/constants"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/gormx"
	"zeebee-backend/domains/zeebee/model"
)

type UserRepo struct {
	*gormx.BaseRepo[*model.User]
	m *model.User
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{
		BaseRepo: gormx.NewBaseRepo[*model.User](db),
		m:        new(model.User),
	}
}

// GetUser return user by id
func (repo *UserRepo) GetUser(ctx context.Context, id int64, userType constants.UserType) (res *model.User, err error) {
	err = repo.BaseRepo.Client(ctx).
		Where("id = ? and type = ?", id, userType).
		Preload("Organizations").
		First(&res).Error
	return res, errorx.WithStack(err)
}

// RegisterWithOrganization
func (repo *UserRepo) RegisterWithOrganization(ctx context.Context, user *model.User, organization *model.Organization, relationStatus constants.OrganizationAndUserStatus) (res *model.User, err error) {
	organizationRepo := NewOrganizationRepo(repo.BaseRepo.Client(ctx))
	relationRepo := NewOrganizationAndUserRepo(repo.BaseRepo.Client(ctx))
	err = repo.Transaction(ctx, func(tx *gorm.DB) error {
		var innerErr error
		if organization.Id == 0 {
			organization, innerErr = organizationRepo.CreateWithTx(ctx, tx, organization, false)
			if innerErr != nil {
				return innerErr
			}
		}
		user.DefaultOrganizationId = organization.Id
		user, innerErr = repo.CreateWithTx(ctx, tx, user, false)
		if innerErr != nil {
			return innerErr
		}
		organization.OwnerId = user.Id
		innerErr = organizationRepo.UpdateOwnerIdWithTx(ctx, tx, organization.OwnerId, user.Id)
		if innerErr != nil {
			return innerErr
		}
		_, innerErr = relationRepo.CreateWithTx(ctx, tx, &model.OrganizationAndUser{
			OrganizationId: organization.Id,
			UserId:         user.Id,
			Status:         relationStatus,
		}, false)

		return innerErr
	})
	res = user

	return res, errorx.WithStack(err)
}

// GetUserByEmail return
func (repo *UserRepo) GetUserByEmail(ctx context.Context, email string, userType constants.UserType) (res *model.User, err error) {
	err = repo.BaseRepo.Client(ctx).
		Where("email = ? and type = ?", email, userType).
		Preload("Organizations").
		First(&res).Error
	return res, errorx.WithStack(err)
}

// update user status
func (repo *UserRepo) UpdateUserStatus(ctx context.Context, id int64, status constants.UserStatus) (err error) {
	err = repo.BaseRepo.Client(ctx).Where("id = ?", id).Update("status", status).Error
	return errorx.WithStack(err)
}

// update user password
func (repo *UserRepo) UpdateUserPassword(ctx context.Context, id int64, password string) (err error) {
	err = repo.BaseRepo.Client(ctx).Where("id = ?", id).Update("password_hash", password).Error
	return errorx.WithStack(err)
}

func (repo *UserRepo) GetLenderUserByEmail(ctx context.Context, email string) (res *model.User, err error) {
	err = repo.BaseRepo.Client(ctx).
		Where("email = ? and type = ?", email, constants.UserTypeLender).
		Preload("Organizations").
		First(&res).Error
	return res, errorx.WithStack(err)
}
