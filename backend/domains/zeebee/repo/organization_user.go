package repo

import (
	"context"
	"fmt"

	"github.com/jinzhu/copier"
	"gorm.io/gorm"

	"zeebee-backend/common/constants"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/gormx"
	"zeebee-backend/domains/zeebee/data/dto"
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

// Page returns paginated organization users
func (repo *OrganizationUserRepo) Page(ctx context.Context, pageReq dto.OrganizationUserPageRequestDto) (page gormx.PageResponse[*model.OrganizationUser], err error) {
	filter := pageReq.Filter
	sorter := pageReq.Sorter
	userTable := model.User{}.TableName()
	organizationUserTable := model.OrganizationUser{}.TableName()
	query := repo.Client(ctx).Model(&model.OrganizationUser{}).
		Joins(fmt.Sprintf("%s.id = %s.user_id", userTable, organizationUserTable)).
		Where(fmt.Sprintf("%s.organization_id = ?", organizationUserTable), pageReq.Filter.OrganizationId)
	if len(filter.Keyword) > 0 {
		query = query.Where(fmt.Sprintf("%s.email like ? or %s.full_name like ? or %s.phone like ?", userTable, organizationUserTable, organizationUserTable), "%"+filter.Keyword+"%", "%"+filter.Keyword+"%", "%"+filter.Keyword+"%")
	}
	if len(sorter.Column) > 0 {
		query = query.Order(fmt.Sprintf("%s.%s %s", organizationUserTable, sorter.Column, sorter.Sort))
	}
	query = query.Order(fmt.Sprintf("case %s.status when '%s' then 1 when '%s' then 2 when '%s' then 3 else 4 end", organizationUserTable,
		constants.OrganizationUserStatusActive,
		constants.OrganizationUserStatusDisabled,
		constants.OrganizationUserStatusNotActivated,
	))
	query = query.Order(fmt.Sprintf("%s.created_at asc", organizationUserTable))
	return gormx.Page[*model.OrganizationUser](query, gormx.PageRequest{
		Page:        pageReq.Page,
		Limit:       pageReq.Limit,
		CountColumn: fmt.Sprintf("%s.id", organizationUserTable),
	})
}

// GetOrganizationUser get organization user by organizationId and userId
func (repo *OrganizationUserRepo) GetOrganizationUser(ctx context.Context, organizationId int64, userId int64) (res *model.OrganizationUser, err error) {
	err = repo.Client(ctx).
		Where("organization_id = ? and user_id = ?", organizationId, userId).
		Preload("User").
		Preload("Roles").
		First(&res).Error
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	return res, nil
}

// GetOrganizationUserById get organization user by id
func (repo *OrganizationUserRepo) GetOrganizationUserById(ctx context.Context, id int64, organizationId int64) (result *model.OrganizationUser, err error) {
	err = repo.Client(ctx).
		Where("id = ? and organization_id = ?", id, organizationId).
		Preload("User").
		Preload("Roles").
		First(&result).Error
	return result, errorx.WithStack(err)
}

func (repo *OrganizationUserRepo) DeleteOrganizationUser(ctx context.Context, id int64, organizationId int64) (err error) {
	err = repo.Client(ctx).WithContext(ctx).
		Where("id = ? and organization_id = ?", id, organizationId).
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

func (repo *OrganizationUserRepo) UpdateRolesWithTx(ctx context.Context, tx *gorm.DB, user *model.OrganizationUser, newRoles []*model.Role, removedRoles []*model.Role) (err error) {
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
