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

func (l *PricingPlanLogic) Create(ctx context.Context, input dto.PricingPlanCreateRequestDto) (result *dto.PricingPlanEntityDto, err error) {
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

	planModel, err := copierx.Trans[dto.PricingPlanCreateRequestDto, model.PricingPlan](&input)
	if err != nil {
		return nil, err
	}
	planModel.Status = constants.CommonStatusActive

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

	if err != nil {
		return nil, err
	}

	return l.transToDetailDto(ctx, planModel)
}

func (l *PricingPlanLogic) Update(ctx context.Context, input dto.PricingPlanUpdateRequestDto) (result *dto.PricingPlanEntityDto, err error) {

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

		innerErr = copierx.Assign[dto.PricingPlanUpdateRequestDto, model.PricingPlan](&input, plan)
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

	return l.transToDetailDto(ctx, plan)
}

func (l *PricingPlanLogic) Delete(ctx context.Context, id int64) (err error) {
	_, err = l.repos.PricingPlanRepo.GetOne(ctx, id, nil)
	if err != nil {
		return errorx.Wrap(err, "fail to find pricing plan")
	}

	return l.repos.PricingPlanRepo.Transaction(ctx, func(tx *gorm.DB) error {
		// check for active subscriptions
		effectedSubscriptions, innerErr := l.repos.OrganizationSubscriptionRepo.FindEffectiveByPlanIdWithTx(ctx, tx, id)
		if innerErr != nil {
			return innerErr
		}
		if len(effectedSubscriptions) > 0 {
			return errorx.New("cannot delete pricing plan with active subscriptions")
		}
		// delete plan
		_, innerErr = l.repos.PricingPlanRepo.DeleteByIDWithTx(ctx, tx, id)
		return innerErr
	})
}

func (l *PricingPlanLogic) Detail(ctx context.Context, id int64) (*dto.PricingPlanEntityDto, error) {
	plan, err := l.repos.PricingPlanRepo.FindOneWithPreloadById(ctx, id, []string{"Policies"})
	if err != nil {
		return nil, errorx.Wrap(err, "fail to find pricing plan")
	}
	return l.transToDetailDto(ctx, plan)
}

func (l *PricingPlanLogic) Page(ctx context.Context, req dto.PricingPlanPageRequestDto) (gormx.PageResponse[dto.PricingPlanPageEntityDto], error) {
	pageResult, err := l.repos.PricingPlanRepo.Page(ctx, req)
	if err != nil {
		return gormx.PageResponse[dto.PricingPlanPageEntityDto]{}, err
	}

	resultPtr, err := copierx.Trans[gormx.PageResponse[model.PricingPlan], gormx.PageResponse[dto.PricingPlanPageEntityDto]](&pageResult)
	if err != nil {
		return gormx.PageResponse[dto.PricingPlanPageEntityDto]{}, err
	}

	return *resultPtr, nil
}

func (l *PricingPlanLogic) SwitchStatus(ctx context.Context, id int64, status constants.CommonStatus) (err error) {
	preloads := []gormx.Preload{
		{Query: "Resources"},
		// {Query: "Subscriptions"},
	}
	plan, err := l.repos.PricingPlanRepo.GetOne(ctx, id, preloads)
	if err != nil {
		return errorx.Wrap(err, "fail to find pricing plan")
	}

	if plan.Status == status {
		return nil
	}

	isDisabling := status == constants.CommonStatusDisabled
	resourceIds := slice.Map(plan.Resources, func(index int, resource *model.Resource) int64 {
		return resource.Id
	})

	var affectedOrganizationIds []int64

	var orgPolicies map[int64]*organizationPolicyChanges
	plan.Status = status
	err = l.repos.PricingPlanRepo.Transaction(ctx, func(tx *gorm.DB) error {
		var innerErr error
		// find affected subscriptions
		effectedSubscriptions, innerErr := l.repos.OrganizationSubscriptionRepo.FindEffectiveByPlanIdWithTx(ctx, tx, id)
		if innerErr != nil {
			return innerErr
		}
		affectedOrganizationIds = slice.Map(effectedSubscriptions, func(index int, subscription *model.OrganizationSubscription) int64 {
			return subscription.OrganizationId
		})
		affectedOrganizationIds = slice.Unique(affectedOrganizationIds)

		_, innerErr = l.repos.PricingPlanRepo.UpdateOneByIdWithTx(ctx, tx, id, plan)
		if innerErr != nil {
			return innerErr
		}

		if len(affectedOrganizationIds) == 0 {
			return nil
		}
		// bind/unbind resources
		if isDisabling {
			orgPolicies, innerErr = l.syncPlanResourcesToOrganizationRoles(ctx, tx, affectedOrganizationIds, resourceIds, plan.Type, false)
		} else {
			orgPolicies, innerErr = l.syncPlanResourcesToOrganizationRoles(ctx, tx, affectedOrganizationIds, resourceIds, plan.Type, true)
		}
		return innerErr
	})
	if err != nil {
		return err
	}

	l.publishOrgRoleReloadEvents(ctx, orgPolicies, plan.Type)
	return
}

