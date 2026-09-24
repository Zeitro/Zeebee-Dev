package logic

import (
	"context"
	"fmt"
	"maps"
	"slices"
	"strings"

	"gorm.io/gorm"

	"zeebee-backend/common/errorx"
	"zeebee-backend/common/gormx"
	"zeebee-backend/common/utils"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
)

func mapProjectStageIDs(source, target []*model.ProjectManagementStage) (map[int64]int64, error) {
	targetByKey := make(map[string]int64, len(target))
	for _, stage := range target {
		key := fmt.Sprintf("%s\x00%s", stage.Name, stage.Phase)
		if _, exists := targetByKey[key]; exists {
			return nil, errorx.New("target stage mapping is ambiguous")
		}
		targetByKey[key] = stage.Id
	}
	mapping := make(map[int64]int64, len(source))
	for _, stage := range source {
		id, ok := targetByKey[fmt.Sprintf("%s\x00%s", stage.Name, stage.Phase)]
		if !ok {
			return nil, errorx.New("target stage mapping is missing")
		}
		mapping[stage.Id] = id
	}
	return mapping, nil
}

func (l *ProjectManagementWorkflowLogic) ProjectWorkflowPage(ctx context.Context, claim utils.JwtUser, projectId, page, limit int64) (gormx.PageResponse[*model.ProjectManagementWorkflow], error) {
	if _, err := l.projectManagementLogic().RequireProject(ctx, projectId, claim.OrganizationId); err != nil {
		return gormx.PageResponse[*model.ProjectManagementWorkflow]{}, err
	}
	if _, err := l.projectManagementLogic().RequireMember(ctx, projectId, claim.OrganizationUserId); err != nil {
		return gormx.PageResponse[*model.ProjectManagementWorkflow]{}, err
	}
	return l.repos.ProjectManagementWorkflowRepo.PageProjectWorkflows(ctx, projectId, claim.OrganizationId, gormx.PageRequest{Page: page, Limit: limit})
}

func (l *ProjectManagementWorkflowLogic) ProjectWorkflowDetail(ctx context.Context, claim utils.JwtUser, projectId, workflowId int64) (*model.ProjectManagementWorkflow, error) {
	if _, err := l.projectManagementLogic().RequireMember(ctx, projectId, claim.OrganizationUserId); err != nil {
		return nil, err
	}
	workflow, err := l.repos.ProjectManagementWorkflowRepo.FindProjectWorkflowByID(ctx, workflowId, projectId, claim.OrganizationId)
	if err != nil {
		return nil, projectManagementNotFound(err, "project workflow not found")
	}
	return workflow, nil
}

