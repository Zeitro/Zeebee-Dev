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

func (l *ProjectManagementTemplateLogic) CreateWorkflow(ctx context.Context, claim utils.JwtUser, templateID int64, input dto.ProjectManagementWorkflowDto) (*model.ProjectManagementWorkflow, error) {
	stageIDs := templateWorkflowStageIDs(input.Transitions)
	if err := l.validateTemplateWorkflow(ctx, templateID, stageIDs, input.Transitions); err != nil {
		return nil, err
	}
	var workflow *model.ProjectManagementWorkflow
	err := l.repos.ProjectManagementTemplateRepo.Transaction(ctx, func(tx *gorm.DB) error {
		if _, err := l.repos.ProjectManagementTemplateRepo.FindPublicByID(ctx, templateID); err != nil {
			return projectManagementNotFound(err, "project template not found")
		}
		workflow = &model.ProjectManagementWorkflow{EntityType: constants.ProjectManagementOwnerTypeTemplate, EntityId: templateID, Name: input.Name, Description: input.Description}
		if err := tx.WithContext(ctx).Create(workflow).Error; err != nil {
			return errorx.WithStack(err)
		}
		return l.replaceTemplateWorkflowDefinition(ctx, tx, templateID, workflow.Id, stageIDs, input.Transitions)
	})
	if err != nil {
		return nil, err
	}
	return l.repos.ProjectManagementTemplateRepo.FindWorkflow(ctx, templateID, workflow.Id)
}

func (l *ProjectManagementTemplateLogic) UpdateWorkflow(ctx context.Context, claim utils.JwtUser, templateID, workflowID int64, input dto.ProjectManagementWorkflowDto) (*model.ProjectManagementWorkflow, error) {
	stageIDs := templateWorkflowStageIDs(input.Transitions)
	if err := l.validateTemplateWorkflow(ctx, templateID, stageIDs, input.Transitions); err != nil {
		return nil, err
	}
	err := l.repos.ProjectManagementTemplateRepo.Transaction(ctx, func(tx *gorm.DB) error {
		workflow, err := l.repos.ProjectManagementTemplateRepo.FindWorkflow(ctx, templateID, workflowID)
		if err != nil {
			return projectManagementNotFound(err, "template workflow not found")
		}
		if err = tx.WithContext(ctx).Model(&model.ProjectManagementWorkflow{}).Where("id = ?", workflow.Id).Updates(map[string]any{"name": input.Name, "description": input.Description}).Error; err != nil {
			return errorx.WithStack(err)
		}
		return l.replaceTemplateWorkflowDefinition(ctx, tx, templateID, workflowID, stageIDs, input.Transitions)
	})
	if err != nil {
		return nil, err
	}
	return l.repos.ProjectManagementTemplateRepo.FindWorkflow(ctx, templateID, workflowID)
}

func (l *ProjectManagementTemplateLogic) DeleteWorkflow(ctx context.Context, claim utils.JwtUser, templateID, workflowID int64) error {
	return l.repos.ProjectManagementTemplateRepo.Transaction(ctx, func(tx *gorm.DB) error {
		workflow, err := l.repos.ProjectManagementTemplateRepo.FindWorkflow(ctx, templateID, workflowID)
		if err != nil {
			return projectManagementNotFound(err, "template workflow not found")
		}
		if workflow.IsDefault {
			return errorx.New("default workflow cannot be deleted")
		}
		if err = tx.WithContext(ctx).Where("workflow_id = ?", workflowID).Delete(&model.ProjectManagementWorkflowStageTransition{}).Error; err != nil {
			return errorx.WithStack(err)
		}
		if err = tx.WithContext(ctx).Where("workflow_id = ?", workflowID).Delete(&model.ProjectManagementWorkflowStage{}).Error; err != nil {
			return errorx.WithStack(err)
		}
		return errorx.WithStack(tx.WithContext(ctx).Delete(&model.ProjectManagementWorkflow{}, workflowID).Error)
	})
}

func (l *ProjectManagementTemplateLogic) WorkflowDetail(ctx context.Context, claim utils.JwtUser, templateID, workflowID int64) (*model.ProjectManagementWorkflow, error) {
	if _, err := l.Detail(ctx, claim, templateID); err != nil {
		return nil, err
	}
	workflow, err := l.repos.ProjectManagementTemplateRepo.FindWorkflow(ctx, templateID, workflowID)
	if err != nil {
		return nil, projectManagementNotFound(err, "template workflow not found")
	}
	return workflow, nil
}

