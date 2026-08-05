package rpc

import (
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"zeebee-backend/app/tenant/config"
	apiServer "zeebee-backend/app/tenant/rpc/internal/server/api"
	authServer "zeebee-backend/app/tenant/rpc/internal/server/auth"
	organizationServer "zeebee-backend/app/tenant/rpc/internal/server/organization"
	organizationEmailServer "zeebee-backend/app/tenant/rpc/internal/server/organizationemail"
	organizationSettingServer "zeebee-backend/app/tenant/rpc/internal/server/organizationsetting"
	organizationSubscriptionServer "zeebee-backend/app/tenant/rpc/internal/server/organizationsubscription"
	organizationTemplateServer "zeebee-backend/app/tenant/rpc/internal/server/organizationtemplate"
	organizationUserServer "zeebee-backend/app/tenant/rpc/internal/server/organizationuser"
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
		// Identity issuance + per-request authorization (direction B).
		tenant.RegisterAuthServer(grpcServer, authServer.NewAuthServer(ctx))
		tenant.RegisterUserServer(grpcServer, userServer.NewUserServer(ctx))

		// RBAC object registry (api endpoints / resources / roles / policies).
		tenant.RegisterApiServer(grpcServer, apiServer.NewApiServer(ctx))

		// Organization modules (consolidated from admin / lender / borrower).
		tenant.RegisterOrganizationServer(grpcServer, organizationServer.NewOrganizationServer(ctx))
		tenant.RegisterOrganizationUserServer(grpcServer, organizationUserServer.NewOrganizationUserServer(ctx))
		tenant.RegisterOrganizationEmailServer(grpcServer, organizationEmailServer.NewOrganizationEmailServer(ctx))
		tenant.RegisterOrganizationSettingServer(grpcServer, organizationSettingServer.NewOrganizationSettingServer(ctx))
		tenant.RegisterOrganizationSubscriptionServer(grpcServer, organizationSubscriptionServer.NewOrganizationSubscriptionServer(ctx))
		tenant.RegisterOrganizationTemplateServer(grpcServer, organizationTemplateServer.NewOrganizationTemplateServer(ctx))

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
