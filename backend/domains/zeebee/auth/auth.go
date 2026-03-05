package auth

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/duke-git/lancet/v2/slice"
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
}

func NewAuth(repos *repo.AllRepos, redisCli *redisx.Redisx, jwtConfig dto.JwtConfig, enforcer *rbac.Enforcer) *Auth {
	return &Auth{
		repos:                       repos,
		redisCli:                    redisCli,
		jwtConfig:                   jwtConfig,
		enforcer:                    enforcer,
		authenticationWhitelistMap:  mapx.NewMap[constants.UserType, []string](),
		authorizationWhitelistMap:   mapx.NewMap[constants.UserType, []string](),
		apiResourceMap:              mapx.NewMap[constants.UserType, *mapx.Map[string, []int64]](),
		organizationRoleResourceMap: mapx.NewMap[int64, *mapx.Map[string, RoleResource]](),
	}
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

// GetDataScope gets the data scope for a user.
func (l *Auth) GetDataScope(ctx context.Context, input dto.AuthCheckRequestDto) (dataScope constants.DataScope, err error) {
	scopes := map[constants.DataScope]int{
		constants.DataScopeAll: 0,
		constants.DataScopeOwn: 1,
	}
	dataScope = constants.DataScopeOwn

	var claim utils.JwtUser
	claim, err = l.GetUserClaimFromCtx(ctx)
	if err != nil {
		return
	}
	userRoles, err := l.getOrganizationUserRoles(ctx, claim.OrganizationUserId, claim.OrganizationId, input.Type)
	if err != nil || len(userRoles) == 0 {
		return
	}

	organizationRoleResources, err := l.getOrganizationRoleResources(ctx, claim.OrganizationId, input.Type)
	if err != nil || organizationRoleResources.Size() == 0 {
		return
	}

	resourceIds, err := l.getResourceIdsByRoute(ctx, stringx.Format(routeFormat, input.Method, input.Path), input.Type)
	if err != nil || len(resourceIds) == 0 {
		return
	}
	slice.ForEach(userRoles, func(index int, role AuthRole) {
		slice.ForEach(resourceIds, func(index int, resourceId int64) {
			key := stringx.Format(roleResourceFormat, role.RoleId, resourceId)
			roleResource, ok := organizationRoleResources.Get(key)
			if !ok {
				return
			}
			_, ok = scopes[roleResource.DataScope]
			if !ok {
				return
			}

			if scopes[roleResource.DataScope] < scopes[dataScope] {
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
		whitelist = result.([]string)
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
		whitelist = result.([]string)
	}
	l.authorizationWhitelistMap.Set(userType, whitelist)
	return whitelist, nil
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
		roles, err = l.loadRoles(ctx, organizationUserId, organizationId, userType)
		if err != nil {
			return nil, err
		}
		return roles, err
	}

	err = json.Unmarshal([]byte(content), &roles)
	return roles, errorx.WithStack(err)
}

// loadRoles load roles from db
func (l *Auth) loadRoles(ctx context.Context, organizationUserId int64, organizationId int64, userType constants.UserType) (roles []AuthRole, err error) {
	organizationUser, _ := l.repos.OrganizationUserRepo.GetOrganizationUserById(ctx, organizationUserId, organizationId)
	if organizationUser == nil || organizationUser.Roles == nil {
		return []AuthRole{}, nil
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

	rolesContent, err := json.Marshal(roles)
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	content := string(rolesContent)
	key := stringx.Format(organizationUserRoleFormat, userType, organizationId, organizationUserId)
	err = l.redisCli.SetExpire(ctx, key, content, organizationUserRoleExpireSeconds)

	return roles, errorx.WithStack(err)
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
