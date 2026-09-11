package repo

import (
	"context"
	"database/sql"
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

type OrganizationSubscriptionPageRow struct {
	Id                    int64                       `gorm:"column:id"`
	Code                  string                      `gorm:"column:code"`
	OrganizationId        int64                       `gorm:"column:organization_id"`
	OrganizationName      string                      `gorm:"column:organization_name"`
	PricingPlanId         int64                       `gorm:"column:pricing_plan_id"`
	PricingPlanName       string                      `gorm:"column:pricing_plan_name"`
	PricingPlanPolicyId   int64                       `gorm:"column:pricing_plan_policy_id"`
	PricingPlanPolicyType constants.PricingPolicyType `gorm:"column:pricing_plan_policy_type"`
	Status                string                      `gorm:"column:status"`
	ExpireDate            *sql.NullTime               `gorm:"column:expire_date"`
	OriginalAmount        float64                     `gorm:"column:original_amount"`
	ActualAmount          float64                     `gorm:"column:actual_amount"`
	IsAutoRenew           bool                        `gorm:"column:is_auto_renew"`
	CreatedAt             time.Time                   `gorm:"column:created_at"`
	UpdatedAt             time.Time                   `gorm:"column:updated_at"`
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

func (repo *OrganizationSubscriptionRepo) FindActiveByOrganizationId(ctx context.Context, organizationId int64) (result []*model.OrganizationSubscription, err error) {
	err = repo.Client(ctx).
		Where("organization_id = ? AND status IN ?",
			organizationId,
			constants.EffectiveStatuses,
		).
		Order("id desc").
		Find(&result).Error
	return result, errorx.WithStack(err)
}

func (repo *OrganizationSubscriptionRepo) FindActiveByOrgIdWithTx(ctx context.Context, tx *gorm.DB, organizationId int64) (result []*model.OrganizationSubscription, err error) {
	err = tx.WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("organization_id = ? AND status IN ?",
			organizationId,
			constants.EffectiveStatuses,
		).
		Order("id desc").
		Find(&result).Error
	return result, errorx.WithStack(err)
}

func (repo *OrganizationSubscriptionRepo) FindActiveSubscriptions(ctx context.Context, organizationType constants.UserType) (result []*model.OrganizationSubscription, err error) {
	organizationRepo := NewOrganizationRepo(repo.GetDB())
	var organizationIds []int64
	organizationFilter := []gormx.WhereRaw{
		{
			Raw: "type = ?",
			Args: []interface{}{
				organizationType,
			},
		},
	}
	err = organizationRepo.PluckByColumns(ctx, "id", &organizationIds, organizationFilter)
	if err != nil {
		return result, errorx.WithStack(err)
	}
	if len(organizationIds) == 0 {
		return
	}

	err = repo.Client(ctx).
		Where("organization_id in ? AND status IN ?",
			organizationIds,
			constants.EffectiveStatuses,
		).
		Order("id desc").
		Find(&result).Error
	return result, errorx.WithStack(err)
}

func (repo *OrganizationSubscriptionRepo) FindActiveSubscriptionsWithTx(ctx context.Context, tx *gorm.DB, organizationType constants.UserType) (result []*model.OrganizationSubscription, err error) {
	organizationRepo := NewOrganizationRepo(tx)
	var organizationIds []int64
	organizationFilter := []gormx.WhereRaw{
		{
			Raw: "type = ?",
			Args: []interface{}{
				organizationType,
			},
		},
	}
	err = organizationRepo.PluckByColumnsWithTx(ctx, tx, "id", &organizationIds, organizationFilter)
	if err != nil {
		return result, errorx.WithStack(err)
	}
	if len(organizationIds) == 0 {
		return
	}

	err = tx.WithContext(ctx).Model(&model.OrganizationSubscription{}).
		Where("organization_id in ? and status IN ?",
			organizationIds,
			constants.EffectiveStatuses,
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
	policyTable := model.PricingPlanPolicy{}.TableName()
	filter := pageReq.Filter

	query := repo.Client(ctx).
		Select(fmt.Sprintf("%s.*, %s.name as organization_name, %s.name as pricing_plan_name, %s.type as pricing_plan_policy_type", subTable, orgTable, planTable, policyTable)).
		Joins(fmt.Sprintf("LEFT JOIN %s ON %s.organization_id = %s.id", orgTable, subTable, orgTable)).
		Joins(fmt.Sprintf("LEFT JOIN %s ON %s.pricing_plan_id = %s.id", planTable, subTable, planTable)).
		Joins(fmt.Sprintf("LEFT JOIN %s ON %s.pricing_plan_policy_id = %s.id", policyTable, subTable, policyTable))

	if len(filter.Keyword) > 0 {
		query = query.Where(
			fmt.Sprintf("%s.code LIKE @keyword OR %s.name LIKE @keyword OR %s.name LIKE @keyword", subTable, orgTable, planTable),
			sql.Named("keyword", "%"+filter.Keyword+"%"),
		)
	}
	if len(filter.Type) > 0 {
		query = query.Where(fmt.Sprintf("%s.type = ?", orgTable), filter.Type)
	}
	if filter.OrganizationId > 0 {
		query = query.Where(fmt.Sprintf("%s.organization_id = ?", subTable), filter.OrganizationId)
	}
	if len(filter.Status) > 0 {
		query = query.Where(fmt.Sprintf("%s.status IN ?", subTable), filter.Status)
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

// FindEffectiveByPlanIdWithTx get all affective subscription by plan id.
// Uses constants.EffectiveStatuses — single source of truth for "effective".
func (repo *OrganizationSubscriptionRepo) FindEffectiveByPlanIdWithTx(ctx context.Context, tx *gorm.DB, pricingPlanId int64) (results []*model.OrganizationSubscription, err error) {
	err = tx.WithContext(ctx).Model(new(model.OrganizationSubscription)).
		Where("pricing_plan_id = ? and status in ?", pricingPlanId, constants.EffectiveStatuses).Find(&results).Error
	return results, errorx.WithStack(err)
}
