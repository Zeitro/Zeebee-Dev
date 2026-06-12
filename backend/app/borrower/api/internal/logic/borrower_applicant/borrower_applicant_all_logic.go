package borrower_applicant

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
)

type BorrowerApplicantAllLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// all
func NewBorrowerApplicantAllLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BorrowerApplicantAllLogic {
	return &BorrowerApplicantAllLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BorrowerApplicantAllLogic) BorrowerApplicantAll(req *types.EmptyRequest) (resp *types.BorrowerApplicantAllResponse, err error) {
	rpcResp, err := l.svcCtx.BorrowerApplicantRpcClient.All(l.ctx, &borrower.BorrowerApplicantAllRequest{
		Data: &borrower.BorrowerApplicantAllRequestData{},
	})
	if err != nil {
		return nil, err
	}

	resp = &types.BorrowerApplicantAllResponse{}
	resp.Lists = make([]*types.BorrowerApplicantPageEntity, 0, len(rpcResp.Data))
	listPtr, err := copierx.Trans[[]*borrower.BorrowerApplicantPageEntity, []*types.BorrowerApplicantPageEntity](&rpcResp.Data)
	if err != nil {
		return nil, err
	}
	resp.Lists = *listPtr

	return resp, nil
}
