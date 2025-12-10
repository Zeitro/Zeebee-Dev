package loan

import (
	"context"

	"github.com/jinzhu/copier"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/lender/api/internal/svc"
	"zeebee-backend/app/lender/api/internal/types"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/errorx"
)

type UpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// update
func NewUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLogic {
	return &UpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateLogic) Update(req *types.LoanUpdateRequest) (resp *types.LoanUpdateResponse, err error) {
	rpcReq := &lender.LoanUpdateRequest{
		Data: &lender.LoanUpdateRequestData{},
	}
	resp = &types.LoanUpdateResponse{}
	err = copier.CopyWithOption(&rpcReq.Data, req, copier.Option{
		IgnoreEmpty:   false,
		DeepCopy:      true,
		CaseSensitive: true,
	})
	if err != nil {
		return nil, errorx.GetFailedTransReqErr(err)
	}
	// claim, err := utils.GetUserClaimFromCtx(l.ctx)
	// if err != nil {
	// 	return nil, err
	// }

	// get lender id
	// rpcReq.Data.OrganizationId = claim.OrganizationId
	// rpcReq.Data.UserId = claim.UserId

	rpcResp, err := l.svcCtx.LoanRpcClient.UpdateLoan(l.ctx, rpcReq)
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	err = copier.Copy(resp, &rpcResp.Data)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	return resp, err
}
