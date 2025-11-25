package borrower_applicant

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

type BorrowerApplicantPageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// page
func NewBorrowerApplicantPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BorrowerApplicantPageLogic {
	return &BorrowerApplicantPageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BorrowerApplicantPageLogic) BorrowerApplicantPage(req *types.BorrowerApplicantPageRequest) (resp *types.BorrowerApplicantPageResponse, err error) {
	rpcReq := &borrower.BorrowerApplicantPageRequest{
		Data: &borrower.BorrowerApplicantPageRequestData{
			Filter: &borrower.BorrowerApplicantPageRequestFilter{},
			Sorter: &borrower.BorrowerApplicantPageRequestSorter{},
		},
	}
	resp = &types.BorrowerApplicantPageResponse{}
	rpcResp := &borrower.BorrowerApplicantPageResponse{
		Data: &borrower.BorrowerApplicantPageResponseData{},
	}

	err = copier.CopyWithOption(&rpcReq.Data, req, copier.Option{DeepCopy: true})
	if err != nil {
		return nil, errorx.GetFailedTransReqErr(err)
	}
	// get borrower id
	claim, err := utils.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	rpcReq.Data.Filter.UserId = claim.UserId
	rpcReq.Data.Filter.OrganizationId = claim.OrganizationId
	if err != nil {
		return nil, err
	}
	rpcResp, err = l.svcCtx.BorrowerApplicantRpcClient.Page(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}
	err = copier.Copy(resp, &rpcResp.Data)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	return resp, err
}
