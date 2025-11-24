package constants

import (
	"zeebee-backend/common/enums"
	"zeebee-backend/common/errorx"
)

// type OrganizationType string
//
// var OrganizationTypeName = "organization_type"
//
// const (
//
//	OrganizationTypeBorrower OrganizationType = "Borrower"
//	OrganizationTypeLender   OrganizationType = "Lender"
//	OrganizationTypeAdmin    OrganizationType = "Admin"
//
// )
//
//	func (v OrganizationType) Validate() error {
//		if v == OrganizationTypeBorrower || v == OrganizationTypeLender || v == OrganizationTypeAdmin {
//			return nil
//		}
//		return errorx.Errorf("Invalid Organization Type: %v", v)
//	}
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

func init() {
	// register organization status
	// enums.Register(OrganizationTypeName, []enums.EnumOpt{
	// 	{
	// 		Label: string(OrganizationTypeBorrower),
	// 		Value: string(OrganizationTypeBorrower),
	// 	},
	// 	{
	// 		Label: string(OrganizationTypeLender),
	// 		Value: string(OrganizationTypeLender),
	// 	},
	// 	{
	// 		Label: string(OrganizationTypeAdmin),
	// 		Value: string(OrganizationTypeAdmin),
	// 	},
	// })
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
