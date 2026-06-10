package loandocumentlogic

import (
	"context"

	"github.com/duke-git/lancet/v2/slice"
	"gorm.io/gorm"

	"zeebee-backend/app/lender/rpc/internal/logic"
	"zeebee-backend/app/lender/rpc/internal/svc"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data"
	"zeebee-backend/domains/zeebee/model"

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
		baseLogic: logic.GetBaseLogic(svcCtx),
		Logger:    logx.WithContext(ctx),
	}
}

func (l *DeleteRequestLogic) DeleteRequest(in *lender.LoanDocumentDeleteRequestRequest) (resp *lender.EmptyResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(
				logx.Field("stacks", errorx.GetErrorStack(err)),
				logx.Field("service_name", l.svcCtx.Config.Name),
			).Error(err)
		}
	}()
	var loan *model.Loan
	var loanDocument *model.LoanDocument

	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	// validate
	loan, loanDocument, err = l.prepare(in, claim.OrganizationId)
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

		switch loanDocument.Category {
		case constants.LoanDocumentCategoryExtend:
			_, innerErr = l.svcCtx.LoanDocumentRepo.DeleteByIDWithTx(l.ctx, tx, loanDocument.Id)
			if innerErr != nil {
				return errorx.Wrap(innerErr, "Failed To Delete Loan Document")
			}
		case constants.LoanDocumentCategoryBorrower:
			innerErr = l.svcCtx.LoanDocumentRepo.UpdateColumnsByIdWithTx(l.ctx, tx, loanDocument.Id, map[string]interface{}{
				"status": constants.LoanDocumentStatusPendingUpload,
			})
			if innerErr != nil {
				return errorx.Wrap(innerErr, "Failed To Update Loan Document Status")
			}
		default:
			return errorx.WithStack(errorx.New("Invalid loan document category"))
		}
		operationLoanDocument := &data.OperationLoanDocument{}
		operationLoanDocument, err = copierx.Trans[model.LoanDocument, data.OperationLoanDocument](loanDocument)
		if err != nil {
			return errorx.Wrap(err, "Failed To Copy Loan Document")
		}

		innerErr = l.baseLogic.CreateLoanOperationLog(
			l.ctx,
			tx,
			loanDocument.LoanId,
			claim.OrganizationUserId,
			data.LoanOperationStatus{
				Stage:         loan.Stage,
				PendingStatus: loan.PendingStatus,
			},
			data.LoanOperationStatus{
				Stage:         loan.Stage,
				PendingStatus: loan.PendingStatus,
			},
			constants.LoanOperationActionDeleteDocumentRequest,
			&data.LoanOperationMetadata{
				LoanDocuments: []*data.OperationLoanDocument{
					operationLoanDocument,
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

func (l *DeleteRequestLogic) prepare(in *lender.LoanDocumentDeleteRequestRequest, organizationId int64) (loan *model.Loan, loanDocument *model.LoanDocument, err error) {
	return l.baseLogic.GetLoanAndDocumentByDocumentId(l.ctx, in.Data.Id, organizationId)
}
