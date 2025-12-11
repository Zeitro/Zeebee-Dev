package loan

import (
	"context"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/jinzhu/copier"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/errorx"

	"github.com/zeromicro/go-zero/core/logx"
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

	err = copier.CopyWithOption(&rpcReq.Data, req, copier.Option{
		DeepCopy:      true,
		IgnoreEmpty:   true,
		CaseSensitive: true,
	})
	if err != nil {
		return nil, errorx.GetFailedTransReqErr(err)
	}

	rpcApplicants := slice.UniqueByComparator(rpcReq.Data.Applicants, func(item *borrower.LoanApplicantApplyEntry, other *borrower.LoanApplicantApplyEntry) bool {
		return item.BorrowerApplicantId == other.BorrowerApplicantId
	})
	rpcReq.Data.Applicants = rpcApplicants

	rpcReq.Data.LenderCompanyId = 1

	rpcResp, err = l.svcCtx.LoanRpcClient.SaveDraftLoan(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}
	err = copier.Copy(resp, &rpcResp.Data)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}

	return
}
