package logic

import (
	"context"

	"github.com/duke-git/lancet/v2/slice"
	"gorm.io/gorm"

	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/gormx"
	"zeebee-backend/domains/zeebee/auth"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"

	"github.com/zeromicro/go-zero/core/threading"
)

type PricingPlanLogic struct {
	*BaseLogic
	authTopicMap map[constants.UserType]string
}

func NewPricingPlanLogic(baseLogic *BaseLogic, authTopicMap map[constants.UserType]string) *PricingPlanLogic {
	return &PricingPlanLogic{
		BaseLogic:    baseLogic,
		authTopicMap: authTopicMap,
	}
}

func (l *PricingPlanLogic) Create(ctx context.Context, input dto.PricingPlanCreateDto) (result *dto.PricingPlanEntityDto, err error) {
	var isExisted bool
	isExisted, err = l.repos.PricingPlanRepo.IsExistedByColumns(ctx, map[string]interface{}{
		"name": input.Name,
		"type": input.Type,
	})
	if err != nil {
		return
	}

	if isExisted {
		return nil, errorx.New("pricing plan name already exists under this type")
	}

	planModel, err := copierx.Trans[dto.PricingPlanCreateDto, model.PricingPlan](&input)
	if err != nil {
		return nil, err
	}

	err = l.repos.PricingPlanRepo.Transaction(ctx, func(tx *gorm.DB) error {
		var innerErr error
		if planModel.IsDefault {
			// set all other plans to non-default
			columns := map[string]interface{}{"is_default": false}
			filterColumns := map[string]interface{}{"is_default": true}
			innerErr = l.repos.PricingPlanRepo.UpdateColumnsWithTx(ctx, tx, columns, filterColumns, nil)
			if innerErr != nil {
				return innerErr
			}
		}
		planModel, innerErr = l.repos.PricingPlanRepo.CreateWithTx(ctx, tx, planModel, false)

		return innerErr
	})
	// plan, err := l.repos.PricingPlanRepo.Create(ctx, planModel)
	if err != nil {
		return nil, err
	}

	l.publishResourceReloadEvent(ctx, planModel.Type)
	return l.transToDetailDto(ctx, planModel)
}

func (l *PricingPlanLogic) Update(ctx context.Context, input dto.PricingPlanUpdateDto) (result *dto.PricingPlanEntityDto, err error) {

	var plan *model.PricingPlan
	err = l.repos.PricingPlanRepo.Transaction(ctx, func(tx *gorm.DB) error {
		var innerErr error
		plan, innerErr = l.repos.PricingPlanRepo.FindOneByIdWithTx(ctx, tx, input.Id)
		if innerErr != nil {
			return errorx.Wrap(innerErr, "fail to find pricing plan")
		}
		// set all other plans to non-default
		if !plan.IsDefault && input.IsDefault {
			columns := map[string]interface{}{"is_default": false}
			filterColumns := map[string]interface{}{"is_default": true}
			innerErr = l.repos.PricingPlanRepo.UpdateColumnsWithTx(ctx, tx, columns, filterColumns, nil)
			if innerErr != nil {
				return innerErr
			}

		}

		innerErr = copierx.Assign[dto.PricingPlanUpdateDto, model.PricingPlan](&input, plan)
		if innerErr != nil {
			return innerErr
		}
		plan, innerErr = l.repos.PricingPlanRepo.UpdateOneByIdWithTx(ctx, tx, input.Id, plan)

		return innerErr
	})

	if err != nil {
		err = errorx.Wrap(err, "fail to update pricing plan")
		return
	}

	l.publishResourceReloadEvent(ctx, plan.Type)
	return l.transToDetailDto(ctx, plan)
}

func (l *PricingPlanLogic) Delete(ctx context.Context, id int64) error {
	plan, err := l.repos.PricingPlanRepo.GetOne(ctx, id, nil)
	if err != nil {
		return errorx.Wrap(err, "pricing plan not found")
	}

	subscriptions, err := l.repos.OrganizationSubscriptionRepo.FindByPricingPlanId(ctx, id)
	if err != nil {
		return err
	}
	for _, sub := range subscriptions {
		if sub.Status != constants.OrganizationSubscriptionStatusExpired {
			return errorx.New("cannot delete pricing plan with active subscriptions")
		}
	}

	_, err = l.repos.PricingPlanRepo.DeleteByIDWithTx(ctx, l.repos.PricingPlanRepo.GetDB(), id)
	if err != nil {
		return err
	}

	l.publishResourceReloadEvent(ctx, plan.Type)
	return nil
}

