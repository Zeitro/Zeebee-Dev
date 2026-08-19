package repo

import (
	"context"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"zeebee-backend/common/errorx"
	"zeebee-backend/common/gormx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/model"
)

type ProjectManagementWorkflowRepo struct {
	*gormx.BaseRepo[*model.ProjectManagementWorkflow]
}

func NewProjectManagementWorkflowRepo(db *gorm.DB) *ProjectManagementWorkflowRepo {
	return &ProjectManagementWorkflowRepo{BaseRepo: gormx.NewBaseRepo[*model.ProjectManagementWorkflow](db)}
}

func (r *ProjectManagementWorkflowRepo) FindByOrganization(ctx context.Context, id, organizationId int64) (*model.ProjectManagementWorkflow, error) {
	workflow := new(model.ProjectManagementWorkflow)
	err := r.Client(ctx).Where("id = ? AND organization_id = ?", id, organizationId).First(workflow).Error
	return workflow, errorx.WithStack(err)
}

func (r *ProjectManagementWorkflowRepo) FindDefinitionByOrganization(ctx context.Context, id, organizationId int64) (*model.ProjectManagementWorkflow, error) {
	workflow := new(model.ProjectManagementWorkflow)
	err := r.Client(ctx).
		Preload("Stages", func(db *gorm.DB) *gorm.DB { return db.Order("sort_order ASC, id ASC") }).
		Preload("Transitions", func(db *gorm.DB) *gorm.DB { return db.Order("id ASC") }).
		Where("id = ? AND organization_id = ? AND kind = ?", id, organizationId, constants.ProjectManagementWorkflowKindOrganizationTemplate).
		First(workflow).Error
	return workflow, errorx.WithStack(err)
}

func (r *ProjectManagementWorkflowRepo) ListEnabledByOrganization(ctx context.Context, organizationId int64) ([]*model.ProjectManagementWorkflow, error) {
	workflows := make([]*model.ProjectManagementWorkflow, 0)
	err := r.Client(ctx).
		Preload("Stages", func(db *gorm.DB) *gorm.DB { return db.Order("sort_order ASC, id ASC") }).
		Preload("Transitions", func(db *gorm.DB) *gorm.DB { return db.Order("id ASC") }).
		Where("organization_id = ? AND kind = ? AND enabled = true", organizationId, constants.ProjectManagementWorkflowKindOrganizationTemplate).
		Order("is_default DESC, id ASC").
		Find(&workflows).Error
	return workflows, errorx.WithStack(err)
}

func (r *ProjectManagementWorkflowRepo) Page(ctx context.Context, organizationId int64, pageReq gormx.PageRequest, keyword string, enabled *bool) (gormx.PageResponse[*model.ProjectManagementWorkflow], error) {
	query := r.Client(ctx).
		Where("organization_id = ? AND kind = ?", organizationId, constants.ProjectManagementWorkflowKindOrganizationTemplate).
		Order("is_default DESC, id DESC")
	if strings.TrimSpace(keyword) != "" {
		query = query.Where("name ILIKE ?", "%"+keyword+"%")
	}
	if enabled != nil {
		query = query.Where("enabled = ?", *enabled)
	}
	return r.BaseRepo.Page(query, pageReq)
}

func (r *ProjectManagementWorkflowRepo) FindPendingStageWithTx(ctx context.Context, tx *gorm.DB, workflowId int64) (*model.ProjectManagementWorkflowStage, error) {
	stage := new(model.ProjectManagementWorkflowStage)
	err := tx.WithContext(ctx).Where("workflow_id = ? AND phase = ?", workflowId, constants.ProjectManagementStagePhasePending).First(stage).Error
	return stage, errorx.WithStack(err)
}

func (r *ProjectManagementWorkflowRepo) FindProjectDefinitionWithTx(ctx context.Context, tx *gorm.DB, projectId, organizationId int64) (*model.ProjectManagementWorkflow, error) {
	workflow := new(model.ProjectManagementWorkflow)
	err := tx.WithContext(ctx).
		Preload("Stages", func(db *gorm.DB) *gorm.DB { return db.Order("sort_order ASC, id ASC") }).
		Preload("Transitions", func(db *gorm.DB) *gorm.DB { return db.Order("id ASC") }).
		Where("project_id = ? AND organization_id = ? AND kind = ?", projectId, organizationId, constants.ProjectManagementWorkflowKindProject).
		First(workflow).Error
	return workflow, errorx.WithStack(err)
}

