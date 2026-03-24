package borrower_applicant

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
)

type BorrowerApplicantCheckExistedLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// existed
func NewBorrowerApplicantCheckExistedLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BorrowerApplicantCheckExistedLogic {
	return &BorrowerApplicantCheckExistedLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BorrowerApplicantCheckExistedLogic) BorrowerApplicantCheckExisted(req *types.BorrowerApplicantCkeckExistedRequest) (resp *types.BorrowerApplicantCkeckExistedResponse, err error) {
	rpcReq := &borrower.BorrowerApplicantCheckExistedRequest{
		Data: &borrower.BorrowerApplicantCheckExistedRequestData{
			FirstName: req.FirstName,
			LastName:  req.LastName,
			MidName:   req.MidName,
			Id:        req.Id,
		},
	}

	rpcResp, err := l.svcCtx.BorrowerApplicantRpcClient.CheckExisted(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	resp = &types.BorrowerApplicantCkeckExistedResponse{}
	resp, err = copierx.Trans[borrower.BorrowerApplicantCheckExistedResponseData, types.BorrowerApplicantCkeckExistedResponse](rpcResp.Data)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
