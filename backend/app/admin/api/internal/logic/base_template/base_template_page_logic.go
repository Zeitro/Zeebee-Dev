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

type BaseTemplatePageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Base template page list
func NewBaseTemplatePageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BaseTemplatePageLogic {
	return &BaseTemplatePageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BaseTemplatePageLogic) BaseTemplatePage(req *types.BaseTemplatePageRequest) (resp *types.BaseTemplatePageResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()

	// Convert types to pb request
	rpcReq := &admin.BaseTemplatePageRequest{
		Data: &admin.BaseTemplatePageRequestData{},
	}
	rpcReq.Data, err = copierx.Trans[types.BaseTemplatePageRequest, admin.BaseTemplatePageRequestData](req)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}

	// Call RPC server
	rpcResp, err := l.svcCtx.BaseTemplateRpcClient.Page(l.ctx, rpcReq)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	// Convert pb response to types
	resp = &types.BaseTemplatePageResponse{}
	resp.Page = rpcResp.Data.Page
	resp.Limit = rpcResp.Data.Limit
	resp.Total = rpcResp.Data.Total
	resp.TotalPage = rpcResp.Data.TotalPage

	// Convert list items
	for _, item := range rpcResp.Data.List {
		entity, err := copierx.Trans[admin.BaseTemplatePageEntity, types.BaseTemplatePageEntity](item, copierx.TimeToRFC3339Str())
		if err != nil {
			return nil, errorx.GetFailedTransRespErr(err)
		}
		resp.List = append(resp.List, entity)
	}

	return resp, nil
}
