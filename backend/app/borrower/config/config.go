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
	LimiterConf                  LimiterConf
	VerifyEmailExpireSeconds     int
	ResetPwdExpireSeconds        int
	PreApprovalLoanExpireSeconds int
	EmailConf                    EmailConf
	RedisConf                    RedisConf
	JwtAuth                      types.JwtAuth `json:",optional"`
	DbConf                       DbConf
	AwsS3Conf                    AwsS3Conf
	LoginUrl                     string
	VerifyEmailUrl               string
	ResetPwdUrl                  string
	PreApprovalLoanConf          []PreApprovalLoanConf
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
	Login             LimiterItem
	SendLoginEmail    LimiterItem
	SendVerityEmail   LimiterItem
	SendResetPwdEmail LimiterItem
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

// PreApprovalLoanConf pre-approve loan config
type PreApprovalLoanConf struct {
	LenderCompanyId   int64
	LenderCompanyName string
	TotalAmount       PreApprovalLoanTotalAmount
	Rate              PreApprovalLoanRate
	LoanTerm          PreApprovalLoanLoanTerm
}

type PreApprovalLoanTotalAmount struct {
	Min float64
	Max float64
}
type PreApprovalLoanRate struct {
	Min float64
	Max float64
}
type PreApprovalLoanLoanTerm struct {
	Min int
	Max int
}
