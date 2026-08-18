package logic

import (
	"context"
	"errors"
	"slices"

	"gorm.io/gorm"

	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
)

func (l *ProjectManagementLogic) RequireActor(ctx context.Context, organizationId, userId int64) (*dto.ProjectManagementActor, error) {
	organizationUser, err := l.repos.OrganizationUserRepo.GetOrganizationUserByUserId(ctx, organizationId, userId)
	if err != nil {
		return nil, projectManagementNotFound(err, "organization member not found")
	}
	return &dto.ProjectManagementActor{
		OrganizationId:     organizationId,
		OrganizationUserId: organizationUser.Id,
		Name:               organizationUser.FullName,
	}, nil
}

func (l *ProjectManagementLogic) RequireProject(ctx context.Context, projectId, organizationId int64) (*model.ProjectManagementProject, error) {
	project, err := l.repos.ProjectManagementProjectRepo.FindByOrganization(ctx, projectId, organizationId)
	if err != nil {
		return nil, projectManagementNotFound(err, "project not found")
	}
	return project, nil
}

func (l *ProjectManagementLogic) RequireMember(ctx context.Context, projectId, organizationUserId int64) (*model.ProjectManagementProjectMember, error) {
	member, err := l.repos.ProjectManagementProjectRepo.FindMember(ctx, projectId, organizationUserId)
	if err != nil {
		return nil, projectManagementNotFound(err, "project member required")
	}
	return member, nil
}

func projectManagementNotFound(err error, message string) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errorx.New(message)
	}
	return err
}

func (l *ProjectManagementLogic) RequireMemberRole(ctx context.Context, projectId, organizationUserId int64, roles ...constants.ProjectManagementMemberRole) (*model.ProjectManagementProjectMember, error) {
	member, err := l.RequireMember(ctx, projectId, organizationUserId)
	if err != nil {
		return nil, err
	}
	if slices.Contains(roles, member.Role) {
		return member, nil
	}
	return nil, errorx.New("project member role is not allowed")
}
