package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

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

// rpaFields: 5 contract fields. Dates MM/DD/YYYY ("" absent); amounts float (0 absent).
type rpaFields struct {
	PurchaseAmount         float64 `json:"purchase_amount"`
	CoeDate                string  `json:"coe_date"`
	EmdAmount              float64 `json:"emd_amount"`
	EmdDueDate             string  `json:"emd_due_date"`
	ContractDate           string  `json:"contract_date"`
	ContingencyRemovalDate string  `json:"contingency_removal_date"`
}

// hasAny reports whether any field carries a real value. All-absent → no write-back.
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

// rpaConsolidatedPrompt: read ALL attached documents and arbitrate across files
// in this single call — pick the authoritative value on conflict, null if undeterminable.
const rpaConsolidatedPrompt = `You are given one or more contract documents (PDF/image). Read ALL of them and extract these fields, arbitrating across files when they disagree:
- purchase_amount: purchase price as a JSON number (float). Strip thousands separators/currency. 0 if absent.
- emd_amount: earnest money deposit as a JSON number (float). 0 if absent.
- contract_date: MM/DD/YYYY, or null if absent.
- coe_date: closing of escrow date as MM/DD/YYYY, or null if absent.
- emd_due_date: earnest money deposit due date as MM/DD/YYYY, or null if absent.
- contingency_removal_date: contingency removal date as MM/DD/YYYY, or null if absent.
When files conflict, choose the most authoritative (e.g. signed contract over a draft); if you cannot determine which is correct, return null for that field — never guess. Output every date as MM/DD/YYYY regardless of source format.

CRITICAL — output format: your ENTIRE response must be exactly this JSON object and NOTHING else. Start with { and end with }. Do NOT wrap it in markdown code fences. Do NOT add any explanation, reasoning, notes, or text before or after the JSON:
{"purchase_amount": <number>, "coe_date": <string|null>, "emd_amount": <number>, "emd_due_date": <string|null>, "contingency_removal_date": <string|null>, "contract_date": <string|null>}`

// ---- collect (in-scope leaves by rule.Source) ----

// CollectSubtree returns Leaf documents in the subtree rooted at anchorId.
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

// CollectBySelector returns Leaf documents whose tags contain every selector tag (AND).
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

// applyRPAFields copies only real values (non-zero amount, parseable date).
func applyRPAFields(tx *model.BorrowerPropertyTransaction, f rpaFields) {
	if f.PurchaseAmount > 0 {
		tx.PurchaseAmount = f.PurchaseAmount
	}
	if f.EmdAmount > 0 {
		tx.EmdAmount = f.EmdAmount
	}
	if t, ok := parseDate(f.CoeDate); ok {
		tx.CoeDate = t
	}
	if t, ok := parseDate(f.EmdDueDate); ok {
		tx.EmdDueDate = t
	}
	if t, ok := parseDate(f.ContingencyRemovalDate); ok {
		tx.ContingencyRemovalDate = t
	}
	if t, ok := parseDate(f.ContractDate); ok {
		tx.ContractDate = t
	}
}

// parseDate parses MM/DD/YYYY (rpa LLM output contract). Empty/garbage → not ok.
func parseDate(s string) (*sql.NullTime, bool) {
	if s == "" {
		return nil, false
	}
	t, err := time.Parse("01/02/2006", s)
	if err != nil {
		return nil, false
	}
	return &sql.NullTime{Time: t, Valid: true}, true
}

// ---- handler ----

// RPAHandler serves the rpa parser: manual trigger runs Collect → Extract (upload
// all → single consolidated call with in-call arbitration) → Writeback. Pure worker.
type RPAHandler struct {
	base
}

func NewRPAHandler(repos *repo.AllRepos, logger logx.Logger, client llm.Client) *RPAHandler {
	return &RPAHandler{base{repos: repos, logger: logger, llm: client}}
}

// Stages returns the 3-stage pipeline (Collect → Extract → Writeback).
func (h *RPAHandler) Stages(_ []model.File) []progress.StageState {
	return []progress.StageState{
		{Name: progress.StageCollect, Status: progress.StageStatusPending},
		{Name: progress.StageExtract, Status: progress.StageStatusPending},
		{Name: progress.StageWriteback, Status: progress.StageStatusPending},
	}
}

