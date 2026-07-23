package biz_helper

import (
	"context"
	"encoding/json"

	"gorm.io/gorm"

	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/repo"
)

type BorrowerPropertyTransaction struct {
	repos *repo.AllRepos
}

func NewBorrowerPropertyTransaction(repos *repo.AllRepos) *BorrowerPropertyTransaction {
	return &BorrowerPropertyTransaction{
		repos: repos,
	}
}

func (h *BorrowerPropertyTransaction) RecalculateBorrowerPropertyTransactionChecklist(ctx context.Context, tx *gorm.DB, transactionId int64) (checklist dto.ChecklistRoot, err error) {

	entityDoc, err := h.repos.EntityDocumentRepo.FindOneByColumnsWithTx(ctx, tx, map[string]interface{}{
		"entity_type": constants.EntityDocumentTypeBorrowerPropertyTransaction,
		"entity_id":   transactionId,
	})
	if err != nil {
		err = errorx.Wrap(err, "entity document not found")
		return
	}

	allDocs, err := h.repos.DocumentRepo.GetAllByEntityDocumentIdWithTx(ctx, tx, entityDoc.Id)
	if err != nil {
		err = errorx.Wrap(err, "fail to load documents")
		return
	}

	childrenMap := buildChildrenMap(allDocs)
	leafIds := collectLeafIds(allDocs, childrenMap)
	fileCountMap, err := h.repos.FileRepo.BatchFileCountByEntityIds(ctx, tx, leafIds, constants.FileEntityTypeDocument)
	if err != nil {
		err = errorx.Wrap(err, "fail to load file counts")
		return
	}

	checklist = BuildChecklistTree(allDocs, fileCountMap)

	existing, err := h.repos.BorrowerPropertyTransactionRepo.FindOneByIdWithTx(ctx, tx, transactionId)
	if err != nil {
		err = errorx.WithStack(err)
		return
	}
	if existing == nil {
		err = errorx.New("transaction not found")
		return
	}

	checklist.AllChecked = existing.Checklist.AllChecked
	existing.Checklist = checklist
	_, err = h.repos.BorrowerPropertyTransactionRepo.UpdateOneByIdWithTx(ctx, tx, transactionId, existing)
	return existing.Checklist, errorx.WithStack(err)
}

// BuildChecklistTree builds a complete ChecklistRoot from loaded documents and file counts.
// Pure computation — no DB dependency, suitable for unit testing.
func BuildChecklistTree(allDocs []*model.Document, fileCountMap map[int64]int32) dto.ChecklistRoot {
	childrenMap := buildChildrenMap(allDocs)
	var rootId int64
	for _, doc := range allDocs {
		if doc.ParentId == 0 {
			rootId = doc.Id
			break
		}
	}

	var requiredNodes, optionalNodes []*dto.ChecklistNode
	var totalReqCount, totalReqUpdatedCount, totalOptCount, totalOptUpdatedCount int32

	for _, child := range childrenMap[rootId] {
		node := buildChecklistNode(child, childrenMap, fileCountMap)
		if node.IsRequired {
			requiredNodes = append(requiredNodes, node)
		} else {
			optionalNodes = append(optionalNodes, node)
		}
		totalReqCount += node.RequiredDocumentCount
		totalOptCount += node.OptionalDocumentCount
		totalReqUpdatedCount += countUploadedLeaves(node, true)
		totalOptUpdatedCount += countUploadedLeaves(node, false)
	}

	return dto.ChecklistRoot{
		TotalDocumentCount:                 totalReqCount + totalOptCount,
		TotalRequiredDocumentCount:         totalReqCount,
		TotalRequiredUploadedDocumentCount: totalReqUpdatedCount,
		TotalOptionalDocumentCount:         totalOptCount,
		TotalOptionalUploadedDocumentCount: totalOptUpdatedCount,
		RequiredDocuments:                  requiredNodes,
		OptionalDocuments:                  optionalNodes,
	}
}

func buildChildrenMap(allDocs []*model.Document) map[int64][]*model.Document {
	childrenMap := make(map[int64][]*model.Document)
	for _, doc := range allDocs {
		if doc.ParentId != 0 {
			childrenMap[doc.ParentId] = append(childrenMap[doc.ParentId], doc)
		}
	}
	return childrenMap
}

