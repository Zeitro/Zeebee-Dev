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
	JwtAuth types.JwtAuth      `json:",optional"`
	RpcConf zrpc.RpcClientConf `json:",optional"`
}

type RpcConf struct {
	zrpc.RpcServerConf
	LimiterConf LimiterConf
	EmailConf   EmailConf
	RedisConf   RedisConf
	JwtAuth     types.JwtAuth `json:",optional"`
	DbConf      DbConf
	AwsS3Conf   AwsS3Conf
	CasbinConf  CasbinConf
}

type DbConf struct {
	Dsn string `json:"dsn"`
}

type RedisConf struct {
	redis.RedisConf
	DB     int    `json:",optional"`
	Prefix string // key prefix
}

type AwsS3Conf struct {
	AccessKey                string
	AccessSecret             string
	Bucket                   string
	Region                   string
	FileFolder               string
	PreSigneUrlExpireSeconds int
}

type LimiterConf struct {
	Login LimiterItem
}

type LimiterItem struct {
	Rate   int
	Period int
}

type EmailConf struct {
	Host     string
	Port     int
	Password string
	Username string
}

type CasbinConf struct {
	ModelPath string
	Watcher   CasbinWatcher
}

type CasbinWatcher struct {
	Host            string
	Pwd             string
	AdminChannel    string
	BorrowerChannel string
	LenderChannel   string
}
