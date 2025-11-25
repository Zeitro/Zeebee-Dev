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
