package commonlogic

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/duke-git/lancet/v2/slice"
	"gorm.io/gorm"

	"zeebee-backend/app/admin/rpc/internal/logic"
	"zeebee-backend/app/admin/rpc/internal/static"
	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/common/errorx"
)

type SyncDataLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	baseLogic *logic.BaseLogic
}

func NewSyncDataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SyncDataLogic {
	return &SyncDataLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *SyncDataLogic) SyncData(_ *admin.EmptyRequest) (resp *admin.EmptyResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(
				logx.Field("stacks", errorx.GetErrorStack(err)),
				logx.Field("service_name", l.svcCtx.Config.Name),
			).Error(err)
		}
	}()
	err = l.syncApi()
	if err != nil {
		return nil, err
	}

	// err = l.createMenuButton()
	// if err != nil {
	// 	return nil, err
	// }
	err = l.syncResource()
	if err != nil {
		return nil, err
	}

	return &admin.EmptyResponse{}, nil
}

// ========== sync api start ==========
type ApiGroup struct {
	Group string     `json:"group"`
	Apis  []*ApiItem `json:"apis"`
}

type ApiItem struct {
	Name               string              `json:"name"`
	Path               string              `json:"path"`
	Method             constants.ApiMethod `json:"method"`
	Description        string              `json:"description"`
	AuthRequired       bool                `json:"auth_required"`
	PermissionRequired bool                `json:"permission_required"`
}

func (l *SyncDataLogic) syncApi() (err error) {
	apiMap := map[constants.UserType][]ApiGroup{}
	apiContent := static.ApiStatic
	err = json.Unmarshal(apiContent, &apiMap)
	if err != nil {
		return errorx.Wrap(err, "Failed to unmarshal api content")
	}
	apiRepo := l.svcCtx.ApiRepo
	for apiType, apiGroup := range apiMap {
		err = apiRepo.Transaction(l.ctx, func(tx *gorm.DB) error {
			var innerError error
			var parentApiSort int32 = 1
			for _, group := range apiGroup {
				// get or create parent api
				parentApi := &model.Api{
					Name:        group.Group,
					Path:        "",
					Method:      "",
					Description: "",
					Type:        apiType,
					ParentId:    0,
					SortOrder:   parentApiSort,
				}
				parentApi, innerError = apiRepo.GetOrCreateApiWithTx(l.ctx, tx, parentApi)
				if innerError != nil {
					return innerError
				}
				// get or create child api
				childApiSort := int32(1)
				for _, api := range group.Apis {
					childApi := &model.Api{
						Name:               api.Name,
						Path:               api.Path,
						Method:             api.Method,
						Description:        api.Description,
						Type:               apiType,
						ParentId:           parentApi.Id,
						AuthRequired:       api.AuthRequired,
						PermissionRequired: api.PermissionRequired,
						SortOrder:          childApiSort,
					}
					_, innerError = apiRepo.GetOrCreateApiWithTx(l.ctx, tx, childApi)
					if innerError != nil {
						return innerError
					}
					childApiSort++
				}
				parentApiSort++
			}

			return innerError
		})
		if err != nil {
			return err
		}
	}

	return err
}

// ========== sync api end ==========

// ========== create button for menu ==========
func (l *SyncDataLogic) createMenuButton() (err error) {
	types := []constants.UserType{
		constants.UserTypeLender,
		constants.UserTypeBorrower,
	}
	ctx := l.ctx
	for _, userType := range types {
		var resources []*dto.ResourceTreeDto
		resources, err = l.baseLogic.ResourceLogic.ResourceTree(ctx, userType)
		if err != nil {
			return err
		}
		for _, resource := range resources {
			err = l.handleMenu(l.ctx, resource)
			if err != nil {
				return err
			}
		}
	}
	return err
}

func (l *SyncDataLogic) handleMenu(ctx context.Context, resource *dto.ResourceTreeDto) (err error) {
	if resource == nil || resource.ResourceType != constants.ResourceTypeMenu || !resource.IsShow {
		return
	}
	hasSubMenu := slice.Some(resource.Children, func(index int, item *dto.ResourceTreeDto) bool {
		return item.ResourceType == constants.ResourceTypeMenu
	})
	// handle sub menu
	if hasSubMenu {
		for _, child := range resource.Children {
			if child.ResourceType == constants.ResourceTypeMenu {
				err = l.handleMenu(ctx, child)
				if err != nil {
					return err
				}
			}
		}
		return
	}
	var addButtonName = "Add"
	var editButtonName = "Edit"
	var deleteButtonName = "Delete"

	// handle menu
	hasCreateButton := false
	hasEditButton := false
	hasDeleteButton := false
	slice.ForEach(resource.Children, func(index int, item *dto.ResourceTreeDto) {
		if item.ResourceType != constants.ResourceTypeButton {
			return
		}
		if item.Name == addButtonName {
			hasCreateButton = true
		}
		if item.Name == editButtonName {
			hasEditButton = true
		}
		if item.Name == deleteButtonName {
			hasDeleteButton = true
		}
	})
	var input *dto.ResourceCreateRequestDto
	if !hasCreateButton {
		input = &dto.ResourceCreateRequestDto{
			ParentId:        resource.Id,
			Type:            resource.Type,
			ResourceType:    constants.ResourceTypeButton,
			Name:            addButtonName,
			Tag:             fmt.Sprintf("%s:%s", resource.Tag, addButtonName),
			IsShow:          true,
			Status:          constants.CommonStatusActive,
			SortOrder:       1,
			ApiIds:          []int64{},
			EnableDataScope: false,
		}
		_, err = l.baseLogic.ResourceLogic.CreateResource(ctx, input)
		if err != nil {
			return err
		}
	}
	if !hasEditButton {
		input = &dto.ResourceCreateRequestDto{
			ParentId:        resource.Id,
			Type:            resource.Type,
			ResourceType:    constants.ResourceTypeButton,
			Name:            editButtonName,
			Tag:             fmt.Sprintf("%s:%s", resource.Tag, editButtonName),
			IsShow:          true,
			Status:          constants.CommonStatusActive,
			SortOrder:       2,
			ApiIds:          []int64{},
			EnableDataScope: false,
		}
		_, err = l.baseLogic.ResourceLogic.CreateResource(ctx, input)
		if err != nil {
			return err
		}
	}
	if !hasDeleteButton {
		input = &dto.ResourceCreateRequestDto{
			ParentId:        resource.Id,
			Type:            resource.Type,
			ResourceType:    constants.ResourceTypeButton,
			Name:            deleteButtonName,
			Tag:             fmt.Sprintf("%s:%s", resource.Tag, deleteButtonName),
			IsShow:          true,
			Status:          constants.CommonStatusActive,
			SortOrder:       3,
			ApiIds:          []int64{},
			EnableDataScope: false,
		}
		_, err = l.baseLogic.ResourceLogic.CreateResource(ctx, input)
		if err != nil {
			return err
		}
	}
	return
}

