package rpc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/zeromicro/go-zero/core/service"

	"zeebee-backend/app/lender/config"
	"zeebee-backend/app/lender/rpc/internal/global"
	commonServer "zeebee-backend/app/lender/rpc/internal/server/common"
	fileServer "zeebee-backend/app/lender/rpc/internal/server/file"
	loanServer "zeebee-backend/app/lender/rpc/internal/server/loan"
	loanDocumentServer "zeebee-backend/app/lender/rpc/internal/server/loandocument"
	organizationServer "zeebee-backend/app/lender/rpc/internal/server/organization"
	organizationUserServer "zeebee-backend/app/lender/rpc/internal/server/organizationuser"
	roleServer "zeebee-backend/app/lender/rpc/internal/server/role"
	storageServer "zeebee-backend/app/lender/rpc/internal/server/storage"
	userServer "zeebee-backend/app/lender/rpc/internal/server/user"
	"zeebee-backend/app/lender/rpc/internal/svc"
	"zeebee-backend/app/lender/rpc/pb/lender"
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
		lender.RegisterUserServer(grpcServer, userServer.NewUserServer(ctx))
		lender.RegisterCommonServer(grpcServer, commonServer.NewCommonServer(ctx))
		lender.RegisterStorageServer(grpcServer, storageServer.NewStorageServer(ctx))
		lender.RegisterLoanServer(grpcServer, loanServer.NewLoanServer(ctx))
		lender.RegisterFileServer(grpcServer, fileServer.NewFileServer(ctx))
		lender.RegisterLoanDocumentServer(grpcServer, loanDocumentServer.NewLoanDocumentServer(ctx))
		lender.RegisterRoleServer(grpcServer, roleServer.NewRoleServer(ctx))
		lender.RegisterOrganizationUserServer(grpcServer, organizationUserServer.NewOrganizationUserServer(ctx))
		lender.RegisterOrganizationServer(grpcServer, organizationServer.NewOrganizationServer(ctx))

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
