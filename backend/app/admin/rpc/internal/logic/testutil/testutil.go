// Package testutil wires a full admin ServiceContext for characterization
// baseline tests: in-memory SQLite (shared domain tables), miniredis (verify
// codes / limiter / casbin watcher), and the real Casbin enforcer. No external
// services required.
//
// Like lender, admin logic reads package-level globals (Config / GetRedis /
// AdminEnforcer / EmailClient), so NewEnv injects the isolated environment
// into those globals too.
package testutil

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/require"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"

	"zeebee-backend/app/admin/config"
	"zeebee-backend/app/admin/rpc/internal/global"
	"zeebee-backend/app/admin/rpc/internal/logic"
	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/common/email"
	"zeebee-backend/common/rbac"
	"zeebee-backend/common/redisx"
	"zeebee-backend/common/types"
	domainConf "zeebee-backend/domains/zeebee/config"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/repo"
)

const adminChannel = "admin-test-events"

// Env is a self-contained admin test environment.
type Env struct {
	SvcCtx *svc.ServiceContext
	Redis  *miniredis.Miniredis
	DB     *gorm.DB
}

// NewEnv builds a wired admin ServiceContext backed by sqlite + miniredis and
// injects it into the admin global singletons. Each test gets an isolated
// in-memory DB (unique DSN from the test name) and its own baseLogic.
func NewEnv(t *testing.T) *Env {
	t.Helper()
	mr := miniredis.RunT(t)
	conf := TestRpcConf()

	dsn := fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, rewritePgDefaultsForSQLite(db))
	require.NoError(t, db.AutoMigrate(adminModels()...))

	redisCli, err := redisx.GetRedis(redisx.Opt{Host: mr.Addr(), Prefix: "admin-test"})
	require.NoError(t, err)

	enforcer, err := newTestEnforcer(t, db, mr)
	require.NoError(t, err)

	// Inject the isolated environment into the package globals admin logic reads.
	global.SetGlobalConfig(conf)
	global.SetRedis(redisCli)
	global.SetGormDB(db)
	global.SetAdminEnforcer(enforcer)
	global.SetEmailClient(email.NewEmail(email.SMTPConfig{}))

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
// model uses PostgreSQL-style default:now() tags sqlite cannot parse.
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
	modelPath := filepath.Join(repoRoot(t), "app", "admin", "etc", "model.conf")
	modelPolicy, err := os.ReadFile(modelPath) //nolint:gosec // test reads a fixed fixture inside the repo
	if err != nil {
		return nil, err
	}
	return rbac.NewEnforcer(rbac.EnforcerConfig{
		DB:          db,
		ModelPolicy: string(modelPolicy),
		Watcher: rbac.WatcherConfig{
			RedisHost: mr.Addr(),
			Channel:   adminChannel,
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

// adminModels lists the shared-domain tables admin identity baselines touch.
func adminModels() []any {
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

// TestRpcConf returns an admin config usable in tests: shared JWT key, admin
// casbin event topic, login limiter item. No SMTP (email client is a no-op
// stub set in NewEnv).
func TestRpcConf() config.RpcConf {
	return config.RpcConf{
		JwtConf: types.JwtConf{
			AccessSecret:  "test-access-secret",
			AccessExpire:  3600,
			RefreshSecret: "test-refresh-secret",
			RefreshExpire: 86400,
		},
		LimiterConf: config.LimiterConf{
			Login: domainConf.LimiterItem{Rate: 10, Period: 300},
		},
		CasbinConf: config.CasbinConf{
			Watcher:    config.CasbinWatcher{AdminChannel: adminChannel},
			EventTopic: config.CasbinEventTopic{AdminTopic: adminChannel},
		},
	}
}
