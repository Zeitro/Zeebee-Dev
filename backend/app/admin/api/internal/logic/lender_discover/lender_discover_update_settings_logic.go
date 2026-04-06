package lender_discover

import (
	"context"

	"zeebee-backend/app/admin/api/internal/svc"
	"zeebee-backend/app/admin/api/internal/types"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type LenderDiscoverUpdateSettingsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// lender discover update settings
func NewLenderDiscoverUpdateSettingsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LenderDiscoverUpdateSettingsLogic {
	return &LenderDiscoverUpdateSettingsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LenderDiscoverUpdateSettingsLogic) LenderDiscoverUpdateSettings(req *types.LenderDiscoverUpdateSettingsRequest) (resp *types.LenderDiscoverSettingsResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	if len(req.OpenaiModel) == 0 {
		req.OpenaiModel = "gpt-4o-mini"
	}
	rpcResp, err := l.svcCtx.LenderDiscoverRpcClient.UpdateSettings(l.ctx, &admin.LenderDiscoverUpdateSettingsRequest{
		Data: &admin.LenderDiscoverUpdateSettingsRequestData{
			ScheduleEnabled:    req.ScheduleEnabled,
			IntervalDays:       int32(req.IntervalDays),
			ScopeState:         req.ScopeState,
			ScopeExtra:         req.ScopeExtra,
			UserPrompt:         req.UserPrompt,
			OpenaiModel:        req.OpenaiModel,
			UpdatedByUserId:    0,
		},
	})
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	return &types.LenderDiscoverSettingsResponse{
		Data: protoSettingsToTypes(rpcResp.GetData()),
	}, nil
}
