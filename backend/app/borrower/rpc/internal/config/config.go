package config

import (
	"github.com/zeromicro/go-zero/zrpc"

	"zeebee-backend/common/types"
)

type Config struct {
	zrpc.RpcServerConf
	JwtAuth types.JwtAuth `json:",optional"`
	DbConf  DbConf
}

type DbConf struct {
	Dsn string `json:"dsn"`
}
