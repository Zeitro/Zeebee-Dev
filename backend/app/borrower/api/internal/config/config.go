package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"

	"zeebee-backend/common/types"
)

type Config struct {
	rest.RestConf
	SwaggerPath string             `json:",optional"`
	JwtConf     types.JwtConf      `json:",optional"`
	RpcConf     zrpc.RpcClientConf `json:",optional"`
}
