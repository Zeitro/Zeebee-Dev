package api

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

type CreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// create
func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	return &CreateLogic{
		Logger: svcCtx.Logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateLogic) Create(req *types.ApiCreateRequest) (resp *types.ApiEntity, err error) {
	resp = &types.ApiEntity{}
	rpcReq := &admin.ApiCreateRequest{
		Data: &admin.ApiCreateRequestData{},
	}

	rpcReq.Data, err = copierx.Trans[types.ApiCreateRequest, admin.ApiCreateRequestData](req)
	if err != nil {
		return nil, err
	}

	rpcResp, err := l.svcCtx.ApiRpcClient.Create(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// trans response
	resp, err = copierx.Trans[proto.ApiEntity, types.ApiEntity](rpcResp.Data)

	if err != nil {
		return nil, errorx.GetFailedTransReqErr(err)
	}
	return resp, err
}
