package commonlogic

import (
	"context"

	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/common/errorx"
)

type GetUsStateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUsStateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUsStateLogic {
	return &GetUsStateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUsStateLogic) GetUsState(in *admin.EmptyRequest) (resp *admin.StringSliceResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	states, err := l.svcCtx.UsGeographyRepo.GetStateUsps(l.ctx)
	if err != nil {
		return nil, err
	}

	return &admin.StringSliceResponse{
		Data: states,
	}, nil
}
