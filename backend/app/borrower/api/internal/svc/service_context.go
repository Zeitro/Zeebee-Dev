package svc

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"

	"zeebee-backend/app/borrower/api/internal/config"
	"zeebee-backend/app/borrower/api/internal/middleware"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
)

type ServiceContext struct {
	Config config.Config
	// middleware
	AuthMiddleware rest.Middleware

	// rpc client
	RpcClients RpcClients
}

type RpcClients struct {
	UserRpcClient borrower.UserClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:         c,
		AuthMiddleware: middleware.NewAuthMiddleware(c.JwtAuth.AccessSecret).Handle,
		RpcClients: RpcClients{
			UserRpcClient: borrower.NewUserClient(zrpc.MustNewClient(c.RpcConf.UserConf).Conn()),
		},
	}
}
