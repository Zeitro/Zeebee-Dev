package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/singleflight"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"zeebee-backend/common/errorx"
	"zeebee-backend/common/mapx"
	"zeebee-backend/common/rbac"
	"zeebee-backend/common/redisx"
	"zeebee-backend/common/stringx"
	"zeebee-backend/common/utils"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/repo"
)

// policyEnforcer is the *rbac.Enforcer subset used by rebuild/auth-check paths.
// Interface so tests can substitute a fake; production passes the real enforcer.
type policyEnforcer interface {
	AllPolicies() ([][]string, error)
	AddPolicies([][]string) error
	RemovePolicies([][]string) error
	Enforce(rvals ...interface{}) (bool, error)
	SavePolicy() error
	LoadDefaultFilteredPolicy() error
}

type Auth struct {
	repos                       *repo.AllRepos
	redisCli                    *redisx.Client
	enforcer                    policyEnforcer
	jwtConfig                   dto.JwtConfig
	authenticationWhitelistMap  *mapx.Map[constants.UserType, []string] // login api whitelist
	authorizationWhitelistMap   *mapx.Map[constants.UserType, []string] // permission api whitelist
	apiResourceMap              *mapx.Map[constants.UserType, *mapx.Map[string, []int64]]
	organizationRoleResourceMap *mapx.Map[int64, *mapx.Map[string, RoleResource]]
	organizationPlanResourceMap *mapx.Map[int64, *mapx.Map[int64, bool]] // orgId -> set of plan resource IDs
	sf                          singleflight.Group
	eventTopic                  string
	eventSubScribeCancelFn      context.CancelFunc
}

func NewAuth(repos *repo.AllRepos, redisCli *redisx.Client, jwtConfig dto.JwtConfig, enforcer *rbac.Enforcer, authEventTopic string) *Auth {
	ctx, cancel := context.WithCancel(context.Background())
	auth := &Auth{
		repos:                       repos,
		redisCli:                    redisCli,
		jwtConfig:                   jwtConfig,
		enforcer:                    enforcer,
		authenticationWhitelistMap:  mapx.NewMap[constants.UserType, []string](),
		authorizationWhitelistMap:   mapx.NewMap[constants.UserType, []string](),
		apiResourceMap:              mapx.NewMap[constants.UserType, *mapx.Map[string, []int64]](),
		organizationRoleResourceMap: mapx.NewMap[int64, *mapx.Map[string, RoleResource]](),
		organizationPlanResourceMap: mapx.NewMap[int64, *mapx.Map[int64, bool]](),
		eventTopic:                  authEventTopic,
		eventSubScribeCancelFn:      cancel,
	}

	go func() {
		auth.subscribeEvent(ctx)
	}()

	return auth
}

// newAuthWithEnforcer is a test seam: Auth with a fake enforcer and no redis
// loop, for unit-testing rebuild/auth-check without DB or redis.
func newAuthWithEnforcer(repos *repo.AllRepos, redisCli *redisx.Client, jwtConfig dto.JwtConfig, enforcer policyEnforcer, authEventTopic string) *Auth {
	return &Auth{
		repos:                       repos,
		redisCli:                    redisCli,
		jwtConfig:                   jwtConfig,
		enforcer:                    enforcer,
		authenticationWhitelistMap:  mapx.NewMap[constants.UserType, []string](),
		authorizationWhitelistMap:   mapx.NewMap[constants.UserType, []string](),
		apiResourceMap:              mapx.NewMap[constants.UserType, *mapx.Map[string, []int64]](),
		organizationRoleResourceMap: mapx.NewMap[int64, *mapx.Map[string, RoleResource]](),
		organizationPlanResourceMap: mapx.NewMap[int64, *mapx.Map[int64, bool]](),
		eventTopic:                  authEventTopic,
	}
}

// subscribeEvent
func (l *Auth) subscribeEvent(ctx context.Context) {
	if len(l.eventTopic) == 0 {
		return
	}

	retryCount := 0
	for {
		retryCount++
		if retryCount > maxRetryTimes {
			return
		}
		if ctx.Err() != nil {
			return
		}
		subPub := l.redisCli.GetClient().Subscribe(ctx, l.eventTopic)
		err := subPub.Ping(ctx, "ping")
		if err != nil {
			logx.WithContext(ctx).WithFields(
				logx.Field("stack", errorx.GetErrorStack(err)),
				logx.Field("type", "auth_message_event"),
			).Errorf("subscribeEvent error")
			continue
		}
		retryCount = 0
	handleLoop:
		for {
			select {
			case <-ctx.Done():
				if subPub != nil {
					_ = subPub.Close()
				}
				return
			case msg, ok := <-subPub.Channel():
				if !ok {
					break handleLoop
				}

				err = l.handleAuthMsg(ctx, msg)
				if err != nil {
					logx.WithContext(ctx).WithFields(
						logx.Field("stack", errorx.GetErrorStack(err)),
						logx.Field("msg", msg),
						logx.Field("type", "auth_message_event"),
						logx.Field("error", err.Error()),
					).Errorf("handleAuthMsg error")
				}
			}
		}
		_ = subPub.Close()
		sleepInterval := retryIntervalTime * (1 << retryCount)
		if sleepInterval > maxRetryIntervalTime {
			sleepInterval = maxRetryIntervalTime
		}
		time.Sleep(sleepInterval)
	}
}

