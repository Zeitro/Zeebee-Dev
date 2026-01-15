package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/jinzhu/copier"
	"gorm.io/gorm"

	"zeebee-backend/common/errorx"
	"zeebee-backend/common/gormx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
)

type OrganizationUserPageEntity struct {
	Id             int64                            `gorm:"column:id;default:0" json:"id"`
	CreatedAt      time.Time                        `gorm:"column:created_at;default:''" json:"created_at"`
	UpdatedAt      time.Time                        `gorm:"column:updated_at;default:''" json:"updated_at"`
	Email          string                           `gorm:"column:email;default:''" json:"email"`
	OrganizationId int64                            `gorm:"column:organization_id;default:0" json:"organization_id"`
	UserId         int64                            `gorm:"column:user_id;default:0" json:"user_id"`
	Status         constants.OrganizationUserStatus `gorm:"column:status;default:'Not Activated'" json:"status"`
	FirstName      string                           `gorm:"column:first_name;default:''" json:"first_name"`
	MidName        string                           `gorm:"column:mid_name;default:''" json:"mid_name"`
	LastName       string                           `gorm:"column:last_name;default:''" json:"last_name"`
	FullName       string                           `gorm:"column:full_name;default:''" json:"full_name"`
	Phone          string                           `gorm:"column:phone;default:''" json:"phone"`
	Remark         string                           `gorm:"column:remark;default:''" json:"remark"`

	Roles []*model.Role `gorm:"-" json:"roles"`
}

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
func (repo *OrganizationUserRepo) Page(ctx context.Context, pageReq dto.OrganizationUserPageRequestDto) (page gormx.PageResponse[*OrganizationUserPageEntity], err error) {
	filter := pageReq.Filter
	sorter := pageReq.Sorter
	userTable := model.User{}.TableName()
	organizationUserAndRoleTable := "organization_user_and_roles"
	organizationUserTable := model.OrganizationUser{}.TableName()
	fields := []string{
		fmt.Sprintf("%s.id", organizationUserTable),
		fmt.Sprintf("%s.user_id", organizationUserTable),
		fmt.Sprintf("%s.organization_id", organizationUserTable),
		fmt.Sprintf("%s.created_at", organizationUserTable),
		fmt.Sprintf("%s.updated_at", organizationUserTable),
		fmt.Sprintf("%s.email", userTable),
		fmt.Sprintf("%s.phone", organizationUserTable),
		fmt.Sprintf("%s.status", organizationUserTable),
		fmt.Sprintf("%s.first_name", organizationUserTable),
		fmt.Sprintf("%s.mid_name", organizationUserTable),
		fmt.Sprintf("%s.last_name", organizationUserTable),
		fmt.Sprintf("%s.full_name", organizationUserTable),
		fmt.Sprintf("%s.remark", organizationUserTable),
	}

	query := repo.Client(ctx).Model(&model.OrganizationUser{}).Select(fields).
		Joins(fmt.Sprintf("left join %s on %s.id = %s.user_id", userTable, userTable, organizationUserTable)).
		Where(fmt.Sprintf("%s.organization_id = ?", organizationUserTable), pageReq.Filter.OrganizationId)
	if len(filter.Keyword) > 0 {
		query = query.Where(fmt.Sprintf("%s.email like ? or %s.full_name like ? or %s.phone like ?", userTable, organizationUserTable, organizationUserTable), "%"+filter.Keyword+"%", "%"+filter.Keyword+"%", "%"+filter.Keyword+"%")
	}

	if filter.RoleId > 0 {
		query = query.Joins(fmt.Sprintf("left join %s on %s.id = %s.organization_user_id", organizationUserAndRoleTable, organizationUserTable, organizationUserAndRoleTable)).
			Where(fmt.Sprintf("%s.role_id = ?", organizationUserAndRoleTable), filter.RoleId)
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
	return gormx.Page[*OrganizationUserPageEntity](query, gormx.PageRequest{
		Page:        pageReq.Page,
		Limit:       pageReq.Limit,
		CountColumn: fmt.Sprintf("%s.id", organizationUserTable),
	})
}

// GetOrganizationUserByUserId get organization user by organizationId and userId
func (repo *OrganizationUserRepo) GetOrganizationUserByUserId(ctx context.Context, organizationId int64, userId int64) (res *model.OrganizationUser, err error) {
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

// ActiveByUserIdWithTx active organization user by userId with transaction
func (repo *OrganizationUserRepo) ActiveByUserIdWithTx(ctx context.Context, tx *gorm.DB, userId int64) (err error) {
	return tx.WithContext(ctx).Model(&model.OrganizationUser{}).
		Where("user_id = ? and status = ?", userId, constants.OrganizationUserStatusNotActivated).
		Update("status", constants.OrganizationUserStatusActive).Error
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

// GetOrganizationUserByIdWithTx get organization user by id with transaction
func (repo *OrganizationUserRepo) GetOrganizationUserByIdWithTx(ctx context.Context, tx *gorm.DB, id int64, organizationId int64) (result *model.OrganizationUser, err error) {
	err = tx.WithContext(ctx).Model(&model.OrganizationUser{}).
		Where("id = ? and organization_id = ?", id, organizationId).
		Preload("User").
		Preload("Roles").
		First(&result).Error
	return result, errorx.WithStack(err)
}

func (repo *OrganizationUserRepo) DeleteOrganizationUserWithTx(ctx context.Context, tx *gorm.DB, organizationUser *model.OrganizationUser) (err error) {
	query := tx.WithContext(ctx).Model(organizationUser)
	err = query.Association("Roles").Clear()
	if err != nil {
		return errorx.WithStack(err)
	}
	err = tx.WithContext(ctx).Delete(organizationUser).Error
	return errorx.WithStack(err)
}

// UpdateOrganizationUserStatus update relation status
func (repo *OrganizationUserRepo) UpdateOrganizationUserStatus(ctx context.Context, id int64, organizationId int64, status constants.OrganizationUserStatus) (err error) {
	err = repo.Client(ctx).Model(&model.OrganizationUser{}).
		Where("id = ? and organization_id = ?", id, organizationId).
		Update("status", status).Error
	return errorx.WithStack(err)
}

// CheckOrganizationUserExistedByEmail check organization user exist
func (repo *OrganizationUserRepo) CheckOrganizationUserExistedByEmail(ctx context.Context, organizationId int64, email string, userType constants.UserType) (existed bool, err error) {
	userTable := model.User{}.TableName()
	organizationUserTable := model.OrganizationUser{}.TableName()
	var count int64
	err = repo.Client(ctx).
		Joins(fmt.Sprintf("left join %s on %s.id = %s.user_id", userTable, userTable, organizationUserTable)).
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

// UpdateRolesWithTx update organization user's roles with transaction
func (repo *OrganizationUserRepo) UpdateRolesWithTx(ctx context.Context, tx *gorm.DB, user *model.OrganizationUser, newRoles []*model.Role, removedRoles []*model.Role) (err error) {
	if user == nil {
		return errorx.New("User Is Empty")
	}
	copyUser := &model.OrganizationUser{}
	_ = copier.CopyWithOption(copyUser, user, copier.Option{DeepCopy: true})

	copyUser.RemoveRelationships()
	query := tx.WithContext(ctx).Model(copyUser)
	if len(newRoles) > 0 {
		err = query.Association("Roles").Append(newRoles)
		if err != nil {
			return errorx.WithStack(err)
		}
	}
	copyUser.RemoveRelationships()
	if len(removedRoles) > 0 {
		err = query.Association("Roles").Delete(removedRoles)
		if err != nil {
			return errorx.WithStack(err)
		}
	}
	return nil
}

// UpdateStatusByOrganizationId batch update organization user status by organization id
// Used for organization status switching - single SQL for performance with 1000+ users
func (repo *OrganizationUserRepo) UpdateStatusByOrganizationId(ctx context.Context, organizationId int64, status constants.OrganizationUserStatus) (err error) {
	err = repo.Client(ctx).Model(&model.OrganizationUser{}).
		Where("organization_id = ?", organizationId).
		Update("status", status).Error
	return errorx.WithStack(err)
}

// DeleteByOrganizationIdWithTx delete all organization users by organization id with transaction
// Used for organization deletion - cascading delete of all organization user relationships
func (repo *OrganizationUserRepo) DeleteByOrganizationIdWithTx(ctx context.Context, tx *gorm.DB, organizationId int64) (err error) {
	var organizationUser *model.OrganizationUser
	err = tx.WithContext(ctx).
		Where("organization_id = ?", organizationId).
		Preload("Roles").
		First(&organizationUser).Error
	if err != nil {
		return errorx.WithStack(err)
	}

	err = tx.WithContext(ctx).Model(organizationUser).Association("Roles").Clear()
	if err != nil {
		return errorx.WithStack(err)
	}
	organizationUser.RemoveRelationships()
	err = tx.WithContext(ctx).Delete(organizationUser).Error

	return errorx.WithStack(err)
}
