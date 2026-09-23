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

type ProjectManagementTemplateRepo struct {
	*gormx.BaseRepo[*model.ProjectManagementTemplate]
}

func NewProjectManagementTemplateRepo(db *gorm.DB) *ProjectManagementTemplateRepo {
	return &ProjectManagementTemplateRepo{BaseRepo: gormx.NewBaseRepo[*model.ProjectManagementTemplate](db)}
}

func (r *ProjectManagementTemplateRepo) FindByOrganization(ctx context.Context, id, organizationId int64) (*model.ProjectManagementTemplate, error) {
	return r.FindOneByColumns(ctx, map[string]any{"id": id, "organization_id": organizationId})
}

func (r *ProjectManagementTemplateRepo) LockByOrganizationWithTx(ctx context.Context, tx *gorm.DB, id, organizationId int64) (*model.ProjectManagementTemplate, error) {
	template := new(model.ProjectManagementTemplate)
	err := tx.WithContext(ctx).Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ? AND organization_id = ?", id, organizationId).First(template).Error
	return template, errorx.WithStack(err)
}

func (r *ProjectManagementTemplateRepo) ListWorkflowSources(ctx context.Context, organizationId int64) ([]*model.ProjectManagementTemplate, error) {
	templates := make([]*model.ProjectManagementTemplate, 0)
	err := r.Client(ctx).Where("organization_id = ? AND template_type = ? AND (project_type IS NULL OR project_type = '')", organizationId, constants.ProjectManagementTemplateTypeOrganization).
		Order("is_default DESC, sort_order ASC, id ASC").Find(&templates).Error
	return templates, errorx.WithStack(err)
}

func (r *ProjectManagementTemplateRepo) ListProjectTemplates(ctx context.Context, organizationId int64) ([]*model.ProjectManagementTemplate, error) {
	templates := make([]*model.ProjectManagementTemplate, 0)
	err := r.Client(ctx).Model(&model.ProjectManagementTemplate{}).
		Where("organization_id = ? AND template_type = ? AND project_type <> ''", organizationId, constants.ProjectManagementTemplateTypeOrganization).
		Order("sort_order ASC, id ASC").Find(&templates).Error
	return templates, errorx.WithStack(err)
}

func (r *ProjectManagementTemplateRepo) ListPublicTemplates(ctx context.Context, status string) ([]*model.ProjectManagementTemplate, error) {
	templates := make([]*model.ProjectManagementTemplate, 0)
	err := r.Client(ctx).Model(&model.ProjectManagementTemplate{}).
		Where("template_type = ? AND organization_id = 0 AND status = ?", constants.ProjectManagementTemplateTypePublic, status).
		Order("sort_order ASC, id ASC").Find(&templates).Error
	return templates, errorx.WithStack(err)
}

func (r *ProjectManagementTemplateRepo) PagePublic(ctx context.Context, page, limit int64) (gormx.PageResponse[*model.ProjectManagementTemplate], error) {
	query := r.Client(ctx).Model(&model.ProjectManagementTemplate{}).
		Where("template_type = ? AND organization_id = 0", constants.ProjectManagementTemplateTypePublic).
		Order("id DESC")
	return r.BaseRepo.Page(query, gormx.PageRequest{Page: page, Limit: limit})
}

func (r *ProjectManagementTemplateRepo) FindPublicByID(ctx context.Context, id int64) (*model.ProjectManagementTemplate, error) {
	template := new(model.ProjectManagementTemplate)
	err := r.Client(ctx).Model(&model.ProjectManagementTemplate{}).
		Where("id = ? AND template_type = ? AND organization_id = 0", id, constants.ProjectManagementTemplateTypePublic).
		First(template).Error
	return template, errorx.WithStack(err)
}

func (r *ProjectManagementTemplateRepo) FindVisibleProjectTemplate(ctx context.Context, id, organizationId int64) (*model.ProjectManagementTemplate, error) {
	template := new(model.ProjectManagementTemplate)
	err := r.Client(ctx).Model(&model.ProjectManagementTemplate{}).
		Where("id = ? AND ((template_type = ? AND organization_id = 0) OR (template_type = ? AND organization_id = ? AND project_type <> ''))", id, constants.ProjectManagementTemplateTypePublic, constants.ProjectManagementTemplateTypeOrganization, organizationId).
		First(template).Error
	return template, errorx.WithStack(err)
}

func (r *ProjectManagementTemplateRepo) FindVisibleProjectTemplateWithTx(ctx context.Context, tx *gorm.DB, id, organizationId int64) (*model.ProjectManagementTemplate, error) {
	template := new(model.ProjectManagementTemplate)
	err := tx.WithContext(ctx).Model(&model.ProjectManagementTemplate{}).
		Where("id = ? AND ((template_type = ? AND organization_id = 0) OR (template_type = ? AND organization_id = ? AND project_type <> ''))", id, constants.ProjectManagementTemplateTypePublic, constants.ProjectManagementTemplateTypeOrganization, organizationId).
		First(template).Error
	return template, errorx.WithStack(err)
}

func (r *ProjectManagementTemplateRepo) ListWorkflows(ctx context.Context, templateID int64) ([]*model.ProjectManagementWorkflow, error) {
	workflows := make([]*model.ProjectManagementWorkflow, 0)
	err := r.Client(ctx).Model(&model.ProjectManagementWorkflow{}).
		Preload("Stages", func(db *gorm.DB) *gorm.DB { return db.Order("sort_order ASC, id ASC") }).
		Preload("Transitions", preloadWorkflowTransitions).
		Where("entity_type = ? AND entity_id = ?", constants.ProjectManagementOwnerTypeTemplate, templateID).
		Order("is_default DESC, id ASC").Find(&workflows).Error
	return workflows, errorx.WithStack(err)
}

func (r *ProjectManagementTemplateRepo) FindWorkflow(ctx context.Context, templateID, workflowID int64) (*model.ProjectManagementWorkflow, error) {
	workflow := new(model.ProjectManagementWorkflow)
	err := r.Client(ctx).Model(&model.ProjectManagementWorkflow{}).
		Preload("Stages", func(db *gorm.DB) *gorm.DB { return db.Order("sort_order ASC, id ASC") }).
		Preload("Transitions", preloadWorkflowTransitions).
		Where("id = ? AND entity_type = ? AND entity_id = ?", workflowID, constants.ProjectManagementOwnerTypeTemplate, templateID).
		First(workflow).Error
	return workflow, errorx.WithStack(err)
}

func (r *ProjectManagementTemplateRepo) ListStages(ctx context.Context, templateID int64) ([]*model.ProjectManagementStage, error) {
	stages := make([]*model.ProjectManagementStage, 0)
	err := r.Client(ctx).Where("entity_type = ? AND entity_id = ?", constants.ProjectManagementOwnerTypeTemplate, templateID).
		Order("sort_order ASC, id ASC").Find(&stages).Error
	return stages, errorx.WithStack(err)
}
