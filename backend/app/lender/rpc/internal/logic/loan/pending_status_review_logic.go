package loanlogic

import (
	"context"

	"gorm.io/gorm"

	"zeebee-backend/app/lender/rpc/internal/logic"
	"zeebee-backend/app/lender/rpc/internal/svc"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/proto"
)

type PendingStatusReviewLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewPendingStatusReviewLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PendingStatusReviewLogic {
	return &PendingStatusReviewLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *PendingStatusReviewLogic) PendingStatusReview(in *lender.LoanPendingStatusReviewRequest) (resp *proto.EmptyResponse, err error) {

	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx, constants.UserTypeLender)
	if err != nil {
		return nil, err
	}

	loan, err := l.svcCtx.LoanRepo.GetLenderLoan(l.ctx, in.Data.Id, claim.OrganizationId, false)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	currentStatus := constants.LoanPendingStatus(in.Data.CurrentPendingStatus)
	action := constants.LoanOperationAction(in.Data.Action)

	originalStage := loan.Stage
	originalPendingStatus := loan.PendingStatus
	targetStage := constants.LoanStagePending

	// validate loan pending status
	if loan.Stage != constants.LoanStagePending {
		return nil, errorx.WithStack(errorx.New("loan is not in pending stage"))
	}
	if loan.Source != constants.LoanSourceSystem && loan.Source != constants.LoanSourcePreapproval {
		return nil, errorx.WithStack(errorx.New("invalid loan source"))
	}

	if loan.PendingStatus != currentStatus {
		return nil, errorx.WithStack(errorx.New("current pending status is not valid"))
	}
	if action == constants.LoanOperationActionReject && in.Data.Comment == "" {
		return nil, errorx.WithStack(errorx.New("comment is empty"))
	}
	var targetStatus constants.LoanPendingStatus
	loanBaseLogic := NewLoanBaseLogic(l.ctx, l.svcCtx)
	targetStatus, err = loanBaseLogic.GetNextPendingStatus(currentStatus, action)
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	if targetStatus == constants.LoanPendingStatusNone && currentStatus == constants.LoanPendingStatusClosing {
		targetStage = constants.LoanStageActive
	}

	err = l.svcCtx.LoanRepo.Transaction(l.ctx, func(tx *gorm.DB) error {
		var innerErr error

		// update loan pending status
		innerErr = l.svcCtx.LoanRepo.UpdatePendingStatusWithTx(l.ctx, tx, in.Data.Id, targetStatus, targetStage)
		if innerErr != nil {
			return errorx.WithStack(innerErr)
		}

		// add operation
		operatorId := claim.OrganizationUserId
		fromStatus := dto.LoanOperationStatus{
			Stage:         originalStage,
			PendingStatus: originalPendingStatus,
		}
		toStatus := dto.LoanOperationStatus{
			Stage:         targetStage,
			PendingStatus: targetStatus,
		}
		metadata := &dto.LoanOperationMetadata{
			Comment: in.Data.Comment,
		}

		return l.baseLogic.CreateLoanOperationLog(l.ctx, tx, loan.Id, operatorId, fromStatus, toStatus, action, metadata)
	})

	if err != nil {
		return nil, err
	}

	return &proto.EmptyResponse{}, nil
}
