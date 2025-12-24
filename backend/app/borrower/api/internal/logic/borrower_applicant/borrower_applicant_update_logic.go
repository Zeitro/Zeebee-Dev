package borrower_applicant

import (
	"context"

	"github.com/jinzhu/copier"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/errorx"
)

type BorrowerApplicantUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// update
func NewBorrowerApplicantUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BorrowerApplicantUpdateLogic {
	return &BorrowerApplicantUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BorrowerApplicantUpdateLogic) BorrowerApplicantUpdate(req *types.BorrowerApplicantEditRequest) (resp *types.EmptyResponse, err error) {
	rpcReq := &borrower.BorrowerApplicantUpdateRequest{
		Data: &borrower.BorrowerApplicantEditEntity{},
	}

	err = copier.CopyWithOption(&rpcReq.Data, req, copier.Option{
		DeepCopy:      true,
		IgnoreEmpty:   true,
		CaseSensitive: true,
	})
	if err != nil {
		return nil, errorx.GetFailedTransReqErr(err)
	}

	_, err = l.svcCtx.BorrowerApplicantRpcClient.Update(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	return &types.EmptyResponse{}, nil
}
