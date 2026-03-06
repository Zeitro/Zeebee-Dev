package logic

import (
	"context"
	"fmt"

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
	results, err := l.repos.ResourceRepo.All(ctx, resourceType)
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
		// check existed
		isExisted, innerErr := l.repos.ResourceRepo.IsExistedWithTx(ctx, tx, resourceModel)
		if innerErr != nil {
			return innerErr
		}
		if isExisted {
			return errorx.New("Resource already exists")
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
		_ = l.BaseLogic.publishAuthEvent(subCtx, l.authTopicMap[input.Type], msg)
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
			return errorx.New("Resource already exists")
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

		_ = l.BaseLogic.publishAuthEvent(subCtx, l.authTopicMap[input.Type], msg)
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

		_ = l.BaseLogic.publishAuthEvent(subCtx, l.authTopicMap[resourceType], msg)
	})

	return nil
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

// createResourcePolices build resource policies
func (l *ResourceLogic) createResourcePolices(resourceId int64, resourceType constants.UserType, roles []*model.Role) (polices [][]string) {
	polices = slice.FilterMap(roles, func(index int, item *model.Role) ([]string, bool) {

		if item == nil || item.Status != constants.CommonStatusActive {
			return nil, false
		}
		return []string{
			fmt.Sprintf("Role:%d", item.Id),
			string(resourceType),
			fmt.Sprintf("Resource:%d", resourceId),
			"allow",
		}, true
	})

	return
}
