package base_template

import (
	"context"

	"zeebee-backend/app/admin/api/internal/svc"
	"zeebee-backend/app/admin/api/internal/types"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type BaseTemplatePreviewLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Preview
func NewBaseTemplatePreviewLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BaseTemplatePreviewLogic {
	return &BaseTemplatePreviewLogic{
		Logger: svcCtx.Logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BaseTemplatePreviewLogic) BaseTemplatePreview(req *types.BaseTemplatePreviewRequest) (resp *types.BaseTemplatePreviewResponse, err error) {
	// Convert types to pb request
	rpcReq := &admin.BaseTemplatePreviewRequest{
		Data: &admin.BaseTemplatePreviewRequestData{},
	}
	rpcReq.Data, err = copierx.Trans[types.BaseTemplatePreviewRequest, admin.BaseTemplatePreviewRequestData](req)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}

	// Call RPC server
	rpcResp, err := l.svcCtx.BaseTemplateRpcClient.Preview(l.ctx, rpcReq)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	// Convert pb response to types
	resp, err = copierx.Trans[admin.BaseTemplatePreviewResponseData, types.BaseTemplatePreviewResponse](rpcResp.Data)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
