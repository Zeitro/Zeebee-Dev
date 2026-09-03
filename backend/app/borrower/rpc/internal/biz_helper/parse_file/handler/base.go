package handler

import (
	"bytes"
	"context"
	"io"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/rpc/internal/biz_helper/parse_file/file"
	"zeebee-backend/app/borrower/rpc/internal/biz_helper/parse_file/progress"
	pdfclient "zeebee-backend/app/pdf_generator/client"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/llm"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/repo"
)

// extractSystemPrompt forces a single raw JSON object — no fences, no prose.
const extractSystemPrompt = "You extract data from documents. Your entire response MUST be a single raw JSON object and nothing else: start with { and end with }. Do NOT wrap it in markdown code fences (no ```json, no ```). Do NOT add prose, explanations, backticks, or any text before or after the JSON."

const (
	// uploadMaxRetry bounds convert+upload retries (overload/timeout). Backoff
	// doubles per attempt; exhaustion fails the file soft (handler skips it).
	uploadMaxRetry     = 3
	uploadRetryBackoff = 2 * time.Second
)

// pdfConverter converts a doc/docx URL to pdf bytes. Abstracted so base tests
// can inject a fake without a live pdf_generator.
type pdfConverter interface {
	DocToPdf(ctx context.Context, typ, url string, data []byte, ext string) ([]byte, error)
}

// fileStreamResult bundles a fetched file's stream, MIME, and the name downstream
// uploads it under (Word files are renamed to .pdf — see file.PdfDisplayName).
type fileStreamResult struct {
	Reader      io.Reader
	ContentType string
	DisplayName string
}

// base holds the parsing-only deps every handler shares; the coordinator owns
// the run lifecycle.
type base struct {
	repos     *repo.AllRepos
	logger    logx.Logger
	llm       llm.Client
	pdfClient pdfConverter
}

// uploadOne slices PDFs over maxPages (0 = whole) then uploads, retrying the
// convert+upload on transient failures (overload/timeout). Cache hit needs
// provider+MaxPages match (a cap change re-slices); one hold → RMW-safe.
//
// Retry is unconditional (not gated on error type): pdf_generator does not
// distinguish overload from permanent failures in its response code, so any
// convert/upload error is retried up to uploadMaxRetry then fail-soft skipped.
func (b *base) uploadOne(ctx context.Context, rep Reporter, f collectedFile, maxPages int) (progress.FileCheckpoint, error) {
	expected := maxPagesForFile(f.Name, maxPages)
	if cached, ok := rep.Files(ctx)[f.FileId]; ok && cached.Provider == b.llm.Provider() && cached.MaxPages == expected {
		return cached, nil
	}
	var lastErr error
	for attempt := range uploadMaxRetry {
		fc, err := b.convertAndUpload(ctx, rep, f, maxPages, expected)
		if err == nil {
			return fc, nil
		}
		lastErr = err
		if attempt == uploadMaxRetry-1 {
			break
		}
		select {
		case <-ctx.Done():
			return progress.FileCheckpoint{}, ctx.Err()
		case <-time.After(uploadRetryBackoff << attempt):
		}
	}
	return progress.FileCheckpoint{}, lastErr
}

// convertAndUpload fetches/converts one file, slices pdf bytes over maxPages,
// uploads under the display name, and checkpoints the provider file id.
func (b *base) convertAndUpload(ctx context.Context, rep Reporter, f collectedFile, maxPages, expected int) (progress.FileCheckpoint, error) {
	res, err := b.fileStream(ctx, f)
	if err != nil {
		return progress.FileCheckpoint{}, err
	}
	full, err := io.ReadAll(res.Reader)
	if err != nil {
		return progress.FileCheckpoint{}, errorx.Wrap(err, "read file")
	}

	uploadBytes := full
	// Word→pdf and native pdf both stream as application/pdf, so the page cap
	// + pdfium slicing apply uniformly.
	if res.ContentType == "application/pdf" && maxPages > 0 {
		if sliced, serr := file.FirstPagesPDF(bytes.NewReader(full), maxPages); serr != nil {
			return progress.FileCheckpoint{}, errorx.Wrap(serr, "slice pdf")
		} else if sliced != nil {
			uploadBytes = sliced
		}
	}

	provider, fileID, err := b.llm.UploadFile(ctx, bytes.NewReader(uploadBytes), res.DisplayName, res.ContentType)
	if err != nil {
		return progress.FileCheckpoint{}, errorx.Wrap(err, "upload file")
	}
	fc := progress.FileCheckpoint{Provider: provider, ProviderFileID: fileID, MaxPages: expected}
	rep.SaveFile(ctx, f.FileId, fc)
	return fc, nil
}

