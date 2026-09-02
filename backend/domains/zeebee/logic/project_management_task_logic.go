package logic

import (
	"context"
	"database/sql"
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
	claim        utils.JwtUser
	projectId    int64
	taskId       int64
	workflowId   int64
	stageId      int64
	transitionId int64
	fileId       int64
	action       constants.OperationAction
	before       map[string]any
	after        map[string]any
}

func NewProjectManagementTaskLogic(baseLogic *BaseLogic) *ProjectManagementTaskLogic {
	return &ProjectManagementTaskLogic{BaseLogic: baseLogic}
}

func (l *ProjectManagementTaskLogic) Create(ctx context.Context, claim utils.JwtUser, input dto.ProjectManagementTaskCreateDto) (task *model.ProjectManagementTask, err error) {
	if err = validateProjectManagementTaskCreate(input); err != nil {
		return nil, err
	}
	if _, err = l.projectManagementLogic().RequireMember(ctx, input.ProjectId, claim.OrganizationUserId); err != nil {
		return nil, err
	}
	taskData := input.ProjectManagementTaskDto
	taskData.StartAt, taskData.EndAt, taskData.TaskType = nil, nil, nil
	task, err = copierx.Trans[dto.ProjectManagementTaskDto, model.ProjectManagementTask](&taskData)
	if err != nil {
		return nil, err
	}
	task.OrganizationId, task.CreatorId, task.CreatorName = claim.OrganizationId, claim.OrganizationUserId, claim.Username
	task.StartAt, task.EndAt = input.StartAt, input.EndAt
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
		task.StageId, task.Phase = pendingStage.Id, pendingStage.Phase
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
		assignees, err := l.projectManagementTaskAssignees(ctx, tx, projectManagementTaskAssigneeInput{projectId: input.ProjectId, taskId: task.Id, stageId: task.StageId, organizationId: claim.OrganizationId, organizationUserIds: input.AssigneeOrganizationUserIds})
		if err != nil {
			return err
		}
		if err = l.repos.ProjectManagementTaskRepo.CreateAssigneesWithTx(ctx, tx, assignees); err != nil {
			return err
		}
		if err = l.recordProjectTaskOperation(ctx, tx, projectManagementTaskOperation{claim: claim, projectId: task.ProjectId, taskId: task.Id, action: constants.OperationActionProjectManagementTaskCreate, after: map[string]any{"id": task.Id, "name": task.Title}}); err != nil {
			return err
		}
		for _, assignee := range assignees {
			if err = l.recordProjectTaskOperation(ctx, tx, projectManagementTaskOperation{claim: claim, projectId: task.ProjectId, taskId: task.Id, action: constants.OperationActionProjectManagementTaskAssigneeAdd, after: assigneeOperationValues(assignee)}); err != nil {
				return err
			}
		}
		return nil
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
	return validateProjectManagementTaskTimeRange(input.StartAt, input.EndAt)
}

func validateProjectManagementTaskTimeRange(startAt, endAt *sql.NullTime) error {
	if startAt != nil && startAt.Valid && endAt != nil && endAt.Valid && endAt.Time.Before(startAt.Time) {
		return errorx.New("task end time must not be before start time")
	}
	return nil
}

func (l *ProjectManagementTaskLogic) Detail(ctx context.Context, claim utils.JwtUser, taskId int64) (*model.ProjectManagementTask, error) {
	task, err := l.repos.ProjectManagementTaskRepo.FindDetailByOrganization(ctx, taskId, claim.OrganizationId)
	if err != nil {
		return nil, projectManagementNotFound(err, "task not found")
	}
	if _, err = l.projectManagementLogic().RequireMember(ctx, task.ProjectId, claim.OrganizationUserId); err != nil {
		return nil, err
	}
	return task, nil
}

