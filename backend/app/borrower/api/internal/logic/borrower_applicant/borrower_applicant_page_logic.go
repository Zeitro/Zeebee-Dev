package borrower_applicant

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/proto"
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
	rpcReq := &proto.CommonPageRequest{
		Data: &proto.CommonPageRequestData{
			Filter: &proto.CommonPageFilter{},
			Sorter: &proto.CommonPageSorter{},
		},
	}
	resp = &types.BorrowerApplicantPageResponse{}
	rpcResp := &borrower.BorrowerApplicantPageResponse{
		Data: &borrower.BorrowerApplicantPageResponseData{},
	}
	rpcReq.Data, err = copierx.Trans[types.BorrowerApplicantPageRequest, proto.CommonPageRequestData](req)
	if err != nil {
		return nil, err
	}
	rpcResp, err = l.svcCtx.BorrowerApplicantRpcClient.Page(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}
	resp, err = copierx.Trans[borrower.BorrowerApplicantPageResponseData, types.BorrowerApplicantPageResponse](rpcResp.Data)
	if err != nil {
		return nil, err
	}
	return resp, err
}
