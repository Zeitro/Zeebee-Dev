package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/jinzhu/copier"
	"gorm.io/gorm"

	"zeebee-backend/common/constants"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/gormx"
	"zeebee-backend/domains/zeebee/model"
)

type UserPageRequest struct {
	Limit  int64
	Page   int64
	Filter *UserPageRequestFilter
	Sorter *UserPageRequestSorter
}
type UserPageRequestFilter struct {
	Keyword        string
	Status         string
	OrganizationId int64
	UserType       constants.UserType
}
type UserPageRequestSorter struct {
	Column string
	Sort   string
}

type UserPageEntity struct {
	Id                      int64                               `json:"id"`
	Type                    constants.UserType                  `json:"type"`
	Email                   string                              `json:"email"`
	FirstName               string                              `json:"firstName"`
	LastName                string                              `json:"lastName"`
	FullName                string                              `json:"fullName"`
	Phone                   string                              `json:"phone"`
	CreatedAt               time.Time                           `json:"created_at"`
	UpdatedAt               time.Time                           `json:"updated_at"`
	PivotOrganizationStatus constants.OrganizationAndUserStatus `json:"pivot_organization_status"`
	Roles                   []*model.Role                       `json:"roles" gorm:"-"`
}

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

// Page returns paginated users
func (repo *UserRepo) Page(ctx context.Context, pageReq *UserPageRequest) (page gormx.PageResponse[*UserPageEntity], err error) {
	userTable := model.User{}.TableName()
	organizationAndUserTable := model.OrganizationAndUser{}.TableName()
	filter := pageReq.Filter
	sorter := pageReq.Sorter

	query := repo.Client(ctx).Model(&model.User{}).
		Joins(fmt.Sprintf("left join %s on %s.user_id = %s.id", organizationAndUserTable, organizationAndUserTable, userTable)).
		Where(fmt.Sprintf("%s.type = ?", userTable), filter.UserType).
		Where(fmt.Sprintf("%s.organization_id = ?", organizationAndUserTable), filter.OrganizationId)
	if len(filter.Keyword) > 0 {
		query = query.Where(fmt.Sprintf("%s.email like ? or %s.full_name like ? or %s.phone like ?", userTable, userTable, userTable),
			"%"+filter.Keyword+"%",
			"%"+filter.Keyword+"%",
			"%"+filter.Keyword+"%",
		)
	}
	if len(filter.Status) > 0 {
		query = query.Where(fmt.Sprintf("%s.status = ?", organizationAndUserTable), filter.Status)
	}

	if len(sorter.Column) == 0 {
		sorter.Column = fmt.Sprintf("%s.id", userTable)
	}
	if len(sorter.Sort) == 0 {
		sorter.Sort = "desc"
	}
	fields := []string{
		fmt.Sprintf("%s.id", userTable),
		fmt.Sprintf("%s.type", userTable),
		fmt.Sprintf("%s.email", userTable),
		fmt.Sprintf("%s.first_name", userTable),
		fmt.Sprintf("%s.last_name", userTable),
		fmt.Sprintf("%s.full_name", userTable),
		fmt.Sprintf("%s.phone", userTable),
		fmt.Sprintf("%s.status as pivot_organization_status", organizationAndUserTable),
		fmt.Sprintf("%s.created_at", userTable),
		fmt.Sprintf("%s.updated_at", userTable),
	}

	query = query.Select(fields).Order(fmt.Sprintf("%s %s", sorter.Column, sorter.Sort))
	page, err = gormx.Page[*UserPageEntity](query, gormx.PageRequest{
		Page:        pageReq.Page,
		Limit:       pageReq.Limit,
		CountColumn: fmt.Sprintf("%s.id", userTable),
	})
	roleRepo := NewRoleRepo(repo.BaseRepo.Client(ctx))
	for i, entity := range page.List {
		roles, err := roleRepo.GetUserRoles(ctx, entity.Id, filter.OrganizationId)
		if err != nil {
			return page, err
		}
		page.List[i].Roles = roles
	}

	return page, errorx.WithStack(err)
}

// GetUser return user by id
func (repo *UserRepo) GetUser(ctx context.Context, id int64, userType constants.UserType) (res *model.User, err error) {
	err = repo.BaseRepo.Client(ctx).
		Where("id = ? and type = ?", id, userType).
		Preload("Organizations").
		First(&res).Error
	return res, errorx.WithStack(err)
}

// GetOrganizationUser get user by id and organization id
func (repo *UserRepo) GetOrganizationUser(ctx context.Context, id int64, organizationId int64, userType constants.UserType) (res *model.User, err error) {
	err = repo.BaseRepo.Client(ctx).
		Where("id = ? and type = ?", id, userType).
		Preload("Organizations", func(db *gorm.DB) *gorm.DB {
			return db.Where("id = ?", organizationId)
		}).
		Preload("Roles", func(db *gorm.DB) *gorm.DB {
			return db.Where("organization_id = ?", organizationId)
		}).
		First(&res).Error
	if err != nil {
		return res, errorx.WithStack(err)
	}

	if len(res.Organizations) == 0 {
		return res, errorx.New("User Not Belongs To This Organization")
	}
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

func (repo *UserRepo) GetUserByEmailWithTx(ctx context.Context, tx *gorm.DB, email string, userType constants.UserType) (res *model.User, err error) {
	err = tx.WithContext(ctx).
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

// CheckOrganizationUserExistWithTx check organization user exist by organization id and email
func (repo *UserRepo) CheckOrganizationUserExistWithTx(ctx context.Context, tx *gorm.DB, organizationId int64, email string) (isExisted bool, err error) {
	userTable := model.User{}.TableName()
	pivotTable := model.OrganizationAndUser{}.TableName()
	var count int64
	err = tx.WithContext(ctx).Model(&model.User{}).
		Joins(fmt.Sprintf("left join %s on %s.user_id = %s.id", pivotTable, pivotTable, userTable)).
		Where(fmt.Sprintf("%s.organization_id = ? and %s.email = ?", pivotTable, userTable), organizationId, email).
		Count(&count).Error
	return count > 0, errorx.WithStack(err)
}

// CheckUserExistWithTx check user exist by email
func (repo *UserRepo) CheckUserExistWithTx(ctx context.Context, tx *gorm.DB, email string, userType constants.UserType) (isExisted bool, err error) {
	var count int64
	err = tx.WithContext(ctx).Model(&model.User{}).
		Where("email = ? and type = ?", email, userType).
		Count(&count).Error
	return count > 0, errorx.WithStack(err)
}

// UpdateRolesWithTx update user roles with transaction
func (repo *UserRepo) UpdateRolesWithTx(ctx context.Context, tx *gorm.DB, user *model.User, newRoles []*model.Role, removedRoles []*model.Role) (err error) {
	if user == nil {
		return errorx.New("User Is Empty")
	}
	copyUser := &model.User{}
	_ = copier.CopyWithOption(copyUser, user, copier.Option{DeepCopy: true})

	copyUser.RemoveRelationships()
	query := tx.WithContext(ctx).Model(copyUser)
	if len(newRoles) > 0 {
		err = query.Association("Roles").Append(newRoles)
		if err != nil {
			return errorx.WithStack(err)
		}
	}
	if len(removedRoles) > 0 {
		err = query.Association("Roles").Delete(removedRoles)
		if err != nil {
			return errorx.WithStack(err)
		}
	}
	return nil
}
