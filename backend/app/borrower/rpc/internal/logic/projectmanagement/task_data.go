package projectmanagementlogic

import (
	"encoding/json"
	"strings"
	"time"

	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/data/dto"
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

func taskEntity(task *model.ProjectManagementTask) (*borrower.ProjectManagementTaskEntity, error) {
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
	if task.TaskType != nil && task.TaskType.Id != 0 {
		entity.TaskType, err = copierx.Trans[dto.ProjectManagementTaskTypeSnapshotDto, borrower.ProjectManagementTaskTypeSnapshotEntity](task.TaskType)
		if err != nil {
			return nil, err
		}
	}
	assignees, err := copierx.Trans[[]*model.ProjectManagementTaskAssignee, []*borrower.ProjectManagementTaskAssigneeEntity](&task.Assignees)
	if err != nil {
		return nil, err
	}
	entity.Assignees = *assignees
	return entity, nil
}

func attachmentEntity(file *model.File) (*borrower.ProjectManagementAttachmentEntity, error) {
	entity, err := copierx.Trans[model.File, borrower.ProjectManagementAttachmentEntity](file)
	if err != nil {
		return nil, err
	}
	entity.FileId = file.Id
	return entity, nil
}

func dependencyEntity(dependency *model.ProjectManagementTaskDependency) (*borrower.ProjectManagementTaskDependencyEntity, error) {
	return copierx.Trans[model.ProjectManagementTaskDependency, borrower.ProjectManagementTaskDependencyEntity](dependency)
}

func operationEntity(operation *model.Operation) (*borrower.ProjectManagementOperationEntity, error) {
	entity, err := copierx.Trans[model.Operation, borrower.ProjectManagementOperationEntity](operation)
	if err != nil {
		return nil, err
	}
	entity.ProjectId, entity.CreatedAt = operation.EntityId, operation.CreatedAt.Format(time.RFC3339)
	var metadata map[string]int64
	if len(operation.Metadata) > 0 && json.Unmarshal(operation.Metadata, &metadata) == nil {
		entity.TaskId = metadata["task_id"]
	}
	return entity, nil
}
