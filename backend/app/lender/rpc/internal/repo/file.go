package repo

import (
	"context"

	"github.com/pkg/errors"
	"gorm.io/gorm"

	"zeebee-backend/app/lender/rpc/internal/repo/model"
	"zeebee-backend/common/constants"
	"zeebee-backend/common/gormx"
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

// GetFileCountByEntityIdWithTx
func (r *FileRepo) GetFileCountByEntityIdWithTx(ctx context.Context, tx *gorm.DB, entityId int64, entityType constants.FileEntityType) (count int64, err error) {
	err = tx.WithContext(ctx).Model(&model.File{}).
		Where("entity_id = ?", entityId).
		Where("entity_type = ?", entityType).
		Count(&count).Error

	return count, errors.WithStack(err)
}
