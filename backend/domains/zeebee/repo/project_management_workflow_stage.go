package repo

import (
	"context"

	"gorm.io/gorm"

	"zeebee-backend/common/errorx"
	"zeebee-backend/common/gormx"
	"zeebee-backend/domains/zeebee/model"
)

type ProjectManagementWorkflowStageRepo struct {
	*gormx.BaseRepo[*model.ProjectManagementWorkflowStage]
}

func NewProjectManagementWorkflowStageRepo(db *gorm.DB) *ProjectManagementWorkflowStageRepo {
	return &ProjectManagementWorkflowStageRepo{BaseRepo: gormx.NewBaseRepo[*model.ProjectManagementWorkflowStage](db)}
}

func (r *ProjectManagementWorkflowStageRepo) MaxSortOrderWithTx(ctx context.Context, tx *gorm.DB, workflowId int64) (int, error) {
	var maxSortOrder int
	err := tx.WithContext(ctx).Model(&model.ProjectManagementWorkflowStage{}).
		Where("workflow_id = ?", workflowId).
		Select("COALESCE(MAX(sort_order), 0)").
		Scan(&maxSortOrder).Error
	return maxSortOrder, errorx.WithStack(err)
}
