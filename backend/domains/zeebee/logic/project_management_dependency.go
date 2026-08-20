package logic

import (
	"context"
	"encoding/json"

	"gorm.io/gorm"

	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/gormx"
	"zeebee-backend/common/utils"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
)

type ProjectManagementDependencyLogic struct{ *BaseLogic }

func NewProjectManagementDependencyLogic(baseLogic *BaseLogic) *ProjectManagementDependencyLogic {
	return &ProjectManagementDependencyLogic{BaseLogic: baseLogic}
}

func validateProjectManagementDependency(input dto.ProjectManagementTaskDependencyDto) error {
	if input.TaskId == 0 || input.DependencyTaskId == 0 {
		return errorx.New("task and dependency task are required")
	}
	if input.TaskId == input.DependencyTaskId {
		return errorx.New("task cannot depend on itself")
	}
	if input.DependencyType != constants.ProjectManagementDependencyTypePreTask && input.DependencyType != constants.ProjectManagementDependencyTypePostTask {
		return errorx.New("invalid task dependency direction")
	}
	if input.DependencyRule != constants.ProjectManagementDependencyRuleFS && input.DependencyRule != constants.ProjectManagementDependencyRuleFF && input.DependencyRule != constants.ProjectManagementDependencyRuleSS && input.DependencyRule != constants.ProjectManagementDependencyRuleSF {
		return errorx.New("invalid task dependency rule")
	}
	return nil
}

func dependencyEndpoints(taskId, dependencyTaskId int64, dependencyType constants.ProjectManagementDependencyType) (successorTaskId, predecessorTaskId int64) {
	if dependencyType == constants.ProjectManagementDependencyTypePreTask {
		return taskId, dependencyTaskId
	}
	return dependencyTaskId, taskId
}

func dependencyTasks(tasks []*model.ProjectManagementTask, taskId, dependencyTaskId int64) (task, dependencyTask *model.ProjectManagementTask) {
	for _, item := range tasks {
		switch item.Id {
		case taskId:
			task = item
		case dependencyTaskId:
			dependencyTask = item
		}
	}
	return task, dependencyTask
}