func (l *ProjectManagementWorkflowLogic) ProjectWorkflowBatch(ctx context.Context, claim utils.JwtUser, projectId int64, sourceWorkflowIds []int64) (created []*model.ProjectManagementWorkflow, err error) {
	if len(sourceWorkflowIds) == 0 {
		return nil, errorx.New("source workflow ids are required")
	}
	err = l.repos.ProjectManagementWorkflowRepo.Transaction(ctx, func(tx *gorm.DB) error {
		if _, innerErr := l.projectManagementLogic().requireProjectManagerWithTx(ctx, tx, claim, projectId); innerErr != nil {
			return innerErr
		}
		var existing []*model.ProjectManagementWorkflow
		if innerErr := tx.WithContext(ctx).Model(&model.ProjectManagementWorkflow{}).
			Select("id, created_at, updated_at, is_del, entity_type, entity_id, source_workflow_id, source_workflow_name, name, description, is_default").
			Where("entity_type = ? AND entity_id = ?", constants.ProjectManagementOwnerTypeProject, projectId).
			Order("id ASC").Find(&existing).Error; innerErr != nil {
			return errorx.WithStack(innerErr)
		}
		seen := make(map[int64]struct{}, len(sourceWorkflowIds))
		for _, workflow := range existing {
			seen[workflow.SourceWorkflowId] = struct{}{}
		}
		for i, sourceId := range sourceWorkflowIds {
			if _, ok := seen[sourceId]; ok {
				return errorx.New("duplicate source workflow")
			}
			seen[sourceId] = struct{}{}
			source, innerErr := l.repos.ProjectManagementWorkflowRepo.FindWorkflowSourceByIDWithTx(ctx, tx, sourceId, claim.OrganizationId)
			if innerErr != nil {
				return projectManagementNotFound(innerErr, "workflow source not found")
			}
			ownerStages, innerErr := l.repos.ProjectManagementWorkflowRepo.ListWorkflowStagesWithTx(ctx, tx, source.Id)
			if innerErr != nil {
				return innerErr
			}
			if err = validateWorkflowSourceStages(ownerStages); err != nil {
				return err
			}
			sourceTransitions, innerErr := l.repos.ProjectManagementWorkflowRepo.ListStageTransitionsWithTx(ctx, tx, source.Id)
			if innerErr != nil {
				return innerErr
			}
			local := &model.ProjectManagementWorkflow{EntityType: constants.ProjectManagementOwnerTypeProject, EntityId: projectId, SourceWorkflowId: source.Id, SourceWorkflowName: source.Name, Name: source.Name, Description: source.Description, IsDefault: len(existing) == 0 && i == 0}
			if innerErr = tx.WithContext(ctx).Omit("Kind", "OrganizationId", "ProjectId", "SourceTemplateId", "SourceTemplateName").Create(local).Error; innerErr != nil {
				return errorx.WithStack(innerErr)
			}
			stageIds, innerErr := l.reconcileProjectStagesWithTx(ctx, tx, projectId, ownerStages)
			if innerErr != nil {
				return innerErr
			}
			local.Stages = make([]*model.ProjectManagementStage, 0, len(ownerStages))
			for _, sourceStage := range ownerStages {
				local.Stages = append(local.Stages, &model.ProjectManagementStage{Base: model.Base{Id: stageIds[sourceStage.Id]}})
			}
			for _, sourceStage := range ownerStages {
				link := &model.ProjectManagementWorkflowStage{WorkflowId: local.Id, StageId: stageIds[sourceStage.Id]}
				if innerErr = tx.WithContext(ctx).Create(link).Error; innerErr != nil {
					return errorx.WithStack(innerErr)
				}
			}
			local.Transitions = make([]*model.ProjectManagementWorkflowTransition, 0, len(sourceTransitions))
			for _, sourceTransition := range sourceTransitions {
				fromStageId, toStageId := stageIds[sourceTransition.FromStageId], stageIds[sourceTransition.ToStageId]
				if fromStageId == 0 || toStageId == 0 {
					return errorx.New("workflow source transition endpoint is invalid")
				}
				transition := &model.ProjectManagementWorkflowTransition{WorkflowId: local.Id, Name: sourceTransition.Name, FromStageId: fromStageId, ToStageId: toStageId, AllowedActorTypes: sourceTransition.AllowedActorTypes}
				if innerErr = tx.WithContext(ctx).Create(transition).Error; innerErr != nil {
					return errorx.WithStack(innerErr)
				}
				local.Transitions = append(local.Transitions, transition)
			}
			created = append(created, local)
		}
		if err := l.validateProjectWorkflowInvariantWithTx(ctx, tx, projectId); err != nil {
			return err
		}
		return l.recordProjectWorkflowOperation(ctx, tx, claim, projectId, 0, 0, constants.OperationActionProjectManagementWorkflowCreate, nil, map[string]any{"workflow_ids": workflowIDs(created)})
	})
	return created, err
}

func (l *ProjectManagementWorkflowLogic) reconcileProjectStagesWithTx(ctx context.Context, tx *gorm.DB, projectId int64, source []*model.ProjectManagementStage) (map[int64]int64, error) {
	stages := make([]*model.ProjectManagementStage, 0)
	if err := tx.WithContext(ctx).Model(&model.ProjectManagementStage{}).Where("entity_type = ? AND entity_id = ?", constants.ProjectManagementOwnerTypeProject, projectId).Find(&stages).Error; err != nil {
		return nil, errorx.WithStack(err)
	}
	byName := make(map[string]*model.ProjectManagementStage, len(stages))
	var pending *model.ProjectManagementStage
	for _, stage := range stages {
		if byName[stage.Name] != nil {
			return nil, errorx.New("project has duplicate stage name")
		}
		byName[stage.Name] = stage
		if stage.Phase == constants.ProjectManagementStagePhasePending {
			if pending != nil {
				return nil, errorx.New("project has multiple pending stages")
			}
			pending = stage
		}
	}
	mapping := make(map[int64]int64, len(source))
	for _, stage := range source {
		if stage.Phase == constants.ProjectManagementStagePhasePending && pending != nil {
			mapping[stage.Id] = pending.Id
			continue
		}
		if existing := byName[stage.Name]; existing != nil {
			if existing.Phase != stage.Phase {
				return nil, errorx.New("project stage name has conflicting phase")
			}
			mapping[stage.Id] = existing.Id
			continue
		}
		copy := &model.ProjectManagementStage{EntityType: constants.ProjectManagementOwnerTypeProject, EntityId: projectId, Name: stage.Name, Description: stage.Description, Phase: stage.Phase, Color: stage.Color, SortOrder: stage.SortOrder}
		if err := tx.WithContext(ctx).Create(copy).Error; err != nil {
			return nil, errorx.WithStack(err)
		}
		mapping[stage.Id] = copy.Id
		byName[copy.Name] = copy
		if copy.Phase == constants.ProjectManagementStagePhasePending {
			pending = copy
		}
	}
	return mapping, nil
}

