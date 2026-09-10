// Package testserver boots the real tenant gRPC services in-process for
// cross-service integration tests (borrower / admin characterization
// baselines). The caller shares its sqlite DB + miniredis and the JWT secrets,
// so a backend proxy call exercises the full backend→tenant hop against the
// same data and keys production uses.
package testserver

import (
	"context"
	"net"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis_rate/v10"
	"github.com/stretchr/testify/require"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"gorm.io/gorm"

	"zeebee-backend/app/tenant/config"
	"zeebee-backend/app/tenant/rpc/client"
	"zeebee-backend/app/tenant/rpc/internal/logic"
	apiServer "zeebee-backend/app/tenant/rpc/internal/server/api"
	authServer "zeebee-backend/app/tenant/rpc/internal/server/auth"
	organizationServer "zeebee-backend/app/tenant/rpc/internal/server/organization"
	organizationEmailServer "zeebee-backend/app/tenant/rpc/internal/server/organizationemail"
	organizationSettingServer "zeebee-backend/app/tenant/rpc/internal/server/organizationsetting"
	organizationSubscriptionServer "zeebee-backend/app/tenant/rpc/internal/server/organizationsubscription"
	organizationTemplateServer "zeebee-backend/app/tenant/rpc/internal/server/organizationtemplate"
	organizationUserServer "zeebee-backend/app/tenant/rpc/internal/server/organizationuser"
	pricingPlanServer "zeebee-backend/app/tenant/rpc/internal/server/pricingplan"
	pricingPlanPolicyServer "zeebee-backend/app/tenant/rpc/internal/server/pricingplanpolicy"
	resourceServer "zeebee-backend/app/tenant/rpc/internal/server/resource"
	roleServer "zeebee-backend/app/tenant/rpc/internal/server/role"
	userServer "zeebee-backend/app/tenant/rpc/internal/server/user"
	"zeebee-backend/app/tenant/rpc/internal/svc"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	commonlimiter "zeebee-backend/common/limiter"
	"zeebee-backend/common/rbac"
	"zeebee-backend/common/redisx"
	"zeebee-backend/common/types"
	domainConf "zeebee-backend/domains/zeebee/config"
	domainlimiter "zeebee-backend/domains/zeebee/limiter"
	"zeebee-backend/domains/zeebee/repo"
)

// Start builds a tenant ServiceContext over the caller's DB + redis, starts an
// in-process gRPC server with every tenant service registered, and returns a
// composite tenant client wired to it. The server and client are torn down
// when the test finishes. jwt must match the consuming backend's keys so tokens
// signed on either side parse on both (ODD-5 shared-secret period).
func Start(t *testing.T, db *gorm.DB, mr *miniredis.Miniredis, jwt types.JwtConf) *client.Client {
	t.Helper()
	logic.ResetBaseLogicForTest()

	redisCli, err := redisx.GetRedis(redisx.Opt{Host: mr.Addr()})
	require.NoError(t, err)
	enforcer, err := newEnforcer(t, db, mr)
	require.NoError(t, err)

	conf := testRpcConf()
	conf.JwtConf = jwt
	svcCtx := &svc.ServiceContext{
		Config:     conf,
		Logger:     logx.WithContext(context.Background()),
		AllRepos:   repo.NewAllRepos(db),
		RedisCli:   redisCli,
		Enforcer:   enforcer,
		LimiterMgr: newLimiterMgr(mr),
	}

	lis, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	grpcServer := grpc.NewServer()
	tenant.RegisterUserServer(grpcServer, userServer.NewUserServer(svcCtx))
	tenant.RegisterAuthServer(grpcServer, authServer.NewAuthServer(svcCtx))
	tenant.RegisterOrganizationServer(grpcServer, organizationServer.NewOrganizationServer(svcCtx))
	tenant.RegisterOrganizationUserServer(grpcServer, organizationUserServer.NewOrganizationUserServer(svcCtx))
	tenant.RegisterOrganizationEmailServer(grpcServer, organizationEmailServer.NewOrganizationEmailServer(svcCtx))
	tenant.RegisterOrganizationSettingServer(grpcServer, organizationSettingServer.NewOrganizationSettingServer(svcCtx))
	tenant.RegisterOrganizationSubscriptionServer(grpcServer, organizationSubscriptionServer.NewOrganizationSubscriptionServer(svcCtx))
	tenant.RegisterOrganizationTemplateServer(grpcServer, organizationTemplateServer.NewOrganizationTemplateServer(svcCtx))
	tenant.RegisterRoleServer(grpcServer, roleServer.NewRoleServer(svcCtx))
	tenant.RegisterResourceServer(grpcServer, resourceServer.NewResourceServer(svcCtx))
	tenant.RegisterApiServer(grpcServer, apiServer.NewApiServer(svcCtx))
	tenant.RegisterPricingPlanServer(grpcServer, pricingPlanServer.NewPricingPlanServer(svcCtx))
	tenant.RegisterPricingPlanPolicyServer(grpcServer, pricingPlanPolicyServer.NewPricingPlanPolicyServer(svcCtx))
	go func() { _ = grpcServer.Serve(lis) }()
	t.Cleanup(grpcServer.Stop)

	return client.New(zrpc.MustNewClient(zrpc.RpcClientConf{
		Endpoints: []string{lis.Addr().String()},
		NonBlock:  true,
	}))
}