// handleAuthMsg
func (l *Auth) handleAuthMsg(ctx context.Context, msg *redis.Message) (err error) {
	var eventMsg Msg
	if msg == nil {
		err = errorx.New("auth message is nil")
		return
	}
	eventMsg, err = ParseMsg(msg.Payload)
	if err != nil {
		return err
	}
	switch eventMsg.Action {
	case ActionApiReload:
		err = l.handleApiReload(ctx, eventMsg)
	case ActionRoleReload:
		err = l.handleRoleReload(ctx, eventMsg)
	case ActionResourceReload:
		err = l.handleResourceReload(ctx, eventMsg)
	case ActionAllPolicyReload:
		err = l.handleAllPolicyReload(ctx, eventMsg)
	case ActionPolicyRebuild:
		err = l.handlePolicyRebuild(ctx, eventMsg)
	default:
		err = errorx.Errorf("invalid action %v", eventMsg.Action)
	}

	return err
}

// handlePolicyRebuild recomputes policies for the message's scope on this
// replica. All replicas rebuild independently; idempotent policy ops converge.
func (l *Auth) handlePolicyRebuild(ctx context.Context, msg Msg) (err error) {
	switch msg.Scope {
	case ScopeOrganization:
		if msg.OrganizationId <= 0 {
			return nil
		}
		return l.RebuildPolicyForOrganization(ctx, msg.OrganizationId, msg.Type)
	case ScopeGlobal:
		return l.RebuildPolicy(ctx, msg.Type)
	default:
		return nil
	}
}

// handleApiReload
func (l *Auth) handleApiReload(ctx context.Context, msg Msg) (err error) {
	_, err = l.LoadAuthenticationWhitelist(ctx, msg.Type)
	if err != nil {
		return
	}
	_, err = l.LoadAuthorizationWhitelist(ctx, msg.Type)
	if err != nil {
		return
	}
	_, err = l.LoadApiResourceIds(ctx, msg.Type)
	if err != nil {
		return
	}
	return nil
}

// handleRoleReload
func (l *Auth) handleRoleReload(ctx context.Context, msg Msg) (err error) {
	switch msg.Scope {
	case ScopeOrganization:
		_, err = l.loadOrganizationRoles(ctx, msg.OrganizationId, msg.Type)
		if err != nil {
			return
		}
		// Invalidate plan resource cache so RebuildPolicy picks up changes
		l.organizationPlanResourceMap.Delete(msg.OrganizationId)
	case ScopeOrganizationUser:
		_, err = l.loadOrganizationUserRoles(ctx, msg.OrganizationUserId, msg.OrganizationId, msg.Type)
		if err != nil {
			return
		}
	}

	if len(msg.AddedPolicies) > 0 {
		err = l.enforcer.AddPolicies(msg.AddedPolicies)
		if err != nil {
			return
		}
	}
	if len(msg.RemovedPolicies) > 0 {
		err = l.enforcer.RemovePolicies(msg.RemovedPolicies)
		if err != nil {
			return
		}
	}
	return nil
}

// handleResourceReload
func (l *Auth) handleResourceReload(ctx context.Context, msg Msg) (err error) {
	_, err = l.LoadApiResourceIds(ctx, msg.Type)
	if err != nil {
		return
	}
	switch msg.Scope {
	case ScopeGlobal:
		_, err = l.LoadAuthorizationWhitelist(ctx, msg.Type)
		if err != nil {
			return
		}
		// Clear all organization plan resource caches — plan/resource changes
		// affect all subscribed organizations
		l.organizationPlanResourceMap.Clear()
	}

	if len(msg.AddedPolicies) > 0 {
		err = l.enforcer.AddPolicies(msg.AddedPolicies)
		if err != nil {
			return
		}
	}
	if len(msg.RemovedPolicies) > 0 {
		err = l.enforcer.RemovePolicies(msg.RemovedPolicies)
		if err != nil {
			return
		}
	}
	return nil
}

// handleAllPolicyReload
func (l *Auth) handleAllPolicyReload(ctx context.Context, msg Msg) (err error) {
	return l.enforcer.LoadDefaultFilteredPolicy()
}

