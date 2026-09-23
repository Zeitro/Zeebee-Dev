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
	"zeebee-backend/domains/zeebee/repo"
)

type ProjectManagementProjectLogic struct{ *BaseLogic }

// CalculateProjectEndAt returns the inclusive end date for a planned project.
func CalculateProjectEndAt(startAt time.Time, periodValue int, periodUnit constants.ProjectManagementPeriodUnit) (time.Time, error) {
	if periodValue <= 0 {
		return time.Time{}, errorx.New("project period value must be positive")
	}
	switch periodUnit {
	case constants.ProjectManagementPeriodUnitDay:
		return startAt.AddDate(0, 0, periodValue-1), nil
	case constants.ProjectManagementPeriodUnitWeek:
		return startAt.AddDate(0, 0, periodValue*7-1), nil
	case constants.ProjectManagementPeriodUnitMonth:
		target := time.Date(startAt.Year(), startAt.Month()+time.Month(periodValue-1), 1, startAt.Hour(), startAt.Minute(), startAt.Second(), startAt.Nanosecond(), startAt.Location())
		lastDay := time.Date(target.Year(), target.Month()+1, 0, startAt.Hour(), startAt.Minute(), startAt.Second(), startAt.Nanosecond(), startAt.Location()).Day()
		return time.Date(target.Year(), target.Month(), min(startAt.Day(), lastDay), startAt.Hour(), startAt.Minute(), startAt.Second(), startAt.Nanosecond(), startAt.Location()), nil
	default:
		return time.Time{}, errorx.New("invalid project period unit")
	}
}

type projectManagementProjectOperation struct {
	claim            utils.JwtUser
	projectId        int64
	action           constants.OperationAction
	templateId       int64
	templateType     string
	copyTasks        *bool
	skippedAssignees []map[string]any
	before           map[string]any
	after            map[string]any
}

func NewProjectManagementProjectLogic(baseLogic *BaseLogic) *ProjectManagementProjectLogic {
	return &ProjectManagementProjectLogic{BaseLogic: baseLogic}
}

func (l *ProjectManagementProjectLogic) Page(ctx context.Context, claim utils.JwtUser, input dto.ProjectManagementProjectPageRequestDto) (gormx.PageResponse[*model.ProjectManagementProject], error) {
	return l.repos.ProjectManagementProjectRepo.Page(ctx, claim.OrganizationId, input)
}

func (l *ProjectManagementMemberLogic) Page(ctx context.Context, claim utils.JwtUser, input dto.ProjectManagementMemberPageRequestDto) (gormx.PageResponse[*dto.ProjectManagementMemberDetailDto], error) {
	if _, err := l.projectManagementLogic().RequireProject(ctx, input.Filter.ProjectId, claim.OrganizationId); err != nil {
		return gormx.PageResponse[*dto.ProjectManagementMemberDetailDto]{}, err
	}
	if _, err := l.projectManagementLogic().RequireMember(ctx, input.Filter.ProjectId, claim.OrganizationUserId); err != nil {
		return gormx.PageResponse[*dto.ProjectManagementMemberDetailDto]{}, err
	}
	page, err := l.repos.ProjectManagementProjectRepo.PageMembers(ctx, input.Filter.ProjectId, input)
	if err != nil {
		return gormx.PageResponse[*dto.ProjectManagementMemberDetailDto]{}, err
	}
	list, err := projectManagementMemberDtos(page.List)
	if err != nil {
		return gormx.PageResponse[*dto.ProjectManagementMemberDetailDto]{}, err
	}
	return gormx.PageResponse[*dto.ProjectManagementMemberDetailDto]{Page: page.Page, Limit: page.Limit, Total: page.Total, TotalPage: page.TotalPage, List: list}, nil
}

