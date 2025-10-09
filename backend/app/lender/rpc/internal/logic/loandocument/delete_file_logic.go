package loandocumentlogic

import (
	"context"

	"github.com/pkg/errors"
	"gorm.io/gorm"

	"zeebee-backend/app/lender/rpc/internal/logic"
	"zeebee-backend/app/lender/rpc/internal/repo/model"
	"zeebee-backend/app/lender/rpc/internal/svc"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/constants"
	"zeebee-backend/common/types"

	"github.com/zeromicro/go-zero/core/logx"

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

func (l *DeleteFileLogic) DeleteFile(in *lender.LoanDocumentDeleteFileRequest) (resp *lender.EmptyResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	var loan *model.Loan
	var loanDocument *model.LoanDocument

	loan, loanDocument, _, err = l.prepare(in)
	if err != nil {
		return nil, err
	}

	err = l.svcCtx.FileRepo.Transaction(l.ctx, func(tx *gorm.DB) error {
		var innerErr error
		var file *model.File
		file, innerErr = l.svcCtx.FileRepo.FindOneByIdWithTx(l.ctx, tx, in.Data.FileId)
		if innerErr != nil {
			return errors.Wrap(innerErr, "failed to find file")
		}
		_, innerErr = l.svcCtx.FileRepo.DeleteByIDWithTx(l.ctx, tx, in.Data.FileId)
		if innerErr != nil {
			return errors.Wrap(err, "failed to delete file")
		}

		var fileCount int64
		fileCount, innerErr = l.svcCtx.FileRepo.GetFileCountByEntityIdWithTx(l.ctx, tx, in.Data.FileId, constants.FileEntityTypeLoanDocument)
		if innerErr != nil {
			return errors.Wrap(innerErr, "failed to check remaining files")
		}
		// change LoanDocument status to PendingUpload
		if fileCount == 0 {
			if innerErr = l.svcCtx.LoanDocumentRepo.UpdateColumnsByIdWithTx(l.ctx, tx, loanDocument.Id, map[string]interface{}{
				"status": constants.LoanDocumentStatusPendingUpload,
			}); innerErr != nil {
				return errors.Wrap(innerErr, "failed to update loan document status")
			}
		}

		if innerErr = l.baseLogic.CreateLoanOperationLog(
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
			constants.LoanOperationActionDeleteDocumentFile,
			&types.LoanOperationMetadata{
				LoanDocuments: []*types.LoanDocument{
					{
						Category: loanDocument.Category,
						FileType: loanDocument.FileType,
						Status:   loanDocument.Status,
						Files: []*types.OperationFile{
							{
								EntityId:   file.EntityId,
								EntityType: file.EntityType,
								FileUrl:    file.FileUrl,
							},
						},
					},
				},
			},
		); innerErr != nil {
			return errors.Wrap(innerErr, "failed to create operation log")
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &lender.EmptyResponse{}, nil
}

// prepare
func (l *DeleteFileLogic) prepare(in *lender.LoanDocumentDeleteFileRequest) (loan *model.Loan, loanDocument *model.LoanDocument, file *model.File, err error) {
	file, err = l.svcCtx.FileRepo.FindOneById(l.ctx, in.Data.FileId)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "failed to find file")
	}

	if file.EntityType != constants.FileEntityTypeLoanDocument {
		return nil, nil, nil, errors.WithStack(errors.New("file is not associated with a loan document"))
	}
	loan, loanDocument, err = l.baseLogic.GetLoanAndDocumentByDocumentId(file.EntityId, in.Data.LenderCompanyId)

	return
}
