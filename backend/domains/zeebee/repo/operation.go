package repo

import (
	"context"
	"fmt"

	"github.com/duke-git/lancet/v2/datetime"
	"gorm.io/gorm"

	"zeebee-backend/common/gormx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
)

type OperationRepo struct {
	*gormx.BaseRepo[*model.Operation]
}

func (repo *OperationRepo) PageProjectManagementOperations(ctx context.Context, organizationId int64, input dto.ProjectManagementOperationPageRequestDto) (gormx.PageResponse[*model.Operation], error) {
	query := repo.Client(ctx).Where("entity_id = ? AND entity_type = ? AND organization_id = ?", input.Filter.ProjectId, constants.OperationEntityTypeProjectManagementProject, organizationId).Order("created_at DESC, id DESC")
	if input.Filter.OperationType != "" {
		query = query.Where("operation = ?", input.Filter.OperationType)
	}
	if input.Filter.TaskId != 0 {
		query = query.Where(fmt.Sprintf("CAST(metadata->>'task_id' AS TEXT) = CAST(? AS TEXT)"), input.Filter.TaskId)
	}
	return repo.BaseRepo.Page(query, gormx.PageRequest{Page: input.Page, Limit: input.Limit})
}

func NewOperationRepo(db *gorm.DB) *OperationRepo {
	return &OperationRepo{
		BaseRepo: gormx.NewBaseRepo[*model.Operation](db),
	}
}

type OperationStatsItem struct {
	Type  string `json:"type"`
	Count int64  `json:"count"`
}

func (repo *OperationRepo) GetOperationStats(ctx context.Context, input dto.OperationStatsRequestDto) (items []OperationStatsItem, total int64, err error) {
	query := repo.Client(ctx).
		Select("operation as type, count(*) as count").
		Where("entity_id = ? AND entity_type = ? AND is_del = 0", input.EntityId, input.EntityType)
	if input.StartDate != nil && input.StartDate.Valid {
		startDate := datetime.BeginOfDay(input.StartDate.Time)
		query = query.Where("created_at >= ?", startDate)
	}
	if input.EndDate != nil && input.EndDate.Valid {
		endDate := datetime.EndOfDay(input.EndDate.Time)
		query = query.Where("created_at <= ?", endDate)
	}

	err = query.Group("operation").
		Find(&items).Error
	if err != nil {
		return
	}
	for _, item := range items {
		total += item.Count
	}
	return
}

func (repo *OperationRepo) GetOperationList(ctx context.Context, input dto.OperationListRequestDto) (list []*model.Operation, err error) {
	query := repo.Client(ctx).
		Where("entity_id = ? AND entity_type = ?", input.EntityId, input.EntityType)
	if len(input.Operation) > 0 {
		query = query.Where("operation = ?", input.Operation)
	}

	if input.StartDate != nil && input.StartDate.Valid {
		startDate := datetime.BeginOfDay(input.StartDate.Time)
		query = query.Where("created_at >= ?", startDate)
	}
	if input.EndDate != nil && input.EndDate.Valid {
		endDate := datetime.EndOfDay(input.EndDate.Time)
		query = query.Where("created_at <= ?", endDate)
	}

	err = query.Order("created_at DESC").Find(&list).Error
	return
}
