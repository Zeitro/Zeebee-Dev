package rpc

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"

	"github.com/zeromicro/go-zero/core/service"

	"zeebee-backend/app/admin/config"
	"zeebee-backend/app/admin/rpc/internal/global"
	apiServer "zeebee-backend/app/admin/rpc/internal/server/api"
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
	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/cron"
	"zeebee-backend/common/errorx"
	domainlenderdiscover "zeebee-backend/domains/zeebee/lenderdiscover"

	"zeebee-backend/app/admin/rpc/internal/schedule"
)

type RpcServer struct {
	s      *zrpc.RpcServer
	svcCtx *svc.ServiceContext
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
		admin.RegisterPricingPlanServer(grpcServer, pricingPlanServer.NewPricingPlanServer(ctx))
		admin.RegisterPricingPlanPolicyServer(grpcServer, pricingPlanPolicyServer.NewPricingPlanPolicyServer(ctx))
		admin.RegisterOrganizationSubscriptionServer(grpcServer, organizationSubscriptionServer.NewOrganizationSubscriptionServer(ctx))
		if rpcConfig.Mode == service.DevMode || rpcConfig.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	// Outermost interceptor: tag every rpc ctx with service_name, and propagate
	// incoming gRPC metadata (jwt_token / user_type / data_scope) to the outgoing
	// context so the US-12 proxy calls to tenant rpc carry the caller's identity
	// across the hop. grpc.UnaryInterceptor is prepended by grpc-go to the chain
	// (outermost), so it runs above go-zero's stat interceptor and the access-log
	// line carries the field. gRPC does not mirror incoming→outgoing metadata
	// automatically; go-zero only forwards trace headers.
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
		s:      rpcServer,
		svcCtx: ctx,
	}
	err = RegisterCronTask(ctx, rpcConfig, logger)
	if err != nil {
		return nil, err
	}

	return rpcService, nil
}

func (s *RpcServer) Start() {
	go func() {
		s.s.Start()
	}()
	cron.Start()
}

func (s *RpcServer) Stop() {
	s.s.Stop()
	cron.Stop()
}

func RegisterCronTask(svcCtx *svc.ServiceContext, conf config.RpcConf, logger logx.Logger) (err error) {
	// lender discover: spec derived from DB config. DiscoverUpdateSettings calls
	// cron.Replace/Unregister to keep the timer in sync on config change.
	if cfg, _, cfgErr := domainlenderdiscover.LoadConfig(context.Background(), svcCtx.OrganizationSettingRepo); cfgErr != nil {
		logger.Errorf("lender discover cron register: load settings: %v", cfgErr)
	} else if cfg.ScheduleEnabled {
		if regErr := cron.Register(schedule.TaskName, schedule.SpecFromInterval(cfg.IntervalDays), func() {
			schedule.RunScheduled(context.Background(), svcCtx)
		}); regErr != nil {
			logger.Errorf("lender discover cron register: %v", regErr)
		}
	}
	return
}
