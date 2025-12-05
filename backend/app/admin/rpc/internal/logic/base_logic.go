package logic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/admin/rpc/internal/global"
	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/domains/zeebee/logic"
)

type BaseLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	UserLogic     *logic.UserLogic
	ApiLogic      *logic.ApiLogic
	ResourceLogic *logic.ResourceLogic
	RoleLogic     *logic.RoleLogic
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
		ctx:           ctx,
		svcCtx:        svcCtx,
		Logger:        logx.WithContext(ctx),
		UserLogic:     logic.NewUserLogic(svcCtx.AllRepos, redisCli, jwtConfig),
		ApiLogic:      logic.NewApiLogic(svcCtx.AllRepos),
		ResourceLogic: logic.NewResourceLogic(svcCtx.AllRepos),
		RoleLogic:     logic.NewRoleLogic(svcCtx.AllRepos),
	}
}