func (l *ProjectManagementWorkflowLogic) validateProjectWorkflowInvariantWithTx(ctx context.Context, tx *gorm.DB, projectId int64) error {
	workflows := make([]*model.ProjectManagementWorkflow, 0)
	stages := make([]*model.ProjectManagementStage, 0)
	links := make([]*model.ProjectManagementWorkflowStage, 0)
	if err := tx.WithContext(ctx).Model(&model.ProjectManagementWorkflow{}).Where("entity_type = ? AND entity_id = ?", constants.ProjectManagementOwnerTypeProject, projectId).Find(&workflows).Error; err != nil {
		return errorx.WithStack(err)
	}
	if err := tx.WithContext(ctx).Model(&model.ProjectManagementStage{}).Where("entity_type = ? AND entity_id = ?", constants.ProjectManagementOwnerTypeProject, projectId).Find(&stages).Error; err != nil {
		return errorx.WithStack(err)
	}
	if err := tx.WithContext(ctx).Model(&model.ProjectManagementWorkflowStage{}).Where("workflow_id IN ?", workflowIDs(workflows)).Find(&links).Error; err != nil {
		return errorx.WithStack(err)
	}
	return validateProjectWorkflowInvariant(projectId, workflows, stages, links)
}

func validateWorkflowSourceStages(stages []*model.ProjectManagementStage) error {
	pending, finished := 0, 0
	for i, stage := range stages {
		if i == 0 && stage.Phase != constants.ProjectManagementStagePhasePending {
			return errorx.New("first workflow stage must be pending")
		}
		switch stage.Phase {
		case constants.ProjectManagementStagePhasePending:
			pending++
		case constants.ProjectManagementStagePhaseFinished:
			finished++
		}
	}
	if pending != 1 || finished == 0 {
		return errorx.New("workflow source requires exactly one pending stage and at least one finished stage")
	}
	if stages[len(stages)-1].Phase != constants.ProjectManagementStagePhaseFinished {
		return errorx.New("last workflow stage must be finished")
	}
	return nil
}

func (l *ProjectManagementWorkflowLogic) ProjectWorkflowCreate(ctx context.Context, claim utils.JwtUser, projectId int64, input dto.ProjectManagementProjectWorkflowCreateDto) (*model.ProjectManagementWorkflow, error) {
	if strings.TrimSpace(input.Name) == "" {
		return nil, errorx.New("workflow name is required")
	}
	var workflow *model.ProjectManagementWorkflow
	err := l.repos.ProjectManagementWorkflowRepo.Transaction(ctx, func(tx *gorm.DB) error {
		if _, err := l.projectManagementLogic().requireProjectManagerWithTx(ctx, tx, claim, projectId); err != nil {
			return err
		}
		stages, err := l.repos.ProjectManagementStageRepo.ListByOwnerWithTx(ctx, tx, constants.ProjectManagementOwnerTypeProject, projectId)
		if err != nil {
			return err
		}
		if len(stages) < 2 {
			return errorx.New("project workflow requires at least two stages")
		}
		workflows, err := l.repos.ProjectManagementWorkflowRepo.ListProjectWorkflowsWithTx(ctx, tx, projectId)
		if err != nil {
			return err
		}
		workflow = &model.ProjectManagementWorkflow{EntityType: constants.ProjectManagementOwnerTypeProject, EntityId: projectId, Name: input.Name, Description: input.Description, IsDefault: len(workflows) == 0}
		if err = tx.WithContext(ctx).Create(workflow).Error; err != nil {
			return errorx.WithStack(err)
		}
		if err = l.replaceProjectWorkflowTransitionsWithTx(ctx, tx, projectId, workflow.Id, projectWorkflowTransitions(stages)); err != nil {
			return err
		}
		return l.recordProjectWorkflowOperation(ctx, tx, claim, projectId, workflow.Id, 0, constants.OperationActionProjectManagementWorkflowCreate, nil, workflowOperationValues(workflow))
	})
	if err != nil {
		return nil, err
	}
	return l.ProjectWorkflowDetail(ctx, claim, projectId, workflow.Id)
}

