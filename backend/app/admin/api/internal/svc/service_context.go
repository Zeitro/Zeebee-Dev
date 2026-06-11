package svc

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"

	"zeebee-backend/app/admin/api/internal/middleware"
	"zeebee-backend/app/admin/config"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/domains/zeebee/constants"
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
	OrganizationRpcClient     admin.OrganizationClient
	StorageRpcClient          admin.StorageClient
	CommonRpcClient           admin.CommonClient
	BaseTemplateRpcClient     admin.BaseTemplateClient
	LoanRpcClient             admin.LoanClient

	AdminApiContent string
}

func NewServiceContext(c config.ApiConf, adminApiContent string) *ServiceContext {
	rpcClient := zrpc.MustNewClient(c.RpcConf)

	authRpcClient := admin.NewAuthClient(rpcClient.Conn())
	return &ServiceContext{
		Config:         c,
		AuthMiddleware: middleware.NewAuthMiddleware(constants.UserTypeAdmin, authRpcClient).Handle,

		// rpc client
		UserRpcClient:             admin.NewUserClient(rpcClient.Conn()),
		ApiRpcClient:              admin.NewApiClient(rpcClient.Conn()),
		ResourceRpcClient:         admin.NewResourceClient(rpcClient.Conn()),
		RoleRpcClient:             admin.NewRoleClient(rpcClient.Conn()),
		OrganizationUserRpcClient: admin.NewOrganizationUserClient(rpcClient.Conn()),
		OrganizationRpcClient:     admin.NewOrganizationClient(rpcClient.Conn()),
		StorageRpcClient:          admin.NewStorageClient(rpcClient.Conn()),
		CommonRpcClient:           admin.NewCommonClient(rpcClient.Conn()),
		BaseTemplateRpcClient:     admin.NewBaseTemplateClient(rpcClient.Conn()),
		LoanRpcClient:             admin.NewLoanClient(rpcClient.Conn()),
		AdminApiContent:           adminApiContent,
	}
}
