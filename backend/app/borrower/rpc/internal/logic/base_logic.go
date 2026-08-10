package logic

import (
	"sync"

	"zeebee-backend/app/borrower/rpc/internal/global"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/logic"
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

	logicConfig := logic.AllLogicConfig{
		Repos:                 svcCtx.AllRepos,
		RedisCli:              redisCli,
		JwtConfig:             jwtConfig,
		EmailClient:           global.EmailClient,
		// US-12: tenant is the sole casbin executor; borrower no longer runs a
		// local enforcer, so it subscribes to no reload channel (AuthLogic is
		// never invoked here anymore).
		CurrentEnforcer:       nil,
		AuthEventTopicMap:     authEventTopicMap,
		Logger:                svcCtx.Logger,
	}

	l := &BaseLogic{

		svcCtx:    svcCtx,
		AllLogics: logic.NewAllLogics(logicConfig),
	}

	return l
}
