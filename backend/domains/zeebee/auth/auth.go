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

type Auth struct {
	repos                       *repo.AllRepos
	redisCli                    *redisx.Redisx
	enforcer                    *rbac.Enforcer
	jwtConfig                   dto.JwtConfig
	authenticationWhitelistMap  *mapx.Map[constants.UserType, []string]
	authorizationWhitelistMap   *mapx.Map[constants.UserType, []string]
	apiResourceMap              *mapx.Map[constants.UserType, *mapx.Map[string, []int64]]
	organizationRoleResourceMap *mapx.Map[int64, *mapx.Map[string, RoleResource]]
	sf                          singleflight.Group
	eventTopic                  string
	eventSubScribeCancelFn      context.CancelFunc
}

func NewAuth(repos *repo.AllRepos, redisCli *redisx.Redisx, jwtConfig dto.JwtConfig, enforcer *rbac.Enforcer, authEventTopic string) *Auth {
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
		eventTopic:                  authEventTopic,
		eventSubScribeCancelFn:      cancel,
	}

	go func() {
		auth.subscribeEvent(ctx)
	}()

	return auth
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
						logx.Field("error", err.Error()),
					).Errorf("handleAuthMsg error")
				}

			}
		}
		_ = subPub.Close()
		time.Sleep(retryInterval)
	}
}

