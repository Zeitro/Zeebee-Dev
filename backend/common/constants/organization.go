package constants

import (
	"zeebee-backend/common/enums"
	"zeebee-backend/common/errorx"
)

type OrganizationOwnerType string

const OrganizationOwnerTypeName = "lender_company_owner_type"
const (
	OrganizationOwnerTypeNone       OrganizationOwnerType = ""
	OrganizationOwnerTypeSelf       OrganizationOwnerType = "Self"
	OrganizationOwnerTypeThirdParty OrganizationOwnerType = "ThirdParty"
	OrganizationOwnerTypePartner    OrganizationOwnerType = "Partner"
)

type OrganizationStatus string

var OrganizationStatusName = "organization_status"

const (
	OrganizationStatusActive   OrganizationStatus = "Active"
	OrganizationStatusInactive OrganizationStatus = "Inactive"
)

func (v OrganizationStatus) Validate() error {
	if v == OrganizationStatusActive || v == OrganizationStatusInactive {
		return nil
	}
	return errorx.Errorf("Invalid Organization Status: %v", v)
}

type OrganizationUserStatus string

var OrganizationUserStatusName = "organization_user_status"

const (
	OrganizationUserStatusActive       OrganizationUserStatus = "Active"
	OrganizationUserStatusNotActivated OrganizationUserStatus = "Not activated"
	OrganizationUserStatusDisabled     OrganizationUserStatus = "Disabled"
)

func (v OrganizationUserStatus) Validate() error {
	if v == OrganizationUserStatusActive || v == OrganizationUserStatusNotActivated || v == OrganizationUserStatusDisabled {
		return nil
	}
	return errorx.Errorf("Invalid Organization And User Status: %v", v)
}

func init() {
	// register organization owner type
	enums.Register(OrganizationOwnerTypeName, []enums.EnumOpt{
		{
			Label: string(OrganizationOwnerTypeNone),
			Value: string(OrganizationOwnerTypeNone),
		},
		{
			Label: string(OrganizationOwnerTypeSelf),
			Value: string(OrganizationOwnerTypeSelf),
		},
		{
			Label: string(OrganizationOwnerTypeThirdParty),
			Value: string(OrganizationOwnerTypeThirdParty),
		},
		{
			Label: string(OrganizationOwnerTypePartner),
			Value: string(OrganizationOwnerTypePartner),
		},
	})
	// register organization status
	enums.Register(OrganizationStatusName, []enums.EnumOpt{
		{
			Label: string(OrganizationStatusActive),
			Value: string(OrganizationStatusActive),
		},
		{
			Label: string(OrganizationStatusInactive),
			Value: string(OrganizationStatusInactive),
		},
	})
	// register organization and user status
	enums.Register(OrganizationUserStatusName, []enums.EnumOpt{
		{
			Label: string(OrganizationUserStatusActive),
			Value: string(OrganizationUserStatusActive),
		},
		// {
		// 	Label: string(OrganizationUserStatusNotActivated),
		// 	Value: string(OrganizationUserStatusNotActivated),
		// },
		{
			Label: string(OrganizationUserStatusDisabled),
			Value: string(OrganizationUserStatusDisabled),
		},
	})

}
