package borrowerpropertytransactionlogic

import (
	"context"
	"encoding/json"
	"strings"

	"gorm.io/gorm"

	"zeebee-backend/app/borrower/rpc/internal/biz_helper"
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

type CreateDocumentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewCreateDocumentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateDocumentLogic {
	return &CreateDocumentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *CreateDocumentLogic) CreateDocument(in *borrower.BorrowerPropertyTransactionDocumentCreateRequest) (resp *borrower.BorrowerPropertyTransactionDocumentCreateResponse, err error) {
	claim, err := l.baseLogic.GetUserClaimFromCtx(l.ctx)
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
		return nil, errorx.Wrap(err, "failed to find parent document")
	}
	if !isDescendantOf(l.ctx, l.svcCtx.DocumentRepo, parent, entityDoc.RootDocumentId) {
		return nil, errorx.New("parent document does not belong to this transaction")
	}
	if parent.Level != 1 {
		return nil, errorx.New("parent document must be a branch")
	}

	// FR-014: name uniqueness among siblings — check all requested names
	existingNames := map[string]bool{}
	for _, name := range in.Data.Names {
		isExisted, checkErr := l.svcCtx.DocumentRepo.IsExistedByColumns(l.ctx, map[string]interface{}{
			"name":      name,
			"parent_id": parent.Id,
		})
		if checkErr != nil {
			return nil, errorx.Wrap(checkErr, "failed to check document name existence")
		}
		if isExisted {
			existingNames[name] = true
		}
	}

	maxSortOrder, err := l.svcCtx.DocumentRepo.GetMaxSortOrderByParentId(l.ctx, parent.Id)
	if err != nil {
		return nil, errorx.Wrap(err, "failed to get max sort order")
	}

	// unmarshal parent metadata once
	metadata := map[string]interface{}{}
	if err = json.Unmarshal(parent.Metadata, &metadata); err != nil {
		return nil, errorx.Wrap(err, "failed to unmarshal parent document metadata")
	}
	metadata["create_type"] = "Manual"
	metadataContent, _ := json.Marshal(metadata)

	var createdDocs []*model.Document
	var createdNames []string

	err = l.svcCtx.DocumentRepo.Transaction(l.ctx, func(tx *gorm.DB) error {
		var innerErr error
		for i, name := range in.Data.Names {
			if existingNames[name] {
				continue
			}
			doc := &model.Document{
				ParentId:         parent.Id,
				EntityDocumentId: parent.EntityDocumentId,
				Type:             constants.DocumentTypeLeaf,
				Name:             name,
				Level:            parent.Level + 1,
				SortOrder:        int32(maxSortOrder) + int32(i) + 1,
				Metadata:         metadataContent,
			}

			doc, innerErr = l.svcCtx.DocumentRepo.CreateWithTx(l.ctx, tx, doc, false)
			if innerErr != nil {
				return innerErr
			}
			createdDocs = append(createdDocs, doc)
			createdNames = append(createdNames, name)
		}

		// check parent document type
		if parent.Type != constants.DocumentTypeBranch && parent.Type != constants.DocumentTypeRoot {
			innerErr = l.svcCtx.DocumentRepo.UpdateColumnsByIdWithTx(l.ctx, tx, parent.Id, map[string]interface{}{
				"type": constants.DocumentTypeBranch,
			})
			if innerErr != nil {
				return innerErr
			}
		}

		if len(createdDocs) == 0 {
			return errorx.New("all document names already exist at this level")
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
			Note:           "create documents: " + strings.Join(createdNames, ", "),
		})
		if innerErr != nil {
			return innerErr
		}
		transHelper := biz_helper.NewBorrowerPropertyTransaction(l.svcCtx.AllRepos)
		// recalculate checklist
		_, innerErr = transHelper.RecalculateBorrowerPropertyTransactionChecklist(l.ctx, tx, in.Data.TransactionId)
		return innerErr
	})
	if err != nil {
		return nil, err
	}

	var protoDocs []*proto.DocumentEntity
	for _, doc := range createdDocs {
		protoDoc, transErr := copierx.Trans[model.Document, proto.DocumentEntity](doc, copierx.TimeToRFC3339Str())
		if transErr != nil {
			return nil, errorx.GetFailedTransRespErr(transErr)
		}
		protoDocs = append(protoDocs, protoDoc)
	}

	return &borrower.BorrowerPropertyTransactionDocumentCreateResponse{
		Data: protoDocs,
	}, nil
}
