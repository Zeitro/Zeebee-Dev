package repo

import (
	"context"
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
	err := tx.WithContext(ctx).Where("successor_task_id = ? AND predecessor_task_id = ?", successorTaskId, predecessorTaskId).First(dependency).Error
	return dependency, errorx.WithStack(err)
}

func (r *ProjectManagementDependencyRepo) LockWithTx(ctx context.Context, tx *gorm.DB, id int64) (*model.ProjectManagementTaskDependency, error) {
	dependency := new(model.ProjectManagementTaskDependency)
	err := tx.WithContext(ctx).Clauses(clause.Locking{Strength: "UPDATE"}).First(dependency, id).Error
	return dependency, errorx.WithStack(err)
}

func (r *ProjectManagementDependencyRepo) HasPathWithTx(ctx context.Context, tx *gorm.DB, fromTaskId, toTaskId int64) (bool, error) {
	var found bool
	err := tx.WithContext(ctx).Raw(`WITH RECURSIVE reachable(id) AS (
    SELECT predecessor_task_id FROM project_management_task_dependencies WHERE successor_task_id = ? AND is_del = 0
    UNION
    SELECT dependency.predecessor_task_id FROM project_management_task_dependencies dependency JOIN reachable ON dependency.successor_task_id = reachable.id WHERE dependency.is_del = 0
) SELECT EXISTS (SELECT 1 FROM reachable WHERE id = ?)`, fromTaskId, toTaskId).Scan(&found).Error
	return found, errorx.WithStack(err)
}

func (r *ProjectManagementDependencyRepo) CreateWithTx(ctx context.Context, tx *gorm.DB, dependency *model.ProjectManagementTaskDependency) (*model.ProjectManagementTaskDependency, error) {
	if err := tx.WithContext(ctx).Create(dependency).Error; err != nil {
		return nil, errorx.WithStack(err)
	}
	return dependency, nil
}

func (r *ProjectManagementDependencyRepo) UpdateTypeWithTx(ctx context.Context, tx *gorm.DB, dependency *model.ProjectManagementTaskDependency, dependencyType constants.ProjectManagementDependencyType) error {
	err := tx.WithContext(ctx).Model(dependency).Update("dependency_type", dependencyType).Error
	return errorx.WithStack(err)
}

func (r *ProjectManagementDependencyRepo) ListBySuccessorWithTx(ctx context.Context, tx *gorm.DB, taskId int64) ([]*model.ProjectManagementTaskDependency, error) {
	dependencies := make([]*model.ProjectManagementTaskDependency, 0)
	err := tx.WithContext(ctx).Where("successor_task_id = ?", taskId).Order("id ASC").Find(&dependencies).Error
	return dependencies, errorx.WithStack(err)
}

func (r *ProjectManagementDependencyRepo) PageByTask(ctx context.Context, taskId int64, predecessor bool, pageReq gormx.PageRequest) (gormx.PageResponse[*model.ProjectManagementTaskDependency], error) {
	column := "predecessor_task_id"
	if predecessor {
		column = "successor_task_id"
	}
	query := r.Client(ctx).Where(column+" = ?", taskId).Order("id DESC")
	return r.BaseRepo.Page(query, pageReq)
}

func (r *ProjectManagementDependencyRepo) DeleteWithTx(ctx context.Context, tx *gorm.DB, dependency *model.ProjectManagementTaskDependency) error {
	err := tx.WithContext(ctx).Delete(dependency).Error
	return errorx.WithStack(err)
}
