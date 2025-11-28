package svc

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"

	"zeebee-backend/app/admin/api/internal/middleware"
	"zeebee-backend/app/admin/config"
	"zeebee-backend/app/admin/rpc/pb/admin"
)

type ServiceContext struct {
	Config config.ApiConf
	// middleware
	AuthMiddleware rest.Middleware

	// rpc client
	UserRpcClient     admin.UserClient
	ApiRpcClient      admin.ApiClient
	ResourceRpcClient admin.ResourceClient
}

func NewServiceContext(c config.ApiConf) *ServiceContext {
	return &ServiceContext{
		Config:         c,
		AuthMiddleware: middleware.NewAuthMiddleware(c.JwtAuth.AccessSecret).Handle,

		// rpc client
		UserRpcClient:     admin.NewUserClient(zrpc.MustNewClient(c.RpcConf).Conn()),
		ApiRpcClient:      admin.NewApiClient(zrpc.MustNewClient(c.RpcConf).Conn()),
		ResourceRpcClient: admin.NewResourceClient(zrpc.MustNewClient(c.RpcConf).Conn()),
	}
}