func (l *ProjectManagementTaskLogic) Page(ctx context.Context, claim utils.JwtUser, input dto.ProjectManagementTaskPageRequestDto) (gormx.PageResponse[*model.ProjectManagementTask], error) {
	if _, err := l.projectManagementLogic().RequireProject(ctx, input.Filter.ProjectId, claim.OrganizationId); err != nil {
		return gormx.PageResponse[*model.ProjectManagementTask]{}, err
	}
	if _, err := l.projectManagementLogic().RequireMember(ctx, input.Filter.ProjectId, claim.OrganizationUserId); err != nil {
		return gormx.PageResponse[*model.ProjectManagementTask]{}, err
	}
	page, err := l.repos.ProjectManagementTaskRepo.PageByProject(ctx, claim.OrganizationId, input)
	if err != nil {
		return page, err
	}
	return page, nil
}

func (l *ProjectManagementTaskLogic) Candidate(ctx context.Context, claim utils.JwtUser, input dto.ProjectManagementTaskCandidateDto) ([]*dto.ProjectManagementTaskCandidateEntityDto, error) {
	if _, err := l.projectManagementLogic().RequireProject(ctx, input.ProjectId, claim.OrganizationId); err != nil {
		return nil, err
	}
	if _, err := l.projectManagementLogic().RequireMember(ctx, input.ProjectId, claim.OrganizationUserId); err != nil {
		return nil, err
	}
	return l.repos.ProjectManagementTaskRepo.CandidateByProject(ctx, claim.OrganizationId, input)
}

func (l *ProjectManagementTaskLogic) Update(ctx context.Context, claim utils.JwtUser, taskId int64, input dto.ProjectManagementTaskUpdateDto) (*model.ProjectManagementTask, error) {
	if strings.TrimSpace(input.Title) == "" {
		return nil, errorx.New("task title is required")
	}
	if err := validateProjectManagementTaskTimeRange(input.StartAt, input.EndAt); err != nil {
		return nil, err
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
		before := taskOperationValues(task)
		base, taskType := task.Base, task.TaskType
		if err = copierx.Assign(&input, task); err != nil {
			return err
		}
		task.Base, task.TaskType = base, taskType
		if input.TaskTypeId != 0 {
			taskType, err := l.repos.ProjectManagementTaskTypeRepo.FindByOrganizationWithTx(ctx, tx, input.TaskTypeId, claim.OrganizationId)
			if err != nil {
				return projectManagementNotFound(err, "task type not found")
			}
			task.TaskType = ProjectManagementTaskTypeSnapshot(taskType)
		}
		assignees, err := l.projectManagementTaskAssignees(ctx, tx, projectManagementTaskAssigneeInput{projectId: task.ProjectId, taskId: task.Id, stageId: task.StageId, organizationId: claim.OrganizationId, organizationUserIds: input.AssigneeOrganizationUserIds})
		if err != nil {
			return err
		}
		existingAssignees, err := l.repos.ProjectManagementTaskRepo.ListAssigneesWithTx(ctx, tx, task.Id, task.StageId)
		if err != nil {
			return err
		}
		if err = l.repos.ProjectManagementTaskRepo.ReplaceAssigneesWithTx(ctx, tx, task.Id, task.StageId, assignees); err != nil {
			return err
		}
		requestedAssignees := make(map[int64]struct{}, len(assignees))
		for _, assignee := range assignees {
			requestedAssignees[assignee.OrganizationUserId] = struct{}{}
		}
		for _, assignee := range existingAssignees {
			if _, exists := requestedAssignees[assignee.OrganizationUserId]; exists {
				continue
			}
			if err = l.recordProjectTaskOperation(ctx, tx, projectManagementTaskOperation{claim: claim, projectId: task.ProjectId, taskId: task.Id, action: constants.OperationActionProjectManagementTaskAssigneeRemove, before: assigneeOperationValues(assignee)}); err != nil {
				return err
			}
		}
		existingAssigneeIds := make(map[int64]struct{}, len(existingAssignees))
		for _, assignee := range existingAssignees {
			existingAssigneeIds[assignee.OrganizationUserId] = struct{}{}
		}
		for _, assignee := range assignees {
			if _, exists := existingAssigneeIds[assignee.OrganizationUserId]; exists {
				continue
			}
			if err = l.recordProjectTaskOperation(ctx, tx, projectManagementTaskOperation{claim: claim, projectId: task.ProjectId, taskId: task.Id, action: constants.OperationActionProjectManagementTaskAssigneeAdd, after: assigneeOperationValues(assignee)}); err != nil {
				return err
			}
		}
		if err = l.syncAttachmentsWithTx(ctx, tx, claim, task, input.Attachments); err != nil {
			return err
		}
		if err = NewProjectManagementDependencyLogic(l.BaseLogic).SyncWithTx(ctx, tx, claim, task.Id, input.Dependencies); err != nil {
			return err
		}
		task.RemoveRelationships()
		if _, err = l.repos.ProjectManagementTaskRepo.UpdateOneWithTx(ctx, tx, task); err != nil {
			return err
		}
		after := taskOperationValues(task)
		before, after = operationValuesDiff(before, after)
		if len(before) == 0 {
			return nil
		}
		before["id"], after["id"] = task.Id, task.Id
		return l.recordProjectTaskOperation(ctx, tx, projectManagementTaskOperation{claim: claim, projectId: task.ProjectId, taskId: task.Id, action: constants.OperationActionProjectManagementTaskUpdate, before: before, after: after})
	})
	if err != nil {
		return nil, err
	}
	return l.Detail(ctx, claim, task.Id)
}

