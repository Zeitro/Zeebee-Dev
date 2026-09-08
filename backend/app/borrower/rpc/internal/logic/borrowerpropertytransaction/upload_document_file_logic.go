package borrowerpropertytransactionlogic

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"

	"zeebee-backend/app/borrower/rpc/internal/biz_helper"
	"zeebee-backend/app/borrower/rpc/internal/biz_helper/parse_file"
	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/proto"
)

type UploadDocumentFileLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewUploadDocumentFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadDocumentFileLogic {
	return &UploadDocumentFileLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *UploadDocumentFileLogic) UploadDocumentFile(in *borrower.BorrowerPropertyTransactionDocumentUploadFileRequest) (resp *borrower.BorrowerPropertyTransactionDocumentUploadFileResponse, err error) {
	claim, err := l.baseLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	// FR-009: validate document exists and is leaf node
	doc, err := l.svcCtx.DocumentRepo.FindOneById(l.ctx, in.Data.DocumentId)
	if err != nil {
		err = errorx.Wrap(err, "failed to find document")
		return
	}
	if doc.Type != constants.DocumentTypeLeaf {
		err = errorx.New("can only upload files to leaf documents")
		return
	}

	// validate document belongs to the transaction's document tree
	entityDoc, err := l.svcCtx.EntityDocumentRepo.GetByEntity(l.ctx, constants.EntityDocumentTypeBorrowerPropertyTransaction, in.Data.TransactionId)
	if err != nil {
		err = errorx.New("transaction document binding not found")
		return
	}
	if entityDoc.OrganizationId != claim.OrganizationId {
		err = errorx.New("the document does not belong to your organization")
		return
	}

	if !isDescendantOf(l.ctx, l.svcCtx.DocumentRepo, doc, entityDoc.RootDocumentId) {
		err = errorx.New("document does not belong to this transaction")
		return
	}

	now := &sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}

	var uploadedFiles []model.File

	err = l.svcCtx.DocumentRepo.Transaction(l.ctx, func(tx *gorm.DB) error {
		for _, item := range in.Data.GetFiles() {
			storageType := constants.StorageType(item.StorageType)
			if !storageType.IsValid() {
				return errorx.New("invalid storage type")
			}
			uploadType := constants.UploadType(item.UploadType)
			if !uploadType.IsValid() {
				return errorx.New("invalid upload type")
			}

			fileRecord := &model.File{
				EntityType:  constants.FileEntityTypeDocument,
				EntityId:    doc.Id,
				FileUrl:     item.FileUrl,
				FileName:    item.FileName,
				StorageType: storageType,
				UploadType:  uploadType,
			}

			created, innerErr := l.svcCtx.FileRepo.CreateWithTx(l.ctx, tx, fileRecord, false)
			if innerErr != nil {
				return innerErr
			}
			uploadedFiles = append(uploadedFiles, *created)
		}

		// propagate last_uploaded_at to document and all ancestors
		innerErr := l.propagateLastUploadedAt(tx, doc, now)
		if innerErr != nil {
			return innerErr
		}

		innerErr = l.updateEscrowChecklist(tx, in.Data.TransactionId, doc)
		if innerErr != nil {
			return innerErr
		}

		metadata, _ := json.Marshal(map[string]any{
			"after": map[string]any{
				"files": uploadedFiles,
			},
		})

		// operation log
		innerErr = l.baseLogic.OperationLogic.CreateOperationLog(l.ctx, tx, dto.OperationCreateDto{
			EntityId:       in.Data.TransactionId,
			EntityType:     constants.OperationEntityTypeBorrowerPropertyTransaction,
			OperatorId:     claim.OrganizationUserId,
			OperatorName:   claim.Username,
			OperatorType:   claim.UserType,
			Operation:      constants.OperationActionBorrowerPropertyTransactionUploadDocument,
			OrganizationId: claim.OrganizationId,
			Note:           fmt.Sprintf("upload %d file(s) to document %s", len(uploadedFiles), doc.Name),
			Metadata:       metadata,
		})
		if innerErr != nil {
			return innerErr
		}

		// recalculate checklist
		transHelper := biz_helper.NewBorrowerPropertyTransaction(l.svcCtx.AllRepos)
		_, innerErr = transHelper.RecalculateBorrowerPropertyTransactionChecklist(l.ctx, tx, in.Data.TransactionId)
		return innerErr
	})
	if err != nil {
		return
	}

	// dispatch async extraction on the just-uploaded files.

	requestId := l.parseUploadedFiles(in.Data.TransactionId, doc, uploadedFiles)

	protoFiles, err := copierx.Trans[[]model.File, []*proto.FileEntity](&uploadedFiles, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}

	return &borrower.BorrowerPropertyTransactionDocumentUploadFileResponse{
		Data:      *protoFiles,
		RequestId: requestId,
	}, nil
}

// parseUploadedFiles dispatches async extraction on the just-uploaded files and
// returns the requestId tracking this batch (empty when nothing fired).
func (l *UploadDocumentFileLogic) parseUploadedFiles(transactionId int64, doc *model.Document, files []model.File) (requestId string) {
	var err error
	newCtx := context.WithoutCancel(l.ctx)
	defer func() {
		if err != nil {
			l.svcCtx.Logger.WithContext(newCtx).WithFields(
				logx.Field("stack", errorx.GetErrorStack(err))).Error("fail to parse uploaded files")
		}
	}()
	requestId, err = l.svcCtx.FileParser.Dispatch(newCtx, parse_file.DispatchParam{
		TransactionId: transactionId,
		DocumentId:    doc.Id,
		Files:         files,
		Trigger:       parse_file.TriggerOnUpload,
	})
	return requestId
}

func (l *UploadDocumentFileLogic) propagateLastUploadedAt(tx *gorm.DB, doc *model.Document, t *sql.NullTime) error {
	err := l.svcCtx.DocumentRepo.UpdateColumnsByIdWithTx(l.ctx, tx, doc.Id, map[string]any{
		"last_uploaded_at": t,
	})
	if err != nil {
		return err
	}

	if doc.ParentId != 0 {
		parent, err := l.svcCtx.DocumentRepo.FindOneByIdWithTx(l.ctx, tx, doc.ParentId)
		if err != nil {
			return err
		}
		return l.propagateLastUploadedAt(tx, parent, t)
	}
	return nil
}

// updateEscrowChecklist flags escrow checklist items by uploaded doc name.
func (l *UploadDocumentFileLogic) updateEscrowChecklist(tx *gorm.DB, transactionId int64, doc *model.Document) (err error) {
	escrow, err := l.svcCtx.PropertyTransactionEscrowRepo.FindOneByColumnsWithTx(l.ctx, tx, map[string]any{
		"borrower_property_transaction_id": transactionId,
	})
	if err != nil {
		// no escrow yet is valid (skip); any other lookup error fails the tx.
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}
	checklist := escrow.Checklist
	if !checklist.SetByName(doc.Name, true) {
		return nil // not an escrow doc
	}
	escrow.Checklist = checklist
	_, err = l.svcCtx.PropertyTransactionEscrowRepo.UpdateOneByIdWithTx(l.ctx, tx, escrow.Id, escrow)

	return
}
