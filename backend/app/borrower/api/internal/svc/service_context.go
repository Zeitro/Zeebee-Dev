package svc

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"

	"zeebee-backend/app/borrower/api/internal/middleware"
	"zeebee-backend/app/borrower/config"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
)

type ServiceContext struct {
	Config config.ApiConf
	// middleware
	AuthMiddleware rest.Middleware

	// rpc client
	UserRpcClient      borrower.UserClient
	LoanRpcClient      borrower.LoanClient
	CommonRpcClient    borrower.CommonClient
	PropertyRpcClient  borrower.PropertyClient
	ApplicantRpcClient borrower.ApplicantClient
}

func NewServiceContext(c config.ApiConf) *ServiceContext {
	return &ServiceContext{
		Config:         c,
		AuthMiddleware: middleware.NewAuthMiddleware(c.JwtAuth.AccessSecret).Handle,
		// rpc client
		UserRpcClient:      borrower.NewUserClient(zrpc.MustNewClient(c.RpcConf).Conn()),
		LoanRpcClient:      borrower.NewLoanClient(zrpc.MustNewClient(c.RpcConf).Conn()),
		CommonRpcClient:    borrower.NewCommonClient(zrpc.MustNewClient(c.RpcConf).Conn()),
		PropertyRpcClient:  borrower.NewPropertyClient(zrpc.MustNewClient(c.RpcConf).Conn()),
		ApplicantRpcClient: borrower.NewApplicantClient(zrpc.MustNewClient(c.RpcConf).Conn()),
	}
}
