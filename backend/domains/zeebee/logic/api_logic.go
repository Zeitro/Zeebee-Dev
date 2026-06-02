package logic

import (
	"context"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/threading"
	"gorm.io/gorm"

	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/timex"
	"zeebee-backend/domains/zeebee/auth"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
)

// ApiLogic provides API business logic
type ApiLogic struct {
	*BaseLogic
	authTopicMap map[constants.UserType]string
}

// NewApiLogic creates a new API logic instance
func NewApiLogic(baseLogic *BaseLogic, authTopicMap map[constants.UserType]string) *ApiLogic {
	return &ApiLogic{
		BaseLogic:    baseLogic,
		authTopicMap: authTopicMap,
	}
}

func (l *ApiLogic) ApiTree(ctx context.Context, apiType constants.UserType) (trees []*dto.ApiEntityDto, err error) {
	results, err := l.repos.ApiRepo.All(ctx, apiType)
	if err != nil {
		return
	}
	trees = l.buildTree(results, 0)
	return
}

// CreateApi creates a new api
func (l *ApiLogic) CreateApi(ctx context.Context, input *dto.ApiCreateRequestDto) (result *dto.ApiEntityDto, err error) {
	apiModel, err := copierx.Trans[dto.ApiCreateRequestDto, model.Api](input)
	if err != nil {
		return nil, err
	}
	// check existed
	isExisted, err := l.repos.ApiRepo.IsExisted(ctx, apiModel)
	if err != nil {
		return nil, err
	}
	if isExisted {
		err = errorx.New("API already exists")
		return nil, err
	}
	apiModel, err = l.repos.ApiRepo.Create(ctx, apiModel)
	if err != nil {
		return nil, err
	}
	result, err = copierx.Trans[model.Api, dto.ApiEntityDto](apiModel)
	if err != nil {
		return nil, err
	}

	threading.GoSafe(func() {
		subCtx := context.WithoutCancel(ctx)
		msg := auth.Msg{
			Action: auth.ActionApiReload,
			Scope:  auth.ScopeGlobal,
			Type:   input.Type,
		}
		_ = l.publishAuthEvent(subCtx, l.authTopicMap[input.Type], msg)
	})

	return result, nil
}

// UpdateApi updates an existing api
func (l *ApiLogic) UpdateApi(ctx context.Context, input *dto.ApiUpdateRequestDto) (result *dto.ApiEntityDto, err error) {
	apiModel, err := l.repos.ApiRepo.GetApi(ctx, input.Id, input.Type, false)
	if err != nil {
		return nil, err
	}

	err = copier.CopyWithOption(apiModel, input, copier.Option{DeepCopy: true})
	if err != nil {
		return nil, errorx.GetFailedTransReqErr(err)
	}

	// check existed
	isExisted, err := l.repos.ApiRepo.IsExisted(ctx, apiModel)
	if err != nil {
		return nil, err
	}
	if isExisted {
		err = errorx.New("API already exists")
		return nil, err
	}
	apiModel, err = l.repos.ApiRepo.UpdateOneById(ctx, input.Id, apiModel)
	if err != nil {
		return nil, err
	}
	result, err = copierx.Trans[model.Api, dto.ApiEntityDto](apiModel)
	// publish auth event
	threading.GoSafe(func() {
		subCtx := context.WithoutCancel(ctx)
		msg := auth.Msg{
			Action: auth.ActionApiReload,
			Scope:  auth.ScopeGlobal,
			Type:   input.Type,
		}
		_ = l.publishAuthEvent(subCtx, l.authTopicMap[input.Type], msg)
	})
	return result, err
}

// Detail returns api detail
func (l *ApiLogic) Detail(ctx context.Context, id int64, apiType constants.UserType) (result *dto.ApiEntityDto, err error) {
	api, err := l.repos.ApiRepo.GetApi(ctx, id, apiType, true)
	if err != nil {
		return nil, err
	}
	result, err = copierx.Trans[model.Api, dto.ApiEntityDto](api)
	return result, err
}

// DeleteApi deletes an existing api
func (l *ApiLogic) DeleteApi(ctx context.Context, id int64, apiType constants.UserType) (err error) {
	err = l.repos.ApiRepo.Transaction(ctx, func(tx *gorm.DB) error {
		return l.deleteApi(ctx, tx, id, apiType)
	})
	if err != nil {
		return err
	}
	// publish auth event
	threading.GoSafe(func() {
		subCtx := context.WithoutCancel(ctx)
		msg := auth.Msg{
			Action: auth.ActionApiReload,
			Scope:  auth.ScopeGlobal,
			Type:   apiType,
		}
		_ = l.publishAuthEvent(subCtx, l.authTopicMap[apiType], msg)
	})
	return
}

