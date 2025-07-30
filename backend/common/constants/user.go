package constants

import "zeebee-backend/common/enums"

type UserStatus string

var UserStatusName = "user_status"

const (
	UserStatusActive   UserStatus = "active"
	UserStatusInactive UserStatus = "inactive"
)

type UserType string

const (
	UserTypeBorrower UserType = "borrower"
	UserTypeLender   UserType = "lender"
)

const (
	BorrowerUserIdKey = "borrower_user_id"
	LenderUserIdKey   = "lender_user_id"

	BorrowerUserAccountIdKey = "borrower_user_account_id"
	LenderUserAccountIdKey   = "lender_user_account_id"

	BorrowerUserKey = "borrower_user"
)

func init() {
	// register user status
	enums.Register(UserStatusName, []enums.EnumOpt{
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
