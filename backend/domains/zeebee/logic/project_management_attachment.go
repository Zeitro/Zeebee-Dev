package logic

import (
	"context"
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

type ProjectManagementAttachmentLogic struct{ *BaseLogic }

func NewProjectManagementAttachmentLogic(baseLogic *BaseLogic) *ProjectManagementAttachmentLogic {
	return &ProjectManagementAttachmentLogic{BaseLogic: baseLogic}
}

func (l *ProjectManagementAttachmentLogic) Page(ctx context.Context, claim utils.JwtUser, taskId int64, input dto.ProjectManagementAttachmentPageRequestDto) (gormx.PageResponse[*model.File], error) {
	if _, err := NewProjectManagementTaskLogic(l.BaseLogic).Detail(ctx, claim, taskId); err != nil {
		return gormx.PageResponse[*model.File]{}, err
	}
	return l.repos.FileRepo.PageByEntity(ctx, taskId, constants.FileEntityTypeProjectManagementTask, gormx.PageRequest{Page: input.Page, Limit: input.Limit})
}

func (l *ProjectManagementAttachmentLogic) Create(ctx context.Context, claim utils.JwtUser, input dto.ProjectManagementAttachmentCreateDto) (file *model.File, err error) {
	if input.TaskId == 0 || strings.TrimSpace(input.FileUrl) == "" {
		return nil, errorx.New("task and file url are required")
	}
	if input.StorageType == "" {
		input.StorageType = constants.StorageTypeAwsS3
	}
	if !input.StorageType.IsValid() {
		return nil, errorx.New("invalid attachment storage type")
	}
	err = l.repos.FileRepo.Transaction(ctx, func(tx *gorm.DB) error {
		task, err := l.repos.ProjectManagementTaskRepo.LockWithTx(ctx, tx, input.TaskId)
		if err != nil {
			return projectManagementNotFound(err, "task not found")
		}
		if task.OrganizationId != claim.OrganizationId {
			return errorx.New("task not found")
		}
		if _, err = l.repos.ProjectManagementProjectRepo.LockMemberWithTx(ctx, tx, task.ProjectId, claim.OrganizationUserId); err != nil {
			return projectManagementNotFound(err, "project member required")
		}
		file, err = copierx.Trans[dto.ProjectManagementAttachmentCreateDto, model.File](&input)
		if err != nil {
			return err
		}
		file.EntityId, file.EntityType = task.Id, constants.FileEntityTypeProjectManagementTask
		file, err = l.repos.FileRepo.CreateWithTx(ctx, tx, file, false)
		if err != nil {
			return err
		}
		return NewProjectManagementTaskLogic(l.BaseLogic).recordProjectTaskOperation(ctx, tx, projectManagementTaskOperation{claim: claim, projectId: task.ProjectId, taskId: task.Id, action: constants.OperationActionProjectManagementAttachmentAdd})
	})
	return file, err
}

func (l *ProjectManagementAttachmentLogic) Delete(ctx context.Context, claim utils.JwtUser, taskId, fileId int64) error {
	return l.repos.FileRepo.Transaction(ctx, func(tx *gorm.DB) error {
		task, err := l.repos.ProjectManagementTaskRepo.LockWithTx(ctx, tx, taskId)
		if err != nil {
			return projectManagementNotFound(err, "task not found")
		}
		if task.OrganizationId != claim.OrganizationId {
			return errorx.New("task not found")
		}
		if _, err = l.repos.ProjectManagementProjectRepo.LockMemberWithTx(ctx, tx, task.ProjectId, claim.OrganizationUserId); err != nil {
			return projectManagementNotFound(err, "project member required")
		}
		file, err := l.repos.FileRepo.FindOneByIdWithTx(ctx, tx, fileId)
		if err != nil {
			return projectManagementNotFound(err, "task attachment not found")
		}
		if file.EntityId != task.Id || file.EntityType != constants.FileEntityTypeProjectManagementTask {
			return errorx.New("task attachment not found")
		}
		if _, err = l.repos.FileRepo.DeleteByIDWithTx(ctx, tx, file.Id); err != nil {
			return err
		}
		return NewProjectManagementTaskLogic(l.BaseLogic).recordProjectTaskOperation(ctx, tx, projectManagementTaskOperation{claim: claim, projectId: task.ProjectId, taskId: task.Id, action: constants.OperationActionProjectManagementAttachmentDelete})
	})
}
