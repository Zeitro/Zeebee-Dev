package repo

import (
	"context"

	"gorm.io/gorm"

	"zeebee-backend/common/errorx"
	"zeebee-backend/common/gormx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/model"
)

type DocumentTemplateRepo struct {
	*gormx.BaseRepo[*model.DocumentTemplate]
}

func NewDocumentTemplateRepo(db *gorm.DB) *DocumentTemplateRepo {
	return &DocumentTemplateRepo{
		BaseRepo: gormx.NewBaseRepo[*model.DocumentTemplate](db),
	}
}

func (repo *DocumentTemplateRepo) GetByType(ctx context.Context, templateType constants.DocumentTemplateType) (result *model.DocumentTemplate, err error) {
	err = repo.Client(ctx).
		Where("type = ?", templateType).
		First(&result).Error
	return result, errorx.WithStack(err)
}
