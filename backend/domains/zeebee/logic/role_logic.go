package logic

import (
	"context"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/zeromicro/go-zero/core/threading"
	"gorm.io/gorm"

	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/gormx"
	"zeebee-backend/domains/zeebee/auth"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
)

// RoleLogic provides role business logic
type RoleLogic struct {
	*BaseLogic
	// enforcer *rbac.Enforcer
	authTopicMap map[constants.UserType]string
}

// NewRoleLogic creates a new role logic instance
func NewRoleLogic(baseLogic *BaseLogic, authTopicMap map[constants.UserType]string) *RoleLogic {
	return &RoleLogic{
		BaseLogic:    baseLogic,
		authTopicMap: authTopicMap,
		// enforcer:  enforcer,
	}
}

func (l *RoleLogic) checkRoleName(roleName string) (err error) {
	if roleName == constants.AdminRoleName {
		return errorx.New("Role name is reserved")
	}
	return err
}

// Page returns paginated roles
func (l *RoleLogic) Page(ctx context.Context, req dto.RolePageRequestDto) (result gormx.PageResponse[dto.RoleEntityDto], err error) {
	pageResult, err := l.repos.RoleRepo.Page(ctx, req)
	if err != nil {
		return gormx.PageResponse[dto.RoleEntityDto]{}, err
	}

	resultPtr, err := copierx.Trans[gormx.PageResponse[*model.Role], gormx.PageResponse[dto.RoleEntityDto]](&pageResult)
	if err != nil {
		return gormx.PageResponse[dto.RoleEntityDto]{}, err
	}

	return *resultPtr, err
}

