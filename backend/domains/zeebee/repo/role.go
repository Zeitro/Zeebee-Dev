package repo

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"zeebee-backend/common/constants"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/gormx"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
)

type RoleRepo struct {
	*gormx.BaseRepo[*model.Role]
	m *model.Role
}

func NewRoleRepo(db *gorm.DB) *RoleRepo {
	return &RoleRepo{
		BaseRepo: gormx.NewBaseRepo[*model.Role](db),
		m:        new(model.Role),
	}
}

// All returns all roles by type
func (repo *RoleRepo) All(ctx context.Context, organizationId int64, roleType constants.UserType) (result []*model.Role, err error) {
	err = repo.Client(ctx).
		Where("organization_id = ? and type = ? and status = ?", organizationId, roleType, constants.CommonStatusActive).
		Order("id asc").
		Find(&result).Error
	return result, errorx.WithStack(err)
}

// IsExisted check role is existed
func (repo *RoleRepo) IsExisted(ctx context.Context, data *model.Role) (isExisted bool, err error) {
	query := repo.Client(ctx).
		Where("name = ?", data.Name).
		Where("type = ?", data.Type).
		Where("organization_id = ?", data.OrganizationId)
	if data.Id > 0 {
		query = query.Where("id != ?", data.Id)
	}
	var count int64
	err = query.Count(&count).Error
	return count > 0, errorx.WithStack(err)
}

// GetRole return role by id
func (repo *RoleRepo) GetRole(ctx context.Context, id int64, organizationId int64, roleType constants.UserType) (result *model.Role, err error) {
	err = repo.Client(ctx).
		Where("id = ?", id).
		Where("organization_id = ?", organizationId).
		Where("type = ?", roleType).
		Preload("Resources").
		First(&result).Error
	return result, errorx.WithStack(err)
}

// GetOrganizationRolesByIdsWithTx return roles by ids
func (repo *RoleRepo) GetOrganizationRolesByIdsWithTx(ctx context.Context, tx *gorm.DB, organizationId int64, roleIds []int64) (result []*model.Role, err error) {
	err = tx.WithContext(ctx).Model(new(model.Role)).
		Where("id in (?) and organization_id = ?", roleIds, organizationId).
		Preload("Resources").
		Find(&result).Error
	return result, errorx.WithStack(err)

}

// Page returns paginated roles
func (repo *RoleRepo) Page(ctx context.Context, req dto.RolePageRequestDto) (page gormx.PageResponse[*model.Role], err error) {
	filter := req.Filter
	sorter := req.Sorter

	query := repo.Client(ctx).Model(new(model.Role)).
		Where("type = ?", filter.RoleType).
		Where("organization_id = ?", filter.OrganizationId)

	if len(filter.Keyword) > 0 {
		query = query.Where("name LIKE ?", "%"+filter.Keyword+"%")
	}
	if len(sorter.Column) == 0 {
		sorter.Column = "id"
	}
	if len(sorter.Sort) == 0 {
		sorter.Sort = "asc"
	}
	query = query.Order("CASE status WHEN 'Active' THEN 1 WHEN 'Inactive' THEN 2 ELSE 3 END").
		Order(fmt.Sprintf("%s %s", sorter.Column, sorter.Sort))

	result, err := gormx.Page[*model.Role](query, gormx.PageRequest{
		Page:        req.Page,
		Limit:       req.Limit,
		CountColumn: "id",
	})
	return result, err
}

// GetUserRoles
func (repo *RoleRepo) GetOrganizationUserRoles(ctx context.Context, organizationUserId int64) (result []*model.Role, err error) {
	roleTable := model.Role{}.TableName()
	pivotTable := "organization_user_and_roles"
	err = repo.Client(ctx).
		Joins(fmt.Sprintf("left join %s on %s.id = %s.role_id", pivotTable, roleTable, pivotTable)).
		Where(fmt.Sprintf("%s.organization_user_id = ?", pivotTable), organizationUserId).
		Find(&result).Error
	return result, errorx.WithStack(err)
}
