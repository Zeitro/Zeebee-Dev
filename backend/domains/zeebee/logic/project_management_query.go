package logic

import (
	"context"
	"time"

	"zeebee-backend/common/gormx"
	"zeebee-backend/common/utils"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
)

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

func (l *ProjectManagementTaskLogic) Overdue(ctx context.Context, task *model.ProjectManagementTask) (*dto.ProjectManagementTaskOverdueDto, error) {
	stage, err := l.repos.ProjectManagementWorkflowRepo.FindStage(ctx, task.StageId)
	if err != nil {
		return nil, projectManagementNotFound(err, "task stage not found")
	}
	overdue, overdueSeconds := projectManagementTaskOverdue(task, stage.Phase, time.Now())
	return &dto.ProjectManagementTaskOverdueDto{Overdue: overdue, OverdueSeconds: overdueSeconds}, nil
}

func (l *ProjectManagementTaskLogic) MyPage(ctx context.Context, claim utils.JwtUser, input dto.ProjectManagementMyTaskPageRequestDto) (gormx.PageResponse[*model.ProjectManagementTask], error) {
	page, err := l.repos.ProjectManagementTaskRepo.PageMine(ctx, claim.OrganizationId, claim.OrganizationUserId, input)
	if err != nil {
		return page, err
	}
	return page, nil
}
