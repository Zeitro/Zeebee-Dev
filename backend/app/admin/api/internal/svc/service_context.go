package svc

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"

	"zeebee-backend/app/admin/config"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/middleware"
)

type ServiceContext struct {
	Config config.ApiConf
	// middleware
	AuthMiddleware rest.Middleware

	// rpc client
	UserRpcClient             admin.UserClient
	ApiRpcClient              admin.ApiClient
	ResourceRpcClient         admin.ResourceClient
	RoleRpcClient             admin.RoleClient
	OrganizationUserRpcClient admin.OrganizationUserClient
}

func NewServiceContext(c config.ApiConf) *ServiceContext {
	return &ServiceContext{
		Config:         c,
		AuthMiddleware: middleware.NewHttpAuthMiddleware(c.JwtAuth.AccessSecret, constants.UserTypeAdmin).Handle,

		// rpc client
		UserRpcClient:             admin.NewUserClient(zrpc.MustNewClient(c.RpcConf).Conn()),
		ApiRpcClient:              admin.NewApiClient(zrpc.MustNewClient(c.RpcConf).Conn()),
		ResourceRpcClient:         admin.NewResourceClient(zrpc.MustNewClient(c.RpcConf).Conn()),
		RoleRpcClient:             admin.NewRoleClient(zrpc.MustNewClient(c.RpcConf).Conn()),
		OrganizationUserRpcClient: admin.NewOrganizationUserClient(zrpc.MustNewClient(c.RpcConf).Conn()),
	}
}
