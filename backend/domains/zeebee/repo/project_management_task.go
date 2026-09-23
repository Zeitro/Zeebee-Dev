package repo

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"zeebee-backend/common/errorx"
	"zeebee-backend/common/gormx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
)

type ProjectManagementTaskRepo struct {
	*gormx.BaseRepo[*model.ProjectManagementTask]
	assigneeRepo *gormx.BaseRepo[*model.ProjectManagementTaskAssignee]
}

func NewProjectManagementTaskRepo(db *gorm.DB) *ProjectManagementTaskRepo {
	return &ProjectManagementTaskRepo{
		BaseRepo:     gormx.NewBaseRepo[*model.ProjectManagementTask](db),
		assigneeRepo: gormx.NewBaseRepo[*model.ProjectManagementTaskAssignee](db),
	}
}

func (r *ProjectManagementTaskRepo) FindByOrganization(ctx context.Context, id, organizationId int64) (*model.ProjectManagementTask, error) {
	task, err := r.FindOneByColumns(ctx, map[string]any{"id": id, "organization_id": organizationId, "entity_type": constants.ProjectManagementOwnerTypeProject})
	return task, err
}

func (r *ProjectManagementTaskRepo) ListByWorkflowWithTx(ctx context.Context, tx *gorm.DB, workflowId, organizationId int64, history bool) ([]*model.ProjectManagementTask, error) {
	tasks := make([]*model.ProjectManagementTask, 0)
	query := tx.WithContext(ctx).Model(&model.ProjectManagementTask{}).
		Select("id, created_at, updated_at, is_del, entity_id, entity_type, workflow_id, organization_id, parent_task_id, stage_id, phase, creator_id, creator_name, title, description, priority, work_hours, start_at, end_at, task_type").
		Where("workflow_id = ? AND organization_id = ? AND entity_type = ?", workflowId, organizationId, constants.ProjectManagementOwnerTypeProject)
	if history {
		query = query.Unscoped().Where("phase = ?", constants.ProjectManagementStagePhaseFinished)
	} else {
		query = query.Where("phase <> ?", constants.ProjectManagementStagePhaseFinished)
	}
	err := query.Order("id ASC").Find(&tasks).Error
	return tasks, errorx.WithStack(err)
}

func (r *ProjectManagementTaskRepo) ListByProjectForView(ctx context.Context, projectId, organizationId int64) ([]*model.ProjectManagementTask, error) {
	tasks := make([]*model.ProjectManagementTask, 0)
	err := r.Client(ctx).Model(&model.ProjectManagementTask{}).
		Select("project_management_tasks.*, project_management_workflows.name AS workflow_name").
		Joins("LEFT JOIN project_management_workflows ON project_management_workflows.id = project_management_tasks.workflow_id").
		Where("project_management_tasks.entity_type = ? AND project_management_tasks.entity_id = ? AND project_management_tasks.organization_id = ?", constants.ProjectManagementOwnerTypeProject, projectId, organizationId).
		Order("project_management_tasks.parent_task_id ASC, project_management_tasks.id ASC").Find(&tasks).Error
	return tasks, errorx.WithStack(err)
}

func (r *ProjectManagementTaskRepo) ListByOwnerWithTx(ctx context.Context, tx *gorm.DB, entityType constants.ProjectManagementOwnerType, entityId int64) ([]*model.ProjectManagementTask, error) {
	tasks := make([]*model.ProjectManagementTask, 0)
	err := tx.WithContext(ctx).Model(&model.ProjectManagementTask{}).
		Where("entity_type = ? AND entity_id = ?", entityType, entityId).
		Order("id ASC").Find(&tasks).Error
	return tasks, errorx.WithStack(err)
}

func (r *ProjectManagementTaskRepo) UpdateParentWithTx(ctx context.Context, tx *gorm.DB, taskId, parentTaskId int64) error {
	return r.UpdateColumnsWithTx(ctx, tx, map[string]any{"parent_task_id": parentTaskId}, map[string]any{"id": taskId}, nil)
}

