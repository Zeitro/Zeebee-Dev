package svc

import (
	"github.com/zeromicro/go-zero/zrpc"

	"zeebee-backend/app/lender/config"
	"zeebee-backend/app/lender/rpc/pb/lender"
)

type ServiceContext struct {
	Config config.ApiConf

	// rpc client
	UserRpcClient lender.UserClient
	// LoanRpcClient      borrower.LoanClient
	CommonRpcClient lender.CommonClient
}

func NewServiceContext(c config.ApiConf) *ServiceContext {
	return &ServiceContext{
		Config: c,

		// rpc client
		UserRpcClient:   lender.NewUserClient(zrpc.MustNewClient(c.RpcConf).Conn()),
		CommonRpcClient: lender.NewCommonClient(zrpc.MustNewClient(c.RpcConf).Conn()),
	}
}
