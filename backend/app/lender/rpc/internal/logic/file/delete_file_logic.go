package filelogic

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

	"zeebee-backend/app/lender/rpc/internal/logic"
	"zeebee-backend/common/errorx"
)

type DeleteFileLogic struct {
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	baseLogic *logic.BaseLogic
	logx.Logger
}

func NewDeleteFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteFileLogic {
	return &DeleteFileLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		baseLogic: logic.NewBaseLogic(ctx, svcCtx),
		Logger:    logx.WithContext(ctx),
	}
}

func (l *DeleteFileLogic) DeleteFile(in *lender.FileDeleteRequest) (resp *lender.EmptyResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()

	var loan *model.Loan
	var loanDocument *model.LoanDocument
	var file *model.File

	file, err = l.svcCtx.FileRepo.FindOneById(l.ctx, in.Data.Id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find file")
	}

	if file.EntityType != constants.FileEntityTypeLoanDocument {
		return nil, errors.WithStack(errors.New("file is not associated with a loan document"))
	}

	loan, loanDocument, err = l.baseLogic.GetLoanAndDocumentByDocumentId(file.EntityId, in.Data.LenderCompanyId)
	if err != nil {
		return nil, err
	}

	err = l.svcCtx.FileRepo.Transaction(l.ctx, func(tx *gorm.DB) error {
		var innerErr error

		// delete file
		_, innerErr = l.svcCtx.FileRepo.DeleteByIDWithTx(l.ctx, tx, in.Data.Id)
		if innerErr != nil {
			return errors.Wrap(innerErr, "failed to delete file")
		}

		// get file count
		var fileCount int64
		fileCount, innerErr = l.svcCtx.FileRepo.GetFileCountByEntityIdWithTx(l.ctx, tx, loanDocument.Id, constants.FileEntityTypeLoanDocument)
		if innerErr != nil {
			return errors.Wrap(innerErr, "failed to check remaining files")
		}

		// update loan document status
		if fileCount == 0 {
			innerErr = l.svcCtx.LoanDocumentRepo.UpdateColumnsByIdWithTx(l.ctx, tx, loanDocument.Id, map[string]interface{}{
				"status": constants.LoanDocumentStatusPendingUpload,
			})
			if innerErr != nil {
				return errors.Wrap(innerErr, "failed to update loan document status")
			}
		}

		// create operation log
		operationFile := &types.OperationFile{
			EntityId:   file.EntityId,
			EntityType: file.EntityType,
			FileUrl:    file.FileUrl,
		}

		innerErr = l.baseLogic.CreateLoanOperationLog(
			tx,
			loanDocument.LoanId,
			in.Data.LenderEmployeeId,
			types.LoanOperationStatus{
				Stage:         loan.Stage,
				PendingStatus: loan.PendingStatus,
			},
			types.LoanOperationStatus{
				Stage:         loan.Stage,
				PendingStatus: loan.PendingStatus,
			},
			constants.LoanOperationActionDeleteDocumentFile,
			&types.LoanOperationMetadata{
				LoanDocuments: []*types.LoanDocument{
					{
						Category: loanDocument.Category,
						FileType: loanDocument.FileType,
						Status:   loanDocument.Status,
						Files:    []*types.OperationFile{operationFile},
					},
				},
			},
		)
		if innerErr != nil {
			return errors.Wrap(innerErr, "failed to create operation log")
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &lender.EmptyResponse{}, nil
}
