package lender_discover

import (
	"context"

	"zeebee-backend/app/admin/api/internal/svc"
	"zeebee-backend/app/admin/api/internal/types"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type LenderDiscoverPreviewLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// lender discover preview (OpenAI + dedupe, no insert)
func NewLenderDiscoverPreviewLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LenderDiscoverPreviewLogic {
	return &LenderDiscoverPreviewLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LenderDiscoverPreviewLogic) LenderDiscoverPreview(_ *types.EmptyRequest) (resp *types.LenderDiscoverPreviewResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	rpcResp, err := l.svcCtx.LenderDiscoverRpcClient.Preview(context.WithoutCancel(l.ctx), &admin.EmptyRequest{})
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	var data types.LenderDiscoverPreviewData
	if rpcResp.GetData() != nil {
		data = protoPreviewDataToTypes(rpcResp.GetData())
	}
	return &types.LenderDiscoverPreviewResponse{Data: data}, nil
}
