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

func (l *SyncDataLogic) SyncData(in *admin.EmptyRequest) (resp *admin.EmptyResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
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

	return &admin.EmptyResponse{}, nil
}

// ========== sync api start ==========
type ApiGroup struct {
	Group string     `json:"group"`
	Apis  []*ApiItem `json:"apis"`
}

type ApiItem struct {
	Name               string `json:"name"`
	Path               string `json:"path"`
	Method             string `json:"method"`
	Description        string `json:"description"`
	AuthRequired       bool   `json:"auth_required"`
	PermissionRequired bool   `json:"permission_required"`
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
	var addButtonName string = "Add"
	var editButtonName string = "Edit"
	var deleteButtonName string = "Delete"

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
