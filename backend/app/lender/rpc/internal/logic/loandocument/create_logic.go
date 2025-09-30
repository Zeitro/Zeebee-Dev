package loandocumentlogic

import (
	"context"

	"github.com/duke-git/lancet/v2/slice"
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

type CreateLogic struct {
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	baseLogic *BaseLogic
	logx.Logger
}

func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	return &CreateLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		baseLogic: NewBaseLogic(ctx, svcCtx),
		Logger:    logx.WithContext(ctx),
	}
}

func (l *CreateLogic) Create(in *lender.LoanDocumentCreateExtendTypeRequest) (resp *lender.EmptyResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	var loan *model.Loan
	loan, err = l.prepare(in)
	if err != nil {
		return nil, err
	}

	err = l.svcCtx.LoanDocumentRepo.Transaction(l.ctx, func(tx *gorm.DB) error {
		var innerErr error
		fileTypes := slice.Unique(in.Data.FileTypes)
		loanDocuments := slice.FlatMap[string, *model.LoanDocument](fileTypes, func(_ int, fileType string) []*model.LoanDocument {
			return []*model.LoanDocument{
				{
					LoanId:   in.Data.LoanId,
					Category: constants.LoanDocumentUploaderTypeExtend,
					FileType: constants.LoanDocumentType(fileType),
					Status:   constants.LoanDocumentStatusPendingUpload,
				},
			}
		})
		loanDocuments, innerErr = l.svcCtx.LoanDocumentRepo.BatchCreateWithTx(l.ctx, tx, loanDocuments)
		if innerErr != nil {
			return innerErr
		}
		operationDocs := []*types.LoanDocument{}
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
			in.Data.LenderEmployeeId,
			types.LoanOperationStatus{
				Stage:         loan.Stage,
				PendingStatus: loan.PendingStatus,
			},
			types.LoanOperationStatus{
				Stage:         loan.Stage,
				PendingStatus: loan.PendingStatus,
			},
			constants.LoanOperationActionCreateDocumentType,
			&types.LoanOperationMetadata{
				LoanDocuments: operationDocs,
			},
		)

		return innerErr
	})
	if err != nil {
		return nil, err
	}

	return &lender.EmptyResponse{}, nil
}

func (l *CreateLogic) prepare(in *lender.LoanDocumentCreateExtendTypeRequest) (loan *model.Loan, err error) {

	loan, err = l.baseLogic.GetLenderLoan(in.Data.LoanId, in.Data.LenderCompanyId, nil)
	if err != nil {
		return nil, err
	}

	fileTypes := in.Data.FileTypes
	// check unique
	fileTypes = slice.Unique(fileTypes)
	if len(fileTypes) == 0 {
		return nil, errors.WithStack(errors.New("file types cannot be empty"))
	}

	// checked repeated file types
	var allFileTypes []string
	allFileTypes, err = l.svcCtx.LoanDocumentRepo.GetFileTypesByLoanId(l.ctx, in.Data.LoanId)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	exitedFileTypes := slice.Difference(fileTypes, allFileTypes)
	if len(exitedFileTypes) > 0 {
		return nil, errors.Errorf("file type %s is existed", exitedFileTypes[0])
	}
	return nil, nil
}
