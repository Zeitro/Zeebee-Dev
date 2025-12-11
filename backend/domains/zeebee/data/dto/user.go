package dto

import "zeebee-backend/common/constants"

type UserChangePasswordRequestDto struct {
	OldPassword string             `json:"old_password"`
	NewPassword string             `json:"new_password"`
	UserId      int64              `json:"user_id"`
	UserType    constants.UserType `json:"user_type"`
}

type UserResetPasswordByEmailRequestDto struct {
	UserType    constants.UserType `json:"user_type"`
	NewPassword string             `json:"new_password"`
	Email       string             `json:"email"`
	Code        string             `json:"code"`
}
