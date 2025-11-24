package repo

import (
	"gorm.io/gorm"

	"zeebee-backend/common/gormx"
	"zeebee-backend/domains/zeebee/model"
)

type OrganizationRepo struct {
	*gormx.BaseRepo[*model.Organization]
	m *model.Organization
}

func NewOrganizationRepo(db *gorm.DB) *OrganizationRepo {
	return &OrganizationRepo{
		BaseRepo: gormx.NewBaseRepo[*model.Organization](db),
		m:        new(model.Organization),
	}
}
