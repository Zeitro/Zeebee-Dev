package projectmanagementtasklogic

import (
	"context"
	"strings"
	"time"

	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/data/dto"
	domainlogic "zeebee-backend/domains/zeebee/logic"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/proto"
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

func taskEntity(task *model.ProjectManagementTask, overdue *dto.ProjectManagementTaskOverdueDto) (*proto.ProjectManagementTaskEntity, error) {
	entity, err := copierx.Trans[model.ProjectManagementTask, proto.ProjectManagementTaskEntity](task)
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
		entity.TaskType, err = copierx.Trans[dto.ProjectManagementTaskTypeSnapshotDto, proto.ProjectManagementTaskTypeSnapshotEntity](task.TaskType)
		if err != nil {
			return nil, err
		}
	}
	if overdue != nil {
		entity.Overdue, entity.OverdueSeconds = overdue.Overdue, overdue.OverdueSeconds
	}
	assignees, err := copierx.Trans[[]*model.ProjectManagementTaskAssignee, []*proto.ProjectManagementTaskAssigneeEntity](&task.Assignees)
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

func taskEntityWithOverdue(ctx context.Context, taskLogic *domainlogic.ProjectManagementTaskLogic, task *model.ProjectManagementTask) (*proto.ProjectManagementTaskEntity, error) {
	overdue, err := taskLogic.Overdue(ctx, task)
	if err != nil {
		return nil, err
	}
	return taskEntity(task, overdue)
}

func attachmentEntity(file *model.File) (*proto.FileEntity, error) {
	entity, err := copierx.Trans[model.File, proto.FileEntity](file)
	if err != nil {
		return nil, err
	}
	entity.Id = file.Id
	entity.CreatedAt, entity.UpdatedAt = formatEntityTime(file.CreatedAt), formatEntityTime(file.UpdatedAt)
	return entity, nil
}
