package handler

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/common/errorx"
	"zeebee-backend/common/llm"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/repo"
)

// extractSystemPrompt is the shared system prompt every extraction/arbitration
// call uses: force a single raw JSON object, no fences, no prose.
const extractSystemPrompt = "You extract data from documents. Your entire response MUST be a single raw JSON object and nothing else: start with { and end with }. Do NOT wrap it in markdown code fences (no ```json, no ```). Do NOT add prose, explanations, backticks, or any text before or after the JSON."

// base holds the parsing-only deps every handler shares: repos, logger, llm. It
// deliberately holds NO progress manager and NO nodeId — the coordinator (parent
// parse_file package) owns the run lifecycle (lock/job/hold/heartbeat/finalize/
// progress). Embedding base keeps a handler focused on extraction + writeback.
type base struct {
	repos  *repo.AllRepos
	logger logx.Logger
	llm    llm.Client
}

// collectedFile is the flat file list handed from scope resolution to extraction
// so extraction needs no extra DB round-trip per file. Type-agnostic — PDF today,
// image/doc later; the handler's extract step dispatches per type.
type collectedFile struct {
	FileId int64  `json:"file_id"`
	Url    string `json:"url"`
	Name   string `json:"name"`
}

// resolveFiles picks the in-scope files for a parse per rule.Source (the SCOPE),
// independent of how the run was triggered and of file type:
//   - SingleFile / Files → the batch the dispatcher passed in (an upload batch;
//     a manual trigger passes the document's current files);
//   - Directory          → leaf files under the anchor document's subtree;
//   - TagFiltered        → leaf files across the transaction tagged with every
//     selector tag.
//
// Shared so every handler resolves scope the same way; a handler only owns its
// extract + writeback (and dispatches per file type there).
func (b *base) resolveFiles(ctx context.Context, txId, docId int64, files []model.File, rule dto.DocumentTemplateFileParseRule) ([]collectedFile, error) {
	switch rule.Source {
	case constants.DocumentTemplateFileParseSourceTagFiltered:
		return b.collectBySelector(ctx, txId, rule.Selector)
	case constants.DocumentTemplateFileParseSourceDirectory:
		return b.collectSubtree(ctx, txId, docId)
	default: // SingleFile / Files — use the passed batch verbatim.
		out := make([]collectedFile, 0, len(files))
		for _, f := range files {
			out = append(out, collectedFile{FileId: f.Id, Url: f.FileUrl, Name: f.FileName})
		}
		return out, nil
	}
}

// collectSubtree gathers leaf files under the anchor document's subtree.
func (b *base) collectSubtree(ctx context.Context, txId, anchorId int64) ([]collectedFile, error) {
	docs, err := b.loadTxDocs(ctx, txId)
	if err != nil {
		return nil, err
	}
	return b.gatherFiles(ctx, CollectSubtree(docs, anchorId)), nil
}

// collectBySelector gathers leaf files across the transaction whose tags contain
// every selector tag (AND); can span subtrees.
func (b *base) collectBySelector(ctx context.Context, txId int64, selector *dto.DocumentTemplateTagSelector) ([]collectedFile, error) {
	docs, err := b.loadTxDocs(ctx, txId)
	if err != nil {
		return nil, err
	}
	return b.gatherFiles(ctx, CollectBySelector(docs, selector)), nil
}

// gatherFiles loads each leaf's files, logging (not failing) on a per-document
// load error. Type-agnostic — it returns every file; per-type handling is the
// extract step's job.
func (b *base) gatherFiles(ctx context.Context, leaves []*model.Document) []collectedFile {
	var out []collectedFile
	for _, d := range leaves {
		files, err := b.repos.FileRepo.GetFiles(ctx, d.Id, constants.FileEntityTypeDocument)
		if err != nil {
			b.logger.WithContext(ctx).Errorf("load files for doc %d: %s", d.Id, err)
			continue
		}
		for _, f := range files {
			out = append(out, collectedFile{FileId: f.Id, Url: f.FileUrl, Name: f.FileName})
		}
	}
	return out
}

// loadTxDocs returns all documents attached to the transaction's entity document.
func (b *base) loadTxDocs(ctx context.Context, txId int64) ([]*model.Document, error) {
	entDoc, err := b.repos.EntityDocumentRepo.GetByEntity(ctx, constants.EntityDocumentTypeBorrowerPropertyTransaction, txId)
	if err != nil {
		return nil, err
	}
	if entDoc == nil {
		return nil, errorx.New("entity document not found")
	}
	return b.repos.DocumentRepo.GetAllByEntityDocumentId(ctx, entDoc.Id)
}
