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
	"zeebee-backend/domains/zeebee/repo"

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
	plan, err := l.repos.PricingPlanRepo.GetOne(ctx, id, []gormx.Preload{
		{Query: "Policies"},
		{Query: "Resources"},
		{Query: "Subscriptions"},
	})
	if err != nil {
		return errorx.Wrap(err, "fail to find pricing plan")
	}

	err = l.repos.PricingPlanRepo.Transaction(ctx, func(tx *gorm.DB) error {
		var innerErr error
		// check for active subscriptions
		// effectedSubscriptions, innerErr := l.repos.OrganizationSubscriptionRepo.FindEffectiveByPlanIdWithTx(ctx, tx, id)
		// get active subscriptions
		effectedSubscriptions := slice.Filter(plan.Subscriptions, func(index int, item *model.OrganizationSubscription) bool {
			return item.Status == constants.OrganizationSubscriptionStatusActive ||
				item.Status == constants.OrganizationSubscriptionStatusExpiring ||
				item.Status == constants.OrganizationSubscriptionStatusSuspended ||
				item.Status == constants.OrganizationSubscriptionStatusTrialing
		})

		// if innerErr != nil {
		// 	return innerErr
		// }
		if len(effectedSubscriptions) > 0 {
			return errorx.New("cannot delete pricing plan with active subscriptions")
		}
		// delete plan
		innerErr = l.repos.PricingPlanRepo.DeleteWithTx(ctx, tx, plan)
		return innerErr
	})

	return errorx.WithStack(err)
}

func (l *PricingPlanLogic) Detail(ctx context.Context, id int64) (*dto.PricingPlanEntityDto, error) {
	plan, err := l.repos.PricingPlanRepo.FindOneWithPreloadById(ctx, id, []string{"Policies", "Resources"})
	if err != nil {
		return nil, errorx.Wrap(err, "fail to find pricing plan")
	}

	policies := plan.Policies
	policyTypeMap := map[constants.PricingPolicyType]int{
		constants.PricingPolicyTypeMonthly:   0,
		constants.PricingPolicyTypeQuarterly: 1,
		constants.PricingPolicyTypeYearly:    2,
	}
	slice.SortBy(policies, func(a, b *model.PricingPlanPolicy) bool {
		return policyTypeMap[a.Type] < policyTypeMap[b.Type]
	})
	plan.Policies = policies

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

// All returns all active plans of the given type as option entities (id/name/type/sort_order/policies),
// ordered by sort_order with policies sorted monthly -> quarterly -> yearly.
func (l *PricingPlanLogic) All(ctx context.Context, planType constants.UserType, isPublic bool) ([]*dto.PricingPlanOptEntityDto, error) {
	plans, err := l.repos.PricingPlanRepo.FindAllByType(ctx, planType, isPublic)
	if err != nil {
		return nil, errorx.Wrap(err, "fail to find pricing plans")
	}

	policyTypeMap := map[constants.PricingPolicyType]int{
		constants.PricingPolicyTypeMonthly:   0,
		constants.PricingPolicyTypeQuarterly: 1,
		constants.PricingPolicyTypeYearly:    2,
	}

	result := make([]*dto.PricingPlanOptEntityDto, 0, len(plans))
	for _, plan := range plans {
		slice.SortBy(plan.Policies, func(a, b *model.PricingPlanPolicy) bool {
			return policyTypeMap[a.Type] < policyTypeMap[b.Type]
		})
		opt, err := copierx.Trans[model.PricingPlan, dto.PricingPlanOptEntityDto](plan)
		if err != nil {
			return nil, err
		}
		result = append(result, opt)
	}
	return result, nil
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
			orgPolicies, innerErr = syncPlanResourcesToOrganizationRoles(ctx, l.repos, tx, planResourceSyncReq{
				OrganizationIds: affectedOrganizationIds,
				ResourceIds:     resourceIds,
				UserType:        plan.Type,
				IsBind:          false,
			})
		} else {
			orgPolicies, innerErr = syncPlanResourcesToOrganizationRoles(ctx, l.repos, tx, planResourceSyncReq{
				OrganizationIds: affectedOrganizationIds,
				ResourceIds:     resourceIds,
				UserType:        plan.Type,
				IsBind:          true,
			})
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
			removedOrgPolicies, innerErr := syncPlanResourcesToOrganizationRoles(ctx, l.repos, tx, planResourceSyncReq{
				OrganizationIds: organizationIds,
				ResourceIds:     removedIds,
				UserType:        plan.Type,
				IsBind:          false,
			})
			if innerErr != nil {
				return innerErr
			}
			mergeOrganizationPolicies(organizationPolicies, removedOrgPolicies)
		}

		if len(newIds) > 0 {
			addedOrgPolicies, innerErr := syncPlanResourcesToOrganizationRoles(ctx, l.repos, tx, planResourceSyncReq{
				OrganizationIds: organizationIds,
				ResourceIds:     newIds,
				UserType:        plan.Type,
				IsBind:          true,
			})
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
		return nil, errorx.Wrap(err, "fail to find pricing plan")
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

// planResourceSyncReq bundles inputs for syncPlanResourcesToOrganizationRoles.
type planResourceSyncReq struct {
	OrganizationIds []int64
	ResourceIds     []int64
	UserType        constants.UserType
	IsBind          bool
}

// syncPlanResourcesToOrganizationRoles binds/unbinds plan resources to organization roles.
func syncPlanResourcesToOrganizationRoles(ctx context.Context, repos *repo.AllRepos, tx *gorm.DB, req planResourceSyncReq) (map[int64]*organizationPolicyChanges, error) {
	if len(req.OrganizationIds) == 0 || len(req.ResourceIds) == 0 {
		return nil, nil
	}
	roles, err := repos.RoleRepo.FindByOrganizationIdsWithTx(ctx, tx, req.OrganizationIds, req.UserType)
	if err != nil {
		return nil, err
	}
	var removedResources []*model.Resource
	var addedResources []*model.Resource
	if req.IsBind {
		addedResources = resourceStubs(req.ResourceIds)
	} else {
		removedResources = resourceStubs(req.ResourceIds)
	}
	organizationPolicies := make(map[int64]*organizationPolicyChanges)

	for _, role := range roles {
		err = repos.RoleRepo.UpdateResourcesWithTx(ctx, tx, role, addedResources, removedResources)
		if err != nil {
			return nil, err
		}
		if organizationPolicies[role.OrganizationId] == nil {
			organizationPolicies[role.OrganizationId] = &organizationPolicyChanges{}
		}
		for _, resourceId := range req.ResourceIds {
			if req.IsBind {
				organizationPolicies[role.OrganizationId].added = append(organizationPolicies[role.OrganizationId].added, auth.BuildPolicy(role.Id, resourceId, req.UserType))
			} else {
				organizationPolicies[role.OrganizationId].removed = append(organizationPolicies[role.OrganizationId].removed, auth.BuildPolicy(role.Id, resourceId, req.UserType))
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
