package logic

import (
	"context"

	"gorm.io/gorm"

	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/model"
)

func (l *ProjectManagementWorkflowLogic) List(ctx context.Context, organizationId int64) ([]*model.ProjectManagementWorkflow, error) {
	workflows, err := l.repos.ProjectManagementWorkflowRepo.ListByOrganization(ctx, organizationId)
	if err != nil {
		return workflows, err
	}
	for _, workflow := range workflows {
		if workflow.IsDefault {
			return workflows, nil
		}
	}
	if err = l.ensureDefault(ctx, organizationId); err != nil {
		workflows, retryErr := l.repos.ProjectManagementWorkflowRepo.ListByOrganization(ctx, organizationId)
		if retryErr == nil {
			for _, workflow := range workflows {
				if workflow.IsDefault {
					return workflows, nil
				}
			}
		}
		return nil, err
	}
	return l.repos.ProjectManagementWorkflowRepo.ListByOrganization(ctx, organizationId)
}

func (l *ProjectManagementWorkflowLogic) ensureDefault(ctx context.Context, organizationId int64) error {
	return l.repos.ProjectManagementWorkflowRepo.Transaction(ctx, func(tx *gorm.DB) error {
		workflow, err := l.repos.ProjectManagementWorkflowRepo.CreateWithTx(ctx, tx, &model.ProjectManagementWorkflow{
			Kind:           constants.ProjectManagementWorkflowKindOrganizationTemplate,
			OrganizationId: organizationId,
			Name:           "Default workflow",
			IsDefault:      true,
		}, false)
		if err != nil {
			return err
		}
		stages := []*model.ProjectManagementWorkflowStage{
			{WorkflowId: workflow.Id, Name: "Pending", Phase: constants.ProjectManagementStagePhasePending, SortOrder: 10},
			{WorkflowId: workflow.Id, Name: "In progress", Phase: constants.ProjectManagementStagePhaseInProgress, SortOrder: 20},
			{WorkflowId: workflow.Id, Name: "Finished", Phase: constants.ProjectManagementStagePhaseFinished, SortOrder: 30},
		}
		if err = tx.WithContext(ctx).Create(&stages).Error; err != nil {
			return errorx.WithStack(err)
		}
		transitions := []*model.ProjectManagementWorkflowTransition{
			{WorkflowId: workflow.Id, FromStageId: stages[0].Id, ToStageId: stages[1].Id},
			{WorkflowId: workflow.Id, FromStageId: stages[1].Id, ToStageId: stages[2].Id},
		}
		return errorx.WithStack(tx.WithContext(ctx).Create(&transitions).Error)
	})
}
