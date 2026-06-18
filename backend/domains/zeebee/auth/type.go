package auth

import (
	"encoding/json"
	"time"

	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/constants"
)

const (
	routeFormat                = "%s:%s"                                        // <method>:<path>
	organizationUserRoleFormat = "%s:organization:%d:organizationuser:%d:roles" // <type>:organization:<organizationId>:organizationuser:<organizationUserId>:roles
	roleResourceFormat         = "%d:%d"                                        // <roleId>:<resourceId>
)

const (
	organizationUserRoleExpireSeconds = 3600
	maxRetryTimes                     = 10
	retryIntervalTime                 = time.Millisecond * 50
	maxRetryIntervalTime              = time.Minute
)

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
	ActionApiReload       Action = "ApiReload"
	ActionResourceReload  Action = "ResourceReload"
	ActionRoleReload      Action = "RoleReload"
	ActionAllPolicyReload Action = "AllPolicyReload"
	ActionPolicyRebuild   Action = "PolicyRebuild"
)

type Scope string

const (
	ScopeGlobal           Scope = "Global"
	ScopeOrganization     Scope = "Organization"
	ScopeOrganizationUser Scope = "OrganizationUser"
)

type Msg struct {
	Action             Action             `json:"action"`
	Scope              Scope              `json:"scope"`
	Type               constants.UserType `json:"type"`
	OrganizationId     int64              `json:"organization_id"`
	OrganizationUserId int64              `json:"organization_user_id"`
	AddedPolicies      [][]string         `json:"added_policies"`
	RemovedPolicies    [][]string         `json:"removed_policies"`
}

func (msg Msg) String() string {
	content, _ := json.Marshal(msg)
	return string(content)
}

func ParseMsg(msg string) (Msg, error) {
	var m Msg
	err := json.Unmarshal([]byte(msg), &m)
	if err != nil {
		return m, errorx.Wrap(err, "Failed to parse auth message")
	}
	return m, err
}
