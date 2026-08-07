package svc

import (
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"

	"zeebee-backend/app/lender/config"
	"zeebee-backend/app/lender/rpc/pb/lender"
	tenantclient "zeebee-backend/app/tenant/rpc/client"
	"zeebee-backend/app/tenant/rpc/middleware"
	"zeebee-backend/domains/zeebee/constants"
)

type ServiceContext struct {
	Config config.ApiConf
	// Logger
	Logger logx.Logger
	// middleware
	AuthMiddleware rest.Middleware

	// rpc client
	UserRpcClient                 lender.UserClient
	LoanRpcClient                 lender.LoanClient
	CommonRpcClient               lender.CommonClient
	StorageRpcClient              lender.StorageClient
	FileRpcClient                 lender.FileClient
	LoanDocumentRpcClient         lender.LoanDocumentClient
	RoleRpcClient                 lender.RoleClient
	OrganizationUserRpcClient     lender.OrganizationUserClient
	OrganizationRpcClient         lender.OrganizationClient
	OrganizationTemplateRpcClient lender.OrganizationTemplateClient
	OrganizationSettingRpcClient  lender.OrganizationSettingClient
	ResourceRpcClient             lender.ResourceClient
	ApiRpcClient                  lender.ApiClient

	BorrowerApiContent string
}

func NewServiceContext(c config.ApiConf, logger logx.Logger, borrowerApiContent string) *ServiceContext {
	rpcClient := zrpc.MustNewClient(c.RpcConf)
	return &ServiceContext{
		Config:         c,
		Logger:         logger,
		AuthMiddleware: middleware.NewAuthMiddleware(constants.UserTypeLender, tenantclient.New(zrpc.MustNewClient(c.TenantRpc)).Auth, c.JwtConf.AccessSecret).Handle,

		// rpc client
		UserRpcClient:                 lender.NewUserClient(rpcClient.Conn()),
		LoanRpcClient:                 lender.NewLoanClient(rpcClient.Conn()),
		CommonRpcClient:               lender.NewCommonClient(rpcClient.Conn()),
		StorageRpcClient:              lender.NewStorageClient(rpcClient.Conn()),
		FileRpcClient:                 lender.NewFileClient(rpcClient.Conn()),
		LoanDocumentRpcClient:         lender.NewLoanDocumentClient(rpcClient.Conn()),
		RoleRpcClient:                 lender.NewRoleClient(rpcClient.Conn()),
		OrganizationUserRpcClient:     lender.NewOrganizationUserClient(rpcClient.Conn()),
		OrganizationRpcClient:         lender.NewOrganizationClient(rpcClient.Conn()),
		OrganizationTemplateRpcClient: lender.NewOrganizationTemplateClient(rpcClient.Conn()),
		OrganizationSettingRpcClient:  lender.NewOrganizationSettingClient(rpcClient.Conn()),
		ResourceRpcClient:             lender.NewResourceClient(rpcClient.Conn()),
		ApiRpcClient:                  lender.NewApiClient(rpcClient.Conn()),
		BorrowerApiContent:            borrowerApiContent,
	}
}
