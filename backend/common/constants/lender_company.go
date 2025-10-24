package constants

import "zeebee-backend/common/enums"

type LenderCompanyStatus string

const LenderCompanyStatusName = "lender_company_status"
const (
	LenderCompanyStatusUnverified LenderCompanyStatus = "Unverified"
	LenderCompanyStatusActive     LenderCompanyStatus = "Active"
	LenderCompanyStatusInactive   LenderCompanyStatus = "Inactive"
)

// Self,Third Party, Partner
type LenderCompanyOwnerType string

const LenderCompanyOwnerTypeName = "lender_company_owner_type"
const (
	LenderCompanyOwnerTypeSelf       LenderCompanyOwnerType = "Self"
	LenderCompanyOwnerTypeThirdParty LenderCompanyOwnerType = "ThirdParty"
	LenderCompanyOwnerTypePartner    LenderCompanyOwnerType = "Partner"
)

func init() {
	// register lender company status
	enums.Register(LenderCompanyStatusName, []enums.EnumOpt{
		{
			Label: string(LenderCompanyStatusUnverified),
			Value: string(LenderCompanyStatusUnverified),
		},
		{
			Label: string(LenderCompanyStatusActive),
			Value: string(LenderCompanyStatusActive),
		},
		{
			Label: string(LenderCompanyStatusInactive),
			Value: string(LenderCompanyStatusInactive),
		},
	})

	// register lender company owner type
	enums.Register(LenderCompanyOwnerTypeName, []enums.EnumOpt{
		{
			Label: string(LenderCompanyOwnerTypeSelf),
			Value: string(LenderCompanyOwnerTypeSelf),
		},
		{
			Label: string(LenderCompanyOwnerTypeThirdParty),
			Value: string(LenderCompanyOwnerTypeThirdParty),
		},
		{
			Label: string(LenderCompanyOwnerTypePartner),
			Value: string(LenderCompanyOwnerTypePartner),
		},
	})
}