func (l *ProjectManagementTemplateLogic) SetWorkflowDefault(ctx context.Context, claim utils.JwtUser, templateID, workflowID int64) error {
	return l.repos.ProjectManagementTemplateRepo.Transaction(ctx, func(tx *gorm.DB) error {
		if _, err := l.repos.ProjectManagementTemplateRepo.FindWorkflow(ctx, templateID, workflowID); err != nil {
			return projectManagementNotFound(err, "template workflow not found")
		}
		if err := tx.WithContext(ctx).Model(&model.ProjectManagementWorkflow{}).Where("entity_type = ? AND entity_id = ?", constants.ProjectManagementOwnerTypeTemplate, templateID).Update("is_default", false).Error; err != nil {
			return errorx.WithStack(err)
		}
		return errorx.WithStack(tx.WithContext(ctx).Model(&model.ProjectManagementWorkflow{}).Where("id = ?", workflowID).Update("is_default", true).Error)
	})
}

func (l *ProjectManagementTemplateLogic) validateTemplateWorkflow(ctx context.Context, templateID int64, stageIDs []int64, transitions []dto.ProjectManagementWorkflowTransitionDto) error {
	stages, err := l.repos.ProjectManagementTemplateRepo.ListStages(ctx, templateID)
	if err != nil {
		return err
	}
	byID := make(map[int64]*model.ProjectManagementStage, len(stages))
	for _, stage := range stages {
		byID[stage.Id] = stage
	}
	definition := dto.ProjectManagementWorkflowDto{Name: "template workflow", Stages: make([]dto.ProjectManagementWorkflowStageDto, 0, len(stageIDs)), Transitions: transitions}
	seen := make(map[int64]struct{}, len(stageIDs))
	for _, id := range stageIDs {
		stage := byID[id]
		if _, exists := seen[id]; stage == nil || exists {
			return errorx.New("template workflow stage is invalid")
		}
		seen[id] = struct{}{}
		definition.Stages = append(definition.Stages, dto.ProjectManagementWorkflowStageDto{Id: id, Name: stage.Name, Description: stage.Description, Phase: stage.Phase, Color: stage.Color, SortOrder: stage.SortOrder})
	}
	return ValidateProjectManagementWorkflow(definition)
}

func templateWorkflowStageIDs(transitions []dto.ProjectManagementWorkflowTransitionDto) []int64 {
	seen := make(map[int64]struct{}, len(transitions)*2)
	stageIDs := make([]int64, 0, len(transitions)*2)
	for _, transition := range transitions {
		for _, stageID := range []int64{transition.FromStageId, transition.ToStageId} {
			if _, exists := seen[stageID]; exists {
				continue
			}
			seen[stageID] = struct{}{}
			stageIDs = append(stageIDs, stageID)
		}
	}
	return stageIDs
}

func (l *ProjectManagementTemplateLogic) replaceTemplateWorkflowDefinition(ctx context.Context, tx *gorm.DB, templateID, workflowID int64, stageIDs []int64, transitions []dto.ProjectManagementWorkflowTransitionDto) error {
	if err := tx.WithContext(ctx).Where("workflow_id = ?", workflowID).Delete(&model.ProjectManagementWorkflowStageTransition{}).Error; err != nil {
		return errorx.WithStack(err)
	}
	if err := tx.WithContext(ctx).Where("workflow_id = ?", workflowID).Delete(&model.ProjectManagementWorkflowStage{}).Error; err != nil {
		return errorx.WithStack(err)
	}
	for _, stageID := range stageIDs {
		link := &model.ProjectManagementWorkflowStage{WorkflowId: workflowID, StageId: stageID}
		if err := tx.WithContext(ctx).Create(link).Error; err != nil {
			return errorx.WithStack(err)
		}
	}
	for _, transition := range transitions {
		item := &model.ProjectManagementWorkflowStageTransition{WorkflowId: workflowID, Name: transition.Name, FromStageId: transition.FromStageId, ToStageId: transition.ToStageId, AllowedActorTypes: transition.AllowedActorTypes}
		if err := tx.WithContext(ctx).Create(item).Error; err != nil {
			return errorx.WithStack(err)
		}
	}
	return nil
}
