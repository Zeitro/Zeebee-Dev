package utils

import (
	"context"
	"strconv"
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
