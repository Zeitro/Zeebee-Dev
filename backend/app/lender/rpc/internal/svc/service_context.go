package svc

import (
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/zrpc"

	"zeebee-backend/app/lender/config"
	"zeebee-backend/app/lender/rpc/internal/global"
	tenantauth "zeebee-backend/app/tenant/rpc/client/auth"
	tenantuser "zeebee-backend/app/tenant/rpc/client/user"
	"zeebee-backend/domains/zeebee/repo"
)

type ServiceContext struct {
	Config config.RpcConf
	Logger logx.Logger
	*repo.AllRepos
	// TenantUser / TenantAuth are outbound clients to the tenant service
	// (identity + authz authority, direction B); used by the US-12 cutover.
	TenantUser tenantuser.User
	TenantAuth tenantauth.Auth
}

func NewServiceContext(c config.RpcConf, logger logx.Logger) *ServiceContext {
	tenantCli := zrpc.MustNewClient(c.RpcServiceConf.Tenant)
	return &ServiceContext{
		Config:     c,
		Logger:     logger,
		AllRepos:   repo.NewAllRepos(global.GetGormDB()),
		TenantUser: tenantuser.NewUser(tenantCli),
		TenantAuth: tenantauth.NewAuth(tenantCli),
	}
}
