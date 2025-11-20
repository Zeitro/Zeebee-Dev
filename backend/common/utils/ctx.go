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

func GetUserClaimFromCtx(ctx context.Context) (claim *JWTClaims, err error) {
	val, ok := ctx.Value(constants.BorrowerUserKey).(*JWTClaims)
	if !ok {
		return nil, errorx.WithStack(errorx.New("Invalid user claim"))
	}
	return val, nil
}

func GetAccountIdFromCtx(ctx context.Context) (accountId int64, err error) {
	claim, err := GetUserClaimFromCtx(ctx)
	if err != nil {
		return 0, err
	}
	accountId, err = strconv.ParseInt(claim.AccountId, 10, 64)
	if err != nil {
		return 0, errorx.Wrap(err, "Invalid account id")
	}

	return accountId, nil
}

func GetLenderUserClaimFromCtx(ctx context.Context) (claim *JWTClaims, err error) {
	val, ok := ctx.Value(constants.LenderUserKey).(*JWTClaims)
	if !ok {
		return nil, errorx.WithStack(errorx.New("Invalid user claim"))
	}
	return val, nil
}

func GetLenderAccountIdFromCtx(ctx context.Context) (accountId int64, err error) {
	claim, err := GetLenderUserClaimFromCtx(ctx)
	if err != nil {
		return 0, err
	}
	accountId, err = strconv.ParseInt(claim.AccountId, 10, 64)
	if err != nil {
		return 0, errorx.Wrap(err, "Invalid account id")
	}

	return accountId, nil
}

func GetLenderCompanyIdFromCtx(ctx context.Context) (lenderCompanyId int64, err error) {
	claim, err := GetLenderUserClaimFromCtx(ctx)
	if err != nil {
		return 0, err
	}

	lenderCompanyId = claim.LenderCompanyId
	if lenderCompanyId == 0 {
		return 0, errorx.WithStack(errorx.New("Invalid lender company id"))
	}

	return lenderCompanyId, nil
}

func GetLenderCompanyAndEmployeeId(ctx context.Context) (companyId, employeeId int64, err error) {
	claim, err := GetLenderUserClaimFromCtx(ctx)
	if err != nil {
		return 0, 0, err
	}
	companyId = claim.LenderCompanyId
	employeeId, err = strconv.ParseInt(claim.AccountId, 10, 64)
	if err != nil {
		return 0, 0, errorx.Wrap(err, "Invalid employee Id")
	}
	return companyId, employeeId, nil
}
