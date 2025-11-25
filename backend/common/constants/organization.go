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

const (
	OrganizationAndUserStatusActive   OrganizationAndUserStatus = "Active"
	OrganizationAndUserStatusInactive OrganizationAndUserStatus = "Inactive"
	OrganizationAndUserStatusDisabled OrganizationAndUserStatus = "Disabled"
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
}
