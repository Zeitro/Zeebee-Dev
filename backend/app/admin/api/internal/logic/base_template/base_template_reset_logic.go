package base_template

import (
	"context"

	"zeebee-backend/app/admin/api/internal/svc"
	"zeebee-backend/app/admin/api/internal/types"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type BaseTemplateResetLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Reset
func NewBaseTemplateResetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BaseTemplateResetLogic {
	return &BaseTemplateResetLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BaseTemplateResetLogic) BaseTemplateReset(req *types.PathIdRequest) (resp *types.EmptyResponse, err error) {
	// Convert types to pb request
	rpcReq := &admin.IdRequest{
		Data: req.Id,
	}

	// Call RPC server
	_, err = l.svcCtx.BaseTemplateRpcClient.Reset(l.ctx, rpcReq)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	return &types.EmptyResponse{}, nil
}
