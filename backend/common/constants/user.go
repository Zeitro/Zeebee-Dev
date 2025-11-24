package constants

import "zeebee-backend/common/enums"

type UserStatus string

var UserStatusName = "user_status"

const (
	UserStatusUnverified UserStatus = "Unverified"
	UserStatusActive     UserStatus = "Active"
	UserStatusInactive   UserStatus = "Inactive"
)

type UserType string

const (
	UserTypeBorrower UserType = "Borrower"
	UserTypeLender   UserType = "Lender"
	UserTypeAdmin    UserType = "Admin"
)

const (
	UserIdCtxKey    = "user_id"
	UserClaimCtxKey = "user_claim"

	// BorrowerUserIdKey = "borrower_user_id"
	UserJwtIdKey = "user_jwt"
	// LenderUserIdKey   = "lender_user_id"
	// AdminUserIdKey   = "admin_user_id"

	// LenderUserCompanyIdKey = "lender_user_company_id"

	BorrowerUserAccountIdKey = "borrower_user_account_id"
	LenderUserAccountIdKey   = "lender_user_account_id"

	// BorrowerUserKey = "borrower_user"
	// LenderUserKey   = "lender_user"
	// AdminUserKey   = "admin_user"
)

func init() {
	// register user status
	enums.Register(UserStatusName, []enums.EnumOpt{
		{
			Label: string(UserStatusUnverified),
			Value: string(UserStatusUnverified),
		},
		{
			Label: string(UserStatusActive),
			Value: string(UserStatusActive),
		},
		{
			Label: string(UserStatusInactive),
			Value: string(UserStatusInactive),
		},
	})
}
