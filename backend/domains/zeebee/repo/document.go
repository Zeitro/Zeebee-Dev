package repo

import (
	"context"

	"gorm.io/gorm"

	"zeebee-backend/common/errorx"
	"zeebee-backend/common/gormx"
	"zeebee-backend/domains/zeebee/model"
)

type DocumentRepo struct {
	*gormx.BaseRepo[*model.Document]
}

func NewDocumentRepo(db *gorm.DB) *DocumentRepo {
	return &DocumentRepo{
		BaseRepo: gormx.NewBaseRepo[*model.Document](db),
	}
}

func (repo *DocumentRepo) GetChildren(ctx context.Context, parentId int64) (result []*model.Document, err error) {
	err = repo.Client(ctx).
		Where("parent_id = ?", parentId).
		Order("sort_order asc").
		Find(&result).Error
	return result, errorx.WithStack(err)
}

func (repo *DocumentRepo) GetDescendants(ctx context.Context, rootIds []int64) (result []*model.Document, err error) {
	if len(rootIds) == 0 {
		return nil, nil
	}
	err = repo.Client(ctx).
		Where("id in ? or parent_id in ?", rootIds, rootIds).
		Order("level asc, sort_order asc").
		Find(&result).Error
	return result, errorx.WithStack(err)
}

func (repo *DocumentRepo) GetAllByEntityDocumentId(ctx context.Context, entityDocumentId int64) (result []*model.Document, err error) {
	err = repo.Client(ctx).
		Where("entity_document_id = ?", entityDocumentId).
		Order("level asc, sort_order asc").
		Find(&result).Error
	return result, errorx.WithStack(err)
}

func (repo *DocumentRepo) GetAllByEntityDocumentIdWithTx(ctx context.Context, tx *gorm.DB, entityDocumentId int64) (result []*model.Document, err error) {
	err = tx.WithContext(ctx).
		Where("entity_document_id = ?", entityDocumentId).
		Order("level asc, sort_order asc").
		Find(&result).Error
	return result, errorx.WithStack(err)
}