// AuthCheck performs authentication and authorization checks for the given request.
// It returns a response containing the status of authentication, authorization,
// and the applicable data scope.
func (l *Auth) AuthCheck(ctx context.Context, input dto.AuthCheckRequestDto) (result *dto.AuthCheckResponseDto, err error) {
	// Initialize result with default values requiring both authentication and authorization
	result = &dto.AuthCheckResponseDto{
		AuthenticationRequired: true,
		Authenticated:          false,
		AuthorizationRequired:  true,
		Authorized:             false,
		DataScope:              constants.DataScopeAll,
	}
	// skip check admin
	// isAdmin, err := l.IsAdmin(ctx, input.Type)
	// if err != nil {
	// 	return
	// }
	// if isAdmin && input.Type == constants.UserTypeAdmin {
	// 	result.Authenticated = true
	// 	result.Authorized = true
	// 	result.DataScope = constants.DataScopeAll
	// 	return
	// }

	// Check authentication status and requirements
	ok, isRequired, err := l.CheckAuthentication(ctx, input)
	if err != nil {
		return
	}
	result.Authenticated = ok
	result.AuthenticationRequired = isRequired
	if !isRequired {
		// If authentication is not required, skip authorization
		result.AuthenticationRequired = false
		result.AuthorizationRequired = false
		return
	}
	// If authentication failed, return
	if !ok {
		return
	}

	// Check authorization status and requirements
	ok, isRequired, err = l.CheckAuthorization(ctx, input)
	if err != nil {
		return
	}
	result.AuthorizationRequired = isRequired
	result.Authorized = ok

	// If authorization is not required, return
	if !isRequired {
		return
	}

	// If authorization failed, return
	if !ok {
		return
	}

	// Get data scope for authorized user
	var dataScope constants.DataScope
	dataScope, err = l.GetDataScope(ctx, input)
	if err != nil {
		return
	}
	result.DataScope = dataScope

	return
}

// CheckAuthentication verifies if the incoming request is authenticated or requires authentication.
// It checks against a whitelist of routes that don't require authentication.
func (l *Auth) CheckAuthentication(ctx context.Context, input dto.AuthCheckRequestDto) (authenticated bool, authenticationRequired bool, err error) {
	authenticationRequired = true
	var authenticationWhitelist []string
	authenticationWhitelist, err = l.getAuthenticationWhitelist(ctx, input.Type)
	if err != nil {
		return
	}
	inWhitelist := slice.Contain(authenticationWhitelist, fmt.Sprintf(routeFormat, input.Method, input.Path))
	if inWhitelist {
		authenticationRequired = false
		return
	}
	_, err = l.GetUserClaimFromCtx(ctx, input.Type)
	if err == nil {
		authenticated = true
	}
	return
}

// CheckAuthorization checks if a user is authorized to access a specific resource.
// It first checks if the request path is in the authorization whitelist.
// If not, it verifies the user's roles and permissions using an enforcer.
func (l *Auth) CheckAuthorization(ctx context.Context, input dto.AuthCheckRequestDto) (authorized bool, authorizationRequired bool, err error) {
	claim, err := l.GetUserClaimFromCtx(ctx, input.Type)
	if err != nil {
		return
	}
	authorizationRequired = true
	var authorizationWhitelist []string
	authorizationWhitelist, err = l.getAuthorizationWhitelist(ctx, input.Type)
	if err != nil {
		return
	}
	inWhitelist := slice.Contain(authorizationWhitelist, fmt.Sprintf(routeFormat, input.Method, input.Path))
	if inWhitelist {
		authorizationRequired = false
		return
	}
	var roles []AuthRole
	roles, err = l.getOrganizationUserRoles(ctx, claim.OrganizationUserId, claim.OrganizationId, input.Type)
	if err != nil {
		return
	}
	if len(roles) == 0 {
		return
	}

	route := stringx.Format(routeFormat, input.Method, input.Path)
	resourceIds, err := l.getResourceIdsByRoute(ctx, route, input.Type)
	if err != nil {
		return
	}
	if len(resourceIds) == 0 {
		return
	}

	// Gate resources by the organization's active pricing plan: only resources
	// covered by the plan are eligible, regardless of role.
	var planResources *mapx.Map[int64, bool]
	planResources, err = l.getOrganizationPlanResourceIds(ctx, claim.OrganizationId)
	if err != nil {
		return
	}
	planResourceIds := slice.Filter(resourceIds, func(_ int, id int64) bool {
		inPlan, _ := planResources.Get(id)
		return inPlan
	})
	if len(planResourceIds) == 0 {
		return
	}

	// Admins are authorized for any resource within their plan scope.
	isAdmin := slice.Some(roles, func(_ int, item AuthRole) bool {
		return item.IsAdmin
	})
	if isAdmin {
		authorized = true
		return
	}

	// Non-admins: a role-resource policy must exist for a plan-covered resource.
	authorized = slice.Some(roles, func(_ int, item AuthRole) bool {
		return slice.Some(planResourceIds, func(_ int, resourceId int64) bool {
			policy := BuildInterfacePolicy(item.RoleId, resourceId, input.Type)
			ok, _ := l.enforcer.Enforce(policy...)
			return ok
		})
	})

	return
}

