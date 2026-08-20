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

func (l *ProjectManagementWorkflowLogic) ProjectDetail(ctx context.Context, claim *utils.JwtUser, projectId int64) (*model.ProjectManagementWorkflow, error) {
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

func (l *ProjectManagementWorkflowLogic) ProjectUpdate(ctx context.Context, claim *utils.JwtUser, projectId int64, input dto.ProjectManagementWorkflowDto) (workflow *model.ProjectManagementWorkflow, err error) {
	if err = validateProjectManagementProjectWorkflow(input); err != nil {
		return nil, err
	}
	definition, stages, transitions, err := projectManagementTemplateModels(input)
	if err != nil {
		return nil, err
	}
	err = l.repos.ProjectManagementWorkflowRepo.Transaction(ctx, func(tx *gorm.DB) error {
		if _, err = l.projectManagementAccess().requireProjectManagerWithTx(ctx, tx, claim, projectId); err != nil {
			return err
		}
		current, err := l.repos.ProjectManagementWorkflowRepo.LockProjectWithTx(ctx, tx, projectId, claim.OrganizationId)
		if err != nil {
			return projectManagementNotFound(err, "project workflow not found")
		}
		definition.Id, definition.OrganizationId, definition.ProjectId = current.Id, claim.OrganizationId, projectId
		if err = l.repos.ProjectManagementWorkflowRepo.UpdateProjectWithTx(ctx, tx, definition); err != nil {
			return err
		}
		if err = l.repos.ProjectManagementWorkflowRepo.ReplaceProjectDefinitionWithTx(ctx, tx, current.Id, stages, transitions); err != nil {
			return err
		}
		return NewProjectManagementTaskLogic(l.BaseLogic).recordProjectTaskOperation(ctx, tx, claim, projectId, 0, constants.OperationActionProjectManagementWorkflowUpdate)
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
