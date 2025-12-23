package global

import (
	"time"

	"github.com/go-redis/redis_rate/v10"

	"zeebee-backend/app/borrower/config"
	"zeebee-backend/common/errorx"
	commonlimiter "zeebee-backend/common/limiter"
	"zeebee-backend/domains/zeebee/limiter"
)

var BorrowerTryLoginLimiter *commonlimiter.RedisLimiter
var SendLoginEmailLimiter *commonlimiter.RedisLimiter
var SendRegisterEmailLimiter *commonlimiter.RedisLimiter
var SendResetPasswordEmailLimiter *commonlimiter.RedisLimiter

var LimiterMgr *limiter.LimiterManager

func initLimiter(cnf config.RpcConf) (err error) {
	rdb, err := commonlimiter.GetLimiterRedisClient(commonlimiter.LimiterRedisConf{
		Host: cnf.RedisConf.Host,
		Pass: cnf.RedisConf.Pass,
		DB:   cnf.RedisConf.DB,
		Tls:  cnf.RedisConf.Tls,
	})
	if err != nil {
		return errorx.Wrap(err, "Init Limiter Error")
	}

	limiterConf := cnf.LimiterConf

	LimiterMgr = limiter.NewLimiterManager(commonlimiter.NewRedisLimiterWithRedis(rdb), map[limiter.LimiterStrategy]redis_rate.Limit{
		limiter.StrategySendLoginEmail: redis_rate.Limit{
			Rate:   limiterConf.SendLoginEmail.Rate,
			Burst:  limiterConf.SendLoginEmail.Rate,
			Period: time.Duration(int64(limiterConf.SendLoginEmail.Period)) * time.Second,
		},
	})

	SendLoginEmailLimiter = commonlimiter.NewRedisLimiterWithRedis(rdb)
	SendRegisterEmailLimiter = commonlimiter.NewRedisLimiterWithRedis(rdb)
	SendResetPasswordEmailLimiter = commonlimiter.NewRedisLimiterWithRedis(rdb)
	BorrowerTryLoginLimiter = commonlimiter.NewRedisLimiterWithRedis(rdb)

	return
}
