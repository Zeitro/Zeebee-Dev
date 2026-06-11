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

type OrganizationSubscriptionPageRow struct {
	Id                  int64     `gorm:"column:id"`
	Code                string    `gorm:"column:code"`
	OrganizationId      int64     `gorm:"column:organization_id"`
	OrganizationName    string    `gorm:"column:organization_name"`
	PricingPlanId       int64     `gorm:"column:pricing_plan_id"`
	PricingPlanName     string    `gorm:"column:pricing_plan_name"`
	PricingPlanPolicyId int64     `gorm:"column:pricing_plan_policy_id"`
	Status              string    `gorm:"column:status"`
	ExpireDate          time.Time `gorm:"column:expire_date"`
	OriginalAmount      float64   `gorm:"column:original_amount"`
	ActualAmount        float64   `gorm:"column:actual_amount"`
	IsAutoRenew         bool      `gorm:"column:is_auto_renew"`
	CreatedAt           time.Time `gorm:"column:created_at"`
	UpdatedAt           time.Time `gorm:"column:updated_at"`
}

type OrganizationSubscriptionRepo struct {
	*gormx.BaseRepo[*model.OrganizationSubscription]
}

func NewOrganizationSubscriptionRepo(db *gorm.DB) *OrganizationSubscriptionRepo {
	return &OrganizationSubscriptionRepo{
		BaseRepo: gormx.NewBaseRepo[*model.OrganizationSubscription](db),
	}
}

func (repo *OrganizationSubscriptionRepo) GetOne(ctx context.Context, id int64) (result *model.OrganizationSubscription, err error) {
	err = repo.Client(ctx).Where("id = ?", id).First(&result).Error
	return result, errorx.WithStack(err)
}

func (repo *OrganizationSubscriptionRepo) ExistsByCode(ctx context.Context, code string) (bool, error) {
	return repo.IsExistedByColumns(ctx, map[string]interface{}{
		"code": code,
	})
}

func (repo *OrganizationSubscriptionRepo) FindActiveByOrgId(ctx context.Context, organizationId int64) (result []*model.OrganizationSubscription, err error) {
	err = repo.Client(ctx).
		Where("organization_id = ? AND status NOT IN ?",
			organizationId,
			[]constants.OrganizationSubscriptionStatus{constants.OrganizationSubscriptionStatusExpired},
		).
		Order("id desc").
		Find(&result).Error
	return result, errorx.WithStack(err)
}

func (repo *OrganizationSubscriptionRepo) FindActiveByOrgIdWithTx(ctx context.Context, tx *gorm.DB, organizationId int64) (result []*model.OrganizationSubscription, err error) {
	err = tx.WithContext(ctx).
		Where("organization_id = ? AND status NOT IN ?",
			organizationId,
			[]constants.OrganizationSubscriptionStatus{constants.OrganizationSubscriptionStatusExpired},
		).
		Order("id desc").
		Find(&result).Error
	return result, errorx.WithStack(err)
}

func (repo *OrganizationSubscriptionRepo) FindByPricingPlanId(ctx context.Context, pricingPlanId int64) (result []*model.OrganizationSubscription, err error) {
	err = repo.Client(ctx).
		Where("pricing_plan_id = ?", pricingPlanId).
		Find(&result).Error
	return result, errorx.WithStack(err)
}

func (repo *OrganizationSubscriptionRepo) Page(ctx context.Context, pageReq dto.OrganizationSubscriptionPageRequestDto) (page gormx.PageResponse[OrganizationSubscriptionPageRow], err error) {
	subTable := model.OrganizationSubscription{}.TableName()
	orgTable := model.Organization{}.TableName()
	planTable := model.PricingPlan{}.TableName()
	filter := pageReq.Filter

	query := repo.GetDB().Table(subTable).WithContext(ctx).
		Select(fmt.Sprintf("%s.*, %s.name as organization_name, %s.name as pricing_plan_name", subTable, orgTable, planTable)).
		Joins(fmt.Sprintf("LEFT JOIN %s ON %s.organization_id = %s.id", orgTable, subTable, orgTable)).
		Joins(fmt.Sprintf("LEFT JOIN %s ON %s.pricing_plan_id = %s.id", planTable, subTable, planTable))

	if filter.OrganizationId > 0 {
		query = query.Where(fmt.Sprintf("%s.organization_id = ?", subTable), filter.OrganizationId)
	}
	if filter.PricingPlanId > 0 {
		query = query.Where(fmt.Sprintf("%s.pricing_plan_id = ?", subTable), filter.PricingPlanId)
	}
	if len(filter.Status) > 0 {
		query = query.Where(fmt.Sprintf("%s.status = ?", subTable), filter.Status)
	}

	if len(pageReq.Sorter.Column) > 0 {
		query = query.Order(fmt.Sprintf("%s.%s %s", subTable, pageReq.Sorter.Column, pageReq.Sorter.Sort))
	}
	query = query.Order(fmt.Sprintf("%s.created_at DESC", subTable))

	page, err = gormx.Page[OrganizationSubscriptionPageRow](query, gormx.PageRequest{
		Page:        pageReq.Page,
		Limit:       pageReq.Limit,
		CountColumn: fmt.Sprintf("%s.id", subTable),
	})
	return page, errorx.WithStack(err)
}
