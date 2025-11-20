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
	// get borrowerId
	borrowerId, err := utils.GetAccountIdFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	rpcResp, err := l.svcCtx.BorrowerApplicantRpcClient.All(l.ctx, &borrower.BorrowerApplicantAllRequest{
		Data: &borrower.BorrowerApplicantAllRequestData{
			BorrowerId: borrowerId,
		},
	})
	if err != nil {
		return nil, err
	}

	resp = &types.BorrowerApplicantAllResponse{}
	resp.Lists = make([]*types.BorrowerApplicantPageEntry, 0, len(rpcResp.Data))

	err = copier.CopyWithOption(&resp.Lists, rpcResp.Data, copier.Option{
		DeepCopy:      true,
		CaseSensitive: true,
	})
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}

	return resp, nil
}
