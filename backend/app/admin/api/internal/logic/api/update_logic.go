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

type UpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// update
func NewUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLogic {
	return &UpdateLogic{
		Logger: svcCtx.Logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateLogic) Update(req *types.ApiUpdateRequest) (resp *types.ApiEntity, err error) {
	resp = &types.ApiEntity{}
	rpcReq := &admin.ApiUpdateRequest{
		Data: &admin.ApiUpdateRequestData{},
	}
	rpcReq.Data, err = copierx.Trans[types.ApiUpdateRequest, admin.ApiUpdateRequestData](req)
	if err != nil {
		return nil, err
	}

	rpcResp, err := l.svcCtx.ApiRpcClient.Update(l.ctx, rpcReq)
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
