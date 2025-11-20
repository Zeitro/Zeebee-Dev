package borrower_applicant

import (
	"context"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/utils"

	"github.com/zeromicro/go-zero/core/logx"
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

	// get borrower id
	borrowerId, err := utils.GetAccountIdFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	rpcReq.Data.BorrowerId = borrowerId

	_, err = l.svcCtx.BorrowerApplicantRpcClient.Delete(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	resp = &types.EmptyResponse{}
	return resp, nil
}
