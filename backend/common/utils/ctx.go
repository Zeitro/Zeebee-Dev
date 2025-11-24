package utils

import (
	"context"
	"strconv"

	"zeebee-backend/common/errorx"

	"zeebee-backend/common/constants"
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

// func GetClaimFromCtx(ctx context.Context) (claim *JWTClaims, err error) {
// 	val, ok := ctx.Value(constants.BorrowerUserKey).(*JWTClaims)
// 	if !ok {
// 		return nil, errorx.WithStack(errorx.New("Invalid user claim"))
// 	}
// 	return val, nil
// }

func GetUserClaimFromCtx(ctx context.Context) (claim *JWTClaims, err error) {
	val, ok := ctx.Value(constants.UserClaimCtxKey).(*JWTClaims)
	if !ok {
		return nil, errorx.WithStack(errorx.New("Invalid User Claim"))
	}
	return val, nil
}

func GetOrganizationIdFromCtx(ctx context.Context) (organizationId int64, err error) {
	claim, err := GetUserClaimFromCtx(ctx)
	if err != nil {
		return 0, err
	}
	return claim.OrganizationId, nil
}

// func GetLenderUserClaimFromCtx(ctx context.Context) (claim *JWTClaims, err error) {
// 	val, ok := ctx.Value(constants.LenderUserKey).(*JWTClaims)
// 	if !ok {
// 		return nil, errorx.WithStack(errorx.New("Invalid user claim"))
// 	}
// 	return val, nil
// }

// func GetLenderAccountIdFromCtx(ctx context.Context) (accountId int64, err error) {
// 	claim, err := GetLenderUserClaimFromCtx(ctx)
// 	if err != nil {
// 		return 0, err
// 	}
// 	accountId, err = strconv.ParseInt(claim.AccountId, 10, 64)
// 	if err != nil {
// 		return 0, errorx.Wrap(err, "Invalid account id")
// 	}
//
// 	return accountId, nil
// }

// func GetOrganizationIdFromCtx(ctx context.Context) (organizationId int64, err error) {
// 	claim, err := GetLenderUserClaimFromCtx(ctx)
// 	if err != nil {
// 		claim, err = GetUserClaimFromCtx(ctx)
// 		if err != nil {
// 			return 0, err
// 		}
// 	}
//
// 	organizationId = claim.OrganizationId
// 	if organizationId == 0 {
// 		return 0, errorx.WithStack(errorx.New("Invalid organization id"))
// 	}
//
// 	return organizationId, nil
// }

// 兼容性函数，保持向后兼容
// func GetLenderCompanyIdFromCtx(ctx context.Context) (lenderCompanyId int64, err error) {
// 	return GetOrganizationIdFromCtx(ctx)
// }

func GetOrganizationAndUserId(ctx context.Context) (organizationId, userId int64, err error) {
	claim, err := GetUserClaimFromCtx(ctx)
	if err != nil {
		return 0, 0, err
	}
	return claim.OrganizationId, claim.UserId, nil
}

// 新增通用函数：获取组织和用户ID
// func GetOrganizationAndUserId(ctx context.Context) (organizationId, userId int64, err error) {
// 	claim, err := GetLenderUserClaimFromCtx(ctx)
// 	if err != nil {
// 		// 尝试从 Borrower 上下文获取
// 		claim, err = GetUserClaimFromCtx(ctx)
// 		if err != nil {
// 			return 0, 0, err
// 		}
// 	}
// 	organizationId = claim.OrganizationId
// 	userId, err = strconv.ParseInt(claim.AccountId, 10, 64)
// 	if err != nil {
// 		return 0, 0, errorx.Wrap(err, "Invalid user Id")
// 	}
// 	return organizationId, userId, nil
// }
