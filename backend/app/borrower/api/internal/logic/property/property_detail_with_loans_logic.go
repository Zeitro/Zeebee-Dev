package property

import (
	"context"

	"github.com/jinzhu/copier"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type PropertyDetailWithLoansLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// detail
func NewPropertyDetailWithLoansLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PropertyDetailWithLoansLogic {
	return &PropertyDetailWithLoansLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PropertyDetailWithLoansLogic) PropertyDetailWithLoans(req *types.PropertyDetailRequest) (resp *types.PropertyDetailWithLoansResponse, err error) {
	rpcReq := &borrower.PropertyDetailRequest{
		Data: &borrower.PropertyDetailRequestData{
			Id: req.Id,
		},
	}
	// get borrower id
	borrowerId, err := utils.GetAccountIdFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	rpcReq.Data.BorrowerId = borrowerId

	rpcResp, err := l.svcCtx.PropertyRpcClient.PropertyDetailWithLoans(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	resp = &types.PropertyDetailWithLoansResponse{}
	err = copier.CopyWithOption(resp, &rpcResp.Data, copier.Option{DeepCopy: true})
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}

	return resp, nil
}