func projectWorkflowTransitions(stages []*model.ProjectManagementStage) []dto.ProjectManagementWorkflowTransitionDto {
	transitions := make([]dto.ProjectManagementWorkflowTransitionDto, 0, len(stages)*(len(stages)-1))
	for _, from := range stages {
		for _, to := range stages {
			if from.Id == to.Id {
				continue
			}
			transitions = append(transitions, dto.ProjectManagementWorkflowTransitionDto{
				Name:              fmt.Sprintf("%s → %s", from.Name, to.Name),
				FromStageId:       from.Id,
				ToStageId:         to.Id,
				AllowedActorTypes: defaultWorkflowActorTypes(),
			})
		}
	}
	return transitions
}

func (l *ProjectManagementWorkflowLogic) ProjectWorkflowUpdate(ctx context.Context, claim utils.JwtUser, projectId, workflowId int64, input dto.ProjectManagementProjectWorkflowUpdateDto) (*model.ProjectManagementWorkflow, error) {
	var workflow *model.ProjectManagementWorkflow
	err := l.repos.ProjectManagementWorkflowRepo.Transaction(ctx, func(tx *gorm.DB) error {
		if _, err := l.projectManagementLogic().requireProjectManagerWithTx(ctx, tx, claim, projectId); err != nil {
			return err
		}
		var err error
		workflow, err = l.repos.ProjectManagementWorkflowRepo.FindProjectWorkflowByIDWithTx(ctx, tx, workflowId, projectId, claim.OrganizationId)
		if err != nil {
			return projectManagementNotFound(err, "project workflow not found")
		}
		if err = l.replaceProjectWorkflowTransitionsWithTx(ctx, tx, projectId, workflowId, input.Transitions); err != nil {
			return err
		}
		return l.recordProjectWorkflowOperation(ctx, tx, claim, projectId, workflowId, 0, constants.OperationActionProjectManagementWorkflowStageUpdate, nil, map[string]any{"stage_ids": projectWorkflowStageIDs(input.Transitions)})
	})
	if err != nil {
		return nil, err
	}
	return l.ProjectWorkflowDetail(ctx, claim, projectId, workflow.Id)
}

