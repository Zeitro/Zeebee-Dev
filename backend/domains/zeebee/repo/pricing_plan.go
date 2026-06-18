package repo

import (
	"context"
	"fmt"
	"sort"

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

func (repo *PricingPlanRepo) GetOneWithTx(ctx context.Context, tx *gorm.DB, id int64, preloads []gormx.Preload) (result *model.PricingPlan, err error) {
	query := tx.WithContext(ctx).Model(&model.PricingPlan{}).Where("id = ?", id)
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

// FindAllByType returns all active plans of the given type with policies preloaded,
// ordered by sort_order for option/dropdown lists.
func (repo *PricingPlanRepo) FindAllByType(ctx context.Context, planType constants.UserType, isPublic bool) (result []*model.PricingPlan, err error) {
	query := repo.Client(ctx).
		// Select("id", "name", "status", "type", "is_public", "sort_order").
		Where("type = ? AND status = ?", planType, constants.CommonStatusActive)
	if isPublic {
		query = query.Where("is_public = ?", true)
	}

	err = query.Preload("Policies").
		Order("sort_order asc").
		Find(&result).Error
	return result, errorx.WithStack(err)
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

	query := repo.Client(ctx)

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
	query = query.Order(fmt.Sprintf("%s.sort_order asc", table))
	if len(pageReq.Sorter.Column) > 0 {
		query = query.Order(fmt.Sprintf("%s.%s %s", table, pageReq.Sorter.Column, pageReq.Sorter.Sort))
	}
	query = query.Order(fmt.Sprintf("%s.created_at DESC", table))

	page, err = gormx.Page[model.PricingPlan](query, gormx.PageRequest{
		Page:        pageReq.Page,
		Limit:       pageReq.Limit,
		CountColumn: fmt.Sprintf("%s.id", table),
	})
	if err != nil {
		return page, errorx.WithStack(err)
	}

	if len(page.List) > 0 {
		policyTypeOrder := map[constants.PricingPolicyType]int{
			constants.PricingPolicyTypeMonthly:   0,
			constants.PricingPolicyTypeQuarterly: 1,
			constants.PricingPolicyTypeYearly:    2,
		}
		policyRepo := NewPricingPlanPolicyRepo(repo.GetDB())

		for i := range page.List {
			var policies []*model.PricingPlanPolicy
			policyErr := policyRepo.Client(ctx).
				Where("pricing_plan_id = ?", page.List[i].Id).
				Find(&policies).Error
			if policyErr != nil {
				return page, errorx.WithStack(policyErr)
			}
			sort.Slice(policies, func(a, b int) bool {
				return policyTypeOrder[policies[a].Type] < policyTypeOrder[policies[b].Type]
			})
			page.List[i].Policies = policies
		}
	}

	return page, nil
}

// DeleteWithTx deletes a pricing plan within a transaction.
// It removes the plan's resource associations, deletes the plan's pricing
// policies, then deletes the plan itself.
func (repo *PricingPlanRepo) DeleteWithTx(ctx context.Context, tx *gorm.DB, plan *model.PricingPlan) (err error) {
	if plan == nil {
		return
	}

	// nil out loaded associations first: GORM's Association().Clear() would
	// otherwise encode the fully-loaded Resource structs (with model.Base
	// timestamps) into the join-table delete clause → OID 0 encoding error.
	plan.RemoveRelationships()

	// remove resource associations (join table has no soft-delete column → hard delete)
	err = tx.WithContext(ctx).Model(plan).Association("Resources").Clear()
	if err != nil {
		return errorx.WithStack(err)
	}

	// delete the plan's pricing policies (soft-delete via model.Base)
	err = tx.WithContext(ctx).
		Where("pricing_plan_id = ?", plan.Id).
		Delete(&model.PricingPlanPolicy{}).Error
	if err != nil {
		return errorx.WithStack(err)
	}

	// delete the plan (soft-delete)
	err = tx.WithContext(ctx).Delete(plan).Error
	return errorx.WithStack(err)
}
