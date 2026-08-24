package repo

import (
	"context"
	"encoding/json"

	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/gormx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/model"
)

type ProjectManagementWorkflowRepo struct {
	*gormx.BaseRepo[*model.ProjectManagementWorkflow]
	stageRepo      *gormx.BaseRepo[*model.ProjectManagementWorkflowStage]
	transitionRepo *gormx.BaseRepo[*model.ProjectManagementWorkflowTransition]
}

func NewProjectManagementWorkflowRepo(db *gorm.DB) *ProjectManagementWorkflowRepo {
	return &ProjectManagementWorkflowRepo{
		BaseRepo:       gormx.NewBaseRepo[*model.ProjectManagementWorkflow](db),
		stageRepo:      gormx.NewBaseRepo[*model.ProjectManagementWorkflowStage](db),
		transitionRepo: gormx.NewBaseRepo[*model.ProjectManagementWorkflowTransition](db),
	}
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

func (r *ProjectManagementWorkflowRepo) ListByOrganization(ctx context.Context, organizationId int64) ([]*model.ProjectManagementWorkflow, error) {
	workflows := make([]*model.ProjectManagementWorkflow, 0)
	err := r.Client(ctx).
		Preload("Stages", func(db *gorm.DB) *gorm.DB { return db.Order("sort_order ASC, id ASC") }).
		// Preload("Transitions", func(db *gorm.DB) *gorm.DB { return db.Order("id ASC") }).
		Where("organization_id = ? AND kind = ?", organizationId, constants.ProjectManagementWorkflowKindOrganizationTemplate).
		Order("is_default DESC, id ASC").
		Find(&workflows).Error
	return workflows, errorx.WithStack(err)
}

func (r *ProjectManagementWorkflowRepo) FindPendingStageWithTx(ctx context.Context, tx *gorm.DB, workflowId int64) (*model.ProjectManagementWorkflowStage, error) {
	return r.stageRepo.FindOneByColumnsWithTx(ctx, tx, map[string]any{"workflow_id": workflowId, "phase": constants.ProjectManagementStagePhasePending})
}

func (r *ProjectManagementWorkflowRepo) ListStagesWithTx(ctx context.Context, tx *gorm.DB, workflowId int64) ([]*model.ProjectManagementWorkflowStage, error) {
	return r.stageRepo.FindAllWithTx(ctx, tx, map[string]any{"workflow_id": workflowId})
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
	return r.stageRepo.FindOneByIdWithTx(ctx, tx, id)
}

func (r *ProjectManagementWorkflowRepo) FindStage(ctx context.Context, id int64) (*model.ProjectManagementWorkflowStage, error) {
	return r.stageRepo.FindOneById(ctx, id)
}

func (r *ProjectManagementWorkflowRepo) ListTemplateStages(ctx context.Context, workflowId, organizationId int64) ([]*model.ProjectManagementWorkflowStage, error) {
	stages := make([]*model.ProjectManagementWorkflowStage, 0)
	err := r.stageRepo.Client(ctx).Where("workflow_id = ? AND EXISTS (SELECT 1 FROM project_management_workflows WHERE id = project_management_workflow_stages.workflow_id AND organization_id = ? AND kind = ?)", workflowId, organizationId, constants.ProjectManagementWorkflowKindOrganizationTemplate).
		Order("sort_order ASC, id ASC").Find(&stages).Error
	return stages, errorx.WithStack(err)
}

func (r *ProjectManagementWorkflowRepo) FindTemplateStage(ctx context.Context, id, workflowId, organizationId int64) (*model.ProjectManagementWorkflowStage, error) {
	stage := new(model.ProjectManagementWorkflowStage)
	err := r.stageRepo.Client(ctx).
		Preload("Transitions", func(db *gorm.DB) *gorm.DB { return db.Order("id ASC") }).
		Where(
			"id = ? AND workflow_id = ? AND EXISTS (SELECT 1 FROM project_management_workflows WHERE id = project_management_workflow_stages.workflow_id AND organization_id = ? AND kind = ?)",
			id, workflowId, organizationId, constants.ProjectManagementWorkflowKindOrganizationTemplate).
		First(stage).Error
	return stage, errorx.WithStack(err)
}

func (r *ProjectManagementWorkflowRepo) FindTemplateStageWithTx(ctx context.Context, tx *gorm.DB, id, workflowId, organizationId int64) (*model.ProjectManagementWorkflowStage, error) {
	stage := new(model.ProjectManagementWorkflowStage)
	err := tx.WithContext(ctx).Model(new(model.ProjectManagementWorkflowStage)).Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ? AND workflow_id = ? AND EXISTS (SELECT 1 FROM project_management_workflows WHERE id = project_management_workflow_stages.workflow_id AND organization_id = ? AND kind = ?)", id, workflowId, organizationId, constants.ProjectManagementWorkflowKindOrganizationTemplate).First(stage).Error
	return stage, errorx.WithStack(err)
}

func (r *ProjectManagementWorkflowRepo) ListTransitionsByFromStage(ctx context.Context, fromStageId int64) ([]*model.ProjectManagementWorkflowTransition, error) {
	transitions := make([]*model.ProjectManagementWorkflowTransition, 0)
	err := r.transitionRepo.Client(ctx).
		Preload("FromStage").
		Preload("ToStage").
		Where("from_stage_id = ?", fromStageId).
		Order("id ASC").
		Find(&transitions).Error
	return transitions, errorx.WithStack(err)
}

func (r *ProjectManagementWorkflowRepo) FindTemplateTransition(ctx context.Context, id, workflowId, organizationId int64) (*model.ProjectManagementWorkflowTransition, error) {
	transition := new(model.ProjectManagementWorkflowTransition)
	err := r.transitionRepo.Client(ctx).
		Preload("FromStage").
		Preload("ToStage").
		Where("id = ? AND workflow_id = ? AND EXISTS (SELECT 1 FROM project_management_workflows WHERE id = project_management_workflow_transitions.workflow_id AND organization_id = ? AND kind = ?)", id, workflowId, organizationId, constants.ProjectManagementWorkflowKindOrganizationTemplate).
		First(transition).Error
	return transition, errorx.WithStack(err)
}

func (r *ProjectManagementWorkflowRepo) FindTemplateTransitionWithTx(ctx context.Context, tx *gorm.DB, id, workflowId, organizationId int64) (*model.ProjectManagementWorkflowTransition, error) {
	transition := new(model.ProjectManagementWorkflowTransition)
	err := tx.WithContext(ctx).Model(new(model.ProjectManagementWorkflowTransition)).Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ? AND workflow_id = ? AND EXISTS (SELECT 1 FROM project_management_workflows WHERE id = project_management_workflow_transitions.workflow_id AND organization_id = ? AND kind = ?)", id, workflowId, organizationId, constants.ProjectManagementWorkflowKindOrganizationTemplate).First(transition).Error
	return transition, errorx.WithStack(err)
}

func (r *ProjectManagementWorkflowRepo) CreateTemplateStageWithTx(ctx context.Context, tx *gorm.DB, stage *model.ProjectManagementWorkflowStage) error {
	_, err := r.stageRepo.CreateWithTx(ctx, tx, stage, false)
	return err
}

func (r *ProjectManagementWorkflowRepo) UpdateTemplateStageWithTx(ctx context.Context, tx *gorm.DB, stage *model.ProjectManagementWorkflowStage) error {
	return r.stageRepo.UpdateColumnsWithTx(ctx, tx,
		map[string]any{"name": stage.Name, "phase": stage.Phase, "color": stage.Color, "sort_order": stage.SortOrder},
		map[string]any{"id": stage.Id, "workflow_id": stage.WorkflowId}, nil)
}

func (r *ProjectManagementWorkflowRepo) CreateTemplateTransitionWithTx(ctx context.Context, tx *gorm.DB, transition *model.ProjectManagementWorkflowTransition) error {
	_, err := r.transitionRepo.CreateWithTx(ctx, tx, transition, false)
	return err
}

func (r *ProjectManagementWorkflowRepo) UpdateTemplateTransitionWithTx(ctx context.Context, tx *gorm.DB, transition *model.ProjectManagementWorkflowTransition) error {
	allowedActorTypes, err := json.Marshal(transition.AllowedActorTypes)
	if err != nil {
		return errorx.WithStack(err)
	}
	return r.transitionRepo.UpdateColumnsWithTx(ctx, tx,
		map[string]any{"name": transition.Name, "from_stage_id": transition.FromStageId, "to_stage_id": transition.ToStageId, "allowed_actor_types": datatypes.JSON(allowedActorTypes)},
		map[string]any{"id": transition.Id, "workflow_id": transition.WorkflowId}, nil)
}

func (r *ProjectManagementWorkflowRepo) DeleteTemplateTransitionWithTx(ctx context.Context, tx *gorm.DB, transition *model.ProjectManagementWorkflowTransition) error {
	_, err := r.transitionRepo.DeleteWithTx(ctx, tx, map[string]any{"id": transition.Id, "workflow_id": transition.WorkflowId})
	return err
}

func (r *ProjectManagementWorkflowRepo) FindTransitionWithTx(ctx context.Context, tx *gorm.DB, workflowId, fromStageId, toStageId int64) (*model.ProjectManagementWorkflowTransition, error) {
	return r.transitionRepo.FindOneByColumnsWithTx(ctx, tx, map[string]any{"workflow_id": workflowId, "from_stage_id": fromStageId, "to_stage_id": toStageId})
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

func (r *ProjectManagementWorkflowRepo) LockDefaultTemplateDefinitionWithTx(ctx context.Context, tx *gorm.DB, organizationId int64) (*model.ProjectManagementWorkflow, error) {
	workflow := new(model.ProjectManagementWorkflow)
	err := tx.WithContext(ctx).Clauses(clause.Locking{Strength: "UPDATE"}).
		Preload("Stages", func(db *gorm.DB) *gorm.DB { return db.Order("sort_order ASC, id ASC") }).
		Preload("Transitions", func(db *gorm.DB) *gorm.DB { return db.Order("id ASC") }).
		Where("organization_id = ? AND kind = ? AND is_default = true", organizationId, constants.ProjectManagementWorkflowKindOrganizationTemplate).
		First(workflow).Error
	return workflow, errorx.WithStack(err)
}

func (r *ProjectManagementWorkflowRepo) LockTemplatesWithTx(ctx context.Context, tx *gorm.DB, organizationId int64) ([]*model.ProjectManagementWorkflow, error) {
	workflows := make([]*model.ProjectManagementWorkflow, 0)
	err := tx.WithContext(ctx).Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("organization_id = ? AND kind = ?", organizationId, constants.ProjectManagementWorkflowKindOrganizationTemplate).
		Order("id ASC").Find(&workflows).Error
	return workflows, errorx.WithStack(err)
}

func (r *ProjectManagementWorkflowRepo) DemoteDefaultWithTx(ctx context.Context, tx *gorm.DB, organizationId, exceptId int64) error {
	err := tx.WithContext(ctx).Model(&model.ProjectManagementWorkflow{}).
		Where("organization_id = ? AND kind = ? AND id <> ?", organizationId, constants.ProjectManagementWorkflowKindOrganizationTemplate, exceptId).
		Update("is_default", false).Error
	return errorx.WithStack(err)
}

func (r *ProjectManagementWorkflowRepo) SetDefaultWithTx(ctx context.Context, tx *gorm.DB, id, organizationId int64) error {
	err := tx.WithContext(ctx).Model(&model.ProjectManagementWorkflow{}).
		Where("id = ? AND organization_id = ? AND kind = ?", id, organizationId, constants.ProjectManagementWorkflowKindOrganizationTemplate).
		Update("is_default", true).Error
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
		if _, err := r.stageRepo.CreateWithTx(ctx, tx, stage, false); err != nil {
			return err
		}
		stageIds[referenceId] = stage.Id
	}
	for _, transition := range transitions {
		transition.Id = 0
		transition.WorkflowId = workflowId
		transition.FromStageId = stageIds[transition.FromStageId]
		transition.ToStageId = stageIds[transition.ToStageId]
		if _, err := r.transitionRepo.CreateWithTx(ctx, tx, transition, false); err != nil {
			return err
		}
	}
	return nil
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
			if err := r.stageRepo.UpdateColumnsWithTx(ctx, tx,
				map[string]any{"name": stage.Name, "phase": stage.Phase, "color": stage.Color, "sort_order": stage.SortOrder},
				map[string]any{"id": referenceId}, nil); err != nil {
				return err
			}
			stageIds[referenceId], retained[referenceId] = referenceId, struct{}{}
			continue
		}
		stage.Id, stage.WorkflowId = 0, workflowId
		if _, err := r.stageRepo.CreateWithTx(ctx, tx, stage, false); err != nil {
			return err
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
			if _, err := r.stageRepo.DeleteByIDWithTx(ctx, tx, stage.Id); err != nil {
				return err
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
		if _, err := r.transitionRepo.CreateWithTx(ctx, tx, transition, false); err != nil {
			return err
		}
	}
	return nil
}

func (r *ProjectManagementWorkflowRepo) DeleteTemplateWithTx(ctx context.Context, tx *gorm.DB, workflowId int64) error {
	if _, err := r.transitionRepo.DeleteWithTx(ctx, tx, map[string]any{"workflow_id": workflowId}); err != nil {
		return err
	}
	if _, err := r.stageRepo.DeleteWithTx(ctx, tx, map[string]any{"workflow_id": workflowId}); err != nil {
		return err
	}
	_, err := r.DeleteByIDWithTx(ctx, tx, workflowId)
	return err
}

func (r *ProjectManagementWorkflowRepo) DeleteTemplateStageWithTx(ctx context.Context, tx *gorm.DB, stage *model.ProjectManagementWorkflowStage) error {
	if err := tx.WithContext(ctx).Where("workflow_id = ? AND (from_stage_id = ? OR to_stage_id = ?)", stage.WorkflowId, stage.Id, stage.Id).Delete(&model.ProjectManagementWorkflowTransition{}).Error; err != nil {
		return errorx.WithStack(err)
	}
	_, err := r.stageRepo.DeleteWithTx(ctx, tx, map[string]any{"id": stage.Id, "workflow_id": stage.WorkflowId})
	return err
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
		Where("id = ? AND organization_id = ? AND kind = ?", templateId, organizationId, constants.ProjectManagementWorkflowKindOrganizationTemplate).
		First(template).Error
	if err != nil {
		return errorx.WithStack(err)
	}
	workflow, err := copierx.Trans[model.ProjectManagementWorkflow, model.ProjectManagementWorkflow](template)
	if err != nil {
		return err
	}
	workflow.Base = model.Base{}
	workflow.RemoveRelationships()
	workflow.Kind = constants.ProjectManagementWorkflowKindProject
	workflow.OrganizationId = organizationId
	workflow.ProjectId = projectId
	workflow.IsDefault = false
	workflow.SourceTemplateId = template.Id
	workflow.SourceTemplateName = template.Name
	if _, err = r.CreateWithTx(ctx, tx, workflow, false); err != nil {
		return err
	}
	stageIds := make(map[int64]int64, len(template.Stages))
	for _, stage := range template.Stages {
		copyStage, err := copierx.Trans[model.ProjectManagementWorkflowStage, model.ProjectManagementWorkflowStage](stage)
		if err != nil {
			return err
		}
		copyStage.Base = model.Base{}
		copyStage.RemoveRelationships()
		copyStage.WorkflowId = workflow.Id
		if _, err = r.stageRepo.CreateWithTx(ctx, tx, copyStage, false); err != nil {
			return err
		}
		stageIds[stage.Id] = copyStage.Id
	}
	for _, transition := range template.Transitions {
		copyTransition, err := copierx.Trans[model.ProjectManagementWorkflowTransition, model.ProjectManagementWorkflowTransition](transition)
		if err != nil {
			return err
		}
		copyTransition.Base = model.Base{}
		copyTransition.RemoveRelationships()
		copyTransition.WorkflowId = workflow.Id
		copyTransition.FromStageId = stageIds[transition.FromStageId]
		copyTransition.ToStageId = stageIds[transition.ToStageId]
		if _, err = r.transitionRepo.CreateWithTx(ctx, tx, copyTransition, false); err != nil {
			return err
		}
	}
	return nil
}
