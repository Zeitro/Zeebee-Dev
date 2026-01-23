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

type BaseTemplateUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Update base template
func NewBaseTemplateUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BaseTemplateUpdateLogic {
	return &BaseTemplateUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BaseTemplateUpdateLogic) BaseTemplateUpdate(req *types.UpdateBaseTemplateRequest) (resp *types.BaseTemplateDetailResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()

	// Convert types to pb request
	rpcReq := &admin.UpdateBaseTemplateRequest{
		Data: &admin.UpdateBaseTemplateRequestData{},
	}
	rpcReq.Data, err = copierx.Trans[types.UpdateBaseTemplateRequest, admin.UpdateBaseTemplateRequestData](req)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}

	// Call RPC server
	rpcResp, err := l.svcCtx.BaseTemplateRpcClient.Update(l.ctx, rpcReq)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	// Convert pb response to types
	resp = &types.BaseTemplateDetailResponse{}
	entity, err := copierx.Trans[admin.BaseTemplateEntity, types.BaseTemplateEntity](rpcResp.Data, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	resp.BaseTemplateEntity = *entity

	return resp, nil
}
