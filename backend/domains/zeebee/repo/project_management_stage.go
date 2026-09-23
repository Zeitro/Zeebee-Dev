package repo

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"zeebee-backend/common/errorx"
	"zeebee-backend/common/gormx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/model"
)

type ProjectManagementStageRepo struct {
	*gormx.BaseRepo[*model.ProjectManagementStage]
}

func NewProjectManagementStageRepo(db *gorm.DB) *ProjectManagementStageRepo {
	return &ProjectManagementStageRepo{BaseRepo: gormx.NewBaseRepo[*model.ProjectManagementStage](db)}
}

func (r *ProjectManagementStageRepo) ListByOwner(ctx context.Context, entityType constants.ProjectManagementOwnerType, entityId int64) ([]*model.ProjectManagementStage, error) {
	stages := make([]*model.ProjectManagementStage, 0)
	err := r.Client(ctx).Where("entity_type = ? AND entity_id = ?", entityType, entityId).Order("sort_order ASC, id ASC").Find(&stages).Error
	return stages, errorx.WithStack(err)
}

func (r *ProjectManagementStageRepo) ListByOwnerWithTx(ctx context.Context, tx *gorm.DB, entityType constants.ProjectManagementOwnerType, entityId int64) ([]*model.ProjectManagementStage, error) {
	stages := make([]*model.ProjectManagementStage, 0)
	err := tx.WithContext(ctx).Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("entity_type = ? AND entity_id = ?", entityType, entityId).Order("sort_order ASC, id ASC").Find(&stages).Error
	return stages, errorx.WithStack(err)
}

func (r *ProjectManagementStageRepo) IsUsedByTemplateWorkflowWithTx(ctx context.Context, tx *gorm.DB, templateId, stageId int64) (bool, error) {
	var count int64
	err := tx.WithContext(ctx).Model(&model.ProjectManagementWorkflowStage{}).
		Joins("JOIN project_management_workflows ON project_management_workflows.id = project_management_workflow_stages.workflow_id").
		Where("project_management_workflows.entity_type = ? AND project_management_workflows.entity_id = ? AND project_management_workflow_stages.stage_id = ?", constants.ProjectManagementOwnerTypeTemplate, templateId, stageId).
		Count(&count).Error
	return count > 0, errorx.WithStack(err)
}

func (r *ProjectManagementStageRepo) LockByOwnerWithTx(ctx context.Context, tx *gorm.DB, id int64, entityType constants.ProjectManagementOwnerType, entityId int64) (*model.ProjectManagementStage, error) {
	stage := new(model.ProjectManagementStage)
	err := tx.WithContext(ctx).Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ? AND entity_type = ? AND entity_id = ?", id, entityType, entityId).First(stage).Error
	return stage, errorx.WithStack(err)
}
