package config

import (
	"zeebee-backend/domains/zeebee/config"

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
	// OpenAIConf   OpenAIConf `json:",optional"`
	LlmConf      LlmConf `json:",optional"`
	CasbinConf   CasbinConf
	CronTaskConf []CronTaskConf
}

// OpenAIConf optional keys for lender discover; prefer env OPENAI_API_KEY in production.
// Deprecated: lender discover now uses LlmConf (claude). Kept for config-file compat.
// type OpenAIConf struct {
// 	APIKey  string `json:",optional"`
// 	Model   string `json:",optional"`
// 	BaseURL string `json:",optional"`
// }

// LlmConf wires the shared LLM provider list + per-biz tuning. Lender discover
// reads LenderDiscoverConf to pick provider + model + token/timeout budgets.
type LlmConf struct {
	config.LlmConf

	LenderDiscoverConf config.LlmBizConf
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
	Login config.LimiterItem
}

type EmailConf struct {
	Host     string
	Port     int
	Password string
	Username string
}

type CasbinConf struct {
	ModelPath  string
	Watcher    CasbinWatcher
	EventTopic CasbinEventTopic
}

type CasbinWatcher struct {
	Host            string
	Pwd             string
	Tls             bool `json:",default=false"`
	AdminChannel    string
	BorrowerChannel string
	LenderChannel   string
}

type CasbinEventTopic struct {
	AdminTopic    string
	BorrowerTopic string
	LenderTopic   string
}

type CronTaskConf struct {
	Name     string `json:"name"`
	Schedule string `json:"schedule"`
}
