package commonlogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"

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

func (l *GetUsStateLogic) GetUsState(in *borrower.EmptyRequest) (resp *borrower.StringSliceResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()

	states, err := l.svcCtx.UsGeographyRepo.GetStateUsps(l.ctx)
	if err != nil {
		return nil, err
	}

	return &borrower.StringSliceResponse{
		Data: states,
	}, nil
}
