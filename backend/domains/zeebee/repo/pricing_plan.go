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

type PricingPlanRepo struct {
	*gormx.BaseRepo[*model.PricingPlan]
}

func NewPricingPlanRepo(db *gorm.DB) *PricingPlanRepo {
	return &PricingPlanRepo{
		BaseRepo: gormx.NewBaseRepo[*model.PricingPlan](db),
	}
}

func (repo *PricingPlanRepo) GetOne(ctx context.Context, id int64) (result *model.PricingPlan, err error) {
	err = repo.Client(ctx).Where("id = ?", id).First(&result).Error
	return result, errorx.WithStack(err)
}

func (repo *PricingPlanRepo) ExistsByNameAndType(ctx context.Context, name string, planType constants.UserType) (bool, error) {
	return repo.IsExistedByColumns(ctx, map[string]interface{}{
		"name": name,
		"type": planType,
	})
}

func (repo *PricingPlanRepo) Page(ctx context.Context, pageReq dto.PricingPlanPageRequestDto) (page gormx.PageResponse[model.PricingPlan], err error) {
	table := model.PricingPlan{}.TableName()
	filter := pageReq.Filter

	query := repo.GetDB().Table(table).WithContext(ctx)

	if len(filter.Keyword) > 0 {
		keyword := "%" + filter.Keyword + "%"
		query = query.Where("name LIKE ?", keyword)
	}
	if len(filter.Status) > 0 {
		query = query.Where("status = ?", filter.Status)
	}
	if len(filter.Type) > 0 {
		query = query.Where("type = ?", filter.Type)
	}
	if len(filter.PlanType) > 0 {
		query = query.Where("plan_type = ?", filter.PlanType)
	}

	if len(pageReq.Sorter.Column) > 0 {
		query = query.Order(fmt.Sprintf("%s.%s %s", table, pageReq.Sorter.Column, pageReq.Sorter.Sort))
	}
	query = query.Order(fmt.Sprintf("%s.created_at DESC", table))

	page, err = gormx.Page[model.PricingPlan](query, gormx.PageRequest{
		Page:        pageReq.Page,
		Limit:       pageReq.Limit,
		CountColumn: fmt.Sprintf("%s.id", table),
	})
	return page, errorx.WithStack(err)
}
