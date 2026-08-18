package logic

import (
	"context"
	"encoding/json"

	"gorm.io/gorm"

	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/gormx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
)

func validateProjectManagementDependency(input dto.ProjectManagementTaskDependencyDto) error {
	if input.SuccessorTaskId == 0 || input.PredecessorTaskId == 0 {
		return errorx.New("successor and predecessor tasks are required")
	}
	if input.SuccessorTaskId == input.PredecessorTaskId {
		return errorx.New("task cannot depend on itself")
	}
	if input.DependencyType != constants.ProjectManagementDependencyTypeFS && input.DependencyType != constants.ProjectManagementDependencyTypeFF && input.DependencyType != constants.ProjectManagementDependencyTypeSS && input.DependencyType != constants.ProjectManagementDependencyTypeSF {
		return errorx.New("invalid task dependency type")
	}
	return nil
}

func (l *ProjectManagementLogic) CreateTaskDependency(ctx context.Context, actor *dto.ProjectManagementActor, input dto.ProjectManagementTaskDependencyDto) (*model.ProjectManagementTaskDependency, error) {
	if err := validateProjectManagementDependency(input); err != nil {
		return nil, err
	}
	var dependency *model.ProjectManagementTaskDependency
	err := l.repos.ProjectManagementDependencyRepo.Transaction(ctx, func(tx *gorm.DB) error {
		tasks, err := l.repos.ProjectManagementDependencyRepo.LockTasksWithTx(ctx, tx, input.SuccessorTaskId, input.PredecessorTaskId)
		if err != nil {
			return projectManagementNotFound(err, "task not found")
		}
		successor, predecessor := tasks[0], tasks[1]
		if successor.Id != input.SuccessorTaskId {
			successor, predecessor = predecessor, successor
		}
		if successor.OrganizationId != actor.OrganizationId || predecessor.OrganizationId != actor.OrganizationId || successor.ProjectId != predecessor.ProjectId {
			return errorx.New("dependency tasks must belong to the same project")
		}
		if _, err = l.repos.ProjectManagementProjectRepo.LockMemberWithTx(ctx, tx, successor.ProjectId, actor.OrganizationUserId); err != nil {
			return projectManagementNotFound(err, "project member required")
		}
		stage, err := l.repos.ProjectManagementWorkflowRepo.FindStageWithTx(ctx, tx, successor.StageId)
		if err != nil {
			return projectManagementNotFound(err, "successor task stage not found")
		}
		if stage.Phase != constants.ProjectManagementStagePhasePending {
			return errorx.New("dependency successor task must be pending")
		}
		if _, err = l.repos.ProjectManagementDependencyRepo.FindByTasksWithTx(ctx, tx, successor.Id, predecessor.Id); err == nil {
			return errorx.New("task dependency already exists")
		} else if !isProjectManagementNotFound(err) {
			return err
		}
		cycle, err := l.repos.ProjectManagementDependencyRepo.HasPathWithTx(ctx, tx, predecessor.Id, successor.Id)
		if err != nil {
			return err
		}
		if cycle {
			return errorx.New("task dependency would create a cycle")
		}
		dependency, err = copierx.Trans[dto.ProjectManagementTaskDependencyDto, model.ProjectManagementTaskDependency](&input)
		if err != nil {
			return err
		}
		dependency, err = l.repos.ProjectManagementDependencyRepo.CreateWithTx(ctx, tx, dependency)
		if err != nil {
			return err
		}
		return l.recordProjectDependencyOperation(ctx, tx, actor, successor.ProjectId, dependency, constants.OperationActionProjectManagementDependencyCreate)
	})
	return dependency, err
}

func (l *ProjectManagementLogic) PageTaskDependencies(ctx context.Context, actor *dto.ProjectManagementActor, input dto.ProjectManagementTaskDependencyPageRequestDto) (gormx.PageResponse[*model.ProjectManagementTaskDependency], error) {
	if _, err := l.GetTask(ctx, actor, input.TaskId); err != nil {
		return gormx.PageResponse[*model.ProjectManagementTaskDependency]{}, err
	}
	return l.repos.ProjectManagementDependencyRepo.PageByTask(ctx, input.TaskId, input.Predecessor, gormx.PageRequest{Page: input.Page, Limit: input.Limit})
}

