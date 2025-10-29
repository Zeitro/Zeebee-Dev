package common

import (
	"context"

	"zeebee-backend/app/lender/api/internal/svc"
	"zeebee-backend/app/lender/api/internal/types"
	"zeebee-backend/app/lender/rpc/pb/lender"

	"github.com/zeromicro/go-zero/core/logx"
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
	var rpcResp *lender.StringSliceResponse
	rpcResp, err = l.svcCtx.CommonRpcClient.GetUsCities(l.ctx, &lender.UsCitiesRequest{
		Data: &lender.UsCitiesRequestData{
			StateUsps: req.State,
		},
	})
	return &types.GetUsCitiesResponse{
		List: rpcResp.Data,
	}, nil
}
