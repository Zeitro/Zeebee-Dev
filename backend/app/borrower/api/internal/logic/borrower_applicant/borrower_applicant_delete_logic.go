package borrower_applicant

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
)

type BorrowerApplicantDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// delete
func NewBorrowerApplicantDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BorrowerApplicantDeleteLogic {
	return &BorrowerApplicantDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BorrowerApplicantDeleteLogic) BorrowerApplicantDelete(req *types.PathIdRequest) (resp *types.EmptyResponse, err error) {
	rpcReq := &borrower.BorrowerApplicantDeleteRequest{
		Data: &borrower.BorrowerApplicantDeleteRequestData{
			Id: req.Id,
		},
	}

	_, err = l.svcCtx.BorrowerApplicantRpcClient.Delete(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	resp = &types.EmptyResponse{}
	return resp, nil
}
