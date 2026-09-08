package logic

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"zeebee-backend/common/errorx"
	"zeebee-backend/common/utils"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
)

type projectManagementTransitionContext struct {
	ctx          context.Context
	tx           *gorm.DB
	claim        utils.JwtUser
	task         *model.ProjectManagementTask
	project      *model.ProjectManagementProject
	member       *model.ProjectManagementProjectMember
	currentStage *model.ProjectManagementWorkflowStage
	targetStage  *model.ProjectManagementWorkflowStage
	transition   *model.ProjectManagementWorkflowTransition
}

func (l *ProjectManagementTaskLogic) Transition(ctx context.Context, claim utils.JwtUser, taskId, targetStageId int64) (result *dto.ProjectManagementTaskDto, err error) {
	if taskId == 0 || targetStageId == 0 {
		return nil, errorx.New("task and target stage are required")
	}
	var task *model.ProjectManagementTask
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
			return projectManagementNotFound(err, fmt.Sprintf("workflow has no transition configured from %q to %q", currentStage.Name, targetStage.Name))
		}
		transitionCtx := &projectManagementTransitionContext{
			ctx:          ctx,
			tx:           tx,
			claim:        claim,
			task:         task,
			project:      project,
			member:       member,
			currentStage: currentStage,
			targetStage:  targetStage,
			transition:   transition,
		}
		if err = l.requireTransitionActor(transitionCtx); err != nil {
			return err
		}
		if err = l.validateTransitionDependencies(transitionCtx); err != nil {
			return err
		}
		if targetStage.Phase == constants.ProjectManagementStagePhaseFinished {
			if err = l.validateFinishedChildren(transitionCtx); err != nil {
				return err
			}
		}
		if err = l.repos.ProjectManagementTaskRepo.CopyCurrentStageAssigneesIfTargetEmptyWithTx(ctx, tx, task, targetStage.Id); err != nil {
			return err
		}
		task.StageId, task.Phase = targetStage.Id, targetStage.Phase
		task.RemoveRelationships()
		if _, err = l.repos.ProjectManagementTaskRepo.UpdateOneWithTx(ctx, tx, task); err != nil {
			return err
		}
		return l.recordProjectTaskTransition(transitionCtx)
	})
	if err != nil {
		return nil, err
	}
	return l.Detail(ctx, claim, task.Id)
}

func (l *ProjectManagementTaskLogic) requireTransitionActor(transitionCtx *projectManagementTransitionContext) error {
	if len(transitionCtx.transition.AllowedActorTypes) == 0 {
		return nil
	}
	for _, actorType := range transitionCtx.transition.AllowedActorTypes {
		switch constants.ProjectManagementMemberRole(actorType) {
		case constants.ProjectManagementMemberRoleProjectOwner:
			if transitionCtx.project.OwnerId == transitionCtx.claim.OrganizationUserId || transitionCtx.member.Role == constants.ProjectManagementMemberRoleProjectOwner {
				return nil
			}
		case constants.ProjectManagementMemberRoleProjectAdmin:
			if transitionCtx.member.Role == constants.ProjectManagementMemberRoleProjectAdmin {
				return nil
			}
		case constants.ProjectManagementMemberRoleTaskCreator:
			if transitionCtx.task.CreatorId == transitionCtx.claim.OrganizationUserId {
				return nil
			}
		case constants.ProjectManagementMemberRoleTaskExecutor:
			if transitionCtx.member.Role == constants.ProjectManagementMemberRoleTaskExecutor {
				return nil
			}
		}
	}
	return errorx.New("project member is not allowed to perform this transition")
}

