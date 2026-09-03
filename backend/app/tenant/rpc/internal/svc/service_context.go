package svc

import (
	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/tenant/config"
	"zeebee-backend/common/cron"
	"zeebee-backend/common/llm"
	"zeebee-backend/common/rbac"
	"zeebee-backend/common/redisx"
	"zeebee-backend/domains/zeebee/limiter"
	"zeebee-backend/domains/zeebee/repo"
)

// ServiceContext carries shared connections as injected fields (no globals).
// AllLogics construction lives in internal/logic/base_logic.go (mirrors borrower);
// svc owns infra (gorm / redis / casbin / limiter), repos, and discover runtime dependencies.
type ServiceContext struct {
	Config config.RpcConf
	Logger logx.Logger
	*repo.AllRepos
	RedisCli    *redisx.Client
	Enforcer    *rbac.Enforcer
	LimiterMgr  *limiter.LimiterManager
	DiscoverLLM llm.Client
	Cron        *cron.Engine
}

// NewServiceContext builds gorm / redis / casbin / limiter from config and returns
// a fully wired ServiceContext. Connections are owned here, passed by parameter.
func NewServiceContext(c config.RpcConf, logger logx.Logger) (*ServiceContext, error) {
	gormDB, err := newGormDB(c.DbConf.Dsn, !proMode(c))
	if err != nil {
		return nil, err
	}

	redisCli, err := newRedis(c.RedisConf)
	if err != nil {
		return nil, err
	}

	enforcer, err := newEnforcer(c.CasbinConf, gormDB)
	if err != nil {
		return nil, err
	}

	limiterMgr, err := newLimiterMgr(c)
	if err != nil {
		return nil, err
	}
	discoverLLM, err := llm.NewClient(c.LlmConf.Providers, c.LlmConf.LenderDiscoverConf, c.Mode)
	if err != nil {
		logger.Errorf("lender discover llm client not initialized: %v", err)
	}

	return &ServiceContext{
		Config:      c,
		Logger:      logger,
		AllRepos:    repo.NewAllRepos(gormDB),
		RedisCli:    redisCli,
		Enforcer:    enforcer,
		LimiterMgr:  limiterMgr,
		DiscoverLLM: discoverLLM,
	}, nil
}
