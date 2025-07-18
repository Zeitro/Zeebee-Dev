package global

import "zeebee-backend/app/borrower/rpc/internal/config"

var Config config.Config

func InitGlobal(c config.Config) (err error) {
	Config = c
	err = initDB(Config.DbConf.Dsn, Config.Mode == "dev")
	if err != nil {
		return err
	}
	return err
}
