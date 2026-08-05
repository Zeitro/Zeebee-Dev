// Package testutil wires a full tenant ServiceContext for integration tests:
// in-memory SQLite (shared domain tables), miniredis (codes / limiter / casbin
// watcher), and the real Casbin enforcer. No external services required.
package testutil

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/glebarez/sqlite"
	"github.com/go-redis/redis_rate/v10"
	"github.com/stretchr/testify/require"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"

	"zeebee-backend/app/tenant/config"
	"zeebee-backend/app/tenant/rpc/internal/svc"
	commonlimiter "zeebee-backend/common/limiter"
	"zeebee-backend/common/rbac"
	"zeebee-backend/common/redisx"
	"zeebee-backend/common/types"
	domainConf "zeebee-backend/domains/zeebee/config"
	domainlimiter "zeebee-backend/domains/zeebee/limiter"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/repo"
)

// Env is a self-contained tenant test environment.
type Env struct {
	SvcCtx *svc.ServiceContext
	Redis  *miniredis.Miniredis
	DB     *gorm.DB
}

// NewEnv builds a wired tenant ServiceContext backed by sqlite + miniredis.
// Each test gets an isolated in-memory DB (unique DSN from the test name).
func NewEnv(t *testing.T) *Env {
	t.Helper()
	mr := miniredis.RunT(t)
	conf := TestRpcConf()

	dsn := fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	require.NoError(t, err)
	rewritePgDefaultsForSQLite(db)
	require.NoError(t, db.AutoMigrate(tenantModels()...))

	redisCli, err := redisx.GetRedis(redisx.Opt{Host: mr.Addr(), Prefix: "tenant-test"})
	require.NoError(t, err)

	enforcer, err := newTestEnforcer(t, db, mr, conf.CasbinConf.Watcher.Channel)
	require.NoError(t, err)

	limiterRdb, err := commonlimiter.GetLimiterRedisClient(commonlimiter.LimiterRedisConf{Host: mr.Addr()})
	require.NoError(t, err)
	limiterMgr := domainlimiter.NewLimiterManager(
		commonlimiter.NewRedisLimiterWithRedis(limiterRdb),
		testLimiterRates(),
	)

	svcCtx := &svc.ServiceContext{
		Config:     conf,
		Logger:     logx.WithContext(context.Background()),
		AllRepos:   repo.NewAllRepos(db),
		RedisCli:   redisCli,
		Enforcer:   enforcer,
		LimiterMgr: limiterMgr,
	}
	return &Env{SvcCtx: svcCtx, Redis: mr, DB: db}
}

// rewritePgDefaultsForSQLite makes AutoMigrate work on SQLite: the shared Base
// model uses PostgreSQL-style default:now() tags, which sqlite cannot parse.
// GORM resolves the full CREATE TABLE DDL before the gorm:raw callback runs, so
// we normalize the default there instead of touching production models.
func rewritePgDefaultsForSQLite(db *gorm.DB) {
	db.Callback().Raw().Before("gorm:raw").Register("sqlite:default-now", func(tx *gorm.DB) {
		if tx.Statement.SQL.Len() == 0 {
			return
		}
		sql := tx.Statement.SQL.String()
		if strings.Contains(sql, "DEFAULT now()") {
			sql = strings.ReplaceAll(sql, "DEFAULT now()", "DEFAULT CURRENT_TIMESTAMP")
			tx.Statement.SQL.Reset()
			tx.Statement.SQL.WriteString(sql)
		}
	})
}

func newTestEnforcer(t *testing.T, db *gorm.DB, mr *miniredis.Miniredis, channel string) (*rbac.Enforcer, error) {
	t.Helper()
	modelPath := filepath.Join(repoRoot(t), "app", "tenant", "etc", "model.conf")
	modelPolicy, err := os.ReadFile(modelPath)
	if err != nil {
		return nil, err
	}
	return rbac.NewEnforcer(rbac.EnforcerConfig{
		DB:          db,
		ModelPolicy: string(modelPolicy),
		Watcher: rbac.WatcherConfig{
			RedisHost: mr.Addr(),
			Channel:   channel,
		},
	})
}

// repoRoot walks up from the package dir until the go.mod of zeebee-backend.
func repoRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	require.NoError(t, err)
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("repo root not found")
		}
		dir = parent
	}
}

// tenantModels lists the shared-domain tables tenant touches. jsonb columns are
// TEXT affinity in sqlite; gorm serializers handle JSON round-trips.
func tenantModels() []any {
	return []any{
		&model.User{},
		&model.Organization{},
		&model.OrganizationUser{},
		&model.Role{},
		&model.Resource{},
		&model.Api{},
		&model.RoleAndResource{},
		&model.PricingPlan{},
		&model.PricingPlanPolicy{},
		&model.PricingPlanAndResource{},
		&model.OrganizationSubscription{},
		&model.OrganizationSetting{},
		&model.OrganizationTemplate{},
		&model.BaseTemplate{},
		&model.OrganizationEmail{},
		&model.EmailAuthorization{},
	}
}

func testLimiterRates() map[domainlimiter.LimiterStrategy]redis_rate.Limit {
	per := func(rate, seconds int) redis_rate.Limit {
		return redis_rate.Limit{Rate: rate, Burst: rate, Period: time.Duration(seconds) * time.Second}
	}
	return map[domainlimiter.LimiterStrategy]redis_rate.Limit{
		domainlimiter.StrategyLogin:          per(10, 300),
		domainlimiter.StrategySendLoginEmail: per(5, 300),
		domainlimiter.StrategySendVerifyEmail: per(5, 300),
		domainlimiter.StrategySendResetPwdEmail: per(5, 300),
	}
}

// TestRpcConf returns a tenant config with per-platform JWT / URL / expiry
// profiles and no SMTP (tests exercise DB + redis paths, not real sending).
func TestRpcConf() config.RpcConf {
	profile := func(accessSecret, refreshSecret string) config.PlatformProfile {
		return config.PlatformProfile{
			JwtAuth: types.JwtAuth{
				AccessSecret:  accessSecret,
				AccessExpire:  3600,
				RefreshSecret: refreshSecret,
				RefreshExpire: 86400,
			},
			LoginUrl:                    "http://test.local/login",
			VerifyEmailUrl:              "http://test.local/verify-email",
			ResetPwdUrl:                 "http://test.local/reset-password",
			InviteUserUrl:               "http://test.local/invite",
			VerifyEmailExpireSeconds:    86400,
			ResetPwdExpireSeconds:       3600,
			InviteUserExpireSeconds:     259200,
			SendInviteEmailPeriodSeconds: 60,
			SenderName:                  "Zeebee",
			SenderEmail:                 "info@test.local",
		}
	}
	return config.RpcConf{
		CasbinConf: config.CasbinConf{
			Watcher: config.CasbinWatcher{Channel: "tenant-test-events"},
		},
		LimiterConf: config.LimiterConf{
			Login:             domainConf.LimiterItem{Rate: 10, Period: 300, LockPeriod: 120},
			SendLoginEmail:    domainConf.LimiterItem{Rate: 5, Period: 300},
			SendVerifyEmail:   domainConf.LimiterItem{Rate: 5, Period: 300},
			SendResetPwdEmail: domainConf.LimiterItem{Rate: 5, Period: 300},
		},
		PlatformConf: config.PlatformConf{
			Borrower: profile("borrower-access-secret", "borrower-refresh-secret"),
			Lender:   profile("lender-access-secret", "lender-refresh-secret"),
			Admin:    profile("admin-access-secret", "admin-refresh-secret"),
		},
	}
}