// deleteApi deletes api and its children recursively
func (l *ApiLogic) deleteApi(ctx context.Context, tx *gorm.DB, id int64, apiType constants.UserType) (err error) {
	api, err := l.repos.ApiRepo.GetApiWithTx(ctx, tx, id, apiType, true)
	if err != nil {
		return err
	}
	// Get children and delete them
	children, err := l.repos.ApiRepo.GetApisByParentId(ctx, api.Id, apiType)
	if err != nil {
		return err
	}
	for _, child := range children {
		err = l.deleteApi(ctx, tx, child.Id, apiType)
		if err != nil {
			return err
		}
	}
	// delete resource and api relation
	err = l.repos.ApiRepo.DeleteApiWithTx(ctx, tx, api)
	return err
}

// buildTree builds a tree structure from a list of apis
func (l *ApiLogic) buildTree(apis []*model.Api, parentId int64) (results []*dto.ApiEntityDto) {
	results = slice.FilterMap[*model.Api, *dto.ApiEntityDto](apis, func(index int, item *model.Api) (*dto.ApiEntityDto, bool) {
		if item.ParentId != parentId {
			return nil, false
		}
		tree, err := copierx.Trans[model.Api, dto.ApiEntityDto](item, copierx.SqlTimePtrToStr(timex.LayoutUsStyle), copierx.TimeToRFC3339Str())
		if err != nil {
			return nil, false
		}
		return tree, true
	})

	for i, result := range results {
		results[i].Children = l.buildTree(apis, result.Id)
	}
	return results
}

// SyncApi syncs API definitions from parsed .api files
func (l *ApiLogic) SyncApi(ctx context.Context, input *dto.ApiSyncRequestDto) (err error) {
	apiRepo := l.repos.ApiRepo

	err = apiRepo.Transaction(ctx, func(tx *gorm.DB) error {
		for groupName, apis := range input.Apis {
			// 1. Get or create parent category
			parent, innerErr := apiRepo.FindApiCategoryByNameAndTypeWithTx(ctx, tx, groupName, input.Type)
			if innerErr != nil {
				return innerErr
			}
			if parent == nil {
				maxSort, innerErr := apiRepo.GetMaxSortOrderWithTx(ctx, tx, 0, input.Type)
				if innerErr != nil {
					return innerErr
				}
				parent = &model.Api{
					Name:               groupName,
					Path:               "",
					Method:             "",
					Description:        groupName,
					Type:               input.Type,
					ParentId:           0,
					Status:             constants.CommonStatusActive,
					SortOrder:          maxSort + 1,
					AuthRequired:       false,
					PermissionRequired: false,
				}
				innerErr = tx.WithContext(ctx).Create(parent).Error
				if innerErr != nil {
					return innerErr
				}
			}

			// 2. Get or create child APIs
			var apiMaxSort int32
			apiMaxSort, innerErr = apiRepo.GetMaxSortOrderWithTx(ctx, tx, parent.Id, input.Type)
			if innerErr != nil {
				return innerErr
			}
			for _, api := range apis {
				if len(api.Path) == 0 || len(api.Method) == 0 {
					continue
				}
				existing, innerErr := apiRepo.FindApiByPathMethodAndTypeWithTx(ctx, tx, api.Path, api.Method, input.Type)
				if innerErr != nil {
					return innerErr
				}
				if existing == nil {
					apiMaxSort++
					apiMethod := constants.ApiMethod(api.Method)
					if !apiMethod.IsValid() {
						innerErr = errorx.Errorf("api[%s-%s] invalid method: %v", api.Name, api.Path, api.Method)
					}
					child := &model.Api{
						Name:               api.Name,
						Path:               api.Path,
						Method:             apiMethod,
						Description:        api.Description,
						Type:               input.Type,
						ParentId:           parent.Id,
						Status:             constants.CommonStatusActive,
						SortOrder:          apiMaxSort,
						AuthRequired:       true,
						PermissionRequired: true,
					}
					innerErr = tx.WithContext(ctx).Create(child).Error
					if innerErr != nil {
						return innerErr
					}
				}
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	// publish auth reload event
	threading.GoSafe(func() {
		subCtx := context.WithoutCancel(ctx)
		msg := auth.Msg{
			Action: auth.ActionApiReload,
			Scope:  auth.ScopeGlobal,
			Type:   input.Type,
		}
		_ = l.publishAuthEvent(subCtx, l.authTopicMap[input.Type], msg)
	})

	return nil
}