func (l *ProjectManagementWorkflowLogic) replaceProjectWorkflowTransitionsWithTx(ctx context.Context, tx *gorm.DB, projectId, workflowId int64, transitions []dto.ProjectManagementWorkflowTransitionDto) error {
	stageIDs, err := validateProjectWorkflowTransitions(transitions)
	if err != nil {
		return err
	}
	stages, err := l.repos.ProjectManagementStageRepo.ListByOwnerWithTx(ctx, tx, constants.ProjectManagementOwnerTypeProject, projectId)
	if err != nil {
		return err
	}
	bound := make([]*model.ProjectManagementStage, 0, len(stageIDs))
	for _, stage := range stages {
		if _, ok := stageIDs[stage.Id]; ok {
			bound = append(bound, stage)
			delete(stageIDs, stage.Id)
		}
	}
	if len(stageIDs) != 0 {
		return errorx.New("workflow stage owner mismatch")
	}
	if err = validateWorkflowSourceStages(bound); err != nil {
		return err
	}
	links, err := l.repos.ProjectManagementWorkflowRepo.ListStageLinksWithTx(ctx, tx, workflowId)
	if err != nil {
		return err
	}
	boundIDs := make(map[int64]struct{}, len(bound))
	for _, stage := range bound {
		boundIDs[stage.Id] = struct{}{}
	}
	for _, link := range links {
		if _, ok := boundIDs[link.StageId]; ok {
			delete(boundIDs, link.StageId)
			continue
		}
		var tasks int64
		if err = tx.WithContext(ctx).Model(&model.ProjectManagementTask{}).
			Where("workflow_id = ? AND entity_type = ? AND entity_id = ? AND stage_id = ? AND (phase IS NULL OR phase <> ?)", workflowId, constants.ProjectManagementOwnerTypeProject, projectId, link.StageId, constants.ProjectManagementStagePhaseFinished).
			Count(&tasks).Error; err != nil {
			return errorx.WithStack(err)
		}
		if tasks > 0 {
			return errorx.New("workflow stage with unfinished tasks cannot be removed")
		}
		if err = tx.WithContext(ctx).Delete(&model.ProjectManagementWorkflowStage{}, link.Id).Error; err != nil {
			return errorx.WithStack(err)
		}
	}
	for stageID := range boundIDs {
		if err = tx.WithContext(ctx).Create(&model.ProjectManagementWorkflowStage{WorkflowId: workflowId, StageId: stageID}).Error; err != nil {
			return errorx.WithStack(err)
		}
	}
	if err = tx.WithContext(ctx).Where("workflow_id = ?", workflowId).Delete(&model.ProjectManagementWorkflowTransition{}).Error; err != nil {
		return errorx.WithStack(err)
	}
	for _, transition := range transitions {
		if err = tx.WithContext(ctx).Create(&model.ProjectManagementWorkflowTransition{WorkflowId: workflowId, Name: transition.Name, FromStageId: transition.FromStageId, ToStageId: transition.ToStageId, AllowedActorTypes: transition.AllowedActorTypes}).Error; err != nil {
			return errorx.WithStack(err)
		}
	}
	return nil
}

func validateProjectWorkflowTransitions(transitions []dto.ProjectManagementWorkflowTransitionDto) (map[int64]struct{}, error) {
	if len(transitions) == 0 {
		return nil, errorx.New("workflow transitions are required")
	}
	stageIDs := make(map[int64]struct{}, len(transitions)*2)
	edges := make(map[[2]int64]struct{}, len(transitions))
	for _, transition := range transitions {
		if err := validateProjectManagementWorkflowTransition(transition); err != nil {
			return nil, err
		}
		edge := [2]int64{transition.FromStageId, transition.ToStageId}
		if _, exists := edges[edge]; exists {
			return nil, errorx.New("workflow transition must be unique")
		}
		edges[edge] = struct{}{}
		stageIDs[transition.FromStageId] = struct{}{}
		stageIDs[transition.ToStageId] = struct{}{}
	}
	return stageIDs, nil
}

func projectWorkflowStageIDs(transitions []dto.ProjectManagementWorkflowTransitionDto) []int64 {
	ids, err := validateProjectWorkflowTransitions(transitions)
	if err != nil {
		return nil
	}
	return slices.Sorted(maps.Keys(ids))
}

func (l *ProjectManagementWorkflowLogic) ProjectWorkflowDelete(ctx context.Context, claim utils.JwtUser, projectId, workflowId int64) error {
	return l.repos.ProjectManagementWorkflowRepo.Transaction(ctx, func(tx *gorm.DB) error {
		if _, err := l.projectManagementLogic().requireProjectManagerWithTx(ctx, tx, claim, projectId); err != nil {
			return err
		}
		workflow, err := l.repos.ProjectManagementWorkflowRepo.FindProjectWorkflowByIDWithTx(ctx, tx, workflowId, projectId, claim.OrganizationId)
		if err != nil {
			return projectManagementNotFound(err, "project workflow not found")
		}
		if workflow.IsDefault {
			return errorx.New("default workflow cannot be deleted")
		}
		tasks, err := l.repos.ProjectManagementTaskRepo.ListByWorkflowWithTx(ctx, tx, workflowId, claim.OrganizationId, false)
		if err != nil {
			return err
		}
		if len(tasks) > 0 {
			return errorx.New("workflow with unfinished tasks cannot be deleted")
		}
		if err = tx.WithContext(ctx).Where("workflow_id = ?", workflowId).Delete(&model.ProjectManagementWorkflowStage{}).Error; err != nil {
			return errorx.WithStack(err)
		}
		if err = tx.WithContext(ctx).Where("workflow_id = ?", workflowId).Delete(&model.ProjectManagementWorkflowTransition{}).Error; err != nil {
			return errorx.WithStack(err)
		}
		if err = tx.WithContext(ctx).Where("id = ?", workflowId).Delete(&model.ProjectManagementWorkflow{}).Error; err != nil {
			return errorx.WithStack(err)
		}
		return l.recordProjectWorkflowOperation(ctx, tx, claim, projectId, workflowId, 0, constants.OperationActionProjectManagementWorkflowDelete, workflowOperationValues(workflow), nil)
	})
}