// GetDataScope returns the most restrictive data scope across the user's
// applicable role-resource combinations (lower Priority = more restrictive).
func (l *Auth) GetDataScope(ctx context.Context, input dto.AuthCheckRequestDto) (dataScope constants.DataScope, err error) {
	dataScope = constants.DataScopeOwn

	isAdmin, err := l.IsAdmin(ctx, input.Type)
	if err != nil {
		return
	}
	if isAdmin {
		dataScope = constants.DataScopeAll
		return
	}

	// Extract JWT user claims from the context
	var claim utils.JwtUser
	claim, err = l.GetUserClaimFromCtx(ctx, input.Type)
	if err != nil {
		return
	}

	// Get user roles for the specific organization and resource type
	userRoles, err := l.getOrganizationUserRoles(ctx, claim.OrganizationUserId, claim.OrganizationId, input.Type)
	if err != nil || len(userRoles) == 0 {
		return
	}

	// Retrieve all organization role resources for the given resource type
	organizationRoleResources, err := l.getOrganizationRoleResources(ctx, claim.OrganizationId, input.Type)
	if err != nil || organizationRoleResources.Size() == 0 {
		return
	}

	// Get resource IDs matching the specific route (method and path)
	resourceIds, err := l.getResourceIdsByRoute(ctx, stringx.Format(routeFormat, input.Method, input.Path), input.Type)
	if err != nil || len(resourceIds) == 0 {
		return
	}

	// Iterate through user roles and resource IDs to find the most restrictive data scope
	slice.ForEach(userRoles, func(index int, role AuthRole) {
		slice.ForEach(resourceIds, func(index int, resourceId int64) {
			// Create a key to look up role-resource permissions
			key := stringx.Format(roleResourceFormat, role.RoleId, resourceId)
			roleResource, ok := organizationRoleResources.Get(key)
			if !ok {
				return
			}
			// Update data scope if a more restrictive one is found
			if roleResource.DataScope.Priority() < dataScope.Priority() {
				dataScope = roleResource.DataScope
			}
		})
	})

	return
}

// IsAdmin checks if the current user is an admin.
func (l *Auth) IsAdmin(ctx context.Context, userType constants.UserType) (isAdmin bool, err error) {
	claim, err := l.GetUserClaimFromCtx(ctx, userType)
	if err != nil {
		return false, nil
	}

	if userType != claim.UserType {
		return
	}
	organizationUserId := claim.OrganizationUserId
	organizationId := claim.OrganizationId
	var roles []AuthRole
	roles, err = l.getOrganizationUserRoles(ctx, organizationUserId, organizationId, userType)
	if err != nil {
		return
	}
	isAdmin = slice.Some(roles, func(index int, item AuthRole) bool {
		return item.IsAdmin
	})
	return
}

// LoadAuthenticationWhitelist load authentication whitelist for a given user type.
func (l *Auth) LoadAuthenticationWhitelist(ctx context.Context, userType constants.UserType) (whitelist []string, err error) {
	// get apis from db
	result, err, _ := l.sf.Do(stringx.Format(authenticationWhitelistFormat, userType), func() (interface{}, error) {
		apis, err := l.repos.ApiRepo.GetAuthenticationWhitelist(ctx, userType)
		if err != nil {
			return nil, err
		}
		results := slice.FilterMap(apis, func(index int, item *model.Api) (string, bool) {
			if len(item.Path) == 0 || len(item.Method) == 0 {
				return "", false
			}
			return stringx.Format(routeFormat, item.Method, item.Path), true
		})
		return results, nil
	})
	if err != nil {
		return nil, err
	}
	if result == nil {
		whitelist = []string{}
	} else {
		var ok bool
		whitelist, ok = result.([]string)
		if !ok {
			return nil, errorx.New("invalid authentication whitelist")
		}
	}
	l.authenticationWhitelistMap.Set(userType, whitelist)

	return whitelist, nil
}

