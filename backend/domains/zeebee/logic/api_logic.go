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
