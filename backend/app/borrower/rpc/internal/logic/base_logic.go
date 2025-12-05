package logic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/rpc/internal/global"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/domains/zeebee/logic"
)

type BaseLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	UserLogic *logic.UserLogic
	RoleLogic *logic.RoleLogic
}

func NewBaseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BaseLogic {
	redisCli := global.GetRedis()
	jwtConfig := logic.JwtConfig{
		AccessSecret:  svcCtx.Config.JwtAuth.AccessSecret,
		AccessExpire:  svcCtx.Config.JwtAuth.AccessExpire,
		RefreshSecret: svcCtx.Config.JwtAuth.RefreshSecret,
		RefreshExpire: svcCtx.Config.JwtAuth.RefreshExpire,
	}
	return &BaseLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		UserLogic: logic.NewUserLogic(svcCtx.AllRepos, redisCli, jwtConfig),
		RoleLogic: logic.NewRoleLogic(svcCtx.AllRepos),
	}
}
