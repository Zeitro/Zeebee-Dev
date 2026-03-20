package loandocumentlogic

import (
	"context"

	"gorm.io/gorm"

	"zeebee-backend/app/lender/rpc/internal/logic"
	"zeebee-backend/app/lender/rpc/internal/svc"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/utils"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data"
	"zeebee-backend/domains/zeebee/model"

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
	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	loan, loanDocument, _, err = l.prepare(in, claim)
	if err != nil {
		return nil, err
	}

	err = l.svcCtx.FileRepo.Transaction(l.ctx, func(tx *gorm.DB) error {
		var innerErr error
		var file *model.File
		file, innerErr = l.svcCtx.FileRepo.FindOneByIdWithTx(l.ctx, tx, in.Data.FileId)
		if innerErr != nil {
			return errorx.Wrap(innerErr, "Failed To Find File")
		}
		_, innerErr = l.svcCtx.FileRepo.DeleteByIDWithTx(l.ctx, tx, in.Data.FileId)
		if innerErr != nil {
			return errorx.Wrap(err, "Failed To Delete File")
		}

		var fileCount int64
		fileCount, innerErr = l.svcCtx.FileRepo.GetFileCountByEntityIdWithTx(l.ctx, tx, in.Data.FileId, constants.FileEntityTypeLoanDocument)
		if innerErr != nil {
			return errorx.Wrap(innerErr, "Failed To Check Remaining Files")
		}
		// change LoanDocument status to PendingUpload
		if fileCount == 0 {
			documentStatus := constants.LoanDocumentStatusPendingUpload
			if loanDocument.Category == constants.LoanDocumentCategoryBorrower || loanDocument.Category == constants.LoanDocumentCategoryExtend {
				documentStatus = constants.LoanDocumentStatusPendingRequested
			}

			if innerErr = l.svcCtx.LoanDocumentRepo.UpdateColumnsByIdWithTx(l.ctx, tx, loanDocument.Id, map[string]interface{}{
				"status": documentStatus,
			}); innerErr != nil {
				return errorx.Wrap(innerErr, "Failed To Update Loan Document Status")
			}
		}

		if innerErr = l.baseLogic.CreateLoanOperationLog(
			tx,
			in.Data.Id,
			claim.OrganizationUserId,
			data.LoanOperationStatus{
				Stage:         loan.Stage,
				PendingStatus: loan.PendingStatus,
			},
			data.LoanOperationStatus{
				Stage:         loan.Stage,
				PendingStatus: loan.PendingStatus,
			},
			constants.LoanOperationActionDeleteDocumentFile,
			&data.LoanOperationMetadata{
				LoanDocuments: []*data.OperationLoanDocument{
					{
						Category: loanDocument.Category,
						FileType: loanDocument.FileType,
						Status:   loanDocument.Status,
						Files: []*data.OperationFile{
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
			return errorx.Wrap(innerErr, "Failed To Create Operation Log")
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &lender.EmptyResponse{}, nil
}

// prepare
func (l *DeleteFileLogic) prepare(in *lender.LoanDocumentDeleteFileRequest, claim utils.JwtUser) (loan *model.Loan, loanDocument *model.LoanDocument, file *model.File, err error) {
	file, err = l.svcCtx.FileRepo.FindOneById(l.ctx, in.Data.FileId)
	if err != nil {
		return nil, nil, nil, errorx.Wrap(err, "Failed To Find File")
	}

	if file.EntityType != constants.FileEntityTypeLoanDocument {
		return nil, nil, nil, errorx.WithStack(errorx.New("File is not associated with a loan document"))
	}
	loan, loanDocument, err = l.baseLogic.GetLoanAndDocumentByDocumentId(file.EntityId, claim.OrganizationId)

	return
}
