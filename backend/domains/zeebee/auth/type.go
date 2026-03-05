package auth

import "zeebee-backend/domains/zeebee/constants"

var routeFormat = "%s:%s"
var organizationUserRoleFormat = "%s:organization:%d:organizationuser:%d:roles" // <type>:organization:<organizationId>:organizationuser:<organizationUserId>:roles
var roleResourceFormat = "%d:%d"                                                // <roleId>:<resourceId>
var organizationUserRoleExpireSeconds = 3600

type AuthRole struct {
	RoleId  int64 `json:"role_id"`
	IsAdmin bool  `json:"is_admin"`
}

type AuthResource struct {
	ResourceId     int64 `json:"resource_id"`
	OrganizationId int64 `json:"organization_id"`
}

type RoleResource struct {
	RoleId     int64               `json:"role_id"`
	ResourceId int64               `json:"resource_id"`
	DataScope  constants.DataScope `json:"data_scope"`
}
