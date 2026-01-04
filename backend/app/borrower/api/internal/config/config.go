package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"

	"zeebee-backend/common/types"
)

type Config struct {
	rest.RestConf
	SwaggerPath string             `json:",optional"`
	JwtAuth     types.JwtAuth      `json:",optional"`
	RpcConf     zrpc.RpcClientConf `json:",optional"`
}

// type RpcConf struct {
// 	UserConf zrpc.RpcClientConf
// }
