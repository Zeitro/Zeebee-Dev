package svc

import (
	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/tenant/config"
	"zeebee-backend/common/rbac"
	"zeebee-backend/common/redisx"
	"zeebee-backend/domains/zeebee/limiter"
	"zeebee-backend/domains/zeebee/repo"
)

// ServiceContext carries shared connections as injected fields (no globals).
// AllLogics construction lives in internal/logic/base_logic.go (mirrors borrower);
// svc only owns infra (gorm / redis / casbin / limiter) + repos.
type ServiceContext struct {
	Config config.RpcConf
	Logger logx.Logger
	*repo.AllRepos
	RedisCli   *redisx.Client
	Enforcer   *rbac.Enforcer
	LimiterMgr *limiter.LimiterManager
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

	return &ServiceContext{
		Config:     c,
		Logger:     logger,
		AllRepos:   repo.NewAllRepos(gormDB),
		RedisCli:   redisCli,
		Enforcer:   enforcer,
		LimiterMgr: limiterMgr,
	}, nil
}