// ========== create button for menu end ==========

// =========== sync resource start ==========
type Resource struct {
	// Id              int64                  `json:"id"`
	// ParentId        int64                  `json:"parent_id"`
	Type            constants.UserType     `json:"type"`
	ResourceType    constants.ResourceType `json:"resource_type"`
	Name            string                 `json:"name"`
	Path            string                 `json:"path"`
	Icon            string                 `json:"icon"`
	Tag             string                 `json:"tag"`
	IsShow          bool                   `json:"is_show"`
	Status          constants.CommonStatus `json:"status"`
	SortOrder       int32                  `json:"sort_order"`
	EnableDataScope bool                   `json:"enable_data_scope"`
	Description     string                 `json:"description"`
	IsPublic        bool                   `json:"is_public"`
	Children        []*Resource            `json:"children"`
	Apis            []*ResourceApi         `json:"apis"`
}

type ResourceApi struct {
	// Id                 int64                  `json:"id"`
	// ParentId           int64                  `json:"parent_id"`
	Name               string                 `json:"name"`
	Path               string                 `json:"path"`
	Method             constants.ApiMethod    `json:"method"`
	Status             constants.CommonStatus `json:"status"`
	Description        string                 `json:"description"`
	Type               constants.UserType     `json:"type"`
	SortOrder          int32                  `json:"sort_order"`
	AuthRequired       bool                   `json:"auth_required"`
	PermissionRequired bool                   `json:"permission_required"`
}

func (l *SyncDataLogic) syncResource() (err error) {
	resourceMap := map[constants.UserType][]*Resource{}
	err = json.Unmarshal(static.ResourceStatic, &resourceMap)
	if err != nil {
		return errorx.Wrap(err, "Failed to unmarshal resource content")
	}
	for _, resources := range resourceMap {
		for _, resource := range resources {
			err = l.handleResource(l.ctx, resource, 0)
			if err != nil {
				return err
			}
		}
	}

	return err
}

func (l *SyncDataLogic) handleResource(ctx context.Context, resource *Resource, parentId int64) (err error) {
	apiModels := slice.FilterMap(resource.Apis, func(index int, resourceApi *ResourceApi) (*model.Api, bool) {
		if len(resourceApi.Path) == 0 || len(resourceApi.Method) == 0 {
			return nil, false
		}
		api, _ := l.svcCtx.ApiRepo.FindOneByColumns(ctx, map[string]interface{}{
			"path":   resourceApi.Path,
			"method": resourceApi.Method,
			"type":   resource.Type,
		})
		if api == nil || api.Id == 0 {
			return nil, false
		}
		return api, true
	})
	resourceModel, _ := l.svcCtx.ResourceRepo.FindOneByColumns(ctx, map[string]interface{}{
		"type":          resource.Type,
		"resource_type": resource.ResourceType,
		"tag":           resource.Tag,
		"path":          resource.Path,
		"name":          resource.Name,
	})
	if resourceModel == nil || resourceModel.Id == 0 {
		//  create new one
		input := &model.Resource{
			ParentId:        parentId,
			Type:            resource.Type,
			ResourceType:    resource.ResourceType,
			Name:            resource.Name,
			Path:            resource.Path,
			Icon:            resource.Icon,
			Tag:             resource.Tag,
			IsShow:          resource.IsShow,
			Status:          resource.Status,
			SortOrder:       resource.SortOrder,
			EnableDataScope: resource.EnableDataScope,
			Description:     resource.Description,
			IsPublic:        resource.IsPublic,
			Apis:            apiModels,
		}
		resourceModel, err = l.svcCtx.ResourceRepo.Create(ctx, input)
		if err != nil {
			return errorx.Wrap(err, "Failed to create resource")
		}
	}
	for _, child := range resource.Children {
		err = l.handleResource(ctx, child, resourceModel.Id)
		if err != nil {
			return err
		}
	}

	return err
}

// =========== sync resource end ==========
