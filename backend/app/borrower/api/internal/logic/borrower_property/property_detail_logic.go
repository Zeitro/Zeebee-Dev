package borrower_property

import (
	"context"

	"github.com/jinzhu/copier"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type PropertyDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// detail
func NewPropertyDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PropertyDetailLogic {
	return &PropertyDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PropertyDetailLogic) PropertyDetail(req *types.PropertyDetailRequest) (resp *types.PropertyDetailResponse, err error) {
	rpcReq := &borrower.PropertyDetailRequest{
		Data: &borrower.PropertyDetailRequestData{
			Id: req.Id,
		},
	}

	rpcResp, err := l.svcCtx.BorrowerPropertyRpcClient.Detail(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	resp = &types.PropertyDetailResponse{}
	err = copier.CopyWithOption(resp, &rpcResp.Data, copier.Option{DeepCopy: true})
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}

	return resp, nil
}
