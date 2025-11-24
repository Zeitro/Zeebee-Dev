package loan

import (
	"context"

	"github.com/jinzhu/copier"

	"zeebee-backend/app/lender/api/internal/svc"
	"zeebee-backend/app/lender/api/internal/types"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type PageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// page
func NewPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PageLogic {
	return &PageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PageLogic) Page(req *types.LoanPageRequest) (resp *types.LoanPageResponse, err error) {
	rpcReq := &lender.LoanPageRequest{
		Data: &lender.LoanPageRequestData{
			Filter: &lender.LoanPageRequestFilter{},
		},
	}
	resp = &types.LoanPageResponse{}

	err = copier.CopyWithOption(&rpcReq.Data, req, copier.Option{
		IgnoreEmpty: true,
		DeepCopy:    true,
	})
	if err != nil {
		return nil, errorx.GetFailedTransReqErr(err)
	}
	// get lender id
	claim, err := utils.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	rpcReq.Data.Filter.OrganizationId = claim.OrganizationId
	rpcResp, err := l.svcCtx.LoanRpcClient.LoanPage(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}
	err = copier.Copy(resp, &rpcResp.Data)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	return resp, err
}
