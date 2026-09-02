package logic

import (
	"context"
	"slices"

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

type projectManagementDependencyOperation struct {
	claim      utils.JwtUser
	projectId  int64
	dependency *model.ProjectManagementTaskDependency
	action     constants.OperationAction
	before     map[string]any
	after      map[string]any
}

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

func (l *ProjectManagementDependencyLogic) validateTaskDependencyHierarchyWithTx(ctx context.Context, tx *gorm.DB, taskId, dependencyTaskId int64) error {
	descendantIds, err := l.repos.ProjectManagementTaskRepo.DescendantIdsWithTx(ctx, tx, taskId)
	if err != nil {
		return err
	}
	if slices.Contains(descendantIds, dependencyTaskId) {
		return errorx.New("task dependency cannot be between ancestor and descendant tasks")
	}
	descendantIds, err = l.repos.ProjectManagementTaskRepo.DescendantIdsWithTx(ctx, tx, dependencyTaskId)
	if err != nil {
		return err
	}
	if slices.Contains(descendantIds, taskId) {
		return errorx.New("task dependency cannot be between ancestor and descendant tasks")
	}
	return nil
}

func (l *ProjectManagementDependencyLogic) Create(ctx context.Context, claim utils.JwtUser, input dto.ProjectManagementTaskDependencyDto) (*model.ProjectManagementTaskDependency, error) {
	if err := validateProjectManagementDependency(input); err != nil {
		return nil, err
	}
	var dependency *model.ProjectManagementTaskDependency
	err := l.repos.ProjectManagementDependencyRepo.Transaction(ctx, func(tx *gorm.DB) error {
		var err error
		dependency, err = l.createWithTx(ctx, tx, claim, input)
		return err
	})
	return dependency, err
}

func (l *ProjectManagementDependencyLogic) createWithTx(ctx context.Context, tx *gorm.DB, claim utils.JwtUser, input dto.ProjectManagementTaskDependencyDto) (*model.ProjectManagementTaskDependency, error) {
	tasks, err := l.repos.ProjectManagementDependencyRepo.LockTasksWithTx(ctx, tx, input.TaskId, input.DependencyTaskId)
	if err != nil {
		return nil, projectManagementNotFound(err, "task not found")
	}
	task, dependencyTask := dependencyTasks(tasks, input.TaskId, input.DependencyTaskId)
	if task == nil || dependencyTask == nil {
		return nil, projectManagementNotFound(gorm.ErrRecordNotFound, "task not found")
	}
	successorTaskId, _ := dependencyEndpoints(input.TaskId, input.DependencyTaskId, input.DependencyType)
	successor, predecessor := task, dependencyTask
	if successor.Id != successorTaskId {
		successor, predecessor = predecessor, successor
	}
	if successor.OrganizationId != claim.OrganizationId || predecessor.OrganizationId != claim.OrganizationId || successor.ProjectId != predecessor.ProjectId {
		return nil, errorx.New("dependency tasks must belong to the same project")
	}
	if _, err = l.repos.ProjectManagementProjectRepo.LockMemberWithTx(ctx, tx, successor.ProjectId, claim.OrganizationUserId); err != nil {
		return nil, projectManagementNotFound(err, "project member required")
	}
	if err = l.validateTaskDependencyHierarchyWithTx(ctx, tx, successor.Id, predecessor.Id); err != nil {
		return nil, err
	}
	stage, err := l.repos.ProjectManagementWorkflowRepo.FindStageWithTx(ctx, tx, successor.StageId)
	if err != nil {
		return nil, projectManagementNotFound(err, "successor task stage not found")
	}
	if stage.Phase != constants.ProjectManagementStagePhasePending {
		return nil, errorx.New("dependency successor task must be pending")
	}
	if _, err = l.repos.ProjectManagementDependencyRepo.FindByTasksWithTx(ctx, tx, successor.Id, predecessor.Id); err == nil {
		return nil, errorx.New("task dependency already exists")
	} else if !isProjectManagementNotFound(err) {
		return nil, err
	}
	cycle, err := l.repos.ProjectManagementDependencyRepo.HasPathWithTx(ctx, tx, successor.Id, predecessor.Id)
	if err != nil {
		return nil, err
	}
	if cycle {
		return nil, errorx.New("task dependency would create a cycle")
	}
	dependency, err := copierx.Trans[dto.ProjectManagementTaskDependencyDto, model.ProjectManagementTaskDependency](&input)
	if err != nil {
		return nil, err
	}
	dependency, err = l.repos.ProjectManagementDependencyRepo.CreateWithTx(ctx, tx, dependency)
	if err != nil {
		return nil, err
	}
	return dependency, l.recordOperation(ctx, tx, projectManagementDependencyOperation{claim: claim, projectId: successor.ProjectId, dependency: dependency, action: constants.OperationActionProjectManagementDependencyCreate, after: dependencyOperationIdentity(dependency, predecessor)})
}

func (l *ProjectManagementDependencyLogic) Page(ctx context.Context, claim utils.JwtUser, input dto.ProjectManagementTaskDependencyPageRequestDto) (gormx.PageResponse[*model.ProjectManagementTaskDependency], error) {
	if _, err := NewProjectManagementTaskLogic(l.BaseLogic).Detail(ctx, claim, input.Filter.TaskId); err != nil {
		return gormx.PageResponse[*model.ProjectManagementTaskDependency]{}, err
	}
	return l.repos.ProjectManagementDependencyRepo.PageByTask(ctx, input.Filter.TaskId, input.Filter.Predecessor, gormx.PageRequest{Page: input.Page, Limit: input.Limit})
}

func (l *ProjectManagementDependencyLogic) Update(ctx context.Context, claim utils.JwtUser, input dto.ProjectManagementTaskDependencyDto) (*model.ProjectManagementTaskDependency, error) {
	var dependency *model.ProjectManagementTaskDependency
	err := l.repos.ProjectManagementDependencyRepo.Transaction(ctx, func(tx *gorm.DB) error {
		var err error
		dependency, err = l.updateWithTx(ctx, tx, claim, input)
		return err
	})
	return dependency, err
}

func (l *ProjectManagementDependencyLogic) updateWithTx(ctx context.Context, tx *gorm.DB, claim utils.JwtUser, input dto.ProjectManagementTaskDependencyDto) (*model.ProjectManagementTaskDependency, error) {
	var dependency *model.ProjectManagementTaskDependency
	var err error
	dependency, err = l.repos.ProjectManagementDependencyRepo.LockWithTx(ctx, tx, input.Id)
	if err != nil {
		return nil, projectManagementNotFound(err, "task dependency not found")
	}
	if dependency.TaskId != input.TaskId && dependency.DependencyTaskId != input.TaskId {
		return nil, errorx.New("task dependency not found")
	}
	oldTasks, err := l.repos.ProjectManagementDependencyRepo.LockTasksWithTx(ctx, tx, dependency.TaskId, dependency.DependencyTaskId)
	if err != nil {
		return nil, projectManagementNotFound(err, "task not found")
	}
	oldTask, oldRelatedTask := dependencyTasks(oldTasks, dependency.TaskId, dependency.DependencyTaskId)
	if oldTask == nil || oldRelatedTask == nil {
		return nil, projectManagementNotFound(gorm.ErrRecordNotFound, "task not found")
	}
	_, oldPredecessor := dependencyEndpoints(dependency.TaskId, dependency.DependencyTaskId, dependency.DependencyType)
	oldPredecessorTask := oldRelatedTask
	if oldPredecessor == oldTask.Id {
		oldPredecessorTask = oldTask
	}
	before := dependencyOperationValues(dependency, oldPredecessorTask)
	update := dto.ProjectManagementTaskDependencyDto{TaskId: input.TaskId, DependencyTaskId: input.DependencyTaskId, DependencyType: input.DependencyType, DependencyRule: input.DependencyRule}
	if err = validateProjectManagementDependency(update); err != nil {
		return nil, err
	}
	tasks, err := l.repos.ProjectManagementDependencyRepo.LockTasksWithTx(ctx, tx, input.TaskId, input.DependencyTaskId)
	if err != nil {
		return nil, projectManagementNotFound(err, "task not found")
	}
	task, relatedTask := dependencyTasks(tasks, input.TaskId, input.DependencyTaskId)
	if task == nil || relatedTask == nil {
		return nil, projectManagementNotFound(gorm.ErrRecordNotFound, "task not found")
	}
	successorTaskId, _ := dependencyEndpoints(input.TaskId, input.DependencyTaskId, input.DependencyType)
	successor, predecessor := task, relatedTask
	if successor.Id != successorTaskId {
		successor, predecessor = predecessor, successor
	}
	if successor == nil || predecessor == nil || successor.OrganizationId != claim.OrganizationId || predecessor.OrganizationId != claim.OrganizationId || successor.ProjectId != predecessor.ProjectId {
		return nil, errorx.New("task dependency not found")
	}
	if _, err = l.repos.ProjectManagementProjectRepo.LockMemberWithTx(ctx, tx, successor.ProjectId, claim.OrganizationUserId); err != nil {
		return nil, projectManagementNotFound(err, "project member required")
	}
	if err = l.validateTaskDependencyHierarchyWithTx(ctx, tx, successor.Id, predecessor.Id); err != nil {
		return nil, err
	}
	stage, err := l.repos.ProjectManagementWorkflowRepo.FindStageWithTx(ctx, tx, successor.StageId)
	if err != nil {
		return nil, projectManagementNotFound(err, "successor task stage not found")
	}
	if stage.Phase != constants.ProjectManagementStagePhasePending {
		return nil, errorx.New("dependency successor task must be pending")
	}
	if existing, findErr := l.repos.ProjectManagementDependencyRepo.FindByTasksWithTx(ctx, tx, successor.Id, predecessor.Id); findErr == nil && existing.Id != dependency.Id {
		return nil, errorx.New("task dependency already exists")
	} else if findErr != nil && !isProjectManagementNotFound(findErr) {
		return nil, findErr
	}
	cycle, err := l.repos.ProjectManagementDependencyRepo.HasPathWithTx(ctx, tx, successor.Id, predecessor.Id, dependency.Id)
	if err != nil {
		return nil, err
	}
	if cycle {
		return nil, errorx.New("task dependency would create a cycle")
	}
	dependency.TaskId = input.TaskId
	dependency.DependencyTaskId = input.DependencyTaskId
	dependency.DependencyType = input.DependencyType
	dependency.DependencyRule = input.DependencyRule
	if err = l.repos.ProjectManagementDependencyRepo.UpdateWithTx(ctx, tx, dependency); err != nil {
		return nil, err
	}
	after := dependencyOperationValues(dependency, predecessor)
	before, after = operationValuesDiff(before, after)
	if len(before) == 0 {
		return dependency, nil
	}
	return dependency, l.recordOperation(ctx, tx, projectManagementDependencyOperation{claim: claim, projectId: successor.ProjectId, dependency: dependency, action: constants.OperationActionProjectManagementDependencyUpdate, before: before, after: after})
}

func (l *ProjectManagementDependencyLogic) Delete(ctx context.Context, claim utils.JwtUser, id int64) error {
	return l.repos.ProjectManagementDependencyRepo.Transaction(ctx, func(tx *gorm.DB) error {
		return l.deleteWithTx(ctx, tx, claim, id)
	})
}

func (l *ProjectManagementDependencyLogic) deleteWithTx(ctx context.Context, tx *gorm.DB, claim utils.JwtUser, id int64) error {
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
	return l.recordOperation(ctx, tx, projectManagementDependencyOperation{claim: claim, projectId: successor.ProjectId, dependency: dependency, action: constants.OperationActionProjectManagementDependencyDelete, before: dependencyOperationIdentity(dependency, predecessor)})
}

func (l *ProjectManagementDependencyLogic) SyncWithTx(ctx context.Context, tx *gorm.DB, claim utils.JwtUser, taskId int64, inputs []*dto.ProjectManagementTaskDependencyDto) error {
	existing, err := l.repos.ProjectManagementDependencyRepo.ListByTaskWithTx(ctx, tx, taskId)
	if err != nil {
		return err
	}
	existingById := make(map[int64]*model.ProjectManagementTaskDependency, len(existing))
	for _, dependency := range existing {
		if dependency.TaskId != taskId {
			continue
		}
		existingById[dependency.Id] = dependency
	}
	requested := make(map[int64]struct{}, len(inputs))
	for _, input := range inputs {
		if input == nil || input.TaskId != taskId {
			return errorx.New("dependency task id must match task")
		}
		if input.Id == 0 {
			if err := validateProjectManagementDependency(*input); err != nil {
				return err
			}
			if _, err := l.createWithTx(ctx, tx, claim, *input); err != nil {
				return err
			}
			continue
		}
		existingDependency, ok := existingById[input.Id]
		if !ok {
			return projectManagementNotFound(gorm.ErrRecordNotFound, "task dependency not found")
		}
		if _, ok := requested[input.Id]; ok {
			return errorx.New("duplicate task dependency")
		}
		requested[input.Id] = struct{}{}
		if existingDependency.DependencyTaskId != input.DependencyTaskId || existingDependency.DependencyType != input.DependencyType || existingDependency.DependencyRule != input.DependencyRule {
			return errorx.New("existing task dependency cannot be modified")
		}
	}
	for id := range existingById {
		if _, ok := requested[id]; ok {
			continue
		}
		if err := l.deleteWithTx(ctx, tx, claim, id); err != nil {
			return err
		}
	}
	return nil
}

func (l *ProjectManagementDependencyLogic) recordOperation(ctx context.Context, tx *gorm.DB, input projectManagementDependencyOperation) error {
	successorTaskId, _ := dependencyEndpoints(input.dependency.TaskId, input.dependency.DependencyTaskId, input.dependency.DependencyType)
	metadata, err := marshalProjectManagementOperationMetadata(projectManagementOperationMetadata{ProjectID: input.projectId, TaskID: successorTaskId, DependencyID: input.dependency.Id}, input.before, input.after)
	if err != nil {
		return errorx.WithStack(err)
	}
	return NewOperationLogic(l.BaseLogic).CreateOperationLog(ctx, tx, dto.OperationCreateDto{EntityId: input.projectId, EntityType: constants.OperationEntityTypeProjectManagementProject, OperatorId: input.claim.OrganizationUserId, OperatorName: input.claim.Username, OperatorType: constants.UserTypeBorrower, Operation: input.action, Metadata: metadata, OrganizationId: input.claim.OrganizationId})
}

func dependencyOperationValues(dependency *model.ProjectManagementTaskDependency, dependencyTask *model.ProjectManagementTask) map[string]any {
	return map[string]any{"dependency_task_id": dependencyTask.Id, "dependency_task_name": dependencyTask.Title, "dependency_type": dependency.DependencyType, "dependency_rule": dependency.DependencyRule}
}

func dependencyOperationIdentity(dependency *model.ProjectManagementTaskDependency, dependencyTask *model.ProjectManagementTask) map[string]any {
	return map[string]any{"id": dependency.Id, "name": dependencyTask.Title}
}
