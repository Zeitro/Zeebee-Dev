package repo

import (
	"context"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"zeebee-backend/common/errorx"
	"zeebee-backend/common/gormx"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
)

type ProjectManagementTaskRepo struct {
	*gormx.BaseRepo[*model.ProjectManagementTask]
}

func NewProjectManagementTaskRepo(db *gorm.DB) *ProjectManagementTaskRepo {
	return &ProjectManagementTaskRepo{BaseRepo: gormx.NewBaseRepo[*model.ProjectManagementTask](db)}
}

func (r *ProjectManagementTaskRepo) FindByOrganization(ctx context.Context, id, organizationId int64) (*model.ProjectManagementTask, error) {
	task := new(model.ProjectManagementTask)
	err := r.Client(ctx).Where("id = ? AND organization_id = ?", id, organizationId).First(task).Error
	return task, errorx.WithStack(err)
}

func (r *ProjectManagementTaskRepo) FindDetailByOrganization(ctx context.Context, id, organizationId int64) (*model.ProjectManagementTask, error) {
	task := new(model.ProjectManagementTask)
	err := r.Client(ctx).Preload("Assignees", func(db *gorm.DB) *gorm.DB { return db.Order("id ASC") }).
		Where("id = ? AND organization_id = ?", id, organizationId).First(task).Error
	return task, errorx.WithStack(err)
}

func (r *ProjectManagementTaskRepo) PageByProject(ctx context.Context, organizationId int64, input dto.ProjectManagementTaskPageRequestDto) (gormx.PageResponse[*model.ProjectManagementTask], error) {
	tasksTable := model.ProjectManagementTask{}.TableName()
	assigneesTable := model.ProjectManagementTaskAssignee{}.TableName()
	stagesTable := model.ProjectManagementWorkflowStage{}.TableName()
	query := r.Client(ctx).
		Joins(fmt.Sprintf("JOIN %s ON %s.id = %s.stage_id AND %s.is_del = 0", stagesTable, stagesTable, tasksTable, stagesTable)).
		Where(fmt.Sprintf("%s.organization_id = ? AND %s.project_id = ?", tasksTable, tasksTable), organizationId, input.Filter.ProjectId).
		Order(fmt.Sprintf("%s.sort_order ASC, %s.id DESC", stagesTable, tasksTable))
	if input.Filter.StageId != 0 {
		query = query.Where(fmt.Sprintf("%s.stage_id = ?", tasksTable), input.Filter.StageId)
	}
	if input.Filter.TaskTypeId != 0 {
		query = query.Where(fmt.Sprintf("CAST(%s.task_type->>'id' AS TEXT) = CAST(? AS TEXT)", tasksTable), input.Filter.TaskTypeId)
	}
	if input.Filter.AssigneeId != 0 {
		query = query.Where(fmt.Sprintf("EXISTS (SELECT 1 FROM %s WHERE %s.task_id = %s.id AND %s.organization_user_id = ? AND %s.is_del = 0)", assigneesTable, assigneesTable, tasksTable, assigneesTable, assigneesTable), input.Filter.AssigneeId)
	}
	if input.Filter.Keyword != "" {
		query = query.Where(fmt.Sprintf("LOWER(%s.title) LIKE LOWER(?)", tasksTable), "%"+input.Filter.Keyword+"%")
	}
	if input.Filter.Overdue != nil {
		if *input.Filter.Overdue {
			query = query.Where(fmt.Sprintf("%s.end_at IS NOT NULL AND %s.end_at < CURRENT_TIMESTAMP AND %s.phase <> ?", tasksTable, tasksTable, stagesTable), "finished")
		} else {
			query = query.Where(fmt.Sprintf("(%s.end_at IS NULL OR %s.end_at >= CURRENT_TIMESTAMP OR %s.phase = ?)", tasksTable, tasksTable, stagesTable), "finished")
		}
	}
	query = query.Preload("Assignees", func(db *gorm.DB) *gorm.DB { return db.Order("id ASC") })
	return r.BaseRepo.Page(query, gormx.PageRequest{Page: input.Page, Limit: input.Limit})
}

