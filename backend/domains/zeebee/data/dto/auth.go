package dto

import "zeebee-backend/domains/zeebee/constants"

type AuthCheckRequestDto struct {
	Method string             `json:"method"`
	Path   string             `json:"path"`
	Type   constants.UserType `json:"type"`
}

type AuthCheckResponseDto struct {
	AuthenticationRequired bool                `json:"authentication_required"`
	Authenticated          bool                `json:"authenticated"`
	AuthorizationRequired  bool                `json:"authorization_required"`
	Authorized             bool                `json:"authorized"`
	DataScope              constants.DataScope `json:"data_scope"`
}
