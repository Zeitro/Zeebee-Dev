package global

import (
	"zeebee-backend/app/admin/config"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/limiter"
)

var AdminTryLoginLimiter *limiter.RedisLimiter

// var SendLoginEmailLimiter *limiter.RedisLimiter
// var SendRegisterEmailLimiter *limiter.RedisLimiter
// var SendResetPasswordEmailLimiter *limiter.RedisLimiter

func initLimiter(cnf config.RedisConf) (err error) {
	rdb, err := limiter.GetLimiterRedisClient(limiter.LimiterRedisConf{
		Host: cnf.Host,
		Pass: cnf.Pass,
		DB:   cnf.DB,
		Tls:  cnf.Tls,
	})
	if err != nil {
		return errorx.Wrap(err, "Init Limiter Error")
	}
	// SendLoginEmailLimiter = limiter.NewRedisLimiterWithRedis(rdb)
	// SendRegisterEmailLimiter = limiter.NewRedisLimiterWithRedis(rdb)
	// SendResetPasswordEmailLimiter = limiter.NewRedisLimiterWithRedis(rdb)
	AdminTryLoginLimiter = limiter.NewRedisLimiterWithRedis(rdb)

	return
}
