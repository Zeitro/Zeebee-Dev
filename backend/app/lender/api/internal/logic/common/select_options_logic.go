package common

import (
	"context"

	"github.com/jinzhu/copier"

	"zeebee-backend/app/lender/api/internal/svc"
	"zeebee-backend/app/lender/api/internal/types"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type SelectOptionsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// select options
func NewSelectOptionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SelectOptionsLogic {
	return &SelectOptionsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SelectOptionsLogic) SelectOptions(req *types.GetSelectOptionRequest) (resp *types.GetSelectOptionResponse, err error) {
	var rpcResp *lender.OptionResponse
	resp = &types.GetSelectOptionResponse{
		Options: make([]types.SelectOptionEntity, 0),
	}
	rpcResp, err = l.svcCtx.CommonRpcClient.GetOpt(l.ctx, &lender.OptionRequest{
		Data: &lender.OptionRequestData{
			Name: req.Name,
		},
	})
	if err != nil {
		return nil, err
	}
	err = copier.Copy(&resp.Options, rpcResp.Data)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}

	return
}