func (l *ProjectManagementWorkflowLogic) ProjectWorkflowSetDefault(ctx context.Context, claim utils.JwtUser, projectId, workflowId int64) error {
	return l.repos.ProjectManagementWorkflowRepo.Transaction(ctx, func(tx *gorm.DB) error {
		if _, err := l.projectManagementLogic().requireProjectManagerWithTx(ctx, tx, claim, projectId); err != nil {
			return err
		}
		workflow, err := l.repos.ProjectManagementWorkflowRepo.FindProjectWorkflowByIDWithTx(ctx, tx, workflowId, projectId, claim.OrganizationId)
		if err != nil {
			return projectManagementNotFound(err, "project workflow not found")
		}
		if workflow.IsDefault {
			return nil
		}
		if err = tx.WithContext(ctx).Model(&model.ProjectManagementWorkflow{}).
			Where("entity_type = ? AND entity_id = ?", constants.ProjectManagementOwnerTypeProject, projectId).
			Update("is_default", false).Error; err != nil {
			return errorx.WithStack(err)
		}
		if err = tx.WithContext(ctx).Model(&model.ProjectManagementWorkflow{}).Where("id = ?", workflowId).Update("is_default", true).Error; err != nil {
			return errorx.WithStack(err)
		}
		return l.recordProjectWorkflowOperation(ctx, tx, claim, projectId, workflowId, 0, constants.OperationActionProjectManagementWorkflowDefault, map[string]any{"is_default": false}, map[string]any{"is_default": true})
	})
}

