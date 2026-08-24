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
	"zeebee-backend/common/utils"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
)

type ProjectManagementTaskLogic struct{ *BaseLogic }

type projectManagementTaskAssigneeInput struct {
	projectId           int64
	taskId              int64
	stageId             int64
	organizationId      int64
	organizationUserIds []int64
}

type projectManagementTaskOperation struct {
	claim     utils.JwtUser
	projectId int64
	taskId    int64
	action    constants.OperationAction
}

func NewProjectManagementTaskLogic(baseLogic *BaseLogic) *ProjectManagementTaskLogic {
	return &ProjectManagementTaskLogic{BaseLogic: baseLogic}
}

func (l *ProjectManagementTaskLogic) Create(ctx context.Context, claim utils.JwtUser, input dto.ProjectManagementTaskCreateDto) (task *model.ProjectManagementTask, err error) {
	if err = validateProjectManagementTaskCreate(input); err != nil {
		return nil, err
	}
	if _, err = l.projectManagementAccess().RequireMember(ctx, input.ProjectId, claim.OrganizationUserId); err != nil {
		return nil, err
	}
	taskData := input.ProjectManagementTaskDto
	taskData.StartAt, taskData.EndAt, taskData.TaskType = nil, nil, nil
	task, err = copierx.Trans[dto.ProjectManagementTaskDto, model.ProjectManagementTask](&taskData)
	if err != nil {
		return nil, err
	}
	task.OrganizationId, task.CreatorId, task.CreatorName = claim.OrganizationId, claim.OrganizationUserId, claim.Username
	task.StartAt, task.EndAt = projectManagementNullTime(input.StartAt), projectManagementNullTime(input.EndAt)
	task.TaskType = &dto.ProjectManagementTaskTypeSnapshotDto{}
	err = l.repos.ProjectManagementTaskRepo.Transaction(ctx, func(tx *gorm.DB) error {
		if _, err = l.repos.ProjectManagementProjectRepo.LockByOrganizationWithTx(ctx, tx, input.ProjectId, claim.OrganizationId); err != nil {
			return projectManagementNotFound(err, "project not found")
		}
		workflow, err := l.repos.ProjectManagementWorkflowRepo.FindProjectDefinitionWithTx(ctx, tx, input.ProjectId, claim.OrganizationId)
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
			taskType, err := l.repos.ProjectManagementTaskTypeRepo.FindByOrganizationWithTx(ctx, tx, input.TaskTypeId, claim.OrganizationId)
			if err != nil {
				return projectManagementNotFound(err, "task type not found")
			}
			task.TaskType = ProjectManagementTaskTypeSnapshot(taskType)
		}
		task, err = l.repos.ProjectManagementTaskRepo.CreateWithTx(ctx, tx, task, false)
		if err != nil {
			return err
		}
		assignees, err := l.projectManagementTaskAssignees(ctx, tx, projectManagementTaskAssigneeInput{projectId: input.ProjectId, taskId: task.Id, stageId: task.StageId, organizationId: claim.OrganizationId, organizationUserIds: input.AssigneeIds})
		if err != nil {
			return err
		}
		if err = l.repos.ProjectManagementTaskRepo.CreateAssigneesWithTx(ctx, tx, assignees); err != nil {
			return err
		}
		return l.recordProjectTaskOperation(ctx, tx, projectManagementTaskOperation{claim: claim, projectId: task.ProjectId, taskId: task.Id, action: constants.OperationActionProjectManagementTaskCreate})
	})
	if err != nil {
		return nil, err
	}
	return l.Detail(ctx, claim, task.Id)
}

func (l *ProjectManagementTaskLogic) DetailForProject(ctx context.Context, claim utils.JwtUser, projectId, taskId int64) (*model.ProjectManagementTask, error) {
	task, err := l.Detail(ctx, claim, taskId)
	if err != nil {
		return nil, err
	}
	if task.ProjectId != projectId {
		return nil, errorx.New("task not found")
	}
	return task, nil
}

