package lenderdiscover

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/errorx"
)

type LatestRunLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLatestRunLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LatestRunLogic {
	return &LatestRunLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *LatestRunLogic) LatestRun(_ *admin.EmptyRequest) (resp *admin.LenderDiscoverRunResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	row, err := l.svcCtx.LenderDiscoverRunRepo.Latest(l.ctx)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return &admin.LenderDiscoverRunResponse{Data: &admin.LenderDiscoverRunEntity{}}, nil
	}
	return &admin.LenderDiscoverRunResponse{Data: runToProto(row)}, nil
}
