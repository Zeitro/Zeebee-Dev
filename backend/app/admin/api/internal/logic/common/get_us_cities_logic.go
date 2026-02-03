package common

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/admin/api/internal/svc"
	"zeebee-backend/app/admin/api/internal/types"
	"zeebee-backend/app/admin/rpc/pb/admin"
)

type GetUsCitiesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// get us city
func NewGetUsCitiesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUsCitiesLogic {
	return &GetUsCitiesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUsCitiesLogic) GetUsCities(req *types.GetUsCitiesRequest) (resp *types.GetUsCitiesResponse, err error) {
	var rpcResp *admin.StringSliceResponse
	rpcResp, err = l.svcCtx.CommonRpcClient.GetUsCities(l.ctx, &admin.UsCitiesRequest{
		Data: &admin.UsCitiesRequestData{
			StateUsps: req.State,
		},
	})
	return &types.GetUsCitiesResponse{
		List: rpcResp.Data,
	}, nil
}