func (r *ProjectManagementWorkflowRepo) FindProjectDefinition(ctx context.Context, projectId, organizationId int64) (*model.ProjectManagementWorkflow, error) {
	workflow := new(model.ProjectManagementWorkflow)
	err := r.Client(ctx).
		Preload("Stages", func(db *gorm.DB) *gorm.DB { return db.Order("sort_order ASC, id ASC") }).
		Preload("Transitions", func(db *gorm.DB) *gorm.DB { return db.Order("id ASC") }).
		Where("project_id = ? AND organization_id = ? AND kind = ?", projectId, organizationId, constants.ProjectManagementWorkflowKindProject).
		First(workflow).Error
	return workflow, errorx.WithStack(err)
}

func (r *ProjectManagementWorkflowRepo) FindStageWithTx(ctx context.Context, tx *gorm.DB, id int64) (*model.ProjectManagementWorkflowStage, error) {
	stage := new(model.ProjectManagementWorkflowStage)
	err := tx.WithContext(ctx).Where("id = ?", id).First(stage).Error
	return stage, errorx.WithStack(err)
}

func (r *ProjectManagementWorkflowRepo) FindStage(ctx context.Context, id int64) (*model.ProjectManagementWorkflowStage, error) {
	stage := new(model.ProjectManagementWorkflowStage)
	err := r.Client(ctx).Model(&model.ProjectManagementWorkflowStage{}).Where("id = ?", id).First(stage).Error
	return stage, errorx.WithStack(err)
}

func (r *ProjectManagementWorkflowRepo) FindTransitionWithTx(ctx context.Context, tx *gorm.DB, workflowId, fromStageId, toStageId int64) (*model.ProjectManagementWorkflowTransition, error) {
	transition := new(model.ProjectManagementWorkflowTransition)
	err := tx.WithContext(ctx).Where("workflow_id = ? AND from_stage_id = ? AND to_stage_id = ?", workflowId, fromStageId, toStageId).First(transition).Error
	return transition, errorx.WithStack(err)
}

func (r *ProjectManagementWorkflowRepo) LockProjectWithTx(ctx context.Context, tx *gorm.DB, projectId, organizationId int64) (*model.ProjectManagementWorkflow, error) {
	workflow := new(model.ProjectManagementWorkflow)
	err := tx.WithContext(ctx).Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("project_id = ? AND organization_id = ? AND kind = ?", projectId, organizationId, constants.ProjectManagementWorkflowKindProject).
		First(workflow).Error
	return workflow, errorx.WithStack(err)
}

func (r *ProjectManagementWorkflowRepo) LockWithTx(ctx context.Context, tx *gorm.DB, id int64) (*model.ProjectManagementWorkflow, error) {
	workflow := new(model.ProjectManagementWorkflow)
	err := tx.WithContext(ctx).Clauses(clause.Locking{Strength: "UPDATE"}).First(workflow, id).Error
	return workflow, errorx.WithStack(err)
}

func (r *ProjectManagementWorkflowRepo) LockTemplateWithTx(ctx context.Context, tx *gorm.DB, id, organizationId int64) (*model.ProjectManagementWorkflow, error) {
	workflow := new(model.ProjectManagementWorkflow)
	err := tx.WithContext(ctx).Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ? AND organization_id = ? AND kind = ?", id, organizationId, constants.ProjectManagementWorkflowKindOrganizationTemplate).
		First(workflow).Error
	return workflow, errorx.WithStack(err)
}

func (r *ProjectManagementWorkflowRepo) DemoteDefaultWithTx(ctx context.Context, tx *gorm.DB, organizationId, exceptId int64) error {
	err := tx.WithContext(ctx).Model(&model.ProjectManagementWorkflow{}).
		Where("organization_id = ? AND kind = ? AND id <> ?", organizationId, constants.ProjectManagementWorkflowKindOrganizationTemplate, exceptId).
		Update("is_default", false).Error
	return errorx.WithStack(err)
}

func (r *ProjectManagementWorkflowRepo) UpdateTemplateWithTx(ctx context.Context, tx *gorm.DB, workflow *model.ProjectManagementWorkflow) error {
	err := tx.WithContext(ctx).Model(workflow).
		Where("id = ? AND organization_id = ? AND kind = ?", workflow.Id, workflow.OrganizationId, constants.ProjectManagementWorkflowKindOrganizationTemplate).
		Select("name", "description", "enabled", "is_default").Updates(workflow).Error
	return errorx.WithStack(err)
}

