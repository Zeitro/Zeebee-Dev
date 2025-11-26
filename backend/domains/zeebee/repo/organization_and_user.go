package repo

import (
	"context"

	"gorm.io/gorm"

	"zeebee-backend/common/errorx"
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

func (repo *OrganizationAndUserRepo) GetRelation(ctx context.Context, organizationId int64, userId int64) (res *model.OrganizationAndUser, err error) {
	err = repo.Client(ctx).
		Where("organization_id = ? and user_id = ?", organizationId, userId).
		Find(&res).Error
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	return res, nil
}