func (l *ProjectManagementLogic) UpdateTaskDependency(ctx context.Context, actor *dto.ProjectManagementActor, id int64, dependencyType constants.ProjectManagementDependencyType) (*model.ProjectManagementTaskDependency, error) {
	if err := validateProjectManagementDependency(dto.ProjectManagementTaskDependencyDto{SuccessorTaskId: 1, PredecessorTaskId: 2, DependencyType: dependencyType}); err != nil {
		return nil, err
	}
	var dependency *model.ProjectManagementTaskDependency
	err := l.repos.ProjectManagementDependencyRepo.Transaction(ctx, func(tx *gorm.DB) error {
		var err error
		dependency, err = l.repos.ProjectManagementDependencyRepo.LockWithTx(ctx, tx, id)
		if err != nil {
			return projectManagementNotFound(err, "task dependency not found")
		}
		tasks, err := l.repos.ProjectManagementDependencyRepo.LockTasksWithTx(ctx, tx, dependency.SuccessorTaskId, dependency.PredecessorTaskId)
		if err != nil {
			return projectManagementNotFound(err, "task not found")
		}
		successor := tasks[0]
		if successor.Id != dependency.SuccessorTaskId {
			successor = tasks[1]
		}
		if successor.OrganizationId != actor.OrganizationId {
			return errorx.New("task dependency not found")
		}
		if _, err = l.repos.ProjectManagementProjectRepo.LockMemberWithTx(ctx, tx, successor.ProjectId, actor.OrganizationUserId); err != nil {
			return projectManagementNotFound(err, "project member required")
		}
		stage, err := l.repos.ProjectManagementWorkflowRepo.FindStageWithTx(ctx, tx, successor.StageId)
		if err != nil {
			return projectManagementNotFound(err, "successor task stage not found")
		}
		if stage.Phase != constants.ProjectManagementStagePhasePending {
			return errorx.New("dependency successor task must be pending")
		}
		if err = l.repos.ProjectManagementDependencyRepo.UpdateTypeWithTx(ctx, tx, dependency, dependencyType); err != nil {
			return err
		}
		dependency.DependencyType = dependencyType
		return l.recordProjectDependencyOperation(ctx, tx, actor, successor.ProjectId, dependency, constants.OperationActionProjectManagementDependencyUpdate)
	})
	return dependency, err
}

func (l *ProjectManagementLogic) DeleteTaskDependency(ctx context.Context, actor *dto.ProjectManagementActor, id int64) error {
	return l.repos.ProjectManagementDependencyRepo.Transaction(ctx, func(tx *gorm.DB) error {
		dependency, err := l.repos.ProjectManagementDependencyRepo.LockWithTx(ctx, tx, id)
		if err != nil {
			return projectManagementNotFound(err, "task dependency not found")
		}
		tasks, err := l.repos.ProjectManagementDependencyRepo.LockTasksWithTx(ctx, tx, dependency.SuccessorTaskId, dependency.PredecessorTaskId)
		if err != nil {
			return projectManagementNotFound(err, "task not found")
		}
		successor := tasks[0]
		if successor.Id != dependency.SuccessorTaskId {
			successor = tasks[1]
		}
		if successor.OrganizationId != actor.OrganizationId {
			return errorx.New("task dependency not found")
		}
		if _, err = l.repos.ProjectManagementProjectRepo.LockMemberWithTx(ctx, tx, successor.ProjectId, actor.OrganizationUserId); err != nil {
			return projectManagementNotFound(err, "project member required")
		}
		if err = l.repos.ProjectManagementDependencyRepo.DeleteWithTx(ctx, tx, dependency); err != nil {
			return err
		}
		return l.recordProjectDependencyOperation(ctx, tx, actor, successor.ProjectId, dependency, constants.OperationActionProjectManagementDependencyDelete)
	})
}

func (l *ProjectManagementLogic) recordProjectDependencyOperation(ctx context.Context, tx *gorm.DB, actor *dto.ProjectManagementActor, projectId int64, dependency *model.ProjectManagementTaskDependency, operation constants.OperationAction) error {
	metadata, err := json.Marshal(map[string]int64{"project_id": projectId, "task_id": dependency.SuccessorTaskId, "dependency_id": dependency.Id})
	if err != nil {
		return errorx.WithStack(err)
	}
	return NewOperationLogic(l.BaseLogic).CreateOperationLog(ctx, tx, dto.OperationCreateDto{EntityId: projectId, EntityType: constants.OperationEntityTypeProjectManagementProject, OperatorId: actor.OrganizationUserId, OperatorName: actor.Name, OperatorType: constants.UserTypeBorrower, Operation: operation, Metadata: metadata, OrganizationId: actor.OrganizationId})
}
