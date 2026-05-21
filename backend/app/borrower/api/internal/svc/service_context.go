package svc

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"

	"zeebee-backend/app/borrower/api/internal/middleware"
	"zeebee-backend/app/borrower/config"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/domains/zeebee/constants"
)

type ServiceContext struct {
	Config config.ApiConf
	// middleware
	AuthMiddleware rest.Middleware

	// rpc client
	UserRpcClient              borrower.UserClient
	LoanRpcClient              borrower.LoanClient
	CommonRpcClient            borrower.CommonClient
	BorrowerPropertyRpcClient  borrower.BorrowerPropertyClient
	BorrowerEntityRpcClient    borrower.BorrowerEntityClient
	BorrowerApplicantRpcClient borrower.BorrowerApplicantClient
	LoanApplicantRpcClient     borrower.LoanApplicantClient
	StorageRpcClient           borrower.StorageClient
	FileRpcClient              borrower.FileClient
	LoanDocumentRpcClient      borrower.LoanDocumentClient
	RoleRpcClient              borrower.RoleClient
	OrganizationUserRpcClient  borrower.OrganizationUserClient
	OrganizationRpcClient      borrower.OrganizationClient
	OrganizationEmailRpcClient borrower.OrganizationEmailClient
	ZeitroRpcClient                       borrower.ZeitroClient
	ResourceRpcClient                     borrower.ResourceClient
	BorrowerPropertyTransactionRpcClient  borrower.BorrowerPropertyTransactionClient
	DocumentRpcClient                     borrower.DocumentClient
}

func NewServiceContext(c config.ApiConf) *ServiceContext {
	rpcClient := zrpc.MustNewClient(c.RpcConf)
	authRpcClient := borrower.NewAuthClient(rpcClient.Conn())

	return &ServiceContext{
		Config:         c,
		AuthMiddleware: middleware.NewAuthMiddleware(constants.UserTypeBorrower, authRpcClient).Handle,
		// rpc client
		UserRpcClient:              borrower.NewUserClient(rpcClient.Conn()),
		LoanRpcClient:              borrower.NewLoanClient(rpcClient.Conn()),
		CommonRpcClient:            borrower.NewCommonClient(rpcClient.Conn()),
		BorrowerPropertyRpcClient:  borrower.NewBorrowerPropertyClient(rpcClient.Conn()),
		BorrowerEntityRpcClient:    borrower.NewBorrowerEntityClient(rpcClient.Conn()),
		BorrowerApplicantRpcClient: borrower.NewBorrowerApplicantClient(rpcClient.Conn()),
		LoanApplicantRpcClient:     borrower.NewLoanApplicantClient(rpcClient.Conn()),
		StorageRpcClient:           borrower.NewStorageClient(rpcClient.Conn()),
		FileRpcClient:              borrower.NewFileClient(rpcClient.Conn()),
		LoanDocumentRpcClient:      borrower.NewLoanDocumentClient(rpcClient.Conn()),
		RoleRpcClient:              borrower.NewRoleClient(rpcClient.Conn()),
		OrganizationUserRpcClient:  borrower.NewOrganizationUserClient(rpcClient.Conn()),
		OrganizationRpcClient:      borrower.NewOrganizationClient(rpcClient.Conn()),
		OrganizationEmailRpcClient: borrower.NewOrganizationEmailClient(rpcClient.Conn()),
		ZeitroRpcClient:                       borrower.NewZeitroClient(rpcClient.Conn()),
		ResourceRpcClient:                     borrower.NewResourceClient(rpcClient.Conn()),
		BorrowerPropertyTransactionRpcClient:  borrower.NewBorrowerPropertyTransactionClient(rpcClient.Conn()),
		DocumentRpcClient:                     borrower.NewDocumentClient(rpcClient.Conn()),
	}
}
