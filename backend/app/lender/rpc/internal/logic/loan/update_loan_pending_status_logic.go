package loanlogic

import (
	"context"

	"github.com/pkg/errors"
	"gorm.io/gorm"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/lender/rpc/internal/svc"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/constants"
	"zeebee-backend/domains/zeebee/model"

	"zeebee-backend/common/errorx"
)

type UpdateLoanPendingStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateLoanPendingStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLoanPendingStatusLogic {
	return &UpdateLoanPendingStatusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateLoanPendingStatusLogic) UpdateLoanPendingStatus(in *lender.LoanUpdatePendingStatusRequest) (resp *lender.EmptyResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()

	// Get loan
	loan, err := l.svcCtx.LoanRepo.GetLenderLoan(l.ctx, in.Data.Id, in.Data.LenderCompanyId, false)
	if err != nil {
		return nil, errors.Wrap(err, "Failed To Get Loan")
	}

	// Validate loan is in pending stage
	if loan.Stage != constants.LoanStagePending {
		return nil, errors.New("Loan Is Not In Pending Stage")
	}
	if loan.Source != constants.LoanSourceSystem && loan.Source != constants.LoanSourcePreApproval {
		return nil, errors.WithStack(errors.New("Invalid Loan Source"))
	}
	// Save original status for operation log
	originalStage := loan.Stage
	originalPendingStatus := loan.PendingStatus

	// Update loan pending status
	err = l.svcCtx.LoanRepo.Transaction(l.ctx, func(tx *gorm.DB) error {
		// Update pending status
		pendingStatus := constants.LoanPendingStatus(in.Data.PendingStatus)

		// 根据状态规则进行验证
		switch pendingStatus {
		case constants.LoanPendingStatusQuoting,
			constants.LoanPendingStatusProcessing,
			constants.LoanPendingStatusUnderWriting:
			// 这些是中间状态，可以直接设置
		case constants.LoanPendingStatusClosing:
			// Closing 是最终状态，可以直接设置
		default:
			return errors.Errorf("Invalid Pending Status: %s", pendingStatus)
		}

		err = l.svcCtx.LoanRepo.UpdatePendingStatusWithTx(l.ctx, tx, in.Data.Id, pendingStatus, loan.Stage)
		if err != nil {
			return errors.Wrap(err, "Failed To Update Loan")
		}

		// Create operation log
		operatorId := in.Data.LenderEmployeeId

		fromStatus := model.LoanOperationStatus{
			Stage:         originalStage,
			PendingStatus: originalPendingStatus,
		}

		toStatus := model.LoanOperationStatus{
			Stage:         loan.Stage,
			PendingStatus: pendingStatus,
		}

		loanBaseLogic := NewLoanBaseLogic(l.ctx, l.svcCtx)
		return loanBaseLogic.CreateLoanOperationLog(tx, loan.Id, operatorId, fromStatus, toStatus, constants.LoanOperationActionUpdate, nil)
	})

	if err != nil {
		return nil, err
	}

	return &lender.EmptyResponse{}, nil
}