func (l *ProjectManagementTaskLogic) validateTransitionDependencies(transitionCtx *projectManagementTransitionContext) error {
	dependencies, err := l.repos.ProjectManagementDependencyRepo.ListBySuccessorWithTx(transitionCtx.ctx, transitionCtx.tx, transitionCtx.task.Id)
	if err != nil {
		return err
	}
	for _, dependency := range dependencies {
		_, predecessorTaskId := dependencyEndpoints(dependency.TaskId, dependency.DependencyTaskId, dependency.DependencyType)
		predecessor, err := l.repos.ProjectManagementTaskRepo.LockWithTx(transitionCtx.ctx, transitionCtx.tx, predecessorTaskId)
		if err != nil {
			return projectManagementNotFound(err, "dependency predecessor task not found")
		}
		stage, err := l.repos.ProjectManagementWorkflowRepo.FindStageWithTx(transitionCtx.ctx, transitionCtx.tx, predecessor.StageId)
		if err != nil {
			return projectManagementNotFound(err, "dependency predecessor stage not found")
		}
		finished, started := stage.Phase == constants.ProjectManagementStagePhaseFinished, stage.Phase != constants.ProjectManagementStagePhasePending
		if transitionCtx.targetStage.Phase == constants.ProjectManagementStagePhaseInProgress && ((dependency.DependencyRule == constants.ProjectManagementDependencyRuleFS && !finished) || (dependency.DependencyRule == constants.ProjectManagementDependencyRuleSS && !started)) {
			state := "started"
			if dependency.DependencyRule == constants.ProjectManagementDependencyRuleFS {
				state = "finished"
			}
			return errorx.New(fmt.Sprintf("cannot move task to InProgress: predecessor task %q must be %s (%s dependency)", predecessor.Title, state, dependency.DependencyRule))
		}
		if transitionCtx.targetStage.Phase == constants.ProjectManagementStagePhaseFinished && ((dependency.DependencyRule == constants.ProjectManagementDependencyRuleFF && !finished) || (dependency.DependencyRule == constants.ProjectManagementDependencyRuleSF && !started)) {
			state := "started"
			if dependency.DependencyRule == constants.ProjectManagementDependencyRuleFF {
				state = "finished"
			}
			return errorx.New(fmt.Sprintf("cannot move task to Finished: predecessor task %q must be %s (%s dependency)", predecessor.Title, state, dependency.DependencyRule))
		}
	}
	return nil
}

func (l *ProjectManagementTaskLogic) validateFinishedChildren(transitionCtx *projectManagementTransitionContext) error {
	children, err := l.repos.ProjectManagementTaskRepo.DirectChildrenWithTx(transitionCtx.ctx, transitionCtx.tx, transitionCtx.task.Id)
	if err != nil {
		return err
	}
	for _, child := range children {
		stage, err := l.repos.ProjectManagementWorkflowRepo.FindStageWithTx(transitionCtx.ctx, transitionCtx.tx, child.StageId)
		if err != nil {
			return projectManagementNotFound(err, "child task stage not found")
		}
		if stage.Phase != constants.ProjectManagementStagePhaseFinished {
			return errorx.New("unfinished child task blocks transition to finished")
		}
	}
	return nil
}

func (l *ProjectManagementTaskLogic) recordProjectTaskTransition(transitionCtx *projectManagementTransitionContext) error {
	metadata, err := marshalProjectManagementOperationMetadata(
		projectManagementOperationMetadata{ProjectID: transitionCtx.task.ProjectId, TaskID: transitionCtx.task.Id, WorkflowID: transitionCtx.currentStage.WorkflowId, StageID: transitionCtx.currentStage.Id, TransitionID: transitionCtx.transition.Id},
		stageOperationValues(transitionCtx.currentStage),
		stageOperationValues(transitionCtx.targetStage),
	)
	if err != nil {
		return errorx.WithStack(err)
	}
	input := dto.OperationCreateDto{
		EntityId:       transitionCtx.task.ProjectId,
		EntityType:     constants.OperationEntityTypeProjectManagementProject,
		OperatorId:     transitionCtx.claim.OrganizationUserId,
		OperatorName:   transitionCtx.claim.Username,
		OperatorType:   constants.UserTypeBorrower,
		Operation:      constants.OperationActionProjectManagementTaskTransition,
		Metadata:       metadata,
		OrganizationId: transitionCtx.claim.OrganizationId,
	}
	return NewOperationLogic(l.BaseLogic).CreateOperationLog(transitionCtx.ctx, transitionCtx.tx, input)
}

func stageOperationValues(stage *model.ProjectManagementWorkflowStage) map[string]any {
	return map[string]any{"stage_id": stage.Id, "stage_name": stage.Name}
}