func maxPagesForFile(name string, cap int) int {
	if file.IsPDF(name) || file.IsWord(name) {
		return cap
	}
	return 0
}

// fileStream fetches one file's bytes + content type + display name by
// extension. Word files convert to pdf via pdf_generator first.
func (b *base) fileStream(ctx context.Context, f collectedFile) (fileStreamResult, error) {
	switch {
	case file.IsPDF(f.Name):
		reader, err := file.DownloadPDF(ctx, f.Url)
		if err != nil {
			return fileStreamResult{}, errorx.Wrap(err, "download pdf")
		}
		return fileStreamResult{Reader: reader, ContentType: "application/pdf", DisplayName: f.Name}, nil
	case file.IsImage(f.Name):
		img, err := file.DownloadImage(ctx, f.Url)
		if err != nil {
			return fileStreamResult{}, errorx.Wrap(err, "download image")
		}
		return fileStreamResult{Reader: bytes.NewReader(*img), ContentType: file.ImageContentType(f.Name), DisplayName: f.Name}, nil
	case file.IsWord(f.Name):
		return b.wordToPdf(ctx, f)
	}
	return fileStreamResult{}, errorx.New("unsupported file type")
}

// wordToPdf converts a Word file to pdf bytes via pdf_generator. Absent client
// (base URL unconfigured) → error so the caller fail-soft skips the file.
func (b *base) wordToPdf(ctx context.Context, f collectedFile) (fileStreamResult, error) {
	if b.pdfClient == nil {
		return fileStreamResult{}, errorx.New("pdf generator unavailable")
	}
	pdf, err := b.pdfClient.DocToPdf(ctx, pdfclient.SourceTypeUrl, f.Url, nil, file.WordExt(f.Name))
	if err != nil {
		return fileStreamResult{}, errorx.Wrap(err, "convert word to pdf")
	}
	return fileStreamResult{Reader: bytes.NewReader(pdf), ContentType: "application/pdf", DisplayName: file.PdfDisplayName(f.Name)}, nil
}

// collectedFile is the flat file list handed from scope resolution to extraction
// (avoids a DB round-trip per file). Type-agnostic; extract dispatches per type.
type collectedFile struct {
	FileId int64  `json:"file_id"`
	Url    string `json:"url"`
	Name   string `json:"name"`
}

// resolveFiles picks in-scope files by rule.Source (SCOPE), independent of
// trigger and file type:
//   - SingleFile/Files → the passed batch (upload batch, or the doc's current files on manual);
//   - Directory → leaf files under the anchor's subtree;
//   - TagFiltered → leaves tagged with every selector tag (AND).
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

func (b *base) collectSubtree(ctx context.Context, txId, anchorId int64) ([]collectedFile, error) {
	docs, err := b.loadTxDocs(ctx, txId)
	if err != nil {
		return nil, err
	}
	return b.gatherFiles(ctx, CollectSubtree(docs, anchorId)), nil
}

// collectBySelector: leaves across the transaction tagged with every selector tag (AND).
func (b *base) collectBySelector(ctx context.Context, txId int64, selector *dto.DocumentTemplateTagSelector) ([]collectedFile, error) {
	docs, err := b.loadTxDocs(ctx, txId)
	if err != nil {
		return nil, err
	}
	return b.gatherFiles(ctx, CollectBySelector(docs, selector)), nil
}

// gatherFiles loads each leaf's files, fail-soft per document. Type-agnostic.
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