func (l *ProjectManagementMemberLogic) Candidate(ctx context.Context, claim utils.JwtUser, input dto.ProjectManagementMemberCandidateDto) ([]*dto.ProjectManagementMemberCandidateEntityDto, error) {
	if input.Purpose != constants.ProjectManagementMemberCandidatePurposeAddMember && input.Purpose != constants.ProjectManagementMemberCandidatePurposeTransferOwner && input.Purpose != constants.ProjectManagementMemberCandidatePurposeProjectMember {
		return nil, errorx.New("invalid project member candidate purpose")
	}
	if _, err := l.projectManagementLogic().RequireProject(ctx, input.ProjectId, claim.OrganizationId); err != nil {
		return nil, err
	}
	candidates, err := l.repos.ProjectManagementProjectRepo.CandidateMembers(ctx, claim.OrganizationId, input.ProjectId, input.Keyword, input.Purpose)
	if err != nil {
		return nil, err
	}
	return candidates, nil
}

func (l *ProjectManagementProjectLogic) RequirePropertyDeletable(ctx context.Context, propertyId int64) error {
	blocked, err := l.repos.ProjectManagementProjectRepo.HasUnfinishedTasksForEntity(ctx, propertyId, constants.ProjectManagementEntityTypeBorrowerProperty)
	if err != nil {
		return err
	}
	if blocked {
		return errorx.New("cannot delete a property with unfinished project tasks")
	}
	return nil
}

func projectManagementMemberDtos(models []*repo.ProjectManagementProjectMemberEntity) ([]*dto.ProjectManagementMemberDetailDto, error) {
	result := make([]*dto.ProjectManagementMemberDetailDto, 0, len(models))
	for _, item := range models {
		member, err := copierx.Trans[repo.ProjectManagementProjectMemberEntity, dto.ProjectManagementMemberDetailDto](item)
		if err != nil {
			return nil, err
		}
		result = append(result, member)
	}
	return result, nil
}

func (l *ProjectManagementProjectLogic) Create(ctx context.Context, claim utils.JwtUser, input dto.ProjectManagementProjectCreateDto) (*model.ProjectManagementProject, error) {
	project, _, err := l.CreateWithAssigneeResult(ctx, claim, input)
	return project, err
}

