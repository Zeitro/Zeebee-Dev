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

type CreateDraftLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// create draft
func NewDraftLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateDraftLogic {
	return &CreateDraftLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateDraftLogic) CreateDraft(req *types.LoanApplyRequest) (resp *types.LoanApplyResponse, err error) {
	rpcReq := &borrower.LoanApplyRequest{
		Data: &borrower.LoanApplyRequestData{},
	}
	resp = &types.LoanApplyResponse{}
	rpcResp := &borrower.LoanApplyResponse{}

	err = copier.CopyWithOption(&rpcReq.Data, req, copier.Option{
		DeepCopy:      true,
		IgnoreEmpty:   false,
		CaseSensitive: true,
		Converters:    []copier.TypeConverter{},
	})
	if err != nil {
		return nil, errorx.GetFailedTransReqErr(err)
	}

	// handle applicants
	rpcApplicants := slice.UniqueByComparator(rpcReq.Data.Applicants, func(item *borrower.LoanApplicantApplyEntry, other *borrower.LoanApplicantApplyEntry) bool {
		return item.BorrowerApplicantId == other.BorrowerApplicantId
	})
	rpcReq.Data.Applicants = rpcApplicants

	// get borrower id
	claim, err := utils.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	rpcReq.Data.UserId = claim.UserId
	rpcReq.Data.OrganizationId = claim.OrganizationId

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
