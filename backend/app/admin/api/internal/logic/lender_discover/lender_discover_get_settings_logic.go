package lender_discover

import (
	"context"

	"zeebee-backend/app/admin/api/internal/svc"
	"zeebee-backend/app/admin/api/internal/types"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type LenderDiscoverGetSettingsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// lender discover get settings
func NewLenderDiscoverGetSettingsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LenderDiscoverGetSettingsLogic {
	return &LenderDiscoverGetSettingsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LenderDiscoverGetSettingsLogic) LenderDiscoverGetSettings(_ *types.EmptyRequest) (resp *types.LenderDiscoverSettingsResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	rpcResp, err := l.svcCtx.LenderDiscoverRpcClient.GetSettings(l.ctx, &admin.EmptyRequest{})
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	return &types.LenderDiscoverSettingsResponse{
		Data: protoSettingsToTypes(rpcResp.GetData()),
	}, nil
}
