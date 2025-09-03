package utils

import (
	"context"
	"strconv"

	"github.com/pkg/errors"

	"zeebee-backend/common/constants"
)

func GetInt64FromCtx(ctx context.Context, key string) int64 {
	val := ctx.Value(key)
	switch val.(type) {
	case string:
		i, err := strconv.Atoi(val.(string))
		if err != nil {
			return 0
		}
		return int64(i)
	case int64:
		return val.(int64)
	case int:
		return int64(val.(int))
	case int32:
		return int64(val.(int32))
	case int16:
		return int64(val.(int16))
	case int8:
		return int64(val.(int8))
	default:
		return 0
	}
}

func GetUserClaimFromCtx(ctx context.Context) (claim *JWTClaims, err error) {
	val, ok := ctx.Value(constants.BorrowerUserKey).(*JWTClaims)
	if !ok {
		return nil, errors.WithStack(errors.New("Invalid user claim"))
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
		return 0, errors.Wrap(err, "Invalid account id")
	}

	return accountId, nil
}
