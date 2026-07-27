package handler

import (
	"bytes"
	"context"
	"io"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/rpc/internal/biz_helper/parse_file/file"
	"zeebee-backend/app/borrower/rpc/internal/biz_helper/parse_file/progress"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/llm"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/repo"
)

// extractSystemPrompt forces a single raw JSON object — no fences, no prose.
const extractSystemPrompt = "You extract data from documents. Your entire response MUST be a single raw JSON object and nothing else: start with { and end with }. Do NOT wrap it in markdown code fences (no ```json, no ```). Do NOT add prose, explanations, backticks, or any text before or after the JSON."

// base holds the parsing-only deps every handler shares. No progress manager,
// no nodeId — the coordinator owns the run lifecycle.
type base struct {
	repos  *repo.AllRepos
	logger logx.Logger
	llm    llm.Client
}

// uploadOne slices PDFs over maxPages (0 = whole file) before upload. Cache hit
// needs provider+MaxPages match so a cap change re-slices; one hold → RMW-safe.
func (b *base) uploadOne(ctx context.Context, rep Reporter, f collectedFile, maxPages int) (progress.FileCheckpoint, error) {
	expected := maxPagesForFile(f.Name, maxPages)
	if cached, ok := rep.Files(ctx)[f.FileId]; ok && cached.Provider == b.llm.Provider() && cached.MaxPages == expected {
		return cached, nil
	}
	reader, contentType, err := b.fileStream(ctx, f)
	if err != nil {
		return progress.FileCheckpoint{}, err
	}
	full, err := io.ReadAll(reader)
	if err != nil {
		return progress.FileCheckpoint{}, errorx.Wrap(err, "read file")
	}

	uploadBytes := full
	if file.IsPDF(f.Name) && maxPages > 0 {
		if sliced, serr := file.FirstPagesPDF(bytes.NewReader(full), maxPages); serr != nil {
			return progress.FileCheckpoint{}, errorx.Wrap(serr, "slice pdf")
		} else if sliced != nil {
			uploadBytes = sliced
		}
	}

	provider, fileID, err := b.llm.UploadFile(ctx, bytes.NewReader(uploadBytes), f.Name, contentType)
	if err != nil {
		return progress.FileCheckpoint{}, errorx.Wrap(err, "upload file")
	}
	fc := progress.FileCheckpoint{Provider: provider, ProviderFileID: fileID, MaxPages: expected}
	rep.SaveFile(ctx, f.FileId, fc)
	return fc, nil
}

// maxPagesForFile: cap for PDFs, 0 for images (cap-agnostic).
func maxPagesForFile(name string, cap int) int {
	if file.IsPDF(name) {
		return cap
	}
	return 0
}

// fileStream fetches one file's bytes + content type by extension (pdf/image).
// Word/other types are excluded from parsing scope.
func (b *base) fileStream(ctx context.Context, f collectedFile) (io.Reader, string, error) {
	switch {
	case file.IsPDF(f.Name):
		reader, err := file.DownloadPDF(ctx, f.Url)
		if err != nil {
			return nil, "", errorx.Wrap(err, "download pdf")
		}
		return reader, "application/pdf", nil
	case file.IsImage(f.Name):
		img, err := file.DownloadImage(ctx, f.Url)
		if err != nil {
			return nil, "", errorx.Wrap(err, "download image")
		}
		return bytes.NewReader(*img), file.ImageContentType(f.Name), nil
	}
	return nil, "", errorx.New("unsupported file type")
}

// collectedFile is the flat file list handed from scope resolution to extraction
// so extraction needs no extra DB round-trip per file. Type-agnostic — PDF today,
// image/doc later; the handler's extract step dispatches per type.
type collectedFile struct {
	FileId int64
	Url    string
	Name   string
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
