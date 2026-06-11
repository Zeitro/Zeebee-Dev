package repo

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"zeebee-backend/common/copierx"
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

func (repo *PricingPlanRepo) GetOne(ctx context.Context, id int64, preloads []gormx.Preload) (result *model.PricingPlan, err error) {
	query := repo.Client(ctx).Where("id = ?", id)
	for i := range preloads {
		args := make([]interface{}, 0, len(preloads[i].Args))
		for _, arg := range preloads[i].Args {
			args = append(args, arg)
		}
		query = query.Preload(preloads[i].Query, args...)
	}
	err = query.First(&result).Error
	return result, errorx.WithStack(err)
}

func (repo *PricingPlanRepo) ExistsByNameAndType(ctx context.Context, name string, planType constants.UserType) (bool, error) {
	return repo.IsExistedByColumns(ctx, map[string]interface{}{
		"name": name,
		"type": planType,
	})
}

func (repo *PricingPlanRepo) FindDefaultByType(ctx context.Context, planType constants.UserType) (result *model.PricingPlan, err error) {
	err = repo.Client(ctx).Where("type = ? AND is_default = ? AND status = ?", planType, true, constants.CommonStatusActive).First(&result).Error
	if err != nil {
		return nil, nil
	}
	return result, nil
}

// UpdateResourcesWithTx updates a plan's resource associations within a transaction.
// New resources are appended, removed resources are deleted, unchanged ones are left alone.
func (repo *PricingPlanRepo) UpdateResourcesWithTx(ctx context.Context, tx *gorm.DB, plan *model.PricingPlan, newResources []*model.Resource, removedResources []*model.Resource) error {
	copyPlan, err := copierx.Trans[model.PricingPlan, model.PricingPlan](plan)
	if err != nil {
		return err
	}
	copyPlan.RemoveRelationships()
	query := tx.WithContext(ctx).Model(copyPlan)
	if len(newResources) > 0 {
		err = query.Association("Resources").Append(newResources)
		if err != nil {
			return errorx.WithStack(err)
		}
	}
	copyPlan.RemoveRelationships()
	if len(removedResources) > 0 {
		err = query.Association("Resources").Delete(removedResources)
		if err != nil {
			return errorx.WithStack(err)
		}
	}
	return nil
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
