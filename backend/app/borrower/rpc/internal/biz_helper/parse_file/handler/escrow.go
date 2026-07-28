package handler

import (
	"context"
	"encoding/json"
	"fmt"

	"gorm.io/gorm"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/rpc/internal/biz_helper/parse_file/file"
	"zeebee-backend/app/borrower/rpc/internal/biz_helper/parse_file/progress"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/llm"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/repo"
)

const maxEscrowPages = 1

// escrowFields is the extracted escrow row. Empty fields = absent, skipped on write-back.
type escrowFields struct {
	EscrowCompany string `json:"escrow_company"`
	OfficerName   string `json:"officer_name"`
	OfficerEmail  string `json:"officer_email"`
	OfficerPhone  string `json:"officer_phone"`
	EscrowNumber  string `json:"escrow_number"`
}

func parseEscrowFields(content string) (escrowFields, error) {
	var f escrowFields
	if err := json.Unmarshal([]byte(content), &f); err != nil {
		return f, fmt.Errorf("escrow response not json: %s", content)
	}
	return f, nil
}

// escrowPrompt asks for the 5 escrow fields off the attached file.
const escrowPrompt = `Extract these fields from the attached wire/escrow instruction document.
- escrow_company: the escrow/title company name in standard Title Case (e.g. "Mission Escrow, Inc.").
- officer_name: the escrow officer's full name only (e.g. "Renee Lo").
- officer_email: the escrow officer's email address, if printed.
- officer_phone: the escrow officer's phone number as DIGITS ONLY — strip spaces, parentheses, dashes, dots, plus signs.
- escrow_number: the value next to the "Escrow No." label (e.g. "202656-RL").
Only extract what is actually present. If a field is absent, unreadable, or uncertain, use null — never guess.

CRITICAL — output format: your ENTIRE response must be exactly this JSON object and NOTHING else. Start with an open brace and end with a close brace. Do NOT wrap it in markdown code fences. Do NOT add any explanation, reasoning, notes, or text before or after the JSON:
{"escrow_company": <string|null>, "officer_name": <string|null>, "officer_email": <string|null>, "officer_phone": <string|null>, "escrow_number": <string|null>}`

// ApplyEscrowFields copies only non-empty extracted fields (skip-empty merge).
func ApplyEscrowFields(esc *model.PropertyTransactionEscrow, f escrowFields) {
	if f.EscrowCompany != "" {
		esc.EscrowCompany = f.EscrowCompany
	}
	if f.OfficerName != "" {
		esc.OfficerName = f.OfficerName
	}
	if f.OfficerEmail != "" {
		esc.OfficerEmail = f.OfficerEmail
	}
	if f.OfficerPhone != "" {
		esc.OfficerPhone = f.OfficerPhone
	}
	if f.EscrowNumber != "" {
		esc.No = f.EscrowNumber
	}
}

// EscrowHandler serves the escrow parser (on_upload). Pure worker: the
// coordinator hands it the rule + files; this handler uploads each file, runs a
// single consolidated call per file, and skip-empty upserts into the escrow row.
type EscrowHandler struct {
	base
}

func NewEscrowHandler(repos *repo.AllRepos, logger logx.Logger, client llm.Client) *EscrowHandler {
	return &EscrowHandler{base{repos: repos, logger: logger, llm: client}}
}

// Stages returns a single Extract stage; total is seeded from the upload batch.
func (h *EscrowHandler) Stages(files []model.File) []progress.StageState {
	return []progress.StageState{{
		Name: progress.StageExtract, Status: progress.StageStatusRunning,
		Total: len(files),
	}}
}

