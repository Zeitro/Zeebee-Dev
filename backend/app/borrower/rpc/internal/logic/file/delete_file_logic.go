package filelogic

import (
	"context"

	"github.com/pkg/errors"
	"gorm.io/gorm"

	"zeebee-backend/app/borrower/rpc/internal/repo/model"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/constants"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/common/errorx"
)

type DeleteFileLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteFileLogic {
	return &DeleteFileLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteFileLogic) DeleteFile(in *borrower.FileDeleteRequest) (resp *borrower.EmptyResponse, err error) {
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

	loanDocument, err = l.svcCtx.LoanDocumentRepo.FindOneWithPreloadById(l.ctx, file.EntityId, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find loan document")
	}

	loan, err = l.svcCtx.LoanRepo.GetBorrowerLoan(l.ctx, loanDocument.LoanId, in.Data.BorrowerId)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get loan")
	}
	if loan.BorrowerId != in.Data.BorrowerId {
		return nil, errors.WithStack(errors.New("loan does not belong to borrower"))
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
			loanDocumentStatus := constants.LoanDocumentStatusPendingUpload
			if loanDocument.Category == constants.LoanDocumentCategoryBorrower || loanDocument.Category == constants.LoanDocumentCategoryExtend {
				loanDocumentStatus = constants.LoanDocumentStatusPendingRequested
			}

			innerErr = l.svcCtx.LoanDocumentRepo.UpdateColumnsByIdWithTx(l.ctx, tx, loanDocument.Id, map[string]interface{}{
				"status": loanDocumentStatus,
			})
			if innerErr != nil {
				return errors.Wrap(innerErr, "failed to update loan document status")
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &borrower.EmptyResponse{}, nil
}
