package repo

import (
	"context"

	"gorm.io/gorm"

	"zeebee-backend/common/errorx"
	"zeebee-backend/common/gormx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/model"
)

type EntityDocumentRepo struct {
	*gormx.BaseRepo[*model.EntityDocument]
}

func NewEntityDocumentRepo(db *gorm.DB) *EntityDocumentRepo {
	return &EntityDocumentRepo{
		BaseRepo: gormx.NewBaseRepo[*model.EntityDocument](db),
	}
}

func (repo *EntityDocumentRepo) GetByEntity(ctx context.Context, entityType constants.EntityDocumentType, entityId int64) (result *model.EntityDocument, err error) {
	err = repo.Client(ctx).
		Where("entity_type = ? and entity_id = ?", entityType, entityId).
		Preload("RootDocument").
		First(&result).Error
	return result, errorx.WithStack(err)
}

func (repo *EntityDocumentRepo) GetByEntityWithTx(ctx context.Context, tx *gorm.DB, entityType constants.EntityDocumentType, entityId int64) (result *model.EntityDocument, err error) {

	err = tx.WithContext(ctx).
		Model(model.EntityDocument{}).
		Where("entity_type = ? and entity_id = ?", entityType, entityId).
		First(&result).Error
	return result, errorx.WithStack(err)
}
