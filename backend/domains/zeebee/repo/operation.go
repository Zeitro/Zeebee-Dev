package repo

import (
	"context"

	"gorm.io/gorm"

	"zeebee-backend/common/gormx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/model"
)

type OperationRepo struct {
	*gormx.BaseRepo[*model.Operation]
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

func (repo *OperationRepo) GetOperationStats(ctx context.Context, entityId int64, entityType constants.OperationEntityType) (items []OperationStatsItem, total int64, err error) {
	err = repo.Client(ctx).
		Select("operation as type, count(*) as count").
		Where("entity_id = ? AND entity_type = ? AND is_del = 0", entityId, entityType).
		Group("operation").
		Find(&items).Error
	if err != nil {
		return
	}
	for _, item := range items {
		total += item.Count
	}
	return
}

func (repo *OperationRepo) GetOperationList(ctx context.Context, entityId int64, entityType constants.OperationEntityType, operationType string) (list []*model.Operation, err error) {
	query := repo.Client(ctx).
		Where("entity_id = ? AND entity_type = ? AND is_del = 0", entityId, entityType)
	if operationType != "" {
		query = query.Where("operation = ?", operationType)
	}
	err = query.Order("created_at DESC").Find(&list).Error
	return
}
