package global

import (
	"github.com/pkg/errors"

	"zeebee-backend/app/borrower/config"
	"zeebee-backend/common/limiter"
)

var BorrowerTryLoginLimiter *limiter.RedisLimiter
var SendLoginEmailLimiter *limiter.RedisLimiter
var SendRegisterEmailLimiter *limiter.RedisLimiter
var SendResetPasswordEmailLimiter *limiter.RedisLimiter

func initLimiter(cnf config.RedisConf) (err error) {
	rdb, err := limiter.GetLimiterRedisClient(limiter.LimiterRedisConf{
		Host: cnf.Host,
		Pass: cnf.Pass,
		DB:   cnf.DB,
		Tls:  cnf.Tls,
	})
	if err != nil {
		return errors.Wrap(err, "init limiter error")
	}
	SendLoginEmailLimiter = limiter.NewRedisLimiterWithRedis(rdb)
	SendRegisterEmailLimiter = limiter.NewRedisLimiterWithRedis(rdb)
	SendResetPasswordEmailLimiter = limiter.NewRedisLimiterWithRedis(rdb)
	BorrowerTryLoginLimiter = limiter.NewRedisLimiterWithRedis(rdb)

	return
}
