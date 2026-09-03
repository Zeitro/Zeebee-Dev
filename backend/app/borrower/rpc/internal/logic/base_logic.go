package logic

import (
	"context"
	"sync"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"zeebee-backend/app/borrower/rpc/internal/global"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/utils"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/logic"
	"zeebee-backend/domains/zeebee/payment"
)

var baseLogic *BaseLogic
var onceBase sync.Once

type BaseLogic struct {
	// ctx    context.Context
	svcCtx *svc.ServiceContext
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
	authEventTopicMap := map[constants.UserType]string{
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
		EmailClient: global.EmailClient,
		// US-12: tenant is the sole casbin executor; borrower no longer runs a
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
	jwtClaim, err := copierx.Trans[tenant.GetClaimResponseData, utils.JwtUser](claim.Data)
	if err != nil {
		return utils.JwtUser{}, err
	}
	return *jwtClaim, nil
}