// LoadAuthorizationWhitelist load authorization whitelist for a given user type.
func (l *Auth) LoadAuthorizationWhitelist(ctx context.Context, userType constants.UserType) (whitelist []string, err error) {
	result, err, _ := l.sf.Do(stringx.Format(authorizationWhitelistFormat, userType), func() (interface{}, error) {
		apis, err := l.repos.ApiRepo.GetAuthorizationWhitelist(ctx, userType)
		if err != nil {
			return nil, err
		}
		results := slice.FilterMap(apis, func(index int, item *model.Api) (string, bool) {
			if len(item.Path) == 0 || len(item.Method) == 0 {
				return "", false
			}
			return stringx.Format(routeFormat, item.Method, item.Path), true
		})
		// public resource
		publicResource, err := l.repos.ResourceRepo.PublicResources(ctx, userType, true)
		if err != nil {
			return nil, err
		}
		slice.ForEach(publicResource, func(_ int, resource *model.Resource) {
			publicResults := slice.FilterMap(resource.Apis, func(_ int, item *model.Api) (string, bool) {
				if len(item.Path) == 0 || len(item.Method) == 0 {
					return "", false
				}
				return stringx.Format(routeFormat, item.Method, item.Path), true
			})
			if len(publicResults) > 0 {
				results = append(results, publicResults...)
			}
		})
		return results, nil
	})
	if err != nil {
		return nil, err
	}
	if result == nil {
		whitelist = []string{}
	} else {
		var ok bool
		whitelist, ok = result.([]string)
		if !ok {
			return nil, errorx.New("invalid authorization whitelist")
		}
	}
	l.authorizationWhitelistMap.Set(userType, whitelist)
	return whitelist, nil
}

// RebuildPolicy rebuilds policies for ALL orgs of the given user type. Runs at
// service startup (see app/*/rpc/internal/logic/base_logic.go); for runtime
// changes use RebuildPolicyForOrganization or publish an ActionPolicyRebuild.
// Recomputes from resource↔role bindings gated by each org's active plan, diffs
// against the enforcer, and (full path only) self-publishes ActionAllPolicyReload.
func (l *Auth) RebuildPolicy(ctx context.Context, organizationType constants.UserType) (err error) {
	return l.rebuildPolicy(ctx, organizationType, 0)
}

// RebuildPolicyForOrganization recomputes a single organization's policies for
// the given user type. It only touches policies whose subject ("Role:<id>")
// belongs to one of this org's roles; other orgs' policies are never modified.
// Triggered by an ActionPolicyRebuild ScopeOrganization event on every replica.
func (l *Auth) RebuildPolicyForOrganization(ctx context.Context, organizationId int64, organizationType constants.UserType) (err error) {
	if organizationId <= 0 {
		return nil
	}
	return l.rebuildPolicy(ctx, organizationType, organizationId)
}

// rebuildPolicy core. orgId == 0: full rebuild; orgId > 0: org-scoped.
// singleflight dedupes concurrent same-scope rebuilds. Distinct key prefixes
// avoid collision; idempotent policy ops let full + org-scoped rebuilds converge.
func (l *Auth) rebuildPolicy(ctx context.Context, organizationType constants.UserType, organizationId int64) (err error) {
	sfKey := stringx.Format(policyRebuildFormat, organizationType)
	if organizationId > 0 {
		sfKey = stringx.Format(organizationPolicyRebuildFormat, organizationId, organizationType)
	}
	_, err, _ = l.sf.Do(sfKey, func() (interface{}, error) {
		return nil, l.doRebuild(ctx, organizationType, organizationId)
	})
	return err
}

