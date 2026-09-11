package repo

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"

	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/gormx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
)

type RoleWithResource struct {
	Id              int64                  `json:"id"`
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at"`
	Name            string                 `json:"name"`
	Type            constants.UserType     `json:"type"`
	Status          constants.CommonStatus `json:"status"`
	OrganizationId  int64                  `json:"organization_id"`
	ResourceId      int64                  `json:"resource_id"`
	ResourceType    string                 `json:"resource_type"`
	ResourceName    string                 `json:"resource_name"`
	ResourceTag     string                 `json:"resource_tag"`
	ResourcePath    string                 `json:"resource_path"`
	ResourceStatus  constants.CommonStatus `json:"resource_status"`
	EnableDataScope bool                   `json:"enable_data_scope"`
	DataScope       constants.DataScope    `json:"data_scope"`
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
		Preload("OrganizationUsers").
		First(&result).Error
	return result, errorx.WithStack(err)
}

func (repo *RoleRepo) GetOrganizationRolesByIds(ctx context.Context, organizationId int64, roleIds []int64, roleType constants.UserType, withPreload bool) (result []*model.Role, err error) {
	if len(roleIds) == 0 {
		return []*model.Role{}, nil
	}

	query := repo.Client(ctx).Model(new(model.Role)).
		Where("id in (?) and organization_id = ? and type = ?", roleIds, organizationId, roleType)
	if withPreload {
		query = query.
			Preload("Resources")
	}
	err = query.Find(&result).Error
	return result, errorx.WithStack(err)
}

// GetOrganizationRolesByIdsWithTx return roles by ids
func (repo *RoleRepo) GetOrganizationRolesByIdsWithTx(ctx context.Context, tx *gorm.DB, organizationId int64, roleIds []int64, roleType constants.UserType) (result []*model.Role, err error) {
	err = tx.WithContext(ctx).Model(new(model.Role)).
		Where("id in (?) and organization_id = ? and type = ?", roleIds, organizationId, roleType).
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
	query = query.Order("CASE status WHEN 'Active' THEN 1 WHEN 'Disabled' THEN 2 ELSE 3 END").
		Order(fmt.Sprintf("%s %s", sorter.Column, sorter.Sort))

	result, err := gormx.Page[*model.Role](query, gormx.PageRequest{
		Page:        req.Page,
		Limit:       req.Limit,
		CountColumn: "id",
	})
	return result, err
}

// GetRolesByOrganizationUserId
func (repo *RoleRepo) GetRolesByOrganizationUserId(ctx context.Context, organizationUserId int64) (result []*model.Role, err error) {
	roleTable := model.Role{}.TableName()
	pivotTable := "organization_user_and_roles"
	err = repo.Client(ctx).
		Joins(fmt.Sprintf("left join %s on %s.id = %s.role_id", pivotTable, roleTable, pivotTable)).
		Where(fmt.Sprintf("%s.organization_user_id = ?", pivotTable), organizationUserId).
		Find(&result).Error
	return result, errorx.WithStack(err)
}

// FindByOrganizationIdsWithTx returns all roles for the given organization IDs within a transaction.
func (repo *RoleRepo) FindByOrganizationIdsWithTx(ctx context.Context, tx *gorm.DB, organizationIds []int64, roleType constants.UserType) (result []*model.Role, err error) {
	if len(organizationIds) == 0 {
		return []*model.Role{}, nil
	}
	err = tx.WithContext(ctx).
		Where("organization_id IN ? AND type = ?", organizationIds, roleType).
		Preload("Resources").
		Find(&result).Error
	return result, errorx.WithStack(err)
}

// UpdateResourcesWithTx updates a role's resource associations within a transaction.
// New resources are appended, removed resources are deleted, unchanged ones are left alone.
// Uses GORM Association to manage the role_and_resources junction table.
func (repo *RoleRepo) UpdateResourcesWithTx(ctx context.Context, tx *gorm.DB, role *model.Role, newResources []*model.Resource, removedResources []*model.Resource) error {
	copyRole, err := copierx.Trans[model.Role, model.Role](role)
	if err != nil {
		return err
	}
	copyRole.RemoveRelationships()
	query := tx.WithContext(ctx).Model(copyRole)
	if len(newResources) > 0 {
		err = query.Association("Resources").Append(newResources)
		if err != nil {
			return errorx.WithStack(err)
		}
	}
	copyRole.RemoveRelationships()
	if len(removedResources) > 0 {
		err = query.Association("Resources").Delete(removedResources)
		if err != nil {
			return errorx.WithStack(err)
		}
	}
	return nil
}

