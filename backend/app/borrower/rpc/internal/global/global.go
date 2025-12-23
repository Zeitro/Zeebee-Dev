package global

import (
	"github.com/zeromicro/go-zero/core/service"

	"zeebee-backend/app/borrower/config"
	"zeebee-backend/app/borrower/rpc/internal/statics"
	"zeebee-backend/common/errorx"
)

// InitGlobal initializes all global components and services required by the application.
// It takes an RPC configuration parameter and performs initialization of various components
// including time, database, Redis, rate limiter, S3 client, email client, and high-cost limit manager.
// Returns an error if any of the initialization steps fail.
//
// The function performs the following initializations:
// - System time configuration
// - Database connection with debug mode based on service mode
// - Redis connection
// - Rate limiter setup
// - AWS S3 client
// - Email client
// - High-cost limit manager with county loan limits
func InitGlobal(c config.RpcConf) (err error) {
	SetGlobalConfig(c)
	eg := errorx.NewErrorGroup()
	eg.Do(func() error {
		return initTime()
	})
	eg.Do(func() error {
		return initDB(Config.DbConf.Dsn, Config.Mode != service.ProMode)
	})
	eg.Do(func() error {
		return initRedis(c.RedisConf)
	})
	eg.Do(func() error {
		return initLimiter(c)
	})

	eg.Do(func() error {
		return initS3Client(c.AwsS3Conf)
	})
	eg.Do(func() error {
		return initEmailClient(c)
	})
	eg.Do(func() error {
		return initStateHighCostLimitMgr(statics.FullCountyLoanLimitListContent, c.HighCostLimitValue)
	})

	err = eg.Error()
	if err != nil {
		return err
	}
	return err
}
