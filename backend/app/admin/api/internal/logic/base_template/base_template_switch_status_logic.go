package base_template

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/admin/api/internal/svc"
	"zeebee-backend/app/admin/api/internal/types"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
)

type BaseTemplateSwitchStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Switch base template status
func NewBaseTemplateSwitchStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BaseTemplateSwitchStatusLogic {
	return &BaseTemplateSwitchStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BaseTemplateSwitchStatusLogic) BaseTemplateSwitchStatus(req *types.BaseTemplateSwitchStatusRequest) (resp *types.EmptyResponse, err error) {

	// Convert types to pb request
	rpcReq := &admin.SwitchBaseTemplateStatusRequest{
		Data: &admin.SwitchBaseTemplateStatusRequestData{},
	}
	rpcReq.Data, err = copierx.Trans[types.BaseTemplateSwitchStatusRequest, admin.SwitchBaseTemplateStatusRequestData](req)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}

	// Call RPC server
	_, err = l.svcCtx.BaseTemplateRpcClient.SwitchStatus(l.ctx, rpcReq)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	return &types.EmptyResponse{}, nil
}
