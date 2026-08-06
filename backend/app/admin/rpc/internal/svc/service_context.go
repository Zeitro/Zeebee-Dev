package svc

import (
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/zrpc"

	"zeebee-backend/app/admin/config"
	"zeebee-backend/app/admin/rpc/internal/global"
	tenantauth "zeebee-backend/app/tenant/rpc/client/auth"
	tenantuser "zeebee-backend/app/tenant/rpc/client/user"
	"zeebee-backend/common/llm"
	"zeebee-backend/domains/zeebee/repo"
)

type ServiceContext struct {
	Logger logx.Logger
	Config config.RpcConf
	*repo.AllRepos
	// DiscoverLLM is the lender-discover LLM client (built from LlmConf).
	// nil if the provider is unconfigured; discover calls then fail at runtime.
	DiscoverLLM llm.Client
	// TenantUser / TenantAuth are outbound clients to the tenant service
	// (identity + authz authority, direction B); used by the US-12 cutover.
	TenantUser tenantuser.User
	TenantAuth tenantauth.Auth
}

func NewServiceContext(c config.RpcConf, logger logx.Logger) *ServiceContext {
	discoverLLM, err := llm.NewClient(c.LlmConf.Providers, c.LlmConf.LenderDiscoverConf, c.Mode)
	if err != nil {
		logger.Errorf("lender discover llm client not initialized: %v", err)
	}
	tenantCli := zrpc.MustNewClient(c.RpcServiceConf.Tenant)
	return &ServiceContext{
		Config:      c,
		Logger:      logger,
		AllRepos:    repo.NewAllRepos(global.GetGormDB()),
		DiscoverLLM: discoverLLM,
		TenantUser:  tenantuser.NewUser(tenantCli),
		TenantAuth:  tenantauth.NewAuth(tenantCli),
	}
}
