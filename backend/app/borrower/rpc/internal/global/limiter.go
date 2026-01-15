package global

import (
	"time"

	"github.com/go-redis/redis_rate/v10"

	"zeebee-backend/app/borrower/config"
	"zeebee-backend/common/errorx"
	commonlimiter "zeebee-backend/common/limiter"
	"zeebee-backend/domains/zeebee/limiter"
)

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
		limiter.StrategyLogin: {
			Rate:   limiterConf.Login.Rate,
			Burst:  limiterConf.Login.Rate,
			Period: time.Duration(int64(limiterConf.Login.Period)) * time.Second,
		},
		limiter.StrategySendLoginEmail: {
			Rate:   limiterConf.SendLoginEmail.Rate,
			Burst:  limiterConf.SendLoginEmail.Rate,
			Period: time.Duration(int64(limiterConf.SendLoginEmail.Period)) * time.Second,
		},
		limiter.StrategySendVerifyEmail: {
			Rate:   limiterConf.SendVerityEmail.Rate,
			Burst:  limiterConf.SendVerityEmail.Rate,
			Period: time.Duration(int64(limiterConf.SendVerityEmail.Period)) * time.Second,
		},
		limiter.StrategySendResetPwdEmail: {
			Rate:   limiterConf.SendResetPwdEmail.Rate,
			Burst:  limiterConf.SendResetPwdEmail.Rate,
			Period: time.Duration(int64(limiterConf.SendResetPwdEmail.Period)) * time.Second,
		},
		limiter.StrategySendLoanQuoteEmail: {
			Rate:   limiterConf.SendLoanQuoteEmail.Rate,
			Burst:  limiterConf.SendLoanQuoteEmail.Rate,
			Period: time.Duration(int64(limiterConf.SendLoanQuoteEmail.Period)) * time.Second,
		},
	})

	return
}
