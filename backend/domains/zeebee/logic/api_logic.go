package logic

import (
	"context"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/jinzhu/copier"

	"zeebee-backend/common/copierx"
	"zeebee-backend/common/timex"
	"zeebee-backend/domains/zeebee/constants"
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

type ApiTree struct {
	Id          int64      `json:"id"`
	ParentId    int64      `json:"parent_id"`
	Name        string     `json:"name"`
	Path        string     `json:"path"`
	Method      string     `json:"method"`
	Status      string     `json:"status"`
	Description string     `json:"description"`
	Type        string     `json:"type"`
	SortOrder   int32      `json:"sort_order"`
	CreatedAt   string     `json:"created_at"`
	UpdatedAt   string     `json:"updated_at"`
	Children    []*ApiTree `json:"children,omitempty"`
}

func (l *ApiLogic) ApiTree(ctx context.Context, apiType constants.UserType) (trees []*ApiTree, err error) {
	results, err := l.repos.ApiRepo.All(ctx, apiType)
	if err != nil {
		return
	}
	trees = l.buildTree(results, 0)
	return
}

// DeleteApi delete api and its children recursively
func (l *ApiLogic) DeleteApi(ctx context.Context, id int64, apiType constants.UserType) (err error) {
	api, err := l.repos.ApiRepo.GetApi(ctx, id, apiType)
	if err != nil {
		return err
	}
	// get children and delete them
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
	// todo delete resource and api relation

	_, err = l.repos.ApiRepo.DeleteByID(ctx, api.Id)
	return
}

func (l *ApiLogic) buildTree(apis []*model.Api, parentId int64) (results []*ApiTree) {
	results = slice.FilterMap[*model.Api, *ApiTree](apis, func(index int, item *model.Api) (*ApiTree, bool) {
		if item.ParentId != parentId {
			return nil, false
		}
		tree := &ApiTree{}
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
