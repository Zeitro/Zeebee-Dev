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
)
