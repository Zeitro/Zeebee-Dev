package repo

import (
	"context"

	"github.com/pkg/errors"
	"gorm.io/gorm"

	"zeebee-backend/app/borrower/rpc/internal/repo/model"
	"zeebee-backend/common/constants"
	"zeebee-backend/common/gormx"
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
	return res, errors.WithStack(err)
}

func (repo *UserRepo) GetBorrowerUserByEmail(ctx context.Context, email string) (res *model.User, err error) {
	err = repo.BaseRepo.Client(ctx).
		Where("email = ? and type = ?", email, constants.UserTypeBorrower).
		Preload("Borrower").
		First(&res).Error
	return res, errors.WithStack(err)
}

// update user status
func (repo *UserRepo) UpdateUserStatus(ctx context.Context, id int64, status constants.UserStatus) (err error) {
	err = repo.BaseRepo.Client(ctx).Where("id = ?", id).Update("status", status).Error
	return errors.WithStack(err)
}

// update user password
func (repo *UserRepo) UpdateUserPassword(ctx context.Context, id int64, password string) (err error) {
	err = repo.BaseRepo.Client(ctx).Where("id = ?", id).Update("password_hash", password).Error
	return errors.WithStack(err)
}
