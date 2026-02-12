package loandocumentlogic

import (
	"context"
	"time"

	"github.com/duke-git/lancet/v2/slice"
	"gorm.io/gorm"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/model"
)

type AddFilesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	baseLogic *logic.BaseLogic
}

func NewAddFilesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddFilesLogic {
	return &AddFilesLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		baseLogic: logic.NewBaseLogic(ctx, svcCtx),
	}
}

func (l *AddFilesLogic) AddFiles(in *borrower.LoanDocumentAddFilesRequest) (resp *borrower.EmptyResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()

	// var loan *model.Loan
	var loanDocument *model.LoanDocument

	_, loanDocument, err = l.prepare(in)
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
				return errorx.Wrap(innerErr, "Failed To Create File Records")
			}
		}

		// Update loan document status to ReadyForReview
		innerErr = l.svcCtx.LoanDocumentRepo.UpdateColumnsByIdWithTx(l.ctx, tx, in.Data.Id, map[string]interface{}{
			"status": constants.LoanDocumentStatusReadyForReview,
		})
		if innerErr != nil {
			return errorx.Wrap(innerErr, "Failed To Update Loan Document Status")
		}

		// upload loan document upload time
		innerErr = l.svcCtx.LoanDocumentRepo.UpdateColumnsByIdWithTx(l.ctx, tx, loanDocument.Id, map[string]interface{}{
			"last_uploaded_at": time.Now(),
		})

		return innerErr
	})

	if err != nil {
		return nil, err
	}

	return &borrower.EmptyResponse{}, nil
}

// prepare
func (l *AddFilesLogic) prepare(in *borrower.LoanDocumentAddFilesRequest) (loan *model.Loan, loanDocument *model.LoanDocument, err error) {
	loanDocument, err = l.svcCtx.LoanDocumentRepo.FindOneWithPreloadById(l.ctx, in.Data.Id, []string{"Files"})
	if err != nil {
		return nil, nil, errorx.Wrap(err, "Failed To Find Loan Document")
	}

	// Validate file count - existing files + new files should not exceed 10
	existingFileCount := int64(len(loanDocument.Files))
	newFileCount := int64(len(in.Data.FileUrls))
	if existingFileCount+newFileCount > 10 {
		return nil, nil, errorx.WithStack(errorx.New("Total number of files (existing + new) cannot exceed 10"))
	}

	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, nil, err
	}

	loan, err = l.svcCtx.LoanRepo.GetBorrowerLoan(l.ctx, loanDocument.LoanId, claim.OrganizationId)
	if err != nil {
		return nil, nil, errorx.Wrap(err, "Failed To Get Loan")
	}

	if loan == nil || loanDocument == nil {
		return nil, nil, errorx.WithStack(errorx.New("Loan or loan document not found"))
	}
	validCategories := []constants.LoanDocumentCategory{
		constants.LoanDocumentCategoryBorrower,
		constants.LoanDocumentCategoryExtend,
	}
	if !slice.Contain(validCategories, loanDocument.Category) {
		return nil, nil, errorx.WithStack(errorx.New("Loan document category not supported"))
	}

	if loanDocument.Status != constants.LoanDocumentStatusRequested && loanDocument.Status != constants.LoanDocumentStatusReadyForReview {
		return nil, nil, errorx.WithStack(errorx.New("Loan document is not in requested or ready for review status"))
	}

	if loan.Stage != constants.LoanStagePending {
		return nil, nil, errorx.WithStack(errorx.New("Loan is not in pending stage"))
	}
	if loan.Source != constants.LoanSourceSystem && loan.Source != constants.LoanSourcePreapproval {
		return nil, nil, errorx.WithStack(errorx.New("Loan source is not System"))
	}

	if loan.BorrowerOrganizationId != claim.OrganizationId {
		return nil, nil, errorx.WithStack(errorx.New("Loan does not belong to borrower"))
	}

	return loan, loanDocument, nil
}
