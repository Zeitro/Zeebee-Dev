package svc

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"

	"zeebee-backend/app/lender/api/internal/middleware"
	"zeebee-backend/app/lender/config"
	"zeebee-backend/app/lender/rpc/pb/lender"
)

type ServiceContext struct {
	Config config.ApiConf
	// middleware
	AuthMiddleware rest.Middleware

	// rpc client
	UserRpcClient    lender.UserClient
	LoanRpcClient    lender.LoanClient
	CommonRpcClient  lender.CommonClient
	StorageRpcClient lender.StorageClient
}

func NewServiceContext(c config.ApiConf) *ServiceContext {
	return &ServiceContext{
		Config:         c,
		AuthMiddleware: middleware.NewAuthMiddleware(c.JwtAuth.AccessSecret).Handle,

		// rpc client
		UserRpcClient:    lender.NewUserClient(zrpc.MustNewClient(c.RpcConf).Conn()),
		LoanRpcClient:    lender.NewLoanClient(zrpc.MustNewClient(c.RpcConf).Conn()),
		CommonRpcClient:  lender.NewCommonClient(zrpc.MustNewClient(c.RpcConf).Conn()),
		StorageRpcClient: lender.NewStorageClient(zrpc.MustNewClient(c.RpcConf).Conn()),
	}
}
