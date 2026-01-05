package loandocumentlogic

import (
	"context"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"

	"zeebee-backend/app/lender/rpc/internal/logic"
	"zeebee-backend/app/lender/rpc/internal/svc"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/utils"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/model"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/common/errorx"
)

type CreateLogic struct {
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	baseLogic *logic.BaseLogic
	logx.Logger
}

func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	return &CreateLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		baseLogic: logic.NewBaseLogic(ctx, svcCtx),
		Logger:    logx.WithContext(ctx),
	}
}

func (l *CreateLogic) Create(in *lender.LoanDocumentCreateExtendTypeRequest) (resp *lender.LoanDocumentCreateExtendTypeResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	var loan *model.Loan
	results := []*lender.LoanDocumentOptEntity{}
	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	loan, err = l.prepare(in, claim)
	if err != nil {
		return nil, err
	}

	err = l.svcCtx.LoanDocumentRepo.Transaction(l.ctx, func(tx *gorm.DB) error {
		var innerErr error
		category := constants.LoanDocumentCategory(in.Data.Category)
		status := constants.LoanDocumentStatusPendingRequested
		if category == constants.LoanDocumentCategoryLenderClosingExtend {
			status = constants.LoanDocumentStatusPendingUpload
		}
		fileTypes := slice.Unique(in.Data.FileTypes)
		loanDocuments := slice.FlatMap[string, *model.LoanDocument](fileTypes, func(_ int, fileType string) []*model.LoanDocument {
			return []*model.LoanDocument{
				{
					LoanId:   in.Data.LoanId,
					Category: category,
					FileType: constants.LoanDocumentFileType(fileType),
					Status:   status,
				},
			}
		})
		loanDocuments, innerErr = l.svcCtx.LoanDocumentRepo.BatchCreateWithTx(l.ctx, tx, loanDocuments)
		if innerErr != nil {
			return innerErr
		}

		innerErr = copier.CopyWithOption(&results, loanDocuments, copier.Option{
			IgnoreEmpty:   false,
			DeepCopy:      true,
			CaseSensitive: true,
		})
		if innerErr != nil {
			return errorx.GetFailedTransErr(innerErr)
		}

		operationDocs := []*model.OperationLoanDocument{}
		innerErr = copier.CopyWithOption(&operationDocs, loanDocuments, copier.Option{
			IgnoreEmpty:   false,
			DeepCopy:      true,
			CaseSensitive: true,
		})
		if innerErr != nil {
			return errorx.GetFailedTransErr(innerErr)
		}

		innerErr = l.baseLogic.CreateLoanOperationLog(
			tx,
			in.Data.LoanId,
			claim.OrganizationUserId,
			model.LoanOperationStatus{
				Stage:         loan.Stage,
				PendingStatus: loan.PendingStatus,
			},
			model.LoanOperationStatus{
				Stage:         loan.Stage,
				PendingStatus: loan.PendingStatus,
			},
			constants.LoanOperationActionCreateDocumentType,
			&model.LoanOperationMetadata{
				LoanDocuments: operationDocs,
			},
		)

		return innerErr
	})
	if err != nil {
		return nil, err
	}

	return &lender.LoanDocumentCreateExtendTypeResponse{
		Data: results,
	}, nil
}

func (l *CreateLogic) prepare(in *lender.LoanDocumentCreateExtendTypeRequest, claim utils.JwtUser) (loan *model.Loan, err error) {

	category := constants.LoanDocumentCategory(in.Data.Category)
	if !slice.Contain([]constants.LoanDocumentCategory{
		constants.LoanDocumentCategoryExtend,
		constants.LoanDocumentCategoryLenderClosingExtend,
	}, category) {
		return nil, errorx.WithStack(errorx.New("Loan Document Category Not Supported"))
	}

	loan, err = l.baseLogic.GetLenderLoan(in.Data.LoanId, claim.OrganizationId, nil)
	if err != nil {
		return nil, err
	}

	fileTypes := in.Data.FileTypes
	// check unique
	fileTypes = slice.Unique(fileTypes)
	if len(fileTypes) == 0 {
		return nil, errorx.WithStack(errorx.New("File Types Cannot Be Empty"))
	}

	// checked repeated file types
	var allFileTypes []string
	allFileTypes, err = l.svcCtx.LoanDocumentRepo.GetFileTypesByLoanId(l.ctx, in.Data.LoanId)
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	exitedFileTypes := slice.Intersection(fileTypes, allFileTypes)
	if len(exitedFileTypes) > 0 {
		return nil, errorx.Errorf("File Type %s Is Existed", exitedFileTypes[0])
	}
	return loan, nil
}
