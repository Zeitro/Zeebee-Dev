package logic

import (
	"context"
	"strings"

	"gorm.io/gorm"

	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/utils"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
)

type ProjectManagementTaskTypeLogic struct{ *BaseLogic }

type projectManagementTaskTypeOperation struct {
	claim      utils.JwtUser
	taskTypeId int64
	action     constants.OperationAction
	before     map[string]any
	after      map[string]any
}

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

func (l *ProjectManagementTaskTypeLogic) Create(ctx context.Context, claim utils.JwtUser, input dto.ProjectManagementTaskTypeDto) (taskType *model.ProjectManagementTaskType, err error) {
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
		return l.recordOperation(ctx, tx, projectManagementTaskTypeOperation{claim: claim, taskTypeId: taskType.Id, action: constants.OperationActionProjectManagementTaskTypeCreate, after: map[string]any{"id": taskType.Id, "name": taskType.Name}})
	})
	return taskType, err
}

func (l *ProjectManagementTaskTypeLogic) All(ctx context.Context, organizationId int64) ([]*model.ProjectManagementTaskType, error) {
	taskTypes, err := l.repos.ProjectManagementTaskTypeRepo.ListByOrganization(ctx, organizationId)
	if err != nil || len(taskTypes) != 0 {
		return taskTypes, err
	}
	if err := l.repos.ProjectManagementTaskTypeRepo.Transaction(ctx, func(tx *gorm.DB) error {
		existing, err := l.repos.ProjectManagementTaskTypeRepo.ListByOrganizationWithTx(ctx, tx, organizationId)
		if err != nil || len(existing) != 0 {
			return err
		}
		templates, err := l.repos.ProjectManagementTaskTypeRepo.ListByOrganizationWithTx(ctx, tx, 0)
		if err != nil {
			return err
		}
		if len(templates) == 0 {
			return errorx.New("global task type templates not found")
		}
		copies := make([]*model.ProjectManagementTaskType, len(templates))
		for i, template := range templates {
			clone := *template
			clone.Base = model.Base{}
			clone.OrganizationId = organizationId
			copies[i] = &clone
		}
		return errorx.WithStack(tx.WithContext(ctx).Create(&copies).Error)
	}); err != nil {
		return nil, err
	}
	return l.repos.ProjectManagementTaskTypeRepo.ListByOrganization(ctx, organizationId)
}

func (l *ProjectManagementTaskTypeLogic) Detail(ctx context.Context, organizationId, id int64) (*model.ProjectManagementTaskType, error) {
	taskType, err := l.repos.ProjectManagementTaskTypeRepo.FindByOrganization(ctx, id, organizationId)
	if err != nil {
		return nil, projectManagementNotFound(err, "task type not found")
	}
	return taskType, nil
}

