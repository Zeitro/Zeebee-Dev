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

type CreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewCreateLogic create loan
func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	return &CreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateLogic) Create(req *types.LoanCreateRequest) (resp *types.LoanCreateResponse, err error) {
	rpcReq := &borrower.LoanCreateRequest{
		Data: &borrower.LoanCreateRequestData{},
	}
	resp = &types.LoanCreateResponse{}
	rpcResp := &borrower.LoanCreateResponse{}

	rpcReq.Data, err = copierx.Trans[types.LoanCreateRequest, borrower.LoanCreateRequestData](req)
	if err != nil {
		return nil, err
	}
	// handle applicants
	rpcApplicants := slice.UniqueByComparator(rpcReq.Data.Applicants, func(item *borrower.ThirdPartyLoanApplicantEditEntity, other *borrower.ThirdPartyLoanApplicantEditEntity) bool {
		return item.BorrowerApplicantId == other.BorrowerApplicantId
	})
	if len(rpcApplicants) == 0 {
		err = errorx.WithStack(errorx.New("No applicant"))
		return nil, err
	}
	rpcReq.Data.Applicants = rpcApplicants

	rpcResp, err = l.svcCtx.LoanRpcClient.CreateBorrowerLoan(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}
	resp, err = copierx.Trans[borrower.LoanCreateResponseData, types.LoanCreateResponse](rpcResp.Data)
	if err != nil {
		return nil, err
	}
	return resp, err
}
