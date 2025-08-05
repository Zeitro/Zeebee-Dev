package global

import (
	"zeebee-backend/app/borrower/config"
)

var Config config.RpcConf

func SetGlobalConfig(cnf config.RpcConf) {
	Config = cnf
}
