package rpc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/zeromicro/go-zero/core/service"

	borrowerApplicantServer "zeebee-backend/app/borrower/rpc/internal/server/borrowerapplicant"
	borrowerEntityServer "zeebee-backend/app/borrower/rpc/internal/server/borrowerentity"
	borrowerPropertyServer "zeebee-backend/app/borrower/rpc/internal/server/borrowerproperty"

	loanApplicantServer "zeebee-backend/app/borrower/rpc/internal/server/loanapplicant"

	"zeebee-backend/app/borrower/config"
	"zeebee-backend/app/borrower/rpc/internal/global"
	commonServer "zeebee-backend/app/borrower/rpc/internal/server/common"
	loanServer "zeebee-backend/app/borrower/rpc/internal/server/loan"
	loanDocumentServer "zeebee-backend/app/borrower/rpc/internal/server/loandocument"
	roleServer "zeebee-backend/app/borrower/rpc/internal/server/role"
	storageServer "zeebee-backend/app/borrower/rpc/internal/server/storage"
	userServer "zeebee-backend/app/borrower/rpc/internal/server/user"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
)

type RpcServer struct {
	s *zrpc.RpcServer
}

func NewRpcServer(c config.Config) (s service.Service, err error) {
	rpcConfig := c.RpcConf

	err = global.InitGlobal(rpcConfig)
	if err != nil {
		return nil, err
	}

	ctx := svc.NewServiceContext(rpcConfig)
	rpcServer := zrpc.MustNewServer(rpcConfig.RpcServerConf, func(grpcServer *grpc.Server) {
		borrower.RegisterUserServer(grpcServer, userServer.NewUserServer(ctx))
		borrower.RegisterLoanServer(grpcServer, loanServer.NewLoanServer(ctx))
		borrower.RegisterCommonServer(grpcServer, commonServer.NewCommonServer(ctx))
		borrower.RegisterBorrowerPropertyServer(grpcServer, borrowerPropertyServer.NewBorrowerPropertyServer(ctx))
		borrower.RegisterLoanApplicantServer(grpcServer, loanApplicantServer.NewLoanApplicantServer(ctx))
		borrower.RegisterStorageServer(grpcServer, storageServer.NewStorageServer(ctx))
		borrower.RegisterLoanDocumentServer(grpcServer, loanDocumentServer.NewLoanDocumentServer(ctx))
		borrower.RegisterBorrowerEntityServer(grpcServer, borrowerEntityServer.NewBorrowerEntityServer(ctx))
		borrower.RegisterBorrowerApplicantServer(grpcServer, borrowerApplicantServer.NewBorrowerApplicantServer(ctx))
		borrower.RegisterRoleServer(grpcServer, roleServer.NewRoleServer(ctx))

		if rpcConfig.Mode == service.DevMode || rpcConfig.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	return &RpcServer{
		s: rpcServer,
	}, nil
}

func (s *RpcServer) Start() {
	go func() {
		s.s.Start()
	}()
}

func (s *RpcServer) Stop() {
	s.s.Stop()
}