func (l *ProjectManagementDependencyLogic) Create(ctx context.Context, claim *utils.JwtUser, input dto.ProjectManagementTaskDependencyDto) (*model.ProjectManagementTaskDependency, error) {
	if err := validateProjectManagementDependency(input); err != nil {
		return nil, err
	}
	var dependency *model.ProjectManagementTaskDependency
	err := l.repos.ProjectManagementDependencyRepo.Transaction(ctx, func(tx *gorm.DB) error {
		tasks, err := l.repos.ProjectManagementDependencyRepo.LockTasksWithTx(ctx, tx, input.TaskId, input.DependencyTaskId)
		if err != nil {
			return projectManagementNotFound(err, "task not found")
		}
		task, dependencyTask := dependencyTasks(tasks, input.TaskId, input.DependencyTaskId)
		if task == nil || dependencyTask == nil {
			return projectManagementNotFound(gorm.ErrRecordNotFound, "task not found")
		}
		successorTaskId, _ := dependencyEndpoints(input.TaskId, input.DependencyTaskId, input.DependencyType)
		successor, predecessor := task, dependencyTask
		if successor.Id != successorTaskId {
			successor, predecessor = predecessor, successor
		}
		if successor.OrganizationId != claim.OrganizationId || predecessor.OrganizationId != claim.OrganizationId || successor.ProjectId != predecessor.ProjectId {
			return errorx.New("dependency tasks must belong to the same project")
		}
		if _, err = l.repos.ProjectManagementProjectRepo.LockMemberWithTx(ctx, tx, successor.ProjectId, claim.OrganizationUserId); err != nil {
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
		cycle, err := l.repos.ProjectManagementDependencyRepo.HasPathWithTx(ctx, tx, successor.Id, predecessor.Id)
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
		return l.recordOperation(ctx, tx, claim, successor.ProjectId, dependency, constants.OperationActionProjectManagementDependencyCreate)
	})
	return dependency, err
}

func (l *ProjectManagementDependencyLogic) Page(ctx context.Context, claim *utils.JwtUser, input dto.ProjectManagementTaskDependencyPageRequestDto) (gormx.PageResponse[*model.ProjectManagementTaskDependency], error) {
	if _, err := NewProjectManagementTaskLogic(l.BaseLogic).Detail(ctx, claim, input.TaskId); err != nil {
		return gormx.PageResponse[*model.ProjectManagementTaskDependency]{}, err
	}
	return l.repos.ProjectManagementDependencyRepo.PageByTask(ctx, input.TaskId, input.Predecessor, gormx.PageRequest{Page: input.Page, Limit: input.Limit})
}

func (l *ProjectManagementDependencyLogic) Update(ctx context.Context, claim *utils.JwtUser, taskId, id int64, dependencyType constants.ProjectManagementDependencyType, dependencyRule constants.ProjectManagementDependencyRule) (*model.ProjectManagementTaskDependency, error) {
	var dependency *model.ProjectManagementTaskDependency
	err := l.repos.ProjectManagementDependencyRepo.Transaction(ctx, func(tx *gorm.DB) error {
		var err error
		dependency, err = l.repos.ProjectManagementDependencyRepo.LockWithTx(ctx, tx, id)
		if err != nil {
			return projectManagementNotFound(err, "task dependency not found")
		}
		if dependency.TaskId != taskId && dependency.DependencyTaskId != taskId {
			return errorx.New("task dependency not found")
		}
		dependencyTaskId := dependency.DependencyTaskId
		if dependency.TaskId != taskId {
			dependencyTaskId = dependency.TaskId
		}
		input := dto.ProjectManagementTaskDependencyDto{TaskId: taskId, DependencyTaskId: dependencyTaskId, DependencyType: dependencyType, DependencyRule: dependencyRule}
		if err = validateProjectManagementDependency(input); err != nil {
			return err
		}
		tasks, err := l.repos.ProjectManagementDependencyRepo.LockTasksWithTx(ctx, tx, taskId, dependencyTaskId)
		if err != nil {
			return projectManagementNotFound(err, "task not found")
		}
		task, relatedTask := dependencyTasks(tasks, taskId, dependencyTaskId)
		if task == nil || relatedTask == nil {
			return projectManagementNotFound(gorm.ErrRecordNotFound, "task not found")
		}
		successorTaskId, _ := dependencyEndpoints(taskId, dependencyTaskId, dependencyType)
		successor, predecessor := task, relatedTask
		if successor.Id != successorTaskId {
			successor, predecessor = predecessor, successor
		}
		if successor == nil || predecessor == nil || successor.OrganizationId != claim.OrganizationId || predecessor.OrganizationId != claim.OrganizationId || successor.ProjectId != predecessor.ProjectId {
			return errorx.New("task dependency not found")
		}
		if _, err = l.repos.ProjectManagementProjectRepo.LockMemberWithTx(ctx, tx, successor.ProjectId, claim.OrganizationUserId); err != nil {
			return projectManagementNotFound(err, "project member required")
		}
		stage, err := l.repos.ProjectManagementWorkflowRepo.FindStageWithTx(ctx, tx, successor.StageId)
		if err != nil {
			return projectManagementNotFound(err, "successor task stage not found")
		}
		if stage.Phase != constants.ProjectManagementStagePhasePending {
			return errorx.New("dependency successor task must be pending")
		}
		if existing, findErr := l.repos.ProjectManagementDependencyRepo.FindByTasksWithTx(ctx, tx, successor.Id, predecessor.Id); findErr == nil && existing.Id != dependency.Id {
			return errorx.New("task dependency already exists")
		} else if findErr != nil && !isProjectManagementNotFound(findErr) {
			return findErr
		}
		cycle, err := l.repos.ProjectManagementDependencyRepo.HasPathWithTx(ctx, tx, successor.Id, predecessor.Id, dependency.Id)
		if err != nil {
			return err
		}
		if cycle {
			return errorx.New("task dependency would create a cycle")
		}
		dependency.TaskId = taskId
		dependency.DependencyTaskId = dependencyTaskId
		dependency.DependencyType = dependencyType
		dependency.DependencyRule = dependencyRule
		if err = l.repos.ProjectManagementDependencyRepo.UpdateWithTx(ctx, tx, dependency); err != nil {
			return err
		}
		return l.recordOperation(ctx, tx, claim, successor.ProjectId, dependency, constants.OperationActionProjectManagementDependencyUpdate)
	})
	return dependency, err
}

func (l *ProjectManagementDependencyLogic) Delete(ctx context.Context, claim *utils.JwtUser, id int64) error {
	return l.repos.ProjectManagementDependencyRepo.Transaction(ctx, func(tx *gorm.DB) error {
		dependency, err := l.repos.ProjectManagementDependencyRepo.LockWithTx(ctx, tx, id)
		if err != nil {
			return projectManagementNotFound(err, "task dependency not found")
		}
		successorTaskId, _ := dependencyEndpoints(dependency.TaskId, dependency.DependencyTaskId, dependency.DependencyType)
		tasks, err := l.repos.ProjectManagementDependencyRepo.LockTasksWithTx(ctx, tx, dependency.TaskId, dependency.DependencyTaskId)
		if err != nil {
			return projectManagementNotFound(err, "task not found")
		}
		task, dependencyTask := dependencyTasks(tasks, dependency.TaskId, dependency.DependencyTaskId)
		if task == nil || dependencyTask == nil {
			return projectManagementNotFound(gorm.ErrRecordNotFound, "task not found")
		}
		successor, predecessor := task, dependencyTask
		if successor.Id != successorTaskId {
			successor, predecessor = predecessor, successor
		}
		if successor == nil || predecessor == nil || successor.OrganizationId != claim.OrganizationId {
			return errorx.New("task dependency not found")
		}
		if _, err = l.repos.ProjectManagementProjectRepo.LockMemberWithTx(ctx, tx, successor.ProjectId, claim.OrganizationUserId); err != nil {
			return projectManagementNotFound(err, "project member required")
		}
		if err = l.repos.ProjectManagementDependencyRepo.DeleteWithTx(ctx, tx, dependency); err != nil {
			return err
		}
		return l.recordOperation(ctx, tx, claim, successor.ProjectId, dependency, constants.OperationActionProjectManagementDependencyDelete)
	})
}

func (l *ProjectManagementDependencyLogic) recordOperation(ctx context.Context, tx *gorm.DB, claim *utils.JwtUser, projectId int64, dependency *model.ProjectManagementTaskDependency, operation constants.OperationAction) error {
	successorTaskId, _ := dependencyEndpoints(dependency.TaskId, dependency.DependencyTaskId, dependency.DependencyType)
	metadata, err := json.Marshal(map[string]int64{"project_id": projectId, "task_id": successorTaskId, "dependency_id": dependency.Id})
	if err != nil {
		return errorx.WithStack(err)
	}
	return NewOperationLogic(l.BaseLogic).CreateOperationLog(ctx, tx, dto.OperationCreateDto{EntityId: projectId, EntityType: constants.OperationEntityTypeProjectManagementProject, OperatorId: claim.OrganizationUserId, OperatorName: claim.Username, OperatorType: constants.UserTypeBorrower, Operation: operation, Metadata: metadata, OrganizationId: claim.OrganizationId})
}
