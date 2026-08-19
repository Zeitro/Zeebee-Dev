package logic

import (
	"context"
	"database/sql"
	"encoding/json"
	"strings"
	"time"

	"gorm.io/gorm"

	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/gormx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
)

func (l *ProjectManagementLogic) CreateTask(ctx context.Context, actor *dto.ProjectManagementActor, input dto.ProjectManagementTaskCreateDto) (task *model.ProjectManagementTask, err error) {
	if err = validateProjectManagementTaskCreate(input); err != nil {
		return nil, err
	}
	if _, err = l.RequireMember(ctx, input.ProjectId, actor.OrganizationUserId); err != nil {
		return nil, err
	}
	taskData := input.ProjectManagementTaskDto
	taskData.StartAt, taskData.EndAt, taskData.TaskType = nil, nil, nil
	task, err = copierx.Trans[dto.ProjectManagementTaskDto, model.ProjectManagementTask](&taskData)
	if err != nil {
		return nil, err
	}
	task.OrganizationId, task.CreatorId, task.CreatorName = actor.OrganizationId, actor.OrganizationUserId, actor.Name
	task.StartAt, task.EndAt = projectManagementNullTime(input.StartAt), projectManagementNullTime(input.EndAt)
	task.TaskType = &dto.ProjectManagementTaskTypeSnapshotDto{}
	err = l.repos.ProjectManagementTaskRepo.Transaction(ctx, func(tx *gorm.DB) error {
		if _, err = l.repos.ProjectManagementProjectRepo.LockByOrganizationWithTx(ctx, tx, input.ProjectId, actor.OrganizationId); err != nil {
			return projectManagementNotFound(err, "project not found")
		}
		workflow, err := l.repos.ProjectManagementWorkflowRepo.FindProjectDefinitionWithTx(ctx, tx, input.ProjectId, actor.OrganizationId)
		if err != nil {
			return projectManagementNotFound(err, "project workflow not found")
		}
		pendingStage, err := l.repos.ProjectManagementWorkflowRepo.FindPendingStageWithTx(ctx, tx, workflow.Id)
		if err != nil {
			return projectManagementNotFound(err, "project pending stage not found")
		}
		if input.ParentTaskId != 0 {
			parent, err := l.repos.ProjectManagementTaskRepo.LockWithTx(ctx, tx, input.ParentTaskId)
			if err != nil {
				return projectManagementNotFound(err, "parent task not found")
			}
			if parent.ProjectId != input.ProjectId {
				return errorx.New("parent task must belong to the project")
			}
			parentStage, err := l.repos.ProjectManagementWorkflowRepo.FindStageWithTx(ctx, tx, parent.StageId)
			if err != nil {
				return projectManagementNotFound(err, "parent task stage not found")
			}
			if parentStage.Phase == constants.ProjectManagementStagePhaseFinished {
				return errorx.New("cannot add child task to a finished parent task")
			}
		}
		task.StageId = pendingStage.Id
		if input.TaskTypeId != 0 {
			taskType, err := l.repos.ProjectManagementTaskTypeRepo.FindEnabledByOrganizationWithTx(ctx, tx, input.TaskTypeId, actor.OrganizationId)
			if err != nil {
				return projectManagementNotFound(err, "enabled task type not found")
			}
			task.TaskType = ProjectManagementTaskTypeSnapshot(taskType)
		}
		task, err = l.repos.ProjectManagementTaskRepo.CreateWithTx(ctx, tx, task, false)
		if err != nil {
			return err
		}
		assignees, err := l.projectManagementTaskAssignees(ctx, tx, input.ProjectId, task.Id, actor.OrganizationId, input.AssigneeIds)
		if err != nil {
			return err
		}
		if err = l.repos.ProjectManagementTaskRepo.CreateAssigneesWithTx(ctx, tx, assignees); err != nil {
			return err
		}
		return l.recordProjectTaskOperation(ctx, tx, actor, task.ProjectId, task.Id, constants.OperationActionProjectManagementTaskCreate)
	})
	if err != nil {
		return nil, err
	}
	return l.GetTask(ctx, actor, task.Id)
}

func validateProjectManagementTaskCreate(input dto.ProjectManagementTaskCreateDto) error {
	if input.ProjectId == 0 || strings.TrimSpace(input.Title) == "" {
		return errorx.New("project and task title are required")
	}
	return nil
}

