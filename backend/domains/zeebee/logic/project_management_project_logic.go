package logic

import (
	"context"
	"strings"

	"gorm.io/gorm"

	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/gormx"
	"zeebee-backend/common/utils"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/repo"
)

type ProjectManagementProjectLogic struct{ *BaseLogic }

type projectManagementProjectOperation struct {
	claim     utils.JwtUser
	projectId int64
	action    constants.OperationAction
	before    map[string]any
	after     map[string]any
}

func NewProjectManagementProjectLogic(baseLogic *BaseLogic) *ProjectManagementProjectLogic {
	return &ProjectManagementProjectLogic{BaseLogic: baseLogic}
}

func (l *ProjectManagementProjectLogic) Page(ctx context.Context, claim utils.JwtUser, input dto.ProjectManagementProjectPageRequestDto) (gormx.PageResponse[*model.ProjectManagementProject], error) {
	return l.repos.ProjectManagementProjectRepo.Page(ctx, claim.OrganizationId, input)
}

func (l *ProjectManagementMemberLogic) Page(ctx context.Context, claim utils.JwtUser, input dto.ProjectManagementMemberPageRequestDto) (gormx.PageResponse[*dto.ProjectManagementMemberDetailDto], error) {
	if _, err := l.projectManagementLogic().RequireProject(ctx, input.Filter.ProjectId, claim.OrganizationId); err != nil {
		return gormx.PageResponse[*dto.ProjectManagementMemberDetailDto]{}, err
	}
	if _, err := l.projectManagementLogic().RequireMember(ctx, input.Filter.ProjectId, claim.OrganizationUserId); err != nil {
		return gormx.PageResponse[*dto.ProjectManagementMemberDetailDto]{}, err
	}
	page, err := l.repos.ProjectManagementProjectRepo.PageMembers(ctx, input.Filter.ProjectId, input)
	if err != nil {
		return gormx.PageResponse[*dto.ProjectManagementMemberDetailDto]{}, err
	}
	list, err := projectManagementMemberDtos(page.List)
	if err != nil {
		return gormx.PageResponse[*dto.ProjectManagementMemberDetailDto]{}, err
	}
	return gormx.PageResponse[*dto.ProjectManagementMemberDetailDto]{Page: page.Page, Limit: page.Limit, Total: page.Total, TotalPage: page.TotalPage, List: list}, nil
}

func (l *ProjectManagementMemberLogic) Candidate(ctx context.Context, claim utils.JwtUser, input dto.ProjectManagementMemberCandidateDto) ([]*dto.ProjectManagementMemberCandidateEntityDto, error) {
	if input.Purpose != constants.ProjectManagementMemberCandidatePurposeAddMember && input.Purpose != constants.ProjectManagementMemberCandidatePurposeTransferOwner && input.Purpose != constants.ProjectManagementMemberCandidatePurposeProjectMember {
		return nil, errorx.New("invalid project member candidate purpose")
	}
	if _, err := l.projectManagementLogic().RequireProject(ctx, input.ProjectId, claim.OrganizationId); err != nil {
		return nil, err
	}
	if _, err := l.projectManagementLogic().RequireMemberRole(ctx, input.ProjectId, claim.OrganizationUserId, constants.ProjectManagementMemberRoleProjectOwner, constants.ProjectManagementMemberRoleProjectAdmin); err != nil {
		return nil, err
	}
	candidates, err := l.repos.ProjectManagementProjectRepo.CandidateMembers(ctx, claim.OrganizationId, input.ProjectId, input.Keyword, input.Purpose)
	if err != nil {
		return nil, err
	}
	return candidates, nil
}

func (l *ProjectManagementProjectLogic) RequirePropertyDeletable(ctx context.Context, propertyId int64) error {
	blocked, err := l.repos.ProjectManagementProjectRepo.HasUnfinishedTasksForEntity(ctx, propertyId, constants.ProjectManagementEntityTypeBorrowerProperty)
	if err != nil {
		return err
	}
	if blocked {
		return errorx.New("cannot delete a property with unfinished project tasks")
	}
	return nil
}

func projectManagementMemberDtos(models []*repo.ProjectManagementProjectMemberEntity) ([]*dto.ProjectManagementMemberDetailDto, error) {
	result := make([]*dto.ProjectManagementMemberDetailDto, 0, len(models))
	for _, item := range models {
		member, err := copierx.Trans[repo.ProjectManagementProjectMemberEntity, dto.ProjectManagementMemberDetailDto](item)
		if err != nil {
			return nil, err
		}
		result = append(result, member)
	}
	return result, nil
}

func (l *ProjectManagementProjectLogic) Create(ctx context.Context, claim utils.JwtUser, input dto.ProjectManagementProjectCreateDto) (project *model.ProjectManagementProject, err error) {
	if err = validateProjectManagementProjectCreate(input); err != nil {
		return nil, err
	}
	// Validate entity
	if input.EntityType == constants.ProjectManagementEntityTypeBorrowerProperty && input.EntityId > 0 {
		exists, err := l.repos.BorrowerPropertyRepo.IsExistedByColumns(ctx, map[string]any{"id": input.EntityId, "organization_id": claim.OrganizationId})
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, errorx.New("property not found")
		}
	}
	project, err = copierx.Trans[dto.ProjectManagementProjectDto, model.ProjectManagementProject](&input.ProjectManagementProjectDto)
	if err != nil {
		return nil, err
	}
	project.OrganizationId, project.OwnerId, project.OwnerName = claim.OrganizationId, claim.OrganizationUserId, claim.Username
	err = l.repos.ProjectManagementProjectRepo.Transaction(ctx, func(tx *gorm.DB) error {
		var innerErr error
		project, innerErr = l.repos.ProjectManagementProjectRepo.CreateWithTx(ctx, tx, project, false)
		if innerErr != nil {
			return innerErr
		}
		if innerErr = l.repos.ProjectManagementProjectRepo.CreateMemberWithTx(ctx, tx, &model.ProjectManagementProjectMember{
			ProjectId:          project.Id,
			OrganizationUserId: claim.OrganizationUserId,
			Role:               constants.ProjectManagementMemberRoleProjectOwner,
		}); innerErr != nil {
			return innerErr
		}
		if innerErr = l.repos.ProjectManagementWorkflowRepo.CloneTemplateToProjectWithTx(ctx, tx, input.WorkflowTemplateId, claim.OrganizationId, project.Id); innerErr != nil {
			return projectManagementNotFound(innerErr, "workflow template not found")
		}
		return l.recordOperation(ctx, tx, projectManagementProjectOperation{claim: claim, projectId: project.Id, action: constants.OperationActionProjectManagementProjectCreate, after: map[string]any{"id": project.Id, "name": project.Name}})
	})
	return project, err
}

