package loandocumentlogic

import (
	"context"

	"github.com/jinzhu/copier"
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

type UploadFilesLogic struct {
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	baseLogic *BaseLogic
	logx.Logger
}

func NewUploadFilesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadFilesLogic {
	return &UploadFilesLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		baseLogic: NewBaseLogic(ctx, svcCtx),
		Logger:    logx.WithContext(ctx),
	}
}

func (l *UploadFilesLogic) UploadFiles(in *lender.LoanDocumentUploadFilesRequest) (resp *lender.EmptyResponse, err error) {
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

	err = l.svcCtx.FileRepo.Transaction(l.ctx, func(tx *gorm.DB) error {
		var innerErr error
		// Create file records for uploaded files
		files := make([]*model.File, 0, len(in.Data.FileUrls))
		for _, fileUrl := range in.Data.FileUrls {
			file := &model.File{
				EntityId:   in.Data.Id,
				EntifyType: constants.FileEntityTypeLoanDocument,
				FileUrl:    fileUrl,
			}
			files = append(files, file)
		}

		if len(files) > 0 {
			files, innerErr = l.svcCtx.FileRepo.BatchCreateWithTx(l.ctx, tx, files, false)
			if innerErr != nil {
				return errors.Wrap(innerErr, "failed to create file records")
			}
		}

		// Update loan document status to ReadyForReview
		innerErr = l.svcCtx.LoanDocumentRepo.UpdateColumnsByIdWithTx(l.ctx, tx, in.Data.Id, map[string]interface{}{
			"status": constants.LoanDocumentStatusReadyForReview,
		})
		if innerErr != nil {
			return errors.Wrap(innerErr, "failed to update loan document status")
		}

		// convert files
		operationFiles := make([]*types.OperationFile, 0, len(files))
		innerErr = copier.CopyWithOption(&operationFiles, files, copier.Option{
			DeepCopy:      true,
			CaseSensitive: true,
		})
		if innerErr != nil {
			return errorx.GetFailedTransErr(innerErr)
		}

		// Create operation log
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
			constants.LoanOperationActionUploadDocumentFile,
			&types.LoanOperationMetadata{
				LoanDocuments: []*types.LoanDocument{
					{
						Category: loanDocument.Category,
						FileType: loanDocument.FileType,
						Status:   constants.LoanDocumentStatusReadyForReview,
						Files:    operationFiles,
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

// prepare
func (l *UploadFilesLogic) prepare(in *lender.LoanDocumentUploadFilesRequest) (loan *model.Loan, loanDocument *model.LoanDocument, err error) {
	loanDocument, err = l.svcCtx.LoanDocumentRepo.FindOneWithPreloadById(l.ctx, in.Data.Id, []string{"Files"})
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to find loan document")
	}

	// Validate file count - existing files + new files should not exceed 10
	existingFileCount := int64(len(loanDocument.Files))
	newFileCount := int64(len(in.Data.FileUrls))
	if existingFileCount+newFileCount > 10 {
		return nil, nil, errors.WithStack(errors.New("total number of files (existing + new) cannot exceed 10"))
	}

	loan, err = l.baseLogic.GetLenderLoan(loanDocument.LoanId, in.Data.LenderCompanyId, nil)
	if err != nil {
		return nil, nil, err
	}

	return loan, loanDocument, nil
}
