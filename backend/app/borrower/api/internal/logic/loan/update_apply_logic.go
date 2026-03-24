package loan

import (
	"context"

	"github.com/duke-git/lancet/v2/slice"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateApplyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// update apply
func NewUpdateApplyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateApplyLogic {
	return &UpdateApplyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateApplyLogic) UpdateApply(req *types.LoanApplyRequest) (resp *types.EmptyResponse, err error) {
	rpcReq := &borrower.LoanApplyRequest{
		Data: &borrower.LoanApplyRequestData{},
	}
	rpcReq.Data, err = copierx.Trans[types.LoanApplyRequest, borrower.LoanApplyRequestData](req)
	if err != nil {
		return nil, err
	}

	rpcApplicants := slice.UniqueByComparator(rpcReq.Data.Applicants, func(item *borrower.LoanApplicantApplyEntity, other *borrower.LoanApplicantApplyEntity) bool {
		return item.BorrowerApplicantId == other.BorrowerApplicantId
	})
	if len(rpcApplicants) == 0 {
		err = errorx.WithStack(errorx.New("No applicant"))
		return nil, err
	}
	rpcReq.Data.Applicants = rpcApplicants

	rpcReq.Data.LenderCompanyId = 1

	_, err = l.svcCtx.LoanRpcClient.UpdateApplyLoan(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}
	return &types.EmptyResponse{}, nil
}