func (l *Auth) doRebuild(ctx context.Context, organizationType constants.UserType, organizationId int64) (err error) {
	buildKeyFn := func(policy []string) string {
		return strings.Join(policy, ",")
	}

	// 1. invalidate the plan-resource cache for this scope so the rebuild reads
	//    the latest subscription/plan state from the DB.
	if organizationId > 0 {
		l.organizationPlanResourceMap.Delete(organizationId)
	} else {
		l.organizationPlanResourceMap.Clear()
	}

	// 2. org-scoped rebuild: fetch this org's role-id set to scope which existing
	//    policies may be touched and to clean orphan policies.
	var subjectSet map[string]int64
	if organizationId > 0 {
		scopeRoles, scopeErr := l.repos.RoleRepo.All(ctx, organizationId, organizationType)
		if scopeErr != nil {
			return scopeErr
		}
		subjectSet = map[string]int64{}
		slice.ForEach(scopeRoles, func(_ int, role *model.Role) {
			if role == nil {
				return
			}
			subjectSet[stringx.Format(roleFormat, role.Id)] = role.Id
		})
	}

	// 3. resources carrying the scope's roles. AllResourcesForAuth preloads Roles
	//    filtered to the scope (orgId>0 ⇒ only this org's active roles).
	resources, err := l.repos.ResourceRepo.AllResourcesForAuth(ctx, organizationId, organizationType)
	if err != nil {
		return err
	}

	// 4. prime the plan-resource cache for the orgs we touch.
	if organizationId > 0 {
		if _, err = l.getOrganizationPlanResourceIds(ctx, organizationId); err != nil {
			return err
		}
	} else {
		orgIds := make(map[int64]bool)
		slice.ForEach(resources, func(_ int, resource *model.Resource) {
			if resource == nil {
				return
			}
			slice.ForEach(resource.Roles, func(_ int, role *model.Role) {
				if role != nil && role.OrganizationId > 0 {
					orgIds[role.OrganizationId] = true
				}
			})
		})
		for orgId := range orgIds {
			if _, err = l.getOrganizationPlanResourceIds(ctx, orgId); err != nil {
				return err
			}
		}
	}

	// 5. desired policy set, gated by plan.
	desired := map[string][]string{}
	slice.ForEach(resources, func(_ int, resource *model.Resource) {
		if resource == nil {
			return
		}
		slice.ForEach(resource.Roles, func(_ int, role *model.Role) {
			if role == nil || role.Status != constants.CommonStatusActive {
				return
			}
			if role.OrganizationId > 0 {
				planResources, _ := l.organizationPlanResourceMap.Get(role.OrganizationId)
				// deny default: no plan entry ⇒ skip (no policy for this pair).
				if planResources == nil {
					return
				}
				if _, ok := planResources.Get(resource.Id); !ok {
					return
				}
			}
			policy := BuildPolicy(role.Id, resource.Id, organizationType)
			desired[buildKeyFn(policy)] = policy
		})
	})

	// 6. existing policies. Org-scoped: only those whose subject belongs to this
	//    org (never touch another org's policies). Full: all current policies,
	//    matching the legacy diff.
	oldPolicies, err := l.enforcer.AllPolicies()
	if err != nil {
		return err
	}
	existing := map[string][]string{}
	slice.ForEach(oldPolicies, func(_ int, item []string) {
		if len(item) == 0 {
			return
		}
		if organizationId > 0 {
			if _, ok := subjectSet[item[0]]; !ok {
				return
			}
		}
		existing[buildKeyFn(item)] = item
	})

	// 7. diff.
	var newPolicies = make([][]string, 0, len(desired))
	var removePolicies = make([][]string, 0, len(existing))
	for key, policy := range desired {
		if _, ok := existing[key]; !ok {
			newPolicies = append(newPolicies, policy)
		}
	}
	for key, policy := range existing {
		if _, ok := desired[key]; !ok {
			removePolicies = append(removePolicies, policy)
		}
	}

	// 8. apply (idempotent; autoSave persists; redis watcher fans out).
	eg := errgroup.Group{}
	eg.Go(func() error {
		if len(removePolicies) > 0 {
			return l.enforcer.RemovePolicies(removePolicies)
		}
		return nil
	})
	eg.Go(func() error {
		if len(newPolicies) > 0 {
			return l.enforcer.AddPolicies(newPolicies)
		}
		return nil
	})
	if err = eg.Wait(); err != nil {
		return err
	}

	// 9. only the full-rebuild path self-publishes ActionAllPolicyReload (legacy
	//    startup fan-out). Org-scoped rebuilds don't: every replica rebuilt itself
	//    from the originating event.
	if organizationId == 0 && (len(newPolicies) > 0 || len(removePolicies) > 0) {
		msg := Msg{
			Action: ActionAllPolicyReload,
			Scope:  ScopeGlobal,
			Type:   organizationType,
		}
		if perr := l.redisCli.GetClient().Publish(ctx, l.eventTopic, msg.String()).Err(); perr != nil {
			return errorx.WithStack(perr)
		}
	}

	return nil
}

func (l *Auth) GetUserClaimFromCtx(ctx context.Context, userType constants.UserType) (claim utils.JwtUser, err error) {
	// Get from metadata
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		err = status.Error(codes.Unauthenticated, "not login")
		return claim, errorx.WithStack(err)
	}
	tokens := md[utils.JwtTokenCtxKey]
	if len(tokens) == 0 {
		err = status.Error(codes.Unauthenticated, "not login")
		return claim, errorx.WithStack(err)
	}
	token := tokens[0]
	jwtConf := l.jwtConfig
	jwtClaim, err := utils.ParseToken(token, jwtConf.AccessSecret)
	if err != nil {
		err = status.Error(codes.InvalidArgument, err.Error())
		return claim, errorx.WithStack(err)
	}
	return jwtClaim.JwtUser, nil
}

// getAuthenticationWhitelist get authentication whitelist
func (l *Auth) getAuthenticationWhitelist(ctx context.Context, userType constants.UserType) (whitelist []string, err error) {
	var ok bool
	whitelist, ok = l.authenticationWhitelistMap.Get(userType)
	if ok {
		return whitelist, nil
	}

	whitelist, err = l.LoadAuthenticationWhitelist(ctx, userType)
	return whitelist, err
}

