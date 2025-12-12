package dto

import (
	"zeebee-backend/common/redisx"
	"zeebee-backend/domains/zeebee/constants"
)

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

type UserLoginRequestDto struct {
	Email    string             `json:"email"`
	Password string             `json:"password"`
	UserType constants.UserType `json:"user_type"`
}

type LoginLimiterDto struct {
	Rate     int
	Period   int
	RedisCli *redisx.Redisx
}
