package config

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"

	"zeebee-backend/common/types"
)

type Config struct {
	zrpc.RpcServerConf
	MaxFailedLoginTimes   int // max failed login times
	ForbiddenLoginSeconds int // seconds for next retry login
	RedisConf             RedisConf

	JwtAuth types.JwtAuth `json:",optional"`
	DbConf  DbConf
}

type DbConf struct {
	Dsn string `json:"dsn"`
}

type RedisConf struct {
	redis.RedisConf
	Prefix string // key prefix
}