// newEnforcer mirrors the tenant test enforcer: tenant is the central
// authority, so it loads all platform policies (no userType filter).
func newEnforcer(t *testing.T, db *gorm.DB, mr *miniredis.Miniredis) (*rbac.Enforcer, error) {
	t.Helper()
	modelPath := filepath.Join(repoRoot(t), "app", "tenant", "etc", "model.conf")
	modelPolicy, err := os.ReadFile(modelPath) //nolint:gosec // fixed test fixture inside the repo
	if err != nil {
		return nil, err
	}
	return rbac.NewEnforcer(rbac.EnforcerConfig{
		DB:          db,
		ModelPolicy: string(modelPolicy),
		Watcher: rbac.WatcherConfig{
			RedisHost: mr.Addr(),
			Channel:   "tenant-test-events",
		},
	})
}

func newLimiterMgr(mr *miniredis.Miniredis) *domainlimiter.LimiterManager {
	rdb, err := commonlimiter.GetLimiterRedisClient(commonlimiter.LimiterRedisConf{Host: mr.Addr()})
	if err != nil {
		panic(err)
	}
	per := func(rate, seconds int) redis_rate.Limit {
		return redis_rate.Limit{Rate: rate, Burst: rate, Period: time.Duration(seconds) * time.Second}
	}
	return domainlimiter.NewLimiterManager(commonlimiter.NewRedisLimiterWithRedis(rdb), map[domainlimiter.LimiterStrategy]redis_rate.Limit{
		domainlimiter.StrategyLogin:             per(10, 300),
		domainlimiter.StrategySendLoginEmail:    per(5, 300),
		domainlimiter.StrategySendVerifyEmail:   per(5, 300),
		domainlimiter.StrategySendResetPwdEmail: per(5, 300),
	})
}

func testRpcConf() config.RpcConf {
	profile := func() config.PlatformProfile {
		return config.PlatformProfile{
			LoginUrl:                     "http://test.local/login",
			VerifyEmailUrl:               "http://test.local/verify-email",
			ResetPwdUrl:                  "http://test.local/reset-password",
			InviteUserUrl:                "http://test.local/invite",
			VerifyEmailExpireSeconds:     86400,
			ResetPwdExpireSeconds:        3600,
			InviteUserExpireSeconds:      259200,
			SendInviteEmailPeriodSeconds: 60,
			SenderName:                   "Zeebee",
			SenderEmail:                  "info@test.local",
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
			Borrower: profile(),
			Lender:   profile(),
			Admin:    profile(),
		},
	}
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
