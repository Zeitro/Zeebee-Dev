package borrowerpropertytransactionlogic

import (
	"context"
	"encoding/json"

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
	ctx    context.Context
	svcCtx *svc.ServiceContext
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
		return nil, errorx.Wrap(err, "failed to find entity document")
	}

	if entityDoc.OrganizationId != claim.OrganizationId {
		return nil, errorx.New("the document does not belong to your organization")
	}

	// validate parent exists and belongs to the transaction's document tree
	parent, err := l.svcCtx.DocumentRepo.FindOneById(l.ctx, in.Data.ParentId)
	if err != nil {
		err = errorx.Wrap(err, "failed to find parent document")
		return
	}
	if !isDescendantOf(l.ctx, l.svcCtx.DocumentRepo, parent, entityDoc.RootDocumentId) {
		err = errorx.New("parent document does not belong to this transaction")
		return
	}
	if parent.Type != constants.DocumentTypeBranch {
		err = errorx.New("parent document must be a branch")
		return
	}

	// FR-014: name uniqueness among siblings
	isExisted, err := l.svcCtx.DocumentRepo.IsExistedByColumns(l.ctx, map[string]interface{}{
		"name":      in.Data.Name,
		"parent_id": parent.Id,
	})
	if err != nil {
		err = errorx.Wrap(err, "failed to check document name existence")
		return
	}
	if isExisted {
		return nil, errorx.New("document name already exists at this level")
	}
	maxSortOrder, err := l.svcCtx.DocumentRepo.GetMaxSortOrderByParentId(l.ctx, parent.Id)
	if err != nil {
		err = errorx.Wrap(err, "failed to get max sort order")
		return
	}

	// FR-013: only create leaf documents at same level as existing leaves
	doc := &model.Document{
		ParentId:         parent.Id,
		EntityDocumentId: parent.EntityDocumentId,
		Type:             constants.DocumentTypeLeaf,
		Name:             in.Data.Name,
		Level:            parent.Level + 1,
		SortOrder:        int32(maxSortOrder) + 1,
	}

	err = l.svcCtx.DocumentRepo.Transaction(l.ctx, func(tx *gorm.DB) error {
		var innerErr error
		// create document
		metadata := map[string]interface{}{}
		innerErr = json.Unmarshal([]byte(parent.Metadata), &metadata)
		if innerErr != nil {
			innerErr = errorx.Wrap(innerErr, "failed to unmarshal parent document metadata")
			return innerErr
		}
		metadata["create_type"] = "Manual"
		metadataContent, _ := json.Marshal(metadata)
		doc.Metadata = metadataContent
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
			OperatorType:   claim.UserType,
			Operation:      constants.OperationActionBorrowerPropertyTransactionAddDocument,
			OrganizationId: claim.OrganizationId,
			Note:           "create document: " + in.Data.Name,
		})
		if innerErr != nil {
			return innerErr
		}

		// recalculate checklist
		_, innerErr = l.baseLogic.RecalculateBorrowerPropertyTransactionChecklist(l.ctx, tx, in.Data.TransactionId)
		return innerErr
	})
	if err != nil {
		return
	}

	protoDoc, err := copierx.Trans[model.Document, borrower.DocumentEntity](doc, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}

	return &borrower.DocumentCreateResponse{
		Data: protoDoc,
	}, nil
}
