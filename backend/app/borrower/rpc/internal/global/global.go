package global

import (
	"time"

	"github.com/zeromicro/go-zero/core/service"

	"zeebee-backend/app/borrower/config"
	"zeebee-backend/app/borrower/rpc/internal/statics"
	"zeebee-backend/common/errorx"
)

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
		return initLimiter(c.RedisConf)
	})

	eg.Do(func() error {
		return initS3Client(c.AwsS3Conf)
	})
	eg.Do(func() error {
		return initEmailClient(c.EmailConf)
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

func initTime() (err error) {
	// set utc time zone
	time.Local = time.UTC

	return err
}
