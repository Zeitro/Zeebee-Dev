package config

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"

	"zeebee-backend/common/types"
)

type Config struct {
	ApiConf ApiConf
	RpcConf RpcConf
}

type ApiConf struct {
	rest.RestConf
	SwaggerPath string             `json:",optional"`
	JwtAuth     types.JwtAuth      `json:",optional"`
	RpcConf     zrpc.RpcClientConf `json:",optional"`
}

type RpcConf struct {
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
