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
	OrganizationRpcClient     admin.OrganizationClient
	StorageRpcClient          admin.StorageClient
	CommonRpcClient           admin.CommonClient
	BaseTemplateRpcClient     admin.BaseTemplateClient
	LoanRpcClient             admin.LoanClient
}

func NewServiceContext(c config.ApiConf) *ServiceContext {
	authRpcClient := admin.NewAuthClient(zrpc.MustNewClient(c.RpcConf).Conn())
	return &ServiceContext{
		Config:         c,
		AuthMiddleware: middleware.NewHttpAuthMiddleware(constants.UserTypeAdmin, authRpcClient).Handle,

		// rpc client
		UserRpcClient:             admin.NewUserClient(zrpc.MustNewClient(c.RpcConf).Conn()),
		ApiRpcClient:              admin.NewApiClient(zrpc.MustNewClient(c.RpcConf).Conn()),
		ResourceRpcClient:         admin.NewResourceClient(zrpc.MustNewClient(c.RpcConf).Conn()),
		RoleRpcClient:             admin.NewRoleClient(zrpc.MustNewClient(c.RpcConf).Conn()),
		OrganizationUserRpcClient: admin.NewOrganizationUserClient(zrpc.MustNewClient(c.RpcConf).Conn()),
		OrganizationRpcClient:     admin.NewOrganizationClient(zrpc.MustNewClient(c.RpcConf).Conn()),
		StorageRpcClient:          admin.NewStorageClient(zrpc.MustNewClient(c.RpcConf).Conn()),
		CommonRpcClient:           admin.NewCommonClient(zrpc.MustNewClient(c.RpcConf).Conn()),
		BaseTemplateRpcClient:     admin.NewBaseTemplateClient(zrpc.MustNewClient(c.RpcConf).Conn()),
		LoanRpcClient:             admin.NewLoanClient(zrpc.MustNewClient(c.RpcConf).Conn()),
	}
}