// Parse resolves in-scope files per rule.Source, then for each PDF/image: upload
// → single consolidated call ([fileId]) → skip-empty upsert. Fail-soft per file.
func (h *EscrowHandler) Parse(ctx context.Context, txId, docId int64, files []model.File, rule dto.DocumentTemplateFileParseRule, rep Reporter) error {
	inScope, err := h.resolveFiles(ctx, txId, docId, files, rule)
	if err != nil {
		return err
	}
	rep.SetStageTotal(ctx, progress.StageExtract, len(inScope))
	rep.SetStageState(ctx, progress.StageExtract, progress.ExtractStateJSON(0, len(inScope)))
	uploaded := 0
	for i, f := range inScope {
		if !file.IsPDF(f.Name) && !file.IsImage(f.Name) {
			rep.SetStageProcessed(ctx, progress.StageExtract, i+1)
			continue // Word/other excluded from scope
		}
		if err := h.parseEscrowFile(ctx, txId, rep, f); err != nil {
			h.logger.WithContext(ctx).WithFields(
				logx.Field("file", f.Name),
				logx.Field("stacks", errorx.GetErrorStack(err)),
			).Errorf("parse escrow file: %s", err)
		}
		uploaded++
		rep.SetStageState(ctx, progress.StageExtract, progress.ExtractStateJSON(uploaded, len(inScope)))
		rep.SetStageProcessed(ctx, progress.StageExtract, i+1)
	}
	return nil
}

// parseEscrowFile uploads one file, runs the consolidated call, skip-empty upserts.
func (h *EscrowHandler) parseEscrowFile(ctx context.Context, txId int64, rep Reporter, f collectedFile) error {
	fc, err := h.uploadOne(ctx, rep, f, maxEscrowPages)
	if err != nil {
		return err
	}
	var fields escrowFields
	if err := progress.RetryDo(ctx, rep, progress.StageExtract, func() error {
		content, usage, e := h.llm.Chat(ctx, extractSystemPrompt, escrowPrompt, llm.WithAttachment(llm.Attachment{
			FileID:      fc.ProviderFileID,
			ContentType: fileContentType(f.Name),
			Filename:    f.Name,
		}))
		if e != nil {
			return e
		}
		h.logger.WithContext(ctx).WithFields(
			logx.Field("biz", "escrow"),
			logx.Field("provider", h.llm.Provider()),
			logx.Field("usage", usage),
			logx.Field("file", f),
		).Info("llm token usage")
		fields, e = parseEscrowFields(content)
		return e
	}); err != nil {
		return errorx.Wrap(err, "extract escrow fields")
	}
	if fields == (escrowFields{}) {
		return nil // nothing extracted; keep existing
	}
	return h.repos.PropertyTransactionEscrowRepo.Transaction(ctx, func(tx *gorm.DB) error {
		return h.upsertEscrowWithTx(ctx, tx, txId, fields)
	})
}

// upsertEscrowWithTx updates the primary (min-id) escrow row, creating one if none.
func (h *EscrowHandler) upsertEscrowWithTx(ctx context.Context, tx *gorm.DB, txId int64, f escrowFields) error {
	rows, err := h.repos.PropertyTransactionEscrowRepo.FindAllWithTx(ctx, tx, map[string]interface{}{
		"borrower_property_transaction_id": txId,
	})
	if err != nil {
		return err
	}
	if len(rows) == 0 {
		escrow := &model.PropertyTransactionEscrow{BorrowerPropertyTransactionId: txId}
		ApplyEscrowFields(escrow, f)
		_, err = h.repos.PropertyTransactionEscrowRepo.CreateWithTx(ctx, tx, escrow, false)
		return err
	}
	escrow := rows[0]
	for _, r := range rows[1:] {
		if r.Id < escrow.Id {
			escrow = r
		}
	}
	ApplyEscrowFields(escrow, f)
	_, err = h.repos.PropertyTransactionEscrowRepo.UpdateOneByIdWithTx(ctx, tx, escrow.Id, escrow)
	return err
}

// fileContentType maps a filename to its MIME type (pdf vs image) for the attachment.
func fileContentType(name string) string {
	if file.IsPDF(name) {
		return "application/pdf"
	}
	return file.ImageContentType(name)
}
