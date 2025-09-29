package loanlogic

import (
	"context"

	"github.com/pkg/errors"
	"gorm.io/gorm"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/lender/rpc/internal/svc"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/constants"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/types"
)

type UpdateLoanStageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateLoanStageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLoanStageLogic {
	return &UpdateLoanStageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateLoanStageLogic) UpdateLoanStage(in *lender.LoanUpdateStageRequest) (resp *lender.EmptyResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()

	// Get loan
	loan, err := l.svcCtx.LoanRepo.GetLenderLoan(l.ctx, in.Data.Id, in.Data.LenderCompanyId, false)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get loan")
	}

	if loan.Source != constants.LoanSourceSystem {
		return nil, errors.WithStack(errors.New("invalid loan source"))
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
			return errors.Wrap(err, "failed to update loan")
		}

		// Create operation log
		operatorId := in.Data.LenderEmployeeId

		fromStatus := types.LoanOperationStatus{
			Stage:         originalStage,
			PendingStatus: originalPendingStatus,
		}

		toStatus := types.LoanOperationStatus{
			Stage:         stage,
			PendingStatus: loan.PendingStatus,
		}

		loanBaseLogic := NewLoanBaseLogic(l.ctx, l.svcCtx)
		return loanBaseLogic.CreateLoanOperationLog(tx, loan.Id, operatorId, fromStatus, toStatus, constants.LoanOperationActionUpdate, nil)
	})

	if err != nil {
		return nil, err
	}

	return &lender.EmptyResponse{}, nil
}
