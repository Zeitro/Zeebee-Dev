package global

import (
	"time"

	"github.com/zeromicro/go-zero/core/service"

	"zeebee-backend/app/lender/config"
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
		return initLimiter(c)
	})

	eg.Do(func() error {
		return initS3Client(c.AwsS3Conf)
	})
	eg.Do(func() error {
		return initEmailClient(c.EmailConf)
	})
	eg.Do(func() error {
		return initCasbin(c.CasbinConf)
	})

	err = eg.Wait()
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
