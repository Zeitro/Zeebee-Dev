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

// GetMaxSortOrderByParentId get max sort order by parent id
func (repo *DocumentRepo) GetMaxSortOrderByParentId(ctx context.Context, parentId int64) (result int64, err error) {
	err = repo.Client(ctx).
		Where("parent_id = ?", parentId).
		Select("COALESCE(MAX(sort_order), 0) as max_sort_order").
		Scan(&result).Error
	return result, errorx.WithStack(err)
}

// FindOneWithFiles get document with files preloaded, ordered by created_at DESC
func (repo *DocumentRepo) FindOneWithFiles(ctx context.Context, id int64) (result *model.Document, err error) {
	err = repo.Client(ctx).
		Preload("Files", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at DESC")
		}).
		First(&result, id).Error
	return result, errorx.WithStack(err)
}