// handleAuthMsg
func (l *Auth) handleAuthMsg(ctx context.Context, msg *redis.Message) (err error) {
	var eventMsg Msg
	if msg == nil {
		err = errorx.New("Auth message is nil")
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
	default:
		err = errorx.Errorf("Invalid action %v", eventMsg.Action)
	}

	return err

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

	// Check authentication status and requirements
	ok, isRequired := l.CheckAuthentication(ctx, input)
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

	// If user is admin, grant all permissions and full data scope
	if l.IsAdmin(ctx, input.Type) {
		result.Authenticated = true
		result.Authorized = true
		result.DataScope = constants.DataScopeAll
		return
	}

	// Check authorization status and requirements
	ok, isRequired = l.CheckAuthorization(ctx, input)
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
func (l *Auth) CheckAuthentication(ctx context.Context, input dto.AuthCheckRequestDto) (authenticated bool, authenticationRequired bool) {
	authenticationRequired = true
	authenticationWhitelist, err := l.getAuthenticationWhitelist(ctx, input.Type)
	if err != nil {
		return
	}
	inWhitelist := slice.Contain(authenticationWhitelist, fmt.Sprintf(routeFormat, input.Method, input.Path))
	if inWhitelist {
		authenticationRequired = false
		return
	}
	_, err = l.GetUserClaimFromCtx(ctx)
	if err == nil {
		authenticated = true
	}
	return
}

// CheckAuthorization checks if a user is authorized to access a specific resource.
// It first checks if the request path is in the authorization whitelist.
// If not, it verifies the user's roles and permissions using an enforcer.
func (l *Auth) CheckAuthorization(ctx context.Context, input dto.AuthCheckRequestDto) (authorized bool, authorizationRequired bool) {
	claim, err := l.GetUserClaimFromCtx(ctx)
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
	roles, err := l.getOrganizationUserRoles(ctx, claim.OrganizationUserId, claim.OrganizationId, input.Type)
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
	authorized = slice.Some(roles, func(index int, item AuthRole) bool {
		return slice.Some(resourceIds, func(index int, subItem int64) bool {
			ok, _ := l.enforcer.Enforce(
				stringx.Format("Role:%d", item.RoleId),
				input.Type,
				stringx.Format("Resource:%d", subItem),
				"allow",
			)
			return ok
		})
	})

	return
}

// GetDataScope retrieves the data scope for a given user based on their roles and resource permissions.
// It takes a context and an AuthCheckRequestDto as input, and returns the most restrictive data scope
// among all applicable role-resource combinations. The function checks JWT claims, user roles,
// organization role resources, and resource IDs to determine the appropriate data scope.
// The data scope hierarchy is defined in the scopes map, where lower values indicate more restrictive scopes.
// Returns the determined data scope and any error encountered during the process.
func (l *Auth) GetDataScope(ctx context.Context, input dto.AuthCheckRequestDto) (dataScope constants.DataScope, err error) {
	dataScope = constants.DataScopeOwn

	// Extract JWT user claims from the context
	var claim utils.JwtUser
	claim, err = l.GetUserClaimFromCtx(ctx)
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
func (l *Auth) IsAdmin(ctx context.Context, userType constants.UserType) (isAdmin bool) {
	claim, err := l.GetUserClaimFromCtx(ctx)
	if err != nil {
		return false
	}

	if userType != claim.UserType {
		return false
	}
	organizationUserId := claim.OrganizationUserId
	organizationId := claim.OrganizationId
	roles, err := l.getOrganizationUserRoles(ctx, organizationUserId, organizationId, userType)
	if err != nil {
		return false
	}
	isAdmin = slice.Some(roles, func(index int, item AuthRole) bool {
		return item.IsAdmin
	})
	return isAdmin
}

// LoadAuthenticationWhitelist load authentication whitelist for a given user type.
func (l *Auth) LoadAuthenticationWhitelist(ctx context.Context, userType constants.UserType) (whitelist []string, err error) {

	// get apis from db
	result, err, _ := l.sf.Do(stringx.Format("AuthenticationWhitelist:%v", userType), func() (interface{}, error) {
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
	result, err, _ := l.sf.Do(stringx.Format("AuthorizationWhitelist:%v", userType), func() (interface{}, error) {
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

// RebuildPolicy rebuilds the authorization policies for a given user type.
// It compares the current policies with the ones derived from the resources in the database,
// identifies new policies to be added and obsolete policies to be removed.
// The function returns an error if any database or policy operation fails.
func (l *Auth) RebuildPolicy(ctx context.Context, organizationType constants.UserType) (err error) {
	buildKeyFn := func(policy []string) string {
		return strings.Join(policy, ",")
	}

	// Get existing policies from the enforcer
	oldPolicies, err := l.enforcer.AllPolicies()
	if err != nil {
		return err
	}
	// Convert old policies to a map for efficient lookup
	oldPoliciesMap := map[string]struct{}{}
	slice.ForEach(oldPolicies, func(index int, item []string) {
		key := buildKeyFn(item)
		oldPoliciesMap[key] = struct{}{}
	})

	// Get all resources for the given user type
	resources, err := l.repos.ResourceRepo.AllResourcesForAuth(ctx, organizationType)
	if err != nil {
		return err
	}
	// Initialize slices for tracking policy changes
	var currentPolicies [][]string
	var newPolicies [][]string
	var removePolicies [][]string
	currentPoliciesMap := map[string]struct{}{}

	// Process each resource and generate policies
	slice.ForEach(resources, func(index int, resource *model.Resource) {
		if resource == nil {
			return
		}
		// Generate policies for each active role associated with the resource
		rolePolicies := slice.FilterMap(resource.Roles, func(index int, role *model.Role) ([]string, bool) {
			if role == nil || role.Status != constants.CommonStatusActive {
				return nil, false
			}
			// Create a policy entry for the role-resource combination
			policy := BuildPolicy(role.Id, resource.Id, organizationType)

			key := buildKeyFn(policy)
			// Check if this is a new policy
			_, ok := oldPoliciesMap[key]
			if !ok {
				newPolicies = append(newPolicies, policy)
			}
			currentPoliciesMap[key] = struct{}{}
			return policy, true
		})
		if len(rolePolicies) == 0 {
			return
		}
		currentPolicies = append(currentPolicies, rolePolicies...)
	})

	// Identify policies that need to be removed
	slice.ForEach(oldPolicies, func(index int, item []string) {
		key := buildKeyFn(item)
		_, ok := currentPoliciesMap[key]
		if !ok {
			removePolicies = append(removePolicies, item)
		}
	})

	// Execute the policy operations in parallel
	eg := errgroup.Group{}
	eg.Go(func() error {
		var err error
		// Remove obsolete policies if any
		if len(removePolicies) > 0 {
			err = l.enforcer.RemovePolicies(removePolicies)
			if err != nil {
				return err
			}
		}
		return err
	})
	eg.Go(func() error {
		var err error
		// Add new policies if any
		if len(newPolicies) > 0 {
			err = l.enforcer.AddPolicies(newPolicies)
			if err != nil {
				return err
			}
		}
		return err
	})
	err = eg.Wait()
	if err != nil {
		return err
	}

	// check if policy is changed
	isPolicyChanged := len(newPolicies) > 0 || len(removePolicies) > 0
	if isPolicyChanged {
		// send event to reload policy
		msg := Msg{
			Action: ActionAllPolicyReload,
			Scope:  ScopeGlobal,
			Type:   organizationType,
		}
		err = l.redisCli.GetClient().Publish(ctx, l.eventTopic, msg.String()).Err()
		if err != nil {
			return errorx.WithStack(err)
		}
	}

	return nil
}

func (l *Auth) GetUserClaimFromCtx(ctx context.Context) (claim utils.JwtUser, err error) {
	// Get from metadata
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		err = status.Error(codes.Unauthenticated, "Not login")
		return claim, errorx.WithStack(err)
	}
	tokens := md[utils.JwtTokenCtxKey]
	if len(tokens) == 0 {
		err = status.Error(codes.Unauthenticated, "Not login")
		return claim, errorx.WithStack(err)
	}
	token := tokens[0]
	jwtClaim, err := utils.ParseToken(token, l.jwtConfig.AccessSecret)
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
			return item.Id, true

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

// Close close auth, save policy and unsubscribe from redis
func (l *Auth) Close() {
	if l.enforcer != nil {
		_ = l.enforcer.SavePolicy()
	}

	if l.eventSubScribeCancelFn != nil {
		l.eventSubScribeCancelFn()
	}
}
