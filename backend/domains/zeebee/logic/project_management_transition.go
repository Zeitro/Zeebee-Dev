package logic

import (
	"context"
	"encoding/json"

	"gorm.io/gorm"

	"zeebee-backend/common/errorx"
	"zeebee-backend/common/utils"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
)

func (l *ProjectManagementTaskLogic) Transition(ctx context.Context, claim *utils.JwtUser, taskId, targetStageId int64) (task *model.ProjectManagementTask, err error) {
	if taskId == 0 || targetStageId == 0 {
		return nil, errorx.New("task and target stage are required")
	}
	err = l.repos.ProjectManagementTaskRepo.Transaction(ctx, func(tx *gorm.DB) error {
		task, err = l.repos.ProjectManagementTaskRepo.LockWithTx(ctx, tx, taskId)
		if err != nil {
			return projectManagementNotFound(err, "task not found")
		}
		if task.OrganizationId != claim.OrganizationId {
			return errorx.New("task not found")
		}
		project, err := l.repos.ProjectManagementProjectRepo.LockByOrganizationWithTx(ctx, tx, task.ProjectId, claim.OrganizationId)
		if err != nil {
			return projectManagementNotFound(err, "project not found")
		}
		member, err := l.repos.ProjectManagementProjectRepo.LockMemberWithTx(ctx, tx, task.ProjectId, claim.OrganizationUserId)
		if err != nil {
			return projectManagementNotFound(err, "project member required")
		}
		workflow, err := l.repos.ProjectManagementWorkflowRepo.FindProjectDefinitionWithTx(ctx, tx, task.ProjectId, claim.OrganizationId)
		if err != nil {
			return projectManagementNotFound(err, "project workflow not found")
		}
		currentStage, err := l.repos.ProjectManagementWorkflowRepo.FindStageWithTx(ctx, tx, task.StageId)
		if err != nil {
			return projectManagementNotFound(err, "task stage not found")
		}
		targetStage, err := l.repos.ProjectManagementWorkflowRepo.FindStageWithTx(ctx, tx, targetStageId)
		if err != nil {
			return projectManagementNotFound(err, "target stage not found")
		}
		if currentStage.WorkflowId != workflow.Id || targetStage.WorkflowId != workflow.Id {
			return errorx.New("task transition stage must belong to the project workflow")
		}
		transition, err := l.repos.ProjectManagementWorkflowRepo.FindTransitionWithTx(ctx, tx, workflow.Id, currentStage.Id, targetStage.Id)
		if err != nil {
			return projectManagementNotFound(err, "task workflow transition not found")
		}
		if err = l.requireTransitionActor(ctx, tx, task, project, member, claim, transition); err != nil {
			return err
		}
		if err = l.validateTransitionDependencies(ctx, tx, task, currentStage, targetStage); err != nil {
			return err
		}
		if targetStage.Phase == constants.ProjectManagementStagePhaseFinished {
			if err = l.validateFinishedChildren(ctx, tx, task); err != nil {
				return err
			}
		}
		if err = l.repos.ProjectManagementTaskRepo.UpdateStageWithTx(ctx, tx, task.Id, targetStage.Id); err != nil {
			return err
		}
		task.StageId = targetStage.Id
		return l.recordProjectTaskTransition(ctx, tx, claim, task.ProjectId, task.Id, currentStage.Id, targetStage.Id)
	})
	if err != nil {
		return nil, err
	}
	return l.Detail(ctx, claim, task.Id)
}

func (l *ProjectManagementTaskLogic) requireTransitionActor(ctx context.Context, tx *gorm.DB, task *model.ProjectManagementTask, project *model.ProjectManagementProject, member *model.ProjectManagementProjectMember, claim *utils.JwtUser, transition *model.ProjectManagementWorkflowTransition) error {
	if len(transition.AllowedActorTypes) == 0 {
		return nil
	}
	for _, actorType := range transition.AllowedActorTypes {
		switch constants.ProjectManagementMemberRole(actorType) {
		case constants.ProjectManagementMemberRoleProjectOwner:
			if project.OwnerId == claim.OrganizationUserId || member.Role == constants.ProjectManagementMemberRoleProjectOwner {
				return nil
			}
		case constants.ProjectManagementMemberRoleProjectAdmin:
			if member.Role == constants.ProjectManagementMemberRoleProjectAdmin {
				return nil
			}
		case constants.ProjectManagementMemberRoleTaskCreator:
			if task.CreatorId == claim.OrganizationUserId {
				return nil
			}
		case constants.ProjectManagementMemberRoleTaskExecutor:
			if member.Role == constants.ProjectManagementMemberRoleTaskExecutor {
				return nil
			}
		}
	}
	return errorx.New("project member is not allowed to perform this transition")
}

