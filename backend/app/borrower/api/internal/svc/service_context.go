package svc

import (
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"

	"zeebee-backend/app/borrower/config"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
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
	UserRpcClient                           borrower.UserClient
	LoanRpcClient                           borrower.LoanClient
	CommonRpcClient                         borrower.CommonClient
	BorrowerPropertyRpcClient               borrower.BorrowerPropertyClient
	BorrowerEntityRpcClient                 borrower.BorrowerEntityClient
	BorrowerApplicantRpcClient              borrower.BorrowerApplicantClient
	LoanApplicantRpcClient                  borrower.LoanApplicantClient
	StorageRpcClient                        borrower.StorageClient
	FileRpcClient                           borrower.FileClient
	LoanDocumentRpcClient                   borrower.LoanDocumentClient
	RoleRpcClient                           borrower.RoleClient
	OrganizationUserRpcClient               borrower.OrganizationUserClient
	OrganizationRpcClient                   borrower.OrganizationClient
	OrganizationEmailRpcClient              borrower.OrganizationEmailClient
	ZeitroRpcClient                         borrower.ZeitroClient
	ResourceRpcClient                       borrower.ResourceClient
	BorrowerPropertyTransactionRpcClient    borrower.BorrowerPropertyTransactionClient
	DocumentRpcClient                       borrower.DocumentClient
	DriveRpcClient                          borrower.DriveClient
	ApiRpcClient                            borrower.ApiClient
	OrganizationSubscriptionRpcClient       borrower.OrganizationSubscriptionClient
	ProjectManagementWorkflowRpcClient      borrower.ProjectManagementWorkflowClient
	ProjectManagementWorkflowStageRpcClient borrower.ProjectManagementWorkflowStageClient
	ProjectManagementTaskTypeRpcClient      borrower.ProjectManagementTaskTypeClient
	ProjectManagementProjectRpcClient       borrower.ProjectManagementProjectClient
	ProjectManagementProjectMemberRpcClient borrower.ProjectManagementProjectMemberClient
	ProjectManagementTaskRpcClient          borrower.ProjectManagementTaskClient
	OperationRpcClient                      borrower.OperationClient
}

func NewServiceContext(c config.ApiConf, logger logx.Logger) *ServiceContext {
	rpcClient := zrpc.MustNewClient(c.RpcConf)

	return &ServiceContext{
		Config:         c,
		Logger:         logger,
		AuthMiddleware: middleware.NewAuthMiddleware(constants.UserTypeBorrower, tenantclient.New(zrpc.MustNewClient(c.TenantRpc)).Auth).Handle,
		// rpc client
		UserRpcClient:                           borrower.NewUserClient(rpcClient.Conn()),
		LoanRpcClient:                           borrower.NewLoanClient(rpcClient.Conn()),
		CommonRpcClient:                         borrower.NewCommonClient(rpcClient.Conn()),
		BorrowerPropertyRpcClient:               borrower.NewBorrowerPropertyClient(rpcClient.Conn()),
		BorrowerEntityRpcClient:                 borrower.NewBorrowerEntityClient(rpcClient.Conn()),
		BorrowerApplicantRpcClient:              borrower.NewBorrowerApplicantClient(rpcClient.Conn()),
		LoanApplicantRpcClient:                  borrower.NewLoanApplicantClient(rpcClient.Conn()),
		StorageRpcClient:                        borrower.NewStorageClient(rpcClient.Conn()),
		FileRpcClient:                           borrower.NewFileClient(rpcClient.Conn()),
		LoanDocumentRpcClient:                   borrower.NewLoanDocumentClient(rpcClient.Conn()),
		RoleRpcClient:                           borrower.NewRoleClient(rpcClient.Conn()),
		OrganizationUserRpcClient:               borrower.NewOrganizationUserClient(rpcClient.Conn()),
		OrganizationRpcClient:                   borrower.NewOrganizationClient(rpcClient.Conn()),
		OrganizationEmailRpcClient:              borrower.NewOrganizationEmailClient(rpcClient.Conn()),
		ZeitroRpcClient:                         borrower.NewZeitroClient(rpcClient.Conn()),
		ResourceRpcClient:                       borrower.NewResourceClient(rpcClient.Conn()),
		BorrowerPropertyTransactionRpcClient:    borrower.NewBorrowerPropertyTransactionClient(rpcClient.Conn()),
		DocumentRpcClient:                       borrower.NewDocumentClient(rpcClient.Conn()),
		DriveRpcClient:                          borrower.NewDriveClient(rpcClient.Conn()),
		ApiRpcClient:                            borrower.NewApiClient(rpcClient.Conn()),
		OrganizationSubscriptionRpcClient:       borrower.NewOrganizationSubscriptionClient(rpcClient.Conn()),
		ProjectManagementWorkflowRpcClient:      borrower.NewProjectManagementWorkflowClient(rpcClient.Conn()),
		ProjectManagementWorkflowStageRpcClient: borrower.NewProjectManagementWorkflowStageClient(rpcClient.Conn()),
		ProjectManagementTaskTypeRpcClient:      borrower.NewProjectManagementTaskTypeClient(rpcClient.Conn()),
		ProjectManagementProjectRpcClient:       borrower.NewProjectManagementProjectClient(rpcClient.Conn()),
		ProjectManagementProjectMemberRpcClient: borrower.NewProjectManagementProjectMemberClient(rpcClient.Conn()),
		ProjectManagementTaskRpcClient:          borrower.NewProjectManagementTaskClient(rpcClient.Conn()),
		OperationRpcClient:                      borrower.NewOperationClient(rpcClient.Conn()),
	}
}
