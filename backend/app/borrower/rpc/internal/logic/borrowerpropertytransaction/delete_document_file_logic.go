package borrowerpropertytransactionlogic

import (
	"context"

	"gorm.io/gorm"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
)

type DeleteDocumentFileLogic struct {
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	logx.Logger
	baseLogic *logic.BaseLogic
}

func NewDeleteDocumentFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteDocumentFileLogic {
	return &DeleteDocumentFileLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *DeleteDocumentFileLogic) DeleteDocumentFile(in *borrower.DocumentDeleteFileRequest) (resp *borrower.EmptyResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()

	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	// validate transaction exists and has document binding
	entityDoc, err := l.svcCtx.EntityDocumentRepo.GetByEntity(l.ctx, constants.EntityDocumentTypeBorrowerPropertyTransaction, in.Data.TransactionId)
	if err != nil {
		return nil, errorx.Wrap(err, "Transaction document binding not found")
	}

	fileRecord, err := l.svcCtx.FileRepo.FindOneById(l.ctx, in.Data.FileId)
	if err != nil {
		return nil, errorx.Wrap(err, "File not found")
	}

	// validate file's document belongs to the transaction's document tree
	fileDoc, err := l.svcCtx.DocumentRepo.FindOneById(l.ctx, fileRecord.EntityId)
	if err != nil {
		return nil, errorx.Wrap(err, "File document not found")
	}
	if !isDescendantOf(l.ctx, l.svcCtx.DocumentRepo, fileDoc, entityDoc.RootDocumentId) {
		return nil, errorx.New("File does not belong to this transaction")
	}

	// FR-012: only allow deleting manually uploaded files
	if fileRecord.UploadType != constants.UploadTypeManual {
		return nil, errorx.New("Can only delete manually uploaded files")
	}

	err = l.svcCtx.FileRepo.Transaction(l.ctx, func(tx *gorm.DB) error {
		// soft delete file
		_, innerErr := l.svcCtx.FileRepo.DeleteByIDWithTx(l.ctx, tx, fileRecord.Id)
		if innerErr != nil {
			return innerErr
		}

		// operation log
		innerErr = l.baseLogic.OperationLogic.CreateOperationLog(l.ctx, tx, dto.OperationCreateDto{
			EntityId:       in.Data.TransactionId,
			EntityType:     constants.OperationEntityTypeBorrowerPropertyTransaction,
			OperatorId:     claim.OrganizationUserId,
			OperatorName:   claim.Username,
			Operation:      constants.OperationActionBorrowerPropertyTransactionDeleteDocument,
			OrganizationId: claim.OrganizationId,
			Note:           "Delete file from document",
		})
		if innerErr != nil {
			return innerErr
		}

		// FR-017a: recalculate checklist
		_, innerErr = l.baseLogic.RecalculateBorrowerPropertyTransactionChecklist(l.ctx, tx, in.Data.TransactionId)
		return innerErr
	})
	if err != nil {
		return nil, errorx.Wrap(err, "Failed to delete file")
	}

	return &borrower.EmptyResponse{}, nil
}
