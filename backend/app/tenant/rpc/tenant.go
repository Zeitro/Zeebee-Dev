package rpc

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"

	"zeebee-backend/app/tenant/config"
	"zeebee-backend/app/tenant/rpc/internal/logic"
	"zeebee-backend/app/tenant/rpc/internal/schedule"
	apiServer "zeebee-backend/app/tenant/rpc/internal/server/api"
	authServer "zeebee-backend/app/tenant/rpc/internal/server/auth"
	organizationServer "zeebee-backend/app/tenant/rpc/internal/server/organization"
	organizationEmailServer "zeebee-backend/app/tenant/rpc/internal/server/organizationemail"
	organizationSettingServer "zeebee-backend/app/tenant/rpc/internal/server/organizationsetting"
	organizationSubscriptionServer "zeebee-backend/app/tenant/rpc/internal/server/organizationsubscription"
	organizationTemplateServer "zeebee-backend/app/tenant/rpc/internal/server/organizationtemplate"
	organizationUserServer "zeebee-backend/app/tenant/rpc/internal/server/organizationuser"
	pricingPlanServer "zeebee-backend/app/tenant/rpc/internal/server/pricingplan"
	pricingPlanPolicyServer "zeebee-backend/app/tenant/rpc/internal/server/pricingplanpolicy"
	resourceServer "zeebee-backend/app/tenant/rpc/internal/server/resource"
	roleServer "zeebee-backend/app/tenant/rpc/internal/server/role"
	userServer "zeebee-backend/app/tenant/rpc/internal/server/user"
	"zeebee-backend/app/tenant/rpc/internal/svc"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/common/cron"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/constants"
	domainlenderdiscover "zeebee-backend/domains/zeebee/lenderdiscover"
)

// RpcServer wraps the tenant gRPC server.
type RpcServer struct {
	s      *zrpc.RpcServer
	svcCtx *svc.ServiceContext
	cron   *cron.Engine
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
		tenant.RegisterResourceServer(grpcServer, resourceServer.NewResourceServer(ctx))
		tenant.RegisterRoleServer(grpcServer, roleServer.NewRoleServer(ctx))
		tenant.RegisterPricingPlanServer(grpcServer, pricingPlanServer.NewPricingPlanServer(ctx))
		tenant.RegisterPricingPlanPolicyServer(grpcServer, pricingPlanPolicyServer.NewPricingPlanPolicyServer(ctx))

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
	// set service name for log and copy metadata to outgoing ctx
	rpcServer.AddOptions(grpc.UnaryInterceptor(func(ctx context.Context, req any,
		info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		ctx = logx.ContextWithFields(ctx, logx.Field("service_name", c.RpcConf.Name))
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			ctx = metadata.NewOutgoingContext(ctx, md.Copy())
		}
		return handler(ctx, req)
	}))
	// set err log stack
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

	server := &RpcServer{
		s:      rpcServer,
		svcCtx: ctx,
		cron:   cron.NewEngine(),
	}
	ctx.Cron = server.cron
	if err := server.RegisterCronTask(ctx, rpcConfig, logger); err != nil {
		return nil, err
	}
	return server, nil
}

func (s *RpcServer) Start() {
	go func() {
		s.s.Start()
	}()
	s.cron.Start()
}

func (s *RpcServer) Stop() {
	s.s.Stop()
	s.cron.Stop()
}

func (s *RpcServer) RegisterCronTask(svcCtx *svc.ServiceContext, conf config.RpcConf, logger logx.Logger) error {
	baseLogic := logic.GetBaseLogic(svcCtx)
	for _, task := range conf.CronTaskConf {
		var userType constants.UserType
		switch task.Name {
		case "update_admin_subscription_status":
			userType = constants.UserTypeAdmin
		case "update_lender_subscription_status":
			userType = constants.UserTypeLender
		case "update_borrower_subscription_status":
			userType = constants.UserTypeBorrower
		default:
			continue
		}
		name, schedule, statusUserType := task.Name, task.Schedule, userType
		if err := s.cron.Register(name, schedule, func() {
			if taskErr := baseLogic.OrganizationSubscriptionLogic.UpdateStatus(context.Background(), statusUserType); taskErr != nil {
				logger.WithFields(logx.Field("type", "cron_task"), logx.Field("name", name), logx.Field("schedule", schedule)).Error(taskErr)
			}
		}); err != nil {
			return err
		}
	}
	if cfg, _, err := domainlenderdiscover.LoadConfig(context.Background(), svcCtx.OrganizationSettingRepo); err != nil {
		logger.Errorf("lender discover cron register: load settings: %v", err)
	} else if cfg.ScheduleEnabled {
		if err := s.cron.Register(schedule.TaskName, schedule.SpecFromInterval(cfg.IntervalDays), func() {
			schedule.RunScheduled(context.Background(), svcCtx)
		}); err != nil {
			logger.Errorf("lender discover cron register: %v", err)
		}
	}
	return nil
}