func (l *ProjectManagementTaskLogic) validateTransitionDependencies(ctx context.Context, tx *gorm.DB, task *model.ProjectManagementTask, _ *model.ProjectManagementWorkflowStage, targetStage *model.ProjectManagementWorkflowStage) error {
	dependencies, err := l.repos.ProjectManagementDependencyRepo.ListBySuccessorWithTx(ctx, tx, task.Id)
	if err != nil {
		return err
	}
	for _, dependency := range dependencies {
		_, predecessorTaskId := dependencyEndpoints(dependency.TaskId, dependency.DependencyTaskId, dependency.DependencyType)
		predecessor, err := l.repos.ProjectManagementTaskRepo.LockWithTx(ctx, tx, predecessorTaskId)
		if err != nil {
			return projectManagementNotFound(err, "dependency predecessor task not found")
		}
		stage, err := l.repos.ProjectManagementWorkflowRepo.FindStageWithTx(ctx, tx, predecessor.StageId)
		if err != nil {
			return projectManagementNotFound(err, "dependency predecessor stage not found")
		}
		finished, started := stage.Phase == constants.ProjectManagementStagePhaseFinished, stage.Phase != constants.ProjectManagementStagePhasePending
		if targetStage.Phase == constants.ProjectManagementStagePhaseInProgress && ((dependency.DependencyRule == constants.ProjectManagementDependencyRuleFS && !finished) || (dependency.DependencyRule == constants.ProjectManagementDependencyRuleSS && !started)) {
			return errorx.New("task dependency blocks transition to in progress")
		}
		if targetStage.Phase == constants.ProjectManagementStagePhaseFinished && ((dependency.DependencyRule == constants.ProjectManagementDependencyRuleFF && !finished) || (dependency.DependencyRule == constants.ProjectManagementDependencyRuleSF && !started)) {
			return errorx.New("task dependency blocks transition to finished")
		}
	}
	return nil
}

func (l *ProjectManagementTaskLogic) validateFinishedChildren(ctx context.Context, tx *gorm.DB, task *model.ProjectManagementTask) error {
	children, err := l.repos.ProjectManagementTaskRepo.DirectChildrenWithTx(ctx, tx, task.Id)
	if err != nil {
		return err
	}
	for _, child := range children {
		stage, err := l.repos.ProjectManagementWorkflowRepo.FindStageWithTx(ctx, tx, child.StageId)
		if err != nil {
			return projectManagementNotFound(err, "child task stage not found")
		}
		if stage.Phase != constants.ProjectManagementStagePhaseFinished {
			return errorx.New("unfinished child task blocks transition to finished")
		}
	}
	return nil
}

func (l *ProjectManagementTaskLogic) recordProjectTaskTransition(ctx context.Context, tx *gorm.DB, claim *utils.JwtUser, projectId, taskId, fromStageId, toStageId int64) error {
	metadata, err := json.Marshal(map[string]int64{"project_id": projectId, "task_id": taskId, "from_stage_id": fromStageId, "to_stage_id": toStageId})
	if err != nil {
		return errorx.WithStack(err)
	}
	return NewOperationLogic(l.BaseLogic).CreateOperationLog(ctx, tx, dto.OperationCreateDto{EntityId: projectId, EntityType: constants.OperationEntityTypeProjectManagementProject, OperatorId: claim.OrganizationUserId, OperatorName: claim.Username, OperatorType: constants.UserTypeBorrower, Operation: constants.OperationActionProjectManagementTaskTransition, Metadata: metadata, OrganizationId: claim.OrganizationId})
}
