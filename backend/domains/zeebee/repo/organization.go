package repo

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"zeebee-backend/common/errorx"
	"zeebee-backend/common/gormx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
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

// UpdateOwnerIdWithTx update owner id with transaction
func (r *OrganizationRepo) UpdateOwnerIdWithTx(ctx context.Context, tx *gorm.DB, id int64, ownerId int64) (err error) {
	err = tx.WithContext(ctx).
		Model(&model.Organization{}).
		Where("id = ?", id).
		Update("owner_id", ownerId).Error
	return errorx.WithStack(err)
}

// UpdateOwnerOrganizationStatusWithTx update owner's organization status with transaction
func (r *OrganizationRepo) UpdateOwnerOrganizationStatusWithTx(ctx context.Context, tx *gorm.DB, ownerId int64, status constants.OrganizationStatus) (err error) {
	err = tx.WithContext(ctx).
		Model(&model.Organization{}).
		Where("owner_id = ?", ownerId).
		Update("status", status).Error
	return errorx.WithStack(err)
}

// GetOrganization return organization by id and type
func (repo *OrganizationRepo) GetOrganization(ctx context.Context, id int64, organizationType constants.UserType, withPreloads bool) (result *model.Organization, err error) {
	query := repo.Client(ctx).
		Where("id = ? and type = ?", id, organizationType)
	if withPreloads {
		query = query.Preload("OrganizationEmails")
	}
	err = query.First(&result).Error
	return result, errorx.WithStack(err)
}

// GetOrganizationWithTx return organization with transaction
func (repo *OrganizationRepo) GetOrganizationWithTx(ctx context.Context, tx *gorm.DB, id int64, organizationType constants.UserType, withPreloads bool) (result *model.Organization, err error) {
	query := tx.WithContext(ctx).Model(&model.Organization{}).
		Where("id = ? and type = ?", id, organizationType)
	if withPreloads {
		query = query.Preload("OrganizationEmails")
	}
	err = query.First(&result).Error
	return result, errorx.WithStack(err)
}

// All
func (repo *OrganizationRepo) All(ctx context.Context, organizationType constants.UserType) (result []*model.Organization, err error) {
	err = repo.Client(ctx).
		Where("type = ? and status = ?", organizationType, constants.OrganizationStatusActive).
		Order("id asc").
		Find(&result).Error
	return result, errorx.WithStack(err)
}

// Page
func (repo *OrganizationRepo) Page(ctx context.Context, pageReq dto.OrganizationPageRequestDto) (page gormx.PageResponse[*model.Organization], err error) {
	filter := pageReq.Filter
	sorter := pageReq.Sorter
	query := repo.Client(ctx)
	if len(filter.Keyword) > 0 {
		keyword := "%" + filter.Keyword + "%"
		query = query.Where("name LIKE @keyword OR contact_email LIKE @keyword OR contact_name LIKE @keyword OR contact_phone LIKE @keyword", map[string]interface{}{
			"keyword": keyword,
		})
	}
	if len(filter.OrganizationType) > 0 {
		query = query.Where("type = ?", filter.OrganizationType)
	}
	if len(filter.Status) > 0 {
		query = query.Where("status = ?", filter.Status)
	}
	if len(sorter.Column) > 0 {
		query = query.Order(fmt.Sprintf("%s %s", sorter.Column, sorter.Sort))
	}
	query = query.Order("id asc")
	page, err = gormx.Page[*model.Organization](query, gormx.PageRequest{
		Page:  pageReq.Page,
		Limit: pageReq.Limit,
	})
	return page, errorx.WithStack(err)
}

// GetUserOrganizations
func (repo *OrganizationRepo) GetUserOrganizations(ctx context.Context, userId int64, userType constants.UserType) (result []*model.Organization, err error) {
	organizationUserTable := model.OrganizationUser{}.TableName()
	organizationTable := model.Organization{}.TableName()
	err = repo.Client(ctx).
		Joins(
			fmt.Sprintf("left join %s on %s.id = %s.organization_id", organizationUserTable, organizationTable, organizationUserTable),
		).
		Where(
			fmt.Sprintf("%s.user_id = ? and %s.type = ? and %s.is_del = 0", organizationUserTable, organizationTable, organizationUserTable),
			userId, userType).
		Find(&result).Error
	return result, errorx.WithStack(err)
}

// GetUserOrganizationsWithTx
func (repo *OrganizationRepo) GetUserOrganizationsWithTx(ctx context.Context, tx *gorm.DB, userId int64, userType constants.UserType) (result []*model.Organization, err error) {
	organizationUserTable := model.OrganizationUser{}.TableName()
	organizationTable := model.Organization{}.TableName()
	err = tx.WithContext(ctx).Model(&model.Organization{}).
		Joins(
			fmt.Sprintf("left join %s on %s.id = %s.organization_id", organizationUserTable, organizationTable, organizationUserTable),
		).
		Where(
			fmt.Sprintf("%s.user_id = ? and %s.type = ? and %s.is_del = 0", organizationUserTable, organizationTable, organizationUserTable),
			userId, userType).
		Find(&result).Error
	return result, errorx.WithStack(err)
}
