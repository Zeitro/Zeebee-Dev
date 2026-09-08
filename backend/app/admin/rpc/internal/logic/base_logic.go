package logic

import (
	"context"
	"sync"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"

	"zeebee-backend/app/admin/rpc/internal/global"
	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/utils"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/logic"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/payment"
)

var ErrorInvalidRequest = errorx.New("invalid request")

var baseLogic *BaseLogic
var onceBase sync.Once

type BaseLogic struct {
	// ctx    context.Context
	svcCtx *svc.ServiceContext
	//
	*logic.AllLogics
}

func GetBaseLogic(svcCtx *svc.ServiceContext) *BaseLogic {
	onceBase.Do(func() {
		baseLogic = NewBaseLogic(svcCtx)
	})
	return baseLogic
}

func NewBaseLogic(svcCtx *svc.ServiceContext) *BaseLogic {
	redisCli := global.GetRedis()
	jwtConfig := dto.JwtConfig{
		AccessSecret:  svcCtx.Config.JwtConf.AccessSecret,
		AccessExpire:  svcCtx.Config.JwtConf.AccessExpire,
		RefreshSecret: svcCtx.Config.JwtConf.RefreshSecret,
		RefreshExpire: svcCtx.Config.JwtConf.RefreshExpire,
	}
	emailClient := global.EmailClient

	authEventTopicMap := map[constants.UserType]string{
		constants.UserTypeAdmin:    svcCtx.Config.CasbinConf.EventTopic.AdminTopic,
		constants.UserTypeLender:   svcCtx.Config.CasbinConf.EventTopic.LenderTopic,
		constants.UserTypeBorrower: svcCtx.Config.CasbinConf.EventTopic.BorrowerTopic,
	}
	paymentGateway, err := payment.NewGateway(payment.Config{
		SecretKey:     svcCtx.Config.StripeConf.SecretKey,
		WebhookSecret: svcCtx.Config.StripeConf.WebhookSecret,
		SuccessUrl:    svcCtx.Config.StripeConf.SuccessUrl,
		CancelUrl:     svcCtx.Config.StripeConf.CancelUrl,
		Currency:      svcCtx.Config.StripeConf.Currency,
	})
	if err != nil {
		panic(err)
	}

	logicConfig := logic.AllLogicConfig{
		Repos:       svcCtx.AllRepos,
		RedisCli:    redisCli,
		JwtConfig:   jwtConfig,
		EmailClient: emailClient,
		// US-12: tenant is the sole casbin executor; admin no longer runs a
		// local enforcer, so it subscribes to no reload channel (AuthLogic is
		// never invoked here anymore).
		CurrentEnforcer:   nil,
		AuthEventTopicMap: authEventTopicMap,
		Logger:            svcCtx.Logger,
		PaymentGateway:        paymentGateway,
	}
	l := &BaseLogic{
		svcCtx:    svcCtx,
		AllLogics: logic.NewAllLogics(logicConfig),
	}

	return l
}

func (l *BaseLogic) GetUserClaimFromCtx(ctx context.Context) (utils.JwtUser, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok || len(md.Get(utils.JwtTokenCtxKey)) == 0 || len(md.Get(utils.UserTypeCtxKey)) == 0 {
		return utils.JwtUser{}, status.Error(codes.Unauthenticated, "not login")
	}
	claim, err := l.svcCtx.TenantClient.User.GetClaim(ctx, &tenant.GetClaimRequest{Data: &tenant.GetClaimRequestData{
		AccessToken: md.Get(utils.JwtTokenCtxKey)[0],
		UserType:    md.Get(utils.UserTypeCtxKey)[0],
	}})
	if err != nil {
		return utils.JwtUser{}, err
	}
	if claim.Data == nil {
		return utils.JwtUser{}, status.Error(codes.Unauthenticated, "invalid token")
	}
	return utils.JwtUser{
		UserId: claim.Data.UserId, OrganizationId: claim.Data.OrganizationId,
		OrganizationUserId: claim.Data.OrganizationUserId, Username: claim.Data.Username,
		Role: claim.Data.Role, UserType: constants.UserType(claim.Data.UserType),
	}, nil
}

func (l *BaseLogic) CreateLoanOperationLog(ctx context.Context, tx *gorm.DB, loanId, operatorId int64, from, to dto.LoanOperationStatus, action constants.LoanOperationAction, metadata *dto.LoanOperationMetadata) (err error) {
	defer func() {
		if err != nil {
			l.svcCtx.Logger.WithContext(ctx).WithFields(
				logx.Field("stacks", errorx.GetErrorStack(err)),
				logx.Field("loan_id", loanId),
				logx.Field("operator_id", operatorId),
				logx.Field("from", from),
				logx.Field("to", to),
				logx.Field("action", action),
				logx.Field("metadata", metadata),
				logx.Field("OperatorType", constants.LoanOperationOperatorTypeAdmin),
			).Error("fail to create loan operation log")
		}
	}()
	operatorName := ""
	if claim, claimErr := l.GetUserClaimFromCtx(ctx); claimErr == nil {
		operatorName = claim.Username
	}

	operation := &model.LoanOperation{
		LoanId:       loanId,
		OperatorId:   operatorId,
		OperatorName: operatorName,
		OperatorType: constants.LoanOperationOperatorTypeAdmin,
		Action:       action,
		From:         from,
		To:           to,
		Metadata:     metadata,
	}
	_, err = l.svcCtx.LoanOperationRepo.CreateWithTx(ctx, tx, operation, false)
	if err != nil {
		return errorx.Wrap(err, "fail to create operation log")
	}

	return nil
}
