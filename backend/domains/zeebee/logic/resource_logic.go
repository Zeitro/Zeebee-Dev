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

type ResourceLogic struct {
	repos *repo.AllRepos
}

func NewResourceLogic(repos *repo.AllRepos) *ResourceLogic {
	return &ResourceLogic{
		repos: repos,
	}
}

type ResourceTree struct {
	Id           int64           `json:"id"`
	ParentId     int64           `json:"parent_id"`
	Name         string          `json:"name"`
	Path         string          `json:"path"`
	ResourceType string          `json:"resource_type"`
	Status       string          `json:"status"`
	Description  string          `json:"description"`
	Type         string          `json:"type"`
	Icon         string          `json:"icon"`
	Tag          string          `json:"tag"`
	SortOrder    int32           `json:"sort_order"`
	IsShow       bool            `json:"is_show"`
	CreatedAt    string          `json:"created_at"`
	UpdatedAt    string          `json:"updated_at"`
	Children     []*ResourceTree `json:"children,omitempty"`
}

type ResourceWithPivot struct {
	Id                 int64               `json:"id"`
	ParentId           int64               `json:"parent_id"`
	Name               string              `json:"name"`
	Path               string              `json:"path"`
	ResourceType       string              `json:"resource_type"`
	Status             string              `json:"status"`
	Description        string              `json:"description"`
	Type               string              `json:"type"`
	Icon               string              `json:"icon"`
	Tag                string              `json:"tag"`
	SortOrder          int32               `json:"sort_order"`
	IsShow             bool                `json:"is_show"`
	PivotRoleDataScope constants.DataScope `json:"pivot_role_data_scope"`
	CreatedAt          string              `json:"created_at"`
	UpdatedAt          string              `json:"updated_at"`
}

func (l *ResourceLogic) ResourceTree(ctx context.Context, resourceType constants.UserType) (trees []*ResourceTree, err error) {
	results, err := l.repos.ResourceRepo.All(ctx, resourceType)
	if err != nil {
		return
	}
	trees = l.buildTree(results, 0)
	return
}

// DeleteResource deletes resource and its children recursively
func (l *ResourceLogic) DeleteResource(ctx context.Context, id int64, resourceType constants.UserType) (err error) {
	resource, err := l.repos.ResourceRepo.GetResource(ctx, id, resourceType)
	if err != nil {
		return
	}
	// Get children and delete them
	children, err := l.repos.ResourceRepo.GetResourceByParentId(ctx, resource.Id, resourceType)
	if err != nil {
		return err
	}
	for _, child := range children {
		err := l.DeleteResource(ctx, child.Id, resourceType)
		if err != nil {
			return err
		}
	}
	// TODO delete resource and api relation

	// TODO delete resource and role relation

	_, err = l.repos.ApiRepo.DeleteByID(ctx, resource.Id)
	return
}

func (l *ResourceLogic) buildTree(resources []*model.Resource, parentId int64) (results []*ResourceTree) {
	results = slice.FilterMap[*model.Resource, *ResourceTree](resources, func(index int, item *model.Resource) (*ResourceTree, bool) {
		if item.ParentId != parentId {
			return nil, false
		}
		tree := &ResourceTree{}
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
		results[i].Children = l.buildTree(resources, result.Id)
	}
	return results
}
