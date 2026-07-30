package loandocumentlogic

import (
	"context"

	"gorm.io/gorm"

	"zeebee-backend/app/lender/rpc/internal/logic"
	"zeebee-backend/app/lender/rpc/internal/svc"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/proto"

	"zeebee-backend/common/errorx"
)

type UpdateCommentLogic struct {
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	baseLogic *logic.BaseLogic
}

func NewUpdateCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateCommentLogic {
	return &UpdateCommentLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *UpdateCommentLogic) UpdateComment(in *lender.LoanDocumentUpdateCommentRequest) (resp *proto.EmptyResponse, err error) {

	var loan *model.Loan
	var loanDocument *model.LoanDocument

	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	loan, loanDocument, err = l.prepare(in, claim.OrganizationId)
	if err != nil {
		return nil, err
	}

	err = l.svcCtx.LoanDocumentRepo.Transaction(l.ctx, func(tx *gorm.DB) error {
		var innerErr error

		innerErr = l.svcCtx.LoanDocumentRepo.UpdateColumnsByIdWithTx(l.ctx, tx, loanDocument.Id, map[string]interface{}{
			"comment": in.Data.Comment,
		})
		if innerErr != nil {
			return errorx.Wrap(innerErr, "fail to update loan document comment")
		}

		innerErr = l.baseLogic.CreateLoanOperationLog(
			l.ctx,
			tx,
			in.Data.Id,
			claim.OrganizationUserId,
			dto.LoanOperationStatus{
				Stage:         loan.Stage,
				PendingStatus: loan.PendingStatus,
			},
			dto.LoanOperationStatus{
				Stage:         loan.Stage,
				PendingStatus: loan.PendingStatus,
			},
			constants.LoanOperationActionComment,
			&dto.LoanOperationMetadata{
				Comment: in.Data.Comment,
			},
		)

		return innerErr
	})

	if err != nil {
		return nil, err
	}

	return &proto.EmptyResponse{}, nil
}

func (l *UpdateCommentLogic) prepare(in *lender.LoanDocumentUpdateCommentRequest, organizationId int64) (loan *model.Loan, loanDocument *model.LoanDocument, err error) {
	return l.baseLogic.GetLoanAndDocumentByDocumentId(l.ctx, in.Data.Id, organizationId)
}