func (r *ProjectManagementWorkflowRepo) ReplaceTemplateDefinitionWithTx(ctx context.Context, tx *gorm.DB, workflowId int64, stages []*model.ProjectManagementWorkflowStage, transitions []*model.ProjectManagementWorkflowTransition) error {
	if err := tx.WithContext(ctx).Where("workflow_id = ?", workflowId).Delete(&model.ProjectManagementWorkflowTransition{}).Error; err != nil {
		return errorx.WithStack(err)
	}
	if err := tx.WithContext(ctx).Where("workflow_id = ?", workflowId).Delete(&model.ProjectManagementWorkflowStage{}).Error; err != nil {
		return errorx.WithStack(err)
	}
	stageIds := make(map[int64]int64, len(stages))
	for _, stage := range stages {
		referenceId := stage.Id
		stage.Id = 0
		stage.WorkflowId = workflowId
		if err := tx.WithContext(ctx).Create(stage).Error; err != nil {
			return errorx.WithStack(err)
		}
		stageIds[referenceId] = stage.Id
	}
	for _, transition := range transitions {
		transition.Id = 0
		transition.WorkflowId = workflowId
		transition.FromStageId = stageIds[transition.FromStageId]
		transition.ToStageId = stageIds[transition.ToStageId]
		if err := tx.WithContext(ctx).Create(transition).Error; err != nil {
			return errorx.WithStack(err)
		}
	}
	return nil
}

func (r *ProjectManagementWorkflowRepo) UpdateProjectWithTx(ctx context.Context, tx *gorm.DB, workflow *model.ProjectManagementWorkflow) error {
	err := tx.WithContext(ctx).Model(workflow).
		Where("id = ? AND organization_id = ? AND kind = ?", workflow.Id, workflow.OrganizationId, constants.ProjectManagementWorkflowKindProject).
		Select("name", "description").Updates(workflow).Error
	return errorx.WithStack(err)
}

// ReplaceProjectDefinitionWithTx preserves submitted positive stage IDs. New stages
// use unique negative IDs for transition references within this replacement request.
func (r *ProjectManagementWorkflowRepo) ReplaceProjectDefinitionWithTx(ctx context.Context, tx *gorm.DB, workflowId int64, stages []*model.ProjectManagementWorkflowStage, transitions []*model.ProjectManagementWorkflowTransition) error {
	existing := make([]*model.ProjectManagementWorkflowStage, 0)
	if err := tx.WithContext(ctx).Where("workflow_id = ?", workflowId).Find(&existing).Error; err != nil {
		return errorx.WithStack(err)
	}
	existingById := make(map[int64]*model.ProjectManagementWorkflowStage, len(existing))
	for _, stage := range existing {
		existingById[stage.Id] = stage
	}
	stageIds := make(map[int64]int64, len(stages))
	retained := make(map[int64]struct{}, len(stages))
	for _, stage := range stages {
		referenceId := stage.Id
		if referenceId > 0 {
			if _, ok := existingById[referenceId]; !ok {
				return errorx.New("workflow stage does not belong to the project workflow")
			}
			if err := tx.WithContext(ctx).Model(stage).Where("id = ?", referenceId).
				Select("name", "phase", "sort_order").Updates(stage).Error; err != nil {
				return errorx.WithStack(err)
			}
			stageIds[referenceId], retained[referenceId] = referenceId, struct{}{}
			continue
		}
		stage.Id, stage.WorkflowId = 0, workflowId
		if err := tx.WithContext(ctx).Create(stage).Error; err != nil {
			return errorx.WithStack(err)
		}
		stageIds[referenceId] = stage.Id
	}
	for _, stage := range existing {
		if _, keep := retained[stage.Id]; keep {
			continue
		}
		var taskCount int64
		if err := tx.WithContext(ctx).Model(&model.ProjectManagementTask{}).Where("stage_id = ?", stage.Id).Count(&taskCount).Error; err != nil {
			return errorx.WithStack(err)
		}
		if taskCount != 0 {
			return errorx.New("workflow stage is referenced by a task")
		}
	}
	if err := tx.WithContext(ctx).Where("workflow_id = ?", workflowId).Delete(&model.ProjectManagementWorkflowTransition{}).Error; err != nil {
		return errorx.WithStack(err)
	}
	for _, stage := range existing {
		if _, keep := retained[stage.Id]; !keep {
			if err := tx.WithContext(ctx).Delete(stage).Error; err != nil {
				return errorx.WithStack(err)
			}
		}
	}
	for _, transition := range transitions {
		fromStageId, fromOK := stageIds[transition.FromStageId]
		toStageId, toOK := stageIds[transition.ToStageId]
		if !fromOK || !toOK {
			return errorx.New("workflow transition stage is invalid")
		}
		transition.Id, transition.WorkflowId = 0, workflowId
		transition.FromStageId, transition.ToStageId = fromStageId, toStageId
		if err := tx.WithContext(ctx).Create(transition).Error; err != nil {
			return errorx.WithStack(err)
		}
	}
	return nil
}

