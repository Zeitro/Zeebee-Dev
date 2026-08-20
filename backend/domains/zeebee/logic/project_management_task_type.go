package logic

import (
	"context"
	"encoding/json"
	"strings"

	"gorm.io/gorm"

	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/gormx"
	"zeebee-backend/common/utils"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
)

type ProjectManagementTaskTypeLogic struct{ *BaseLogic }

func NewProjectManagementTaskTypeLogic(baseLogic *BaseLogic) *ProjectManagementTaskTypeLogic {
	return &ProjectManagementTaskTypeLogic{BaseLogic: baseLogic}
}

func ProjectManagementTaskTypeSnapshot(taskType *model.ProjectManagementTaskType) *dto.ProjectManagementTaskTypeSnapshotDto {
	return &dto.ProjectManagementTaskTypeSnapshotDto{
		Id:          taskType.Id,
		Name:        taskType.Name,
		Description: taskType.Description,
		Color:       taskType.Color,
	}
}

func (l *ProjectManagementTaskTypeLogic) Create(ctx context.Context, claim *utils.JwtUser, input dto.ProjectManagementTaskTypeDto) (taskType *model.ProjectManagementTaskType, err error) {
	if err = validateProjectManagementTaskType(input); err != nil {
		return nil, err
	}
	data, err := copierx.Trans[dto.ProjectManagementTaskTypeDto, model.ProjectManagementTaskType](&input)
	if err != nil {
		return nil, err
	}
	data.OrganizationId = claim.OrganizationId
	err = l.repos.ProjectManagementTaskTypeRepo.Transaction(ctx, func(tx *gorm.DB) error {
		taskType, err = l.repos.ProjectManagementTaskTypeRepo.CreateWithTx(ctx, tx, data, false)
		if err != nil {
			return err
		}
		return l.recordOperation(ctx, tx, claim, taskType.Id, constants.OperationActionProjectManagementTaskTypeCreate)
	})
	return taskType, err
}

func (l *ProjectManagementTaskTypeLogic) Page(ctx context.Context, organizationId int64, input dto.ProjectManagementTaskTypePageRequestDto) (gormx.PageResponse[*model.ProjectManagementTaskType], error) {
	return l.repos.ProjectManagementTaskTypeRepo.Page(ctx, organizationId, input)
}

func (l *ProjectManagementTaskTypeLogic) Update(ctx context.Context, claim *utils.JwtUser, id int64, input dto.ProjectManagementTaskTypeDto) (*model.ProjectManagementTaskType, error) {
	if err := validateProjectManagementTaskType(input); err != nil {
		return nil, err
	}
	if _, err := l.repos.ProjectManagementTaskTypeRepo.FindByOrganization(ctx, id, claim.OrganizationId); err != nil {
		return nil, projectManagementNotFound(err, "task type not found")
	}
	err := l.repos.ProjectManagementTaskTypeRepo.Transaction(ctx, func(tx *gorm.DB) error {
		if err := l.repos.ProjectManagementTaskTypeRepo.UpdateByOrganization(ctx, id, claim.OrganizationId, &model.ProjectManagementTaskType{
			Name: input.Name, Description: input.Description, Color: input.Color, Enabled: input.Enabled,
		}); err != nil {
			return err
		}
		return l.recordOperation(ctx, tx, claim, id, constants.OperationActionProjectManagementTaskTypeUpdate)
	})
	if err != nil {
		return nil, err
	}
	taskType, err := l.repos.ProjectManagementTaskTypeRepo.FindByOrganization(ctx, id, claim.OrganizationId)
	if err != nil {
		return nil, projectManagementNotFound(err, "task type not found")
	}
	return taskType, nil
}

func (l *ProjectManagementTaskTypeLogic) Delete(ctx context.Context, claim *utils.JwtUser, id int64) error {
	if _, err := l.repos.ProjectManagementTaskTypeRepo.FindByOrganization(ctx, id, claim.OrganizationId); err != nil {
		return projectManagementNotFound(err, "task type not found")
	}
	return l.repos.ProjectManagementTaskTypeRepo.Transaction(ctx, func(tx *gorm.DB) error {
		if err := l.repos.ProjectManagementTaskTypeRepo.DeleteByOrganization(ctx, id, claim.OrganizationId); err != nil {
			return err
		}
		return l.recordOperation(ctx, tx, claim, id, constants.OperationActionProjectManagementTaskTypeDelete)
	})
}

func (l *ProjectManagementTaskTypeLogic) RequireEnabled(ctx context.Context, id, organizationId int64) (*model.ProjectManagementTaskType, error) {
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

func (l *ProjectManagementTaskTypeLogic) recordOperation(ctx context.Context, tx *gorm.DB, claim *utils.JwtUser, taskTypeId int64, action constants.OperationAction) error {
	metadata, err := json.Marshal(map[string]int64{"task_type_id": taskTypeId})
	if err != nil {
		return errorx.WithStack(err)
	}
	return NewOperationLogic(l.BaseLogic).CreateOperationLog(ctx, tx, dto.OperationCreateDto{
		EntityId:       claim.OrganizationId,
		EntityType:     constants.OperationEntityTypeProjectManagementOrganization,
		OperatorId:     claim.OrganizationUserId,
		OperatorName:   claim.Username,
		OperatorType:   constants.UserTypeBorrower,
		Operation:      action,
		Metadata:       metadata,
		OrganizationId: claim.OrganizationId,
	})
}
