package repo

import (
	"gorm.io/gorm"

	"zeebee-backend/common/gormx"
	"zeebee-backend/domains/zeebee/model"
)

type OrganizationAndUserRepo struct {
	*gormx.BaseRepo[*model.OrganizationAndUser]
	m *model.OrganizationAndUser
}

func NewOrganizationAndUserRepo(db *gorm.DB) *OrganizationAndUserRepo {
	return &OrganizationAndUserRepo{
		BaseRepo: gormx.NewBaseRepo[*model.OrganizationAndUser](db),
		m:        new(model.OrganizationAndUser),
	}
}
