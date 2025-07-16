package loadlogic

import (
	"context"

	"zeebee-backend/app/load/rpc/internal/svc"
	"zeebee-backend/app/load/rpc/pb/load"

	"github.com/zeromicro/go-zero/core/logx"
)

type PingLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PingLogic {
	return &PingLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PingLogic) Ping(in *load.Request) (*load.Response, error) {
	// todo: add your logic here and delete this line

	return &load.Response{}, nil
}
