package loan

import (
	"context"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
)

type UpdateDraftLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// update draft
func NewUpdateDraftLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateDraftLogic {
	return &UpdateDraftLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateDraftLogic) UpdateDraft(req *types.LoanApplyRequest) (resp *types.LoanApplyResponse, err error) {
	rpcReq := &borrower.LoanApplyRequest{
		Data: &borrower.LoanApplyRequestData{},
	}
	resp = &types.LoanApplyResponse{}
	rpcResp := &borrower.LoanApplyResponse{}
	rpcReq.Data, err = copierx.Trans[types.LoanApplyRequest, borrower.LoanApplyRequestData](req)
	if err != nil {
		return nil, err
	}

	rpcApplicants := slice.UniqueByComparator(rpcReq.Data.Applicants, func(item *borrower.LoanApplicantApplyEntity, other *borrower.LoanApplicantApplyEntity) bool {
		return item.BorrowerApplicantId == other.BorrowerApplicantId
	})
	rpcReq.Data.Applicants = rpcApplicants

	rpcReq.Data.LenderCompanyId = 1

	rpcResp, err = l.svcCtx.LoanRpcClient.SaveDraftLoan(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}
	resp, err = copierx.Trans[borrower.LoanApplyResponse, types.LoanApplyResponse](rpcResp)
	if err != nil {
		return nil, err
	}
	return
}
