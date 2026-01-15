package global

import (
	"zeebee-backend/app/lender/config"
)

var Config config.RpcConf

func SetGlobalConfig(cnf config.RpcConf) {
	Config = cnf
}
