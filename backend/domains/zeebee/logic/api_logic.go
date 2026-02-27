package logic

import (
	"context"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/jinzhu/copier"

	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/timex"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/repo"
)

type ApiLogic struct {
	repos *repo.AllRepos
}

func NewApiLogic(repos *repo.AllRepos) *ApiLogic {
	return &ApiLogic{
		repos: repos,
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

// DeleteApi deletes api and its children recursively
func (l *ApiLogic) DeleteApi(ctx context.Context, id int64, apiType constants.UserType) (err error) {
	api, err := l.repos.ApiRepo.GetApi(ctx, id, apiType, true)
	if err != nil {
		return err
	}
	// Get children and delete them
	children, err := l.repos.ApiRepo.GetApisByParentId(ctx, api.Id, apiType)
	if err != nil {
		return err
	}
	for _, child := range children {
		err := l.DeleteApi(ctx, child.Id, apiType)
		if err != nil {
			return err
		}
	}
	// delete resource and api relation
	err = l.repos.ApiRepo.DeleteApi(ctx, api)
	return
}

// buildTree builds a tree structure from a list of apis
func (l *ApiLogic) buildTree(apis []*model.Api, parentId int64) (results []*dto.ApiEntityDto) {
	results = slice.FilterMap[*model.Api, *dto.ApiEntityDto](apis, func(index int, item *model.Api) (*dto.ApiEntityDto, bool) {
		if item.ParentId != parentId {
			return nil, false
		}
		tree := &dto.ApiEntityDto{}
		err := copier.CopyWithOption(tree, item, copier.Option{
			CaseSensitive: true,
			Converters: []copier.TypeConverter{
				copierx.SqlTimePtrToStr(timex.LayoutUsStyle),
				copierx.TimeToRFC3339Str(),
			},
		})
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
