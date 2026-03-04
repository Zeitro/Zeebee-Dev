package logic

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

var routeFormat = "%s:%s"
var organizationUserRoleFormat = "%s:organization:%d:organizationuser:%d:roles" // <type>:organization:<organizationId>:organizationuser:<organizationUserId>:roles
var organizationUserRoleExpireSeconds = 3600

type AuthRole struct {
	RoleId  int64 `json:"role_id"`
	IsAdmin bool  `json:"is_admin"`
}

type AuthResource struct {
	ResourceId     int64 `json:"resource_id"`
	OrganizationId int64 `json:"organization_id"`
}

type AuthLogic struct {
	repos                      *repo.AllRepos
	redisCli                   *redisx.Redisx
	enforcer                   *rbac.Enforcer
	jwtConfig                  JwtConfig
	authenticationWhitelistMap *mapx.Map[constants.UserType, []string]
	authorizationWhitelistMap  *mapx.Map[constants.UserType, []string]
	organizationApiResourceMap *mapx.Map[int64, *mapx.Map[string, []*model.Resource]]
}

func NewAuthLogic(repos *repo.AllRepos, redisCli *redisx.Redisx, jwtConfig JwtConfig, enforcer *rbac.Enforcer) *AuthLogic {
	return &AuthLogic{
		repos:                      repos,
		redisCli:                   redisCli,
		jwtConfig:                  jwtConfig,
		enforcer:                   enforcer,
		authenticationWhitelistMap: mapx.NewMap[constants.UserType, []string](),
		authorizationWhitelistMap:  mapx.NewMap[constants.UserType, []string](),
		organizationApiResourceMap: mapx.NewMap[int64, *mapx.Map[string, []*model.Resource]](),
	}
}

func (l *AuthLogic) AuthCheck(ctx context.Context, input *dto.AuthCheckRequestDto) (result *dto.AuthCheckResponseDto, err error) {

	result = &dto.AuthCheckResponseDto{
		AuthenticationRequired: true,
		Authenticated:          false,
		AuthorizationRequired:  true,
		Authorized:             false,
		DataScope:              constants.DataScopeAll,
	}
	// check  authentication
	authenticationWhitelist, err := l.getAuthenticationWhitelist(ctx, input.Type)
	if err != nil {
		return
	}
	inWhitelist := slice.Contain(authenticationWhitelist, fmt.Sprintf(routeFormat, input.Method, input.Path))
	if inWhitelist {
		result.AuthenticationRequired = false
		result.AuthorizationRequired = false
		return result, nil
	}
	claim, err := l.GetUserClaimFromCtx(ctx)
	if err != nil {
		return result, nil
	}
	// check authorization
	authorizationWhitelist, err := l.getAuthorizationWhitelist(ctx, input.Type)
	if err != nil {
		return
	}
	inWhitelist = slice.Contain(authorizationWhitelist, fmt.Sprintf(routeFormat, input.Method, input.Path))
	if inWhitelist {
		result.AuthorizationRequired = false
		return result, nil
	}

	organizationUserId := claim.OrganizationUserId
	organizationId := claim.OrganizationId
	roles, err := l.getOrganizationUserRoles(ctx, organizationUserId, organizationId, input.Type)
	if err != nil {
		return nil, err
	}
	if len(roles) == 0 {
		return result, nil
	}
	// check admin
	isAdmin := slice.Some(roles, func(index int, item AuthRole) bool {
		return item.IsAdmin
	})
	if isAdmin {
		result.Authorized = true
		result.DataScope = constants.DataScopeAll
		return result, nil
	}

	return
}

// LoadAuthenticationWhitelist
func (l *AuthLogic) LoadAuthenticationWhitelist(ctx context.Context, userType constants.UserType) (whitelist []string, err error) {
	sf := singleflight.Group{}
	// get apis from db
	result, err, _ := sf.Do(stringx.Format("AuthenticationWhitelist:%v", userType), func() (interface{}, error) {
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

// LoadAuthorizationWhitelist
func (l *AuthLogic) LoadAuthorizationWhitelist(ctx context.Context, userType constants.UserType) (whitelist []string, err error) {
	sf := singleflight.Group{}
	result, err, _ := sf.Do(stringx.Format("AuthorizationWhitelist:%v", userType), func() (interface{}, error) {
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

func (l *AuthLogic) GetUserClaimFromCtx(ctx context.Context) (claim utils.JwtUser, err error) {
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

// get authentication whitelist
func (l *AuthLogic) getAuthenticationWhitelist(ctx context.Context, userType constants.UserType) (whitelist []string, err error) {
	var ok bool
	whitelist, ok = l.authenticationWhitelistMap.Get(userType)
	if ok {
		return whitelist, nil
	}

	whitelist, err = l.LoadAuthenticationWhitelist(ctx, userType)
	return whitelist, err
}

// get authorization whitelist
func (l *AuthLogic) getAuthorizationWhitelist(ctx context.Context, userType constants.UserType) (whitelist []string, err error) {
	var ok bool
	whitelist, ok = l.authorizationWhitelistMap.Get(userType)
	if ok {
		return whitelist, nil
	}
	whitelist, err = l.LoadAuthorizationWhitelist(ctx, userType)
	return whitelist, err
}

// getOrganizationUserRoles get organization user roles
func (l *AuthLogic) getOrganizationUserRoles(ctx context.Context, organizationUserId int64, organizationId int64, userType constants.UserType) (roles []AuthRole, err error) {
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
func (l *AuthLogic) loadRoles(ctx context.Context, organizationUserId int64, organizationId int64, userType constants.UserType) (roles []AuthRole, err error) {
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

// getResourcesByRoute
func (l *AuthLogic) getResourcesByRoute(ctx context.Context, organizationId int64, route string, userType constants.UserType) (resources []*model.Resource, err error) {
	organizationResources, ok := l.organizationApiResourceMap.Get(organizationId)
	if !ok {
		organizationResources, err = l.LoadOrganizationResource(ctx, organizationId, userType)
		if err != nil {
			return nil, err
		}
	}

	resources, ok = organizationResources.Get(route)
	return resources, nil
}

// LoadOrganizationResource
func (l *AuthLogic) LoadOrganizationResource(ctx context.Context, organizationId int64, userType constants.UserType) (resources *mapx.Map[string, []*model.Resource], err error) {
	apis, err := l.repos.ApiRepo.GetAllApis(ctx, userType, true)
	if err != nil {
		return nil, err
	}
	apis = slice.Filter(apis, func(index int, item *model.Api) bool {
		return len(item.Path) > 0 && len(item.Method) > 0
	})
	apiResourceMap := mapx.NewMap[string, []*model.Resource]()
	for _, api := range apis {
		apiResources := slice.Filter(api.Resources, func(index int, item *model.Resource) bool {
			return item.Status == constants.CommonStatusActive
		})
		apiResourceMap.Set(stringx.Format(routeFormat, api.Method, api.Path), apiResources)
	}
	l.organizationApiResourceMap.Set(organizationId, apiResourceMap)
	return apiResourceMap, nil
}
