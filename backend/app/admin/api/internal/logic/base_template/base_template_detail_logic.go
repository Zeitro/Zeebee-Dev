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

type BaseTemplateDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Get base template detail
func NewBaseTemplateDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BaseTemplateDetailLogic {
	return &BaseTemplateDetailLogic{
		Logger: svcCtx.Logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BaseTemplateDetailLogic) BaseTemplateDetail(req *types.PathIdRequest) (resp *types.BaseTemplateDetailResponse, err error) {

	// Convert types to pb request
	rpcReq := &admin.GetBaseTemplateRequest{
		Data: &proto.IdRequest{Data: req.Id},
	}

	// Call RPC server
	rpcResp, err := l.svcCtx.BaseTemplateRpcClient.Detail(l.ctx, rpcReq)
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