func (l *ProjectManagementWorkflowLogic) ProjectWorkflowStageUpdate(ctx context.Context, claim utils.JwtUser, projectId, workflowId int64, input dto.ProjectManagementWorkflowUpdateDto) (*model.ProjectManagementWorkflow, error) {
	if len(input.Stages) == 0 {
		return nil, errorx.New("workflow stages are required")
	}
	if err := validateProjectWorkflowStageUpdate(input.Stages); err != nil {
		return nil, err
	}
	var result *model.ProjectManagementWorkflow
	err := l.repos.ProjectManagementWorkflowRepo.Transaction(ctx, func(tx *gorm.DB) error {
		if _, err := l.projectManagementLogic().requireProjectManagerWithTx(ctx, tx, claim, projectId); err != nil {
			return err
		}
		workflow, err := l.repos.ProjectManagementWorkflowRepo.FindProjectWorkflowByIDWithTx(ctx, tx, workflowId, projectId, claim.OrganizationId)
		if err != nil {
			return projectManagementNotFound(err, "project workflow not found")
		}
		existing, err := l.repos.ProjectManagementWorkflowRepo.ListWorkflowStagesWithTx(ctx, tx, workflowId)
		if err != nil {
			return err
		}
		stageByID := make(map[int64]*model.ProjectManagementStage, len(existing))
		for _, stage := range existing {
			stageByID[stage.Id] = stage
		}
		requested := make(map[int64]struct{}, len(input.Stages))
		for _, item := range input.Stages {
			if item.Id == 0 {
				stage := &model.ProjectManagementStage{EntityType: constants.ProjectManagementOwnerTypeProject, EntityId: projectId, Name: item.Name, Description: item.Description, Phase: item.Phase, Color: item.Color, SortOrder: item.SortOrder}
				if err = tx.WithContext(ctx).Create(stage).Error; err != nil {
					return errorx.WithStack(err)
				}
				item.Id = stage.Id
				stageByID[stage.Id] = stage
			}
			stage, ok := stageByID[item.Id]
			if !ok || stage.EntityType != constants.ProjectManagementOwnerTypeProject || stage.EntityId != projectId {
				return errorx.New("workflow stage owner mismatch")
			}
			if err = tx.WithContext(ctx).Model(&model.ProjectManagementStage{}).Where("id = ?", stage.Id).Updates(map[string]any{"name": item.Name, "description": item.Description, "phase": item.Phase, "color": item.Color, "sort_order": item.SortOrder}).Error; err != nil {
				return errorx.WithStack(err)
			}
			if stage.Phase != item.Phase {
				if err = tx.WithContext(ctx).Model(&model.ProjectManagementTask{}).Where("entity_type = ? AND entity_id = ? AND stage_id = ?", constants.ProjectManagementOwnerTypeProject, projectId, stage.Id).Update("phase", item.Phase).Error; err != nil {
					return errorx.WithStack(err)
				}
			}
			requested[item.Id] = struct{}{}
		}
		linkRows, err := l.repos.ProjectManagementWorkflowRepo.ListStageLinksWithTx(ctx, tx, workflowId)
		if err != nil {
			return err
		}
		linkByStage := make(map[int64]*model.ProjectManagementWorkflowStage, len(input.Stages))
		for _, link := range linkRows {
			if _, ok := requested[link.StageId]; !ok {
				var activeTasks int64
				if err = tx.WithContext(ctx).Model(&model.ProjectManagementTask{}).
					Where("workflow_id = ? AND entity_type = ? AND entity_id = ? AND stage_id = ? AND phase <> ?", workflowId, constants.ProjectManagementOwnerTypeProject, projectId, link.StageId, constants.ProjectManagementStagePhaseFinished).
					Count(&activeTasks).Error; err != nil {
					return errorx.WithStack(err)
				}
				if activeTasks > 0 {
					return errorx.New("workflow stage with unfinished tasks cannot be removed")
				}
				if err = tx.WithContext(ctx).Delete(&model.ProjectManagementWorkflowStage{}, link.Id).Error; err != nil {
					return errorx.WithStack(err)
				}
				continue
			}
			linkByStage[link.StageId] = link
		}
		for stageId := range requested {
			if linkByStage[stageId] != nil {
				continue
			}
			link := &model.ProjectManagementWorkflowStage{WorkflowId: workflowId, StageId: stageId}
			if err = tx.WithContext(ctx).Create(link).Error; err != nil {
				return errorx.WithStack(err)
			}
			linkByStage[stageId] = link
		}
		if err = tx.WithContext(ctx).Where("workflow_id = ?", workflowId).Delete(&model.ProjectManagementWorkflowTransition{}).Error; err != nil {
			return errorx.WithStack(err)
		}
		for _, item := range input.Stages {
			for _, transition := range item.Transitions {
				_, fromOK := linkByStage[transition.FromStageId]
				_, toOK := linkByStage[transition.ToStageId]
				if !fromOK || !toOK {
					return errorx.New("workflow transition endpoint is invalid")
				}
				if err = tx.WithContext(ctx).Create(&model.ProjectManagementWorkflowTransition{WorkflowId: workflowId, Name: transition.Name, FromStageId: transition.FromStageId, ToStageId: transition.ToStageId, AllowedActorTypes: transition.AllowedActorTypes}).Error; err != nil {
					return errorx.WithStack(err)
				}
			}
		}
		result = workflow
		return l.recordProjectWorkflowOperation(ctx, tx, claim, projectId, workflowId, 0, constants.OperationActionProjectManagementWorkflowStageUpdate, nil, map[string]any{"stage_ids": workflowStageIDs(input.Stages)})
	})
	if err != nil {
		return nil, err
	}
	return l.ProjectWorkflowDetail(ctx, claim, projectId, result.Id)
}

func (l *ProjectManagementWorkflowLogic) recordProjectWorkflowOperation(ctx context.Context, tx *gorm.DB, claim utils.JwtUser, projectId, workflowId, stageId int64, action constants.OperationAction, before, after map[string]any) error {
	metadata, err := marshalProjectManagementOperationMetadata(projectManagementOperationMetadata{ProjectID: projectId, WorkflowID: workflowId, StageID: stageId}, before, after)
	if err != nil {
		return errorx.WithStack(err)
	}
	return NewOperationLogic(l.BaseLogic).CreateOperationLog(ctx, tx, dto.OperationCreateDto{EntityId: projectId, EntityType: constants.OperationEntityTypeProjectManagementProject, OperatorId: claim.OrganizationUserId, OperatorName: claim.Username, OperatorType: constants.UserTypeBorrower, Operation: action, Metadata: metadata, OrganizationId: claim.OrganizationId})
}

