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

func validateProjectWorkflowInvariant(projectId int64, workflows []*model.ProjectManagementWorkflow, stages []*model.ProjectManagementStage, links []*model.ProjectManagementWorkflowStage) error {
	if len(workflows) == 0 {
		return errorx.New("project requires a workflow")
	}
	stageByID := make(map[int64]*model.ProjectManagementStage, len(stages))
	pendingCount, finishedCount := 0, 0
	for _, stage := range stages {
		if stage.EntityType != constants.ProjectManagementOwnerTypeProject || stage.EntityId != projectId {
			return errorx.New("project stage owner mismatch")
		}
		stageByID[stage.Id] = stage
		if stage.Phase == constants.ProjectManagementStagePhasePending {
			pendingCount++
		}
		if stage.Phase == constants.ProjectManagementStagePhaseFinished {
			finishedCount++
		}
	}
	if pendingCount != 1 || finishedCount == 0 {
		return errorx.New("project requires exactly one pending stage and at least one finished stage")
	}
	workflowByID := make(map[int64]struct{}, len(workflows))
	defaultCount := 0
	for _, workflow := range workflows {
		if workflow.EntityType != constants.ProjectManagementOwnerTypeProject || workflow.EntityId != projectId {
			return errorx.New("project workflow owner mismatch")
		}
		workflowByID[workflow.Id] = struct{}{}
		if workflow.IsDefault {
			defaultCount++
		}
	}
	if defaultCount != 1 {
		return errorx.New("project requires exactly one default workflow")
	}
	seenLinks := make(map[[2]int64]struct{}, len(links))
	for _, link := range links {
		if _, ok := workflowByID[link.WorkflowId]; !ok {
			return errorx.New("workflow stage link workflow mismatch")
		}
		if _, ok := stageByID[link.StageId]; !ok {
			return errorx.New("workflow stage link stage mismatch")
		}
		key := [2]int64{link.WorkflowId, link.StageId}
		if _, ok := seenLinks[key]; ok {
			return errorx.New("duplicate workflow stage link")
		}
		seenLinks[key] = struct{}{}
	}
	return nil
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
