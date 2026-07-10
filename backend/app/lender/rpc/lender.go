package rpc

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/zeromicro/go-zero/core/service"

	"zeebee-backend/app/lender/config"
	"zeebee-backend/app/lender/rpc/internal/global"
	"zeebee-backend/app/lender/rpc/internal/logic"
	apiServer "zeebee-backend/app/lender/rpc/internal/server/api"
	authServer "zeebee-backend/app/lender/rpc/internal/server/auth"
	commonServer "zeebee-backend/app/lender/rpc/internal/server/common"
	fileServer "zeebee-backend/app/lender/rpc/internal/server/file"
	loanServer "zeebee-backend/app/lender/rpc/internal/server/loan"
	loanDocumentServer "zeebee-backend/app/lender/rpc/internal/server/loandocument"
	organizationServer "zeebee-backend/app/lender/rpc/internal/server/organization"
	organizationSettingServer "zeebee-backend/app/lender/rpc/internal/server/organizationsetting"
	organizationSubscriptionServer "zeebee-backend/app/lender/rpc/internal/server/organizationsubscription"
	organizationTemplateServer "zeebee-backend/app/lender/rpc/internal/server/organizationtemplate"
	organizationUserServer "zeebee-backend/app/lender/rpc/internal/server/organizationuser"
	resourceServer "zeebee-backend/app/lender/rpc/internal/server/resource"
	roleServer "zeebee-backend/app/lender/rpc/internal/server/role"
	storageServer "zeebee-backend/app/lender/rpc/internal/server/storage"
	userServer "zeebee-backend/app/lender/rpc/internal/server/user"
	"zeebee-backend/app/lender/rpc/internal/svc"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/cron"
	"zeebee-backend/domains/zeebee/constants"
)

type RpcServer struct {
	s    *zrpc.RpcServer
	cron *cron.Engine
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
		lender.RegisterOrganizationTemplateServer(grpcServer, organizationTemplateServer.NewOrganizationTemplateServer(ctx))
		lender.RegisterOrganizationSettingServer(grpcServer, organizationSettingServer.NewOrganizationSettingServer(ctx))
		lender.RegisterAuthServer(grpcServer, authServer.NewAuthServer(ctx))
		lender.RegisterResourceServer(grpcServer, resourceServer.NewResourceServer(ctx))
		lender.RegisterApiServer(grpcServer, apiServer.NewApiServer(ctx))
		lender.RegisterOrganizationSubscriptionServer(grpcServer, organizationSubscriptionServer.NewOrganizationSubscriptionServer(ctx))

		if rpcConfig.Mode == service.DevMode || rpcConfig.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	rpcService := &RpcServer{
		s:    rpcServer,
		cron: cron.NewEngine(),
	}
	err = rpcService.RegisterCronTask(ctx, rpcConfig)
	if err != nil {
		return nil, err
	}

	return rpcService, nil
}

func (s *RpcServer) Start() {
	go func() {
		s.s.Start()
	}()
	go func() {
		s.cron.Start()
	}()
}

func (s *RpcServer) Stop() {
	s.s.Stop()
	s.cron.Stop()
}

func (s *RpcServer) RegisterCronTask(svcCtx *svc.ServiceContext, conf config.RpcConf) (err error) {
	if s.cron == nil {
		s.cron = cron.NewEngine()
	}
	baseLogic := logic.GetBaseLogic(svcCtx)
	for _, task := range conf.CronTaskConf {
		switch task.Name {
		case "update_organization_subscription_status":
			err = s.cron.Register(task.Name, task.Schedule, func() {
				ctx := context.Background()
				taskErr := baseLogic.OrganizationSubscriptionLogic.UpdateStatus(ctx, constants.UserTypeLender)
				if taskErr != nil {
					logx.WithContext(ctx).WithFields(
						logx.Field("type", "cron_task"),
						logx.Field("name", task.Name),
						logx.Field("schedule", task.Schedule),
					).Error(taskErr)
				}
			})
		}
	}
	return
}