func (l *ProjectManagementTaskLogic) syncAttachmentsWithTx(ctx context.Context, tx *gorm.DB, claim utils.JwtUser, task *model.ProjectManagementTask, inputs []*dto.ProjectManagementTaskAttachmentDto) error {
	existing, err := l.repos.FileRepo.ListByEntityWithTx(ctx, tx, task.Id, constants.FileEntityTypeProjectManagementTask)
	if err != nil {
		return err
	}
	existingByID := make(map[int64]*model.File, len(existing))
	for _, file := range existing {
		existingByID[file.Id] = file
	}
	requested := make(map[int64]struct{}, len(inputs))
	for _, input := range inputs {
		if input == nil {
			return errorx.New("task attachment is required")
		}
		if input.Id != 0 {
			if _, ok := existingByID[input.Id]; !ok {
				return errorx.New("task attachment not found")
			}
			if _, ok := requested[input.Id]; ok {
				return errorx.New("duplicate task attachment")
			}
			requested[input.Id] = struct{}{}
			continue
		}
		if strings.TrimSpace(input.FileUrl) == "" {
			return errorx.New("attachment file url is required")
		}
		storageType := input.StorageType
		if storageType == "" {
			storageType = constants.StorageTypeAwsS3
		}
		if !storageType.IsValid() {
			return errorx.New("invalid attachment storage type")
		}
		uploadType := input.UploadType
		if uploadType == "" {
			uploadType = constants.UploadTypeManual
		}
		if !uploadType.IsValid() {
			return errorx.New("invalid attachment upload type")
		}
		file, err := l.repos.FileRepo.CreateWithTx(ctx, tx, &model.File{EntityId: task.Id, EntityType: constants.FileEntityTypeProjectManagementTask, FileName: input.FileName, FileUrl: input.FileUrl, StorageType: storageType, UploadType: uploadType}, false)
		if err != nil {
			return err
		}
		if err = l.recordProjectTaskOperation(ctx, tx, projectManagementTaskOperation{claim: claim, projectId: task.ProjectId, taskId: task.Id, fileId: file.Id, action: constants.OperationActionProjectManagementAttachmentAdd, after: attachmentOperationIdentity(file)}); err != nil {
			return err
		}
	}
	for _, file := range existing {
		if _, ok := requested[file.Id]; ok {
			continue
		}
		if _, err = l.repos.FileRepo.DeleteByIDWithTx(ctx, tx, file.Id); err != nil {
			return err
		}
		if err = l.recordProjectTaskOperation(ctx, tx, projectManagementTaskOperation{claim: claim, projectId: task.ProjectId, taskId: task.Id, fileId: file.Id, action: constants.OperationActionProjectManagementAttachmentDelete, before: attachmentOperationIdentity(file)}); err != nil {
			return err
		}
	}
	return nil
}

