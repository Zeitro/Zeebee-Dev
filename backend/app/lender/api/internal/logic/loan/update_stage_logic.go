package loan

import (
	"context"

	"github.com/jinzhu/copier"
	"github.com/pkg/errors"

	"zeebee-backend/app/lender/api/internal/svc"
	"zeebee-backend/app/lender/api/internal/types"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateStageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// update loan stage
func NewUpdateStageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateStageLogic {
	return &UpdateStageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateStageLogic) UpdateStage(req *types.UpdateLoanStageRequest) (resp *types.EmptyResponse, err error) {
	rpcReq := &lender.LoanUpdateStageRequest{
		Data: &lender.LoanUpdateStageRequestData{},
	}
	resp = &types.EmptyResponse{}

	err = copier.CopyWithOption(&rpcReq.Data, req, copier.Option{
		IgnoreEmpty:   false,
		DeepCopy:      true,
		CaseSensitive: true,
	})
	if err != nil {
		return nil, errorx.GetFailedTransReqErr(err)
	}

	// get lender id
	rpcReq.Data.LenderCompanyId, err = utils.GetLenderCompanyIdFromCtx(l.ctx)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// get employee id
	rpcReq.Data.LenderEmployeeId, err = utils.GetLenderAccountIdFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	_, err = l.svcCtx.LoanRpcClient.UpdateLoanStage(l.ctx, rpcReq)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return resp, nil
}
