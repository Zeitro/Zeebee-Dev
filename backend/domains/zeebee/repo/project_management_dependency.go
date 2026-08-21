package repo

import (
	"context"
	"fmt"
	"sort"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"zeebee-backend/common/errorx"
	"zeebee-backend/common/gormx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/model"
)

type ProjectManagementDependencyRepo struct {
	*gormx.BaseRepo[*model.ProjectManagementTaskDependency]
}

func NewProjectManagementDependencyRepo(db *gorm.DB) *ProjectManagementDependencyRepo {
	return &ProjectManagementDependencyRepo{BaseRepo: gormx.NewBaseRepo[*model.ProjectManagementTaskDependency](db)}
}

func (r *ProjectManagementDependencyRepo) LockTasksWithTx(ctx context.Context, tx *gorm.DB, taskIds ...int64) ([]*model.ProjectManagementTask, error) {
	ids := append([]int64(nil), taskIds...)
	sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })
	tasks := make([]*model.ProjectManagementTask, 0, len(ids))
	for _, id := range ids {
		task := new(model.ProjectManagementTask)
		if err := tx.WithContext(ctx).Clauses(clause.Locking{Strength: "UPDATE"}).First(task, id).Error; err != nil {
			return nil, errorx.WithStack(err)
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (r *ProjectManagementDependencyRepo) FindByTasksWithTx(ctx context.Context, tx *gorm.DB, successorTaskId, predecessorTaskId int64) (*model.ProjectManagementTaskDependency, error) {
	dependency := new(model.ProjectManagementTaskDependency)
	err := tx.WithContext(ctx).Where(
		"(task_id = ? AND dependency_task_id = ? AND dependency_type = ?) OR (task_id = ? AND dependency_task_id = ? AND dependency_type = ?)",
		successorTaskId, predecessorTaskId, constants.ProjectManagementDependencyTypePreTask,
		predecessorTaskId, successorTaskId, constants.ProjectManagementDependencyTypePostTask,
	).First(dependency).Error
	return dependency, errorx.WithStack(err)
}

func (r *ProjectManagementDependencyRepo) LockWithTx(ctx context.Context, tx *gorm.DB, id int64) (*model.ProjectManagementTaskDependency, error) {
	dependency := new(model.ProjectManagementTaskDependency)
	err := tx.WithContext(ctx).Clauses(clause.Locking{Strength: "UPDATE"}).First(dependency, id).Error
	return dependency, errorx.WithStack(err)
}

func (r *ProjectManagementDependencyRepo) HasPathWithTx(ctx context.Context, tx *gorm.DB, fromTaskId, toTaskId int64, excludedDependencyId ...int64) (bool, error) {
	var found bool
	excludedId := int64(0)
	if len(excludedDependencyId) > 0 {
		excludedId = excludedDependencyId[0]
	}
	err := tx.WithContext(ctx).Raw(`WITH RECURSIVE edges(from_task_id, to_task_id) AS (
    SELECT dependency_task_id, task_id
    FROM project_management_task_dependencies
    WHERE dependency_type = ? AND is_del = 0 AND (? = 0 OR id <> ?)
    UNION ALL
    SELECT task_id, dependency_task_id
    FROM project_management_task_dependencies
    WHERE dependency_type = ? AND is_del = 0 AND (? = 0 OR id <> ?)
), reachable(id) AS (
    SELECT to_task_id FROM edges WHERE from_task_id = ?
    UNION
    SELECT edges.to_task_id FROM edges JOIN reachable ON edges.from_task_id = reachable.id
)
SELECT EXISTS (SELECT 1 FROM reachable WHERE id = ?)`, constants.ProjectManagementDependencyTypePreTask, excludedId, excludedId, constants.ProjectManagementDependencyTypePostTask, excludedId, excludedId, fromTaskId, toTaskId).Scan(&found).Error
	return found, errorx.WithStack(err)
}

func (r *ProjectManagementDependencyRepo) CreateWithTx(ctx context.Context, tx *gorm.DB, dependency *model.ProjectManagementTaskDependency) (*model.ProjectManagementTaskDependency, error) {
	return r.BaseRepo.CreateWithTx(ctx, tx, dependency, false)
}

func (r *ProjectManagementDependencyRepo) UpdateWithTx(ctx context.Context, tx *gorm.DB, dependency *model.ProjectManagementTaskDependency) error {
	return r.UpdateColumnsWithTx(ctx, tx,
		map[string]any{"task_id": dependency.TaskId, "dependency_task_id": dependency.DependencyTaskId, "dependency_type": dependency.DependencyType, "dependency_rule": dependency.DependencyRule},
		map[string]any{"id": dependency.Id}, nil)
}

func (r *ProjectManagementDependencyRepo) ListBySuccessorWithTx(ctx context.Context, tx *gorm.DB, taskId int64) ([]*model.ProjectManagementTaskDependency, error) {
	dependencies := make([]*model.ProjectManagementTaskDependency, 0)
	err := tx.WithContext(ctx).Where(
		"(task_id = ? AND dependency_type = ?) OR (dependency_task_id = ? AND dependency_type = ?)",
		taskId, constants.ProjectManagementDependencyTypePreTask,
		taskId, constants.ProjectManagementDependencyTypePostTask,
	).Order("id ASC").Find(&dependencies).Error
	return dependencies, errorx.WithStack(err)
}

func (r *ProjectManagementDependencyRepo) PageByTask(ctx context.Context, taskId int64, predecessor bool, pageReq gormx.PageRequest) (gormx.PageResponse[*model.ProjectManagementTaskDependency], error) {
	direction := constants.ProjectManagementDependencyTypePreTask
	reverseDirection := constants.ProjectManagementDependencyTypePostTask
	if !predecessor {
		direction, reverseDirection = reverseDirection, direction
	}
	table := model.ProjectManagementTaskDependency{}.TableName()
	query := r.Client(ctx).Table("(?) AS dependencies", r.Client(ctx).Raw(fmt.Sprintf(`
SELECT id, task_id, dependency_task_id, dependency_type, dependency_rule
FROM %s
WHERE task_id = ? AND dependency_type = ? AND is_del = 0
UNION ALL
SELECT id, ? AS task_id, task_id AS dependency_task_id, ? AS dependency_type, dependency_rule
FROM %s
WHERE dependency_task_id = ? AND dependency_type = ? AND is_del = 0`, table, table), taskId, direction, taskId, direction, taskId, reverseDirection)).Order("id DESC")
	return gormx.Page[*model.ProjectManagementTaskDependency](query, pageReq)
}

func (r *ProjectManagementDependencyRepo) DeleteWithTx(ctx context.Context, tx *gorm.DB, dependency *model.ProjectManagementTaskDependency) error {
	_, err := r.BaseRepo.DeleteWithTx(ctx, tx, map[string]any{"id": dependency.Id})
	return err
}
