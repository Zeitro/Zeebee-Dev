package repo

import (
	"gorm.io/gorm"

	"zeebee-backend/app/lender/rpc/internal/repo/model"
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
