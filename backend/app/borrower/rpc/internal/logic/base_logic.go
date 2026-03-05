package logic

import (
	"sync"

	gormadapter "github.com/casbin/gorm-adapter/v3"

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
	// enforcerMap := map[constants.UserType]*rbac.Enforcer{
	// 	constants.UserTypeAdmin:    global.AdminEnforcer,
	// 	constants.UserTypeLender:   global.LenderEnforcer,
	// 	constants.UserTypeBorrower: global.BorrowerEnforcer,
	// }
	currentEnforcer := global.BorrowerEnforcer
	_ = currentEnforcer.LoadFilteredPolicy(gormadapter.Filter{
		Ptype: []string{"p"},
		V0:    nil,
		V1:    []string{string(constants.UserTypeBorrower)},
		V2:    nil,
		V3:    nil,
		V4:    nil,
		V5:    nil,
	})
	return &BaseLogic{
		// ctx:       ctx,
		svcCtx: svcCtx,
		// Logger:    logx.WithContext(ctx),
		AllLogics: logic.NewAllLogics(svcCtx.AllRepos, redisCli, jwtConfig, global.EmailClient, currentEnforcer, nil),
	}
}
