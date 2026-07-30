package api

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/admin/api/internal/svc"
	"zeebee-backend/app/admin/api/internal/types"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/proto"
)

type DetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// detail
func NewDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DetailLogic {
	return &DetailLogic{
		Logger: svcCtx.Logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DetailLogic) Detail(req *types.ApiDetailRequest) (resp *types.ApiEntity, err error) {
	resp = &types.ApiEntity{}
	rpcReq := &admin.ApiDetailRequest{
		Data: &admin.ApiDetailRequestData{},
	}

	rpcReq.Data, err = copierx.Trans[types.ApiDetailRequest, admin.ApiDetailRequestData](req)
	if err != nil {
		return nil, err
	}

	rpcResp, err := l.svcCtx.ApiRpcClient.Detail(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// trans response
	resp, err = copierx.Trans[proto.ApiEntity, types.ApiEntity](rpcResp.Data)
	if err != nil {
		return nil, err
	}

	return resp, err
}
