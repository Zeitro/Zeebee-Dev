package repo

import (
	"context"

	"gorm.io/gorm"

	"zeebee-backend/common/errorx"

	"zeebee-backend/common/constants"
	"zeebee-backend/common/gormx"
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

func (r *FileRepo) GetFileCountByEntityIdWithTx(ctx context.Context, tx *gorm.DB, entityId int64, entityType constants.FileEntityType) (count int64, err error) {
	err = tx.WithContext(ctx).Model(&model.File{}).
		Where("entity_id = ?", entityId).
		Where("entity_type = ?", entityType).
		Count(&count).Error

	return count, errorx.WithStack(err)
}
