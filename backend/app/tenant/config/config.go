package config

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"

	"zeebee-backend/common/types"
	domainConf "zeebee-backend/domains/zeebee/config"
)

// Config is the tenant service top-level config.
// tenant is RPC-only today; ApiConf is reserved for a future HTTP layer
// (mirrors the admin/lender Config{ApiConf, RpcConf} shape).
type Config struct {
	RpcConf RpcConf
	// ApiConf ApiConf // reserved: add when tenant exposes HTTP directly
}

// RpcConf holds the tenant RPC server + shared infra sub-configs.
type RpcConf struct {
	zrpc.RpcServerConf

	DbConf    DbConf
	RedisConf RedisConf

	// Configurable rate-limit policy (FR-003: not hardcoded).
	LimiterConf LimiterConf

	// Casbin: tenant is the sole runtime executor (direction B).
	CasbinConf CasbinConf

	// JwtConf is the single shared JWT signing config across all platforms
	// (borrower/lender/admin). userType is carried in the token payload, so one
	// key signs and verifies every token (ODD-5: single unified key).
	JwtConf types.JwtConf

	// Per-platform brand + front-end URLs + expiry (ODD-7). JWT signing is shared
	// (see JwtConf above); only brand / URLs / email differ per platform.
	PlatformConf PlatformConf
}

type DbConf struct {
	Dsn string `json:"dsn"`
}

type RedisConf struct {
	redis.RedisConf
	DB     int    `json:",optional"`
	Prefix string `json:",optional"`
}

// EmailConf covers both SMTP and Gmail OAuth sending. All fields optional so a
// platform configures whichever mechanism it uses.
type EmailConf struct {
	Host     string    `json:",optional"`
	Port     int       `json:",optional"`
	Password string    `json:",optional"`
	Username string    `json:",optional"`
	Gmail    GmailConf `json:",optional"`
}

// GmailConf holds Gmail OAuth app credentials (for sending via a connected account).
type GmailConf struct {
	ClientId     string `json:",optional"`
	ClientSecret string `json:",optional"`
	CallbackUrl  string `json:",optional"`
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

// PlatformConf groups per-userType brand, front-end URLs, and one-time token expiry.
// JWT signing is shared (RpcConf.JwtConf); only brand / URLs / email differ per
// platform. Borrower emails must not point at test-environment URLs (ODD-7).
type PlatformConf struct {
	Borrower PlatformProfile
	Lender   PlatformProfile
	Admin    PlatformProfile
}

type PlatformProfile struct {
	// EmailConf is per-platform: SMTP creds and/or Gmail OAuth. A platform uses
	// whichever it is configured for (account emails: registration / reset / invite).
	EmailConf EmailConf

	// Front-end URLs injected into account emails.
	LoginUrl         string `json:",optional"`
	VerifyEmailUrl   string `json:",optional"`
	ResetPwdUrl      string `json:",optional"`
	InviteUserUrl    string `json:",optional"`
	BorrowerLoginUrl string `json:",optional"` // cross-platform link (e.g. borrower portal)
	LenderUrl        string `json:",optional"` // cross-platform link (e.g. lender portal)

	// One-time token / link expiry (seconds).
	VerifyEmailExpireSeconds     int `json:",default=86400"`
	ResetPwdExpireSeconds        int `json:",default=3600"`
	InviteUserExpireSeconds      int `json:",default=259200"`
	SendInviteEmailPeriodSeconds int `json:",default=60"`

	// Brand.
	LogoUrl       string `json:",optional"`
	FooterAddress string `json:",optional"`
	SenderName    string `json:",optional"`
	SenderEmail   string `json:",optional"`
}
