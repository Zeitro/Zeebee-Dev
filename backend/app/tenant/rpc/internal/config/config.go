package config

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"

	"zeebee-backend/common/types"
	domainConf "zeebee-backend/domains/zeebee/config"
)

// Config is the tenant service runtime config.
// tenant = RPC-only identity + RBAC/entitlement authority (port 10401, no HTTP layer).
type Config struct {
	zrpc.RpcServerConf

	DbConf    DbConf
	RedisConf RedisConf
	JwtConf   types.JwtConf `json:",optional"`
	EmailConf EmailConf

	// Configurable rate-limit policy (FR-003: not hardcoded).
	LimiterConf LimiterConf

	// One-time token / link expiry (seconds).
	VerifyEmailExpireSeconds     int `json:",default=1800"`
	ResetPwdExpireSeconds        int `json:",default=1800"`
	InviteUserExpireSeconds      int `json:",default=86400"`
	SendInviteEmailPeriodSeconds int `json:",default=300"`

	// Front-end URLs injected into account emails (per userType via Branding).
	LoginUrl         string `json:",optional"`
	VerifyEmailUrl   string `json:",optional"`
	ResetPwdUrl      string `json:",optional"`
	BorrowerLoginUrl string `json:",optional"`
	InviteUserUrl    string `json:",optional"`

	// Casbin: tenant is the sole runtime executor (direction B).
	CasbinConf CasbinConf

	// Branding per userType + environment (ODD-7).
	Branding BrandingConf
}

type DbConf struct {
	Dsn string `json:"dsn"`
}

type RedisConf struct {
	redis.RedisConf
	DB     int    `json:",optional"`
	Prefix string `json:",optional"`
}

type EmailConf struct {
	Host     string
	Port     int
	Password string
	Username string
}

// LimiterConf holds configurable rate-limit items (FR-003).
type LimiterConf struct {
	Login             domainConf.LimiterItem
	SendLoginEmail    domainConf.LimiterItem
	SendVerifyEmail   domainConf.LimiterItem
	SendResetPwdEmail domainConf.LimiterItem
}

type CasbinConf struct {
	ModelPath string
	Watcher   CasbinWatcher
}

type CasbinWatcher struct {
	Host    string
	Pwd     string
	Channel string
	Tls     bool `json:",default=false"`
}

// BrandingConf holds per-userType brand params (ODD-7). Borrower emails must not
// point at test-environment URLs.
type BrandingConf struct {
	Borrower BrandProfile
	Lender   BrandProfile
	Admin    BrandProfile
}

type BrandProfile struct {
	LogoUrl       string `json:",optional"`
	FooterAddress string `json:",optional"`
	SenderName    string `json:",optional"`
	SenderEmail   string `json:",optional"`
}
