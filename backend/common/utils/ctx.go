package utils

import (
	"context"
	"strconv"

	"zeebee-backend/common/errorx"
)

const (
	UserClaimCtxKey = "user_claim"
	UserIdCtxKey    = "user_id"
)

func GetInt64FromCtx(ctx context.Context, key string) int64 {
	val := ctx.Value(key)
	switch v := val.(type) {
	case string:
		i, err := strconv.Atoi(v)
		if err != nil {
			return 0
		}
		return int64(i)
	case int64:
		return v
	case int:
		return int64(v)
	case int32:
		return int64(v)
	case int16:
		return int64(v)
	case int8:
		return int64(v)
	default:
		return 0
	}
}

func GetUserClaimFromCtx(ctx context.Context) (claim *JWTClaims, err error) {
	// get from ctx
	val, ok := ctx.Value(UserClaimCtxKey).(*JWTClaims)
	if ok {
		return val, nil
	}
	return nil, errorx.New("Invalid User Claim")

}

func GetOrganizationIdFromCtx(ctx context.Context) (organizationId int64, err error) {
	claim, err := GetUserClaimFromCtx(ctx)
	if err != nil {
		return 0, err
	}
	return claim.OrganizationId, nil
}

func GetOrganizationAndUserId(ctx context.Context) (organizationId, userId int64, err error) {
	claim, err := GetUserClaimFromCtx(ctx)
	if err != nil {
		return 0, 0, err
	}
	return claim.OrganizationId, claim.UserId, nil
}