func (r *ProjectManagementTaskRepo) PageMine(ctx context.Context, organizationId, organizationUserId int64, input dto.ProjectManagementMyTaskPageRequestDto) (gormx.PageResponse[*model.ProjectManagementTask], error) {
	tasksTable := model.ProjectManagementTask{}.TableName()
	assigneesTable := model.ProjectManagementTaskAssignee{}.TableName()
	stagesTable := model.ProjectManagementWorkflowStage{}.TableName()
	query := r.Client(ctx).
		Joins(fmt.Sprintf("JOIN %s ON %s.id = %s.stage_id AND %s.is_del = 0", stagesTable, stagesTable, tasksTable, stagesTable)).
		Where(fmt.Sprintf("%s.organization_id = ?", tasksTable), organizationId).
		Where(fmt.Sprintf("(%s.creator_id = ? OR EXISTS (SELECT 1 FROM %s mine WHERE mine.task_id = %s.id AND mine.organization_user_id = ? AND mine.is_del = 0))", tasksTable, assigneesTable, tasksTable), organizationUserId, organizationUserId).
		Order(fmt.Sprintf("%s.sort_order ASC, %s.id DESC", stagesTable, tasksTable))
	if input.Filter.ProjectId != 0 {
		query = query.Where(fmt.Sprintf("%s.project_id = ?", tasksTable), input.Filter.ProjectId)
	}
	if input.Filter.StageId != 0 {
		query = query.Where(fmt.Sprintf("%s.stage_id = ?", tasksTable), input.Filter.StageId)
	}
	if input.Filter.TaskTypeId != 0 {
		query = query.Where(fmt.Sprintf("CAST(%s.task_type->>'id' AS TEXT) = CAST(? AS TEXT)", tasksTable), input.Filter.TaskTypeId)
	}
	if input.Filter.AssigneeId != 0 {
		query = query.Where(fmt.Sprintf("EXISTS (SELECT 1 FROM %s WHERE %s.task_id = %s.id AND %s.organization_user_id = ? AND %s.is_del = 0)", assigneesTable, assigneesTable, tasksTable, assigneesTable, assigneesTable), input.Filter.AssigneeId)
	}
	if input.Filter.Keyword != "" {
		query = query.Where(fmt.Sprintf("LOWER(%s.title) LIKE LOWER(?)", tasksTable), "%"+input.Filter.Keyword+"%")
	}
	if input.Filter.Overdue != nil {
		if *input.Filter.Overdue {
			query = query.Where(fmt.Sprintf("%s.end_at IS NOT NULL AND %s.end_at < CURRENT_TIMESTAMP AND %s.phase <> ?", tasksTable, tasksTable, stagesTable), "finished")
		} else {
			query = query.Where(fmt.Sprintf("(%s.end_at IS NULL OR %s.end_at >= CURRENT_TIMESTAMP OR %s.phase = ?)", tasksTable, tasksTable, stagesTable), "finished")
		}
	}
	query = query.Preload("Assignees", func(db *gorm.DB) *gorm.DB { return db.Order("id ASC") })
	return r.BaseRepo.Page(query, gormx.PageRequest{Page: input.Page, Limit: input.Limit})
}

func (r *ProjectManagementTaskRepo) LockWithTx(ctx context.Context, tx *gorm.DB, id int64) (*model.ProjectManagementTask, error) {
	task := new(model.ProjectManagementTask)
	err := tx.WithContext(ctx).Clauses(clause.Locking{Strength: "UPDATE"}).First(task, id).Error
	return task, errorx.WithStack(err)
}

func (r *ProjectManagementTaskRepo) IsAssigneeWithTx(ctx context.Context, tx *gorm.DB, taskId, organizationUserId int64) (bool, error) {
	var count int64
	err := tx.WithContext(ctx).Model(&model.ProjectManagementTaskAssignee{}).Where("task_id = ? AND organization_user_id = ?", taskId, organizationUserId).Count(&count).Error
	return count > 0, errorx.WithStack(err)
}

func (r *ProjectManagementTaskRepo) DirectChildrenWithTx(ctx context.Context, tx *gorm.DB, taskId int64) ([]*model.ProjectManagementTask, error) {
	children := make([]*model.ProjectManagementTask, 0)
	err := tx.WithContext(ctx).Where("parent_task_id = ?", taskId).Order("id ASC").Find(&children).Error
	return children, errorx.WithStack(err)
}

func (r *ProjectManagementTaskRepo) UpdateStageWithTx(ctx context.Context, tx *gorm.DB, taskId, stageId int64) error {
	err := tx.WithContext(ctx).Model(&model.ProjectManagementTask{}).Where("id = ?", taskId).Update("stage_id", stageId).Error
	return errorx.WithStack(err)
}

func (r *ProjectManagementTaskRepo) IdsByProjectWithTx(ctx context.Context, tx *gorm.DB, projectId int64) ([]int64, error) {
	ids := make([]int64, 0)
	err := tx.WithContext(ctx).Model(&model.ProjectManagementTask{}).Where("project_id = ?", projectId).Pluck("id", &ids).Error
	return ids, errorx.WithStack(err)
}

func (r *ProjectManagementTaskRepo) DeleteByProjectWithTx(ctx context.Context, tx *gorm.DB, projectId int64, taskIds []int64) error {
	if len(taskIds) > 0 {
		if err := tx.WithContext(ctx).Where("task_id IN ?", taskIds).Delete(&model.ProjectManagementTaskAssignee{}).Error; err != nil {
			return errorx.WithStack(err)
		}
		if err := tx.WithContext(ctx).Where("successor_task_id IN ? OR predecessor_task_id IN ?", taskIds, taskIds).Delete(&model.ProjectManagementTaskDependency{}).Error; err != nil {
			return errorx.WithStack(err)
		}
	}
	err := tx.WithContext(ctx).Where("project_id = ?", projectId).Delete(&model.ProjectManagementTask{}).Error
	return errorx.WithStack(err)
}

func (r *ProjectManagementTaskRepo) CreateAssigneesWithTx(ctx context.Context, tx *gorm.DB, assignees []*model.ProjectManagementTaskAssignee) error {
	if len(assignees) == 0 {
		return nil
	}
	err := tx.WithContext(ctx).Create(&assignees).Error
	return errorx.WithStack(err)
}

func (r *ProjectManagementTaskRepo) ReplaceAssigneesWithTx(ctx context.Context, tx *gorm.DB, taskId int64, assignees []*model.ProjectManagementTaskAssignee) error {
	if err := tx.WithContext(ctx).Where("task_id = ?", taskId).Delete(&model.ProjectManagementTaskAssignee{}).Error; err != nil {
		return errorx.WithStack(err)
	}
	return r.CreateAssigneesWithTx(ctx, tx, assignees)
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
	if err := tx.WithContext(ctx).Where("successor_task_id IN ? OR predecessor_task_id IN ?", taskIds, taskIds).Delete(&model.ProjectManagementTaskDependency{}).Error; err != nil {
		return errorx.WithStack(err)
	}
	err := tx.WithContext(ctx).Where("id IN ?", taskIds).Delete(&model.ProjectManagementTask{}).Error
	return errorx.WithStack(err)
}
