package logic

import (
	"context"
	"encoding/json"
	"strings"

	"gorm.io/gorm"

	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/gormx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/repo"
)

func (l *ProjectManagementLogic) PageProjects(ctx context.Context, actor *dto.ProjectManagementActor, input dto.ProjectManagementProjectPageRequestDto) (gormx.PageResponse[*model.ProjectManagementProject], error) {
	return l.repos.ProjectManagementProjectRepo.Page(ctx, actor.OrganizationId, input)
}

func (l *ProjectManagementLogic) PageProjectMembers(ctx context.Context, actor *dto.ProjectManagementActor, projectId int64, input dto.ProjectManagementMemberPageRequestDto) (gormx.PageResponse[*dto.ProjectManagementMemberDetailDto], error) {
	if _, err := l.RequireProject(ctx, projectId, actor.OrganizationId); err != nil {
		return gormx.PageResponse[*dto.ProjectManagementMemberDetailDto]{}, err
	}
	if _, err := l.RequireMember(ctx, projectId, actor.OrganizationUserId); err != nil {
		return gormx.PageResponse[*dto.ProjectManagementMemberDetailDto]{}, err
	}
	page, err := l.repos.ProjectManagementProjectRepo.PageMembers(ctx, projectId, input)
	if err != nil {
		return gormx.PageResponse[*dto.ProjectManagementMemberDetailDto]{}, err
	}
	list, err := projectManagementMemberDtos(page.List)
	if err != nil {
		return gormx.PageResponse[*dto.ProjectManagementMemberDetailDto]{}, err
	}
	return gormx.PageResponse[*dto.ProjectManagementMemberDetailDto]{Page: page.Page, Limit: page.Limit, Total: page.Total, TotalPage: page.TotalPage, List: list}, nil
}

func (l *ProjectManagementLogic) CandidateProjectMembers(ctx context.Context, actor *dto.ProjectManagementActor, projectId int64, keyword string) ([]*dto.ProjectManagementMemberDetailDto, error) {
	if _, err := l.RequireProject(ctx, projectId, actor.OrganizationId); err != nil {
		return nil, err
	}
	if _, err := l.RequireMemberRole(ctx, projectId, actor.OrganizationUserId, constants.ProjectManagementMemberRoleProjectOwner, constants.ProjectManagementMemberRoleProjectAdmin); err != nil {
		return nil, err
	}
	candidates, err := l.repos.ProjectManagementProjectRepo.CandidateMembers(ctx, actor.OrganizationId, projectId, keyword)
	if err != nil {
		return nil, err
	}
	return projectManagementMemberDtos(candidates)
}