func (l *ProjectManagementLogic) GetTask(ctx context.Context, actor *dto.ProjectManagementActor, taskId int64) (*model.ProjectManagementTask, error) {
	task, err := l.repos.ProjectManagementTaskRepo.FindDetailByOrganization(ctx, taskId, actor.OrganizationId)
	if err != nil {
		return nil, projectManagementNotFound(err, "task not found")
	}
	if _, err = l.RequireMember(ctx, task.ProjectId, actor.OrganizationUserId); err != nil {
		return nil, err
	}
	if err = l.fillProjectManagementTaskOverdue(ctx, []*model.ProjectManagementTask{task}); err != nil {
		return nil, err
	}
	return task, nil
}

func (l *ProjectManagementLogic) PageProjectTasks(ctx context.Context, actor *dto.ProjectManagementActor, input dto.ProjectManagementTaskPageRequestDto) (gormx.PageResponse[*model.ProjectManagementTask], error) {
	if _, err := l.RequireProject(ctx, input.Filter.ProjectId, actor.OrganizationId); err != nil {
		return gormx.PageResponse[*model.ProjectManagementTask]{}, err
	}
	if _, err := l.RequireMember(ctx, input.Filter.ProjectId, actor.OrganizationUserId); err != nil {
		return gormx.PageResponse[*model.ProjectManagementTask]{}, err
	}
	page, err := l.repos.ProjectManagementTaskRepo.PageByProject(ctx, actor.OrganizationId, input)
	if err != nil {
		return page, err
	}
	if err = l.fillProjectManagementTaskOverdue(ctx, page.List); err != nil {
		return gormx.PageResponse[*model.ProjectManagementTask]{}, err
	}
	return page, nil
}

func (l *ProjectManagementLogic) UpdateTask(ctx context.Context, actor *dto.ProjectManagementActor, taskId int64, input dto.ProjectManagementTaskUpdateDto) (*model.ProjectManagementTask, error) {
	if strings.TrimSpace(input.Title) == "" {
		return nil, errorx.New("task title is required")
	}
	var task *model.ProjectManagementTask
	err := l.repos.ProjectManagementTaskRepo.Transaction(ctx, func(tx *gorm.DB) error {
		var err error
		task, err = l.repos.ProjectManagementTaskRepo.LockWithTx(ctx, tx, taskId)
		if err != nil {
			return projectManagementNotFound(err, "task not found")
		}
		if task.OrganizationId != actor.OrganizationId {
			return errorx.New("task not found")
		}
		if _, err = l.repos.ProjectManagementProjectRepo.LockMemberWithTx(ctx, tx, task.ProjectId, actor.OrganizationUserId); err != nil {
			return projectManagementNotFound(err, "project member required")
		}
		task.Title, task.Description, task.Priority = input.Title, input.Description, input.Priority
		task.StartAt, task.EndAt = projectManagementNullTime(input.StartAt), projectManagementNullTime(input.EndAt)
		if input.TaskTypeId != 0 {
			taskType, err := l.repos.ProjectManagementTaskTypeRepo.FindEnabledByOrganizationWithTx(ctx, tx, input.TaskTypeId, actor.OrganizationId)
			if err != nil {
				return projectManagementNotFound(err, "enabled task type not found")
			}
			task.TaskType = ProjectManagementTaskTypeSnapshot(taskType)
		}
		if err = tx.WithContext(ctx).Model(&model.ProjectManagementTask{}).Where("id = ?", task.Id).Updates(map[string]interface{}{"title": task.Title, "description": task.Description, "priority": task.Priority, "start_at": task.StartAt, "end_at": task.EndAt, "task_type": task.TaskType}).Error; err != nil {
			return errorx.WithStack(err)
		}
		return l.recordProjectTaskOperation(ctx, tx, actor, task.ProjectId, task.Id, constants.OperationActionProjectManagementTaskUpdate)
	})
	if err != nil {
		return nil, err
	}
	return l.GetTask(ctx, actor, task.Id)
}

func (l *ProjectManagementLogic) ReplaceTaskAssignees(ctx context.Context, actor *dto.ProjectManagementActor, taskId int64, organizationUserIds []int64) (*model.ProjectManagementTask, error) {
	err := l.repos.ProjectManagementTaskRepo.Transaction(ctx, func(tx *gorm.DB) error {
		task, err := l.repos.ProjectManagementTaskRepo.LockWithTx(ctx, tx, taskId)
		if err != nil {
			return projectManagementNotFound(err, "task not found")
		}
		if task.OrganizationId != actor.OrganizationId {
			return errorx.New("task not found")
		}
		if _, err = l.repos.ProjectManagementProjectRepo.LockMemberWithTx(ctx, tx, task.ProjectId, actor.OrganizationUserId); err != nil {
			return projectManagementNotFound(err, "project member required")
		}
		assignees, err := l.projectManagementTaskAssignees(ctx, tx, task.ProjectId, task.Id, actor.OrganizationId, organizationUserIds)
		if err != nil {
			return err
		}
		if err = l.repos.ProjectManagementTaskRepo.ReplaceAssigneesWithTx(ctx, tx, task.Id, assignees); err != nil {
			return err
		}
		return l.recordProjectTaskOperation(ctx, tx, actor, task.ProjectId, task.Id, constants.OperationActionProjectManagementTaskAssignees)
	})
	if err != nil {
		return nil, err
	}
	return l.GetTask(ctx, actor, taskId)
}

