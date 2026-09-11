package repo

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

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
	Source         string                          `json:"source"`
	Order          int                             `json:"order"`
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

// UpdateOrder updates the display order of an organization and cascades the
// existing occupants so that the target org ends up at the requested slot.
//
// Behavior (scoped to the same organization type):
//   - If newOrder == oldOrder: no-op.
//   - If newOrder <  oldOrder (moving up): every other org with order in
//     [newOrder, oldOrder) shifts down by +1 to make room.
//   - If newOrder >  oldOrder (moving down): every other org with order in
//     (oldOrder, newOrder] shifts up by -1 to close the gap.
//
// All writes happen in a single transaction.
func (r *OrganizationRepo) UpdateOrder(ctx context.Context, id int64, newOrder int) (err error) {
	err = r.Transaction(ctx, func(tx *gorm.DB) error {
		var current model.Organization
		if innerErr := tx.WithContext(ctx).
			Select("id", "type", `"order"`).
			Where("id = ?", id).
			First(&current).Error; innerErr != nil {
			return innerErr
		}

		oldOrder := current.Order
		if oldOrder == newOrder {
			return nil
		}

		if newOrder < oldOrder {
			if innerErr := tx.WithContext(ctx).
				Model(&model.Organization{}).
				Where(`type = ? AND id <> ? AND "order" >= ? AND "order" < ?`,
					current.Type, id, newOrder, oldOrder).
				UpdateColumn(`"order"`, gorm.Expr(`"order" + 1`)).Error; innerErr != nil {
				return innerErr
			}
		} else {
			if innerErr := tx.WithContext(ctx).
				Model(&model.Organization{}).
				Where(`type = ? AND id <> ? AND "order" > ? AND "order" <= ?`,
					current.Type, id, oldOrder, newOrder).
				UpdateColumn(`"order"`, gorm.Expr(`"order" - 1`)).Error; innerErr != nil {
				return innerErr
			}
		}

		return tx.WithContext(ctx).
			Model(&model.Organization{}).
			Where("id = ?", id).
			UpdateColumn(`"order"`, newOrder).Error
	})
	return errorx.WithStack(err)
}

// MaxOrderByType is the non-transactional variant of MaxOrderByTypeWithTx.
func (r *OrganizationRepo) MaxOrderByType(ctx context.Context, organizationType constants.UserType) (int, error) {
	return r.MaxOrderByTypeWithTx(ctx, r.Client(ctx), organizationType)
}

// MaxOrderByTypeWithTx returns the largest "order" value among organizations
// of the given type. Returns -1 when no rows exist so callers can safely use
// (max + 1) as the next append slot (starting at 0).
func (r *OrganizationRepo) MaxOrderByTypeWithTx(ctx context.Context, tx *gorm.DB, organizationType constants.UserType) (int, error) {
	var max *int
	err := tx.WithContext(ctx).
		Model(&model.Organization{}).
		Where("type = ?", organizationType).
		Select(`COALESCE(MAX("order"), -1)`).
		Scan(&max).Error
	if err != nil {
		return 0, errorx.WithStack(err)
	}
	if max == nil {
		return -1, nil
	}
	return *max, nil
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

// FindByOwnerIdAndStatusWithTx finds organizations by owner id and status
func (r *OrganizationRepo) FindByOwnerIdAndStatusWithTx(ctx context.Context, tx *gorm.DB, ownerId int64, orgType constants.UserType, status constants.OrganizationStatus) (result []*model.Organization, err error) {
	err = tx.WithContext(ctx).
		Where("owner_id = ? AND type = ? AND status = ?", ownerId, orgType, status).
		Find(&result).Error
	return result, errorx.WithStack(err)
}

// GetOrganization return organization by id and type
func (repo *OrganizationRepo) GetOrganization(ctx context.Context, id int64, organizationType constants.UserType, withPreloads bool) (result *model.Organization, err error) {
	query := repo.Client(ctx).
		Where("id = ? and type = ?", id, organizationType)
	if withPreloads {
		query = query.Preload("OrganizationEmails").
			Preload("Owner").
			Preload("OrganizationSettings").
			Preload("Roles")
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

func (repo *OrganizationRepo) LockByIdWithTx(ctx context.Context, tx *gorm.DB, id int64) (result *model.Organization, err error) {
	result = new(model.Organization)
	err = tx.WithContext(ctx).Clauses(clause.Locking{Strength: "UPDATE"}).First(result, id).Error
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

// All return all organizations
func (repo *OrganizationRepo) All(ctx context.Context, organizationType constants.UserType) (result []*model.Organization, err error) {
	err = repo.Client(ctx).
		Where("type = ? and status = ?", organizationType, constants.OrganizationStatusActive).
		Order("id asc").
		Find(&result).Error
	return result, errorx.WithStack(err)
}

// AllByTypeAllStatuses returns all organizations of a type regardless of status (e.g. discover dedup).
func (repo *OrganizationRepo) AllByTypeAllStatuses(ctx context.Context, organizationType constants.UserType) (result []*model.Organization, err error) {
	err = repo.Client(ctx).
		Where("type = ?", organizationType).
		Order("id asc").
		Find(&result).Error
	return result, errorx.WithStack(err)
}

// AllForOption all for option
func (repo *OrganizationRepo) AllForOption(ctx context.Context, organizationType constants.UserType, status constants.OrganizationStatus) (result []*model.Organization, err error) {
	query := repo.Client(ctx).
		Select("id", "name", "contact_email", "type").
		Where("type = ?", organizationType)
	if status != constants.OrganizationStatusNone {
		query = query.Where("status = ?", status)
	}
	err = query.Order(`"order" asc, id asc`).
		Find(&result).Error
	return result, errorx.WithStack(err)
}

// Page retrieve organizations for pagination
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
		fmt.Sprintf("%s.source", organizationTable),
		fmt.Sprintf(`%s."order"`, organizationTable),
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
	query = query.Order(fmt.Sprintf(`%s."order" asc, %s.created_at desc`, organizationTable, organizationTable))
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

func (repo *OrganizationRepo) GetOrganizationsByOwnerIdWithTx(ctx context.Context, tx *gorm.DB, ownerId int64, organizationType constants.UserType) (results []*model.Organization, err error) {
	err = tx.WithContext(ctx).Model(&model.Organization{}).
		Where("owner_id = ? and type = ?", ownerId, organizationType).
		Find(&results).Error
	return results, errorx.WithStack(err)
}

// SelectOptions retrieves organization options based on the specified type.
// It returns a slice of Organization pointers containing only the id, name, and type fields.
func (repo *OrganizationRepo) SelectOptions(ctx context.Context, organizationType constants.UserType) (result []*model.Organization, err error) {
	err = repo.Client(ctx).Model(&model.Organization{}).
		Select("id", "name", "type", "contact_email").
		Where("type = ? and status = ?", organizationType, constants.OrganizationStatusActive).
		Order("id asc").
		Find(&result).Error
	return result, errorx.WithStack(err)
}
