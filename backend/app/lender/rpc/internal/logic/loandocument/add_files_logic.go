package loandocumentlogic

import (
	"context"
	"time"

	"gorm.io/gorm"

	"zeebee-backend/app/lender/rpc/internal/logic"
	"zeebee-backend/app/lender/rpc/internal/svc"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/utils"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/proto"

	"zeebee-backend/common/errorx"
)

type AddFilesLogic struct {
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	baseLogic *logic.BaseLogic
}

func NewAddFilesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddFilesLogic {
	return &AddFilesLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *AddFilesLogic) AddFiles(in *lender.LoanDocumentAddFilesRequest) (resp *proto.EmptyResponse, err error) {

	var loan *model.Loan
	var loanDocument *model.LoanDocument

	claim, err := l.baseLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	loan, loanDocument, err = l.prepare(in, claim)
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
				EntityType: constants.FileEntityTypeLoanDocument,
				FileUrl:    fileUrl,
			}
			files = append(files, file)
		}

		if len(files) > 0 {
			files, innerErr = l.svcCtx.FileRepo.BatchCreateWithTx(l.ctx, tx, files, false)
			if innerErr != nil {
				return errorx.Wrap(innerErr, "fail to create file records")
			}
		}

		// Update loan document status to ReadyForReview
		innerErr = l.svcCtx.LoanDocumentRepo.UpdateColumnsByIdWithTx(l.ctx, tx, in.Data.Id, map[string]interface{}{
			"status": constants.LoanDocumentStatusReadyForReview,
		})
		if innerErr != nil {
			return errorx.Wrap(innerErr, "fail to update loan document status")
		}

		// convert files
		operationFiles := make([]*dto.OperationFile, 0, len(files))
		operationFilesPtr, err := copierx.Trans[[]*model.File, []*dto.OperationFile](&files)
		if err != nil {
			return err
		}
		operationFiles = *operationFilesPtr

		// Create operation log
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
			constants.LoanOperationActionUploadDocumentFile,
			&dto.LoanOperationMetadata{
				LoanDocuments: []*dto.OperationLoanDocument{
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
			return errorx.Wrap(innerErr, "fail to create operation log")
		}

		// upload loan document upload time
		innerErr = l.svcCtx.LoanDocumentRepo.UpdateColumnsByIdWithTx(l.ctx, tx, loanDocument.Id, map[string]interface{}{
			"last_uploaded_at": time.Now(),
		})

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &proto.EmptyResponse{}, nil
}

// prepare
func (l *AddFilesLogic) prepare(in *lender.LoanDocumentAddFilesRequest, claim utils.JwtUser) (loan *model.Loan, loanDocument *model.LoanDocument, err error) {
	loanDocument, err = l.svcCtx.LoanDocumentRepo.FindOneWithPreloadById(l.ctx, in.Data.Id, []string{"Files"})
	if err != nil {
		return nil, nil, errorx.Wrap(err, "fail to find loan document")
	}

	// Validate file count - existing files + new files should not exceed 10
	existingFileCount := int64(len(loanDocument.Files))
	newFileCount := int64(len(in.Data.FileUrls))
	if existingFileCount+newFileCount > 10 {
		return nil, nil, errorx.WithStack(errorx.New("total number of files (existing + new) cannot exceed 10"))
	}

	loan, err = l.baseLogic.GetLenderLoan(l.ctx, loanDocument.LoanId, claim.OrganizationId, nil)
	if err != nil {
		return nil, nil, err
	}

	return loan, loanDocument, nil
}