// getAuthorizationWhitelist get authorization whitelist
func (l *Auth) getAuthorizationWhitelist(ctx context.Context, userType constants.UserType) (whitelist []string, err error) {
	var ok bool
	whitelist, ok = l.authorizationWhitelistMap.Get(userType)
	if ok {
		return whitelist, nil
	}
	whitelist, err = l.LoadAuthorizationWhitelist(ctx, userType)
	return whitelist, err
}

// getOrganizationUserRoles get organization user roles
func (l *Auth) getOrganizationUserRoles(ctx context.Context, organizationUserId int64, organizationId int64, userType constants.UserType) (roles []AuthRole, err error) {
	key := stringx.Format(organizationUserRoleFormat, userType, organizationId, organizationUserId)
	content, _ := l.redisCli.Get(ctx, key)
	if len(content) == 0 {
		roles, err = l.loadOrganizationUserRoles(ctx, organizationUserId, organizationId, userType)
		if err != nil {
			return nil, err
		}
		return roles, err
	}

	err = json.Unmarshal([]byte(content), &roles)
	return roles, errorx.WithStack(err)
}

// loadOrganizationUserRoles load roles from db
func (l *Auth) loadOrganizationUserRoles(ctx context.Context, organizationUserId int64, organizationId int64, userType constants.UserType) (roles []AuthRole, err error) {
	roles = []AuthRole{}
	organizationUser, _ := l.repos.OrganizationUserRepo.GetOrganizationUserById(ctx, organizationUserId, organizationId)
	if organizationUser == nil || organizationUser.Status != constants.OrganizationUserStatusActive || len(organizationUser.Roles) == 0 {
		goto saveCache
	}

	roles = slice.FilterMap(organizationUser.Roles, func(index int, item *model.Role) (AuthRole, bool) {
		if item == nil {
			return AuthRole{}, false
		}
		if item.Status != constants.CommonStatusActive {
			return AuthRole{}, false
		}
		return AuthRole{RoleId: item.Id, IsAdmin: item.Name == constants.AdminRoleName}, true
	})

saveCache:
	key := stringx.Format(organizationUserRoleFormat, userType, organizationId, organizationUserId)
	rolesContent, err := json.Marshal(roles)
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	content := string(rolesContent)

	err = l.redisCli.SetExpire(ctx, key, content, organizationUserRoleExpireSeconds)

	return roles, errorx.WithStack(err)
}

// loadOrganizationRoles load organization roles
func (l *Auth) loadOrganizationRoles(ctx context.Context, organizationId int64, userType constants.UserType) (roles []AuthRole, err error) {
	organizationUsers, _ := l.repos.OrganizationUserRepo.AllOrganizationUsers(ctx, organizationId)

	slice.ForEach(organizationUsers, func(index int, organizationUser *model.OrganizationUser) {
		if organizationUser == nil || organizationUser.Status != constants.OrganizationUserStatusActive {
			return
		}
		subRoles := slice.FilterMap(organizationUser.Roles, func(index int, role *model.Role) (AuthRole, bool) {
			if role == nil {
				return AuthRole{}, false
			}
			if role.Status != constants.CommonStatusActive {
				return AuthRole{}, false
			}
			return AuthRole{RoleId: role.Id, IsAdmin: role.Name == constants.AdminRoleName}, true
		})
		key := stringx.Format(organizationUserRoleFormat, userType, organizationId, organizationUser.Id)
		roles = append(roles, subRoles...)
		var rolesContent []byte
		rolesContent, err = json.Marshal(subRoles)
		if err != nil {
			logx.WithContext(ctx).WithFields(logx.Field("stack", errorx.GetErrorStack(err))).Error(err)
			return
		}
		content := string(rolesContent)

		_ = l.redisCli.SetExpire(ctx, key, content, organizationUserRoleExpireSeconds)
	})

	return
}

// getResourceIdsByRoute get resource ids by route
func (l *Auth) getResourceIdsByRoute(ctx context.Context, route string, userType constants.UserType) (resourceIds []int64, err error) {
	organizationResources, ok := l.apiResourceMap.Get(userType)
	if !ok {
		organizationResources, err = l.LoadApiResourceIds(ctx, userType)
		if err != nil {
			return nil, err
		}
	}

	resourceIds, ok = organizationResources.Get(route)
	return resourceIds, nil
}

// LoadApiResourceIds load api resource ids
func (l *Auth) LoadApiResourceIds(ctx context.Context, userType constants.UserType) (resources *mapx.Map[string, []int64], err error) {
	apis, err := l.repos.ApiRepo.GetAllApis(ctx, userType, true)
	if err != nil {
		return nil, err
	}
	apis = slice.Filter(apis, func(index int, item *model.Api) bool {
		return len(item.Path) > 0 && len(item.Method) > 0
	})
	apiResourceMap := mapx.NewMap[string, []int64]()
	for _, api := range apis {
		resourceIds := slice.FilterMap(api.Resources, func(index int, item *model.Resource) (int64, bool) {
			if item.Status != constants.CommonStatusActive {
				return 0, false
			}
			return api.Resources[index].Id, true
		})
		apiResourceMap.Set(stringx.Format(routeFormat, api.Method, api.Path), resourceIds)
	}

	l.apiResourceMap.Set(userType, apiResourceMap)
	return apiResourceMap, nil
}