// Create creates a new role
func (l *RoleLogic) Create(ctx context.Context, input dto.RoleCreateRequestDto) (result *dto.RoleEntityDto, err error) {
	if err = l.checkRoleName(input.Name); err != nil {
		return nil, err
	}

	role, err := copierx.Trans[dto.RoleCreateRequestDto, model.Role](&input)
	if err != nil {
		return nil, err
	}
	role.RemoveRelationships()

	// check existed
	isExisted, err := l.repos.RoleRepo.IsExisted(ctx, role)
	if err != nil {
		return nil, err
	}
	if isExisted {
		return nil, errorx.New("Role already exists")
	}

	// C4: a new role only binds the resources the admin explicitly lists in
	// RoleAndResources. Plan resources are NOT auto-inherited from the org's
	// active subscription — the org admin configures role-resource bindings
	// manually (the subscription only defines the available plan-resource scope).
	resourceIds := slice.Map[dto.RoleAndResourceDto, int64](input.RoleAndResources, func(i int, item dto.RoleAndResourceDto) int64 {
		return item.ResourceId
	})
	var resources []*model.Resource
	var rolesAndResources []*model.RoleAndResource
	err = l.repos.RoleRepo.Transaction(ctx, func(tx *gorm.DB) error {
		var innerErr error

		var validResourceIds []int64
		// create role
		role, innerErr = l.repos.RoleRepo.CreateWithTx(ctx, tx, role, false)
		if innerErr != nil {
			return innerErr
		}

		// create role and resource relationship
		resources, innerErr = l.repos.ResourceRepo.GetResourcesWithTx(ctx, tx, resourceIds, input.Type)
		if innerErr != nil {
			return innerErr
		}
		validResourceIds = slice.Map[*model.Resource, int64](resources, func(i int, item *model.Resource) int64 {
			return item.Id
		})
		rolesAndResources = slice.FilterMap[dto.RoleAndResourceDto, *model.RoleAndResource](input.RoleAndResources, func(i int, item dto.RoleAndResourceDto) (*model.RoleAndResource, bool) {
			if slice.Contain(validResourceIds, item.ResourceId) {
				return &model.RoleAndResource{
					RoleId:     role.Id,
					ResourceId: item.ResourceId,
					DataScope:  item.DataScope,
				}, true
			}
			return nil, false
		})
		if len(rolesAndResources) > 0 {
			_, innerErr = l.repos.RoleAndResourceRepo.BatchCreateWithTx(ctx, tx, rolesAndResources, false)
			if innerErr != nil {
				return innerErr
			}
		}

		if innerErr != nil {
			return innerErr
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	threading.GoSafe(func() {
		subCtx := context.WithoutCancel(ctx)
		msg := auth.Msg{
			Action:         auth.ActionRoleReload,
			Scope:          auth.ScopeOrganization,
			Type:           input.Type,
			OrganizationId: input.OrganizationId,
		}
		policies := l.createRolePolicy(role.Id, input.Type, rolesAndResources)
		if len(policies) > 0 && input.Status == constants.CommonStatusActive {
			msg.AddedPolicies = policies
		}
		_ = l.publishAuthEvent(subCtx, l.authTopicMap[input.Type], msg)
	})

	result, err = l.Detail(ctx, role.Id, role.OrganizationId, role.Type)
	if err != nil {
		return nil, err
	}
	return result, err
}

// Update updates a role
func (l *RoleLogic) Update(ctx context.Context, id int64, input dto.RoleUpdateRequestDto) (result *dto.RoleEntityDto, err error) {
	if err = l.checkRoleName(input.Name); err != nil {
		return nil, err
	}
	// First get the existing record
	roleModel, err := l.repos.RoleRepo.GetRole(ctx, id, input.OrganizationId, input.Type)
	if err != nil {
		return nil, err
	}

	err = copierx.Assign[dto.RoleUpdateRequestDto, model.Role](&input, roleModel)
	if err != nil {
		return nil, err
	}

	// check existed
	isExisted, err := l.repos.RoleRepo.IsExisted(ctx, roleModel)
	if err != nil {
		return nil, err
	}
	if isExisted {
		return nil, errorx.New("Role already exists")
	}

	resourceIds := slice.Map[dto.RoleAndResourceDto, int64](input.RoleAndResources, func(i int, item dto.RoleAndResourceDto) int64 {
		return item.ResourceId
	})
	var addPolicies [][]string
	var removedPolicies [][]string
	err = l.repos.RoleRepo.Transaction(ctx, func(tx *gorm.DB) error {
		var innerErr error
		var resources []*model.Resource
		var validResourceIds []int64
		// update role
		roleModel.RemoveRelationships()
		roleModel, innerErr = l.repos.RoleRepo.UpdateOneWithTx(ctx, tx, roleModel)
		if innerErr != nil {
			return innerErr
		}

		// save role and resource relationship
		resources, innerErr = l.repos.ResourceRepo.GetResourcesWithTx(ctx, tx, resourceIds, roleModel.Type)
		if innerErr != nil {
			return innerErr
		}
		validResourceIds = slice.Map[*model.Resource, int64](resources, func(i int, item *model.Resource) int64 {
			return item.Id
		})
		rolesAndResources := slice.FilterMap[dto.RoleAndResourceDto, *model.RoleAndResource](input.RoleAndResources, func(i int, item dto.RoleAndResourceDto) (*model.RoleAndResource, bool) {
			if slice.Contain(validResourceIds, item.ResourceId) {
				return &model.RoleAndResource{
					RoleId:     roleModel.Id,
					ResourceId: item.ResourceId,
					DataScope:  item.DataScope,
				}, true
			}
			return nil, false
		})
		oldRelations, innerErr := l.repos.RoleAndResourceRepo.GetRoleRelationsWithTx(ctx, tx, roleModel.Id)
		if innerErr != nil {
			return innerErr
		}

		updateRelations := []*model.RoleAndResource{}
		newRelations := []*model.RoleAndResource{}
		slice.ForEach(rolesAndResources, func(index int, item *model.RoleAndResource) {
			oldRelation, ok := slice.FindBy(oldRelations, func(subIndex int, subItem *model.RoleAndResource) bool {
				return subItem.ResourceId == item.ResourceId
			})
			if ok {
				rolesAndResources[index].Id = oldRelation.Id
				updateRelations = append(updateRelations, rolesAndResources[index])
			} else {
				newRelations = append(newRelations, rolesAndResources[index])
			}
		})
		removedRelations := slice.FilterMap[*model.RoleAndResource, *model.RoleAndResource](oldRelations, func(index int, item *model.RoleAndResource) (*model.RoleAndResource, bool) {
			_, ok := slice.FindBy(rolesAndResources, func(subIndex int, subItem *model.RoleAndResource) bool {
				return subItem.ResourceId == item.ResourceId
			})
			if ok {
				return nil, false
			}
			return oldRelations[index], true
		})
		if len(newRelations) > 0 {
			_, innerErr = l.repos.RoleAndResourceRepo.BatchCreateWithTx(ctx, tx, newRelations, false)
			if innerErr != nil {
				return innerErr
			}
			// M6: only an active role's new bindings produce AddedPolicies. A disabled
			// role is filtered out during user-role resolution, so its policies would
			// be inert; skip pushing them to avoid stale enforcer state.
			if roleModel.Status == constants.CommonStatusActive {
				addPolicies = l.createRolePolicy(input.Id, input.Type, newRelations)
			}
		}
		if len(updateRelations) > 0 {
			for _, updatedRelation := range updateRelations {
				_, innerErr = l.repos.RoleAndResourceRepo.UpdateOneWithTx(ctx, tx, updatedRelation)
				if innerErr != nil {
					return innerErr
				}
			}
		}
		if len(removedRelations) > 0 {
			removedIds := slice.Map[*model.RoleAndResource, int64](removedRelations, func(i int, item *model.RoleAndResource) int64 {
				return item.Id
			})
			_, innerErr = l.repos.RoleAndResourceRepo.DeleteByIDsWithTx(ctx, tx, removedIds)
			if innerErr != nil {
				return innerErr
			}
			removedPolicies = l.createRolePolicy(input.Id, input.Type, removedRelations)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	threading.GoSafe(func() {
		subCtx := context.WithoutCancel(ctx)
		msg := auth.Msg{
			Action:          auth.ActionRoleReload,
			Scope:           auth.ScopeOrganization,
			Type:            input.Type,
			OrganizationId:  input.OrganizationId,
			AddedPolicies:   addPolicies,
			RemovedPolicies: removedPolicies,
		}
		_ = l.publishAuthEvent(subCtx, l.authTopicMap[input.Type], msg)
	})

	result, err = l.Detail(ctx, roleModel.Id, roleModel.OrganizationId, roleModel.Type)
	return result, err
}

// Detail returns a role by id
func (l *RoleLogic) Detail(ctx context.Context, id int64, organizationId int64, roleType constants.UserType) (result *dto.RoleEntityDto, err error) {
	role, err := l.repos.RoleRepo.GetRole(ctx, id, organizationId, roleType)
	if err != nil {
		return nil, err
	}
	// get admin's resources
	if role.Name == constants.AdminRoleName {
		var resources []*model.Resource
		resources, err = l.repos.ResourceRepo.All(ctx, roleType, false)
		if err != nil {
			return nil, err
		}
		role.Resources = resources
	}

	result, err = copierx.Trans[model.Role, dto.RoleEntityDto](role)
	if err != nil {
		return nil, err
	}
	slice.ForEach(result.Resources, func(index int, item *dto.ResourceEntityDto) {
		relation, err := l.repos.RoleAndResourceRepo.GetRelation(ctx, role.Id, item.Id)
		if err != nil {
			return
		}
		if relation != nil {
			item.DataScope = relation.DataScope
		}
	})

	return result, nil
}

// Delete deletes a role
func (l *RoleLogic) Delete(ctx context.Context, id int64, organizationId int64, roleType constants.UserType) (err error) {
	role, err := l.repos.RoleRepo.GetRole(ctx, id, organizationId, roleType)
	if err != nil {
		return err
	}
	if role.Name == constants.AdminRoleName {
		return errorx.New("Cannot delete admin role")
	}
	var removedPolicies [][]string
	err = l.repos.RoleRepo.Transaction(ctx, func(tx *gorm.DB) error {
		var innerErr error
		oldRelations, innerErr := l.repos.RoleAndResourceRepo.GetRoleRelationsWithTx(ctx, tx, role.Id)
		if innerErr != nil {
			return innerErr
		}
		innerErr = l.repos.RoleRepo.DeleteByIdWithTx(ctx, tx, id)
		if innerErr != nil {
			return innerErr
		}
		// remove role policies
		removedPolicies = l.createRolePolicy(id, roleType, oldRelations)

		return nil
	})
	if err != nil {
		return err
	}

	threading.GoSafe(func() {
		subCtx := context.WithoutCancel(ctx)
		msg := auth.Msg{
			Action:          auth.ActionRoleReload,
			Scope:           auth.ScopeOrganization,
			Type:            roleType,
			OrganizationId:  organizationId,
			RemovedPolicies: removedPolicies,
		}
		_ = l.publishAuthEvent(subCtx, l.authTopicMap[roleType], msg)
	})

	return err
}

// UpdateStatus updates a role's status
func (l *RoleLogic) SwitchStatus(ctx context.Context, id int64, organizationId int64, roleType constants.UserType, status constants.CommonStatus) (err error) {
	// First get the existing record
	result, err := l.repos.RoleRepo.GetRole(ctx, id, organizationId, roleType)
	if err != nil {
		return err
	}
	if result.Name == constants.AdminRoleName {
		return errorx.New("Cannot update admin role")
	}

	result.Status = status
	var addPolicies [][]string
	var removedPolicies [][]string

	err = l.repos.RoleRepo.Transaction(ctx, func(tx *gorm.DB) error {
		var innerErr error
		oldRelations, innerErr := l.repos.RoleAndResourceRepo.GetRoleRelationsWithTx(ctx, tx, id)
		if innerErr != nil {
			return innerErr
		}
		_, innerErr = l.repos.RoleRepo.UpdateOneByIdWithTx(ctx, tx, id, result)
		if innerErr != nil {
			return innerErr
		}

		switch status {
		case constants.CommonStatusActive:
			addPolicies = l.createRolePolicy(id, roleType, oldRelations)
		case constants.CommonStatusDisabled:
			removedPolicies = l.createRolePolicy(id, roleType, oldRelations)
		default:
			innerErr = errorx.Errorf("Invalid status: %v", status)
			return innerErr
		}

		return nil
	})

	if err != nil {
		return err
	}

	threading.GoSafe(func() {
		subCtx := context.WithoutCancel(ctx)
		msg := auth.Msg{
			Action:          auth.ActionRoleReload,
			Scope:           auth.ScopeOrganization,
			Type:            roleType,
			OrganizationId:  organizationId,
			RemovedPolicies: removedPolicies,
			AddedPolicies:   addPolicies,
		}
		_ = l.publishAuthEvent(subCtx, l.authTopicMap[roleType], msg)
	})

	return nil
}

// All returns all roles by type
func (l *RoleLogic) All(ctx context.Context, organizationId int64, roleType constants.UserType) (result []*dto.RoleEntityDto, err error) {
	roles, err := l.repos.RoleRepo.All(ctx, organizationId, roleType)
	if err != nil {
		return nil, err
	}
	resultPtr, err := copierx.Trans[[]*model.Role, []*dto.RoleEntityDto](&roles)
	return *resultPtr, err
}

// AssignOrganizationUsers assign organization users for role
func (l *RoleLogic) AssignOrganizationUsers(ctx context.Context, input dto.RoleAssignOrganizationUsersRequestDto) (result *dto.RoleEntityDto, err error) {
	role, err := l.repos.RoleRepo.GetRole(ctx, input.Id, input.OrganizationId, input.Type)
	if err != nil {
		return nil, err
	}
	if role.Name == constants.AdminRoleName {
		ownerId, err := l.repos.OrganizationUserRepo.GetOwnerId(ctx, input.OrganizationId, input.Type)
		if err != nil {
			return nil, err
		}
		input.OrganizationUserIds = append(input.OrganizationUserIds, ownerId)
	}

	err = l.repos.RoleRepo.Transaction(ctx, func(tx *gorm.DB) error {
		var innerErr error
		organizationUsers, innerErr := l.repos.OrganizationUserRepo.GetOrganizationUsersByIdsWithTx(ctx, tx, input.OrganizationId, input.OrganizationUserIds)
		if innerErr != nil {
			return innerErr
		}
		newUsers := []*model.OrganizationUser{}
		removedUsers := []*model.OrganizationUser{}
		newUsers = slice.DifferenceWith(organizationUsers, role.OrganizationUsers, func(item1, item2 *model.OrganizationUser) bool {
			return item1.Id == item2.Id
		})
		removedUsers = slice.DifferenceWith(role.OrganizationUsers, organizationUsers, func(item1, item2 *model.OrganizationUser) bool {
			return item1.Id == item2.Id
		})
		innerErr = l.repos.RoleRepo.UpdateOrganizationUsersWithTx(ctx, tx, role, newUsers, removedUsers)
		if innerErr != nil {
			return innerErr
		}
		role.OrganizationUsers = organizationUsers
		return innerErr
	})
	if err != nil {
		return nil, err
	}

	threading.GoSafe(func() {
		subCtx := context.WithoutCancel(ctx)
		msg := auth.Msg{
			Action:         auth.ActionRoleReload,
			Scope:          auth.ScopeOrganization,
			Type:           input.Type,
			OrganizationId: input.OrganizationId,
		}
		_ = l.publishAuthEvent(subCtx, l.authTopicMap[input.Type], msg)
	})

	result, err = l.Detail(ctx, role.Id, role.OrganizationId, role.Type)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (l *RoleLogic) createRolePolicy(roleId int64, roleType constants.UserType, relations []*model.RoleAndResource) [][]string {
	return slice.Map[*model.RoleAndResource, []string](relations, func(i int, item *model.RoleAndResource) []string {
		return auth.BuildPolicy(roleId, item.ResourceId, roleType)
	})
}
