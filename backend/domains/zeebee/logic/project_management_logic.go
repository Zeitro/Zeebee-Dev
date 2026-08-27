package logic

import (
	"context"
	"errors"
	"slices"
	"time"

	"gorm.io/gorm"

	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/model"
)

type ProjectManagementLogic struct {
	*BaseLogic
}

func NewProjectManagementLogic(baseLogic *BaseLogic) *ProjectManagementLogic {
	return &ProjectManagementLogic{BaseLogic: baseLogic}
}

func (l *BaseLogic) projectManagementLogic() *ProjectManagementLogic {
	return NewProjectManagementLogic(l)
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

func projectManagementTaskOverdue(task *model.ProjectManagementTask, phase constants.ProjectManagementStagePhase, now time.Time) (bool, int64) {
	if task.EndAt == nil || !task.EndAt.Valid || phase == constants.ProjectManagementStagePhaseFinished {
		return false, 0
	}
	// end_at is date-only (midnight), so +1 day = start of the day after the deadline.
	deadline := task.EndAt.Time.AddDate(0, 0, 1)
	if now.Before(deadline) {
		return false, 0
	}
	return true, int64(now.Sub(deadline).Seconds())
}