func (l *ProjectManagementProjectLogic) CreateWithAssigneeResult(ctx context.Context, claim utils.JwtUser, input dto.ProjectManagementProjectCreateDto) (project *model.ProjectManagementProject, assigneeResult *dto.ProjectManagementProjectCreateAssigneeResultDto, err error) {
	if err = validateProjectManagementProjectCreate(input); err != nil {
		return nil, nil, err
	}
	// Validate entity
	if input.EntityType == constants.ProjectManagementEntityTypeBorrowerProperty && input.EntityId > 0 {
		exists, err := l.repos.BorrowerPropertyRepo.IsExistedByColumns(ctx, map[string]any{"id": input.EntityId, "organization_id": claim.OrganizationId})
		if err != nil {
			return nil, nil, err
		}
		if !exists {
			return nil, nil, errorx.New("property not found")
		}
	}
	project, err = copierx.Trans[dto.ProjectManagementProjectCreateDto, model.ProjectManagementProject](&input)
	if err != nil {
		return nil, nil, err
	}
	project.OrganizationId, project.OwnerId, project.OwnerName = claim.OrganizationId, claim.OrganizationUserId, claim.Username
	if err = prepareProjectSchedule(project); err != nil {
		return nil, nil, err
	}
	if input.TemplateId == 0 {
		if err = NewProjectManagementWorkflowLogic(l.BaseLogic).ensureDefault(ctx, claim.OrganizationId); err != nil {
			return nil, nil, err
		}
	}
	assigneeResult = &dto.ProjectManagementProjectCreateAssigneeResultDto{CopiedOrganizationUserIds: make([]int64, 0), SkippedAssignees: make([]*dto.ProjectManagementTemplateAssigneeResultDto, 0)}
	err = l.repos.ProjectManagementProjectRepo.Transaction(ctx, func(tx *gorm.DB) error {
		var innerErr error
		project, innerErr = l.repos.ProjectManagementProjectRepo.CreateWithTx(ctx, tx, project, false)
		if innerErr != nil {
			return innerErr
		}
		if innerErr = l.repos.ProjectManagementProjectRepo.CreateMemberWithTx(ctx, tx, &model.ProjectManagementProjectMember{
			ProjectId:          project.Id,
			OrganizationUserId: claim.OrganizationUserId,
			Role:               constants.ProjectManagementMemberRoleProjectOwner,
		}); innerErr != nil {
			return innerErr
		}
		if innerErr = l.repos.ProjectManagementWorkflowRepo.CloneTemplateToProjectWithTx(ctx, tx, input.TemplateId, claim.OrganizationId, project.Id); innerErr != nil {
			return projectManagementNotFound(innerErr, "project template not found")
		}
		if innerErr = NewProjectManagementWorkflowLogic(l.BaseLogic).validateProjectWorkflowInvariantWithTx(ctx, tx, project.Id); innerErr != nil {
			return innerErr
		}
		templateType := ""
		if input.TemplateId != 0 {
			template, templateErr := l.repos.ProjectManagementTemplateRepo.FindVisibleProjectTemplateWithTx(ctx, tx, input.TemplateId, claim.OrganizationId)
			if templateErr != nil {
				return projectManagementNotFound(templateErr, "project template not found")
			}
			if innerErr = l.instantiateTemplateTasksWithTx(ctx, tx, claim, project, template, input.CopyAssignees, assigneeResult); innerErr != nil {
				return innerErr
			}
			templateType = string(template.TemplateType)
		}
		skippedAssignees := make([]map[string]any, 0, len(assigneeResult.SkippedAssignees))
		for _, skipped := range assigneeResult.SkippedAssignees {
			skippedAssignees = append(skippedAssignees, map[string]any{"organization_user_id": skipped.OrganizationUserId, "reason": skipped.Reason})
		}
		return l.recordOperation(ctx, tx, projectManagementProjectOperation{
			claim: claim, projectId: project.Id, action: constants.OperationActionProjectManagementProjectCreate,
			templateId: input.TemplateId, templateType: templateType,
			skippedAssignees: skippedAssignees,
			after:            map[string]any{"id": project.Id, "name": project.Name},
		})
	})
	return project, assigneeResult, err
}

