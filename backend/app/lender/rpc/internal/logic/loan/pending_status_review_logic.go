package loanlogic

import (
	"context"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"

	"zeebee-backend/app/lender/rpc/internal/svc"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/constants"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/types"
)

type PendingStatusReviewLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPendingStatusReviewLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PendingStatusReviewLogic {
	return &PendingStatusReviewLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PendingStatusReviewLogic) PendingStatusReview(in *lender.LoanPendingStatusReviewRequest) (resp *lender.EmptyResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()

	loan, err := l.svcCtx.LoanRepo.GetLenderLoan(l.ctx, in.Data.Id, in.Data.LenderCompanyId, false)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	currentStatus := constants.LoanPendingStatus(in.Data.CurrentPendingStatus)
	action := constants.LoanOperationAction(in.Data.Action)

	originalStage := loan.Stage
	originalPendingStatus := loan.PendingStatus
	targetStage := constants.LoanStagePending

	// validate loan pending status
	if loan.Stage != constants.LoanStagePending {
		return nil, errors.WithStack(errors.New("loan is not in pending stage"))
	}
	if loan.Source != constants.LoanSourceSystem {
		return nil, errors.WithStack(errors.New("invalid loan source"))
	}

	if loan.PendingStatus != currentStatus {
		return nil, errors.WithStack(errors.New("current pending status is not valid"))
	}
	if action == constants.LoanOperationActionReject && in.Data.Comment == "" {
		return nil, errors.WithStack(errors.New("Comment is empty"))
	}
	var targetStatus constants.LoanPendingStatus
	loanBaseLogic := NewLoanBaseLogic(l.ctx, l.svcCtx)
	targetStatus, err = loanBaseLogic.GetNextPendingStatus(currentStatus, action)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if targetStatus == constants.LoanPendingStatusNone && currentStatus == constants.LoanPendingStatusClosing {
		targetStage = constants.LoanStageActive
	}

	err = l.svcCtx.LoanRepo.Transaction(l.ctx, func(tx *gorm.DB) error {
		var innerErr error

		// update loan pending status
		innerErr = l.svcCtx.LoanRepo.UpdatePendingStatusWithTx(l.ctx, tx, in.Data.Id, targetStatus, targetStage)
		if innerErr != nil {
			return errors.WithStack(innerErr)
		}

		// add operation
		operatorId := in.Data.LenderEmployeeId
		fromStatus := types.LoanOperationStatus{
			Stage:         originalStage,
			PendingStatus: originalPendingStatus,
		}
		toStatus := types.LoanOperationStatus{
			Stage:         targetStage,
			PendingStatus: targetStatus,
		}
		metadata := &types.LoanOperationMetadata{
			Comment: in.Data.Comment,
		}

		return loanBaseLogic.CreateLoanOperationLog(tx, loan.Id, operatorId, fromStatus, toStatus, action, metadata)
	})

	if err != nil {
		return nil, err
	}

	return &lender.EmptyResponse{}, nil
}
