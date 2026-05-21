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

	var rootDocId int64
	err := tx.WithContext(ctx).Table("entity_documents").
		Select("root_document_id").
		Where("entity_type = ? AND entity_id = ?", string(constants.EntityDocumentTypeBorrowerPropertyTransaction), transactionId).
		Scan(&rootDocId).Error
	if err != nil {
		return errorx.Wrap(err, "Entity document not found")
	}

	children, err := repos.DocumentRepo.GetChildren(ctx, rootDocId)
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

	allChecked := false
	existing, _ := repos.BorrowerPropertyTransactionRepo.GetOne(ctx, transactionId)
	if existing != nil {
		var prev dto.ChecklistRoot
		if json.Unmarshal(existing.Checklist, &prev) == nil {
			allChecked = prev.AllChecked
		}
	}

	checklist := dto.ChecklistRoot{
		AllChecked:                        allChecked,
		TotalDocumentCount:                totalReqCount + totalOptCount,
		TotalRequiredDocumentCount:        totalReqCount,
		TotalRequiredUpdatedDocumentCount: totalReqUpdatedCount,
		TotalOptionalDocumentCount:        totalOptCount,
		TotalOptionalUpdatedDocumentCount: totalOptUpdatedCount,
		RequiredDocuments:                 requiredNodes,
		OptionalDocuments:                 optionalNodes,
	}

	checklistJSON, err := json.Marshal(checklist)
	if err != nil {
		return errorx.Wrap(err, "Failed to marshal checklist")
	}

	return tx.Table("borrower_property_transactions").
		Where("id = ?", transactionId).
		UpdateColumn("checklist", checklistJSON).Error
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
	err = repos.BorrowerPropertyTransactionRepo.Client(ctx).Table("entity_documents").
		Select("entity_id").
		Where("root_document_id = ? AND entity_type = ?", doc.Id, string(constants.EntityDocumentTypeBorrowerPropertyTransaction)).
		Scan(&entityId).Error
	return entityId, errorx.WithStack(err)
}