func (l *ProjectManagementProjectLogic) instantiateTemplateTasksWithTx(ctx context.Context, tx *gorm.DB, claim utils.JwtUser, project *model.ProjectManagementProject, template *model.ProjectManagementTemplate, copyAssignees bool, assigneeResult *dto.ProjectManagementProjectCreateAssigneeResultDto) error {
	tasks, err := l.repos.ProjectManagementTaskRepo.ListByOwnerWithTx(ctx, tx, constants.ProjectManagementOwnerTypeTemplate, template.Id)
	if err != nil || len(tasks) == 0 {
		return err
	}
	workflows, err := l.repos.ProjectManagementWorkflowRepo.ListProjectWorkflowsWithTx(ctx, tx, project.Id)
	if err != nil {
		return err
	}
	workflowIDs := make(map[int64]int64, len(workflows))
	projectWorkflowIDs := make([]int64, 0, len(workflows))
	for _, workflow := range workflows {
		workflowIDs[workflow.SourceWorkflowId] = workflow.Id
		projectWorkflowIDs = append(projectWorkflowIDs, workflow.Id)
	}
	pendingStages, err := l.repos.ProjectManagementWorkflowRepo.ListPendingOwnerStagesWithTx(ctx, tx, projectWorkflowIDs)
	if err != nil {
		return err
	}
	taskIDs := make([]int64, 0, len(tasks))
	newTaskIDs := make(map[int64]int64, len(tasks))
	pendingStageIDs := make(map[int64]int64, len(tasks))
	for _, task := range tasks {
		workflowId := workflowIDs[task.WorkflowId]
		if workflowId == 0 {
			return errorx.New("project workflow mapping is invalid")
		}
		pending := pendingStages[workflowId]
		if pending == nil {
			return errorx.New("project pending stage not found")
		}
		startAt, endAt := shiftTemplateTaskTimes(template, project, task.StartAt, task.EndAt)
		if err = validateProjectManagementTaskTimeRange(startAt, endAt); err != nil {
			return err
		}
		if err = validateTaskWithinProjectSchedule(project, startAt, endAt); err != nil {
			return err
		}
		phase := pending.Phase
		copy := &model.ProjectManagementTask{EntityType: constants.ProjectManagementOwnerTypeProject, EntityId: project.Id, OrganizationId: project.OrganizationId, WorkflowId: workflowId, StageId: pending.Id, Phase: &phase, CreatorId: claim.OrganizationUserId, CreatorName: claim.Username, Title: task.Title, Description: task.Description, Priority: task.Priority, WorkHours: task.WorkHours, StartAt: startAt, EndAt: endAt, TaskType: task.TaskType}
		if err = l.repos.ProjectManagementTaskRepo.CreateEntityWithTx(ctx, tx, copy); err != nil {
			return err
		}
		taskIDs = append(taskIDs, task.Id)
		newTaskIDs[task.Id] = copy.Id
		pendingStageIDs[task.Id] = pending.Id
	}
	for _, task := range tasks {
		if task.ParentTaskId != 0 {
			if err = l.repos.ProjectManagementTaskRepo.UpdateParentWithTx(ctx, tx, newTaskIDs[task.Id], newTaskIDs[task.ParentTaskId]); err != nil {
				return err
			}
		}
	}
	dependencies, err := l.repos.ProjectManagementDependencyRepo.ListByTaskIdsWithTx(ctx, tx, taskIDs)
	if err != nil {
		return err
	}
	for _, dependency := range dependencies {
		if newTaskIDs[dependency.TaskId] == 0 || newTaskIDs[dependency.DependencyTaskId] == 0 {
			return errorx.New("project dependency mapping is invalid")
		}
		if _, err = l.repos.ProjectManagementDependencyRepo.CreateWithTx(ctx, tx, &model.ProjectManagementTaskDependency{TaskId: newTaskIDs[dependency.TaskId], DependencyTaskId: newTaskIDs[dependency.DependencyTaskId], DependencyType: dependency.DependencyType, DependencyRule: dependency.DependencyRule}); err != nil {
			return err
		}
	}
	if !copyAssignees {
		return nil
	}
	assignees, err := l.repos.ProjectManagementTaskRepo.ListTemplateAssigneesByTaskIdsWithTx(ctx, tx, taskIDs)
	if err != nil {
		return err
	}
	seenAssignments := make(map[[2]int64]struct{}, len(assignees))
	seenMembers := make(map[int64]struct{}, len(assignees))
	copies := make([]*model.ProjectManagementTaskAssignee, 0, len(assignees))
	for _, assignee := range assignees {
		active, err := l.repos.OrganizationUserRepo.IsActiveByIdWithTx(ctx, tx, assignee.OrganizationUserId, project.OrganizationId)
		if err != nil {
			return err
		}
		if !active {
			assigneeResult.SkippedAssignees = append(assigneeResult.SkippedAssignees, &dto.ProjectManagementTemplateAssigneeResultDto{OrganizationUserId: assignee.OrganizationUserId, Reason: "organization_user_not_active"})
			continue
		}
		if _, ok := seenMembers[assignee.OrganizationUserId]; !ok {
			if err = l.repos.ProjectManagementProjectRepo.EnsureExecutorMemberWithTx(ctx, tx, project.Id, assignee.OrganizationUserId); err != nil {
				return err
			}
			seenMembers[assignee.OrganizationUserId] = struct{}{}
			assigneeResult.CopiedOrganizationUserIds = append(assigneeResult.CopiedOrganizationUserIds, assignee.OrganizationUserId)
		}
		key := [2]int64{newTaskIDs[assignee.TaskId], assignee.OrganizationUserId}
		if _, ok := seenAssignments[key]; ok {
			continue
		}
		seenAssignments[key] = struct{}{}
		stageID := pendingStageIDs[assignee.TaskId]
		copies = append(copies, &model.ProjectManagementTaskAssignee{TaskId: key[0], StageId: &stageID, OrganizationUserId: assignee.OrganizationUserId, OrganizationUserFullName: assignee.OrganizationUserFullName})
	}
	return l.repos.ProjectManagementTaskRepo.CreateAssigneesWithTx(ctx, tx, copies)
}

