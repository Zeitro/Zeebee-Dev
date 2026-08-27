package repo

import (
	"context"
	"fmt"

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
	return r.FindOneByColumns(ctx, map[string]any{"id": id, "organization_id": organizationId})
}

func (r *ProjectManagementTaskRepo) FindDetailByOrganization(ctx context.Context, id, organizationId int64) (*model.ProjectManagementTask, error) {
	task := new(model.ProjectManagementTask)
	err := r.Client(ctx).
		Preload("Children").
		Preload("Attachments").
		Preload("Dependencies").
		Where("id = ? AND organization_id = ?", id, organizationId).First(task).Error
	if err == nil {
		err = r.assigneeRepo.Client(ctx).Where("task_id = ? AND stage_id = ?", task.Id, task.StageId).Order("id ASC").Find(&task.Assignees).Error
	}
	return task, errorx.WithStack(err)
}

func (r *ProjectManagementTaskRepo) PageByProject(ctx context.Context, organizationId int64, input dto.ProjectManagementTaskPageRequestDto) (gormx.PageResponse[*model.ProjectManagementTask], error) {
	tasksTable := model.ProjectManagementTask{}.TableName()
	assigneesTable := model.ProjectManagementTaskAssignee{}.TableName()
	stagesTable := model.ProjectManagementWorkflowStage{}.TableName()
	projectsTable := model.ProjectManagementProject{}.TableName()
	query := r.Client(ctx).
		Omit("Description").
		Joins(fmt.Sprintf("JOIN %s ON %s.id = %s.stage_id AND %s.is_del = 0", stagesTable, stagesTable, tasksTable, stagesTable)).
		Joins(fmt.Sprintf("JOIN %s ON %s.id = %s.project_id AND %s.is_del = 0", projectsTable, projectsTable, tasksTable, projectsTable)).
		Select(fmt.Sprintf("%s.*, %s.name AS project_name, %s.name AS stage_name", tasksTable, projectsTable, stagesTable)).
		Where(fmt.Sprintf("%s.organization_id = ? AND %s.project_id = ?", tasksTable, tasksTable), organizationId, input.Filter.ProjectId).
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
		Select("id, title").
		Where("organization_id = ? AND project_id = ? AND phase IN ?", organizationId, input.ProjectId, []constants.ProjectManagementStagePhase{constants.ProjectManagementStagePhasePending, constants.ProjectManagementStagePhaseInProgress}).
		Order("title ASC, id ASC")
	if input.TaskId != 0 {
		query = query.Where("id <> ?", input.TaskId)
	}
	if input.Keyword != "" {
		query = query.Where("LOWER(title) LIKE LOWER(?)", "%"+input.Keyword+"%")
	}
	candidates := make([]*dto.ProjectManagementTaskCandidateEntityDto, 0)
	err := query.Find(&candidates).Error
	return candidates, errorx.WithStack(err)
}

func (r *ProjectManagementTaskRepo) PageMine(ctx context.Context, organizationId, organizationUserId int64, input dto.ProjectManagementMyTaskPageRequestDto) (gormx.PageResponse[*model.ProjectManagementTask], error) {
	tasksTable := model.ProjectManagementTask{}.TableName()
	assigneesTable := model.ProjectManagementTaskAssignee{}.TableName()
	stagesTable := model.ProjectManagementWorkflowStage{}.TableName()
	projectsTable := model.ProjectManagementProject{}.TableName()
	query := r.Client(ctx).
		Omit("Description").
		Joins(fmt.Sprintf("JOIN %s ON %s.id = %s.stage_id AND %s.is_del = 0", stagesTable, stagesTable, tasksTable, stagesTable)).
		Joins(fmt.Sprintf("JOIN %s ON %s.id = %s.project_id AND %s.is_del = 0", projectsTable, projectsTable, tasksTable, projectsTable)).
		Select(fmt.Sprintf("%s.*, %s.name AS project_name, %s.name AS stage_name", tasksTable, projectsTable, stagesTable)).
		Where(fmt.Sprintf("%s.organization_id = ?", tasksTable), organizationId).
		Where(fmt.Sprintf("EXISTS (SELECT 1 FROM %s mine WHERE mine.task_id = %s.id AND mine.stage_id = %s.stage_id AND mine.organization_user_id = ? AND mine.is_del = 0)", assigneesTable, tasksTable, tasksTable), organizationUserId).
		Order(fmt.Sprintf("%s.sort_order ASC, %s.id DESC", stagesTable, tasksTable))
	if input.Filter.ProjectId != 0 {
		query = query.Where(fmt.Sprintf("%s.project_id = ?", tasksTable), input.Filter.ProjectId)
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
	err := tx.WithContext(ctx).Clauses(clause.Locking{Strength: "UPDATE"}).First(task, id).Error
	return task, errorx.WithStack(err)
}

func currentStageAssignees(db *gorm.DB) *gorm.DB {
	return db.Joins("JOIN project_management_tasks ON project_management_tasks.id = project_management_task_assignees.task_id AND project_management_tasks.stage_id = project_management_task_assignees.stage_id").Order("project_management_task_assignees.id ASC")
}

func (r *ProjectManagementTaskRepo) IsAssigneeWithTx(ctx context.Context, tx *gorm.DB, taskId, stageId, organizationUserId int64) (bool, error) {
	return r.assigneeRepo.IsExistedByColumnsWithTx(ctx, tx, map[string]any{"task_id": taskId, "stage_id": stageId, "organization_user_id": organizationUserId})
}

func (r *ProjectManagementTaskRepo) DirectChildrenWithTx(ctx context.Context, tx *gorm.DB, taskId int64) ([]*model.ProjectManagementTask, error) {
	children := make([]*model.ProjectManagementTask, 0)
	err := tx.WithContext(ctx).Where("parent_task_id = ?", taskId).Order("id ASC").Find(&children).Error
	return children, errorx.WithStack(err)
}

func (r *ProjectManagementTaskRepo) IdsByProjectWithTx(ctx context.Context, tx *gorm.DB, projectId int64) ([]int64, error) {
	ids := make([]int64, 0)
	err := r.PluckByColumnsWithTx(ctx, tx, "id", &ids, []gormx.WhereRaw{{Raw: "project_id = ?", Args: []interface{}{projectId}}})
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
	_, err := r.DeleteWithTx(ctx, tx, map[string]any{"project_id": projectId})
	return err
}

func (r *ProjectManagementTaskRepo) CreateAssigneesWithTx(ctx context.Context, tx *gorm.DB, assignees []*model.ProjectManagementTaskAssignee) error {
	if len(assignees) == 0 {
		return nil
	}
	_, err := r.assigneeRepo.BatchCreateWithTx(ctx, tx, assignees, false)
	return err
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
		targetAssignees = append(targetAssignees, &model.ProjectManagementTaskAssignee{TaskId: task.Id, StageId: targetStageId, OrganizationUserId: assignee.OrganizationUserId, OrganizationUserName: assignee.OrganizationUserName})
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
	stagesTable := model.ProjectManagementWorkflowStage{}.TableName()
	taskIds := tx.WithContext(ctx).Table(tasksTable).
		Select(fmt.Sprintf("%s.id", tasksTable)).
		Joins(fmt.Sprintf("JOIN %s ON %s.id = %s.stage_id AND %s.is_del = 0", stagesTable, stagesTable, tasksTable, stagesTable)).
		Where(fmt.Sprintf("%s.project_id = ? AND %s.is_del = 0 AND %s.phase <> ?", tasksTable, tasksTable, stagesTable), projectId, constants.ProjectManagementStagePhaseFinished)
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
