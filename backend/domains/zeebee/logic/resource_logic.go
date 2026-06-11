package logic

import (
	"context"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/zeromicro/go-zero/core/threading"
	"gorm.io/gorm"

	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/auth"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
)

// ResourceLogic provides resource business logic
type ResourceLogic struct {
	*BaseLogic
	authTopicMap map[constants.UserType]string
}

// NewResourceLogic creates a new resource logic instance
func NewResourceLogic(baseLogic *BaseLogic, authTopicMap map[constants.UserType]string) *ResourceLogic {
	return &ResourceLogic{
		BaseLogic:    baseLogic,
		authTopicMap: authTopicMap,
	}
}

// ResourceTree returns resource tree
func (l *ResourceLogic) ResourceTree(ctx context.Context, resourceType constants.UserType) (trees []*dto.ResourceTreeDto, err error) {
	results, err := l.repos.ResourceRepo.All(ctx, resourceType, true)
	if err != nil {
		return
	}
	trees = l.buildTree(results, 0)
	return
}

// GetOrganizationResources returns resource tree filtered by the organization's active pricing plan
func (l *ResourceLogic) GetOrganizationResources(ctx context.Context, organizationId int64, resourceType constants.UserType) (trees []*dto.ResourceTreeDto, err error) {
	results, err := l.repos.ResourceRepo.All(ctx, resourceType, true)
	if err != nil {
		return
	}

	results, err = l.filterByPricingPlan(ctx, organizationId, results)
	if err != nil {
		return
	}

	trees = l.buildTree(results, 0)
	return
}

// CreateResource creates a new resource
func (l *ResourceLogic) CreateResource(ctx context.Context, input *dto.ResourceCreateRequestDto) (result *dto.ResourceEntityDto, err error) {
	resourceModel, err := copierx.Trans[dto.ResourceCreateRequestDto, model.Resource](input)
	if err != nil {
		return nil, err
	}

	// Use transaction for create with associations
	err = l.repos.ResourceRepo.Transaction(ctx, func(tx *gorm.DB) error {
		var innerErr error
		// check parent resource
		if input.ParentId != 0 {
			isExisted, innerErr := l.repos.ResourceRepo.IsExistedByColumnsWithTx(ctx, tx, map[string]interface{}{
				"id":   input.ParentId,
				"type": input.Type,
			})
			if innerErr != nil {
				return innerErr
			}
			if !isExisted {
				return errorx.New("parent resource not found")
			}
		}

		// check existed
		isExisted, innerErr := l.repos.ResourceRepo.IsExistedWithTx(ctx, tx, resourceModel)
		if innerErr != nil {
			return innerErr
		}
		if isExisted {
			return errorx.New("resource already exists")
		}

		// check tag
		isTagExisted, innerErr := l.repos.ResourceRepo.IsTagExistedWithTx(ctx, tx, input.Tag, 0, input.Type)
		if innerErr != nil {
			return innerErr
		}
		if isTagExisted {
			return errorx.New("tag already exists")
		}

		// Get apis by ids
		if len(input.ApiIds) > 0 {
			var apis []*model.Api
			apis, innerErr = l.repos.ApiRepo.GetApisWithTx(ctx, tx, input.ApiIds, input.Type)
			if innerErr != nil {
				return innerErr
			}
			resourceModel.Apis = apis
		}

		resourceModel, innerErr = l.repos.ResourceRepo.CreateWithTx(ctx, tx, resourceModel, true)
		return innerErr
	})

	if err != nil {
		return nil, err
	}
	threading.GoSafe(func() {
		subCtx := context.WithoutCancel(ctx)
		msg := auth.Msg{
			Action: auth.ActionResourceReload,
			Scope:  auth.ScopeGlobal,
			Type:   input.Type,
		}
		_ = l.publishAuthEvent(subCtx, l.authTopicMap[input.Type], msg)
	})

	result, err = copierx.Trans[model.Resource, dto.ResourceEntityDto](resourceModel)
	return result, err
}

