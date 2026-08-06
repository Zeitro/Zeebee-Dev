package svc

import (
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"

	"zeebee-backend/app/admin/config"
	"zeebee-backend/app/admin/rpc/pb/admin"
	tenantauth "zeebee-backend/app/tenant/rpc/client/auth"
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
	UserRpcClient                     admin.UserClient
	ApiRpcClient                      admin.ApiClient
	ResourceRpcClient                 admin.ResourceClient
	RoleRpcClient                     admin.RoleClient
	OrganizationUserRpcClient         admin.OrganizationUserClient
	OrganizationRpcClient             admin.OrganizationClient
	StorageRpcClient                  admin.StorageClient
	CommonRpcClient                   admin.CommonClient
	BaseTemplateRpcClient             admin.BaseTemplateClient
	LoanRpcClient                     admin.LoanClient
	PricingPlanRpcClient              admin.PricingPlanClient
	PricingPlanPolicyRpcClient        admin.PricingPlanPolicyClient
	OrganizationSubscriptionRpcClient admin.OrganizationSubscriptionClient

	AdminApiContent string
}

func NewServiceContext(c config.ApiConf, logger logx.Logger, adminApiContent string) *ServiceContext {
	rpcClient := zrpc.MustNewClient(c.RpcConf)

	return &ServiceContext{
		Config:         c,
		AuthMiddleware: middleware.NewAuthMiddleware(constants.UserTypeAdmin, tenantauth.NewAuth(zrpc.MustNewClient(c.TenantRpc)), c.JwtConf.AccessSecret).Handle,
		Logger:         logger,

		// rpc client
		UserRpcClient:                     admin.NewUserClient(rpcClient.Conn()),
		ApiRpcClient:                      admin.NewApiClient(rpcClient.Conn()),
		ResourceRpcClient:                 admin.NewResourceClient(rpcClient.Conn()),
		RoleRpcClient:                     admin.NewRoleClient(rpcClient.Conn()),
		OrganizationUserRpcClient:         admin.NewOrganizationUserClient(rpcClient.Conn()),
		OrganizationRpcClient:             admin.NewOrganizationClient(rpcClient.Conn()),
		StorageRpcClient:                  admin.NewStorageClient(rpcClient.Conn()),
		CommonRpcClient:                   admin.NewCommonClient(rpcClient.Conn()),
		BaseTemplateRpcClient:             admin.NewBaseTemplateClient(rpcClient.Conn()),
		LoanRpcClient:                     admin.NewLoanClient(rpcClient.Conn()),
		PricingPlanRpcClient:              admin.NewPricingPlanClient(rpcClient.Conn()),
		PricingPlanPolicyRpcClient:        admin.NewPricingPlanPolicyClient(rpcClient.Conn()),
		OrganizationSubscriptionRpcClient: admin.NewOrganizationSubscriptionClient(rpcClient.Conn()),
		AdminApiContent:                   adminApiContent,
	}
}