func (r *ProjectManagementTaskRepo) FindDetailByOrganization(ctx context.Context, id, organizationId int64) (*model.ProjectManagementTask, error) {
	task := new(model.ProjectManagementTask)
	err := r.Client(ctx).Model(&model.ProjectManagementTask{}).
		Select("project_management_tasks.*, project_management_workflows.name AS workflow_name").
		Joins("LEFT JOIN project_management_workflows ON project_management_workflows.id = project_management_tasks.workflow_id").
		Preload("Children").
		Preload("Attachments").
		Preload("Dependencies.DependencyTask").
		Where("project_management_tasks.id = ? AND project_management_tasks.organization_id = ?", id, organizationId).First(task).Error
	if err == nil && task.EntityType == constants.ProjectManagementOwnerTypeProject {
		err = r.assigneeRepo.Client(ctx).Where("task_id = ? AND stage_id = ?", task.Id, task.StageId).Order("id ASC").Find(&task.Assignees).Error
		return task, errorx.WithStack(err)
	}
	return task, errorx.WithStack(err)
}

func (r *ProjectManagementTaskRepo) PageByProject(ctx context.Context, organizationId int64, input dto.ProjectManagementTaskPageRequestDto) (gormx.PageResponse[*model.ProjectManagementTask], error) {
	tasksTable := model.ProjectManagementTask{}.TableName()
	assigneesTable := model.ProjectManagementTaskAssignee{}.TableName()
	stagesTable := model.ProjectManagementStage{}.TableName()
	projectsTable := model.ProjectManagementProject{}.TableName()
	workflowsTable := model.ProjectManagementWorkflow{}.TableName()
	query := r.Client(ctx).
		Omit("Description").
		Joins(fmt.Sprintf("JOIN %s ON %s.id = %s.stage_id AND %s.entity_type = ? AND %s.is_del = 0", stagesTable, stagesTable, tasksTable, stagesTable, stagesTable), constants.ProjectManagementOwnerTypeProject).
		Joins(fmt.Sprintf("JOIN %s ON %s.id = %s.entity_id AND %s.is_del = 0", projectsTable, projectsTable, tasksTable, projectsTable)).
		Joins(fmt.Sprintf("JOIN %s ON %s.id = %s.workflow_id", workflowsTable, workflowsTable, tasksTable)).
		Select(fmt.Sprintf("%s.*, %s.name AS project_name, %s.name AS stage_name, %s.name AS workflow_name", tasksTable, projectsTable, stagesTable, workflowsTable)).
		Where(fmt.Sprintf("%s.organization_id = ? AND %s.entity_type = ? AND %s.entity_id = ?", tasksTable, tasksTable, tasksTable), organizationId, constants.ProjectManagementOwnerTypeProject, input.Filter.ProjectId).
		Order(fmt.Sprintf("%s.sort_order ASC, %s.id DESC", stagesTable, tasksTable))
	if input.Filter.StageId != 0 {
		query = query.Where(fmt.Sprintf("%s.stage_id = ?", tasksTable), input.Filter.StageId)
	}
	if len(input.Filter.Phase) > 0 {
		query = query.Where(fmt.Sprintf("%s.phase IN ?", tasksTable), input.Filter.Phase)
	}
	if input.Filter.Priority != "" {
		query = query.Where(fmt.Sprintf("%s.priority = ?", tasksTable), input.Filter.Priority)
	}
	if input.Filter.TaskTypeId != 0 {
		query = query.Where(fmt.Sprintf("CAST(%s.task_type->>'id' AS TEXT) = CAST(? AS TEXT)", tasksTable), input.Filter.TaskTypeId)
	}
	if input.Filter.AssigneeId != 0 {
		query = query.Where(fmt.Sprintf("EXISTS (SELECT 1 FROM %s WHERE %s.task_id = %s.id AND %s.stage_id = %s.stage_id AND %s.organization_user_id = ? AND %s.is_del = 0)", assigneesTable, assigneesTable, tasksTable, assigneesTable, tasksTable, assigneesTable, assigneesTable), input.Filter.AssigneeId)
	}
	if input.Filter.Keyword != "" {
		query = query.Where(fmt.Sprintf("LOWER(%s.title) LIKE LOWER(?)", tasksTable), "%"+input.Filter.Keyword+"%")
	}
	if input.Filter.Overdue != nil {
		if *input.Filter.Overdue {
			query = query.Where(fmt.Sprintf("%s.end_at IS NOT NULL AND %s.end_at < CURRENT_TIMESTAMP AND %s.phase <> ?", tasksTable, tasksTable, stagesTable), constants.ProjectManagementStagePhaseFinished)
		} else {
			query = query.Where(fmt.Sprintf("(%s.end_at IS NULL OR %s.end_at >= CURRENT_TIMESTAMP OR %s.phase = ?)", tasksTable, tasksTable, stagesTable), constants.ProjectManagementStagePhaseFinished)
		}
	}
	return r.BaseRepo.Page(query, gormx.PageRequest{Page: input.Page, Limit: input.Limit})
}

