package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/klippa-app/go-pdfium"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"

	"zeebee-backend/app/borrower/rpc/internal/biz_helper/parse_file/file"
	"zeebee-backend/app/borrower/rpc/internal/biz_helper/parse_file/progress"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/llm"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/repo"
)

// ---- escrow field extraction (vision path, first page) ----

// escrowFields is the vision-extracted escrow row. Empty fields mean "absent"
// and are skipped on write-back (skip-empty merge).
type escrowFields struct {
	EscrowCompany string `json:"escrow_company"`
	OfficerName   string `json:"officer_name"`
	OfficerEmail  string `json:"officer_email"`
	OfficerPhone  string `json:"officer_phone"`
	EscrowNumber  string `json:"escrow_number"`
}

func (e escrowFields) complete() bool {
	return e.EscrowCompany != "" && e.OfficerName != "" && e.OfficerEmail != "" &&
		e.OfficerPhone != "" && e.EscrowNumber != ""
}

func parseEscrowFields(content string) (escrowFields, error) {
	var f escrowFields
	if err := json.Unmarshal([]byte(content), &f); err != nil {
		return f, fmt.Errorf("escrow response not json: %s", content)
	}
	return f, nil
}

// escrowPrompt asks for the 5 escrow fields off the rendered first page.
func escrowPrompt() string {
	return `Extract these fields from the attached first page of a wire/escrow instruction document.
- escrow_company: the escrow/title company name in standard Title Case (e.g. "Mission Escrow, Inc.").
- officer_name: the escrow officer's full name only (e.g. "Renee Lo").
- officer_email: the escrow officer's email address, if printed on the page.
- officer_phone: the escrow officer's phone number as DIGITS ONLY — strip spaces, parentheses, dashes, dots, plus signs (e.g. "(408) 513-6005" becomes "4085136005").
- escrow_number: the value next to the "Escrow No." label (e.g. "202656-RL").
Only extract what is actually printed on the page. If a field is absent, unreadable, or uncertain, use null — never guess.
Respond as JSON: {"escrow_company": <string|null>, "officer_name": <string|null>, "officer_email": <string|null>, "officer_phone": <string|null>, "escrow_number": <string|null>}`
}

// ExtractEscrowFields vision-reads escrow fields off the rendered first page.
func ExtractEscrowFields(ctx context.Context, client llm.Client, r io.ReadSeeker) (fields escrowFields, err error) {
	err = file.WithPdfInstance(func(inst pdfium.Pdfium) error {
		doc, cleanup, e := file.OpenPDF(inst, r)
		if e != nil {
			return e
		}

		defer cleanup()
		png, e := file.RenderPagePNG(inst, doc, 0)
		if e != nil {
			return e
		}
		content, e := client.VisionJSON(ctx, extractSystemPrompt, escrowPrompt(), png)
		if e != nil {
			return e
		}
		fields, e = parseEscrowFields(content)
		return e
	})
	return fields, err
}

// ApplyEscrowFields copies only non-empty extracted fields (skip-empty merge),
// preserving any real value already on the row.
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

// ---- handler ----

// EscrowHandler serves the escrow parser. It is a pure worker: the coordinator
// hands it the rule + files in scope, and this handler resolves scope, then
// downloads + vision-extracts + skip-empty upserts each file. No locks, no
// goroutines, no Redis writes — the coordinator owns all of that.
type EscrowHandler struct {
	base
}

func NewEscrowHandler(repos *repo.AllRepos, logger logx.Logger, client llm.Client) *EscrowHandler {
	return &EscrowHandler{base{repos: repos, logger: logger, llm: client}}
}

// Stages is the single Extract stage. Total is seeded from the passed batch
// (on_upload SingleFile/Files); for collected scopes the handler sets the real
// total via the reporter once it has resolved the file set.
func (h *EscrowHandler) Stages(files []model.File) []progress.StageState {
	return []progress.StageState{{
		Name: progress.StageExtract, Status: progress.StageStatusRunning,
		Total: len(files), Retry: progress.RetryStateFor(0),
	}}
}

// Parse resolves the in-scope files per rule.Source, then extracts + writes
// back each (fail-soft per file: one bad file logs and continues). File-type
// filtering lives here in the extract step (PDF today), not in scope resolution
// — image/doc support lands later without changing the scope contract.
func (h *EscrowHandler) Parse(ctx context.Context, txId, docId int64, files []model.File, rule dto.DocumentTemplateFileParseRule, rep Reporter) error {
	inScope, err := h.resolveFiles(ctx, txId, docId, files, rule)
	if err != nil {
		return err
	}
	rep.SetStageTotal(ctx, progress.StageExtract, len(inScope))
	for i, f := range inScope {
		if !file.IsPDF(f.Name) {
			continue // ponytail: PDF-only for now; image/doc handled in the file pkg later
		}
		if err := h.parseEscrowFile(ctx, txId, f); err != nil {
			h.logger.WithContext(ctx).WithFields(
				logx.Field("file", f.Name),
				logx.Field("stacks", errorx.GetErrorStack(err)),
			).Errorf("parse escrow file: %s", err)
		}
		rep.SetStageProcessed(ctx, progress.StageExtract, i+1)
	}
	return nil
}

// parseEscrowFile downloads + vision-extracts one PDF and writes non-empty
// fields into the primary escrow row (skip-empty merge). Per-file retry wraps
// the extract; a download failure aborts just this file.
func (h *EscrowHandler) parseEscrowFile(ctx context.Context, txId int64, f collectedFile) error {
	reader, err := file.DownloadPDF(ctx, f.Url)
	if err != nil {
		return errorx.Wrap(err, "download pdf")
	}
	var fields escrowFields
	if err := progress.RetryDo(ctx, func() error {
		if _, e := reader.Seek(0, io.SeekStart); e != nil {
			return e
		}
		var extractErr error
		fields, extractErr = ExtractEscrowFields(ctx, h.llm, reader)
		return extractErr
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

// upsertEscrowWithTx updates the primary (min-id) escrow row with non-empty
// fields, creating one if none exists (create_logic seeds an empty row, so this
// usually updates the seed).
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
