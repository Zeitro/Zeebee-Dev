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

// GetBorrowerUser
func (repo *UserRepo) GetBorrowerUser(ctx context.Context, id int64) (res *model.User, err error) {
	err = repo.BaseRepo.Client(ctx).
		Where("id = ? and type = ?", id, constants.UserTypeBorrower).
		Preload("Borrower").
		First(&res).Error
	return res, errorx.WithStack(err)
}

func (repo *UserRepo) GetBorrowerUserByEmail(ctx context.Context, email string) (res *model.User, err error) {
	err = repo.BaseRepo.Client(ctx).
		Where("email = ? and type = ?", email, constants.UserTypeBorrower).
		Preload("Borrower").
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

// GetLenderUser
func (repo *UserRepo) GetLenderUser(ctx context.Context, id int64) (res *model.User, err error) {
	err = repo.BaseRepo.Client(ctx).
		Where("id = ? and type = ?", id, constants.UserTypeLender).
		Preload("Organizations").
		// Preload("LenderEmployee.LenderCompany").
		First(&res).Error
	return res, errorx.WithStack(err)
}

func (repo *UserRepo) GetLenderUserByEmail(ctx context.Context, email string) (res *model.User, err error) {
	err = repo.BaseRepo.Client(ctx).
		Where("email = ? and type = ?", email, constants.UserTypeLender).
		Preload("Organizations").
		// Preload("LenderEmployee.LenderCompany").
		First(&res).Error
	return res, errorx.WithStack(err)
}
