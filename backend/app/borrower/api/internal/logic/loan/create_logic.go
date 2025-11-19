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

	err = copier.CopyWithOption(&rpcReq.Data, req, copier.Option{
		DeepCopy:      true,
		IgnoreEmpty:   false,
		CaseSensitive: true,
	})
	if err != nil {
		return nil, errorx.GetFailedTransReqErr(err)
	}
	// handle applicants
	rpcApplicants := slice.UniqueByComparator(rpcReq.Data.Applicants, func(item *borrower.ThirdPartyLoanApplicantEditEntry, other *borrower.ThirdPartyLoanApplicantEditEntry) bool {
		return item.BorrowerApplicantId == other.BorrowerApplicantId
	})
	if len(rpcApplicants) == 0 {
		err = errorx.WithStack(errorx.New("No Applicant"))
		return nil, err
	}
	rpcReq.Data.Applicants = rpcApplicants

	// get borrower id
	rpcReq.Data.BorrowerId, err = utils.GetAccountIdFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	rpcResp, err = l.svcCtx.LoanRpcClient.CreateBorrowerLoan(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}
	err = copier.Copy(resp, &rpcResp.Data)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	return resp, err
}
