package loanlogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/proto"

	"zeebee-backend/common/errorx"
)

type UpdateBorrowerLoanPendingStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewUpdateBorrowerLoanPendingStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateBorrowerLoanPendingStatusLogic {
	return &UpdateBorrowerLoanPendingStatusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *UpdateBorrowerLoanPendingStatusLogic) UpdateBorrowerLoanPendingStatus(in *borrower.LoanUpdatePendingStatusRequest) (resp *proto.EmptyResponse, err error) {

	claim, err := l.baseLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	// todo check organization
	loanModel, err := l.svcCtx.LoanRepo.FindOneByColumns(l.ctx, map[string]interface{}{
		"id":                       in.Data.Id,
		"borrower_organization_id": claim.OrganizationId,
	})
	if err != nil {
		return nil, errorx.Wrap(err, "fail to get loan")
	}
	if loanModel.Stage != constants.LoanStagePending {
		return nil, errorx.New("loan's stage should be pending")
	}

	if loanModel.Source != constants.LoanSourceBorrower {
		return nil, errorx.WithStack(errorx.New("loan source is not 3rd party"))
	}
	err = l.svcCtx.LoanRepo.UpdateColumnsById(l.ctx, in.Data.Id, map[string]interface{}{"pending_status": in.Data.PendingStatus})

	return &proto.EmptyResponse{}, err
}
