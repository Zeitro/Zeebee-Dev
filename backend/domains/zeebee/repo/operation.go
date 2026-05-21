package repo

import (
	"gorm.io/gorm"

	"zeebee-backend/common/gormx"
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
