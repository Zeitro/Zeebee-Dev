package global

import (
	"zeebee-backend/app/lender/config"
)

var Config config.ApiConf

func SetGlobalConfig(cnf config.ApiConf) {
	Config = cnf
}
