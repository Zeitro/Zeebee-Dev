// Package parse_file orchestrates asynchronous file-content extraction for
// borrower property transactions: staged Collect → Extract → (Arbitrate) →
// Writeback, with Redis-backed progress, dedup locks, per-stage retry, and
// restart recovery. Extraction primitives (pdfium pool, LLM prompts, field
// parsers) are reused from the parent biz_helper package via parameters.
package parse_file

import (
	"encoding/json"

	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/model"
)

// effectiveRules normalizes a FileParseConfig into its rule list. The v2
// multi-rule form passes through; the legacy single-parser form is treated as
// one OnUpload rule on the node itself; an absent config yields no rules.
// Kept in this package (not exported via API) because only the orchestrator
// consumes it — cohesion over premature sharing.
func effectiveRules(cfg *model.FileParseConfig) []model.FileParseRule {
	if cfg == nil {
		return nil
	}
	if len(cfg.Rules) > 0 {
		return cfg.Rules
	}
	if cfg.Parser != "" {
		return []model.FileParseRule{{Parser: cfg.Parser, Mode: model.FileParseModeOnUpload}}
	}
	return nil
}

// onUploadRules returns only the rules that fire on file upload.
func onUploadRules(cfg *model.FileParseConfig) []model.FileParseRule {
	var out []model.FileParseRule
	for _, r := range effectiveRules(cfg) {
		if r.Mode == model.FileParseModeOnUpload {
			out = append(out, r)
		}
	}
	return out
}

// docFileParseConfig reads file_parse_config from a document's metadata jsonb.
// Missing/unreadable metadata yields nil.
func docFileParseConfig(doc *model.Document) *model.FileParseConfig {
	if doc == nil || len(doc.Metadata) == 0 {
		return nil
	}
	var meta struct {
		FileParseConfig *model.FileParseConfig `json:"file_parse_config"`
		Tags            []string               `json:"tags"`
	}
	if err := json.Unmarshal(doc.Metadata, &meta); err != nil {
		return nil
	}
	return meta.FileParseConfig
}

// docTags reads the tags copied from the template node into document metadata.
func docTags(doc *model.Document) []string {
	if doc == nil || len(doc.Metadata) == 0 {
		return nil
	}
	var meta struct {
		Tags []string `json:"tags"`
	}
	if err := json.Unmarshal(doc.Metadata, &meta); err != nil {
		return nil
	}
	return meta.Tags
}

// CollectSubtree returns the Leaf documents in the subtree rooted at anchorId
// (the anchor itself if it is a leaf, plus all descendant leaves). docs must
// contain the full document set for the transaction so parent/child links can
// be walked via ParentId. Files attach to Leaf documents, so only leaves are
// parse targets.
func CollectSubtree(docs []*model.Document, anchorId int64) []*model.Document {
	byId := make(map[int64]*model.Document, len(docs))
	for _, d := range docs {
		byId[d.Id] = d
	}
	root, ok := byId[anchorId]
	if !ok {
		return nil
	}
	var leaves []*model.Document
	var walk func(*model.Document)
	walk = func(d *model.Document) {
		if d.Type == constants.DocumentTypeLeaf {
			leaves = append(leaves, d)
			return
		}
		for _, c := range docs {
			if c.ParentId == d.Id {
				walk(c)
			}
		}
	}
	walk(root)
	return leaves
}

// CollectBySelector returns Leaf documents whose metadata.tags contain every
// tag in selector.Tags (AND). Used by ManualSelect rules; can span subtrees.
// Empty selector tags matches no document (avoid accidental full-scope parse).
func CollectBySelector(docs []*model.Document, selector *model.TagSelector) []*model.Document {
	if selector == nil || len(selector.Tags) == 0 {
		return nil
	}
	var out []*model.Document
	for _, d := range docs {
		if d.Type != constants.DocumentTypeLeaf {
			continue
		}
		if !containsAll(docTags(d), selector.Tags) {
			continue
		}
		out = append(out, d)
	}
	return out
}

// containsAll reports whether have contains every tag in need (order-independent).
func containsAll(have, need []string) bool {
	if len(have) < len(need) {
		return false
	}
	set := make(map[string]struct{}, len(have))
	for _, t := range have {
		set[t] = struct{}{}
	}
	for _, t := range need {
		if _, ok := set[t]; !ok {
			return false
		}
	}
	return true
}