func (l *PricingPlanLogic) Detail(ctx context.Context, id int64) (*dto.PricingPlanEntityDto, error) {
	plan, err := l.repos.PricingPlanRepo.FindOneWithPreloadById(ctx, id, []string{"Policies"})
	if err != nil {
		return nil, errorx.Wrap(err, "pricing plan not found")
	}
	return l.transToDetailDto(ctx, plan)
}

func (l *PricingPlanLogic) Page(ctx context.Context, req dto.PricingPlanPageRequestDto) (gormx.PageResponse[dto.PricingPlanPageEntity], error) {
	pageResult, err := l.repos.PricingPlanRepo.Page(ctx, req)
	if err != nil {
		return gormx.PageResponse[dto.PricingPlanPageEntity]{}, err
	}

	resultPtr, err := copierx.Trans[gormx.PageResponse[model.PricingPlan], gormx.PageResponse[dto.PricingPlanPageEntity]](&pageResult)
	if err != nil {
		return gormx.PageResponse[dto.PricingPlanPageEntity]{}, err
	}

	return *resultPtr, nil
}

func (l *PricingPlanLogic) SwitchStatus(ctx context.Context, id int64, status constants.CommonStatus) error {
	plan, err := l.repos.PricingPlanRepo.GetOne(ctx, id, nil)
	if err != nil {
		return errorx.Wrap(err, "pricing plan not found")
	}

	if plan.Status == status {
		return nil
	}

	plan.Status = status
	_, err = l.repos.PricingPlanRepo.UpdateOneById(ctx, id, plan)
	if err != nil {
		return err
	}

	l.publishResourceReloadEvent(ctx, plan.Type)
	l.reloadSubscribedOrganizationRoles(ctx, id, plan.Type)
	return nil
}

