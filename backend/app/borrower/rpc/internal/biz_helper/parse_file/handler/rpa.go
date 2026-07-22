package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/klippa-app/go-pdfium"
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

// ---- rpa field extraction (text path, leading N pages) ----

// rpaFields is the 5 contract fields extracted from the RPA/contract text.
// Dates are ISO yyyy-mm-dd ("" when absent); amounts are floats (0 when absent).
type rpaFields struct {
	PurchaseAmount         float64 `json:"purchase_amount"`
	CoeDate                string  `json:"coe_date"`
	EmdAmount              float64 `json:"emd_amount"`
	EmdDueDate             string  `json:"emd_due_date"`
	ContingencyRemovalDate string  `json:"contingency_removal_date"`
}

// hasAny reports whether any field carries a real value (non-zero amount or
// non-empty date). All-absent envelopes are dropped before write-back.
func (f rpaFields) hasAny() bool {
	return f.PurchaseAmount > 0 || f.CoeDate != "" || f.EmdAmount > 0 ||
		f.EmdDueDate != "" || f.ContingencyRemovalDate != ""
}

func parseRPAFields(content string) (rpaFields, error) {
	var f rpaFields
	if err := json.Unmarshal([]byte(content), &f); err != nil {
		return f, fmt.Errorf("rpa response not json: %s", content)
	}
	return f, nil
}

// rpaEnvelope is one file's extraction result handed off from Extract to
// Writeback. Fields with no real value are dropped (hasAny) before handoff.
type rpaEnvelope struct {
	FileId int64     `json:"file_id"`
	Fields rpaFields `json:"fields"`
}

// rpaPrompt asks for the 5 contract fields from the leading pages' text.
func rpaPrompt(text string) string {
	return fmt.Sprintf(`Extract these fields from the following contract document text.
- purchase_amount: the purchase price as a JSON number (float). Strip thousands separators (commas), currency symbols, and whitespace. Examples: "$1,750,000.00" -> 1750000. Use 0 if absent.
- coe_date: close-of-escrow date as ISO yyyy-mm-dd (e.g. "2026-09-15"). Use null if absent.
- emd_amount: earnest money deposit amount as a JSON number (float). Use 0 if absent.
- emd_due_date: EMD due date as ISO yyyy-mm-dd. Use null if absent.
- contingency_removal_date: contingency removal date as ISO yyyy-mm-dd. Use null if absent.
Respond as JSON: {"purchase_amount": <number>, "coe_date": <string|null>, "emd_amount": <number>, "emd_due_date": <string|null>, "contingency_removal_date": <string|null>}

Document text:
%s`, text)
}

// ExtractRPAFields reads the leading MaxRPAPages of text and extracts the 5
// contract fields via the text LLM (the only LLM call in the rpa pipeline now
// that arbitration is gone).
func ExtractRPAFields(ctx context.Context, client llm.Client, fileId int64, r io.ReadSeeker) (rpaEnvelope, error) {
	var env rpaEnvelope
	env.FileId = fileId
	err := file.WithPdfInstance(func(inst pdfium.Pdfium) error {
		doc, cleanup, e := file.OpenPDF(inst, r)
		if e != nil {
			return e
		}
		defer cleanup()
		text, e := file.ReadFirstNPages(inst, doc, file.MaxRPAPages)
		if e != nil {
			return e
		}
		content, e := client.ChatJSON(ctx, extractSystemPrompt, rpaPrompt(text))
		if e != nil {
			return e
		}
		env.Fields, e = parseRPAFields(content)
		return e
	})
	return env, err
}

// ---- collect (pick the in-scope leaves by rule.Source) ----