func (l *PricingPlanLogic) BindResources(ctx context.Context, planId int64, resourceIds []int64) (err error) {
	preloads := []gormx.Preload{
		{Query: "Resources"},
		// {Query: "Subscriptions"},
	}
	plan, err := l.repos.PricingPlanRepo.GetOne(ctx, planId, preloads)
	if err != nil {
		return errorx.Wrap(err, "fail to find pricing plan")
	}

	oldResourceIds := slice.Map(plan.Resources, func(index int, resource *model.Resource) int64 {
		return resource.Id
	})

	var currentResourceIds []int64
	if len(resourceIds) > 0 {
		filters := []gormx.WhereRaw{
			{Raw: "id in ? and type = ?", Args: []interface{}{resourceIds, plan.Type}},
		}
		err = l.repos.ResourceRepo.PluckByColumns(ctx, "id", &currentResourceIds, filters)
		if err != nil {
			return errorx.Wrap(err, "fail to find resource")
		}
	}

	newIds, removedIds := diffResourceIds(oldResourceIds, currentResourceIds)

	organizationPolicies := make(map[int64]*organizationPolicyChanges)
	err = l.repos.PricingPlanRepo.Transaction(ctx, func(tx *gorm.DB) error {
		var innerErr error
		innerErr = l.repos.PricingPlanRepo.UpdateResourcesWithTx(ctx, tx, plan,
			resourceStubs(newIds), resourceStubs(removedIds))
		if innerErr != nil {
			return innerErr
		}

		effectedSubscriptions, innerErr := l.repos.OrganizationSubscriptionRepo.FindEffectiveByPlanIdWithTx(ctx, tx, plan.Id)
		if innerErr != nil {
			return innerErr
		}
		organizationIds := slice.Map(effectedSubscriptions, func(index int, subscription *model.OrganizationSubscription) int64 {
			return subscription.OrganizationId
		})
		organizationIds = slice.Unique(organizationIds)

		if len(organizationIds) == 0 {
			return nil
		}

		if len(removedIds) > 0 {
			removedOrgPolicies, innerErr := l.syncPlanResourcesToOrganizationRoles(ctx, tx, organizationIds, removedIds, plan.Type, false)
			if innerErr != nil {
				return innerErr
			}
			mergeOrganizationPolicies(organizationPolicies, removedOrgPolicies)
		}

		if len(newIds) > 0 {
			addedOrgPolicies, innerErr := l.syncPlanResourcesToOrganizationRoles(ctx, tx, organizationIds, newIds, plan.Type, true)
			if innerErr != nil {
				return innerErr
			}
			mergeOrganizationPolicies(organizationPolicies, addedOrgPolicies)
		}

		return nil
	})
	if err != nil {
		return err
	}

	l.publishOrgRoleReloadEvents(ctx, organizationPolicies, plan.Type)
	return
}

// orgPolicyChanges tracks per-org policy changes for auth event publishing.
type organizationPolicyChanges struct {
	added   [][]string
	removed [][]string
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

// syncPlanResourcesToOrganizationRoles binds/unbinds plan resources to organization roles.
func (l *PricingPlanLogic) syncPlanResourcesToOrganizationRoles(ctx context.Context, tx *gorm.DB, organizationIds []int64, resourceIds []int64, userType constants.UserType, isBind bool) (map[int64]*organizationPolicyChanges, error) {
	if len(organizationIds) == 0 || len(resourceIds) == 0 {
		return nil, nil
	}
	roles, err := l.repos.RoleRepo.FindByOrganizationIdsWithTx(ctx, tx, organizationIds, userType)
	if err != nil {
		return nil, err
	}
	var removedResources []*model.Resource
	var addedResources []*model.Resource
	if isBind {
		addedResources = resourceStubs(resourceIds)
	} else {
		removedResources = resourceStubs(resourceIds)
	}
	organizationPolicies := make(map[int64]*organizationPolicyChanges)

	for _, role := range roles {
		err = l.repos.RoleRepo.UpdateResourcesWithTx(ctx, tx, role, addedResources, removedResources)
		if err != nil {
			return nil, err
		}
		if organizationPolicies[role.OrganizationId] == nil {
			organizationPolicies[role.OrganizationId] = &organizationPolicyChanges{}
		}
		for _, resourceId := range resourceIds {
			if isBind {
				organizationPolicies[role.OrganizationId].added = append(organizationPolicies[role.OrganizationId].added, auth.BuildPolicy(role.Id, resourceId, userType))
			} else {
				organizationPolicies[role.OrganizationId].removed = append(organizationPolicies[role.OrganizationId].removed, auth.BuildPolicy(role.Id, resourceId, userType))
			}
		}
	}
	return organizationPolicies, nil
}

// mergeOrganizationPolicies merges source policies into target.
func mergeOrganizationPolicies(target map[int64]*organizationPolicyChanges, source map[int64]*organizationPolicyChanges) {
	for organizationId, srcChanges := range source {
		if target[organizationId] == nil {
			target[organizationId] = &organizationPolicyChanges{}
		}
		target[organizationId].added = append(target[organizationId].added, srcChanges.added...)
		target[organizationId].removed = append(target[organizationId].removed, srcChanges.removed...)
	}
}

// publishOrgRoleReloadEvents sends per-org RoleReload events with explicit policy changes.
func (l *PricingPlanLogic) publishOrgRoleReloadEvents(ctx context.Context, organizationPolicies map[int64]*organizationPolicyChanges, userType constants.UserType) {
	for orgId, changes := range organizationPolicies {
		orgId := orgId
		changes := changes
		threading.GoSafe(func() {
			subCtx := context.WithoutCancel(ctx)
			msg := auth.Msg{
				Action:          auth.ActionRoleReload,
				Scope:           auth.ScopeOrganization,
				Type:            userType,
				OrganizationId:  orgId,
				AddedPolicies:   changes.added,
				RemovedPolicies: changes.removed,
			}
			_ = l.publishAuthEvent(subCtx, l.authTopicMap[userType], msg)
		})
	}
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
	return result, nil
}