func (l *PricingPlanLogic) BindResources(ctx context.Context, planId int64, resourceIds []int64) (err error) {
	preloads := []gormx.Preload{
		{
			Query: "Resources",
		},
		{
			Query: "Subscriptions",
		},
	}
	plan, err := l.repos.PricingPlanRepo.GetOne(ctx, planId, preloads)
	if err != nil {
		return errorx.Wrap(err, "fail to find pricing plan")
	}

	oldResources := plan.Resources
	oldResourceIds := slice.Map(oldResources, func(index int, resource *model.Resource) int64 {
		return resource.Id
	})
	var currentResourceIds []int64

	if len(resourceIds) > 0 {
		filters := []gormx.WhereRaw{
			{
				Raw: "id in ? and type = ?",
				Args: []interface{}{
					resourceIds,
					plan.Type,
				},
			},
		}
		err = l.repos.ResourceRepo.PluckByColumns(ctx, "id", &currentResourceIds, filters)
		if err != nil {
			return errorx.Wrap(err, "fail to find resource")
		}
	}

	newIds, removedIds := diffResourceIds(oldResourceIds, currentResourceIds)
	subscriptions := plan.Subscriptions

	var removedPolicies [][]string
	err = l.repos.PricingPlanRepo.Transaction(ctx, func(tx *gorm.DB) error {
		var innerErr error
		// update pricing plan resources
		innerErr = l.repos.PricingPlanRepo.UpdateResourcesWithTx(ctx, tx, plan,
			resourceStubs(newIds), resourceStubs(removedIds))
		if innerErr != nil {
			return innerErr
		}

		if len(removedIds) == 0 {
			return nil
		}
		// get all non-expired subscriptions
		organizationIds := slice.FilterMap(subscriptions, func(index int, sub *model.OrganizationSubscription) (int64, bool) {
			if sub.Status != constants.OrganizationSubscriptionStatusExpired {
				return sub.OrganizationId, true
			}
			return 0, false
		})

		if len(organizationIds) == 0 {
			return nil
		}

		roles, innerErr := l.repos.RoleRepo.FindByOrganizationIdsWithTx(ctx, tx, organizationIds, plan.Type)
		if innerErr != nil {
			return innerErr
		}

		removedResources := resourceStubs(removedIds)
		// update role resources
		for _, role := range roles {
			innerErr = l.repos.RoleRepo.UpdateResourcesWithTx(ctx, tx, role, nil, removedResources)
			if innerErr != nil {
				return innerErr
			}
		}

		for _, role := range roles {
			for _, resourceId := range removedIds {
				removedPolicies = append(removedPolicies, auth.BuildPolicy(role.Id, resourceId, plan.Type))
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	threading.GoSafe(func() {
		subCtx := context.WithoutCancel(ctx)
		msg := auth.Msg{
			Action:          auth.ActionResourceReload,
			Scope:           auth.ScopeGlobal,
			Type:            plan.Type,
			RemovedPolicies: removedPolicies,
		}
		_ = l.publishAuthEvent(subCtx, l.authTopicMap[plan.Type], msg)
	})
	l.reloadSubscribedOrganizationRoles(ctx, planId, plan.Type)

	return
}

// diffResourceIds computes added and removed resource IDs.
func diffResourceIds(oldIds, newIds []int64) (added, removed []int64) {

	removed = slice.Difference(oldIds, newIds)
	added = slice.Difference(newIds, oldIds)
	return
}

// resourceStubs creates Resource stubs with only IDs set, for GORM association operations.
func resourceStubs(ids []int64) []*model.Resource {
	stubs := make([]*model.Resource, 0, len(ids))
	for _, id := range ids {
		stubs = append(stubs, &model.Resource{Base: model.Base{Id: id}})
	}
	return stubs
}

func (l *PricingPlanLogic) ResourceList(ctx context.Context, planId int64) ([]*dto.PricingPlanResourceEntityDto, error) {
	_, err := l.repos.PricingPlanRepo.GetOne(ctx, planId, nil)
	if err != nil {
		return nil, errorx.Wrap(err, "pricing plan not found")
	}

	bindings, err := l.repos.PricingPlanAndResourceRepo.FindByPlanId(ctx, planId)
	if err != nil {
		return nil, err
	}

	listPtr, err := copierx.Trans[[]*model.PricingPlanAndResource, []*dto.PricingPlanResourceEntityDto](&bindings)
	if err != nil {
		return nil, err
	}
	return *listPtr, nil
}

func (l *PricingPlanLogic) transToDetailDto(ctx context.Context, plan *model.PricingPlan) (*dto.PricingPlanEntityDto, error) {
	result, err := copierx.Trans[model.PricingPlan, dto.PricingPlanEntityDto](plan)
	if err != nil {
		return nil, err
	}

	// load policies if not preloaded
	// if plan.Policies == nil {
	// 	policies, innerErr := l.repos.PricingPlanPolicyRepo.ListByPlanId(ctx, plan.Id)
	// 	if innerErr != nil {
	// 		return nil, innerErr
	// 	}
	// 	plan.Policies = policies
	// }
	//
	// policyDtos, err := copierx.Trans[[]*model.PricingPlanPolicy, []dto.PricingPlanPolicyEntityDto](&plan.Policies)
	// if err != nil {
	// 	return nil, err
	// }
	// result.Policies = *policyDtos

	return result, nil
}

// reloadSubscribedOrganizationRoles triggers organization-level role reload for all orgs
// with active subscriptions to the given plan. This ensures auth state is refreshed when
// plan resources or status change.
func (l *PricingPlanLogic) reloadSubscribedOrganizationRoles(ctx context.Context, planId int64, userType constants.UserType) {
	subscriptions, err := l.repos.OrganizationSubscriptionRepo.FindByPricingPlanId(ctx, planId)
	if err != nil {
		return
	}
	for _, sub := range subscriptions {
		if sub.Status == constants.OrganizationSubscriptionStatusExpired {
			continue
		}
		threading.GoSafe(func() {
			subCtx := context.WithoutCancel(ctx)
			msg := auth.Msg{
				Action:         auth.ActionRoleReload,
				Scope:          auth.ScopeOrganization,
				Type:           userType,
				OrganizationId: sub.OrganizationId,
			}
			_ = l.publishAuthEvent(subCtx, l.authTopicMap[userType], msg)
		})
	}
}

func (l *PricingPlanLogic) publishResourceReloadEvent(ctx context.Context, userType constants.UserType) {
	threading.GoSafe(func() {
		subCtx := context.WithoutCancel(ctx)
		msg := auth.Msg{
			Action: auth.ActionResourceReload,
			Scope:  auth.ScopeGlobal,
			Type:   userType,
		}
		topic := l.authTopicMap[userType]
		_ = l.publishAuthEvent(subCtx, topic, msg)
	})
}
