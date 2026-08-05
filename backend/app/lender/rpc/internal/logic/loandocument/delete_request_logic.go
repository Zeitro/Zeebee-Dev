package loandocumentlogic

import (
	"context"

	"github.com/duke-git/lancet/v2/slice"
	"gorm.io/gorm"

	"zeebee-backend/app/lender/rpc/internal/logic"
	"zeebee-backend/app/lender/rpc/internal/svc"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/proto"

	"zeebee-backend/common/errorx"
)

type DeleteRequestLogic struct {
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	baseLogic *logic.BaseLogic
}

func NewDeleteRequestLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteRequestLogic {
	return &DeleteRequestLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *DeleteRequestLogic) DeleteRequest(in *proto.PathIdRequest) (resp *proto.EmptyResponse, err error) {

	var loan *model.Loan
	var loanDocument *model.LoanDocument

	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx, constants.UserTypeLender)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	// validate
	loan, loanDocument, err = l.prepare(in, claim.OrganizationId)
	if err != nil {
		return nil, err
	}

	err = l.svcCtx.LoanDocumentRepo.Transaction(l.ctx, func(tx *gorm.DB) error {
		var innerErr error

		// delete files
		fileIds := slice.Map[*model.File, int64](loanDocument.Files, func(_ int, file *model.File) int64 {
			return file.Id
		})
		if len(fileIds) > 0 {
			_, innerErr = l.svcCtx.FileRepo.DeleteByIDsWithTx(l.ctx, tx, fileIds)
			if innerErr != nil {
				return innerErr
			}
		}

		switch loanDocument.Category {
		case constants.LoanDocumentCategoryExtend:
			_, innerErr = l.svcCtx.LoanDocumentRepo.DeleteByIDWithTx(l.ctx, tx, loanDocument.Id)
			if innerErr != nil {
				return errorx.Wrap(innerErr, "fail to delete loan document")
			}
		case constants.LoanDocumentCategoryBorrower:
			innerErr = l.svcCtx.LoanDocumentRepo.UpdateColumnsByIdWithTx(l.ctx, tx, loanDocument.Id, map[string]interface{}{
				"status": constants.LoanDocumentStatusPendingUpload,
			})
			if innerErr != nil {
				return errorx.Wrap(innerErr, "fail to update loan document status")
			}
		default:
			return errorx.WithStack(errorx.New("invalid loan document category"))
		}
		operationLoanDocument := &dto.OperationLoanDocument{}
		operationLoanDocument, err = copierx.Trans[model.LoanDocument, dto.OperationLoanDocument](loanDocument)
		if err != nil {
			return errorx.Wrap(err, "fail to copy loan document")
		}

		innerErr = l.baseLogic.CreateLoanOperationLog(
			l.ctx,
			tx,
			loanDocument.LoanId,
			claim.OrganizationUserId,
			dto.LoanOperationStatus{
				Stage:         loan.Stage,
				PendingStatus: loan.PendingStatus,
			},
			dto.LoanOperationStatus{
				Stage:         loan.Stage,
				PendingStatus: loan.PendingStatus,
			},
			constants.LoanOperationActionDeleteDocumentRequest,
			&dto.LoanOperationMetadata{
				LoanDocuments: []*dto.OperationLoanDocument{
					operationLoanDocument,
				},
			},
		)
		if innerErr != nil {
			return errorx.Wrap(innerErr, "fail to create operation log")
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &proto.EmptyResponse{}, nil
}

func (l *DeleteRequestLogic) prepare(in *proto.PathIdRequest, organizationId int64) (loan *model.Loan, loanDocument *model.LoanDocument, err error) {
	return l.baseLogic.GetLoanAndDocumentByDocumentId(l.ctx, in.Data.Id, organizationId)
}
