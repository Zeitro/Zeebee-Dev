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
		Data: &borrower.BorrowerApplicantEditEntity{},
	}
	rpcReq.Data, err = copierx.Trans[types.BorrowerApplicantEditRequest, borrower.BorrowerApplicantEditEntity](req)
	if err != nil {
		return nil, err
	}

	rpcResp, err := l.svcCtx.BorrowerApplicantRpcClient.Create(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	resp = &types.BorrowerApplicantCreateResponse{}
	entityPtr, err := copierx.Trans[proto.BorrowerApplicantEntity, types.BorrowerApplicantEntity](rpcResp.Data)
	if err != nil {
		return nil, err
	}
	resp.BorrowerApplicantEntity = *entityPtr

	return resp, nil
}
