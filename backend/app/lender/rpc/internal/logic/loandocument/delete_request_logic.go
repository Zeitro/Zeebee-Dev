package loandocumentlogic

import (
	"context"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/jinzhu/copier"
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

type DeleteRequestLogic struct {
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	baseLogic *logic.BaseLogic
	logx.Logger
}

func NewDeleteRequestLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteRequestLogic {
	return &DeleteRequestLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		baseLogic: logic.NewBaseLogic(ctx, svcCtx),
		Logger:    logx.WithContext(ctx),
	}
}

func (l *DeleteRequestLogic) DeleteRequest(in *lender.LoanDocumentDeleteRequestRequest) (resp *lender.EmptyResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	var loan *model.Loan
	var loanDocument *model.LoanDocument

	// validate
	loan, loanDocument, err = l.prepare(in)
	if err != nil {
		return nil, err
	}

	err = l.svcCtx.LoanDocumentRepo.Transaction(l.ctx, func(tx *gorm.DB) error {
		var innerErr error

		// delete files
		fileIds := slice.Map[*model.File, int64](loanDocument.Files, func(_ int, file *model.File) int64 {
			return file.Id
		})
		if len(fileIds) > 0 {
			_, innerErr = l.svcCtx.FileRepo.DeleteByIDsWithTx(l.ctx, tx, fileIds)
			if innerErr != nil {
				return innerErr
			}
		}

		if loanDocument.Category == constants.LoanDocumentUploaderTypeExtend {
			// handle extend type
			_, innerErr = l.svcCtx.LoanDocumentRepo.DeleteByIDWithTx(l.ctx, tx, loanDocument.Id)
			if innerErr != nil {
				return errors.Wrap(innerErr, "failed to delete loan document")
			}
		} else {
			// handle buildin type
			innerErr = l.svcCtx.LoanDocumentRepo.UpdateColumnsByIdWithTx(l.ctx, tx, loanDocument.Id, map[string]interface{}{
				"status": constants.LoanDocumentStatusPendingUpload,
			})
			if innerErr != nil {
				return errors.Wrap(innerErr, "failed to update loan document status")
			}
		}
		operationLoanDocument := &types.LoanDocument{}
		innerErr = copier.CopyWithOption(operationLoanDocument, loanDocument, copier.Option{
			CaseSensitive: true,
			DeepCopy:      true,
		})
		if innerErr != nil {
			return errorx.GetFailedTransErr(innerErr)
		}

		innerErr = l.baseLogic.CreateLoanOperationLog(
			tx,
			loanDocument.LoanId,
			in.Data.LenderEmployeeId,
			types.LoanOperationStatus{
				Stage:         loan.Stage,
				PendingStatus: loan.PendingStatus,
			},
			types.LoanOperationStatus{
				Stage:         loan.Stage,
				PendingStatus: loan.PendingStatus,
			},
			constants.LoanOperationActionDeleteDocumentRequest,
			&types.LoanOperationMetadata{
				LoanDocuments: []*types.LoanDocument{
					operationLoanDocument,
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

func (l *DeleteRequestLogic) prepare(in *lender.LoanDocumentDeleteRequestRequest) (loan *model.Loan, loanDocument *model.LoanDocument, err error) {

	return l.baseLogic.GetLoanAndDocumentByDocumentId(in.Data.Id, in.Data.LenderCompanyId)
}
