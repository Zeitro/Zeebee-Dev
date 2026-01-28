package repo

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"

	"zeebee-backend/common/errorx"
	"zeebee-backend/common/gormx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
)

type OrganizationPageEntity struct {
	Id     int64                        `json:"id"`
	Name   string                       `json:"name"`
	Status constants.OrganizationStatus `json:"status"`
	// Address        *model.Address                  `json:"address"`
	ContactEmail   string                          `json:"contact_email"`
	ContactPhone   string                          `json:"contact_phone"`
	ContactName    string                          `json:"contact_name"`
	OwnerType      constants.OrganizationOwnerType `json:"owner_type"`
	Type           constants.UserType              `json:"type"`
	Logo           string                          `json:"logo"`
	Website        string                          `json:"website"`
	Tagline        string                          `json:"tagline"`
	OwnerEmail     string                          `json:"owner_email"`
	OwnerFullName  string                          `json:"owner_full_name"`
	OwnerFirstName string                          `json:"owner_first_name"`
	OwnerMidName   string                          `json:"owner_mid_name"`
	OwnerLastName  string                          `json:"owner_last_name"`
	CreatedAt      time.Time                       `json:"created_at"`
	UpdatedAt      time.Time                       `json:"updated_at"`
}

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
		query = query.Preload("OrganizationEmails").
			Preload("Owner")
	}
	err = query.First(&result).Error
	return result, errorx.WithStack(err)
}

func (repo *OrganizationRepo) GetOrganizationById(ctx context.Context, id int64, withPreloads bool) (result *model.Organization, err error) {
	query := repo.Client(ctx).
		Where("id = ? ", id)
	if withPreloads {
		query = query.Preload("OrganizationEmails").
			Preload("Owner").
			Preload("Roles")
	}
	err = query.First(&result).Error
	return result, errorx.WithStack(err)
}

func (repo *OrganizationRepo) GetOrganizationByIdWithTx(ctx context.Context, tx *gorm.DB, id int64, withPreloads bool) (result *model.Organization, err error) {
	query := tx.WithContext(ctx).Model(&model.Organization{}).
		Where("id = ? ", id)
	if withPreloads {
		query = query.Preload("OrganizationEmails").
			Preload("Owner").
			Preload("Roles")
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
func (repo *OrganizationRepo) Page(ctx context.Context, pageReq dto.OrganizationPageRequestDto) (page gormx.PageResponse[*OrganizationPageEntity], err error) {
	filter := pageReq.Filter
	sorter := pageReq.Sorter
	organizationTable := model.Organization{}.TableName()
	userTable := model.User{}.TableName()
	fields := []string{
		fmt.Sprintf("%s.id", organizationTable),
		fmt.Sprintf("%s.name", organizationTable),
		fmt.Sprintf("%s.status", organizationTable),
		fmt.Sprintf("%s.contact_email", organizationTable),
		fmt.Sprintf("%s.contact_phone", organizationTable),
		fmt.Sprintf("%s.contact_name", organizationTable),
		fmt.Sprintf("%s.owner_type", organizationTable),
		fmt.Sprintf("%s.type", organizationTable),
		fmt.Sprintf("%s.logo", organizationTable),
		fmt.Sprintf("%s.website", organizationTable),
		fmt.Sprintf("%s.tagline", organizationTable),
		fmt.Sprintf("%s.owner_id", organizationTable),
		fmt.Sprintf("%s.created_at", organizationTable),
		fmt.Sprintf("%s.updated_at", organizationTable),
		fmt.Sprintf("%s.email as owner_email", userTable),
		fmt.Sprintf("%s.first_name as owner_first_name", userTable),
		fmt.Sprintf("%s.mid_name as owner_mid_name", userTable),
		fmt.Sprintf("%s.last_name as owner_last_name", userTable),
		fmt.Sprintf("%s.full_name as owner_full_name", userTable),
	}

	query := repo.Client(ctx).
		Model(&model.Organization{}).
		Joins(fmt.Sprintf("left join %s on %s.owner_id = %s.id", userTable, organizationTable, userTable))
	if len(filter.Keyword) > 0 {
		keyword := "%" + filter.Keyword + "%"
		query = query.Where(fmt.Sprintf(
			"%s.name LIKE @keyword OR %s.contact_email LIKE @keyword OR %s.contact_name LIKE @keyword OR %s.contact_phone LIKE @keyword or %s.email LIKE @keyword or %s.full_name LIKE @keyword",
			organizationTable,
			organizationTable,
			organizationTable,
			organizationTable,
			userTable,
			userTable,
		), map[string]interface{}{
			"keyword": keyword,
		})
	}
	if len(filter.OrganizationType) > 0 {
		query = query.Where(fmt.Sprintf("%s.type = ?", organizationTable), filter.OrganizationType)
	}
	if len(filter.Status) > 0 {
		query = query.Where(fmt.Sprintf("%s.status = ?", organizationTable), filter.Status)
	}
	if len(sorter.Column) > 0 {
		query = query.Order(fmt.Sprintf("%s %s", sorter.Column, sorter.Sort))
	}
	query = query.Order(fmt.Sprintf("%s.created_at desc", organizationTable))
	page, err = gormx.Page[*OrganizationPageEntity](query, gormx.PageRequest{
		Page:         pageReq.Page,
		Limit:        pageReq.Limit,
		SelectFields: fields,
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

// DeleteByIdWithTx deletes an organization by its ID within a database transaction.
func (repo *OrganizationRepo) DeleteWithTx(ctx context.Context, tx *gorm.DB, organization *model.Organization) (err error) {
	if organization == nil {
		return errorx.New("organization is nil")
	}
	err = tx.WithContext(ctx).Select("Roles", "OrganizationEmails").Delete(organization).Error

	return errorx.WithStack(err)
}

// GetIdsWithTx
func (repo *OrganizationRepo) GetIdsWithTx(ctx context.Context, tx *gorm.DB, organizationType constants.UserType) (result []int64, err error) {
	err = tx.WithContext(ctx).Model(&model.Organization{}).
		Select("id").
		Where("type = ?", organizationType).
		Pluck("id", &result).
		Error
	return result, errorx.WithStack(err)
}

func (repo *OrganizationRepo) GetActiveIdsWithTx(ctx context.Context, tx *gorm.DB, organizationType constants.UserType) (result []int64, err error) {
	err = tx.WithContext(ctx).Model(&model.Organization{}).
		Select("id").
		Where("type = ? and status = ?", organizationType, constants.OrganizationStatusActive).
		Pluck("id", &result).
		Error
	return result, errorx.WithStack(err)
}

func (repo *OrganizationRepo) GetActiveOrganizationsWithTx(ctx context.Context, tx *gorm.DB, organizationType constants.UserType) (results []*model.Organization, err error) {
	err = tx.WithContext(ctx).Model(&model.Organization{}).
		Where("type = ? and status = ?", organizationType, constants.OrganizationStatusActive).
		Find(&results).
		Error
	return results, errorx.WithStack(err)
}

// GetIdsByOwnerIdWithTx
func (repo *OrganizationRepo) GetIdsByOwnerIdWithTx(ctx context.Context, tx *gorm.DB, ownerId int64, organizationType constants.UserType) (result []int64, err error) {
	err = tx.WithContext(ctx).Model(&model.Organization{}).
		Select("id").
		Where("owner_id = ? and type = ?", ownerId, organizationType).
		Pluck("id", &result).Error
	return result, errorx.WithStack(err)
}
