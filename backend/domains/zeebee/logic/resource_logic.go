package logic

import (
	"context"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"

	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/repo"
)

type ResourceLogic struct {
	repos *repo.AllRepos
}

func NewResourceLogic(repos *repo.AllRepos) *ResourceLogic {
	return &ResourceLogic{
		repos: repos,
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
			apis, innerErr := l.repos.ApiRepo.GetApisWithTx(ctx, tx, input.ApiIds, input.Type)
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

	result, err = copierx.Trans[model.Resource, dto.ResourceEntityDto](resourceModel)
	return result, err
}

// UpdateResource updates an existing resource
func (l *ResourceLogic) UpdateResource(ctx context.Context, input *dto.ResourceUpdateRequestDto) (result *dto.ResourceEntityDto, err error) {
	var resourceModel *model.Resource

	err = l.repos.ResourceRepo.Transaction(ctx, func(tx *gorm.DB) error {
		var innerErr error

		resourceModel, innerErr = l.repos.ResourceRepo.GetResourceWithTx(ctx, tx, input.Id, input.Type, true)
		if innerErr != nil {
			return innerErr
		}

		innerErr = copier.CopyWithOption(resourceModel, input.ResourceCreateRequestDto, copier.Option{DeepCopy: true})
		if innerErr != nil {
			return errorx.GetFailedTransReqErr(innerErr)
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

		// TODO: update user permissions

		return nil
	})

	if err != nil {
		return nil, err
	}

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
	err = l.repos.ResourceRepo.Transaction(ctx, func(tx *gorm.DB) error {
		var innerErr error
		innerErr = l.deleteResource(ctx, tx, id, resourceType)
		return innerErr
	})
	if err != nil {
		return err
	}
	// TODO update user permissions
	return nil
}

// deleteResource
func (l *ResourceLogic) deleteResource(ctx context.Context, tx *gorm.DB, id int64, resourceType constants.UserType) (err error) {
	resource, err := l.repos.ResourceRepo.GetResourceWithTx(ctx, tx, id, resourceType, true)
	if err != nil {
		return err
	}
	children, err := l.repos.ResourceRepo.GetResourceByParentIdWithTx(ctx, tx, resource.Id, resourceType)
	if err != nil {
		return err
	}
	for _, child := range children {
		err = l.deleteResource(ctx, tx, child.Id, resourceType)
		if err != nil {
			return err
		}
	}
	err = l.repos.ResourceRepo.DeleteResourceWithTx(ctx, tx, resource)
	return err
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
