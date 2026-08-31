package logic

import (
	"context"

	"gorm.io/gorm"

	"zeebee-backend/common/errorx"
	"zeebee-backend/common/utils"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
)

func (l *ProjectManagementWorkflowLogic) ProjectDetail(ctx context.Context, claim utils.JwtUser, projectId int64) (*model.ProjectManagementWorkflow, error) {
	if _, err := l.projectManagementLogic().RequireProject(ctx, projectId, claim.OrganizationId); err != nil {
		return nil, err
	}
	if _, err := l.projectManagementLogic().RequireMember(ctx, projectId, claim.OrganizationUserId); err != nil {
		return nil, err
	}
	workflow, err := l.repos.ProjectManagementWorkflowRepo.FindProjectDefinition(ctx, projectId, claim.OrganizationId)
	if err != nil {
		return nil, projectManagementNotFound(err, "project workflow not found")
	}
	return workflow, nil
}

func (l *ProjectManagementWorkflowLogic) ProjectUpdate(ctx context.Context, claim utils.JwtUser, projectId int64, input dto.ProjectManagementWorkflowUpdateDto) (workflow *model.ProjectManagementWorkflow, err error) {
	definition, err := projectManagementWorkflowUpdateDefinition(input)
	if err != nil {
		return nil, err
	}
	err = l.repos.ProjectManagementWorkflowRepo.Transaction(ctx, func(tx *gorm.DB) error {
		var innerErr error
		if _, innerErr = l.projectManagementLogic().requireProjectManagerWithTx(ctx, tx, claim, projectId); innerErr != nil {
			return innerErr
		}
		current, innerErr := l.repos.ProjectManagementWorkflowRepo.LockProjectWithTx(ctx, tx, projectId, claim.OrganizationId)
		if innerErr != nil {
			return projectManagementNotFound(innerErr, "project workflow not found")
		}
		if innerErr = l.assignNewStageSortOrders(ctx, tx, current.Id, &definition); innerErr != nil {
			return innerErr
		}
		if innerErr = validateProjectManagementProjectWorkflow(definition); innerErr != nil {
			return innerErr
		}
		beforeStages, innerErr := l.repos.ProjectManagementWorkflowRepo.ListStagesWithTx(ctx, tx, current.Id)
		if innerErr != nil {
			return innerErr
		}
		beforeTransitions, innerErr := l.repos.ProjectManagementWorkflowRepo.ListTransitionsWithTx(ctx, tx, current.Id)
		if innerErr != nil {
			return innerErr
		}
		_, stages, transitions, innerErr := projectManagementTemplateModels(definition)
		if innerErr != nil {
			return innerErr
		}
		if innerErr = l.repos.ProjectManagementWorkflowRepo.ReplaceProjectDefinitionWithTx(ctx, tx, current.Id, stages, transitions); innerErr != nil {
			return innerErr
		}
		afterStages, innerErr := l.repos.ProjectManagementWorkflowRepo.ListStagesWithTx(ctx, tx, current.Id)
		if innerErr != nil {
			return innerErr
		}
		afterTransitions, innerErr := l.repos.ProjectManagementWorkflowRepo.ListTransitionsWithTx(ctx, tx, current.Id)
		if innerErr != nil {
			return innerErr
		}
		return l.recordProjectWorkflowDefinitionChanges(ctx, tx, claim, projectId, beforeStages, afterStages, beforeTransitions, afterTransitions)
	})
	if err != nil {
		return nil, err
	}
	return l.ProjectDetail(ctx, claim, projectId)
}

func (l *ProjectManagementWorkflowLogic) recordProjectWorkflowDefinitionChanges(ctx context.Context, tx *gorm.DB, claim utils.JwtUser, projectId int64, beforeStages, afterStages []*model.ProjectManagementWorkflowStage, beforeTransitions, afterTransitions []*model.ProjectManagementWorkflowTransition) error {
	record := func(action constants.OperationAction, stageId, transitionId int64, before, after map[string]any) error {
		return NewProjectManagementTaskLogic(l.BaseLogic).recordProjectTaskOperation(ctx, tx, projectManagementTaskOperation{claim: claim, projectId: projectId, stageId: stageId, transitionId: transitionId, action: action, before: before, after: after})
	}
	beforeStageByID := make(map[int64]*model.ProjectManagementWorkflowStage, len(beforeStages))
	for _, stage := range beforeStages {
		beforeStageByID[stage.Id] = stage
	}
	for _, stage := range afterStages {
		if before := beforeStageByID[stage.Id]; before != nil {
			old, next := operationValuesDiff(workflowStageOperationValues(before), workflowStageOperationValues(stage))
			if len(old) != 0 {
				if err := record(constants.OperationActionProjectManagementWorkflowStageUpdate, stage.Id, 0, old, next); err != nil {
					return err
				}
			}
			delete(beforeStageByID, stage.Id)
			continue
		}
		if err := record(constants.OperationActionProjectManagementWorkflowStageCreate, stage.Id, 0, nil, workflowStageOperationIdentity(stage)); err != nil {
			return err
		}
	}
	for _, stage := range beforeStageByID {
		if err := record(constants.OperationActionProjectManagementWorkflowStageDelete, stage.Id, 0, workflowStageOperationIdentity(stage), nil); err != nil {
			return err
		}
	}
	beforeTransitionByID := make(map[int64]*model.ProjectManagementWorkflowTransition, len(beforeTransitions))
	for _, transition := range beforeTransitions {
		beforeTransitionByID[transition.Id] = transition
	}
	for _, transition := range afterTransitions {
		if before := beforeTransitionByID[transition.Id]; before != nil {
			old, next := operationValuesDiff(workflowTransitionOperationValues(before), workflowTransitionOperationValues(transition))
			if len(old) != 0 {
				if err := record(constants.OperationActionProjectManagementWorkflowTransitionUpdate, 0, transition.Id, old, next); err != nil {
					return err
				}
			}
			delete(beforeTransitionByID, transition.Id)
			continue
		}
		if err := record(constants.OperationActionProjectManagementWorkflowTransitionCreate, 0, transition.Id, nil, workflowTransitionOperationIdentity(transition)); err != nil {
			return err
		}
	}
	for _, transition := range beforeTransitionByID {
		if err := record(constants.OperationActionProjectManagementWorkflowTransitionDelete, 0, transition.Id, workflowTransitionOperationIdentity(transition), nil); err != nil {
			return err
		}
	}
	return nil
}

func validateProjectManagementProjectWorkflow(input dto.ProjectManagementWorkflowDto) error {
	if err := validateProjectManagementWorkflowDefinition(input); err != nil {
		return err
	}
	for _, stage := range input.Stages {
		if stage.Id == 0 {
			return errorx.New("project workflow stage id is required")
		}
	}
	return nil
}
