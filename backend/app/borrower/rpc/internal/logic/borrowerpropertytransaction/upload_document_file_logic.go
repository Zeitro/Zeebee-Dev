package borrowerpropertytransactionlogic

import (
	"context"
	"time"

	"gorm.io/gorm"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
)

type UploadDocumentFileLogic struct {
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	logx.Logger
	baseLogic *logic.BaseLogic
}

func NewUploadDocumentFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadDocumentFileLogic {
	return &UploadDocumentFileLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *UploadDocumentFileLogic) UploadDocumentFile(in *borrower.DocumentUploadFileRequest) (resp *borrower.DocumentUploadFileResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()

	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	// FR-009: validate document exists and is leaf node
	doc, err := l.svcCtx.DocumentRepo.FindOneById(l.ctx, in.Data.DocumentId)
	if err != nil {
		return nil, errorx.Wrap(err, "Document not found")
	}
	if doc.Type != "leaf" {
		return nil, errorx.New("Can only upload files to leaf documents")
	}

	// validate document belongs to the transaction's document tree
	entityDoc, err := l.svcCtx.EntityDocumentRepo.GetByEntity(l.ctx, constants.EntityDocumentTypeBorrowerPropertyTransaction, in.Data.TransactionId)
	if err != nil {
		return nil, errorx.Wrap(err, "Transaction document binding not found")
	}
	if !isDescendantOf(l.ctx, l.svcCtx.DocumentRepo, doc, entityDoc.RootDocumentId) {
		return nil, errorx.New("Document does not belong to this transaction")
	}

	// FR-010: create file record with storage_type and upload_type defaults
	fileRecord := &model.File{
		EntityType:  constants.FileEntityTypeDocument,
		EntityId:    doc.Id,
		FileUrl:     in.Data.FileUrl,
		FileName:    in.Data.FileName,
		StorageType: constants.StorageTypeAwsS3,
		UploadType:  constants.UploadTypeManual,
	}

	now := time.Now()

	err = l.svcCtx.DocumentRepo.Transaction(l.ctx, func(tx *gorm.DB) error {
		var innerErr error

		fileRecord, innerErr = l.svcCtx.FileRepo.CreateWithTx(l.ctx, tx, fileRecord, false)
		if innerErr != nil {
			return innerErr
		}

		// FR-011: propagate last_uploaded_at to document and all ancestors
		innerErr = l.propagateLastUploadedAt(tx, doc, now)
		if innerErr != nil {
			return innerErr
		}

		// operation log
		innerErr = l.baseLogic.OperationLogic.CreateOperationLog(l.ctx, tx, dto.OperationCreateDto{
			EntityId:       in.Data.TransactionId,
			EntityType:     constants.OperationEntityTypeBorrowerPropertyTransaction,
			OperatorId:     claim.OrganizationUserId,
			OperatorName:   claim.Username,
			Operation:      constants.OperationActionBorrowerPropertyTransactionUploadDocument,
			OrganizationId: claim.OrganizationId,
			Note:           "Upload file to document: " + doc.Name,
		})
		if innerErr != nil {
			return innerErr
		}

		// FR-017a: recalculate checklist
		_, innerErr = l.baseLogic.RecalculateBorrowerPropertyTransactionChecklist(l.ctx, tx, in.Data.TransactionId)
		return innerErr
	})
	if err != nil {
		return nil, errorx.Wrap(err, "Failed to upload file")
	}

	protoFile, err := copierx.Trans[model.File, borrower.FileEntity](fileRecord, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}

	return &borrower.DocumentUploadFileResponse{
		Data: protoFile,
	}, nil
}

func (l *UploadDocumentFileLogic) propagateLastUploadedAt(tx *gorm.DB, doc *model.Document, t time.Time) error {
	err := l.svcCtx.DocumentRepo.UpdateColumnsByIdWithTx(l.ctx, tx, doc.Id, map[string]interface{}{
		"last_uploaded_at": t,
	})
	if err != nil {
		return err
	}

	if doc.ParentId != 0 {
		parent, err := l.svcCtx.DocumentRepo.FindOneByIdWithTx(l.ctx, tx, doc.ParentId)
		if err != nil {
			return nil
		}
		return l.propagateLastUploadedAt(tx, parent, t)
	}
	return nil
}
