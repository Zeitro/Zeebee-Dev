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

type RolePageReq struct {
	Page   int64              `json:"page"`
	Limit  int64              `json:"limit"`
	Filter *RolePageReqFilter `json:"filter"`
	Sorter *RolePageReqSorter `json:"sorter"`
}

type RolePageReqFilter struct {
	UserId         int64              `json:"user_id"`
	OrganizationId int64              `json:"organization_id"`
	Keyword        string             `json:"keyword"`
	RoleType       constants.UserType `json:"role_type"`
}
type RolePageReqSorter struct {
	Column string `json:"column"`
	Sort   string `json:"sort"`
}

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
func (repo *RoleRepo) Page(ctx context.Context, req *RolePageReq) (page gormx.PageResponse[*model.Role], err error) {
	filter := req.Filter
	sorter := req.Sorter
	if filter == nil {
		filter = &RolePageReqFilter{}
	}
	if sorter == nil {
		sorter = &RolePageReqSorter{}
	}
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
