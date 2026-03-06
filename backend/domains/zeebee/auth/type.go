package auth

import (
	"encoding/json"

	"zeebee-backend/domains/zeebee/constants"
)

const routeFormat = "%s:%s"
const organizationUserRoleFormat = "%s:organization:%d:organizationuser:%d:roles" // <type>:organization:<organizationId>:organizationuser:<organizationUserId>:roles
const roleResourceFormat = "%d:%d"                                                // <roleId>:<resourceId>
const organizationUserRoleExpireSeconds = 3600

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

type Action string

const (
	ActionApiReload      Action = "ApiReload"
	ActionResourceReload Action = "ResourceReload"
	ActionRoleReload     Action = "RoleReload"
)

type Scope string

const (
	ScopeGlobal       Scope = "Global"
	ScopeOrganization Scope = "Organization"
)

// type AuthType string
//
// const (
// 	AuthTypeAdmin    AuthType = "Admin"
// 	AuthTypeBorrower AuthType = "Borrower"
// 	AuthTypeLender   AuthType = "Lender"
// )

type Msg struct {
	Action Action             `json:"action"`
	Scope  Scope              `json:"scope"`
	Type   constants.UserType `json:"type"`
}

func (msg Msg) String() string {
	content, _ := json.Marshal(msg)
	return string(content)
}

func ParseMsg(msg string) (Msg, error) {
	var m Msg
	err := json.Unmarshal([]byte(msg), &m)
	return m, err
}
