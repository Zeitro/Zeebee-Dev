package logic

import (
	"context"
	"time"

	"zeebee-backend/common/gormx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
)

func projectManagementTaskOverdue(task *model.ProjectManagementTask, phase constants.ProjectManagementStagePhase, now time.Time) (bool, int64) {
	if task.EndAt == nil || !task.EndAt.Valid || phase == constants.ProjectManagementStagePhaseFinished || !now.After(task.EndAt.Time) {
		return false, 0
	}
	return true, int64(now.Sub(task.EndAt.Time).Seconds())
}

func (l *ProjectManagementLogic) fillProjectManagementTaskOverdue(ctx context.Context, tasks []*model.ProjectManagementTask) error {
	for _, task := range tasks {
		stage, err := l.repos.ProjectManagementWorkflowRepo.FindStage(ctx, task.StageId)
		if err != nil {
			return projectManagementNotFound(err, "task stage not found")
		}
		task.Overdue, task.OverdueSeconds = projectManagementTaskOverdue(task, stage.Phase, time.Now())
	}
	return nil
}

func (l *ProjectManagementLogic) PageMyTasks(ctx context.Context, actor *dto.ProjectManagementActor, input dto.ProjectManagementMyTaskPageRequestDto) (gormx.PageResponse[*model.ProjectManagementTask], error) {
	page, err := l.repos.ProjectManagementTaskRepo.PageMine(ctx, actor.OrganizationId, actor.OrganizationUserId, input)
	if err != nil {
		return page, err
	}
	if err = l.fillProjectManagementTaskOverdue(ctx, page.List); err != nil {
		return gormx.PageResponse[*model.ProjectManagementTask]{}, err
	}
	return page, nil
}

func (l *ProjectManagementLogic) PageProjectOperations(ctx context.Context, actor *dto.ProjectManagementActor, input dto.ProjectManagementOperationPageRequestDto) (gormx.PageResponse[*model.Operation], error) {
	if _, err := l.RequireProject(ctx, input.Filter.ProjectId, actor.OrganizationId); err != nil {
		return gormx.PageResponse[*model.Operation]{}, err
	}
	if _, err := l.RequireMember(ctx, input.Filter.ProjectId, actor.OrganizationUserId); err != nil {
		return gormx.PageResponse[*model.Operation]{}, err
	}
	return l.repos.OperationRepo.PageProjectManagementOperations(ctx, actor.OrganizationId, input)
}
