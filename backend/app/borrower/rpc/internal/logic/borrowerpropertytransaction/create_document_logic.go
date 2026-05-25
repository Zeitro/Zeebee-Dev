package borrowerpropertytransactionlogic

import (
	"context"

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

type CreateDocumentLogic struct {
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	logx.Logger
	baseLogic *logic.BaseLogic
}

func NewCreateDocumentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateDocumentLogic {
	return &CreateDocumentLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *CreateDocumentLogic) CreateDocument(in *borrower.DocumentCreateRequest) (resp *borrower.DocumentCreateResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()

	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	// validate entity_document binding exists for the transaction
	entityDoc, err := l.svcCtx.EntityDocumentRepo.GetByEntity(l.ctx, constants.EntityDocumentTypeBorrowerPropertyTransaction, in.Data.TransactionId)
	if err != nil {
		return nil, errorx.Wrap(err, "Entity document not found for transaction")
	}

	// validate parent exists and belongs to the transaction's document tree
	parent, err := l.svcCtx.DocumentRepo.FindOneById(l.ctx, in.Data.ParentId)
	if err != nil {
		return nil, errorx.Wrap(err, "Parent document not found")
	}
	if !isDescendantOf(l.ctx, l.svcCtx.DocumentRepo, parent, entityDoc.RootDocumentId) {
		return nil, errorx.New("Parent document does not belong to this transaction")
	}

	// FR-014: name uniqueness among siblings
	siblings, err := l.svcCtx.DocumentRepo.GetChildren(l.ctx, parent.Id)
	if err != nil {
		return nil, err
	}
	for _, s := range siblings {
		if s.Name == in.Data.Name {
			return nil, errorx.New("Document name already exists at this level")
		}
	}

	// FR-013: only create leaf documents at same level as existing leaves
	doc := &model.Document{
		ParentId:  parent.Id,
		Type:      "leaf",
		Name:      in.Data.Name,
		Level:     parent.Level + 1,
		SortOrder: int32(len(siblings)) + 1,
	}

	err = l.svcCtx.DocumentRepo.Transaction(l.ctx, func(tx *gorm.DB) error {
		var innerErr error

		doc, innerErr = l.svcCtx.DocumentRepo.CreateWithTx(l.ctx, tx, doc, false)
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
			Note:           "Create document node: " + in.Data.Name,
		})
		if innerErr != nil {
			return innerErr
		}

		// FR-017a: recalculate checklist
		_, innerErr = l.baseLogic.RecalculateBorrowerPropertyTransactionChecklist(l.ctx, tx, in.Data.TransactionId)
		return innerErr
	})
	if err != nil {
		return nil, errorx.Wrap(err, "Failed to create document")
	}

	protoDoc, err := copierx.Trans[model.Document, borrower.DocumentEntity](doc, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}

	return &borrower.DocumentCreateResponse{
		Data: protoDoc,
	}, nil
}