func (r *ProjectManagementWorkflowRepo) DeleteTemplateWithTx(ctx context.Context, tx *gorm.DB, workflowId int64) error {
	if err := tx.WithContext(ctx).Where("workflow_id = ?", workflowId).Delete(&model.ProjectManagementWorkflowTransition{}).Error; err != nil {
		return errorx.WithStack(err)
	}
	if err := tx.WithContext(ctx).Where("workflow_id = ?", workflowId).Delete(&model.ProjectManagementWorkflowStage{}).Error; err != nil {
		return errorx.WithStack(err)
	}
	err := tx.WithContext(ctx).Delete(&model.ProjectManagementWorkflow{Base: model.Base{Id: workflowId}}).Error
	return errorx.WithStack(err)
}

func (r *ProjectManagementWorkflowRepo) DeleteProjectWorkflowWithTx(ctx context.Context, tx *gorm.DB, projectId int64) error {
	workflows := make([]*model.ProjectManagementWorkflow, 0)
	if err := tx.WithContext(ctx).Where("project_id = ? AND kind = ?", projectId, constants.ProjectManagementWorkflowKindProject).Find(&workflows).Error; err != nil {
		return errorx.WithStack(err)
	}
	for _, workflow := range workflows {
		if err := r.DeleteTemplateWithTx(ctx, tx, workflow.Id); err != nil {
			return err
		}
	}
	return nil
}

func (r *ProjectManagementWorkflowRepo) CloneTemplateToProjectWithTx(ctx context.Context, tx *gorm.DB, templateId, organizationId, projectId int64) error {
	template := new(model.ProjectManagementWorkflow)
	err := tx.WithContext(ctx).Clauses(clause.Locking{Strength: "UPDATE"}).
		Preload("Stages", func(db *gorm.DB) *gorm.DB { return db.Order("sort_order ASC, id ASC") }).
		Preload("Transitions", func(db *gorm.DB) *gorm.DB { return db.Order("id ASC") }).
		Where("id = ? AND organization_id = ? AND kind = ? AND enabled = true", templateId, organizationId, constants.ProjectManagementWorkflowKindOrganizationTemplate).
		First(template).Error
	if err != nil {
		return errorx.WithStack(err)
	}
	workflow := &model.ProjectManagementWorkflow{Kind: constants.ProjectManagementWorkflowKindProject, OrganizationId: organizationId, ProjectId: projectId, Name: template.Name, Description: template.Description, Enabled: true, SourceTemplateId: template.Id, SourceTemplateName: template.Name}
	if err = tx.WithContext(ctx).Create(workflow).Error; err != nil {
		return errorx.WithStack(err)
	}
	stageIds := make(map[int64]int64, len(template.Stages))
	for _, stage := range template.Stages {
		copy := &model.ProjectManagementWorkflowStage{WorkflowId: workflow.Id, Name: stage.Name, Phase: stage.Phase, SortOrder: stage.SortOrder}
		if err = tx.WithContext(ctx).Create(copy).Error; err != nil {
			return errorx.WithStack(err)
		}
		stageIds[stage.Id] = copy.Id
	}
	for _, transition := range template.Transitions {
		copy := &model.ProjectManagementWorkflowTransition{WorkflowId: workflow.Id, FromStageId: stageIds[transition.FromStageId], ToStageId: stageIds[transition.ToStageId], AllowedActorTypes: transition.AllowedActorTypes}
		if err = tx.WithContext(ctx).Create(copy).Error; err != nil {
			return errorx.WithStack(err)
		}
	}
	return nil
}
