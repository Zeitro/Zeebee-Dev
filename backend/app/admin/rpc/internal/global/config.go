package global

import (
	"zeebee-backend/app/admin/config"
)

var Config config.RpcConf

func SetGlobalConfig(cnf config.RpcConf) {
	Config = cnf
}
