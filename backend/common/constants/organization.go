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

type OrganizationAndUserStatus string

var OrganizationAndUserStatusName = "organization_and_user_status"

const (
	OrganizationAndUserStatusActive   OrganizationAndUserStatus = "Active"
	OrganizationAndUserStatusDisabled OrganizationAndUserStatus = "Disabled"
)

type OrganizationAndUserActivationStatus string

var OrganizationAndUserActivationStatusName = "organization_and_user_activation_status"

const (
	OrganizationAndUserActivationStatusActivated    OrganizationAndUserActivationStatus = "Activated"
	OrganizationAndUserActivationStatusNotActivated OrganizationAndUserActivationStatus = "Not Activated"
)

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
	enums.Register(OrganizationAndUserStatusName, []enums.EnumOpt{
		{
			Label: string(OrganizationAndUserStatusActive),
			Value: string(OrganizationAndUserStatusActive),
		},
		{
			Label: string(OrganizationAndUserStatusDisabled),
			Value: string(OrganizationAndUserStatusDisabled),
		},
	})
	// register organization and user activation status
	enums.Register(OrganizationAndUserActivationStatusName, []enums.EnumOpt{
		{
			Label: string(OrganizationAndUserActivationStatusActivated),
			Value: string(OrganizationAndUserActivationStatusActivated),
		},
		{
			Label: string(OrganizationAndUserActivationStatusNotActivated),
			Value: string(OrganizationAndUserActivationStatusNotActivated),
		},
	})

}
