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
	"zeebee-backend/domains/zeebee/data"
)

type UpdateLoanPendingStatusLogic struct {
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	baseLogic *logic.BaseLogic
	logx.Logger
}

func NewUpdateLoanPendingStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLoanPendingStatusLogic {
	return &UpdateLoanPendingStatusLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		baseLogic: logic.NewBaseLogic(ctx, svcCtx),
		Logger:    logx.WithContext(ctx),
	}
}

func (l *UpdateLoanPendingStatusLogic) UpdateLoanPendingStatus(in *lender.LoanUpdatePendingStatusRequest) (resp *lender.EmptyResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()

	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	// Get loan
	loan, err := l.svcCtx.LoanRepo.GetLenderLoan(l.ctx, in.Data.Id, claim.OrganizationId, false)
	if err != nil {
		return nil, errorx.Wrap(err, "Failed To Get Loan")
	}

	// Validate loan is in pending stage
	if loan.Stage != constants.LoanStagePending {
		return nil, errorx.New("Loan is not in pending stage")
	}
	if loan.Source != constants.LoanSourceSystem && loan.Source != constants.LoanSourcePreapproval {
		return nil, errorx.WithStack(errorx.New("Invalid loan source"))
	}
	// Save original status for operation log
	originalStage := loan.Stage
	originalPendingStatus := loan.PendingStatus

	// Update loan pending status
	err = l.svcCtx.LoanRepo.Transaction(l.ctx, func(tx *gorm.DB) error {
		// Update pending status
		pendingStatus := constants.LoanPendingStatus(in.Data.PendingStatus)

		// validate pending status
		switch pendingStatus {
		case constants.LoanPendingStatusQuoting,
			constants.LoanPendingStatusProcessing,
			constants.LoanPendingStatusUnderWriting:
			// these are intermediate statuses, can be set directly
		case constants.LoanPendingStatusClosing:
			// Closing is the final status, can be set directly
		default:
			return errorx.Errorf("Invalid Pending Status: %s", pendingStatus)
		}

		err = l.svcCtx.LoanRepo.UpdatePendingStatusWithTx(l.ctx, tx, in.Data.Id, pendingStatus, loan.Stage)
		if err != nil {
			return errorx.Wrap(err, "Failed To Update Loan")
		}

		// Create operation log
		operatorId := claim.OrganizationUserId

		fromStatus := data.LoanOperationStatus{
			Stage:         originalStage,
			PendingStatus: originalPendingStatus,
		}

		toStatus := data.LoanOperationStatus{
			Stage:         loan.Stage,
			PendingStatus: pendingStatus,
		}

		// loanBaseLogic := NewLoanBaseLogic(l.ctx, l.svcCtx)
		return l.baseLogic.CreateLoanOperationLog(tx, loan.Id, operatorId, fromStatus, toStatus, constants.LoanOperationActionUpdate, nil)
	})

	if err != nil {
		return nil, err
	}

	return &lender.EmptyResponse{}, nil
}