// getOrganizationRoleResources
func (l *Auth) getOrganizationRoleResources(ctx context.Context, organizationId int64, userType constants.UserType) (roleResources *mapx.Map[string, RoleResource], err error) {
	var ok bool
	roleResources, ok = l.organizationRoleResourceMap.Get(organizationId)
	if !ok {
		roleResources, err = l.LoadOrganizationRoleResources(ctx, organizationId, userType)
	}
	return roleResources, err
}

// LoadOrganizationRoleResources load organization role resources
func (l *Auth) LoadOrganizationRoleResources(ctx context.Context, organizationId int64, userType constants.UserType) (resources *mapx.Map[string, RoleResource], err error) {
	roleResources, err := l.repos.RoleRepo.GetOrganizationRoleWithResources(ctx, organizationId, userType)
	if err != nil {
		return
	}
	resources = mapx.NewMap[string, RoleResource]()
	slice.ForEach(roleResources, func(index int, item repo.RoleWithResource) {
		key := stringx.Format(roleResourceFormat, item.Id, item.ResourceId)
		resources.Set(key, RoleResource{
			RoleId:     item.Id,
			ResourceId: item.ResourceId,
			DataScope:  item.DataScope,
		})
	})
	l.organizationRoleResourceMap.Set(organizationId, resources)
	return
}

// getOrganizationPlanResourceIds returns the set of resource IDs allowed by the organization's
// active pricing plan subscriptions. Uses in-memory cache, loading from DB on miss.
// singleflight collapses concurrent cache-miss loads for the same org into one DB hit.
// Returns nil if the organization has no active subscription, meaning no resources are accessible.
func (l *Auth) getOrganizationPlanResourceIds(ctx context.Context, organizationId int64) (*mapx.Map[int64, bool], error) {
	cached, ok := l.organizationPlanResourceMap.Get(organizationId)
	if ok {
		return cached, nil
	}
	result, err, _ := l.sf.Do(stringx.Format(organizationPlanResourceFormat, organizationId), func() (interface{}, error) {
		// re-check cache inside the singleflight call to avoid a redundant
		// DB load when another in-flight call already populated it
		if cached, ok = l.organizationPlanResourceMap.Get(organizationId); ok {
			return cached, nil
		}
		return l.loadOrganizationPlanResources(ctx, organizationId)
	})
	if err != nil {
		return nil, err
	}
	resourceSet, ok := result.(*mapx.Map[int64, bool])
	if !ok {
		return nil, errorx.New("invalid organization plan resource set")
	}
	return resourceSet, nil
}

// loadOrganizationPlanResources loads plan resources for an organization from DB into cache.
func (l *Auth) loadOrganizationPlanResources(ctx context.Context, organizationId int64) (*mapx.Map[int64, bool], error) {
	activeSubs, err := l.repos.OrganizationSubscriptionRepo.FindActiveByOrganizationId(ctx, organizationId)
	if err != nil {
		return nil, err
	}

	resourceSet := mapx.NewMap[int64, bool]()
	for _, activeSub := range activeSubs {
		if activeSub != nil &&
			activeSub.Status.IsEffective() &&
			!isSubscriptionExpired(activeSub, time.Now()) {

			plan, _ := l.repos.PricingPlanRepo.GetOne(ctx, activeSub.PricingPlanId, nil)
			if plan != nil && plan.Status == constants.CommonStatusActive {
				ids, innerErr := l.repos.PricingPlanAndResourceRepo.FindResourceIdsByPlanId(ctx, activeSub.PricingPlanId)
				if innerErr != nil {
					return nil, innerErr
				}
				for _, id := range ids {
					resourceSet.Set(id, true)
				}
			}
		}
	}

	l.organizationPlanResourceMap.Set(organizationId, resourceSet)
	return resourceSet, nil
}

// isSubscriptionExpired treats expire date as valid until the last second of that day
func isSubscriptionExpired(sub *model.OrganizationSubscription, now time.Time) bool {
	if sub.ExpireDate == nil || !sub.ExpireDate.Valid {
		return false
	}
	endOfDay := sub.ExpireDate.Time.Add(24*time.Hour - time.Second)
	return !endOfDay.After(now)
}

// Close close auth, save policy and unsubscribe from redis
func (l *Auth) Close() {
	if l.enforcer != nil {
		_ = l.enforcer.SavePolicy()
	}

	if l.eventSubScribeCancelFn != nil {
		l.eventSubScribeCancelFn()
	}
}