// DeleteByIdWithTx deletes a role by its ID within a transaction.
// It also removes all associated resources and organization users before deletion.
// The function ensures all operations are executed within the provided transaction context.
func (repo *RoleRepo) DeleteByIdWithTx(ctx context.Context, tx *gorm.DB, id int64) (err error) {
	var role *model.Role
	ctxQuery := tx.WithContext(ctx)
	err = ctxQuery.
		Model(&model.Role{}).
		Preload("Resources").
		Preload("OrganizationUsers").
		Where("id = ?", id).
		First(&role).Error
	if err != nil {
		return errorx.WithStack(err)
	}
	// remove resource relations
	err = ctxQuery.Model(role).Association("Resources").Clear()
	if err != nil {
		return errorx.WithStack(err)
	}
	// remove organization user relations
	err = ctxQuery.Model(role).Association("OrganizationUsers").Clear()
	if err != nil {
		return errorx.WithStack(err)
	}
	role.RemoveRelationships()
	err = ctxQuery.Delete(&role).Error

	return errorx.WithStack(err)
}

// GetOrganizationRoleWithResources
func (repo *RoleRepo) GetOrganizationRoleWithResources(ctx context.Context, organizationId int64, roleType constants.UserType) (results []RoleWithResource, err error) {
	roleTable := model.Role{}.TableName()
	resourceTable := model.Resource{}.TableName()
	relationTable := model.RoleAndResource{}.TableName()
	fields := []string{
		fmt.Sprintf("%s.id as id", roleTable),
		fmt.Sprintf("%s.created_at as created_at", roleTable),
		fmt.Sprintf("%s.updated_at as updated_at", roleTable),
		fmt.Sprintf("%s.name as name", roleTable),
		fmt.Sprintf("%s.type as type", roleTable),
		fmt.Sprintf("%s.status as status", roleTable),
		fmt.Sprintf("%s.organization_id as organization_id", roleTable),

		fmt.Sprintf("%s.id as resource_id", resourceTable),
		fmt.Sprintf("%s.name as resource_name", resourceTable),
		fmt.Sprintf("%s.path as resource_path", resourceTable),
		fmt.Sprintf("%s.status as resource_status", resourceTable),
		fmt.Sprintf("%s.tag as resource_tag", resourceTable),
		fmt.Sprintf("%s.type as resource_type", resourceTable),
		fmt.Sprintf("%s.enable_data_scope as enable_data_scope", resourceTable),

		fmt.Sprintf("%s.data_scope as data_scope", relationTable),
	}

	err = repo.Client(ctx).Select(fields).
		Joins(fmt.Sprintf("left join %s on %s.id = %s.role_id", relationTable, roleTable, relationTable)).
		Joins(fmt.Sprintf("left join %s on %s.resource_id = %s.id", resourceTable, relationTable, resourceTable)).
		Where(
			fmt.Sprintf("%s.organization_id = ? and %s.type = ? and %s.status = ? and %s.status = ?", roleTable, roleTable, roleTable, resourceTable),
			organizationId, roleType, constants.CommonStatusActive, constants.CommonStatusActive).
		Where(fmt.Sprintf("%s.is_del = 0 and %s.is_del = 0", roleTable, resourceTable)).
		Scan(&results).Error

	return results, errorx.WithStack(err)
}

// UpdateOrganizationUsersWithTx update role's organization users with transaction
func (repo *RoleRepo) UpdateOrganizationUsersWithTx(ctx context.Context, tx *gorm.DB, role *model.Role, newUsers []*model.OrganizationUser, removedUsers []*model.OrganizationUser) (err error) {
	if role == nil {
		return errorx.New("role is empty")
	}

	copyRole, err := copierx.Trans[model.Role, model.Role](role)
	if err != nil {
		return err
	}

	copyRole.RemoveRelationships()
	query := tx.WithContext(ctx).Model(copyRole)
	if len(newUsers) > 0 {
		err = query.Association("OrganizationUsers").Append(newUsers)
		if err != nil {
			return errorx.WithStack(err)
		}
	}
	copyRole.RemoveRelationships()
	if len(removedUsers) > 0 {
		err = query.Association("OrganizationUsers").Delete(removedUsers)
		if err != nil {
			return errorx.WithStack(err)
		}
	}
	return nil
}