func shiftTemplateTaskTimes(template *model.ProjectManagementTemplate, project *model.ProjectManagementProject, startAt, endAt *sql.NullTime) (*sql.NullTime, *sql.NullTime) {
	if template.StartAt == nil || !template.StartAt.Valid || project.StartAt == nil || !project.StartAt.Valid {
		return startAt, endAt
	}
	templateDate := time.Date(template.StartAt.Time.Year(), template.StartAt.Time.Month(), template.StartAt.Time.Day(), 0, 0, 0, 0, time.UTC)
	projectDate := time.Date(project.StartAt.Time.Year(), project.StartAt.Time.Month(), project.StartAt.Time.Day(), 0, 0, 0, 0, time.UTC)
	days := int(projectDate.Sub(templateDate) / (24 * time.Hour))
	shift := func(value *sql.NullTime) *sql.NullTime {
		if value == nil || !value.Valid {
			return value
		}
		return &sql.NullTime{Time: value.Time.AddDate(0, 0, days), Valid: true}
	}
	return shift(startAt), shift(endAt)
}

func (l *ProjectManagementProjectLogic) Update(ctx context.Context, claim utils.JwtUser, projectId int64, input dto.ProjectManagementProjectUpdateDto) (*model.ProjectManagementProject, error) {
	if err := validateProjectManagementProjectUpdate(input); err != nil {
		return nil, err
	}
	project, err := l.projectManagementLogic().RequireProject(ctx, projectId, claim.OrganizationId)
	if err != nil {
		return nil, err
	}
	if _, err = l.projectManagementLogic().RequireMemberRole(ctx, projectId, claim.OrganizationUserId, constants.ProjectManagementMemberRoleProjectOwner, constants.ProjectManagementMemberRoleProjectAdmin); err != nil {
		return nil, err
	}
	before := projectOperationValues(project)
	startAt, periodValue, periodUnit, endAt := project.StartAt, project.PeriodValue, project.PeriodUnit, project.EndAt
	if err = copierx.Assign(&input, project); err != nil {
		return nil, err
	}
	if input.StartAt == nil {
		project.StartAt = startAt
	}
	if input.PeriodValue == nil {
		project.PeriodValue = periodValue
	}
	if input.PeriodUnit == nil {
		project.PeriodUnit = periodUnit
	}
	if project.ProjectType == constants.ProjectManagementProjectTypePlanned {
		if err = prepareProjectSchedule(project); err != nil {
			return nil, err
		}
	} else {
		project.EndAt = endAt
	}
	project.RemoveRelationships()

	err = l.repos.ProjectManagementProjectRepo.Transaction(ctx, func(tx *gorm.DB) error {
		if _, innerErr := l.repos.ProjectManagementProjectRepo.UpdateOneByIdWithTx(ctx, tx, project.Id, project); innerErr != nil {
			return innerErr
		}
		after := projectOperationValues(project)
		before, after = operationValuesDiff(before, after)
		if len(before) == 0 {
			return nil
		}
		return l.recordOperation(ctx, tx, projectManagementProjectOperation{
			claim:     claim,
			projectId: projectId,
			action:    constants.OperationActionProjectManagementProjectUpdate,
			before:    before,
			after:     after,
		})
	})
	return project, err
}

