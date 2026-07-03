package base_template

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/admin/api/internal/svc"
	"zeebee-backend/app/admin/api/internal/types"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/proto"
)

type BaseTemplateCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Create base template
func NewBaseTemplateCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BaseTemplateCreateLogic {
	return &BaseTemplateCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BaseTemplateCreateLogic) BaseTemplateCreate(req *types.BaseTemplateCreateRequest) (resp *types.BaseTemplateDetailResponse, err error) {

	// Convert types to pb request
	rpcReq := &admin.CreateBaseTemplateRequest{
		Data: &admin.CreateBaseTemplateRequestData{},
	}
	rpcReq.Data, err = copierx.Trans[types.BaseTemplateCreateRequest, admin.CreateBaseTemplateRequestData](req)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}

	// Call RPC server
	rpcResp, err := l.svcCtx.BaseTemplateRpcClient.Create(l.ctx, rpcReq)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	// Convert pb response to types
	resp = &types.BaseTemplateDetailResponse{}
	entity, err := copierx.Trans[proto.BaseTemplateEntity, types.BaseTemplateEntity](rpcResp.Data, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	resp.BaseTemplateEntity = *entity

	return resp, nil
}