func (l *ProjectManagementLogic) DeleteTask(ctx context.Context, actor *dto.ProjectManagementActor, taskId int64, confirm bool) error {
	if err := validateProjectManagementTaskDelete(confirm); err != nil {
		return err
	}
	return l.repos.ProjectManagementTaskRepo.Transaction(ctx, func(tx *gorm.DB) error {
		task, err := l.repos.ProjectManagementTaskRepo.LockWithTx(ctx, tx, taskId)
		if err != nil {
			return projectManagementNotFound(err, "task not found")
		}
		if task.OrganizationId != actor.OrganizationId {
			return errorx.New("task not found")
		}
		if _, err = l.repos.ProjectManagementProjectRepo.LockMemberWithTx(ctx, tx, task.ProjectId, actor.OrganizationUserId); err != nil {
			return projectManagementNotFound(err, "project member required")
		}
		taskIds, err := l.repos.ProjectManagementTaskRepo.DescendantIdsWithTx(ctx, tx, taskId)
		if err != nil {
			return err
		}
		if err = l.repos.FileRepo.DeleteByEntityIdWithTx(ctx, tx, taskIds, constants.FileEntityTypeProjectManagementTask); err != nil {
			return err
		}
		if err = l.repos.ProjectManagementTaskRepo.DeleteByIdsWithTx(ctx, tx, taskIds); err != nil {
			return err
		}
		return l.recordProjectTaskOperation(ctx, tx, actor, task.ProjectId, task.Id, constants.OperationActionProjectManagementTaskDelete)
	})
}

func validateProjectManagementTaskDelete(confirm bool) error {
	if !confirm {
		return errorx.New("task deletion requires confirmation")
	}
	return nil
}

func projectManagementNullTime(value *time.Time) *sql.NullTime {
	if value == nil {
		return nil
	}
	return &sql.NullTime{Time: *value, Valid: true}
}

func (l *ProjectManagementLogic) projectManagementTaskAssignees(ctx context.Context, tx *gorm.DB, projectId, taskId, organizationId int64, organizationUserIds []int64) ([]*model.ProjectManagementTaskAssignee, error) {
	result := make([]*model.ProjectManagementTaskAssignee, 0, len(organizationUserIds))
	seen := make(map[int64]struct{}, len(organizationUserIds))
	for _, organizationUserId := range organizationUserIds {
		if _, exists := seen[organizationUserId]; exists {
			return nil, errorx.New("duplicate task assignee")
		}
		seen[organizationUserId] = struct{}{}
		if _, err := l.repos.ProjectManagementProjectRepo.LockMemberWithTx(ctx, tx, projectId, organizationUserId); err != nil {
			return nil, projectManagementNotFound(err, "task assignee must be a project member")
		}
		organizationUser, err := l.repos.OrganizationUserRepo.GetOrganizationUserByIdWithTx(ctx, tx, organizationUserId, organizationId)
		if err != nil {
			return nil, projectManagementNotFound(err, "organization member not found")
		}
		if organizationUser.Status != constants.OrganizationUserStatusActive {
			return nil, errorx.New("task assignee must be an active organization member")
		}
		result = append(result, &model.ProjectManagementTaskAssignee{TaskId: taskId, OrganizationUserId: organizationUserId, OrganizationUserName: organizationUser.FullName})
	}
	return result, nil
}

func (l *ProjectManagementLogic) recordProjectTaskOperation(ctx context.Context, tx *gorm.DB, actor *dto.ProjectManagementActor, projectId, taskId int64, operation constants.OperationAction) error {
	metadata, err := json.Marshal(map[string]int64{"project_id": projectId, "task_id": taskId})
	if err != nil {
		return errorx.WithStack(err)
	}
	return NewOperationLogic(l.BaseLogic).CreateOperationLog(ctx, tx, dto.OperationCreateDto{EntityId: projectId, EntityType: constants.OperationEntityTypeProjectManagementProject, OperatorId: actor.OrganizationUserId, OperatorName: actor.Name, OperatorType: constants.UserTypeBorrower, Operation: operation, Metadata: metadata, OrganizationId: actor.OrganizationId})
}
