package common

import (
	"context"

	"github.com/jinzhu/copier"

	"zeebee-backend/app/admin/api/internal/svc"
	"zeebee-backend/app/admin/api/internal/types"
	"zeebee-backend/app/admin/rpc/pb/admin"
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
		Logger: svcCtx.Logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SelectOptionsLogic) SelectOptions(req *types.GetSelectOptionRequest) (resp *types.GetSelectOptionResponse, err error) {
	var rpcResp *admin.OptionResponse
	resp = &types.GetSelectOptionResponse{
		Options: make([]types.SelectOptionEntity, 0),
	}
	rpcResp, err = l.svcCtx.CommonRpcClient.GetOpt(l.ctx, &admin.OptionRequest{
		Data: &admin.OptionRequestData{
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
