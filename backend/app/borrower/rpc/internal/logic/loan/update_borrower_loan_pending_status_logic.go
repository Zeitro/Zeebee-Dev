package loanlogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/constants"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/common/errorx"
)

type UpdateBorrowerLoanPendingStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	baseLogic *logic.BaseLogic
}

func NewUpdateBorrowerLoanPendingStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateBorrowerLoanPendingStatusLogic {
	return &UpdateBorrowerLoanPendingStatusLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		baseLogic: logic.NewBaseLogic(ctx, svcCtx),
	}
}

func (l *UpdateBorrowerLoanPendingStatusLogic) UpdateBorrowerLoanPendingStatus(in *borrower.LoanUpdatePendingStatusRequest) (resp *borrower.EmptyResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()

	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	// todo check organization
	loanModel, err := l.svcCtx.LoanRepo.FindOneByColumns(l.ctx, map[string]interface{}{
		"id":                       in.Data.Id,
		"borrower_organization_id": claim.OrganizationId,
	})
	if err != nil {
		return nil, errorx.Wrap(err, "Failed To Get Loan")
	}
	if loanModel.Stage != constants.LoanStagePending {
		return nil, errorx.New("Loan's Stage Should Be Pending")
	}

	if loanModel.Source != constants.LoanSourceBorrower {
		return nil, errorx.WithStack(errorx.New("Loan Source Is Not 3rd Party"))
	}
	err = l.svcCtx.LoanRepo.UpdateColumnsById(l.ctx, in.Data.Id, map[string]interface{}{"pending_status": in.Data.PendingStatus})

	return &borrower.EmptyResponse{}, err
}
