package global

import (
	"zeebee-backend/app/borrower/config"
)

var Config config.RpcConf

// SetGlobalConfig sets the global RPC configuration.
// It takes a config.RpcConf parameter and updates the global Config variable.
func SetGlobalConfig(cnf config.RpcConf) {
	Config = cnf
}