func validateProjectManagementTaskCreate(input dto.ProjectManagementTaskCreateDto) error {
	if input.ProjectId == 0 || strings.TrimSpace(input.Title) == "" {
		return errorx.New("project and task title are required")
	}
	return nil
}

func (l *ProjectManagementTaskLogic) Detail(ctx context.Context, claim utils.JwtUser, taskId int64) (*model.ProjectManagementTask, error) {
	task, err := l.repos.ProjectManagementTaskRepo.FindDetailByOrganization(ctx, taskId, claim.OrganizationId)
	if err != nil {
		return nil, projectManagementNotFound(err, "task not found")
	}
	if _, err = l.projectManagementAccess().RequireMember(ctx, task.ProjectId, claim.OrganizationUserId); err != nil {
		return nil, err
	}
	return task, nil
}

func (l *ProjectManagementTaskLogic) Page(ctx context.Context, claim utils.JwtUser, input dto.ProjectManagementTaskPageRequestDto) (gormx.PageResponse[*model.ProjectManagementTask], error) {
	if _, err := l.projectManagementAccess().RequireProject(ctx, input.Filter.ProjectId, claim.OrganizationId); err != nil {
		return gormx.PageResponse[*model.ProjectManagementTask]{}, err
	}
	if _, err := l.projectManagementAccess().RequireMember(ctx, input.Filter.ProjectId, claim.OrganizationUserId); err != nil {
		return gormx.PageResponse[*model.ProjectManagementTask]{}, err
	}
	page, err := l.repos.ProjectManagementTaskRepo.PageByProject(ctx, claim.OrganizationId, input)
	if err != nil {
		return page, err
	}
	return page, nil
}

func (l *ProjectManagementTaskLogic) Update(ctx context.Context, claim utils.JwtUser, taskId int64, input dto.ProjectManagementTaskUpdateDto) (*model.ProjectManagementTask, error) {
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
		if task.OrganizationId != claim.OrganizationId {
			return errorx.New("task not found")
		}
		if _, err = l.repos.ProjectManagementProjectRepo.LockMemberWithTx(ctx, tx, task.ProjectId, claim.OrganizationUserId); err != nil {
			return projectManagementNotFound(err, "project member required")
		}
		task.Title, task.Description, task.Priority = input.Title, input.Description, input.Priority
		task.StartAt, task.EndAt = projectManagementNullTime(input.StartAt), projectManagementNullTime(input.EndAt)
		if input.TaskTypeId != 0 {
			taskType, err := l.repos.ProjectManagementTaskTypeRepo.FindByOrganizationWithTx(ctx, tx, input.TaskTypeId, claim.OrganizationId)
			if err != nil {
				return projectManagementNotFound(err, "task type not found")
			}
			task.TaskType = ProjectManagementTaskTypeSnapshot(taskType)
		}
		if err = tx.WithContext(ctx).Model(&model.ProjectManagementTask{}).Where("id = ?", task.Id).Updates(map[string]interface{}{"title": task.Title, "description": task.Description, "priority": task.Priority, "start_at": task.StartAt, "end_at": task.EndAt, "task_type": task.TaskType}).Error; err != nil {
			return errorx.WithStack(err)
		}
		return l.recordProjectTaskOperation(ctx, tx, projectManagementTaskOperation{claim: claim, projectId: task.ProjectId, taskId: task.Id, action: constants.OperationActionProjectManagementTaskUpdate})
	})
	if err != nil {
		return nil, err
	}
	return l.Detail(ctx, claim, task.Id)
}

