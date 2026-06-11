package repo

import (
	"context"
	"time"

	"gorm.io/gorm"

	"zeebee-backend/common/errorx"

	"zeebee-backend/common/gormx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/model"
)

type FileRepo struct {
	*gormx.BaseRepo[*model.File]
	m *model.File
}

func NewFileRepo(db *gorm.DB) *FileRepo {
	return &FileRepo{
		BaseRepo: gormx.NewBaseRepo[*model.File](db),
		m:        new(model.File),
	}
}

func (repo *FileRepo) GetFileCountByEntityIdWithTx(ctx context.Context, tx *gorm.DB, entityId int64, entityType constants.FileEntityType) (count int64, err error) {
	err = tx.WithContext(ctx).Model(&model.File{}).
		Where("entity_id = ?", entityId).
		Where("entity_type = ?", entityType).
		Count(&count).Error

	return count, errorx.WithStack(err)
}

// BatchFileCountByEntityIds returns entity_id -> file count map for given entity IDs
func (repo *FileRepo) BatchFileCountByEntityIds(ctx context.Context, tx *gorm.DB, entityIds []int64, entityType constants.FileEntityType) (result map[int64]int32, err error) {
	result = make(map[int64]int32)
	if len(entityIds) == 0 {
		return result, nil
	}

	type row struct {
		EntityId int64 `gorm:"column:entity_id"`
		Count    int32 `gorm:"column:cnt"`
	}
	var rows []row

	err = tx.WithContext(ctx).Model(&model.File{}).
		Select("entity_id, count(*) as cnt").
		Where("entity_id IN ?", entityIds).
		Where("entity_type = ?", entityType).
		Group("entity_id").
		Find(&rows).Error
	if err != nil {
		return result, errorx.WithStack(err)
	}

	for _, r := range rows {
		result[r.EntityId] = r.Count
	}
	return
}

// GetFiles returns files by entity ID
func (repo *FileRepo) GetFiles(ctx context.Context, entityId int64, entityType constants.FileEntityType) (result []*model.File, err error) {
	err = repo.Client(ctx).Where("entity_id = ?", entityId).Where("entity_type = ?", entityType).Find(&result).Error
	return result, errorx.WithStack(err)
}

// MarkDriveSynced records the Google Drive link and sync time for a file.
func (repo *FileRepo) MarkDriveSynced(ctx context.Context, id int64, driveUrl string, syncedAt time.Time) error {
	return repo.UpdateColumnsById(ctx, id, map[string]interface{}{
		"drive_url":       driveUrl,
		"drive_synced_at": syncedAt,
	})
}

// DeleteByEntityIdsWithTx deletes files by entity IDs with transaction
func (repo *FileRepo) DeleteByEntityIdWithTx(ctx context.Context, tx *gorm.DB, entityIds []int64, entityType constants.FileEntityType) (err error) {
	if len(entityIds) == 0 {
		return
	}

	err = tx.WithContext(ctx).
		Where("entity_id IN ? and entity_type = ?", entityIds, entityType).
		Delete(&model.File{}).Error
	return errorx.WithStack(err)
}
