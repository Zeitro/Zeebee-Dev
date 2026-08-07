// Package testutil wires a full lender ServiceContext for characterization
// baseline tests: in-memory SQLite (shared domain tables), miniredis (verify
// codes / limiter / casbin watcher), and the real Casbin enforcer filtered to
// lender policies — matching how lender runs in production. No external
// services required.
//
// Unlike tenant (svcCtx-injected), lender logic reads package-level globals
// (Config / GetRedis / LenderEnforcer / EmailClient / LimiterMgr), so NewEnv
// injects the isolated environment into those globals too.
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
	"github.com/casbin/gorm-adapter/v3"
	"github.com/go-redis/redis_rate/v10"
	"github.com/stretchr/testify/require"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"

	"zeebee-backend/app/lender/config"
	"zeebee-backend/app/lender/rpc/internal/global"
	"zeebee-backend/app/lender/rpc/internal/logic"
	"zeebee-backend/app/lender/rpc/internal/svc"
	"zeebee-backend/common/email"
	commonlimiter "zeebee-backend/common/limiter"
	"zeebee-backend/common/rbac"
	"zeebee-backend/common/redisx"
	"zeebee-backend/common/types"
	domainConf "zeebee-backend/domains/zeebee/config"
	domainlimiter "zeebee-backend/domains/zeebee/limiter"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/repo"
)

// Env is a self-contained lender test environment.
type Env struct {
	SvcCtx *svc.ServiceContext
	Redis  *miniredis.Miniredis
	DB     *gorm.DB
}

const lenderChannel = "lender-test-events"

// NewEnv builds a wired lender ServiceContext backed by sqlite + miniredis and
// injects it into the lender global singletons. Each test gets an isolated
// in-memory DB (unique DSN from the test name) and its own baseLogic.
func NewEnv(t *testing.T) *Env {
	t.Helper()
	mr := miniredis.RunT(t)
	conf := TestRpcConf()

	dsn := fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, rewritePgDefaultsForSQLite(db))
	require.NoError(t, db.AutoMigrate(lenderModels()...))

	redisCli, err := redisx.GetRedis(redisx.Opt{Host: mr.Addr(), Prefix: "lender-test"})
	require.NoError(t, err)

	enforcer, err := newTestEnforcer(t, db, mr)
	require.NoError(t, err)

	limiterMgr := newTestLimiterMgr(mr)

	// Inject the isolated environment into the package globals lender logic reads.
	global.SetGlobalConfig(conf)
	global.SetRedis(redisCli)
	global.SetGormDB(db)
	global.SetLenderEnforcer(enforcer)
	global.SetEmailClient(email.NewEmail(email.SMTPConfig{}))
	global.SetLimiterMgr(limiterMgr)

	// baseLogic is a sync.Once singleton; reset so each test rebuilds it.
	logic.ResetBaseLogicForTest()

	svcCtx := &svc.ServiceContext{
		Config:   conf,
		Logger:   logx.WithContext(context.Background()),
		AllRepos: repo.NewAllRepos(db),
	}
	return &Env{SvcCtx: svcCtx, Redis: mr, DB: db}
}

// rewritePgDefaultsForSQLite makes AutoMigrate work on SQLite: the shared Base
// model uses PostgreSQL-style default:now() tags sqlite cannot parse. Mirrors
// the tenant testutil normalization.
func rewritePgDefaultsForSQLite(db *gorm.DB) error {
	return db.Callback().Raw().Before("gorm:raw").Register("sqlite:default-now", func(tx *gorm.DB) {
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

func newTestEnforcer(t *testing.T, db *gorm.DB, mr *miniredis.Miniredis) (*rbac.Enforcer, error) {
	t.Helper()
	modelPath := filepath.Join(repoRoot(t), "app", "lender", "etc", "model.conf")
	modelPolicy, err := os.ReadFile(modelPath) //nolint:gosec // test reads a fixed fixture inside the repo
	if err != nil {
		return nil, err
	}
	return rbac.NewEnforcer(rbac.EnforcerConfig{
		DB:          db,
		ModelPolicy: string(modelPolicy),
		Watcher: rbac.WatcherConfig{
			RedisHost: mr.Addr(),
			Channel:   lenderChannel,
		},
		// Mirror production: lender enforcer only loads lender (v1) policies.
		DefaultLoadFilteredPolicy: &gormadapter.Filter{
			Ptype: []string{"p"},
			V1:    []string{"lender"},
		},
	})
}

func newTestLimiterMgr(mr *miniredis.Miniredis) *domainlimiter.LimiterManager {
	rdb, err := commonlimiter.GetLimiterRedisClient(commonlimiter.LimiterRedisConf{Host: mr.Addr()})
	if err != nil {
		panic(err)
	}
	per := func(rate int) redis_rate.Limit {
		return redis_rate.Limit{Rate: rate, Burst: rate, Period: 300 * time.Second}
	}
	return domainlimiter.NewLimiterManager(
		commonlimiter.NewRedisLimiterWithRedis(rdb),
		map[domainlimiter.LimiterStrategy]redis_rate.Limit{
			domainlimiter.StrategyLogin:             per(10),
			domainlimiter.StrategySendLoginEmail:    per(5),
			domainlimiter.StrategySendVerifyEmail:   per(5),
			domainlimiter.StrategySendResetPwdEmail: per(5),
		},
	)
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

// lenderModels lists the shared-domain tables lender identity / RBAC baselines
// touch. jsonb columns are TEXT affinity in sqlite; gorm serializers handle
// JSON round-trips.
func lenderModels() []any {
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

// TestRpcConf returns a lender config usable in tests: shared JWT key, lender
// casbin event topic, limiter items, and expiry / URL fields the logic reads.
// No SMTP (email client is a no-op stub set in NewEnv).
func TestRpcConf() config.RpcConf {
	return config.RpcConf{
		JwtConf: types.JwtConf{
			AccessSecret:  "test-access-secret",
			AccessExpire:  3600,
			RefreshSecret: "test-refresh-secret",
			RefreshExpire: 86400,
		},
		LimiterConf: config.LimiterConf{
			Login:             domainConf.LimiterItem{Rate: 10, Period: 300},
			SendLoginEmail:    domainConf.LimiterItem{Rate: 5, Period: 300},
			SendVerityEmail:   domainConf.LimiterItem{Rate: 5, Period: 300},
			SendResetPwdEmail: domainConf.LimiterItem{Rate: 5, Period: 300},
		},
		VerifyEmailExpireSeconds: 86400,
		ResetPwdExpireSeconds:    3600,
		LoginUrl:                 "http://test.local/login",
		VerifyEmailUrl:           "http://test.local/verify-email",
		ResetPwdUrl:              "http://test.local/reset-password",
		CasbinConf: config.CasbinConf{
			Watcher:    config.CasbinWatcher{LenderChannel: lenderChannel},
			EventTopic: config.CasbinEventTopic{LenderTopic: lenderChannel},
		},
	}
}
