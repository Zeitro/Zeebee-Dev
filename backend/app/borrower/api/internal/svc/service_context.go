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
	UserRpcClient borrower.UserClient
	LoanRpcClient borrower.LoanClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:         c,
		AuthMiddleware: middleware.NewAuthMiddleware(c.JwtAuth.AccessSecret).Handle,
		// rpc client
		UserRpcClient: borrower.NewUserClient(zrpc.MustNewClient(c.RpcConf).Conn()),
		LoanRpcClient: borrower.NewLoanClient(zrpc.MustNewClient(c.RpcConf).Conn()),
	}
}