func (l *ProjectManagementTaskTypeLogic) Update(ctx context.Context, claim utils.JwtUser, id int64, input dto.ProjectManagementTaskTypeDto) (*model.ProjectManagementTaskType, error) {
	if err := validateProjectManagementTaskType(input); err != nil {
		return nil, err
	}
	err := l.repos.ProjectManagementTaskTypeRepo.Transaction(ctx, func(tx *gorm.DB) error {
		taskType, innerErr := l.repos.ProjectManagementTaskTypeRepo.FindByOrganizationWithTx(ctx, tx, id, claim.OrganizationId)
		if innerErr != nil {
			return projectManagementNotFound(innerErr, "task type not found")
		}
		before := taskTypeOperationValues(taskType)
		input.Id = taskType.Id
		if innerErr = copierx.Assign(&input, taskType); innerErr != nil {
			return innerErr
		}
		if _, innerErr = l.repos.ProjectManagementTaskTypeRepo.UpdateOneByIdWithTx(ctx, tx, taskType.Id, taskType); innerErr != nil {
			return innerErr
		}
		after := taskTypeOperationValues(taskType)
		before, after = operationValuesDiff(before, after)
		if len(before) == 0 {
			return nil
		}
		return l.recordOperation(ctx, tx, projectManagementTaskTypeOperation{claim: claim, taskTypeId: id, action: constants.OperationActionProjectManagementTaskTypeUpdate, before: before, after: after})
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

func (l *ProjectManagementTaskTypeLogic) Order(ctx context.Context, claim utils.JwtUser, id int64, sortOrder int) error {
	if sortOrder < 0 {
		return errorx.New("task type sort order must be >= 0")
	}
	return l.repos.ProjectManagementTaskTypeRepo.Transaction(ctx, func(tx *gorm.DB) error {
		taskType, err := l.repos.ProjectManagementTaskTypeRepo.FindByOrganizationWithTx(ctx, tx, id, claim.OrganizationId)
		if err != nil {
			return projectManagementNotFound(err, "task type not found")
		}
		if taskType.SortOrder == sortOrder {
			return nil
		}
		if err := l.repos.ProjectManagementTaskTypeRepo.UpdateSortOrderWithTx(ctx, tx, id, claim.OrganizationId, sortOrder); err != nil {
			return err
		}
		return l.recordOperation(ctx, tx, projectManagementTaskTypeOperation{claim: claim, taskTypeId: id, action: constants.OperationActionProjectManagementTaskTypeUpdate, before: map[string]any{"sort_order": taskType.SortOrder}, after: map[string]any{"sort_order": sortOrder}})
	})
}

func (l *ProjectManagementTaskTypeLogic) Delete(ctx context.Context, claim utils.JwtUser, id int64) error {
	return l.repos.ProjectManagementTaskTypeRepo.Transaction(ctx, func(tx *gorm.DB) error {
		taskTypes, err := l.repos.ProjectManagementTaskTypeRepo.LockByOrganizationWithTx(ctx, tx, claim.OrganizationId)
		if err != nil {
			return err
		}
		var target *model.ProjectManagementTaskType
		for _, taskType := range taskTypes {
			if taskType.Id == id {
				target = taskType
				break
			}
		}
		if target == nil {
			return errorx.New("task type not found")
		}
		if len(taskTypes) <= 1 {
			return errorx.New("at least one task type is required")
		}
		if err := l.repos.ProjectManagementTaskTypeRepo.DeleteByOrganizationWithTx(ctx, tx, id, claim.OrganizationId); err != nil {
			return err
		}
		return l.recordOperation(ctx, tx, projectManagementTaskTypeOperation{claim: claim, taskTypeId: id, action: constants.OperationActionProjectManagementTaskTypeDelete, before: taskTypeOperationValues(target)})
	})
}

func validateProjectManagementTaskType(input dto.ProjectManagementTaskTypeDto) error {
	if strings.TrimSpace(input.Name) == "" {
		return errorx.New("task type name is required")
	}
	return nil
}

func (l *ProjectManagementTaskTypeLogic) recordOperation(ctx context.Context, tx *gorm.DB, input projectManagementTaskTypeOperation) error {
	metadata, err := marshalProjectManagementOperationMetadata(projectManagementOperationMetadata{TaskTypeID: input.taskTypeId}, input.before, input.after)
	if err != nil {
		return errorx.WithStack(err)
	}
	return NewOperationLogic(l.BaseLogic).CreateOperationLog(ctx, tx, dto.OperationCreateDto{
		EntityId:       input.claim.OrganizationId,
		EntityType:     constants.OperationEntityTypeProjectManagementOrganization,
		OperatorId:     input.claim.OrganizationUserId,
		OperatorName:   input.claim.Username,
		OperatorType:   constants.UserTypeBorrower,
		Operation:      input.action,
		Metadata:       metadata,
		OrganizationId: input.claim.OrganizationId,
	})
}

func taskTypeOperationValues(taskType *model.ProjectManagementTaskType) map[string]any {
	return map[string]any{"name": taskType.Name, "description": taskType.Description, "color": taskType.Color, "sort_order": taskType.SortOrder}
}
