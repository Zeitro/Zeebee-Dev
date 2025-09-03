package repo

import (
	"context"

	"github.com/pkg/errors"
	"gorm.io/gorm"

	"zeebee-backend/app/borrower/rpc/internal/repo/model"
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

// UpdateOneWithTx
func (repo *FileRepo) UpdateOneWithTx(ctx context.Context, tx *gorm.DB, id int64, data *model.File) (result *model.File, err error) {
	err = tx.WithContext(ctx).Select("*").Omit("id", "created_at", "ref_id").Where("id = ?", id).Updates(data).Error
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return data, nil
}
