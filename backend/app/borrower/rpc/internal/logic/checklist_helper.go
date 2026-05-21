package logic

import (
	"context"
	"encoding/json"

	"gorm.io/gorm"

	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
)

func (b *BaseLogic) RecalculateChecklist(ctx context.Context, tx *gorm.DB, transactionId int64) error {
	repos := b.svcCtx.AllRepos

	entityDoc, err := b.svcCtx.AllRepos.EntityDocumentRepo.FindOneByColumnsWithTx(ctx, tx, map[string]interface{}{
		"entity_type": constants.EntityDocumentTypeBorrowerPropertyTransaction,
		"entity_id":   transactionId,
	})

	if err != nil {
		return errorx.Wrap(err, "Entity document not found")
	}

	children, err := repos.DocumentRepo.GetChildren(ctx, entityDoc.RootDocumentId)
	if err != nil {
		return errorx.Wrap(err, "Failed to get root children")
	}

	var requiredNodes, optionalNodes []dto.ChecklistNode
	var totalReqCount, totalReqUpdatedCount, totalOptCount, totalOptUpdatedCount int32

	for _, child := range children {
		node := b.buildChecklistNode(ctx, child)
		if node.IsRequired {
			requiredNodes = append(requiredNodes, node)
			totalReqCount += 1 + node.SubDocumentCount
			totalReqUpdatedCount += node.UploadedDocumentCount
		} else {
			optionalNodes = append(optionalNodes, node)
			totalOptCount += 1 + node.SubDocumentCount
			totalOptUpdatedCount += node.UploadedDocumentCount
		}
	}

	// allChecked := false
	existing, err := repos.BorrowerPropertyTransactionRepo.FindOneByIdWithTx(ctx, tx, transactionId)
	if err != nil {
		return errorx.WithStack(err)
	}
	if existing == nil {
		return errorx.New("transaction not found")
	}
	
	existing.Checklist = dto.ChecklistRoot{
		AllChecked:                        existing.Checklist.AllChecked,
		TotalDocumentCount:                totalReqCount + totalOptCount,
		TotalRequiredDocumentCount:        totalReqCount,
		TotalRequiredUpdatedDocumentCount: totalReqUpdatedCount,
		TotalOptionalDocumentCount:        totalOptCount,
		TotalOptionalUpdatedDocumentCount: totalOptUpdatedCount,
		RequiredDocuments:                 requiredNodes,
		OptionalDocuments:                 optionalNodes,
	}
	_, err = b.svcCtx.AllRepos.BorrowerPropertyTransactionRepo.UpdateOneByIdWithTx(ctx, tx, transactionId, existing)
	return errorx.WithStack(err)
}

func (b *BaseLogic) buildChecklistNode(ctx context.Context, doc *model.Document) dto.ChecklistNode {
	isRequired := true
	if doc.Metadata != nil {
		var meta struct {
			IsRequired bool `json:"is_required"`
		}
		_ = json.Unmarshal(doc.Metadata, &meta)
		isRequired = meta.IsRequired
	}

	repos := b.svcCtx.AllRepos
	children, _ := repos.DocumentRepo.GetChildren(ctx, doc.Id)

	if len(children) == 0 {
		fileCount, _ := repos.FileRepo.GetFileCountByEntityIdWithTx(ctx, nil, doc.Id, constants.FileEntityTypeDocument)
		return dto.ChecklistNode{
			Name:       doc.Name,
			DocumentId: doc.Id,
			IsRequired: isRequired,
			FileCount:  int32(fileCount),
		}
	}

	var childNodes []dto.ChecklistNode
	var subDocCount, uploadedCount int32
	for _, child := range children {
		childNode := b.buildChecklistNode(ctx, child)
		childNodes = append(childNodes, childNode)
		subDocCount++
		if childNode.FileCount > 0 || childNode.UploadedDocumentCount > 0 {
			uploadedCount++
		}
	}

	return dto.ChecklistNode{
		Name:                  doc.Name,
		DocumentId:            doc.Id,
		IsRequired:            isRequired,
		SubDocumentCount:      subDocCount,
		UploadedDocumentCount: uploadedCount,
		Children:              childNodes,
	}
}

func (b *BaseLogic) FindTransactionIdByDocumentId(ctx context.Context, documentId int64) (int64, error) {
	repos := b.svcCtx.AllRepos
	doc, err := repos.DocumentRepo.FindOneById(ctx, documentId)
	if err != nil {
		return 0, err
	}
	for doc.ParentId != 0 {
		doc, err = repos.DocumentRepo.FindOneById(ctx, doc.ParentId)
		if err != nil {
			return 0, err
		}
	}

	var entityId int64

	err = repos.EntityDocumentRepo.FindColumnByColumns(ctx, "entity_id", &entityId, map[string]interface{}{
		"root_document_id": doc.Id,
		"entity_type":      constants.EntityDocumentTypeBorrowerPropertyTransaction,
	})

	if entityId == 0 {
		return 0, errorx.New("transaction not found")
	}
	return entityId, errorx.WithStack(err)
}
