package loandocumentlogic

import (
	"context"

	"github.com/duke-git/lancet/v2/slice"
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

type CreateLogic struct {
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	baseLogic *logic.BaseLogic
}

func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	return &CreateLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *CreateLogic) Create(in *lender.LoanDocumentCreateExtendTypeRequest) (resp *lender.LoanDocumentCreateExtendTypeResponse, err error) {

	var loan *model.Loan
	results := []*proto.LoanDocumentOptEntity{}
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

		resultsPtr, innerErr := copierx.Trans[[]*model.LoanDocument, []*proto.LoanDocumentOptEntity](&loanDocuments)
		if innerErr != nil {
			return innerErr
		}
		results = *resultsPtr

		operationDocs := []*dto.OperationLoanDocument{}
		operationDocsPtr, innerErr := copierx.Trans[[]*model.LoanDocument, []*dto.OperationLoanDocument](&loanDocuments)
		if innerErr != nil {
			return innerErr
		}
		operationDocs = *operationDocsPtr

		innerErr = l.baseLogic.CreateLoanOperationLog(
			l.ctx,
			tx,
			in.Data.LoanId,
			claim.OrganizationUserId,
			dto.LoanOperationStatus{
				Stage:         loan.Stage,
				PendingStatus: loan.PendingStatus,
			},
			dto.LoanOperationStatus{
				Stage:         loan.Stage,
				PendingStatus: loan.PendingStatus,
			},
			constants.LoanOperationActionCreateDocumentType,
			&dto.LoanOperationMetadata{
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
		return nil, errorx.WithStack(errorx.New("Loan document category not supported"))
	}

	loan, err = l.baseLogic.GetLenderLoan(l.ctx, in.Data.LoanId, claim.OrganizationId, nil)
	if err != nil {
		return nil, err
	}

	fileTypes := in.Data.FileTypes
	// check unique
	fileTypes = slice.Unique(fileTypes)
	if len(fileTypes) == 0 {
		return nil, errorx.WithStack(errorx.New("File types cannot be empty"))
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
