package global

import (
	"time"

	"github.com/zeromicro/go-zero/core/service"

	"zeebee-backend/app/borrower/rpc/internal/config"
)

var Config config.Config

func InitGlobal(c config.Config) (err error) {
	initTime()

	Config = c
	err = initDB(Config.DbConf.Dsn, Config.Mode != service.ProMode)
	if err != nil {
		return err
	}
	return err
}

func initTime() {
	// set utc time zone
	time.Local = time.UTC
}
