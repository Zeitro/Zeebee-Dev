package svc

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"

	"zeebee-backend/app/lender/config"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/constants"
	"zeebee-backend/domains/zeebee/middleware"
)

type ServiceContext struct {
	Config config.ApiConf
	// middleware
	AuthMiddleware rest.Middleware

	// rpc client
	UserRpcClient         lender.UserClient
	LoanRpcClient         lender.LoanClient
	CommonRpcClient       lender.CommonClient
	StorageRpcClient      lender.StorageClient
	FileRpcClient         lender.FileClient
	LoanDocumentRpcClient lender.LoanDocumentClient
}

func NewServiceContext(c config.ApiConf) *ServiceContext {
	return &ServiceContext{
		Config:         c,
		AuthMiddleware: middleware.NewHttpAuthMiddleware(c.JwtAuth.AccessSecret, constants.UserTypeLender).Handle,

		// rpc client
		UserRpcClient:         lender.NewUserClient(zrpc.MustNewClient(c.RpcConf).Conn()),
		LoanRpcClient:         lender.NewLoanClient(zrpc.MustNewClient(c.RpcConf).Conn()),
		CommonRpcClient:       lender.NewCommonClient(zrpc.MustNewClient(c.RpcConf).Conn()),
		StorageRpcClient:      lender.NewStorageClient(zrpc.MustNewClient(c.RpcConf).Conn()),
		FileRpcClient:         lender.NewFileClient(zrpc.MustNewClient(c.RpcConf).Conn()),
		LoanDocumentRpcClient: lender.NewLoanDocumentClient(zrpc.MustNewClient(c.RpcConf).Conn()),
	}
}
