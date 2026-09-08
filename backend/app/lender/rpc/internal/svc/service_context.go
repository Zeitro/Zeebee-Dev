package svc

import (
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/zrpc"

	"zeebee-backend/app/lender/config"
	"zeebee-backend/app/lender/rpc/internal/global"
	tenantclient "zeebee-backend/app/tenant/rpc/client"
	"zeebee-backend/domains/zeebee/repo"
)

type ServiceContext struct {
	Config config.RpcConf
	Logger logx.Logger
	*repo.AllRepos
	// TenantClient bundles tenant's outbound RPC clients (identity + authz +
	// org/role/resource/...); used by the US-12 cutover.
	TenantClient *tenantclient.Client
}

func NewServiceContext(c config.RpcConf, logger logx.Logger) *ServiceContext {
	tenantCli := zrpc.MustNewClient(c.RpcServiceConf.Tenant)
	return &ServiceContext{
		Config:     c,
		Logger:     logger,
		AllRepos:     repo.NewAllRepos(global.GetGormDB()),
		TenantClient: tenantclient.New(tenantCli),
	}
}