func (r *ProjectManagementTaskRepo) CandidateByProject(ctx context.Context, organizationId int64, input dto.ProjectManagementTaskCandidateDto) ([]*dto.ProjectManagementTaskCandidateEntityDto, error) {
	query := r.Client(ctx).
		Select("id, title, phase").
		Where("organization_id = ? AND entity_type = ? AND entity_id = ? AND phase IN ?", organizationId, constants.ProjectManagementOwnerTypeProject, input.EntityId, []constants.ProjectManagementStagePhase{constants.ProjectManagementStagePhasePending, constants.ProjectManagementStagePhaseInProgress}).
		Order("title ASC, id ASC")
	if input.TaskId != 0 {
		query = query.Where("id <> ?", input.TaskId)
	}
	keyword := strings.TrimSpace(input.Keyword)
	if keyword != "" {
		pattern := "%" + keyword + "%"
		if id, err := strconv.ParseInt(keyword, 10, 64); err == nil && id > 0 {
			query = query.Where("id = ? OR LOWER(title) LIKE LOWER(?)", id, pattern)
		} else {
			query = query.Where("LOWER(title) LIKE LOWER(?)", pattern)
		}
	}
	candidates := make([]*dto.ProjectManagementTaskCandidateEntityDto, 0)
	err := query.Find(&candidates).Error
	return candidates, errorx.WithStack(err)
}

