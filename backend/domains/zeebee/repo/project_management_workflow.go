package repo

import (
	"context"
	"errors"

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

func preloadWorkflowTransitions(db *gorm.DB) *gorm.DB {
	return db.Preload("FromStage").Preload("ToStage").Order("project_management_workflow_transitions.id ASC")
}

func NewProjectManagementWorkflowRepo(db *gorm.DB) *ProjectManagementWorkflowRepo {
	return &ProjectManagementWorkflowRepo{BaseRepo: gormx.NewBaseRepo[*model.ProjectManagementWorkflow](db)}
}

func (r *ProjectManagementWorkflowRepo) ListWorkflowStagesWithTx(ctx context.Context, tx *gorm.DB, workflowId int64) ([]*model.ProjectManagementStage, error) {
	stages := make([]*model.ProjectManagementStage, 0)
	err := tx.WithContext(ctx).Model(&model.ProjectManagementStage{}).
		Select("project_management_stages.*").
		Joins("JOIN project_management_workflow_stages ON project_management_workflow_stages.stage_id = project_management_stages.id AND project_management_workflow_stages.is_del = 0").
		Where("project_management_workflow_stages.workflow_id = ?", workflowId).
		Order("project_management_stages.sort_order ASC, project_management_stages.id ASC").Find(&stages).Error
	return stages, errorx.WithStack(err)
}

func (r *ProjectManagementWorkflowRepo) ListStageLinksWithTx(ctx context.Context, tx *gorm.DB, workflowId int64) ([]*model.ProjectManagementWorkflowStage, error) {
	links := make([]*model.ProjectManagementWorkflowStage, 0)
	err := tx.WithContext(ctx).Model(&model.ProjectManagementWorkflowStage{}).Where("workflow_id = ?", workflowId).Order("id ASC").Find(&links).Error
	return links, errorx.WithStack(err)
}

func (r *ProjectManagementWorkflowRepo) ListStageTransitionsWithTx(ctx context.Context, tx *gorm.DB, workflowId int64) ([]*model.ProjectManagementWorkflowTransition, error) {
	transitions := make([]*model.ProjectManagementWorkflowTransition, 0)
	err := tx.WithContext(ctx).Model(&model.ProjectManagementWorkflowTransition{}).Where("workflow_id = ?", workflowId).Order("id ASC").Find(&transitions).Error
	return transitions, errorx.WithStack(err)
}

func (r *ProjectManagementWorkflowRepo) FindProjectWorkflowByIDWithTx(ctx context.Context, tx *gorm.DB, workflowId, projectId, organizationId int64) (*model.ProjectManagementWorkflow, error) {
	workflow := new(model.ProjectManagementWorkflow)
	err := tx.WithContext(ctx).Model(&model.ProjectManagementWorkflow{}).
		Joins("JOIN project_management_projects ON project_management_projects.id = project_management_workflows.entity_id AND project_management_projects.is_del = 0").
		Where("project_management_workflows.id = ? AND project_management_workflows.entity_type = ? AND project_management_workflows.entity_id = ? AND project_management_projects.organization_id = ?", workflowId, constants.ProjectManagementOwnerTypeProject, projectId, organizationId).
		First(workflow).Error
	return workflow, errorx.WithStack(err)
}

func (r *ProjectManagementWorkflowRepo) FindDefaultProjectWorkflowWithTx(ctx context.Context, tx *gorm.DB, projectId, organizationId int64) (*model.ProjectManagementWorkflow, error) {
	workflow := new(model.ProjectManagementWorkflow)
	err := tx.WithContext(ctx).Model(&model.ProjectManagementWorkflow{}).
		Joins("JOIN project_management_projects ON project_management_projects.id = project_management_workflows.entity_id AND project_management_projects.is_del = 0").
		Where("project_management_workflows.entity_type = ? AND project_management_workflows.entity_id = ? AND project_management_projects.organization_id = ? AND project_management_workflows.is_default = 1", constants.ProjectManagementOwnerTypeProject, projectId, organizationId).
		First(workflow).Error
	return workflow, errorx.WithStack(err)
}

func (r *ProjectManagementWorkflowRepo) ListWorkflowSources(ctx context.Context, organizationId int64) ([]*model.ProjectManagementWorkflow, error) {
	workflows := make([]*model.ProjectManagementWorkflow, 0)
	err := r.sourceQuery(r.Client(ctx), organizationId).
		Preload("Stages", func(db *gorm.DB) *gorm.DB { return db.Order("sort_order ASC, id ASC") }).
		Order("project_management_workflows.is_default DESC, project_management_workflows.id ASC").Find(&workflows).Error
	return workflows, errorx.WithStack(err)
}

func (r *ProjectManagementWorkflowRepo) FindWorkflowSourceByID(ctx context.Context, workflowId, organizationId int64) (*model.ProjectManagementWorkflow, error) {
	workflow := new(model.ProjectManagementWorkflow)
	err := r.sourceQuery(r.Client(ctx), organizationId).
		Preload("Stages", func(db *gorm.DB) *gorm.DB { return db.Order("sort_order ASC, id ASC") }).
		Preload("Transitions", preloadWorkflowTransitions).
		Where("project_management_workflows.id = ?", workflowId).First(workflow).Error
	return workflow, errorx.WithStack(err)
}

func (r *ProjectManagementWorkflowRepo) FindWorkflowSourceByIDWithTx(ctx context.Context, tx *gorm.DB, workflowId, organizationId int64) (*model.ProjectManagementWorkflow, error) {
	workflow := new(model.ProjectManagementWorkflow)
	err := r.sourceQuery(tx.WithContext(ctx), organizationId).Where("project_management_workflows.id = ?", workflowId).First(workflow).Error
	return workflow, errorx.WithStack(err)
}

func (r *ProjectManagementWorkflowRepo) FindDefaultWorkflowSourceWithTx(ctx context.Context, tx *gorm.DB, organizationId int64) (*model.ProjectManagementWorkflow, error) {
	workflow := new(model.ProjectManagementWorkflow)
	err := r.sourceQuery(tx.WithContext(ctx), organizationId).Where("project_management_workflows.is_default = 1").First(workflow).Error
	return workflow, errorx.WithStack(err)
}

func (r *ProjectManagementWorkflowRepo) LockWorkflowSourceWithTx(ctx context.Context, tx *gorm.DB, workflowId, organizationId int64) (*model.ProjectManagementWorkflow, error) {
	workflow := new(model.ProjectManagementWorkflow)
	err := r.sourceQuery(tx.WithContext(ctx).Clauses(clause.Locking{Strength: "UPDATE"}), organizationId).Where("project_management_workflows.id = ?", workflowId).First(workflow).Error
	return workflow, errorx.WithStack(err)
}

func (r *ProjectManagementWorkflowRepo) LockWorkflowSourcesWithTx(ctx context.Context, tx *gorm.DB, organizationId int64) ([]*model.ProjectManagementWorkflow, error) {
	workflows := make([]*model.ProjectManagementWorkflow, 0)
	err := r.sourceQuery(tx.WithContext(ctx).Clauses(clause.Locking{Strength: "UPDATE"}), organizationId).Order("project_management_workflows.id ASC").Find(&workflows).Error
	return workflows, errorx.WithStack(err)
}

func (r *ProjectManagementWorkflowRepo) sourceQuery(db *gorm.DB, organizationId int64) *gorm.DB {
	return db.Model(&model.ProjectManagementWorkflow{}).
		Joins("JOIN project_management_templates ON project_management_templates.id = project_management_workflows.entity_id AND project_management_templates.template_type = ? AND project_management_templates.status = ? AND project_management_templates.is_del = 0", constants.ProjectManagementTemplateTypeOrganization, "Active").
		Where("project_management_workflows.entity_type = ? AND project_management_templates.organization_id = ?", constants.ProjectManagementOwnerTypeTemplate, organizationId)
}

func (r *ProjectManagementWorkflowRepo) ListProjectWorkflows(ctx context.Context, projectId, organizationId int64) ([]*model.ProjectManagementWorkflow, error) {
	workflows := make([]*model.ProjectManagementWorkflow, 0)
	err := r.Client(ctx).Model(&model.ProjectManagementWorkflow{}).
		Where("entity_type = ? AND entity_id = ?", constants.ProjectManagementOwnerTypeProject, projectId).
		Where("EXISTS (SELECT 1 FROM project_management_projects WHERE id = ? AND organization_id = ?)", projectId, organizationId).
		Order("is_default DESC, id ASC").Find(&workflows).Error
	return workflows, errorx.WithStack(err)
}

func (r *ProjectManagementWorkflowRepo) ListTemplateWorkflowsWithTx(ctx context.Context, tx *gorm.DB, templateId int64) ([]*model.ProjectManagementWorkflow, error) {
	workflows := make([]*model.ProjectManagementWorkflow, 0)
	err := tx.WithContext(ctx).Model(&model.ProjectManagementWorkflow{}).
		Where("entity_type = ? AND entity_id = ?", constants.ProjectManagementOwnerTypeTemplate, templateId).
		Order("id ASC").Find(&workflows).Error
	return workflows, errorx.WithStack(err)
}

func (r *ProjectManagementWorkflowRepo) ListTemplateWorkflows(ctx context.Context, templateId int64) ([]*model.ProjectManagementWorkflow, error) {
	workflows := make([]*model.ProjectManagementWorkflow, 0)
	err := r.Client(ctx).Model(&model.ProjectManagementWorkflow{}).
		Where("entity_type = ? AND entity_id = ?", constants.ProjectManagementOwnerTypeTemplate, templateId).
		Order("is_default DESC, id ASC").Find(&workflows).Error
	return workflows, errorx.WithStack(err)
}

func (r *ProjectManagementWorkflowRepo) ListProjectWorkflowsWithTx(ctx context.Context, tx *gorm.DB, projectId int64) ([]*model.ProjectManagementWorkflow, error) {
	workflows := make([]*model.ProjectManagementWorkflow, 0)
	err := tx.WithContext(ctx).Model(&model.ProjectManagementWorkflow{}).
		Where("entity_type = ? AND entity_id = ?", constants.ProjectManagementOwnerTypeProject, projectId).
		Order("id ASC").Find(&workflows).Error
	return workflows, errorx.WithStack(err)
}

func (r *ProjectManagementWorkflowRepo) PageProjectWorkflows(ctx context.Context, projectId, organizationId int64, pageReq gormx.PageRequest) (gormx.PageResponse[*model.ProjectManagementWorkflow], error) {
	query := r.Client(ctx).Model(&model.ProjectManagementWorkflow{}).
		Where("entity_type = ? AND entity_id = ?", constants.ProjectManagementOwnerTypeProject, projectId).
		Where("EXISTS (SELECT 1 FROM project_management_projects WHERE id = ? AND organization_id = ?)", projectId, organizationId).
		Order("is_default DESC, id ASC")
	return r.BaseRepo.Page(query, pageReq)
}

func (r *ProjectManagementWorkflowRepo) FindProjectWorkflowByID(ctx context.Context, workflowId, projectId, organizationId int64) (*model.ProjectManagementWorkflow, error) {
	workflow := new(model.ProjectManagementWorkflow)
	err := r.Client(ctx).Model(&model.ProjectManagementWorkflow{}).
		Preload("Stages", func(db *gorm.DB) *gorm.DB { return db.Order("sort_order ASC, id ASC") }).
		Preload("Transitions", preloadWorkflowTransitions).
		Where("id = ? AND entity_type = ? AND entity_id = ?", workflowId, constants.ProjectManagementOwnerTypeProject, projectId).
		Where("EXISTS (SELECT 1 FROM project_management_projects WHERE id = ? AND organization_id = ?)", projectId, organizationId).
		First(workflow).Error
	return workflow, errorx.WithStack(err)
}

func (r *ProjectManagementWorkflowRepo) FindPendingOwnerStageWithTx(ctx context.Context, tx *gorm.DB, workflowId int64) (*model.ProjectManagementStage, error) {
	stage := new(model.ProjectManagementStage)
	err := tx.WithContext(ctx).Model(&model.ProjectManagementStage{}).
		Joins("JOIN project_management_workflow_stages ON project_management_workflow_stages.stage_id = project_management_stages.id AND project_management_workflow_stages.is_del = 0").
		Where("project_management_workflow_stages.workflow_id = ? AND project_management_stages.phase = ?", workflowId, constants.ProjectManagementStagePhasePending).
		First(stage).Error
	return stage, errorx.WithStack(err)
}

func (r *ProjectManagementWorkflowRepo) ListPendingOwnerStagesWithTx(ctx context.Context, tx *gorm.DB, workflowIds []int64) (map[int64]*model.ProjectManagementStage, error) {
	if len(workflowIds) == 0 {
		return map[int64]*model.ProjectManagementStage{}, nil
	}
	type pendingStage struct {
		WorkflowId int64
		model.ProjectManagementStage
	}
	rows := make([]pendingStage, 0, len(workflowIds))
	err := tx.WithContext(ctx).Model(&model.ProjectManagementWorkflowStage{}).
		Select("project_management_workflow_stages.workflow_id, project_management_stages.*").
		Joins("JOIN project_management_stages ON project_management_stages.id = project_management_workflow_stages.stage_id AND project_management_stages.is_del = 0").
		Where("project_management_workflow_stages.workflow_id IN ? AND project_management_stages.phase = ?", workflowIds, constants.ProjectManagementStagePhasePending).
		Find(&rows).Error
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	result := make(map[int64]*model.ProjectManagementStage, len(rows))
	for _, row := range rows {
		stage := row.ProjectManagementStage
		if _, exists := result[row.WorkflowId]; exists {
			return nil, errorx.New("workflow has multiple pending stages")
		}
		result[row.WorkflowId] = &stage
	}
	return result, nil
}

func (r *ProjectManagementWorkflowRepo) FindOwnerStageWithTx(ctx context.Context, tx *gorm.DB, stageId int64, entityType constants.ProjectManagementOwnerType, entityId int64) (*model.ProjectManagementStage, error) {
	stage := new(model.ProjectManagementStage)
	err := tx.WithContext(ctx).Model(&model.ProjectManagementStage{}).Where("id = ? AND entity_type = ? AND entity_id = ?", stageId, entityType, entityId).First(stage).Error
	return stage, errorx.WithStack(err)
}

func (r *ProjectManagementWorkflowRepo) FindOwnerStageInWorkflowWithTx(ctx context.Context, tx *gorm.DB, workflowId, stageId int64, entityType constants.ProjectManagementOwnerType, entityId int64) (*model.ProjectManagementStage, error) {
	stage := new(model.ProjectManagementStage)
	err := tx.WithContext(ctx).Model(&model.ProjectManagementStage{}).
		Joins("JOIN project_management_workflow_stages ON project_management_workflow_stages.stage_id = project_management_stages.id AND project_management_workflow_stages.is_del = 0").
		Where("project_management_workflow_stages.workflow_id = ? AND project_management_stages.id = ? AND project_management_stages.entity_type = ? AND project_management_stages.entity_id = ?", workflowId, stageId, entityType, entityId).
		First(stage).Error
	return stage, errorx.WithStack(err)
}

func (r *ProjectManagementWorkflowRepo) DeleteWorkflowSourceWithTx(ctx context.Context, tx *gorm.DB, workflow *model.ProjectManagementWorkflow) error {
	if err := tx.WithContext(ctx).Where("workflow_id = ?", workflow.Id).Delete(&model.ProjectManagementWorkflowTransition{}).Error; err != nil {
		return errorx.WithStack(err)
	}
	if err := tx.WithContext(ctx).Where("workflow_id = ?", workflow.Id).Delete(&model.ProjectManagementWorkflowStage{}).Error; err != nil {
		return errorx.WithStack(err)
	}
	if err := tx.WithContext(ctx).Where("entity_type = ? AND entity_id = ?", constants.ProjectManagementOwnerTypeTemplate, workflow.EntityId).Delete(&model.ProjectManagementStage{}).Error; err != nil {
		return errorx.WithStack(err)
	}
	if err := tx.WithContext(ctx).Delete(workflow).Error; err != nil {
		return errorx.WithStack(err)
	}
	return errorx.WithStack(tx.WithContext(ctx).Delete(&model.ProjectManagementTemplate{}, workflow.EntityId).Error)
}

func (r *ProjectManagementWorkflowRepo) DeleteProjectWorkflowWithTx(ctx context.Context, tx *gorm.DB, projectId int64) error {
	workflows := make([]*model.ProjectManagementWorkflow, 0)
	if err := tx.WithContext(ctx).Model(&model.ProjectManagementWorkflow{}).Where("entity_type = ? AND entity_id = ?", constants.ProjectManagementOwnerTypeProject, projectId).Find(&workflows).Error; err != nil {
		return errorx.WithStack(err)
	}
	for _, workflow := range workflows {
		if err := tx.WithContext(ctx).Where("workflow_id = ?", workflow.Id).Delete(&model.ProjectManagementWorkflowTransition{}).Error; err != nil {
			return errorx.WithStack(err)
		}
		if err := tx.WithContext(ctx).Where("workflow_id = ?", workflow.Id).Delete(&model.ProjectManagementWorkflowStage{}).Error; err != nil {
			return errorx.WithStack(err)
		}
		if err := tx.WithContext(ctx).Delete(workflow).Error; err != nil {
			return errorx.WithStack(err)
		}
	}
	return errorx.WithStack(tx.WithContext(ctx).Where("entity_type = ? AND entity_id = ?", constants.ProjectManagementOwnerTypeProject, projectId).Delete(&model.ProjectManagementStage{}).Error)
}

func (r *ProjectManagementWorkflowRepo) CloneTemplateToProjectWithTx(ctx context.Context, tx *gorm.DB, templateId, organizationId, projectId int64) error {
	if templateId == 0 {
		source, err := r.FindDefaultWorkflowSourceWithTx(ctx, tx, organizationId)
		if err == nil {
			return r.cloneWorkflowToOwnerWithTx(ctx, tx, source, constants.ProjectManagementOwnerTypeProject, projectId)
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		return r.createDefaultProjectWorkflowWithTx(ctx, tx, projectId)
	}
	template := new(model.ProjectManagementTemplate)
	if err := tx.WithContext(ctx).Model(&model.ProjectManagementTemplate{}).
		Where("id = ? AND ((template_type = ? AND organization_id = 0) OR (template_type = ? AND organization_id = ?)) AND status = ?", templateId, constants.ProjectManagementTemplateTypePublic, constants.ProjectManagementTemplateTypeOrganization, organizationId, "Active").
		First(template).Error; err != nil {
		return errorx.WithStack(err)
	}
	workflows := make([]*model.ProjectManagementWorkflow, 0)
	if err := tx.WithContext(ctx).Model(&model.ProjectManagementWorkflow{}).
		Where("entity_type = ? AND entity_id = ?", constants.ProjectManagementOwnerTypeTemplate, templateId).
		Order("is_default DESC, id ASC").Find(&workflows).Error; err != nil {
		return errorx.WithStack(err)
	}
	if len(workflows) == 0 {
		return errorx.New("project template has no workflows")
	}
	if !workflows[0].IsDefault {
		workflows[0].IsDefault = true
	}
	for _, workflow := range workflows {
		if err := r.cloneWorkflowToOwnerWithTx(ctx, tx, workflow, constants.ProjectManagementOwnerTypeProject, projectId); err != nil {
			return err
		}
	}
	return nil
}

func (r *ProjectManagementWorkflowRepo) CloneProjectToTemplateWithTx(ctx context.Context, tx *gorm.DB, projectId, templateId int64) error {
	workflows := make([]*model.ProjectManagementWorkflow, 0)
	if err := tx.WithContext(ctx).Model(&model.ProjectManagementWorkflow{}).
		Where("entity_type = ? AND entity_id = ?", constants.ProjectManagementOwnerTypeProject, projectId).
		Order("is_default DESC, id ASC").Find(&workflows).Error; err != nil {
		return errorx.WithStack(err)
	}
	if len(workflows) == 0 {
		return errorx.New("project has no workflows")
	}
	for _, workflow := range workflows {
		if err := r.cloneWorkflowToOwnerWithTx(ctx, tx, workflow, constants.ProjectManagementOwnerTypeTemplate, templateId); err != nil {
			return err
		}
	}
	return nil
}

func (r *ProjectManagementWorkflowRepo) cloneWorkflowToOwnerWithTx(ctx context.Context, tx *gorm.DB, source *model.ProjectManagementWorkflow, entityType constants.ProjectManagementOwnerType, entityId int64) error {
	stages, err := r.ListWorkflowStagesWithTx(ctx, tx, source.Id)
	if err != nil {
		return err
	}
	links, err := r.ListStageLinksWithTx(ctx, tx, source.Id)
	if err != nil {
		return err
	}
	transitions, err := r.ListStageTransitionsWithTx(ctx, tx, source.Id)
	if err != nil {
		return err
	}
	workflow := &model.ProjectManagementWorkflow{EntityType: entityType, EntityId: entityId, SourceWorkflowId: source.Id, SourceWorkflowName: source.Name, Name: source.Name, Description: source.Description, IsDefault: source.IsDefault}
	if err = tx.WithContext(ctx).Create(workflow).Error; err != nil {
		return errorx.WithStack(err)
	}
	stageIDs := make(map[int64]int64, len(stages))
	ownerStages := make([]*model.ProjectManagementStage, 0)
	if entityType == constants.ProjectManagementOwnerTypeProject {
		if err = tx.WithContext(ctx).Model(&model.ProjectManagementStage{}).
			Where("entity_type = ? AND entity_id = ?", entityType, entityId).
			Find(&ownerStages).Error; err != nil {
			return errorx.WithStack(err)
		}
	}
	stageByName := make(map[string]*model.ProjectManagementStage, len(ownerStages))
	var pending *model.ProjectManagementStage
	for _, stage := range ownerStages {
		if stageByName[stage.Name] != nil {
			return errorx.New("project has duplicate stage name")
		}
		stageByName[stage.Name] = stage
		if stage.Phase == constants.ProjectManagementStagePhasePending {
			if pending != nil {
				return errorx.New("project has multiple pending stages")
			}
			pending = stage
		}
	}
	for _, sourceStage := range stages {
		if entityType == constants.ProjectManagementOwnerTypeProject {
			if sourceStage.Phase == constants.ProjectManagementStagePhasePending && pending != nil {
				stageIDs[sourceStage.Id] = pending.Id
				continue
			}
			if existing := stageByName[sourceStage.Name]; existing != nil {
				if existing.Phase != sourceStage.Phase {
					return errorx.New("project stage name has conflicting phase")
				}
				stageIDs[sourceStage.Id] = existing.Id
				continue
			}
		}
		stage := &model.ProjectManagementStage{EntityType: entityType, EntityId: entityId, Name: sourceStage.Name, Description: sourceStage.Description, Phase: sourceStage.Phase, Color: sourceStage.Color, SortOrder: sourceStage.SortOrder}
		if err = tx.WithContext(ctx).Create(stage).Error; err != nil {
			return errorx.WithStack(err)
		}
		stageIDs[sourceStage.Id] = stage.Id
		stageByName[stage.Name] = stage
		if stage.Phase == constants.ProjectManagementStagePhasePending {
			pending = stage
		}
	}
	for _, sourceLink := range links {
		link := &model.ProjectManagementWorkflowStage{WorkflowId: workflow.Id, StageId: stageIDs[sourceLink.StageId]}
		if link.StageId == 0 {
			return errorx.New("workflow source link stage is invalid")
		}
		if err = tx.WithContext(ctx).Create(link).Error; err != nil {
			return errorx.WithStack(err)
		}
	}
	for _, sourceTransition := range transitions {
		transition := &model.ProjectManagementWorkflowTransition{WorkflowId: workflow.Id, Name: sourceTransition.Name, FromStageId: stageIDs[sourceTransition.FromStageId], ToStageId: stageIDs[sourceTransition.ToStageId], AllowedActorTypes: sourceTransition.AllowedActorTypes}
		if transition.FromStageId == 0 || transition.ToStageId == 0 {
			return errorx.New("workflow source transition endpoint is invalid")
		}
		if err = tx.WithContext(ctx).Create(transition).Error; err != nil {
			return errorx.WithStack(err)
		}
	}
	return nil
}

func (r *ProjectManagementWorkflowRepo) createDefaultProjectWorkflowWithTx(ctx context.Context, tx *gorm.DB, projectId int64) error {
	return r.createDefaultWorkflowWithTx(ctx, tx, constants.ProjectManagementOwnerTypeProject, projectId)
}

func (r *ProjectManagementWorkflowRepo) CreateDefaultTemplateWorkflowWithTx(ctx context.Context, tx *gorm.DB, templateId int64) error {
	return r.createDefaultWorkflowWithTx(ctx, tx, constants.ProjectManagementOwnerTypeTemplate, templateId)
}

func (r *ProjectManagementWorkflowRepo) createDefaultWorkflowWithTx(ctx context.Context, tx *gorm.DB, entityType constants.ProjectManagementOwnerType, entityId int64) error {
	workflow := &model.ProjectManagementWorkflow{EntityType: entityType, EntityId: entityId, Name: "Default", IsDefault: true}
	if err := tx.WithContext(ctx).Create(workflow).Error; err != nil {
		return errorx.WithStack(err)
	}
	pending := &model.ProjectManagementStage{EntityType: entityType, EntityId: entityId, Name: "Pending", Phase: constants.ProjectManagementStagePhasePending, Color: "#3b82f6", SortOrder: 10}
	inProgress := &model.ProjectManagementStage{EntityType: entityType, EntityId: entityId, Name: "In progress", Phase: constants.ProjectManagementStagePhaseInProgress, Color: "#f59e0b", SortOrder: 20}
	finished := &model.ProjectManagementStage{EntityType: entityType, EntityId: entityId, Name: "Finished", Phase: constants.ProjectManagementStagePhaseFinished, Color: "#22c55e", SortOrder: 30}
	if err := tx.WithContext(ctx).Create([]*model.ProjectManagementStage{pending, inProgress, finished}).Error; err != nil {
		return errorx.WithStack(err)
	}
	from := &model.ProjectManagementWorkflowStage{WorkflowId: workflow.Id, StageId: pending.Id}
	middle := &model.ProjectManagementWorkflowStage{WorkflowId: workflow.Id, StageId: inProgress.Id}
	to := &model.ProjectManagementWorkflowStage{WorkflowId: workflow.Id, StageId: finished.Id}
	if err := tx.WithContext(ctx).Create([]*model.ProjectManagementWorkflowStage{from, middle, to}).Error; err != nil {
		return errorx.WithStack(err)
	}
	roles := []constants.ProjectManagementMemberRole{constants.ProjectManagementMemberRoleProjectOwner, constants.ProjectManagementMemberRoleProjectAdmin, constants.ProjectManagementMemberRoleTaskCreator, constants.ProjectManagementMemberRoleTaskExecutor}
	return errorx.WithStack(tx.WithContext(ctx).Create([]*model.ProjectManagementWorkflowTransition{
		{WorkflowId: workflow.Id, Name: "Start", FromStageId: pending.Id, ToStageId: inProgress.Id, AllowedActorTypes: roles},
		{WorkflowId: workflow.Id, Name: "Finish", FromStageId: inProgress.Id, ToStageId: finished.Id, AllowedActorTypes: roles},
	}).Error)
}
