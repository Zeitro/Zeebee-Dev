package rpc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/zeromicro/go-zero/core/service"

	apiServer "zeebee-backend/app/admin/rpc/internal/server/api"
	authServer "zeebee-backend/app/admin/rpc/internal/server/auth"
	baseTemplateServer "zeebee-backend/app/admin/rpc/internal/server/basetemplate"
	commonServer "zeebee-backend/app/admin/rpc/internal/server/common"
	loanServer "zeebee-backend/app/admin/rpc/internal/server/loan"
	organizationServer "zeebee-backend/app/admin/rpc/internal/server/organization"
	organizationSubscriptionServer "zeebee-backend/app/admin/rpc/internal/server/organizationsubscription"
	organizationUserServer "zeebee-backend/app/admin/rpc/internal/server/organizationuser"
	pricingPlanServer "zeebee-backend/app/admin/rpc/internal/server/pricingplan"
	pricingPlanPolicyServer "zeebee-backend/app/admin/rpc/internal/server/pricingplanpolicy"
	resourceServer "zeebee-backend/app/admin/rpc/internal/server/resource"
	roleServer "zeebee-backend/app/admin/rpc/internal/server/role"
	storageServer "zeebee-backend/app/admin/rpc/internal/server/storage"
	userServer "zeebee-backend/app/admin/rpc/internal/server/user"
	"zeebee-backend/app/admin/rpc/pb/admin"

	"zeebee-backend/app/admin/config"
	"zeebee-backend/app/admin/rpc/internal/global"
	"zeebee-backend/app/admin/rpc/internal/svc"
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
		admin.RegisterUserServer(grpcServer, userServer.NewUserServer(ctx))
		admin.RegisterApiServer(grpcServer, apiServer.NewApiServer(ctx))
		admin.RegisterResourceServer(grpcServer, resourceServer.NewResourceServer(ctx))
		admin.RegisterRoleServer(grpcServer, roleServer.NewRoleServer(ctx))
		admin.RegisterOrganizationUserServer(grpcServer, organizationUserServer.NewOrganizationUserServer(ctx))
		admin.RegisterOrganizationServer(grpcServer, organizationServer.NewOrganizationServer(ctx))
		admin.RegisterStorageServer(grpcServer, storageServer.NewStorageServer(ctx))
		admin.RegisterCommonServer(grpcServer, commonServer.NewCommonServer(ctx))
		admin.RegisterBaseTemplateServer(grpcServer, baseTemplateServer.NewBaseTemplateServer(ctx))
		admin.RegisterLoanServer(grpcServer, loanServer.NewLoanServer(ctx))
		admin.RegisterAuthServer(grpcServer, authServer.NewAuthServer(ctx))
		admin.RegisterPricingPlanServer(grpcServer, pricingPlanServer.NewPricingPlanServer(ctx))
		admin.RegisterPricingPlanPolicyServer(grpcServer, pricingPlanPolicyServer.NewPricingPlanPolicyServer(ctx))
		admin.RegisterOrganizationSubscriptionServer(grpcServer, organizationSubscriptionServer.NewOrganizationSubscriptionServer(ctx))
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