func (r *ProjectManagementTaskRepo) PageMine(ctx context.Context, organizationId, organizationUserId int64, input dto.ProjectManagementMyTaskPageRequestDto) (gormx.PageResponse[*model.ProjectManagementTask], error) {
	tasksTable := model.ProjectManagementTask{}.TableName()
	assigneesTable := model.ProjectManagementTaskAssignee{}.TableName()
	stagesTable := model.ProjectManagementStage{}.TableName()
	projectsTable := model.ProjectManagementProject{}.TableName()
	workflowsTable := model.ProjectManagementWorkflow{}.TableName()
	query := r.Client(ctx).
		Omit("Description").
		Joins(fmt.Sprintf("JOIN %s ON %s.id = %s.stage_id AND %s.entity_type = ? AND %s.is_del = 0", stagesTable, stagesTable, tasksTable, stagesTable, stagesTable), constants.ProjectManagementOwnerTypeProject).
		Joins(fmt.Sprintf("JOIN %s ON %s.id = %s.entity_id AND %s.is_del = 0", projectsTable, projectsTable, tasksTable, projectsTable)).
		Joins(fmt.Sprintf("JOIN %s ON %s.id = %s.workflow_id", workflowsTable, workflowsTable, tasksTable)).
		Select(fmt.Sprintf("%s.*, %s.name AS project_name, %s.name AS stage_name, %s.name AS workflow_name", tasksTable, projectsTable, stagesTable, workflowsTable)).
		Where(fmt.Sprintf("%s.organization_id = ?", tasksTable), organizationId).
		Where(fmt.Sprintf("EXISTS (SELECT 1 FROM %s mine WHERE mine.task_id = %s.id AND mine.stage_id = %s.stage_id AND mine.organization_user_id = ? AND mine.is_del = 0)", assigneesTable, tasksTable, tasksTable), organizationUserId).
		Order(fmt.Sprintf("%s.sort_order ASC, %s.id DESC", stagesTable, tasksTable))
	if input.Filter.ProjectId != 0 {
		query = query.Where(fmt.Sprintf("%s.entity_type = ? AND %s.entity_id = ?", tasksTable, tasksTable), constants.ProjectManagementOwnerTypeProject, input.Filter.ProjectId)
	}
	if input.Filter.StageId != 0 {
		query = query.Where(fmt.Sprintf("%s.stage_id = ?", tasksTable), input.Filter.StageId)
	}
	if len(input.Filter.Phase) > 0 {
		query = query.Where(fmt.Sprintf("%s.phase IN ?", tasksTable), input.Filter.Phase)
	}
	if input.Filter.Priority != "" {
		query = query.Where(fmt.Sprintf("%s.priority = ?", tasksTable), input.Filter.Priority)
	}
	if input.Filter.TaskTypeId != 0 {
		query = query.Where(fmt.Sprintf("CAST(%s.task_type->>'id' AS TEXT) = CAST(? AS TEXT)", tasksTable), input.Filter.TaskTypeId)
	}
	if input.Filter.AssigneeId != 0 {
		query = query.Where(fmt.Sprintf("EXISTS (SELECT 1 FROM %s WHERE %s.task_id = %s.id AND %s.stage_id = %s.stage_id AND %s.organization_user_id = ? AND %s.is_del = 0)", assigneesTable, assigneesTable, tasksTable, assigneesTable, tasksTable, assigneesTable, assigneesTable), input.Filter.AssigneeId)
	}
	if input.Filter.Keyword != "" {
		query = query.Where(fmt.Sprintf("LOWER(%s.title) LIKE LOWER(?)", tasksTable), "%"+input.Filter.Keyword+"%")
	}
	if input.Filter.Overdue != nil {
		if *input.Filter.Overdue {
			query = query.Where(fmt.Sprintf("%s.end_at IS NOT NULL AND %s.end_at < CURRENT_TIMESTAMP AND %s.phase <> ?", tasksTable, tasksTable, stagesTable), constants.ProjectManagementStagePhaseFinished)
		} else {
			query = query.Where(fmt.Sprintf("(%s.end_at IS NULL OR %s.end_at >= CURRENT_TIMESTAMP OR %s.phase = ?)", tasksTable, tasksTable, stagesTable), constants.ProjectManagementStagePhaseFinished)
		}
	}
	return r.BaseRepo.Page(query, gormx.PageRequest{Page: input.Page, Limit: input.Limit})
}

func (r *ProjectManagementTaskRepo) LockWithTx(ctx context.Context, tx *gorm.DB, id int64) (*model.ProjectManagementTask, error) {
	task := new(model.ProjectManagementTask)
	err := tx.WithContext(ctx).Model(&model.ProjectManagementTask{}).
		Select("id, created_at, updated_at, is_del, entity_id, entity_type, workflow_id, organization_id, parent_task_id, stage_id, phase, creator_id, creator_name, title, description, priority, work_hours, start_at, end_at, task_type").
		Clauses(clause.Locking{Strength: "UPDATE"}).First(task, id).Error
	if err == nil && task.EntityType == constants.ProjectManagementOwnerTypeProject {
		return task, nil
	}
	return task, errorx.WithStack(err)
}

