package loanlogic

import (
	"context"

	"gorm.io/gorm"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/lender/rpc/internal/logic"
	"zeebee-backend/app/lender/rpc/internal/svc"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
)

type UpdateLoanStageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	baseLogic *logic.BaseLogic
}

func NewUpdateLoanStageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLoanStageLogic {
	return &UpdateLoanStageLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *UpdateLoanStageLogic) UpdateLoanStage(in *lender.LoanUpdateStageRequest) (resp *lender.EmptyResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(
				logx.Field("stacks", errorx.GetErrorStack(err)),
				logx.Field("service_name", l.svcCtx.Config.Name),
			).Error(err)
		}
	}()
	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	// Get loan
	loan, err := l.svcCtx.LoanRepo.GetLenderLoan(l.ctx, in.Data.Id, claim.OrganizationId, false)
	if err != nil {
		return nil, errorx.Wrap(err, "Failed To Get Loan")
	}

	if loan.Source != constants.LoanSourceSystem && loan.Source != constants.LoanSourcePreapproval {
		return nil, errorx.WithStack(errorx.New("Invalid loan source"))
	}

	// Save original status for operation log
	originalStage := loan.Stage
	originalPendingStatus := loan.PendingStatus

	// Update loan stage
	err = l.svcCtx.LoanRepo.Transaction(l.ctx, func(tx *gorm.DB) error {
		// Update stage
		stage := constants.LoanStage(in.Data.Stage)

		err = l.svcCtx.LoanRepo.UpdateStageWithTx(l.ctx, tx, in.Data.Id, stage)
		if err != nil {
			return errorx.Wrap(err, "Failed To Update Loan")
		}

		// Create operation log
		operatorId := claim.OrganizationUserId

		fromStatus := dto.LoanOperationStatus{
			Stage:         originalStage,
			PendingStatus: originalPendingStatus,
		}

		toStatus := dto.LoanOperationStatus{
			Stage:         stage,
			PendingStatus: loan.PendingStatus,
		}

		// loanBaseLogic := NewLoanBaseLogic(l.ctx, l.svcCtx)
		return l.baseLogic.CreateLoanOperationLog(l.ctx, tx, loan.Id, operatorId, fromStatus, toStatus, constants.LoanOperationActionUpdate, nil)
	})

	if err != nil {
		return nil, err
	}

	return &lender.EmptyResponse{}, nil
}
