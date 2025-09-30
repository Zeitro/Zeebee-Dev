package loandocumentlogic

import (
	"context"

	"github.com/pkg/errors"
	"gorm.io/gorm"

	"zeebee-backend/app/lender/rpc/internal/repo/model"
	"zeebee-backend/app/lender/rpc/internal/svc"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/constants"
	"zeebee-backend/common/types"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/common/errorx"
)

type UpdateCommentLogic struct {
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	baseLogic *BaseLogic
	logx.Logger
}

func NewUpdateCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateCommentLogic {
	return &UpdateCommentLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		baseLogic: NewBaseLogic(ctx, svcCtx),
		Logger:    logx.WithContext(ctx),
	}
}

func (l *UpdateCommentLogic) UpdateComment(in *lender.LoanDocumentUpdateCommentRequest) (resp *lender.EmptyResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	var loan *model.Loan
	var loanDocument *model.LoanDocument

	loan, loanDocument, err = l.prepare(in)
	if err != nil {
		return nil, err
	}

	err = l.svcCtx.LoanDocumentRepo.Transaction(l.ctx, func(tx *gorm.DB) error {
		var innerErr error

		innerErr = l.svcCtx.LoanDocumentRepo.UpdateColumnsByIdWithTx(l.ctx, tx, loanDocument.Id, map[string]interface{}{
			"comment": in.Data.Comment,
		})
		if innerErr != nil {
			return errors.Wrap(innerErr, "failed to update loan document comment")
		}

		innerErr = l.baseLogic.CreateLoanOperationLog(
			tx,
			in.Data.Id,
			in.Data.LenderEmployeeId,
			types.LoanOperationStatus{
				Stage:         loan.Stage,
				PendingStatus: loan.PendingStatus,
			},
			types.LoanOperationStatus{
				Stage:         loan.Stage,
				PendingStatus: loan.PendingStatus,
			},
			constants.LoanOperationActionComment,
			&types.LoanOperationMetadata{
				Comment: in.Data.Comment,
			},
		)

		return innerErr
	})

	if err != nil {
		return nil, err
	}

	return &lender.EmptyResponse{}, nil
}

func (l *UpdateCommentLogic) prepare(in *lender.LoanDocumentUpdateCommentRequest) (loan *model.Loan, loanDocument *model.LoanDocument, err error) {
	return l.baseLogic.GetLoanAndDocumentByDocumentId(in.Data.Id, in.Data.LenderCompanyId)
}
