package svc

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"

	"zeebee-backend/app/borrower/config"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/middleware"
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
	ZeitroRpcClient            borrower.ZeitroClient
}

func NewServiceContext(c config.ApiConf) *ServiceContext {
	return &ServiceContext{
		Config:         c,
		AuthMiddleware: middleware.NewHttpAuthMiddleware(c.JwtAuth.AccessSecret, constants.UserTypeBorrower).Handle,
		// rpc client
		UserRpcClient:              borrower.NewUserClient(zrpc.MustNewClient(c.RpcConf).Conn()),
		LoanRpcClient:              borrower.NewLoanClient(zrpc.MustNewClient(c.RpcConf).Conn()),
		CommonRpcClient:            borrower.NewCommonClient(zrpc.MustNewClient(c.RpcConf).Conn()),
		BorrowerPropertyRpcClient:  borrower.NewBorrowerPropertyClient(zrpc.MustNewClient(c.RpcConf).Conn()),
		BorrowerEntityRpcClient:    borrower.NewBorrowerEntityClient(zrpc.MustNewClient(c.RpcConf).Conn()),
		BorrowerApplicantRpcClient: borrower.NewBorrowerApplicantClient(zrpc.MustNewClient(c.RpcConf).Conn()),
		LoanApplicantRpcClient:     borrower.NewLoanApplicantClient(zrpc.MustNewClient(c.RpcConf).Conn()),
		StorageRpcClient:           borrower.NewStorageClient(zrpc.MustNewClient(c.RpcConf).Conn()),
		FileRpcClient:              borrower.NewFileClient(zrpc.MustNewClient(c.RpcConf).Conn()),
		LoanDocumentRpcClient:      borrower.NewLoanDocumentClient(zrpc.MustNewClient(c.RpcConf).Conn()),
		RoleRpcClient:              borrower.NewRoleClient(zrpc.MustNewClient(c.RpcConf).Conn()),
		OrganizationUserRpcClient:  borrower.NewOrganizationUserClient(zrpc.MustNewClient(c.RpcConf).Conn()),
		OrganizationRpcClient:      borrower.NewOrganizationClient(zrpc.MustNewClient(c.RpcConf).Conn()),
		OrganizationEmailRpcClient: borrower.NewOrganizationEmailClient(zrpc.MustNewClient(c.RpcConf).Conn()),
		ZeitroRpcClient:            borrower.NewZeitroClient(zrpc.MustNewClient(c.RpcConf).Conn()),
	}
}
