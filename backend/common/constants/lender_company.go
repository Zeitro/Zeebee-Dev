package constants

import "zeebee-backend/common/enums"

type LenderCompanyStatus string

const LenderCompanyStatusName = "lender_company_status"
const (
	LenderCompanyStatusUnverified LenderCompanyStatus = "unverified"
	LenderCompanyStatusActive     LenderCompanyStatus = "active"
	LenderCompanyStatusInactive   LenderCompanyStatus = "inactive"
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
}
