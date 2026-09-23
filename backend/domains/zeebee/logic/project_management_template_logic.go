package logic

import (
	"context"
	"database/sql"
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

// ProjectManagementTemplateLogic owns public project-template metadata. Workflow
// definitions remain in ProjectManagementWorkflowLogic and are returned by ID.
type ProjectManagementTemplateLogic struct{ *BaseLogic }

func NewProjectManagementTemplateLogic(baseLogic *BaseLogic) *ProjectManagementTemplateLogic {
	return &ProjectManagementTemplateLogic{BaseLogic: baseLogic}
}

func (l *ProjectManagementTemplateLogic) Page(ctx context.Context, input dto.ProjectManagementTemplatePageRequestDto) (gormx.PageResponse[*model.ProjectManagementTemplate], error) {
	return l.repos.ProjectManagementTemplateRepo.Page(ctx, input)
}

func (l *ProjectManagementTemplateLogic) Create(ctx context.Context, claim utils.JwtUser, input dto.ProjectManagementTemplateCreateDto) (*model.ProjectManagementTemplate, error) {
	if input.CopyAssignees && !input.CopyTasks {
		return nil, errorx.New("copy assignees requires copy tasks")
	}
	if input.TemplateType == "" {
		input.TemplateType = constants.ProjectManagementTemplateTypePublic
	}
	if input.TemplateType != constants.ProjectManagementTemplateTypePublic {
		return nil, errorx.New("only project templates can be created here")
	}
	input.Status = constants.CommonStatusActive
	template, err := copierx.Trans[dto.ProjectManagementTemplateCreateDto, model.ProjectManagementTemplate](&input)
	if err != nil {
		return nil, err
	}
	template.OrganizationId = 0 // Admin-created templates are public.
	if input.SourceProjectId != 0 {
		template.TemplateType = constants.ProjectManagementTemplateTypeOrganization
		if _, err = l.projectManagementLogic().RequireProject(ctx, input.SourceProjectId, claim.OrganizationId); err != nil {
			return nil, err
		}
		if _, err = l.projectManagementLogic().RequireMemberRole(ctx, input.SourceProjectId, claim.OrganizationUserId, constants.ProjectManagementMemberRoleProjectOwner, constants.ProjectManagementMemberRoleProjectAdmin); err != nil {
			return nil, err
		}
		template.OrganizationId = claim.OrganizationId
	} else if err = prepareProjectManagementTemplateSchedule(template); err != nil {
		return nil, err
	}
	err = l.repos.ProjectManagementTemplateRepo.Transaction(ctx, func(tx *gorm.DB) error {
		if input.SourceProjectId != 0 {
			project, err := l.repos.ProjectManagementProjectRepo.LockByOrganizationWithTx(ctx, tx, input.SourceProjectId, claim.OrganizationId)
			if err != nil {
				return projectManagementNotFound(err, "project not found")
			}
			template.ProjectType, template.StartAt, template.EndAt, template.PeriodUnit = project.ProjectType, project.StartAt, project.EndAt, project.PeriodUnit
			template.PeriodValue = project.PeriodValue
		}
		template, err = l.repos.ProjectManagementTemplateRepo.CreateWithTx(ctx, tx, template, false)
		if err != nil {
			return err
		}
		if input.SourceProjectId != 0 {
			if err = l.repos.ProjectManagementWorkflowRepo.CloneProjectToTemplateWithTx(ctx, tx, input.SourceProjectId, template.Id); err != nil {
				return err
			}
			if input.CopyTasks {
				return l.snapshotProjectTasksWithTx(ctx, tx, input.SourceProjectId, template.Id, input.CopyAssignees)
			}
			return nil
		}
		return l.repos.ProjectManagementWorkflowRepo.CreateDefaultTemplateWorkflowWithTx(ctx, tx, template.Id)
	})
	return template, err
}

func (l *ProjectManagementTemplateLogic) snapshotProjectTasksWithTx(ctx context.Context, tx *gorm.DB, projectId, templateId int64, copyAssignees bool) error {
	tasks, err := l.repos.ProjectManagementTaskRepo.ListByOwnerWithTx(ctx, tx, constants.ProjectManagementOwnerTypeProject, projectId)
	if err != nil || len(tasks) == 0 {
		return err
	}
	workflows, err := l.repos.ProjectManagementWorkflowRepo.ListTemplateWorkflowsWithTx(ctx, tx, templateId)
	if err != nil {
		return err
	}
	workflowIDs := make(map[int64]int64, len(workflows))
	for _, workflow := range workflows {
		workflowIDs[workflow.SourceWorkflowId] = workflow.Id
	}
	taskIDs := make([]int64, 0, len(tasks))
	newTaskIDs := make(map[int64]int64, len(tasks))
	for _, task := range tasks {
		if err = validateProjectManagementTaskTimeRange(task.StartAt, task.EndAt); err != nil {
			return err
		}
		workflowId := workflowIDs[task.WorkflowId]
		if workflowId == 0 {
			return errorx.New("template workflow mapping is invalid")
		}
		copy := &model.ProjectManagementTask{EntityType: constants.ProjectManagementOwnerTypeTemplate, EntityId: templateId, OrganizationId: task.OrganizationId, WorkflowId: workflowId, CreatorId: task.CreatorId, CreatorName: task.CreatorName, Title: task.Title, Description: task.Description, Priority: task.Priority, WorkHours: task.WorkHours, StartAt: task.StartAt, EndAt: task.EndAt, TaskType: task.TaskType}
		if err = l.repos.ProjectManagementTaskRepo.CreateTemplateEntityWithTx(ctx, tx, copy); err != nil {
			return err
		}
		taskIDs = append(taskIDs, task.Id)
		newTaskIDs[task.Id] = copy.Id
	}
	for _, task := range tasks {
		if task.ParentTaskId == 0 {
			continue
		}
		if err = l.repos.ProjectManagementTaskRepo.UpdateParentWithTx(ctx, tx, newTaskIDs[task.Id], newTaskIDs[task.ParentTaskId]); err != nil {
			return err
		}
	}
	dependencies, err := l.repos.ProjectManagementDependencyRepo.ListByTaskIdsWithTx(ctx, tx, taskIDs)
	if err != nil {
		return err
	}
	for _, dependency := range dependencies {
		if newTaskIDs[dependency.TaskId] == 0 || newTaskIDs[dependency.DependencyTaskId] == 0 {
			return errorx.New("template dependency mapping is invalid")
		}
		if _, err = l.repos.ProjectManagementDependencyRepo.CreateWithTx(ctx, tx, &model.ProjectManagementTaskDependency{TaskId: newTaskIDs[dependency.TaskId], DependencyTaskId: newTaskIDs[dependency.DependencyTaskId], DependencyType: dependency.DependencyType, DependencyRule: dependency.DependencyRule}); err != nil {
			return err
		}
	}
	if !copyAssignees {
		return nil
	}
	assignees, err := l.repos.ProjectManagementTaskRepo.ListCurrentAssigneesByTaskIdsWithTx(ctx, tx, taskIDs)
	if err != nil {
		return err
	}
	copies := make([]*model.ProjectManagementTaskAssignee, 0, len(assignees))
	for _, assignee := range assignees {
		copies = append(copies, &model.ProjectManagementTaskAssignee{TaskId: newTaskIDs[assignee.TaskId], OrganizationUserId: assignee.OrganizationUserId, OrganizationUserFullName: assignee.OrganizationUserFullName})
	}
	return l.repos.ProjectManagementTaskRepo.CreateTemplateAssigneesWithTx(ctx, tx, copies)
}

func (l *ProjectManagementTemplateLogic) Update(ctx context.Context, claim utils.JwtUser, id int64, input dto.ProjectManagementTemplateUpdateDto) (*model.ProjectManagementTemplate, error) {
	if err := validateProjectManagementTemplateStages(input.Stages); err != nil {
		return nil, err
	}
	if input.Status != "" {
		if err := input.Status.Validate(); err != nil {
			return nil, err
		}
	}
	var updated *model.ProjectManagementTemplate
	err := l.repos.ProjectManagementTemplateRepo.Transaction(ctx, func(tx *gorm.DB) error {
		template, err := l.manageableTemplate(ctx, claim, id)
		if err != nil {
			return err
		}
		if input.ProjectType != template.ProjectType {
			return errorx.New("project type cannot be changed")
		}
		if template.IsDefault && input.Status == constants.CommonStatusDisabled {
			return errorx.New("default project template cannot be disabled")
		}
		updated = new(model.ProjectManagementTemplate)
		*updated = *template
		if err = copierx.Assign(&input, updated); err != nil {
			return err
		}
		if input.Status == "" {
			updated.Status = template.Status
		}
		if err = prepareProjectManagementTemplateSchedule(updated); err != nil {
			return err
		}
		if _, err = l.repos.ProjectManagementTemplateRepo.UpdateOneByIdWithTx(ctx, tx, id, updated); err != nil {
			return err
		}
		if err = l.syncTemplateStagesWithTx(ctx, tx, id, input.Stages); err != nil {
			return err
		}
		return l.validateTemplateWorkflowsWithTx(ctx, tx, id)
	})
	if err != nil {
		return nil, err
	}
	updated.Stages, err = l.repos.ProjectManagementStageRepo.ListByOwner(ctx, constants.ProjectManagementOwnerTypeTemplate, id)
	return updated, err
}

func (l *ProjectManagementTemplateLogic) Delete(ctx context.Context, claim utils.JwtUser, id int64) error {
	template, err := l.manageableTemplate(ctx, claim, id)
	if err != nil {
		return err
	}
	if template.IsDefault {
		return errorx.New("can not delete default project template")
	}
	_, err = l.repos.ProjectManagementTemplateRepo.DeleteByID(ctx, id)
	return err
}

func (l *ProjectManagementTemplateLogic) SetStatus(ctx context.Context, claim utils.JwtUser, id int64, status constants.CommonStatus) (*model.ProjectManagementTemplate, error) {
	if err := status.Validate(); err != nil {
		return nil, err
	}
	template, err := l.manageableTemplate(ctx, claim, id)
	if err != nil {
		return nil, err
	}
	if template.IsDefault && status == constants.CommonStatusDisabled {
		return nil, errorx.New("default project template cannot be disabled")
	}
	if err = l.repos.ProjectManagementTemplateRepo.UpdateColumns(ctx, map[string]any{"status": status}, map[string]any{"id": id}, nil); err != nil {
		return nil, err
	}
	template.Status = status
	return template, nil
}

func (l *ProjectManagementTemplateLogic) SetDefault(ctx context.Context, claim utils.JwtUser, id int64) error {
	template, err := l.manageableTemplate(ctx, claim, id)
	if err != nil {
		return err
	}
	return l.repos.ProjectManagementTemplateRepo.Transaction(ctx, func(tx *gorm.DB) error {
		found, err := l.repos.ProjectManagementTemplateRepo.SetDefaultWithTx(ctx, tx, id, template.OrganizationId, template.TemplateType)
		if err != nil {
			return err
		}
		if !found {
			return projectManagementNotFound(gorm.ErrRecordNotFound, "project template not found")
		}
		return nil
	})
}

func (l *ProjectManagementTemplateLogic) Detail(ctx context.Context, claim utils.JwtUser, id int64) (*model.ProjectManagementTemplate, error) {
	template, err := l.repos.ProjectManagementTemplateRepo.FindVisibleProjectTemplate(ctx, id, claim.OrganizationId)
	if err != nil {
		return nil, projectManagementNotFound(err, "project template not found")
	}
	template.Stages, err = l.repos.ProjectManagementStageRepo.ListByOwner(ctx, constants.ProjectManagementOwnerTypeTemplate, id)
	return template, err
}

func (l *ProjectManagementTemplateLogic) Stages(ctx context.Context, claim utils.JwtUser, templateId int64) ([]*model.ProjectManagementStage, error) {
	template, err := l.Detail(ctx, claim, templateId)
	if err != nil {
		return nil, err
	}
	return template.Stages, nil
}

func (l *ProjectManagementTemplateLogic) SwapStageOrder(ctx context.Context, claim utils.JwtUser, templateId, stageId, targetStageId int64) error {
	if stageId == targetStageId {
		return errorx.New("template stages must be different")
	}
	return l.repos.ProjectManagementTemplateRepo.Transaction(ctx, func(tx *gorm.DB) error {
		if _, err := l.manageableTemplate(ctx, claim, templateId); err != nil {
			return err
		}
		firstID, secondID := stageId, targetStageId
		if secondID < firstID {
			firstID, secondID = secondID, firstID
		}
		first, err := l.repos.ProjectManagementStageRepo.LockByOwnerWithTx(ctx, tx, firstID, constants.ProjectManagementOwnerTypeTemplate, templateId)
		if err != nil {
			return projectManagementNotFound(err, "template stage not found")
		}
		second, err := l.repos.ProjectManagementStageRepo.LockByOwnerWithTx(ctx, tx, secondID, constants.ProjectManagementOwnerTypeTemplate, templateId)
		if err != nil {
			return projectManagementNotFound(err, "template stage not found")
		}
		if err = l.repos.ProjectManagementStageRepo.UpdateColumnsWithTx(ctx, tx, map[string]any{"sort_order": second.SortOrder}, map[string]any{"id": first.Id}, nil); err != nil {
			return err
		}
		return l.repos.ProjectManagementStageRepo.UpdateColumnsWithTx(ctx, tx, map[string]any{"sort_order": first.SortOrder}, map[string]any{"id": second.Id}, nil)
	})
}

func (l *ProjectManagementTemplateLogic) manageableTemplate(ctx context.Context, claim utils.JwtUser, id int64) (*model.ProjectManagementTemplate, error) {
	var (
		template *model.ProjectManagementTemplate
		err      error
	)
	if claim.UserType == constants.UserTypeAdmin {
		template, err = l.repos.ProjectManagementTemplateRepo.FindPublicByID(ctx, id)
	} else if claim.OrganizationId != 0 {
		template, err = l.repos.ProjectManagementTemplateRepo.FindByOrganization(ctx, id, claim.OrganizationId)
		if err == nil && template.TemplateType != constants.ProjectManagementTemplateTypeOrganization {
			err = gorm.ErrRecordNotFound
		}
	} else {
		err = gorm.ErrRecordNotFound
	}
	if err != nil {
		return nil, projectManagementNotFound(err, "project template not found")
	}
	return template, nil
}

func (l *ProjectManagementTemplateLogic) syncTemplateStagesWithTx(ctx context.Context, tx *gorm.DB, templateId int64, input []dto.ProjectManagementTemplateStageDto) error {
	existing, err := l.repos.ProjectManagementStageRepo.ListByOwnerWithTx(ctx, tx, constants.ProjectManagementOwnerTypeTemplate, templateId)
	if err != nil {
		return err
	}
	byID := make(map[int64]*model.ProjectManagementStage, len(existing))
	nextSortOrder := 0
	for _, stage := range existing {
		byID[stage.Id] = stage
		nextSortOrder = max(nextSortOrder, stage.SortOrder+1)
	}
	requested := make(map[int64]struct{}, len(input))
	for _, item := range input {
		if item.Id == 0 {
			continue
		}
		stage := byID[item.Id]
		if stage == nil {
			return errorx.New("template stage owner mismatch")
		}
		requested[item.Id] = struct{}{}
	}
	for _, stage := range existing {
		if _, ok := requested[stage.Id]; ok {
			continue
		}
		if stage.Phase == constants.ProjectManagementStagePhasePending {
			return errorx.New("pending template stage cannot be deleted")
		}
		used, err := l.repos.ProjectManagementStageRepo.IsUsedByTemplateWorkflowWithTx(ctx, tx, templateId, stage.Id)
		if err != nil {
			return err
		}
		if used {
			return errorx.New("template stage used by workflow cannot be deleted")
		}
		if _, err = l.repos.ProjectManagementStageRepo.DeleteByIDWithTx(ctx, tx, stage.Id); err != nil {
			return err
		}
	}
	for _, item := range input {
		if item.Id == 0 {
			stage := &model.ProjectManagementStage{EntityType: constants.ProjectManagementOwnerTypeTemplate, EntityId: templateId, Name: item.Name, Description: item.Description, Phase: item.Phase, Color: item.Color, SortOrder: nextSortOrder}
			if _, err = l.repos.ProjectManagementStageRepo.CreateWithTx(ctx, tx, stage, false); err != nil {
				return err
			}
			nextSortOrder++
			continue
		}
		stage := byID[item.Id]
		if stage.Phase == constants.ProjectManagementStagePhasePending && item.Phase != constants.ProjectManagementStagePhasePending {
			return errorx.New("pending template stage phase cannot be changed")
		}
		stage.Name, stage.Description, stage.Phase, stage.Color, stage.SortOrder = item.Name, item.Description, item.Phase, item.Color, item.SortOrder
		if _, err = l.repos.ProjectManagementStageRepo.UpdateOneByIdWithTx(ctx, tx, stage.Id, stage); err != nil {
			return err
		}
	}
	return nil
}

func validateProjectManagementTemplateStages(stages []dto.ProjectManagementTemplateStageDto) error {
	if len(stages) == 0 {
		return errorx.New("template stages are required")
	}
	ids := make(map[int64]struct{}, len(stages))
	names := make(map[string]struct{}, len(stages))
	pending, finished := 0, 0
	for _, stage := range stages {
		name := strings.TrimSpace(stage.Name)
		if name == "" {
			return errorx.New("template stage name is required")
		}
		if _, exists := names[name]; exists {
			return errorx.New("template stage name must be unique")
		}
		names[name] = struct{}{}
		if !isProjectManagementStagePhase(stage.Phase) {
			return errorx.New("invalid template stage phase")
		}
		if stage.SortOrder < 0 {
			return errorx.New("template stage sort order must be >= 0")
		}
		if stage.Id != 0 {
			if _, exists := ids[stage.Id]; exists {
				return errorx.New("template stage id must be unique")
			}
			ids[stage.Id] = struct{}{}
		}
		if stage.Phase == constants.ProjectManagementStagePhasePending {
			pending++
		}
		if stage.Phase == constants.ProjectManagementStagePhaseFinished {
			finished++
		}
	}
	if pending != 1 {
		return errorx.New("template must have exactly one pending stage")
	}
	if finished == 0 {
		return errorx.New("template must have at least one finished stage")
	}
	return nil
}

func (l *ProjectManagementTemplateLogic) validateTemplateWorkflowsWithTx(ctx context.Context, tx *gorm.DB, templateId int64) error {
	workflows, err := l.repos.ProjectManagementWorkflowRepo.ListTemplateWorkflowsWithTx(ctx, tx, templateId)
	if err != nil {
		return err
	}
	for _, workflow := range workflows {
		stages, err := l.repos.ProjectManagementWorkflowRepo.ListWorkflowStagesWithTx(ctx, tx, workflow.Id)
		if err != nil {
			return err
		}
		definition := dto.ProjectManagementWorkflowDto{Name: workflow.Name, Stages: make([]dto.ProjectManagementWorkflowStageDto, 0, len(stages))}
		for _, stage := range stages {
			definition.Stages = append(definition.Stages, dto.ProjectManagementWorkflowStageDto{Id: stage.Id, Name: stage.Name, Description: stage.Description, Phase: stage.Phase, Color: stage.Color, SortOrder: stage.SortOrder})
		}
		if err = ValidateProjectManagementWorkflow(definition); err != nil {
			return err
		}
	}
	return nil
}

func (l *ProjectManagementTemplateLogic) Workflows(ctx context.Context, claim utils.JwtUser, id int64) ([]*model.ProjectManagementWorkflow, error) {
	if _, err := l.Detail(ctx, claim, id); err != nil {
		return nil, err
	}
	return l.repos.ProjectManagementTemplateRepo.ListWorkflows(ctx, id)
}

func (l *ProjectManagementTemplateLogic) Discover(ctx context.Context, claim utils.JwtUser) ([]*model.ProjectManagementTemplate, error) {
	public, err := l.repos.ProjectManagementTemplateRepo.ListPublicTemplates(ctx, constants.CommonStatusActive)
	if err != nil {
		return nil, err
	}
	organization, err := l.repos.ProjectManagementTemplateRepo.ListProjectTemplates(ctx, claim.OrganizationId)
	if err != nil {
		return nil, err
	}
	result := make([]*model.ProjectManagementTemplate, 0, len(public)+len(organization))
	result = append(result, public...)
	for _, template := range organization {
		if template.Status == constants.CommonStatusActive {
			result = append(result, template)
		}
	}
	return result, nil
}

func prepareProjectManagementTemplateSchedule(template *model.ProjectManagementTemplate) error {
	switch template.ProjectType {
	case constants.ProjectManagementProjectTypeOrdinary:
		template.StartAt, template.EndAt, template.PeriodValue, template.PeriodUnit = nil, nil, 0, ""
	case constants.ProjectManagementProjectTypePlanned:
		if template.StartAt == nil || !template.StartAt.Valid {
			return errorx.New("planned template start date is required")
		}
		if template.PeriodValue <= 0 {
			return errorx.New("planned template period value must be positive")
		}
		endAt, err := CalculateProjectEndAt(template.StartAt.Time, template.PeriodValue, template.PeriodUnit)
		if err != nil {
			return err
		}
		if template.EndAt != nil && template.EndAt.Valid && !template.EndAt.Time.Equal(endAt) {
			return errorx.New("planned template end date must match its start date and period")
		}
		template.EndAt = &sql.NullTime{Time: endAt, Valid: true}
	default:
		return errorx.New("invalid template project type")
	}
	return nil
}
