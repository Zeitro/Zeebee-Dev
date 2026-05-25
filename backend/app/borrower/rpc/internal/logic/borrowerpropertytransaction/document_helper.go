package borrowerpropertytransactionlogic

import (
	"context"

	"zeebee-backend/domains/zeebee/model"
)

// isDescendantOf checks if doc is rootDocId or a descendant of it by walking up the parent chain
func isDescendantOf(ctx context.Context, repo interface {
	FindOneById(context.Context, int64) (*model.Document, error)
}, doc *model.Document, rootDocId int64) bool {
	current := doc
	for current != nil {
		if current.Id == rootDocId {
			return true
		}
		if current.ParentId == 0 {
			break
		}
		parent, err := repo.FindOneById(ctx, current.ParentId)
		if err != nil {
			break
		}
		current = parent
	}
	return false
}
