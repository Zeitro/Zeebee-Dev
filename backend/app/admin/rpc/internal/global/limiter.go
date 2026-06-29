package global

import (
	"zeebee-backend/app/admin/config"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/limiter"
)

var AdminTryLoginLimiter *limiter.RedisLimiter

func initLimiter(cnf config.RedisConf) (err error) {
	rdb, err := limiter.GetLimiterRedisClient(limiter.LimiterRedisConf{
		Host: cnf.Host,
		Pass: cnf.Pass,
		DB:   cnf.DB,
		Tls:  cnf.Tls,
	})
	if err != nil {
		return errorx.Wrap(err, "init limiter error")
	}
	AdminTryLoginLimiter = limiter.NewRedisLimiterWithRedis(rdb)

	return
}
