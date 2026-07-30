package filelogic

import (
	"context"

	"gorm.io/gorm"

	"zeebee-backend/app/lender/rpc/internal/svc"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/proto"

	"zeebee-backend/app/lender/rpc/internal/logic"
	"zeebee-backend/common/errorx"
)

type DeleteFileLogic struct {
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	baseLogic *logic.BaseLogic
}

func NewDeleteFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteFileLogic {
	return &DeleteFileLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *DeleteFileLogic) DeleteFile(in *proto.PathIdRequest) (resp *proto.EmptyResponse, err error) {

	var loan *model.Loan
	var loanDocument *model.LoanDocument
	var file *model.File

	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	file, err = l.svcCtx.FileRepo.FindOneById(l.ctx, in.Data.Id)
	if err != nil {
		return nil, errorx.Wrap(err, "fail to find file")
	}

	if file.EntityType != constants.FileEntityTypeLoanDocument {
		return nil, errorx.WithStack(errorx.New("file is not associated with a loan document"))
	}

	loan, loanDocument, err = l.baseLogic.GetLoanAndDocumentByDocumentId(l.ctx, file.EntityId, claim.OrganizationId)
	if err != nil {
		return nil, err
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
			documentStatus := constants.LoanDocumentStatusPendingUpload
			if loanDocument.Category == constants.LoanDocumentCategoryBorrower || loanDocument.Category == constants.LoanDocumentCategoryExtend {
				documentStatus = constants.LoanDocumentStatusPendingRequested
			}

			innerErr = l.svcCtx.LoanDocumentRepo.UpdateColumnsByIdWithTx(l.ctx, tx, loanDocument.Id, map[string]interface{}{
				"status": documentStatus,
			})
			if innerErr != nil {
				return errorx.Wrap(innerErr, "fail to update loan document status")
			}
		}

		// create operation log
		operationFile := &dto.OperationFile{
			EntityId:   file.EntityId,
			EntityType: file.EntityType,
			FileUrl:    file.FileUrl,
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
			constants.LoanOperationActionDeleteDocumentFile,
			&dto.LoanOperationMetadata{
				LoanDocuments: []*dto.OperationLoanDocument{
					{
						Category: loanDocument.Category,
						FileType: loanDocument.FileType,
						Status:   loanDocument.Status,
						Files:    []*dto.OperationFile{operationFile},
					},
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
