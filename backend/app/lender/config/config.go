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
	JwtAuth         types.JwtAuth      `json:",optional"`
	RpcConf         zrpc.RpcClientConf `json:",optional"`
	PdfGeneratorUrl string
}

type RpcConf struct {
	zrpc.RpcServerConf
	LimiterConf                  LimiterConf
	VerifyEmailExpireSeconds     int
	ResetPwdExpireSeconds        int
	EmailConf                    EmailConf
	RedisConf                    RedisConf
	JwtAuth                      types.JwtAuth `json:",optional"`
	DbConf                       DbConf
	AwsS3Conf                    AwsS3Conf
	LoginUrl                     string
	VerifyEmailUrl               string
	ResetPwdUrl                  string
	BorrowerLoginUrl             string
	InviteUserUrl                string
	InviteUserExpireSeconds      int
	SendInviteEmailPeriodSeconds int
	CasbinConf                   CasbinConf
	CronTaskConf                 []CronTaskConf
	StripeConf                   config.StripeConf `json:",optional"`
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
	Login             config.LimiterItem
	SendLoginEmail    config.LimiterItem
	SendVerityEmail   config.LimiterItem
	SendResetPwdEmail config.LimiterItem
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
	Host          string
	Pwd           string
	LenderChannel string
	Tls           bool `json:",default=false"`
}

type CasbinEventTopic struct {
	LenderTopic string
}

type CronTaskConf struct {
	Name     string `json:"name"`
	Schedule string `json:"schedule"`
}