func workflowStageIDs(stages []dto.ProjectManagementWorkflowUpdateStageDto) []int64 {
	ids := make([]int64, 0, len(stages))
	for _, stage := range stages {
		ids = append(ids, stage.Id)
	}
	return ids
}

func validateProjectWorkflowStageUpdate(stages []dto.ProjectManagementWorkflowUpdateStageDto) error {
	ids := make(map[int64]struct{}, len(stages))
	pending, finished := 0, 0
	minSortOrder := -1
	for _, stage := range stages {
		if strings.TrimSpace(stage.Name) == "" {
			return errorx.New("workflow stage name is required")
		}
		if !isProjectManagementStagePhase(stage.Phase) {
			return errorx.New("invalid workflow stage phase")
		}
		if stage.SortOrder < 0 {
			return errorx.New("workflow stage sort order must be >= 0")
		}
		if minSortOrder < 0 || stage.SortOrder < minSortOrder {
			minSortOrder = stage.SortOrder
		}
		if stage.Id != 0 {
			if _, exists := ids[stage.Id]; exists {
				return errorx.New("workflow stage id must be unique")
			}
			ids[stage.Id] = struct{}{}
		} else if len(stage.Transitions) != 0 {
			return errorx.New("new workflow stage cannot include transitions")
		}
		switch stage.Phase {
		case constants.ProjectManagementStagePhasePending:
			pending++
		case constants.ProjectManagementStagePhaseFinished:
			finished++
		}
		for _, transition := range stage.Transitions {
			if transition.FromStageId != stage.Id || transition.ToStageId == 0 || transition.FromStageId == transition.ToStageId {
				return errorx.New("workflow transition endpoint is invalid")
			}
			if _, ok := ids[transition.ToStageId]; !ok {
				// A later stage may declare the target; validate after collecting all IDs.
				continue
			}
		}
	}
	if pending != 1 {
		return errorx.New("workflow must have exactly one pending stage")
	}
	for _, stage := range stages {
		if stage.Phase == constants.ProjectManagementStagePhasePending && stage.SortOrder != minSortOrder {
			return errorx.New("first workflow stage must be pending")
		}
	}
	if finished == 0 {
		return errorx.New("workflow must have at least one finished stage")
	}
	maxOrder := stages[0].SortOrder
	for _, stage := range stages[1:] {
		maxOrder = max(maxOrder, stage.SortOrder)
	}
	for _, stage := range stages {
		if stage.SortOrder == maxOrder && stage.Phase != constants.ProjectManagementStagePhaseFinished {
			return errorx.New("last workflow stage must be finished")
		}
	}
	for _, stage := range stages {
		for _, transition := range stage.Transitions {
			if transition.ToStageId == 0 {
				return errorx.New("workflow transition endpoint is invalid")
			}
			if _, ok := ids[transition.ToStageId]; !ok {
				return errorx.New("workflow transition endpoint is invalid")
			}
		}
	}
	return nil
}

func (l *ProjectManagementWorkflowLogic) ProjectDetail(ctx context.Context, claim utils.JwtUser, projectId int64) (*model.ProjectManagementWorkflow, error) {
	workflow, err := l.repos.ProjectManagementWorkflowRepo.FindDefaultProjectWorkflowWithTx(ctx, l.repos.ProjectManagementWorkflowRepo.Client(ctx), projectId, claim.OrganizationId)
	if err != nil {
		return nil, projectManagementNotFound(err, "project workflow not found")
	}
	return l.ProjectWorkflowDetail(ctx, claim, projectId, workflow.Id)
}

func (l *ProjectManagementWorkflowLogic) ProjectUpdate(ctx context.Context, claim utils.JwtUser, projectId int64, input dto.ProjectManagementWorkflowUpdateDto) (workflow *model.ProjectManagementWorkflow, err error) {
	workflow, err = l.repos.ProjectManagementWorkflowRepo.FindDefaultProjectWorkflowWithTx(ctx, l.repos.ProjectManagementWorkflowRepo.Client(ctx), projectId, claim.OrganizationId)
	if err != nil {
		return nil, projectManagementNotFound(err, "project workflow not found")
	}
	return l.ProjectWorkflowStageUpdate(ctx, claim, projectId, workflow.Id, input)
}
