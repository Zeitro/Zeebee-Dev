package rpc

import (
	"context"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/zeromicro/go-zero/core/service"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	apiServer "zeebee-backend/app/borrower/rpc/internal/server/api"
	borrowerApplicantServer "zeebee-backend/app/borrower/rpc/internal/server/borrowerapplicant"
	borrowerEntityServer "zeebee-backend/app/borrower/rpc/internal/server/borrowerentity"
	borrowerPropertyServer "zeebee-backend/app/borrower/rpc/internal/server/borrowerproperty"
	borrowerPropertyTransactionServer "zeebee-backend/app/borrower/rpc/internal/server/borrowerpropertytransaction"
	documentServer "zeebee-backend/app/borrower/rpc/internal/server/document"
	driveServer "zeebee-backend/app/borrower/rpc/internal/server/drive"
	organizationServer "zeebee-backend/app/borrower/rpc/internal/server/organization"
	organizationEmailServer "zeebee-backend/app/borrower/rpc/internal/server/organizationemail"
	organizationSubscriptionServer "zeebee-backend/app/borrower/rpc/internal/server/organizationsubscription"
	zeitroServer "zeebee-backend/app/borrower/rpc/internal/server/zeitro"
	"zeebee-backend/common/cron"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/templater"
	"zeebee-backend/domains/zeebee/templater/convertor"
	"zeebee-backend/domains/zeebee/templater/loader"

	loanApplicantServer "zeebee-backend/app/borrower/rpc/internal/server/loanapplicant"

	"zeebee-backend/app/borrower/config"
	"zeebee-backend/app/borrower/rpc/internal/global"
	authServer "zeebee-backend/app/borrower/rpc/internal/server/auth"
	commonServer "zeebee-backend/app/borrower/rpc/internal/server/common"
	loanServer "zeebee-backend/app/borrower/rpc/internal/server/loan"
	loanDocumentServer "zeebee-backend/app/borrower/rpc/internal/server/loandocument"
	organizationUserServer "zeebee-backend/app/borrower/rpc/internal/server/organizationuser"
	resourceServer "zeebee-backend/app/borrower/rpc/internal/server/resource"
	roleServer "zeebee-backend/app/borrower/rpc/internal/server/role"
	storageServer "zeebee-backend/app/borrower/rpc/internal/server/storage"
	userServer "zeebee-backend/app/borrower/rpc/internal/server/user"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
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
	zrpc.SetClientSlowThreshold(time.Second * 2)

	ctx := svc.NewServiceContext(rpcConfig)

	baseLogic := logic.NewBaseLogic(ctx)
	tplLoader := loader.NewDefaultTemplateLoader(baseLogic.BaseTemplateLogic, baseLogic.OrganizationTemplateLogic)
	tplConvertor := convertor.NewDefaultConvertor(rpcConfig.PdfGeneratorUrl)
	tpl := templater.NewEngine(tplLoader, tplConvertor)
	global.Templater = tpl

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
		borrower.RegisterOrganizationUserServer(grpcServer, organizationUserServer.NewOrganizationUserServer(ctx))
		borrower.RegisterRoleServer(grpcServer, roleServer.NewRoleServer(ctx))
		borrower.RegisterOrganizationServer(grpcServer, organizationServer.NewOrganizationServer(ctx))
		borrower.RegisterOrganizationEmailServer(grpcServer, organizationEmailServer.NewOrganizationEmailServer(ctx))
		borrower.RegisterZeitroServer(grpcServer, zeitroServer.NewZeitroServer(ctx))
		borrower.RegisterAuthServer(grpcServer, authServer.NewAuthServer(ctx))
		borrower.RegisterResourceServer(grpcServer, resourceServer.NewResourceServer(ctx))
		borrower.RegisterBorrowerPropertyTransactionServer(grpcServer, borrowerPropertyTransactionServer.NewBorrowerPropertyTransactionServer(ctx))
		borrower.RegisterDocumentServer(grpcServer, documentServer.NewDocumentServer(ctx))
		borrower.RegisterDriveServer(grpcServer, driveServer.NewDriveServer(ctx))
		borrower.RegisterApiServer(grpcServer, apiServer.NewApiServer(ctx))
		borrower.RegisterOrganizationSubscriptionServer(grpcServer, organizationSubscriptionServer.NewOrganizationSubscriptionServer(ctx))

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
				taskErr := baseLogic.OrganizationSubscriptionLogic.UpdateStatus(ctx, constants.UserTypeBorrower)
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
