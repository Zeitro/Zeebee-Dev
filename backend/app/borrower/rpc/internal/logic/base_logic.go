package logic

import (
	"context"
	"sync"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/threading"

	"zeebee-backend/app/borrower/rpc/internal/global"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/logic"
)

var baseLogic *BaseLogic
var onceBase sync.Once

type BaseLogic struct {
	// ctx    context.Context
	svcCtx *svc.ServiceContext
	// logx.Logger
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
		AccessSecret:  svcCtx.Config.JwtAuth.AccessSecret,
		AccessExpire:  svcCtx.Config.JwtAuth.AccessExpire,
		RefreshSecret: svcCtx.Config.JwtAuth.RefreshSecret,
		RefreshExpire: svcCtx.Config.JwtAuth.RefreshExpire,
	}
	currentEnforcer := global.BorrowerEnforcer
	authEventTopicMap := map[constants.UserType]string{
		constants.UserTypeBorrower: svcCtx.Config.CasbinConf.EventTopic.BorrowerTopic,
	}

	logicConfig := logic.AllLogicConfig{
		Repos:                 svcCtx.AllRepos,
		RedisCli:              redisCli,
		JwtConfig:             jwtConfig,
		EmailClient:           global.EmailClient,
		CurrentEnforcer:       currentEnforcer,
		CurrentAuthEventTopic: svcCtx.Config.CasbinConf.EventTopic.BorrowerTopic,
		AuthEventTopicMap:     authEventTopicMap,
	}

	l := &BaseLogic{

		svcCtx:    svcCtx,
		AllLogics: logic.NewAllLogics(logicConfig),
	}

	threading.GoSafe(func() {
		ctx := context.Background()
		err := l.AuthLogic.RebuildPolicy(ctx, constants.UserTypeBorrower)
		if err != nil {
			logx.WithContext(ctx).WithFields(
				logx.Field("stacks", errorx.GetErrorStack(err)),
				logx.Field("user_type", constants.UserTypeBorrower),
			).Error("fail to rebuild borrower policy")
		}
	})

	return l
}
