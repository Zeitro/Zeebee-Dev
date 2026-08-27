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
	if _, err := l.projectManagementAccess().RequireProject(ctx, projectId, claim.OrganizationId); err != nil {
		return nil, err
	}
	if _, err := l.projectManagementAccess().RequireMember(ctx, projectId, claim.OrganizationUserId); err != nil {
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
		if _, innerErr = l.projectManagementAccess().requireProjectManagerWithTx(ctx, tx, claim, projectId); innerErr != nil {
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
		_, stages, transitions, innerErr := projectManagementTemplateModels(definition)
		if innerErr != nil {
			return innerErr
		}
		if innerErr = l.repos.ProjectManagementWorkflowRepo.ReplaceProjectDefinitionWithTx(ctx, tx, current.Id, stages, transitions); innerErr != nil {
			return innerErr
		}
		return NewProjectManagementTaskLogic(l.BaseLogic).recordProjectTaskOperation(ctx, tx, projectManagementTaskOperation{claim: claim, projectId: projectId, action: constants.OperationActionProjectManagementWorkflowUpdate})
	})
	if err != nil {
		return nil, err
	}
	return l.ProjectDetail(ctx, claim, projectId)
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
