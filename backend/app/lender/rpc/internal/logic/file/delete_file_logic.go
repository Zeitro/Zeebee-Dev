package filelogic

import (
	"context"

	"gorm.io/gorm"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/lender/rpc/internal/svc"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/constants"
	"zeebee-backend/domains/zeebee/model"

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
		return nil, errorx.Wrap(err, "Failed To Find File")
	}

	if file.EntityType != constants.FileEntityTypeLoanDocument {
		return nil, errorx.WithStack(errorx.New("File Is Not Associated With A Loan Document"))
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
			return errorx.Wrap(innerErr, "Failed To Delete File")
		}

		// get file count
		var fileCount int64
		fileCount, innerErr = l.svcCtx.FileRepo.GetFileCountByEntityIdWithTx(l.ctx, tx, loanDocument.Id, constants.FileEntityTypeLoanDocument)
		if innerErr != nil {
			return errorx.Wrap(innerErr, "Failed To Check Remaining Files")
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
				return errorx.Wrap(innerErr, "Failed To Update Loan Document Status")
			}
		}

		// create operation log
		operationFile := &model.OperationFile{
			EntityId:   file.EntityId,
			EntityType: file.EntityType,
			FileUrl:    file.FileUrl,
		}

		innerErr = l.baseLogic.CreateLoanOperationLog(
			tx,
			loanDocument.LoanId,
			in.Data.LenderEmployeeId,
			model.LoanOperationStatus{
				Stage:         loan.Stage,
				PendingStatus: loan.PendingStatus,
			},
			model.LoanOperationStatus{
				Stage:         loan.Stage,
				PendingStatus: loan.PendingStatus,
			},
			constants.LoanOperationActionDeleteDocumentFile,
			&model.LoanOperationMetadata{
				LoanDocuments: []*model.OperationLoanDocument{
					{
						Category: loanDocument.Category,
						FileType: loanDocument.FileType,
						Status:   loanDocument.Status,
						Files:    []*model.OperationFile{operationFile},
					},
				},
			},
		)
		if innerErr != nil {
			return errorx.Wrap(innerErr, "Failed To Create Operation Log")
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &lender.EmptyResponse{}, nil
}
