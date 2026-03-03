package global

import (
	"time"

	"github.com/zeromicro/go-zero/core/service"

	"zeebee-backend/app/admin/config"
	"zeebee-backend/common/errorx"
)

func InitGlobal(c config.RpcConf) (err error) {
	SetGlobalConfig(c)
	eg := errorx.NewErrorGroup()
	eg.Do(
		func() error {
			return initTime()
		},
		func() error {
			return initDB(Config.DbConf.Dsn, Config.Mode != service.ProMode)
		},
		func() error {
			return initRedis(c.RedisConf)
		},
		func() error {
			return initLimiter(c.RedisConf)
		},
		func() error {
			return initS3Client(c.AwsS3Conf)
		},
		func() error {
			return initEmailClient(c.EmailConf)
		},
		func() error {
			return initCasbin(c.CasbinConf)
		},
	)

	// eg.Go(func() error {
	// 	return initTime()
	// })
	//
	// eg.Go(func() error {
	// 	return initDB(Config.DbConf.Dsn, Config.Mode != service.ProMode)
	// })
	// eg.Go(func() error {
	// 	return initRedis(c.RedisConf)
	// })
	// eg.Go(func() error {
	// 	return initLimiter(c.RedisConf)
	// })
	//
	// eg.Go(func() error {
	// 	return initS3Client(c.AwsS3Conf)
	// })
	// eg.Go(func() error {
	// 	return initEmailClient(c.EmailConf)
	// })

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