func currentStageAssignees(db *gorm.DB) *gorm.DB {
	return db.Joins("JOIN project_management_tasks ON project_management_tasks.id = project_management_task_assignees.task_id AND project_management_tasks.stage_id = project_management_task_assignees.stage_id").Order("project_management_task_assignees.id ASC")
}

func (r *ProjectManagementTaskRepo) IsAssigneeWithTx(ctx context.Context, tx *gorm.DB, taskId, stageId, organizationUserId int64) (bool, error) {
	return r.assigneeRepo.IsExistedByColumnsWithTx(ctx, tx, map[string]any{"task_id": taskId, "stage_id": stageId, "organization_user_id": organizationUserId})
}

func (r *ProjectManagementTaskRepo) ListAssigneesWithTx(ctx context.Context, tx *gorm.DB, taskId, stageId int64) ([]*model.ProjectManagementTaskAssignee, error) {
	assignees := make([]*model.ProjectManagementTaskAssignee, 0)
	err := tx.WithContext(ctx).Where("task_id = ? AND stage_id = ?", taskId, stageId).Find(&assignees).Error
	return assignees, errorx.WithStack(err)
}

func (r *ProjectManagementTaskRepo) ListCurrentAssigneesByTaskIds(ctx context.Context, taskIds []int64) ([]*model.ProjectManagementTaskAssignee, error) {
	return r.listCurrentAssigneesByTaskIds(ctx, r.assigneeRepo.Client(ctx), taskIds)
}

func (r *ProjectManagementTaskRepo) ListCurrentAssigneesByTaskIdsWithTx(ctx context.Context, tx *gorm.DB, taskIds []int64) ([]*model.ProjectManagementTaskAssignee, error) {
	return r.listCurrentAssigneesByTaskIds(ctx, tx, taskIds)
}

func (r *ProjectManagementTaskRepo) listCurrentAssigneesByTaskIds(ctx context.Context, db *gorm.DB, taskIds []int64) ([]*model.ProjectManagementTaskAssignee, error) {
	if len(taskIds) == 0 {
		return []*model.ProjectManagementTaskAssignee{}, nil
	}
	assignees := make([]*model.ProjectManagementTaskAssignee, 0)
	err := currentStageAssignees(db.WithContext(ctx)).Where("project_management_task_assignees.task_id IN ?", taskIds).Find(&assignees).Error
	return assignees, errorx.WithStack(err)
}

func (r *ProjectManagementTaskRepo) ListTemplateAssigneesByTaskIdsWithTx(ctx context.Context, tx *gorm.DB, taskIds []int64) ([]*model.ProjectManagementTaskAssignee, error) {
	if len(taskIds) == 0 {
		return []*model.ProjectManagementTaskAssignee{}, nil
	}
	assignees := make([]*model.ProjectManagementTaskAssignee, 0)
	err := tx.WithContext(ctx).Model(&model.ProjectManagementTaskAssignee{}).
		Where("task_id IN ? AND stage_id = ?", taskIds, 0).Order("id ASC").Find(&assignees).Error
	return assignees, errorx.WithStack(err)
}

func (r *ProjectManagementTaskRepo) DirectChildrenWithTx(ctx context.Context, tx *gorm.DB, taskId int64) ([]*model.ProjectManagementTask, error) {
	children := make([]*model.ProjectManagementTask, 0)
	err := tx.WithContext(ctx).Where("parent_task_id = ?", taskId).Order("id ASC").Find(&children).Error
	return children, errorx.WithStack(err)
}

func (r *ProjectManagementTaskRepo) IdsByProjectWithTx(ctx context.Context, tx *gorm.DB, projectId int64) ([]int64, error) {
	ids := make([]int64, 0)
	err := r.PluckByColumnsWithTx(ctx, tx, "id", &ids, []gormx.WhereRaw{{Raw: "entity_type = ? AND entity_id = ?", Args: []interface{}{constants.ProjectManagementOwnerTypeProject, projectId}}})
	return ids, err
}