func (l *ProjectManagementProjectLogic) Delete(ctx context.Context, claim utils.JwtUser, projectId int64) error {
	return l.repos.ProjectManagementProjectRepo.Transaction(ctx, func(tx *gorm.DB) error {
		project, err := l.projectManagementLogic().requireProjectManagerWithTx(ctx, tx, claim, projectId)
		if err != nil {
			return err
		}
		taskIds, err := l.repos.ProjectManagementTaskRepo.IdsByProjectWithTx(ctx, tx, projectId)
		if err != nil {
			return err
		}
		if err = l.repos.FileRepo.DeleteByEntityIdWithTx(ctx, tx, taskIds, constants.FileEntityTypeProjectManagementTask); err != nil {
			return err
		}
		if err = l.repos.ProjectManagementTaskRepo.DeleteByProjectWithTx(ctx, tx, projectId, taskIds); err != nil {
			return err
		}
		if err = l.repos.ProjectManagementWorkflowRepo.DeleteProjectWorkflowWithTx(ctx, tx, projectId); err != nil {
			return err
		}
		if err = l.repos.ProjectManagementProjectRepo.DeleteProjectWithTx(ctx, tx, projectId); err != nil {
			return err
		}
		return l.recordOperation(ctx, tx, projectManagementProjectOperation{claim: claim, projectId: projectId, action: constants.OperationActionProjectManagementProjectDelete, before: projectOperationValues(project)})
	})
}

func validateProjectManagementProjectCreate(input dto.ProjectManagementProjectCreateDto) error {
	if strings.TrimSpace(input.Name) == "" {
		return errorx.New("project name is required")
	}
	return input.Priority.Validate()
}

func validateProjectManagementProjectUpdate(input dto.ProjectManagementProjectUpdateDto) error {
	if strings.TrimSpace(input.Name) == "" {
		return errorx.New("project name is required")
	}
	return input.Priority.Validate()
}

func prepareProjectSchedule(project *model.ProjectManagementProject) error {
	if project.ProjectType == "" {
		project.ProjectType = constants.ProjectManagementProjectTypeOrdinary
	}
	switch project.ProjectType {
	case constants.ProjectManagementProjectTypeOrdinary:
		if project.StartAt != nil || project.EndAt != nil || project.PeriodValue != 0 || project.PeriodUnit != "" {
			return errorx.New("ordinary project cannot have schedule fields")
		}
	case constants.ProjectManagementProjectTypePlanned:
		if project.StartAt == nil || !project.StartAt.Valid {
			return errorx.New("planned project start date is required")
		}
		endAt, err := CalculateProjectEndAt(project.StartAt.Time, project.PeriodValue, project.PeriodUnit)
		if err != nil {
			return err
		}
		project.EndAt = &sql.NullTime{Time: endAt, Valid: true}
	default:
		return errorx.New("invalid project type")
	}
	return nil
}

func (l *ProjectManagementProjectLogic) recordOperation(ctx context.Context, tx *gorm.DB, input projectManagementProjectOperation) error {
	metadata, err := marshalProjectManagementOperationMetadata(projectManagementOperationMetadata{
		ProjectID: input.projectId, TemplateID: input.templateId, TemplateType: input.templateType,
		CopyTasks: input.copyTasks, SkippedAssignees: input.skippedAssignees,
	}, input.before, input.after)
	if err != nil {
		return errorx.WithStack(err)
	}
	return NewOperationLogic(l.BaseLogic).CreateOperationLog(ctx, tx, dto.OperationCreateDto{EntityId: input.projectId, EntityType: constants.OperationEntityTypeProjectManagementProject, OperatorId: input.claim.OrganizationUserId, OperatorName: input.claim.Username, OperatorType: constants.UserTypeBorrower, Operation: input.action, Metadata: metadata, OrganizationId: input.claim.OrganizationId})
}

func projectOperationValues(project *model.ProjectManagementProject) map[string]any {
	return map[string]any{
		"name":        project.Name,
		"description": project.Description,
		"priority":    project.Priority,
		"cover_url":   project.CoverUrl,
		"entity_id":   project.EntityId,
		"entity_type": project.EntityType,
		"owner_id":    project.OwnerId,
		"owner_name":  project.OwnerName,
	}
}
