package constants

type UserStatus string

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
