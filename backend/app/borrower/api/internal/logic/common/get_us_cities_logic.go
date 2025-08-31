package common

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
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
	var rpcResp *borrower.StringSliceResponse
	rpcResp, err = l.svcCtx.CommonRpcClient.GetUsCities(l.ctx, &borrower.UsCitiesRequest{
		Data: &borrower.UsCitiesRequestData{
			StateUsps: req.State,
		},
	})
	return &types.GetUsCitiesResponse{
		List: rpcResp.Data,
	}, nil
}
