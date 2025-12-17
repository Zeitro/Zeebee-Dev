package repo

import (
	"gorm.io/gorm"

	"zeebee-backend/common/gormx"
	"zeebee-backend/domains/zeebee/model"
)

type EmailAuthorizationRepo struct {
	*gormx.BaseRepo[*model.EmailAuthorization]
	m *model.EmailAuthorization
}

func NewEmailAuthorizationRepo(db *gorm.DB) *EmailAuthorizationRepo {
	return &EmailAuthorizationRepo{
		BaseRepo: gormx.NewBaseRepo[*model.EmailAuthorization](db),
		m:        new(model.EmailAuthorization),
	}
}
