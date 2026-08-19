package projectmanagementlogic

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/data/dto"
	domainlogic "zeebee-backend/domains/zeebee/logic"
	"zeebee-backend/domains/zeebee/model"
)

func taskTime(value string) (*time.Time, error) {
	if strings.TrimSpace(value) == "" {
		return nil, nil
	}
	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return nil, errorx.New("task time must use RFC3339")
	}
	return &parsed, nil
}

func formatEntityTime(value time.Time) string {
	if value.IsZero() {
		return ""
	}
	return value.Format(time.RFC3339)
}

func taskEntity(task *model.ProjectManagementTask, overdue *dto.ProjectManagementTaskOverdueDto) (*borrower.ProjectManagementTaskEntity, error) {
	entity, err := copierx.Trans[model.ProjectManagementTask, borrower.ProjectManagementTaskEntity](task)
	if err != nil {
		return nil, err
	}
	if task.StartAt != nil && task.StartAt.Valid {
		entity.StartAt = task.StartAt.Time.Format(time.RFC3339)
	}
	if task.EndAt != nil && task.EndAt.Valid {
		entity.EndAt = task.EndAt.Time.Format(time.RFC3339)
	}
	entity.CreatedAt, entity.UpdatedAt = formatEntityTime(task.CreatedAt), formatEntityTime(task.UpdatedAt)
	if task.TaskType != nil && task.TaskType.Id != 0 {
		entity.TaskType, err = copierx.Trans[dto.ProjectManagementTaskTypeSnapshotDto, borrower.ProjectManagementTaskTypeSnapshotEntity](task.TaskType)
		if err != nil {
			return nil, err
		}
	}
	if overdue != nil {
		entity.Overdue, entity.OverdueSeconds = overdue.Overdue, overdue.OverdueSeconds
	}
	assignees, err := copierx.Trans[[]*model.ProjectManagementTaskAssignee, []*borrower.ProjectManagementTaskAssigneeEntity](&task.Assignees)
	if err != nil {
		return nil, err
	}
	entity.Assignees = *assignees
	for index, assignee := range task.Assignees {
		entity.Assignees[index].CreatedAt = formatEntityTime(assignee.CreatedAt)
		entity.Assignees[index].UpdatedAt = formatEntityTime(assignee.UpdatedAt)
	}
	return entity, nil
}

func taskEntityWithOverdue(ctx context.Context, taskLogic *domainlogic.ProjectManagementLogic, task *model.ProjectManagementTask) (*borrower.ProjectManagementTaskEntity, error) {
	overdue, err := taskLogic.TaskOverdue(ctx, task)
	if err != nil {
		return nil, err
	}
	return taskEntity(task, overdue)
}

func attachmentEntity(file *model.File) (*borrower.ProjectManagementAttachmentEntity, error) {
	entity, err := copierx.Trans[model.File, borrower.ProjectManagementAttachmentEntity](file)
	if err != nil {
		return nil, err
	}
	entity.FileId = file.Id
	entity.CreatedAt, entity.UpdatedAt = formatEntityTime(file.CreatedAt), formatEntityTime(file.UpdatedAt)
	return entity, nil
}

func dependencyEntity(dependency *model.ProjectManagementTaskDependency) (*borrower.ProjectManagementTaskDependencyEntity, error) {
	entity, err := copierx.Trans[model.ProjectManagementTaskDependency, borrower.ProjectManagementTaskDependencyEntity](dependency)
	if err != nil {
		return nil, err
	}
	entity.CreatedAt, entity.UpdatedAt = formatEntityTime(dependency.CreatedAt), formatEntityTime(dependency.UpdatedAt)
	return entity, nil
}

func operationEntity(operation *model.Operation) (*borrower.ProjectManagementOperationEntity, error) {
	entity, err := copierx.Trans[model.Operation, borrower.ProjectManagementOperationEntity](operation)
	if err != nil {
		return nil, err
	}
	entity.ProjectId = operation.EntityId
	entity.CreatedAt, entity.UpdatedAt = formatEntityTime(operation.CreatedAt), formatEntityTime(operation.UpdatedAt)
	var metadata map[string]int64
	if len(operation.Metadata) > 0 && json.Unmarshal(operation.Metadata, &metadata) == nil {
		entity.TaskId = metadata["task_id"]
	}
	return entity, nil
}