func (l *ProjectManagementTaskLogic) UpdateAssignees(ctx context.Context, claim utils.JwtUser, taskId, stageId int64, organizationUserIds []int64) (*model.ProjectManagementTask, error) {
	err := l.repos.ProjectManagementTaskRepo.Transaction(ctx, func(tx *gorm.DB) error {
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
		workflow, err := l.repos.ProjectManagementWorkflowRepo.FindProjectDefinitionWithTx(ctx, tx, task.ProjectId, claim.OrganizationId)
		if err != nil {
			return projectManagementNotFound(err, "project workflow not found")
		}
		stage, err := l.repos.ProjectManagementWorkflowRepo.FindStageWithTx(ctx, tx, stageId)
		if err != nil {
			return projectManagementNotFound(err, "workflow stage not found")
		}
		if stage.WorkflowId != workflow.Id {
			return errorx.New("task assignee stage must belong to the project workflow")
		}
		assignees, err := l.projectManagementTaskAssignees(ctx, tx, projectManagementTaskAssigneeInput{projectId: task.ProjectId, taskId: task.Id, stageId: stageId, organizationId: claim.OrganizationId, organizationUserIds: organizationUserIds})
		if err != nil {
			return err
		}
		if err = l.repos.ProjectManagementTaskRepo.ReplaceAssigneesWithTx(ctx, tx, task.Id, stageId, assignees); err != nil {
			return err
		}
		return l.recordProjectTaskOperation(ctx, tx, projectManagementTaskOperation{claim: claim, projectId: task.ProjectId, taskId: task.Id, action: constants.OperationActionProjectManagementTaskAssignees})
	})
	if err != nil {
		return nil, err
	}
	return l.Detail(ctx, claim, taskId)
}

func (l *ProjectManagementTaskLogic) Delete(ctx context.Context, claim utils.JwtUser, taskId int64, confirm bool) error {
	if err := validateProjectManagementTaskDelete(confirm); err != nil {
		return err
	}
	return l.repos.ProjectManagementTaskRepo.Transaction(ctx, func(tx *gorm.DB) error {
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
		return l.recordProjectTaskOperation(ctx, tx, projectManagementTaskOperation{claim: claim, projectId: task.ProjectId, taskId: task.Id, action: constants.OperationActionProjectManagementTaskDelete})
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

func (l *ProjectManagementTaskLogic) projectManagementTaskAssignees(ctx context.Context, tx *gorm.DB, input projectManagementTaskAssigneeInput) ([]*model.ProjectManagementTaskAssignee, error) {
	result := make([]*model.ProjectManagementTaskAssignee, 0, len(input.organizationUserIds))
	seen := make(map[int64]struct{}, len(input.organizationUserIds))
	for _, organizationUserId := range input.organizationUserIds {
		if _, exists := seen[organizationUserId]; exists {
			return nil, errorx.New("duplicate task assignee")
		}
		seen[organizationUserId] = struct{}{}
		if _, err := l.repos.ProjectManagementProjectRepo.LockMemberWithTx(ctx, tx, input.projectId, organizationUserId); err != nil {
			return nil, projectManagementNotFound(err, "task assignee must be a project member")
		}
		organizationUser, err := l.repos.OrganizationUserRepo.GetOrganizationUserByIdWithTx(ctx, tx, organizationUserId, input.organizationId)
		if err != nil {
			return nil, projectManagementNotFound(err, "organization member not found")
		}
		if organizationUser.Status != constants.OrganizationUserStatusActive {
			return nil, errorx.New("task assignee must be an active organization member")
		}
		result = append(result, &model.ProjectManagementTaskAssignee{TaskId: input.taskId, StageId: input.stageId, OrganizationUserId: organizationUserId, OrganizationUserName: organizationUser.FullName})
	}
	return result, nil
}

func (l *ProjectManagementTaskLogic) recordProjectTaskOperation(ctx context.Context, tx *gorm.DB, input projectManagementTaskOperation) error {
	metadata, err := json.Marshal(map[string]int64{"project_id": input.projectId, "task_id": input.taskId})
	if err != nil {
		return errorx.WithStack(err)
	}
	return NewOperationLogic(l.BaseLogic).CreateOperationLog(ctx, tx, dto.OperationCreateDto{EntityId: input.projectId, EntityType: constants.OperationEntityTypeProjectManagementProject, OperatorId: input.claim.OrganizationUserId, OperatorName: input.claim.Username, OperatorType: constants.UserTypeBorrower, Operation: input.action, Metadata: metadata, OrganizationId: input.claim.OrganizationId})
}