func (l *ProjectManagementTaskLogic) Delete(ctx context.Context, claim utils.JwtUser, taskId int64) error {
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
		return l.recordProjectTaskOperation(ctx, tx, projectManagementTaskOperation{claim: claim, projectId: task.ProjectId, taskId: task.Id, action: constants.OperationActionProjectManagementTaskDelete, before: map[string]any{"id": task.Id, "title": task.Title}})
	})
}

func (l *ProjectManagementTaskLogic) Overdue(ctx context.Context, task *model.ProjectManagementTask) (*dto.ProjectManagementTaskOverdueDto, error) {
	stage, err := l.repos.ProjectManagementWorkflowRepo.FindStage(ctx, task.StageId)
	if err != nil {
		return nil, projectManagementNotFound(err, "task stage not found")
	}
	overdue, overdueSeconds := projectManagementTaskOverdue(task, stage.Phase, time.Now())
	return &dto.ProjectManagementTaskOverdueDto{Overdue: overdue, OverdueSeconds: overdueSeconds}, nil
}

func (l *ProjectManagementTaskLogic) MyPage(ctx context.Context, claim utils.JwtUser, input dto.ProjectManagementMyTaskPageRequestDto) (gormx.PageResponse[*model.ProjectManagementTask], error) {
	page, err := l.repos.ProjectManagementTaskRepo.PageMine(ctx, claim.OrganizationId, claim.OrganizationUserId, input)
	if err != nil {
		return page, err
	}
	return page, nil
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
		result = append(result, &model.ProjectManagementTaskAssignee{TaskId: input.taskId, StageId: input.stageId, OrganizationUserId: organizationUserId, OrganizationUserFullName: organizationUser.FullName})
	}
	return result, nil
}

func (l *ProjectManagementTaskLogic) recordProjectTaskOperation(ctx context.Context, tx *gorm.DB, input projectManagementTaskOperation) error {
	metadata, err := marshalProjectManagementOperationMetadata(projectManagementOperationMetadata{
		ProjectID:    input.projectId,
		TaskID:       input.taskId,
		WorkflowID:   input.workflowId,
		StageID:      input.stageId,
		TransitionID: input.transitionId,
		FileID:       input.fileId,
	}, input.before, input.after)
	if err != nil {
		return errorx.WithStack(err)
	}
	return NewOperationLogic(l.BaseLogic).CreateOperationLog(ctx, tx, dto.OperationCreateDto{EntityId: input.projectId, EntityType: constants.OperationEntityTypeProjectManagementProject, OperatorId: input.claim.OrganizationUserId, OperatorName: input.claim.Username, OperatorType: constants.UserTypeBorrower, Operation: input.action, Metadata: metadata, OrganizationId: input.claim.OrganizationId})
}

func taskOperationValues(task *model.ProjectManagementTask) map[string]any {
	values := map[string]any{"title": task.Title, "description": task.Description, "priority": task.Priority, "work_hours": task.WorkHours, "parent_task_id": task.ParentTaskId, "task_type_id": int64(0), "task_type_name": "", "start_at": "", "end_at": ""}
	if task.TaskType != nil {
		values["task_type_id"], values["task_type_name"] = task.TaskType.Id, task.TaskType.Name
	}
	if task.StartAt != nil && task.StartAt.Valid {
		values["start_at"] = task.StartAt.Time.UTC().Format(time.RFC3339)
	}
	if task.EndAt != nil && task.EndAt.Valid {
		values["end_at"] = task.EndAt.Time.UTC().Format(time.RFC3339)
	}
	return values
}

func assigneeOperationValues(assignee *model.ProjectManagementTaskAssignee) map[string]any {
	return map[string]any{"organization_user_id": assignee.OrganizationUserId, "organization_user_full_name": assignee.OrganizationUserFullName}
}
