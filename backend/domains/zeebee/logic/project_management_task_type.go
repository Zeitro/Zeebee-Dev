package logic

import (
	"context"
	"encoding/json"
	"strings"

	"gorm.io/datatypes"
	"gorm.io/gorm"

	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/gormx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
)

func ProjectManagementTaskTypeSnapshot(taskType *model.ProjectManagementTaskType) (datatypes.JSON, error) {
	return json.Marshal(dto.ProjectManagementTaskTypeSnapshotDto{
		Id:          taskType.Id,
		Name:        taskType.Name,
		Description: taskType.Description,
		Color:       taskType.Color,
	})
}

func (l *ProjectManagementLogic) CreateTaskType(ctx context.Context, actor *dto.ProjectManagementActor, input dto.ProjectManagementTaskTypeDto) (taskType *model.ProjectManagementTaskType, err error) {
	if err = validateProjectManagementTaskType(input); err != nil {
		return nil, err
	}
	data, err := copierx.Trans[dto.ProjectManagementTaskTypeDto, model.ProjectManagementTaskType](&input)
	if err != nil {
		return nil, err
	}
	data.OrganizationId = actor.OrganizationId
	err = l.repos.ProjectManagementTaskTypeRepo.Transaction(ctx, func(tx *gorm.DB) error {
		taskType, err = l.repos.ProjectManagementTaskTypeRepo.CreateWithTx(ctx, tx, data, false)
		if err != nil {
			return err
		}
		return l.recordTaskTypeOperation(ctx, tx, actor, taskType.Id, constants.OperationActionProjectManagementTaskTypeCreate)
	})
	return taskType, err
}

func (l *ProjectManagementLogic) PageTaskTypes(ctx context.Context, organizationId int64, input dto.ProjectManagementTaskTypePageRequestDto) (gormx.PageResponse[*model.ProjectManagementTaskType], error) {
	return l.repos.ProjectManagementTaskTypeRepo.Page(ctx, organizationId, input)
}

func (l *ProjectManagementLogic) UpdateTaskType(ctx context.Context, actor *dto.ProjectManagementActor, id int64, input dto.ProjectManagementTaskTypeDto) (*model.ProjectManagementTaskType, error) {
	if err := validateProjectManagementTaskType(input); err != nil {
		return nil, err
	}
	if _, err := l.repos.ProjectManagementTaskTypeRepo.FindByOrganization(ctx, id, actor.OrganizationId); err != nil {
		return nil, projectManagementNotFound(err, "task type not found")
	}
	err := l.repos.ProjectManagementTaskTypeRepo.Transaction(ctx, func(tx *gorm.DB) error {
		if err := l.repos.ProjectManagementTaskTypeRepo.UpdateByOrganization(ctx, id, actor.OrganizationId, map[string]interface{}{
			"name": input.Name, "description": input.Description, "color": input.Color, "enabled": input.Enabled,
		}); err != nil {
			return err
		}
		return l.recordTaskTypeOperation(ctx, tx, actor, id, constants.OperationActionProjectManagementTaskTypeUpdate)
	})
	if err != nil {
		return nil, err
	}
	taskType, err := l.repos.ProjectManagementTaskTypeRepo.FindByOrganization(ctx, id, actor.OrganizationId)
	if err != nil {
		return nil, projectManagementNotFound(err, "task type not found")
	}
	return taskType, nil
}

func (l *ProjectManagementLogic) DeleteTaskType(ctx context.Context, actor *dto.ProjectManagementActor, id int64) error {
	if _, err := l.repos.ProjectManagementTaskTypeRepo.FindByOrganization(ctx, id, actor.OrganizationId); err != nil {
		return projectManagementNotFound(err, "task type not found")
	}
	return l.repos.ProjectManagementTaskTypeRepo.Transaction(ctx, func(tx *gorm.DB) error {
		if err := l.repos.ProjectManagementTaskTypeRepo.DeleteByOrganization(ctx, id, actor.OrganizationId); err != nil {
			return err
		}
		return l.recordTaskTypeOperation(ctx, tx, actor, id, constants.OperationActionProjectManagementTaskTypeDelete)
	})
}

func (l *ProjectManagementLogic) RequireEnabledTaskType(ctx context.Context, id, organizationId int64) (*model.ProjectManagementTaskType, error) {
	taskType, err := l.repos.ProjectManagementTaskTypeRepo.FindEnabledByOrganization(ctx, id, organizationId)
	if err != nil {
		return nil, projectManagementNotFound(err, "enabled task type not found")
	}
	return taskType, nil
}

func validateProjectManagementTaskType(input dto.ProjectManagementTaskTypeDto) error {
	if strings.TrimSpace(input.Name) == "" {
		return errorx.New("task type name is required")
	}
	return nil
}

func (l *ProjectManagementLogic) recordTaskTypeOperation(ctx context.Context, tx *gorm.DB, actor *dto.ProjectManagementActor, taskTypeId int64, action constants.OperationAction) error {
	metadata, err := json.Marshal(map[string]int64{"task_type_id": taskTypeId})
	if err != nil {
		return errorx.WithStack(err)
	}
	return NewOperationLogic(l.BaseLogic).CreateOperationLog(ctx, tx, dto.OperationCreateDto{
		EntityId:       actor.OrganizationId,
		EntityType:     constants.OperationEntityTypeProjectManagementOrganization,
		OperatorId:     actor.OrganizationUserId,
		OperatorName:   actor.Name,
		OperatorType:   constants.UserTypeBorrower,
		Operation:      action,
		Metadata:       metadata,
		OrganizationId: actor.OrganizationId,
	})
}
