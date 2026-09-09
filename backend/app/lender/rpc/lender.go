package rpc

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"

	"github.com/zeromicro/go-zero/core/service"

	"zeebee-backend/app/lender/config"
	"zeebee-backend/app/lender/rpc/internal/global"
	apiServer "zeebee-backend/app/lender/rpc/internal/server/api"
	commonServer "zeebee-backend/app/lender/rpc/internal/server/common"
	fileServer "zeebee-backend/app/lender/rpc/internal/server/file"
	loanServer "zeebee-backend/app/lender/rpc/internal/server/loan"
	loanDocumentServer "zeebee-backend/app/lender/rpc/internal/server/loandocument"
	organizationServer "zeebee-backend/app/lender/rpc/internal/server/organization"
	organizationSettingServer "zeebee-backend/app/lender/rpc/internal/server/organizationsetting"
	organizationTemplateServer "zeebee-backend/app/lender/rpc/internal/server/organizationtemplate"
	organizationUserServer "zeebee-backend/app/lender/rpc/internal/server/organizationuser"
	resourceServer "zeebee-backend/app/lender/rpc/internal/server/resource"
	roleServer "zeebee-backend/app/lender/rpc/internal/server/role"
	storageServer "zeebee-backend/app/lender/rpc/internal/server/storage"
	userServer "zeebee-backend/app/lender/rpc/internal/server/user"
	"zeebee-backend/app/lender/rpc/internal/svc"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/errorx"
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
	logger := logx.WithCallerSkip(0).
		WithFields(logx.Field("service_name", c.RpcConf.Name))
	ctx := svc.NewServiceContext(rpcConfig, logger)
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
		lender.RegisterOrganizationTemplateServer(grpcServer, organizationTemplateServer.NewOrganizationTemplateServer(ctx))
		lender.RegisterOrganizationSettingServer(grpcServer, organizationSettingServer.NewOrganizationSettingServer(ctx))
		lender.RegisterResourceServer(grpcServer, resourceServer.NewResourceServer(ctx))
		lender.RegisterApiServer(grpcServer, apiServer.NewApiServer(ctx))

		if rpcConfig.Mode == service.DevMode || rpcConfig.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	// Outermost interceptor: tag every rpc ctx with service_name, and propagate
	// incoming gRPC metadata (jwt_token / data_scope) to the outgoing context so
	// the US-12 proxy calls to tenant rpc carry the caller's identity across the
	// hop. grpc.UnaryInterceptor is prepended by grpc-go to the chain (outermost),
	// so it runs above go-zero's stat interceptor and the access-log line carries
	// the field. gRPC does not mirror incoming→outgoing metadata automatically;
	// go-zero only forwards trace headers.
	rpcServer.AddOptions(grpc.UnaryInterceptor(func(ctx context.Context, req any,
		info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		ctx = logx.ContextWithFields(ctx, logx.Field("service_name", c.RpcConf.Name))
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			ctx = metadata.NewOutgoingContext(ctx, md.Copy())
		}
		return handler(ctx, req)
	}))
	// Centralized error logging. Registered via AddUnaryInterceptors so it sits at
	// the tail of go-zero's chain — INSIDE the tracing interceptor — and the ctx
	// here already carries trace+span (from tracing) plus service_name (from the
	// outer interceptor above). Logging in the outermost interceptor instead loses
	// trace/span because tracing has not run yet at that layer.
	rpcServer.AddUnaryInterceptors(func(ctx context.Context, req any,
		info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		resp, err := handler(ctx, req)
		if err != nil {
			logger.WithContext(ctx).WithFields(
				logx.Field("stacks", errorx.GetErrorStack(err)),
			).Error(err)
		}
		return resp, err
	})
	rpcService := &RpcServer{
		s: rpcServer,
	}
	return rpcService, nil
}

func (s *RpcServer) Start() {
	go func() {
		s.s.Start()
	}()
}

func (s *RpcServer) Stop() {
	s.s.Stop()
}
