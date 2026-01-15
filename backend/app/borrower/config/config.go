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
	ResetPwdExpireSeconds        int
	PreApprovalLoanExpireSeconds int
	HighCostLimitValue           int64
	EmailConf                    EmailConf
	EmailAuthConf                EmailAuthConf
	RedisConf                    RedisConf
	JwtAuth                      types.JwtAuth `json:",optional"`
	DbConf                       DbConf
	AwsS3Conf                    AwsS3Conf
	LoginUrl                     string
	VerifyEmailUrl               string
	VerifyEmailExpireSeconds     int
	ResetPwdUrl                  string
	LenderUrl                    string
	InviteUserUrl                string
	InviteUserExpireSeconds      int
	SendInviteEmailPeriodSeconds int
	PreApprovalLoanConf          []PreApprovalLoanConf
	ZeebeeAgentEmail             string
	ZeitroConf                   ZeitroConf
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
	Login              LimiterItem
	SendLoginEmail     LimiterItem
	SendVerityEmail    LimiterItem
	SendResetPwdEmail  LimiterItem
	SendLoanQuoteEmail LimiterItem
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

type EmailAuthConf struct {
	Gmail struct {
		ClientId     string
		ClientSecret string
		CallbackUrl  string
	}
}

type ZeitroConf struct {
	Url      string
	Username string
	Password string
}
