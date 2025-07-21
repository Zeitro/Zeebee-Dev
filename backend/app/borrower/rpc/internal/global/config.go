package global

import (
	"zeebee-backend/app/borrower/rpc/internal/config"
)

var Config config.Config

func SetGlobalConfig(cnf config.Config) {
	Config = cnf
}
