package config

import (
	"github.com/zeromicro/go-zero/rest"
	// "github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	DataSource string
	Auth       struct {
		AccessSecret string
		AccessExpire int64
	}
	// Etcd       zrpc.RpcClientConf
	// Streamline zrpc.RpcClientConf
	// Forms      zrpc.RpcClientConf
}