func (r *ProjectManagementTaskRepo) DeleteByProjectWithTx(ctx context.Context, tx *gorm.DB, projectId int64, taskIds []int64) error {
	if len(taskIds) > 0 {
		if err := tx.WithContext(ctx).Where("task_id IN ?", taskIds).Delete(&model.ProjectManagementTaskAssignee{}).Error; err != nil {
			return errorx.WithStack(err)
		}
		if err := tx.WithContext(ctx).Where("task_id IN ? OR dependency_task_id IN ?", taskIds, taskIds).Delete(&model.ProjectManagementTaskDependency{}).Error; err != nil {
			return errorx.WithStack(err)
		}
	}
	_, err := r.DeleteWithTx(ctx, tx, map[string]any{"entity_type": constants.ProjectManagementOwnerTypeProject, "entity_id": projectId})
	return err
}

func (r *ProjectManagementTaskRepo) CreateAssigneesWithTx(ctx context.Context, tx *gorm.DB, assignees []*model.ProjectManagementTaskAssignee) error {
	if len(assignees) == 0 {
		return nil
	}
	_, err := r.assigneeRepo.BatchCreateWithTx(ctx, tx, assignees, false)
	return err
}

func (r *ProjectManagementTaskRepo) CreateEntityWithTx(ctx context.Context, tx *gorm.DB, task *model.ProjectManagementTask) error {
	return errorx.WithStack(tx.WithContext(ctx).Create(task).Error)
}

func (r *ProjectManagementTaskRepo) CreateTemplateEntityWithTx(ctx context.Context, tx *gorm.DB, task *model.ProjectManagementTask) error {
	return errorx.WithStack(tx.WithContext(ctx).Omit("StageId", "Phase").Create(task).Error)
}

func (r *ProjectManagementTaskRepo) CreateTemplateAssigneesWithTx(ctx context.Context, tx *gorm.DB, assignees []*model.ProjectManagementTaskAssignee) error {
	if len(assignees) == 0 {
		return nil
	}
	return errorx.WithStack(tx.WithContext(ctx).Omit("StageId").Create(&assignees).Error)
}

func (r *ProjectManagementTaskRepo) CopyCurrentStageAssigneesIfTargetEmptyWithTx(ctx context.Context, tx *gorm.DB, task *model.ProjectManagementTask, targetStageId int64) error {
	hasTargetAssignees, err := r.assigneeRepo.IsExistedByColumnsWithTx(ctx, tx, map[string]any{"task_id": task.Id, "stage_id": targetStageId})
	if err != nil || hasTargetAssignees {
		return err
	}
	currentAssignees := make([]*model.ProjectManagementTaskAssignee, 0)
	if err = tx.WithContext(ctx).Where("task_id = ? AND stage_id = ?", task.Id, task.StageId).Find(&currentAssignees).Error; err != nil {
		return errorx.WithStack(err)
	}
	targetAssignees := make([]*model.ProjectManagementTaskAssignee, 0, len(currentAssignees))
	for _, assignee := range currentAssignees {
		targetAssignees = append(targetAssignees, &model.ProjectManagementTaskAssignee{TaskId: task.Id, StageId: &targetStageId, OrganizationUserId: assignee.OrganizationUserId, OrganizationUserFullName: assignee.OrganizationUserFullName})
	}
	return r.CreateAssigneesWithTx(ctx, tx, targetAssignees)
}