// UpdateResource updates an existing resource
func (l *ResourceLogic) UpdateResource(ctx context.Context, input *dto.ResourceUpdateRequestDto) (result *dto.ResourceEntityDto, err error) {
	var resourceModel *model.Resource
	var oldRoles []*model.Role
	var isStatusChanged bool

	err = l.repos.ResourceRepo.Transaction(ctx, func(tx *gorm.DB) error {
		var innerErr error

		// check parent resource
		if input.ParentId != 0 {
			isExisted, innerErr := l.repos.ResourceRepo.IsExistedByColumnsWithTx(ctx, tx, map[string]interface{}{
				"id":   input.ParentId,
				"type": input.Type,
			})
			if innerErr != nil {
				return innerErr
			}
			if !isExisted {
				return errorx.New("parent resource not found")
			}
		}

		resourceModel, innerErr = l.repos.ResourceRepo.GetResourceWithTx(ctx, tx, input.Id, input.Type, true)
		if innerErr != nil {
			return innerErr
		}
		oldRoles = resourceModel.Roles
		isStatusChanged = input.Status != resourceModel.Status

		innerErr = copierx.Assign[dto.ResourceUpdateRequestDto, model.Resource](input, resourceModel)
		if innerErr != nil {
			return innerErr
		}

		// check existed
		isExisted, innerErr := l.repos.ResourceRepo.IsExistedWithTx(ctx, tx, resourceModel)
		if innerErr != nil {
			return innerErr
		}
		if isExisted {
			return errorx.New("resource already exists")
		}

		// check tag
		isTagExisted, innerErr := l.repos.ResourceRepo.IsTagExistedWithTx(ctx, tx, input.Tag, resourceModel.Id, input.Type)
		if innerErr != nil {
			return innerErr
		}
		if isTagExisted {
			return errorx.New("tag already exists")
		}

		// handle apis
		oldApis := resourceModel.Apis
		var currentApis []*model.Api
		if len(input.ApiIds) > 0 {
			currentApis, innerErr = l.repos.ApiRepo.GetApisWithTx(ctx, tx, input.ApiIds, input.Type)
			if innerErr != nil {
				return innerErr
			}
		}

		removedApis := slice.DifferenceWith(oldApis, currentApis, func(item1, item2 *model.Api) bool {
			return item1.Id == item2.Id
		})
		newApis := slice.DifferenceWith(currentApis, oldApis, func(item1, item2 *model.Api) bool {
			return item1.Id == item2.Id
		})

		resourceModel.RemoveRelationships()
		innerErr = l.repos.ResourceRepo.RemoveApisWithTx(ctx, tx, resourceModel, removedApis)
		if innerErr != nil {
			return innerErr
		}
		innerErr = l.repos.ResourceRepo.AddApisWithTx(ctx, tx, resourceModel, newApis)
		if innerErr != nil {
			return innerErr
		}
		resourceModel.RemoveRelationships()
		// update resource
		resourceModel, innerErr = l.repos.ResourceRepo.UpdateOneByIdWithTx(ctx, tx, input.Id, resourceModel)
		if innerErr != nil {
			return innerErr
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	threading.GoSafe(func() {
		var removedPolicies [][]string
		var addedPolicies [][]string
		subCtx := context.WithoutCancel(ctx)

		if isStatusChanged && input.Status == constants.CommonStatusDisabled {
			removedPolicies = l.createResourcePolices(input.Id, input.Type, oldRoles)
		}
		if isStatusChanged && input.Status == constants.CommonStatusActive {
			addedPolicies = l.createResourcePolices(input.Id, input.Type, oldRoles)
		}

		msg := auth.Msg{
			Action:          auth.ActionResourceReload,
			Scope:           auth.ScopeGlobal,
			Type:            input.Type,
			AddedPolicies:   addedPolicies,
			RemovedPolicies: removedPolicies,
		}

		_ = l.publishAuthEvent(subCtx, l.authTopicMap[input.Type], msg)
	})

	result, err = l.Detail(ctx, resourceModel.Id, resourceModel.Type)
	return result, err
}

// Detail returns resource detail
func (l *ResourceLogic) Detail(ctx context.Context, id int64, resourceType constants.UserType) (result *dto.ResourceEntityDto, err error) {
	resource, err := l.repos.ResourceRepo.GetResource(ctx, id, resourceType, true)
	if err != nil {
		return nil, err
	}
	result, err = copierx.Trans[model.Resource, dto.ResourceEntityDto](resource)
	return result, err
}

// DeleteResource deletes resource and its children recursively
func (l *ResourceLogic) DeleteResource(ctx context.Context, id int64, resourceType constants.UserType) (err error) {
	var policies [][]string
	err = l.repos.ResourceRepo.Transaction(ctx, func(tx *gorm.DB) error {
		var innerErr error
		policies, innerErr = l.deleteResource(ctx, tx, id, resourceType)
		return innerErr
	})
	if err != nil {
		return err
	}

	// update user permissions
	threading.GoSafe(func() {
		subCtx := context.WithoutCancel(ctx)
		msg := auth.Msg{
			Action:          auth.ActionResourceReload,
			Scope:           auth.ScopeGlobal,
			Type:            resourceType,
			RemovedPolicies: policies,
		}

		_ = l.publishAuthEvent(subCtx, l.authTopicMap[resourceType], msg)
	})

	return nil
}

// GetUserResources returns user's resources as a tree structure
// 1. Only resources with status "Active" are included
// 2. Root level resources must be of type "Menu", nested resources can be "Menu" or "Button"
// 3. User may have multiple roles, resources from all roles are merged and deduplicated
// 4. If user has admin role, all organization resources are returned
func (l *ResourceLogic) GetUserResources(ctx context.Context, input dto.ResourceUserMenuRequestDto) (result []*dto.ResourceTreeDto, err error) {
	// Get user's roles
	organizationUser, err := l.repos.OrganizationUserRepo.GetOrganizationUserById(ctx, input.OrganizationUserId, input.OrganizationId)
	if err != nil {
		return nil, err
	}
	if organizationUser.Status != constants.OrganizationUserStatusActive {
		return nil, errorx.New("user is not active")
	}
	// filter roles
	hasAdminRole := false
	roles := slice.Filter(organizationUser.Roles, func(index int, role *model.Role) bool {
		if role.Status != constants.CommonStatusActive {
			return false
		}
		if role.Name == constants.AdminRoleName {
			hasAdminRole = true
		}
		return true
	})

	var allResources []*model.Resource
	if hasAdminRole {
		// Get all active resources for the organization
		allResources, err = l.repos.ResourceRepo.AllResources(ctx, input.OrganizationId, input.Type, false)
		if err != nil {
			return nil, err
		}
	} else {
		// Get resources from user's roles
		publicResources, err := l.repos.ResourceRepo.PublicResources(ctx, input.Type, false)
		if err != nil {
			return nil, err
		}

		roleIds := slice.Map(roles, func(index int, role *model.Role) int64 {
			return role.Id
		})
		roleResources, err := l.repos.RoleRepo.GetOrganizationRolesByIds(ctx, input.OrganizationId, roleIds, input.Type, true)
		if err != nil {
			return nil, err
		}
		roleResources = slice.Filter(roleResources, func(index int, role *model.Role) bool {
			return role.Status == constants.CommonStatusActive
		})

		allResources = make([]*model.Resource, 0, len(roleResources))
		slice.ForEach(roleResources, func(index int, role *model.Role) {
			if role.Status != constants.CommonStatusActive {
				return
			}
			subResources := slice.Filter(role.Resources, func(index int, resource *model.Resource) bool {
				return resource.Status == constants.CommonStatusActive
			})
			allResources = append(allResources, subResources...)
		})
		allResources = append(allResources, publicResources...)

		allResources = slice.UniqueBy(allResources, func(item *model.Resource) int64 {
			return item.Id
		})
	}

	slice.SortBy(allResources, func(a, b *model.Resource) bool {
		if a.SortOrder != b.SortOrder {
			return a.SortOrder < b.SortOrder
		}
		return a.Id < b.Id
	})

	// Filter by pricing plan resource scope
	allResources, err = l.filterByPricingPlan(ctx, input.OrganizationId, allResources)
	if err != nil {
		return nil, err
	}

	// Filter root level to only menu type, collect all resource IDs
	rootResources := slice.Filter(allResources, func(index int, item *model.Resource) bool {
		return item.ParentId == 0 && item.ResourceType == constants.ResourceTypeMenu
	})

	// Build resource map for quick lookup and collect all descendant resource IDs
	resourceMap := make(map[int64]*model.Resource)
	for _, resource := range allResources {
		resourceMap[resource.Id] = resource
	}
	// Build tree with menu-type roots, allowing menu and button children
	result = l.buildUserResourceTree(resourceMap, rootResources)
	return result, nil
}

// deleteResource
func (l *ResourceLogic) deleteResource(ctx context.Context, tx *gorm.DB, id int64, resourceType constants.UserType) (policies [][]string, err error) {
	resource, err := l.repos.ResourceRepo.GetResourceWithTx(ctx, tx, id, resourceType, true)
	if err != nil {
		return nil, err
	}
	policies = l.createResourcePolices(resource.Id, resourceType, resource.Roles)

	children, err := l.repos.ResourceRepo.GetResourceByParentIdWithTx(ctx, tx, resource.Id, resourceType)
	if err != nil {
		return nil, err
	}
	for _, child := range children {
		subPolices, err := l.deleteResource(ctx, tx, child.Id, resourceType)
		if err != nil {
			return nil, err
		}
		policies = append(policies, subPolices...)
	}
	err = l.repos.ResourceRepo.DeleteResourceWithTx(ctx, tx, resource)
	return policies, err
}

func (l *ResourceLogic) buildTree(resources []*model.Resource, parentId int64) (results []*dto.ResourceTreeDto) {
	results = slice.FilterMap[*model.Resource, *dto.ResourceTreeDto](resources, func(index int, item *model.Resource) (*dto.ResourceTreeDto, bool) {
		if item.ParentId != parentId {
			return nil, false
		}
		tree, err := copierx.Trans[model.Resource, dto.ResourceTreeDto](item)
		if err != nil {
			return nil, false
		}
		return tree, true
	})

	for i, result := range results {
		results[i].Children = l.buildTree(resources, result.Id)
	}
	return results
}

// buildUserResourceTree builds user resource tree from resource map and root resources
// Allows children to be either Menu or Button type
func (l *ResourceLogic) buildUserResourceTree(resourceMap map[int64]*model.Resource, roots []*model.Resource) (results []*dto.ResourceTreeDto) {
	results = slice.FilterMap[*model.Resource, *dto.ResourceTreeDto](roots, func(index int, item *model.Resource) (*dto.ResourceTreeDto, bool) {
		tree, err := copierx.Trans[model.Resource, dto.ResourceTreeDto](item)
		if err != nil {
			return nil, false
		}
		return tree, true
	})

	for i, result := range results {
		results[i].Children = l.buildUserResourceChildren(resourceMap, result.Id)
	}
	return results
}

// buildUserResourceChildren builds children for a resource
// Children can be either Menu or Button type
func (l *ResourceLogic) buildUserResourceChildren(resourceMap map[int64]*model.Resource, parentId int64) (results []*dto.ResourceTreeDto) {
	for _, resource := range resourceMap {
		if resource.ParentId == parentId {
			tree, err := copierx.Trans[model.Resource, dto.ResourceTreeDto](resource)
			if err != nil {
				continue
			}
			tree.Children = l.buildUserResourceChildren(resourceMap, resource.Id)
			results = append(results, tree)
		}
	}
	return results
}

// createResourcePolices build resource policies
func (l *ResourceLogic) createResourcePolices(resourceId int64, resourceType constants.UserType, roles []*model.Role) (polices [][]string) {
	polices = slice.FilterMap(roles, func(index int, item *model.Role) ([]string, bool) {
		if item == nil || item.Status != constants.CommonStatusActive {
			return nil, false
		}
		return auth.BuildPolicy(item.Id, resourceId, resourceType), true
	})

	return
}

// SyncResource sync resource
func (l *ResourceLogic) SyncResource(ctx context.Context, input *dto.ResourceSyncRequestDto) (err error) {
	resourceRepo := l.repos.ResourceRepo

	err = resourceRepo.Transaction(ctx, func(tx *gorm.DB) error {
		for _, item := range input.Resources {
			if innerErr := l.handleSyncResource(ctx, tx, input.Type, item, 0); innerErr != nil {
				return innerErr
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	threading.GoSafe(func() {
		subCtx := context.WithoutCancel(ctx)
		msg := auth.Msg{
			Action: auth.ActionResourceReload,
			Scope:  auth.ScopeGlobal,
			Type:   input.Type,
		}
		_ = l.publishAuthEvent(subCtx, l.authTopicMap[input.Type], msg)
	})

	return nil
}

func (l *ResourceLogic) handleSyncResource(ctx context.Context, tx *gorm.DB, userType constants.UserType, item dto.ResourceSyncItemDto, parentId int64) (err error) {
	resourceRepo := l.repos.ResourceRepo
	apiRepo := l.repos.ApiRepo

	maxSort, err := resourceRepo.GetMaxSortOrderWithTx(ctx, tx, parentId, userType)
	if err != nil {
		return err
	}
	sortOrder := maxSort + 1

	apiModels := make([]*model.Api, 0, len(item.Apis))
	for _, apiItem := range item.Apis {
		if len(apiItem.Path) == 0 || len(apiItem.Method) == 0 {
			continue
		}
		api, _ := apiRepo.FindOneByColumnsWithTx(ctx, tx, map[string]interface{}{
			"path":   apiItem.Path,
			"method": apiItem.Method,
			"type":   userType,
		})
		if api != nil && api.Id != 0 {
			apiModels = append(apiModels, api)
		}
	}

	existing, _ := resourceRepo.FindOneByColumnsWithTx(ctx, tx, map[string]interface{}{
		"type":          userType,
		"resource_type": item.ResourceType,
		"tag":           item.Tag,
		// "path":          item.Path,
		// "name":          item.Name,
	})

	var resourceId int64
	if existing == nil || existing.Id == 0 {
		resourceModel := &model.Resource{
			ParentId:        parentId,
			Type:            userType,
			ResourceType:    item.ResourceType,
			Name:            item.Name,
			Path:            item.Path,
			Icon:            item.Icon,
			Tag:             item.Tag,
			IsShow:          item.IsShow,
			Status:          item.Status,
			SortOrder:       sortOrder,
			EnableDataScope: item.EnableDataScope,
			Description:     item.Description,
			IsPublic:        item.IsPublic,
			Apis:            apiModels,
		}
		created, innerErr := resourceRepo.CreateWithTx(ctx, tx, resourceModel, true)
		if innerErr != nil {
			return errorx.Wrap(innerErr, "failed to create resource")
		}
		resourceId = created.Id
	} else {
		resourceId = existing.Id
	}

	for _, child := range item.Children {
		if err = l.handleSyncResource(ctx, tx, userType, child, resourceId); err != nil {
			return err
		}
	}

	return nil
}

func (l *ResourceLogic) filterByPricingPlan(ctx context.Context, organizationId int64, resources []*model.Resource) ([]*model.Resource, error) {
	activeSubs, err := l.repos.OrganizationSubscriptionRepo.FindActiveByOrgId(ctx, organizationId)
	if err != nil {
		return nil, err
	}
	if len(activeSubs) == 0 {
		return []*model.Resource{}, nil
	}

	planResourceIds := make(map[int64]bool)
	for _, sub := range activeSubs {
		ids, innerErr := l.repos.PricingPlanAndResourceRepo.FindResourceIdsByPlanId(ctx, sub.PricingPlanId)
		if innerErr != nil {
			return nil, innerErr
		}
		for _, id := range ids {
			planResourceIds[id] = true
		}
	}

	return slice.Filter(resources, func(index int, item *model.Resource) bool {
		return planResourceIds[item.Id]
	}), nil
}
