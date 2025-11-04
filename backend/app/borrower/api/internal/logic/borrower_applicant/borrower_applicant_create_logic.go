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

type BorrowerApplicantCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// create
func NewBorrowerApplicantCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BorrowerApplicantCreateLogic {
	return &BorrowerApplicantCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BorrowerApplicantCreateLogic) BorrowerApplicantCreate(req *types.BorrowerApplicantEditRequest) (resp *types.BorrowerApplicantCreateResponse, err error) {
	rpcReq := &borrower.BorrowerApplicantCreateRequest{
		Data: &borrower.BorrowerApplicantEditEntry{},
	}

	err = copier.CopyWithOption(&rpcReq.Data, req, copier.Option{
		DeepCopy:      true,
		IgnoreEmpty:   true,
		CaseSensitive: true,
	})
	if err != nil {
		return nil, errorx.GetFailedTransReqErr(err)
	}

	// get borrower id
	borrowerId, err := utils.GetAccountIdFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	rpcReq.Data.BorrowerId = borrowerId

	rpcResp, err := l.svcCtx.BorrowerApplicantRpcClient.Create(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	resp = &types.BorrowerApplicantCreateResponse{}
	err = copier.CopyWithOption(&resp.BorrowerApplicantEntry, &rpcResp.Data, copier.Option{DeepCopy: true})
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}

	return resp, nil
}
