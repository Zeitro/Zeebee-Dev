package rpc

import (
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	"zeebee-backend/app/tenant/config"
	authzServer "zeebee-backend/app/tenant/rpc/internal/server/authz"
	userServer "zeebee-backend/app/tenant/rpc/internal/server/user"
	"zeebee-backend/app/tenant/rpc/internal/svc"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
)

// RpcServer wraps the tenant gRPC server.
type RpcServer struct {
	s      *zrpc.RpcServer
	svcCtx *svc.ServiceContext
}

// NewRpcServer constructs the tenant RPC server.
// Phase 2 (T007): global.InitGlobal(rpcConfig) for gorm/redis/casbin enforcer.
func NewRpcServer(c config.Config) (service.Service, error) {
	rpcConfig := c.RpcConf

	logger := logx.WithCallerSkip(0).
		WithFields(logx.Field("service_name", rpcConfig.Name))
	ctx, err := svc.NewServiceContext(rpcConfig, logger)
	if err != nil {
		return nil, err
	}

	rpcServer := zrpc.MustNewServer(rpcConfig.RpcServerConf, func(grpcServer *grpc.Server) {
		tenant.RegisterAuthzServer(grpcServer, authzServer.NewAuthzServer(ctx))
		tenant.RegisterUserServer(grpcServer, userServer.NewUserServer(ctx))

		// gRPC health probe (grpc_health_v1) for k8s/deploy liveness (T008).
		healthpb.RegisterHealthServer(grpcServer, health.NewServer())

		if rpcConfig.Mode == service.DevMode || rpcConfig.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})

	return &RpcServer{
		s:      rpcServer,
		svcCtx: ctx,
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
