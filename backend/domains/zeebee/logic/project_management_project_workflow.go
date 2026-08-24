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

func (l *ProjectManagementWorkflowLogic) ProjectUpdate(ctx context.Context, claim utils.JwtUser, projectId int64, input dto.ProjectManagementWorkflowDto) (workflow *model.ProjectManagementWorkflow, err error) {
	if err = validateProjectManagementProjectWorkflow(input); err != nil {
		return nil, err
	}
	definition, stages, transitions, err := projectManagementTemplateModels(input)
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
		current.Name, current.Description = definition.Name, definition.Description
		current.RemoveRelationships()
		if _, innerErr = l.repos.ProjectManagementWorkflowRepo.UpdateOneByIdWithTx(ctx, tx, current.Id, current); innerErr != nil {
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
	if err := ValidateProjectManagementWorkflow(input); err != nil {
		return err
	}
	for _, stage := range input.Stages {
		if stage.Id == 0 {
			return errorx.New("project workflow stage id is required")
		}
	}
	return nil
}