func (l *ProjectManagementProjectLogic) Update(ctx context.Context, claim utils.JwtUser, projectId int64, input dto.ProjectManagementProjectUpdateDto) (*model.ProjectManagementProject, error) {
	if err := validateProjectManagementProjectUpdate(input); err != nil {
		return nil, err
	}
	project, err := l.projectManagementLogic().RequireProject(ctx, projectId, claim.OrganizationId)
	if err != nil {
		return nil, err
	}
	if _, err = l.projectManagementLogic().RequireMemberRole(ctx, projectId, claim.OrganizationUserId, constants.ProjectManagementMemberRoleProjectOwner, constants.ProjectManagementMemberRoleProjectAdmin); err != nil {
		return nil, err
	}
	before := projectOperationValues(project)
	if err = copierx.Assign(&input, project); err != nil {
		return nil, err
	}
	project.RemoveRelationships()

	err = l.repos.ProjectManagementProjectRepo.Transaction(ctx, func(tx *gorm.DB) error {
		if _, innerErr := l.repos.ProjectManagementProjectRepo.UpdateOneByIdWithTx(ctx, tx, project.Id, project); innerErr != nil {
			return innerErr
		}
		after := projectOperationValues(project)
		before, after = operationValuesDiff(before, after)
		if len(before) == 0 {
			return nil
		}
		return l.recordOperation(ctx, tx, projectManagementProjectOperation{
			claim:     claim,
			projectId: projectId,
			action:    constants.OperationActionProjectManagementProjectUpdate,
			before:    before,
			after:     after,
		})
	})
	return project, err
}

func (l *ProjectManagementProjectLogic) Delete(ctx context.Context, claim utils.JwtUser, projectId int64) error {
	return l.repos.ProjectManagementProjectRepo.Transaction(ctx, func(tx *gorm.DB) error {
		project, err := l.projectManagementLogic().requireProjectManagerWithTx(ctx, tx, claim, projectId)
		if err != nil {
			return err
		}
		taskIds, err := l.repos.ProjectManagementTaskRepo.IdsByProjectWithTx(ctx, tx, projectId)
		if err != nil {
			return err
		}
		if err = l.repos.FileRepo.DeleteByEntityIdWithTx(ctx, tx, taskIds, constants.FileEntityTypeProjectManagementTask); err != nil {
			return err
		}
		if err = l.repos.ProjectManagementTaskRepo.DeleteByProjectWithTx(ctx, tx, projectId, taskIds); err != nil {
			return err
		}
		if err = l.repos.ProjectManagementWorkflowRepo.DeleteProjectWorkflowWithTx(ctx, tx, projectId); err != nil {
			return err
		}
		if err = l.repos.ProjectManagementProjectRepo.DeleteProjectWithTx(ctx, tx, projectId); err != nil {
			return err
		}
		return l.recordOperation(ctx, tx, projectManagementProjectOperation{claim: claim, projectId: projectId, action: constants.OperationActionProjectManagementProjectDelete, before: projectOperationValues(project)})
	})
}

func validateProjectManagementProjectCreate(input dto.ProjectManagementProjectCreateDto) error {
	if strings.TrimSpace(input.Name) == "" || input.WorkflowTemplateId == 0 {
		return errorx.New("project name and workflow template are required")
	}
	return input.Priority.Validate()
}

func validateProjectManagementProjectUpdate(input dto.ProjectManagementProjectUpdateDto) error {
	if strings.TrimSpace(input.Name) == "" {
		return errorx.New("project name is required")
	}
	return input.Priority.Validate()
}

func (l *ProjectManagementProjectLogic) recordOperation(ctx context.Context, tx *gorm.DB, input projectManagementProjectOperation) error {
	metadata, err := projectManagementOperationMetadata(map[string]any{"project_id": input.projectId}, input.before, input.after)
	if err != nil {
		return errorx.WithStack(err)
	}
	return NewOperationLogic(l.BaseLogic).CreateOperationLog(ctx, tx, dto.OperationCreateDto{EntityId: input.projectId, EntityType: constants.OperationEntityTypeProjectManagementProject, OperatorId: input.claim.OrganizationUserId, OperatorName: input.claim.Username, OperatorType: constants.UserTypeBorrower, Operation: input.action, Metadata: metadata, OrganizationId: input.claim.OrganizationId})
}

func projectOperationValues(project *model.ProjectManagementProject) map[string]any {
	return map[string]any{
		"name":        project.Name,
		"description": project.Description,
		"priority":    project.Priority,
		"cover_url":   project.CoverUrl,
		"entity_id":   project.EntityId,
		"entity_type": project.EntityType,
		"owner_id":    project.OwnerId,
		"owner_name":  project.OwnerName,
	}
}