func (l *ProjectManagementLogic) RequirePropertyDeletable(ctx context.Context, propertyId int64) error {
	blocked, err := l.repos.ProjectManagementProjectRepo.HasUnfinishedTasksForProperty(ctx, propertyId)
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

func (l *ProjectManagementLogic) CreateProject(ctx context.Context, actor *dto.ProjectManagementActor, input dto.ProjectManagementProjectCreateDto) (project *model.ProjectManagementProject, err error) {
	if err = validateProjectManagementProjectCreate(input); err != nil {
		return nil, err
	}
	if input.PropertyId != 0 {
		if _, err = l.repos.BorrowerPropertyRepo.GetOneByIdForOrganization(ctx, input.PropertyId, actor.OrganizationId); err != nil {
			return nil, projectManagementNotFound(err, "property not found")
		}
	}
	project, err = copierx.Trans[dto.ProjectManagementProjectDto, model.ProjectManagementProject](&input.ProjectManagementProjectDto)
	if err != nil {
		return nil, err
	}
	project.OrganizationId, project.OwnerId = actor.OrganizationId, actor.OrganizationUserId
	err = l.repos.ProjectManagementProjectRepo.Transaction(ctx, func(tx *gorm.DB) error {
		project, err = l.repos.ProjectManagementProjectRepo.CreateWithTx(ctx, tx, project, false)
		if err != nil {
			return err
		}
		if err = l.repos.ProjectManagementProjectRepo.CreateMemberWithTx(ctx, tx, &model.ProjectManagementProjectMember{ProjectId: project.Id, OrganizationUserId: actor.OrganizationUserId, Role: constants.ProjectManagementMemberRoleProjectOwner}); err != nil {
			return err
		}
		if err = l.repos.ProjectManagementWorkflowRepo.CloneTemplateToProjectWithTx(ctx, tx, input.WorkflowTemplateId, actor.OrganizationId, project.Id); err != nil {
			return projectManagementNotFound(err, "enabled workflow template not found")
		}
		return l.recordProjectCreateOperation(ctx, tx, actor, project.Id)
	})
	return project, err
}

func (l *ProjectManagementLogic) UpdateProject(ctx context.Context, actor *dto.ProjectManagementActor, projectId int64, input dto.ProjectManagementProjectUpdateDto) (*model.ProjectManagementProject, error) {
	if err := validateProjectManagementProjectUpdate(input); err != nil {
		return nil, err
	}
	project, err := l.RequireProject(ctx, projectId, actor.OrganizationId)
	if err != nil {
		return nil, err
	}
	if _, err = l.RequireMemberRole(ctx, projectId, actor.OrganizationUserId, constants.ProjectManagementMemberRoleProjectOwner, constants.ProjectManagementMemberRoleProjectAdmin); err != nil {
		return nil, err
	}
	if err = copierx.Assign(&input, project); err != nil {
		return nil, err
	}
	if err = l.repos.ProjectManagementProjectRepo.UpdateByOrganization(ctx, projectId, actor.OrganizationId, map[string]interface{}{"name": project.Name, "description": project.Description, "cover_url": project.CoverUrl}); err != nil {
		return nil, err
	}
	err = l.repos.ProjectManagementProjectRepo.Transaction(ctx, func(tx *gorm.DB) error {
		return l.recordProjectOperation(ctx, tx, actor, projectId, constants.OperationActionProjectManagementProjectUpdate)
	})
	return project, err
}

func (l *ProjectManagementLogic) DeleteProject(ctx context.Context, actor *dto.ProjectManagementActor, projectId int64, confirm bool) error {
	if err := validateProjectManagementProjectDelete(confirm); err != nil {
		return err
	}
	return l.repos.ProjectManagementProjectRepo.Transaction(ctx, func(tx *gorm.DB) error {
		if _, err := l.requireProjectManagerWithTx(ctx, tx, actor, projectId); err != nil {
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
		return l.recordProjectOperation(ctx, tx, actor, projectId, constants.OperationActionProjectManagementProjectDelete)
	})
}

func validateProjectManagementProjectCreate(input dto.ProjectManagementProjectCreateDto) error {
	if strings.TrimSpace(input.Name) == "" || input.WorkflowTemplateId == 0 {
		return errorx.New("project name and workflow template are required")
	}
	return nil
}

func validateProjectManagementProjectUpdate(input dto.ProjectManagementProjectUpdateDto) error {
	if strings.TrimSpace(input.Name) == "" {
		return errorx.New("project name is required")
	}
	return nil
}

func validateProjectManagementProjectDelete(confirm bool) error {
	if !confirm {
		return errorx.New("project deletion requires confirmation")
	}
	return nil
}

func (l *ProjectManagementLogic) recordProjectCreateOperation(ctx context.Context, tx *gorm.DB, actor *dto.ProjectManagementActor, projectId int64) error {
	return l.recordProjectOperation(ctx, tx, actor, projectId, constants.OperationActionProjectManagementProjectCreate)
}

func (l *ProjectManagementLogic) recordProjectOperation(ctx context.Context, tx *gorm.DB, actor *dto.ProjectManagementActor, projectId int64, operation constants.OperationAction) error {
	metadata, err := json.Marshal(map[string]int64{"project_id": projectId})
	if err != nil {
		return errorx.WithStack(err)
	}
	return NewOperationLogic(l.BaseLogic).CreateOperationLog(ctx, tx, dto.OperationCreateDto{EntityId: projectId, EntityType: constants.OperationEntityTypeProjectManagementProject, OperatorId: actor.OrganizationUserId, OperatorName: actor.Name, OperatorType: constants.UserTypeBorrower, Operation: operation, Metadata: metadata, OrganizationId: actor.OrganizationId})
}
