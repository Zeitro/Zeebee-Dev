package documentlogic

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

type DeleteFileLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	baseLogic *logic.BaseLogic
}

func NewDeleteFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteFileLogic {
	return &DeleteFileLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *DeleteFileLogic) DeleteFile(in *borrower.DocumentDeleteFileRequest) (resp *borrower.EmptyResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	fileRecord, err := l.svcCtx.FileRepo.FindOneById(l.ctx, in.Data.FileId)
	if err != nil {
		return nil, errorx.Wrap(err, "File not found")
	}

	// only allow deleting manually uploaded files
	if fileRecord.UploadType != constants.UploadTypeManual {
		return nil, errorx.New("Can only delete manually uploaded files")
	}

	err = l.svcCtx.FileRepo.Transaction(l.ctx, func(tx *gorm.DB) error {
		// soft delete file
		_, innerErr := l.svcCtx.FileRepo.DeleteByIDWithTx(l.ctx, tx, fileRecord.Id)
		if innerErr != nil {
			return innerErr
		}

		// create operation log
		innerErr = l.baseLogic.OperationLogic.CreateOperationLog(l.ctx, tx, dto.OperationCreateDto{
			EntityId:       fileRecord.EntityId,
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

		// recalculate checklist for parent transaction
		transactionId, _ := l.baseLogic.FindTransactionIdByDocumentId(l.ctx, fileRecord.EntityId)
		if transactionId > 0 {
			innerErr = l.baseLogic.RecalculateChecklist(l.ctx, tx, transactionId)
		}
		return innerErr
	})
	if err != nil {
		return nil, errorx.Wrap(err, "Failed to delete file")
	}

	return &borrower.EmptyResponse{}, nil
}
