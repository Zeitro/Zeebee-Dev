package filelogic

import (
	"context"

	"gorm.io/gorm"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/proto"
)

type DeleteFileLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewDeleteFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteFileLogic {
	return &DeleteFileLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *DeleteFileLogic) DeleteFile(in *proto.PathIdRequest) (resp *proto.EmptyResponse, err error) {

	var loan *model.Loan
	var loanDocument *model.LoanDocument
	var file *model.File
	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	file, err = l.svcCtx.FileRepo.FindOneById(l.ctx, in.Data.Id)
	if err != nil {
		return nil, errorx.Wrap(err, "fail to Find File")
	}

	if file.EntityType != constants.FileEntityTypeLoanDocument {
		return nil, errorx.WithStack(errorx.New("file is not associated with a loan document"))
	}

	loanDocument, err = l.svcCtx.LoanDocumentRepo.FindOneWithPreloadById(l.ctx, file.EntityId, nil)
	if err != nil {
		return nil, errorx.Wrap(err, "fail to find loan document")
	}

	loan, err = l.svcCtx.LoanRepo.GetBorrowerLoan(l.ctx, loanDocument.LoanId, claim.OrganizationId)
	if err != nil {
		return nil, errorx.Wrap(err, "fail to get loan")
	}
	if loan.BorrowerOrganizationId != claim.OrganizationId {
		return nil, errorx.WithStack(errorx.New("loan does not belong to borrower"))
	}

	err = l.svcCtx.FileRepo.Transaction(l.ctx, func(tx *gorm.DB) error {
		var innerErr error

		// delete file
		_, innerErr = l.svcCtx.FileRepo.DeleteByIDWithTx(l.ctx, tx, in.Data.Id)
		if innerErr != nil {
			return errorx.Wrap(innerErr, "fail to delete file")
		}

		// get file count
		var fileCount int64
		fileCount, innerErr = l.svcCtx.FileRepo.GetFileCountByEntityIdWithTx(l.ctx, tx, loanDocument.Id, constants.FileEntityTypeLoanDocument)
		if innerErr != nil {
			return errorx.Wrap(innerErr, "fail to check remaining files")
		}

		// update loan document status
		if fileCount == 0 {
			loanDocumentStatus := constants.LoanDocumentStatusPendingUpload
			if loanDocument.Category == constants.LoanDocumentCategoryBorrower || loanDocument.Category == constants.LoanDocumentCategoryExtend {
				loanDocumentStatus = constants.LoanDocumentStatusPendingRequested
			}

			innerErr = l.svcCtx.LoanDocumentRepo.UpdateColumnsByIdWithTx(l.ctx, tx, loanDocument.Id, map[string]interface{}{
				"status": loanDocumentStatus,
			})
			if innerErr != nil {
				return errorx.Wrap(innerErr, "fail to update loan document status")
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &proto.EmptyResponse{}, nil
}