func (r *ProjectManagementTaskRepo) ReplaceAssigneesWithTx(ctx context.Context, tx *gorm.DB, taskId, stageId int64, assignees []*model.ProjectManagementTaskAssignee) error {
	existing := make([]*model.ProjectManagementTaskAssignee, 0)
	if err := tx.WithContext(ctx).Where("task_id = ? AND stage_id = ?", taskId, stageId).Find(&existing).Error; err != nil {
		return errorx.WithStack(err)
	}
	existingByOrganizationUserId := make(map[int64]struct{}, len(existing))
	requestedByOrganizationUserId := make(map[int64]struct{}, len(assignees))
	deleteIds := make([]int64, 0)
	for _, assignee := range existing {
		existingByOrganizationUserId[assignee.OrganizationUserId] = struct{}{}
	}
	for _, assignee := range assignees {
		requestedByOrganizationUserId[assignee.OrganizationUserId] = struct{}{}
	}
	for _, assignee := range existing {
		if _, exists := requestedByOrganizationUserId[assignee.OrganizationUserId]; !exists {
			deleteIds = append(deleteIds, assignee.Id)
		}
	}
	if len(deleteIds) > 0 {
		if err := tx.WithContext(ctx).Where("id IN ?", deleteIds).Delete(&model.ProjectManagementTaskAssignee{}).Error; err != nil {
			return errorx.WithStack(err)
		}
	}
	additions := make([]*model.ProjectManagementTaskAssignee, 0)
	for _, assignee := range assignees {
		if _, exists := existingByOrganizationUserId[assignee.OrganizationUserId]; !exists {
			additions = append(additions, assignee)
		}
	}
	if err := r.CreateAssigneesWithTx(ctx, tx, additions); err != nil {
		return errorx.WithStack(err)
	}
	return nil
}

func (r *ProjectManagementTaskRepo) DeleteUnfinishedAssigneesByProjectMemberWithTx(ctx context.Context, tx *gorm.DB, projectId, organizationUserId int64) error {
	tasksTable := model.ProjectManagementTask{}.TableName()
	stagesTable := model.ProjectManagementStage{}.TableName()
	taskIds := tx.WithContext(ctx).Model(&model.ProjectManagementTask{}).
		Select(fmt.Sprintf("%s.id", tasksTable)).
		Joins(fmt.Sprintf("JOIN %s ON %s.id = %s.stage_id AND %s.entity_type = ? AND %s.is_del = 0", stagesTable, stagesTable, tasksTable, stagesTable, stagesTable), constants.ProjectManagementOwnerTypeProject).
		Where(fmt.Sprintf("%s.entity_type = ? AND %s.entity_id = ? AND %s.is_del = 0 AND %s.phase <> ?", tasksTable, tasksTable, tasksTable, stagesTable), constants.ProjectManagementOwnerTypeProject, projectId, constants.ProjectManagementStagePhaseFinished)
	if err := tx.WithContext(ctx).Where("organization_user_id = ? AND task_id IN (?)", organizationUserId, taskIds).Delete(&model.ProjectManagementTaskAssignee{}).Error; err != nil {
		return errorx.WithStack(err)
	}
	return nil
}

func (r *ProjectManagementTaskRepo) DescendantIdsWithTx(ctx context.Context, tx *gorm.DB, taskId int64) ([]int64, error) {
	ids, frontier := []int64{taskId}, []int64{taskId}
	for len(frontier) > 0 {
		children := make([]int64, 0)
		if err := tx.WithContext(ctx).Model(&model.ProjectManagementTask{}).Where("parent_task_id IN ?", frontier).Pluck("id", &children).Error; err != nil {
			return nil, errorx.WithStack(err)
		}
		ids, frontier = append(ids, children...), children
	}
	return ids, nil
}

func (r *ProjectManagementTaskRepo) DeleteByIdsWithTx(ctx context.Context, tx *gorm.DB, taskIds []int64) error {
	if len(taskIds) == 0 {
		return nil
	}
	if err := tx.WithContext(ctx).Where("task_id IN ?", taskIds).Delete(&model.ProjectManagementTaskAssignee{}).Error; err != nil {
		return errorx.WithStack(err)
	}
	if err := tx.WithContext(ctx).Where("task_id IN ? OR dependency_task_id IN ?", taskIds, taskIds).Delete(&model.ProjectManagementTaskDependency{}).Error; err != nil {
		return errorx.WithStack(err)
	}
	_, err := r.DeleteByIDsWithTx(ctx, tx, taskIds)
	return err
}