// CollectSubtree returns the Leaf documents in the subtree rooted at anchorId.
// docs must contain the full document set for the transaction so parent/child
// links walk via ParentId. Files attach to Leaf documents, so only leaves parse.
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
// tag in selector.Tags (AND). Empty selector tags matches nothing (avoid an
// accidental full-scope parse).
func CollectBySelector(docs []*model.Document, selector *dto.DocumentTemplateTagSelector) []*model.Document {
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

// ---- write-back ----

// applyRPAFields copies only real values (non-zero amount, parseable date),
// preserving any existing manual or prior value (skip-empty merge, FR-014).
func applyRPAFields(tx *model.BorrowerPropertyTransaction, f rpaFields) {
	if f.PurchaseAmount > 0 {
		tx.PurchaseAmount = f.PurchaseAmount
	}
	if f.EmdAmount > 0 {
		tx.EmdAmount = f.EmdAmount
	}
	if t, ok := parseISODate(f.CoeDate); ok {
		tx.CoeDate = t
	}
	if t, ok := parseISODate(f.EmdDueDate); ok {
		tx.EmdDueDate = t
	}
	if t, ok := parseISODate(f.ContingencyRemovalDate); ok {
		tx.ContingencyRemovalDate = t
	}
}

func parseISODate(s string) (*sql.NullTime, bool) {
	if s == "" {
		return nil, false
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return nil, false
	}
	return &sql.NullTime{Time: t, Valid: true}, true
}

// ---- handler ----

// RPAHandler serves the rpa parser: a manual trigger runs Collect → Extract →
// Arbitrate → Writeback over the document subtree (manual_tree) or the
// tag-selected set (manual_select). Pure worker — the coordinator owns the
// lifecycle; this reports stage granularity through rep.
type RPAHandler struct {
	base
}

func NewRPAHandler(repos *repo.AllRepos, logger logx.Logger, client llm.Client) *RPAHandler {
	return &RPAHandler{base{repos: repos, logger: logger, llm: client}}
}

// Stages is the 3-stage pipeline (Collect → Extract → Writeback); files is
// unused (rpa collects its own scope per rule.Source).
func (h *RPAHandler) Stages(_ []model.File) []progress.StageState {
	retry := progress.RetryStateFor(0)
	return []progress.StageState{
		{Name: progress.StageCollect, Status: progress.StageStatusPending, Retry: retry},
		{Name: progress.StageExtract, Status: progress.StageStatusPending, Retry: retry},
		{Name: progress.StageWriteback, Status: progress.StageStatusPending, Retry: retry},
	}
}

// Parse runs the staged pipeline; each error aborts the run (the coordinator
// writes the terminal Failed Done state). Collect resolves scope per rule.Source
// (passed batch / subtree / tag set) through the shared base.resolveFiles, so
// rpa reuses the same scope contract as every other handler.
func (h *RPAHandler) Parse(ctx context.Context, txId, documentId int64, files []model.File, rule dto.DocumentTemplateFileParseRule, rep Reporter) error {
	inScope, err := h.collectStage(ctx, rep, txId, documentId, files, rule)
	if err != nil {
		return err
	}
	envelopes := h.extractStage(ctx, rep, inScope)
	return h.writebackStage(ctx, rep, txId, envelopes)
}

func (h *RPAHandler) collectStage(ctx context.Context, rep Reporter, txId, documentId int64, files []model.File, rule dto.DocumentTemplateFileParseRule) ([]collectedFile, error) {
	rep.MarkStage(ctx, progress.StageCollect, progress.StageStatusRunning)
	out, err := h.resolveFiles(ctx, txId, documentId, files, rule)
	if err != nil {
		return nil, errorx.Wrap(err, "collect files")
	}
	rep.SetStageTotal(ctx, progress.StageCollect, len(out))
	rep.MarkStage(ctx, progress.StageCollect, progress.StageStatusDone)
	return out, nil
}

// extractStage downloads + text-extracts each PDF (fail-soft: one bad file
// logs and is skipped). It never returns an error — a per-file failure must
// not abort the whole parse, so the write-back proceeds with whatever real
// fields were collected.
func (h *RPAHandler) extractStage(ctx context.Context, rep Reporter, files []collectedFile) []rpaEnvelope {
	rep.MarkStage(ctx, progress.StageExtract, progress.StageStatusRunning)
	rep.SetStageTotal(ctx, progress.StageExtract, len(files))
	envelopes := make([]rpaEnvelope, 0, len(files))
	for i, f := range files {
		if !file.IsPDF(f.Name) {
			rep.SetStageProcessed(ctx, progress.StageExtract, i+1)
			continue // ponytail: PDF text path only for now; image/doc land in the file pkg later
		}
		env, err := h.extractOneFile(ctx, f)
		if err != nil {
			// fail-soft: one bad PDF must not abort the whole parse.
			h.logger.WithContext(ctx).WithFields(
				logx.Field("file", f.Name),
				logx.Field("stacks", errorx.GetErrorStack(err)),
			).Errorf("extract rpa file: %s", err)
		} else if env.Fields.hasAny() {
			envelopes = append(envelopes, env)
		}
		rep.SetStageProcessed(ctx, progress.StageExtract, i+1)
	}
	rep.MarkStage(ctx, progress.StageExtract, progress.StageStatusDone)
	return envelopes
}

func (h *RPAHandler) extractOneFile(ctx context.Context, f collectedFile) (rpaEnvelope, error) {
	var env rpaEnvelope
	reader, err := file.DownloadPDF(ctx, f.Url)
	if err != nil {
		return env, errorx.Wrap(err, "download pdf")
	}
	if err := progress.RetryDo(ctx, func() error {
		if _, e := reader.Seek(0, io.SeekStart); e != nil {
			return e
		}
		var extractErr error
		env, extractErr = ExtractRPAFields(ctx, h.llm, f.FileId, reader)
		return extractErr
	}); err != nil {
		return env, errorx.Wrap(err, "extract rpa fields")
	}
	return env, nil
}

// writebackStage applies every envelope's real fields into the transaction row
// in collection order via skip-empty merge — the first real value per field
// wins (applyRPAFields skips empty/zero), so a conflict across files resolves
// deterministically without an arbitration step. The row is loaded + written
// once.
func (h *RPAHandler) writebackStage(ctx context.Context, rep Reporter, txId int64, envelopes []rpaEnvelope) error {
	rep.MarkStage(ctx, progress.StageWriteback, progress.StageStatusRunning)
	rep.SetStageTotal(ctx, progress.StageWriteback, 1)
	tx, err := h.repos.BorrowerPropertyTransactionRepo.GetOne(ctx, txId)
	if err != nil {
		return errorx.Wrap(err, "load transaction")
	}
	if tx == nil {
		return errorx.New("transaction not found")
	}
	for _, env := range envelopes {
		applyRPAFields(tx, env.Fields)
	}
	if _, err := h.repos.BorrowerPropertyTransactionRepo.UpdateOne(ctx, txId, tx); err != nil {
		return errorx.Wrap(err, "update transaction")
	}
	rep.SetStageProcessed(ctx, progress.StageWriteback, 1)
	rep.MarkStage(ctx, progress.StageWriteback, progress.StageStatusDone)
	return nil
}
