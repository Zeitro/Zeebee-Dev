package logic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

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
}

func NewBaseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BaseLogic {
	return &BaseLogic{
		ctx:           ctx,
		svcCtx:        svcCtx,
		Logger:        logx.WithContext(ctx),
		UserLogic:     logic.NewUserLogic(svcCtx.AllRepos),
		ApiLogic:      logic.NewApiLogic(svcCtx.AllRepos),
		ResourceLogic: logic.NewResourceLogic(svcCtx.AllRepos),
	}
}
