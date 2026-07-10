package loanlogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/proto"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/common/errorx"
)

type UpdateBorrowerLoanStageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewUpdateBorrowerLoanStageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateBorrowerLoanStageLogic {
	return &UpdateBorrowerLoanStageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *UpdateBorrowerLoanStageLogic) UpdateBorrowerLoanStage(in *borrower.LoanUpdateStageRequest) (resp *proto.EmptyResponse, err error) {
	defer func() {
		if err != nil {
			l.svcCtx.Logger.WithContext(l.ctx).WithFields(
				logx.Field("stacks", errorx.GetErrorStack(err)),
			).Error(err)
		}
	}()

	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	loanModel, err := l.svcCtx.LoanRepo.FindOneByColumns(l.ctx, map[string]interface{}{
		"id":                       in.Data.Id,
		"borrower_organization_id": claim.OrganizationId,
	})
	if err != nil {
		return nil, errorx.Wrap(err, "failed To Get Loan")
	}

	if loanModel.Source != constants.LoanSourceBorrower {
		return nil, errorx.WithStack(errorx.New("loan source is not 3rd party"))
	}

	err = l.svcCtx.LoanRepo.UpdateColumnsById(l.ctx, in.Data.Id, map[string]interface{}{"stage": in.Data.Stage})

	return &proto.EmptyResponse{}, err
}
