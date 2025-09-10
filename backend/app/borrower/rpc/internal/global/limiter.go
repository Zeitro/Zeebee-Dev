package global

import (
	"zeebee-backend/app/borrower/config"
	"zeebee-backend/common/limiter"
)

var BorrowerTryLoginLimiter *limiter.RedisLimiter
var SendLoginEmailLimiter *limiter.RedisLimiter
var SendRegisterEmailLimiter *limiter.RedisLimiter
var SendResetPasswordEmailLimiter *limiter.RedisLimiter

func initLimiter(cnf config.RedisConf) (err error) {
	rdb, err := limiter.GetLimiterRedisClient(cnf.Host, cnf.Pass, 0)
	if err != nil {
		return err
	}
	SendLoginEmailLimiter = limiter.NewRedisLimiterWithRedis(rdb)
	SendRegisterEmailLimiter = limiter.NewRedisLimiterWithRedis(rdb)
	SendResetPasswordEmailLimiter = limiter.NewRedisLimiterWithRedis(rdb)
	BorrowerTryLoginLimiter = limiter.NewRedisLimiterWithRedis(rdb)

	return
}
