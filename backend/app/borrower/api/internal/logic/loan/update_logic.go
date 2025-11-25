package loan

import (
	"context"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/jinzhu/copier"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// update loan
func NewUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLogic {
	return &UpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateLogic) Update(req *types.LoanUpdateRequest) (resp *types.EmptyResponse, err error) {
	rpcReq := &borrower.LoanUpdateRequest{
		Data: &borrower.LoanUpdateRequestData{},
	}
	resp = &types.EmptyResponse{}

	err = copier.CopyWithOption(&rpcReq.Data, req, copier.Option{
		DeepCopy:      true,
		IgnoreEmpty:   false,
		CaseSensitive: true,
	})
	if err != nil {
		return nil, errorx.GetFailedTransReqErr(err)
	}
	rpcApplicants := slice.UniqueByComparator(rpcReq.Data.Applicants, func(item *borrower.ThirdPartyLoanApplicantEditEntry, other *borrower.ThirdPartyLoanApplicantEditEntry) bool {
		return item.BorrowerApplicantId == other.BorrowerApplicantId
	})
	if len(rpcApplicants) == 0 {
		err = errorx.WithStack(errorx.New("No Applicant"))
		return nil, err
	}
	rpcReq.Data.Applicants = rpcApplicants

	// get borrower id
	claim, err := utils.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	rpcReq.Data.UserId = claim.UserId
	rpcReq.Data.OrganizationId = claim.OrganizationId

	_, err = l.svcCtx.LoanRpcClient.UpdateBorrowerLoan(l.ctx, rpcReq)

	return nil, err
}
