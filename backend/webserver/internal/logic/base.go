package logic

import (
	"context"

	"zeebee-backend/webserver/internal/svc"
)

type BaseLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBaseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BaseLogic {
	return &BaseLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}
