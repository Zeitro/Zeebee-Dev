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

func (b *BaseLogic) RecalculateBorrowerPropertyTransactionChecklist(ctx context.Context, tx *gorm.DB, transactionId int64) (checklist dto.ChecklistRoot,err error) {
	repos := b.svcCtx.AllRepos

	entityDoc, err := b.svcCtx.AllRepos.EntityDocumentRepo.FindOneByColumnsWithTx(ctx, tx, map[string]interface{}{
		"entity_type": constants.EntityDocumentTypeBorrowerPropertyTransaction,
		"entity_id":   transactionId,
	})
	if err != nil {
		err =  errorx.Wrap(err, "Entity document not found")
		return
	}

	// single query: fetch all documents for this entity
	allDocs, err := repos.DocumentRepo.GetAllByEntityDocumentIdWithTx(ctx, tx, entityDoc.Id)
	if err != nil {
		err =  errorx.Wrap(err, "Failed to load documents")
		return
	}

	// build parent-child map and find root
	childrenMap := make(map[int64][]*model.Document)
	var rootId int64
	for _, doc := range allDocs {
		if doc.ParentId == 0 {
			rootId = doc.Id
			continue
		}
		childrenMap[doc.ParentId] = append(childrenMap[doc.ParentId], doc)
	}

	// batch load file counts for all leaf documents
	leafIds := collectLeafIds(allDocs, childrenMap)
	fileCountMap, _ := repos.FileRepo.BatchFileCountByEntityIds(ctx, tx, leafIds, constants.FileEntityTypeDocument)

	// build checklist tree from root's direct children
	rootChildren := childrenMap[rootId]
	var requiredNodes, optionalNodes []dto.ChecklistNode
	var totalReqCount, totalReqUpdatedCount, totalOptCount, totalOptUpdatedCount int32

	for _, child := range rootChildren {
		node := buildChecklistNode(child, childrenMap, fileCountMap)
		if node.IsRequired {
			requiredNodes = append(requiredNodes, node)
			totalReqCount += countLeaves(node)
			totalReqUpdatedCount += countUploadedLeaves(node)
		} else {
			optionalNodes = append(optionalNodes, node)
			totalOptCount += countLeaves(node)
			totalOptUpdatedCount += countUploadedLeaves(node)
		}
	}

	existing, err := repos.BorrowerPropertyTransactionRepo.FindOneByIdWithTx(ctx, tx, transactionId)
	if err != nil {
		err =  errorx.WithStack(err)
		return
	}
	if existing == nil {
		err =  errorx.New("transaction not found")
		return
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
	return existing.Checklist, errorx.WithStack(err)
}

// buildChecklistNode recursively builds checklist node from in-memory maps
func buildChecklistNode(doc *model.Document, childrenMap map[int64][]*model.Document, fileCountMap map[int64]int32) dto.ChecklistNode {
	isRequired := true
	if doc.Metadata != nil {
		var meta struct {
			IsRequired bool `json:"is_required"`
		}
		_ = json.Unmarshal(doc.Metadata, &meta)
		isRequired = meta.IsRequired
	}

	children := childrenMap[doc.Id]
	if len(children) == 0 {
		return dto.ChecklistNode{
			Name:       doc.Name,
			DocumentId: doc.Id,
			IsRequired: isRequired,
			FileCount:  fileCountMap[doc.Id],
		}
	}

	var childNodes []dto.ChecklistNode
	var subDocCount, uploadedCount int32
	for _, child := range children {
		childNode := buildChecklistNode(child, childrenMap, fileCountMap)
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

// collectLeafIds returns document IDs that have no children (leaf nodes)
func collectLeafIds(allDocs []*model.Document, childrenMap map[int64][]*model.Document) []int64 {
	var ids []int64
	for _, doc := range allDocs {
		if doc.ParentId == 0 {
			continue
		}
		if len(childrenMap[doc.Id]) == 0 {
			ids = append(ids, doc.Id)
		}
	}
	return ids
}

func countLeaves(node dto.ChecklistNode) int32 {
	if len(node.Children) == 0 {
		return 1
	}
	var count int32
	for _, child := range node.Children {
		count += countLeaves(child)
	}
	return count
}

func countUploadedLeaves(node dto.ChecklistNode) int32 {
	if len(node.Children) == 0 {
		if node.FileCount > 0 {
			return 1
		}
		return 0
	}
	var count int32
	for _, child := range node.Children {
		count += countUploadedLeaves(child)
	}
	return count
}

func (b *BaseLogic) FindTransactionIdByDocumentId(ctx context.Context, documentId int64) (int64, error) {
	repos := b.svcCtx.AllRepos

	// walk up to root document, then use entity_document_id for single lookup
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

	if doc.EntityDocumentId == 0 {
		return 0, errorx.New("document not bound to entity")
	}

	entityDoc, err := b.svcCtx.AllRepos.EntityDocumentRepo.FindOneById(ctx, doc.EntityDocumentId)
	if err != nil {
		return 0, errorx.Wrap(err, "entity document not found")
	}

	if entityDoc.EntityId == 0 {
		return 0, errorx.New("transaction not found")
	}
	return entityDoc.EntityId, nil
}
