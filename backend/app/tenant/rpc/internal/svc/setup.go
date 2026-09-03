package svc

import (
	"os"
	"time"

	"github.com/go-redis/redis_rate/v10"
	"github.com/zeromicro/go-zero/core/service"
	"gorm.io/gorm"

	"zeebee-backend/app/tenant/config"
	"zeebee-backend/common/gormx"
	commonlimiter "zeebee-backend/common/limiter"
	"zeebee-backend/common/rbac"
	"zeebee-backend/common/redisx"
	"zeebee-backend/domains/zeebee/limiter"
)

// proMode reports whether the service runs in production (debug off for gorm).
func proMode(c config.RpcConf) bool { return c.Mode == service.ProMode }

func newGormDB(dsn string, debug bool) (*gorm.DB, error) {
	return gormx.GetDB(dsn, debug, gormx.DBTypePgsql)
}

func newRedis(c config.RedisConf) (*redisx.Client, error) {
	return redisx.GetRedis(redisx.Opt{
		Host:   c.Host,
		Type:   c.Type,
		Pass:   c.Pass,
		Tls:    c.Tls,
		DB:     c.DB,
		Prefix: c.Prefix,
	})
}

// newEnforcer builds the tenant Casbin enforcer. tenant is the central authority,
// so it loads ALL platform policies (no userType filter).
func newEnforcer(c config.CasbinConf, db *gorm.DB) (*rbac.Enforcer, error) {
	policyContent, err := os.ReadFile(c.ModelPath)
	if err != nil {
		return nil, err
	}
	return rbac.NewEnforcer(rbac.EnforcerConfig{
		DB:          db,
		ModelPolicy: string(policyContent),
		Watcher: rbac.WatcherConfig{
			RedisHost: c.Watcher.Host,
			Pwd:       c.Watcher.Pwd,
			Tls:       c.Watcher.Tls,
			Channel:   c.Watcher.Channel,
		},
		DefaultLoadFilteredPolicy: nil,
	})
}

func newLimiterMgr(c config.RpcConf) (*limiter.LimiterManager, error) {
	rdb, err := commonlimiter.GetLimiterRedisClient(commonlimiter.LimiterRedisConf{
		Host: c.RedisConf.Host,
		Pass: c.RedisConf.Pass,
		DB:   c.RedisConf.DB,
		Tls:  c.RedisConf.Tls,
	})
	if err != nil {
		return nil, err
	}

	lc := c.LimiterConf
	return limiter.NewLimiterManager(commonlimiter.NewRedisLimiterWithRedis(rdb), map[limiter.LimiterStrategy]redis_rate.Limit{
		limiter.StrategyLogin: {
			Rate:   lc.Login.Rate,
			Burst:  lc.Login.Rate,
			Period: time.Duration(int64(lc.Login.Period)) * time.Second,
		},
		limiter.StrategySendLoginEmail: {
			Rate:   lc.SendLoginEmail.Rate,
			Burst:  lc.SendLoginEmail.Rate,
			Period: time.Duration(int64(lc.SendLoginEmail.Period)) * time.Second,
		},
		limiter.StrategySendVerifyEmail: {
			Rate:   lc.SendVerifyEmail.Rate,
			Burst:  lc.SendVerifyEmail.Rate,
			Period: time.Duration(int64(lc.SendVerifyEmail.Period)) * time.Second,
		},
		limiter.StrategySendResetPwdEmail: {
			Rate:   lc.SendResetPwdEmail.Rate,
			Burst:  lc.SendResetPwdEmail.Rate,
			Period: time.Duration(int64(lc.SendResetPwdEmail.Period)) * time.Second,
		},
	}), nil
}
