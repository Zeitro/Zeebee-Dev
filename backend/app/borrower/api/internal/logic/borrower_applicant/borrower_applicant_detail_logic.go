package borrower_applicant

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
)

type BorrowerApplicantDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// detail
func NewBorrowerApplicantDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BorrowerApplicantDetailLogic {
	return &BorrowerApplicantDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BorrowerApplicantDetailLogic) BorrowerApplicantDetail(req *types.PathIdRequest) (resp *types.BorrowerApplicantDetailResponse, err error) {
	rpcReq := &borrower.BorrowerApplicantDetailRequest{
		Data: &borrower.BorrowerApplicantDetailRequestData{
			Id: req.Id,
		},
	}

	rpcResp, err := l.svcCtx.BorrowerApplicantRpcClient.Detail(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	resp = &types.BorrowerApplicantDetailResponse{}
	entityPtr, err := copierx.Trans[borrower.BorrowerApplicantEntity, types.BorrowerApplicantEntity](rpcResp.Data)
	if err != nil {
		return nil, err
	}
	resp.BorrowerApplicantEntity = *entityPtr
	return resp, nil
}