// buildChecklistNode recursively builds checklist node from in-memory maps
func buildChecklistNode(doc *model.Document, childrenMap map[int64][]*model.Document, fileCountMap map[int64]int32) *dto.ChecklistNode {
	isRequired, parseRules, tags := readNodeMetadata(doc)

	children := childrenMap[doc.Id]
	if len(children) == 0 {
		var reqUpdated, optUpdated int32
		if isRequired && fileCountMap[doc.Id] > 0 {
			reqUpdated = 1
		} else if !isRequired && fileCountMap[doc.Id] > 0 {
			optUpdated = 1
		}
		return &dto.ChecklistNode{
			Name:                               doc.Name,
			DocumentId:                         doc.Id,
			IsRequired:                         isRequired,
			FileCount:                          fileCountMap[doc.Id],
			TotalRequiredUploadedDocumentCount: reqUpdated,
			TotalOptionalUploadedDocumentCount: optUpdated,
			ParseRules:                         parseRules,
			Tags:                               tags,
		}
	}

	var childNodes []*dto.ChecklistNode
	var subDocCount, uploadedCount, reqDocCount, optDocCount int32
	var totalReqUpdated, totalOptUpdated int32
	for _, child := range children {
		childNode := buildChecklistNode(child, childrenMap, fileCountMap)
		childNodes = append(childNodes, childNode)
		subDocCount++
		if childNode.IsRequired && (childNode.FileCount > 0 || childNode.UploadedDocumentCount > 0) {
			uploadedCount++
		}
		if childNode.IsRequired {
			reqDocCount++
		} else {
			optDocCount++
		}
		totalReqUpdated += childNode.TotalRequiredUploadedDocumentCount
		totalOptUpdated += childNode.TotalOptionalUploadedDocumentCount
	}

	return &dto.ChecklistNode{
		Name:                               doc.Name,
		DocumentId:                         doc.Id,
		IsRequired:                         isRequired,
		SubDocumentCount:                   subDocCount,
		UploadedDocumentCount:              uploadedCount,
		RequiredDocumentCount:              reqDocCount,
		OptionalDocumentCount:              optDocCount,
		TotalRequiredUploadedDocumentCount: totalReqUpdated,
		TotalOptionalUploadedDocumentCount: totalOptUpdated,
		Children:                           childNodes,
		ParseRules:                         parseRules,
		Tags:                               tags,
	}
}

// readNodeMetadata pulls is_required, the file_parse_config rules, and tags out
// of the document's metadata jsonb in one unmarshal. Missing/unreadable
// metadata yields the defaults (required=true, no rules, no tags).
func readNodeMetadata(doc *model.Document) (bool, []dto.DocumentTemplateFileParseRule, []string) {
	isRequired := true
	if len(doc.Metadata) == 0 {
		return isRequired, nil, nil
	}
	var meta struct {
		IsRequired      bool                                 `json:"is_required"`
		FileParseConfig *dto.DocumentTemplateFileParseConfig `json:"file_parse_config"`
		Tags            []string                             `json:"tags"`
	}
	if err := json.Unmarshal(doc.Metadata, &meta); err != nil {
		return isRequired, nil, nil
	}
	return meta.IsRequired, projectParseRules(meta.FileParseConfig), meta.Tags
}

// projectParseRules returns the stored config's rules for the frontend
// checklist. The full rule flows through (parser, source, triggered_by, label,
// selector) so the UI can render each rule's trigger mode. An absent config
// yields no rules.
func projectParseRules(cfg *dto.DocumentTemplateFileParseConfig) []dto.DocumentTemplateFileParseRule {
	if cfg == nil {
		return nil
	}
	return cfg.Rules
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

func countUploadedLeaves(node *dto.ChecklistNode, required bool) int32 {
	if node == nil {
		return 0
	}
	if len(node.Children) == 0 {
		if node.IsRequired == required && node.FileCount > 0 {
			return 1
		}
		return 0
	}
	var count int32
	for _, child := range node.Children {
		count += countUploadedLeaves(child, required)
	}
	return count
}