// Parse runs the staged pipeline; each error aborts the run (coordinator writes Failed).
func (h *RPAHandler) Parse(ctx context.Context, txId, documentId int64, files []model.File, rule dto.DocumentTemplateFileParseRule, rep Reporter) error {
	inScope, err := h.collectStage(ctx, rep, txId, documentId, files, rule)
	if err != nil {
		return err
	}
	fields, err := h.extractStage(ctx, rep, inScope)
	if err != nil {
		return err
	}
	return h.writebackStage(ctx, rep, txId, fields)
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

// extractStage uploads every in-scope PDF/image (fail-soft: a bad file logs and
// is dropped), then issues ONE consolidated call over all uploaded ids — the LLM
// reads every file and arbitrates conflicts inside that call. Resume reuses
// cached ids (provider match) so a restarted run does not re-upload.
func (h *RPAHandler) extractStage(ctx context.Context, rep Reporter, files []collectedFile) (rpaFields, error) {
	rep.MarkStage(ctx, progress.StageExtract, progress.StageStatusRunning)
	rep.SetStageTotal(ctx, progress.StageExtract, len(files))
	rep.SetStageState(ctx, progress.StageExtract, progress.ExtractStateJSON(0, len(files)))

	type uploaded struct{ id, contentType, name string }
	atts := make([]uploaded, 0, len(files))
	uploadedCount := 0
	for _, f := range files {
		if !file.IsPDF(f.Name) && !file.IsImage(f.Name) {
			continue // Word/other excluded
		}
		fc, err := h.uploadOne(ctx, rep, f, file.MaxRPAPages)
		if err != nil {
			h.logger.WithContext(ctx).WithFields(
				logx.Field("file", f.Name),
				logx.Field("stacks", errorx.GetErrorStack(err)),
			).Errorf("upload rpa file: %s", err)
			continue
		}
		atts = append(atts, uploaded{id: fc.ProviderFileID, contentType: fileContentType(f.Name), name: f.Name})
		uploadedCount++
		rep.SetStageState(ctx, progress.StageExtract, progress.ExtractStateJSON(uploadedCount, len(files)))
	}

	if len(atts) == 0 {
		rep.MarkStage(ctx, progress.StageExtract, progress.StageStatusDone)
		return rpaFields{}, nil
	}
	var fields rpaFields
	err := progress.RetryDo(ctx, rep, progress.StageExtract, func() error {
		opts := make([]llm.ChatOption, 0, len(atts))
		for _, a := range atts {
			opts = append(opts, llm.WithAttachment(llm.Attachment{FileID: a.id, ContentType: a.contentType, Filename: a.name}))
		}
		content, e := h.llm.Chat(ctx, extractSystemPrompt, rpaConsolidatedPrompt, opts...)
		if e != nil {
			return e
		}
		fields, e = parseRPAFields(content)
		return e
	})
	rep.MarkStage(ctx, progress.StageExtract, progress.StageStatusDone)
	if err != nil {
		return rpaFields{}, errorx.Wrap(err, "consolidated rpa extract")
	}
	return fields, nil
}

// writebackStage applies the consolidated fields (authoritative replace) into the
// transaction row once. Empty/zero values preserve any existing manual value.
func (h *RPAHandler) writebackStage(ctx context.Context, rep Reporter, txId int64, fields rpaFields) error {
	rep.MarkStage(ctx, progress.StageWriteback, progress.StageStatusRunning)
	rep.SetStageTotal(ctx, progress.StageWriteback, 1)
	if fields.hasAny() {
		tx, err := h.repos.BorrowerPropertyTransactionRepo.GetOne(ctx, txId)
		if err != nil {
			return errorx.Wrap(err, "load transaction")
		}
		if tx == nil {
			return errorx.New("transaction not found")
		}
		applyRPAFields(tx, fields)
		if _, err := h.repos.BorrowerPropertyTransactionRepo.UpdateOne(ctx, txId, tx); err != nil {
			return errorx.Wrap(err, "update transaction")
		}
	}
	rep.SetStageProcessed(ctx, progress.StageWriteback, 1)
	rep.MarkStage(ctx, progress.StageWriteback, progress.StageStatusDone)
	return nil
}
