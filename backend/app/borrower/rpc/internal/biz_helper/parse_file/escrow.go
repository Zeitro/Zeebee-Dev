package parse_file

import (
	"context"
	"encoding/json"
	"io"
	"strconv"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"

	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/model"
)

// queuedFile is the per-file payload stored in the on_upload per-document
// serial queue so the drainer can fetch+extract without an extra DB lookup.
type queuedFile struct {
	Id   int64  `json:"id"`
	Url  string `json:"url"`
	Name string `json:"name"`
}

// dispatchEscrowOnUpload enqueues the uploaded PDFs into the per-document
// serial queue and kicks a drainer goroutine if none is running for that doc.
func (p *Parser) dispatchEscrowOnUpload(ctx context.Context, txId, docId int64, files []model.File) {
	sig := escrowOnUploadSig(docId)
	for _, f := range files {
		meta, _ := json.Marshal(queuedFile{Id: f.Id, Url: f.FileUrl, Name: f.FileName})
		if err := p.prog.EnqueueFile(ctx, txId, docId, string(meta)); err != nil {
			p.logger.WithContext(ctx).WithFields(
				logx.Field("doc_id", docId),
				logx.Field("file", f.FileName),
				logx.Field("stacks", errorx.GetErrorStack(err)),
			).Errorf("enqueue escrow parse: %s", err)
		}
	}
	// Try to become the drainer for this doc. If another goroutine already holds
	// the per-doc lock, it will drain the queue we just appended to.
	ok, err := p.prog.AcquireDedupLock(ctx, txId, sig, p.nodeId)
	if err != nil {
		p.logger.WithContext(ctx).Errorf("acquire escrow drain lock: %s", err)
		return
	}
	if !ok {
		return
	}
	p.runBackground(ctx, func() { p.drainEscrowOnUpload(txId, docId, sig) })
}

// drainEscrowOnUpload pops queued files one by one and parses each serially
// until the per-document queue is empty, then releases the per-doc lock so a
// later upload can start the next drainer.
func (p *Parser) drainEscrowOnUpload(txId, docId int64, sig string) {
	// ponytail: background ctx — the upload request has returned; parse outlives it.
	ctx := context.Background()
	defer func() {
		if err := p.prog.ReleaseDedupLock(ctx, txId, sig); err != nil {
			p.logger.WithContext(ctx).Errorf("release escrow drain lock: %s", err)
		}
	}()
	for {
		raw, ok, err := p.prog.DequeueFile(ctx, txId, docId)
		if err != nil {
			p.logger.WithContext(ctx).Errorf("dequeue escrow parse: %s", err)
			return
		}
		if !ok {
			return
		}
		var qf queuedFile
		if json.Unmarshal([]byte(raw), &qf) != nil {
			continue
		}
		if err := p.parseEscrowFile(ctx, txId, qf); err != nil {
			// fail-silent: log + continue to next queued file.
			p.logger.WithContext(ctx).WithFields(
				logx.Field("file", qf.Name),
				logx.Field("stacks", errorx.GetErrorStack(err)),
			).Errorf("parse escrow file: %s", err)
		}
	}
}

// parseEscrowFile downloads + vision-extracts one PDF and writes non-empty
// fields into the primary escrow row (skip-empty merge). Per-file retry wraps
// the extract; a download failure aborts just this file.
func (p *Parser) parseEscrowFile(ctx context.Context, txId int64, qf queuedFile) error {
	reader, err := downloadPDF(ctx, qf.Url)
	if err != nil {
		return errorx.Wrap(err, "download pdf")
	}
	var fields escrowFields
	if err := retryDo(ctx, func() error {
		if _, e := reader.Seek(0, io.SeekStart); e != nil {
			return e
		}
		var extractErr error
		fields, extractErr = extractEscrowFields(ctx, p.llm, reader)
		return extractErr
	}); err != nil {
		return errorx.Wrap(err, "extract escrow fields")
	}
	if fields == (escrowFields{}) {
		return nil // nothing extracted; keep existing
	}
	return p.repos.PropertyTransactionEscrowRepo.Transaction(ctx, func(tx *gorm.DB) error {
		return p.upsertEscrowWithTx(ctx, tx, txId, fields)
	})
}

// upsertEscrowWithTx updates the primary (min-id) escrow row with non-empty
// fields, creating one if none exists (create_logic seeds an empty row, so this
// usually updates the seed).
func (p *Parser) upsertEscrowWithTx(ctx context.Context, tx *gorm.DB, txId int64, f escrowFields) error {
	rows, err := p.repos.PropertyTransactionEscrowRepo.FindAllWithTx(ctx, tx, map[string]interface{}{
		"borrower_property_transaction_id": txId,
	})
	if err != nil {
		return err
	}
	if len(rows) == 0 {
		escrow := &model.PropertyTransactionEscrow{
			BorrowerPropertyTransactionId: txId,
		}
		applyEscrowFields(escrow, f)
		_, err = p.repos.PropertyTransactionEscrowRepo.CreateWithTx(ctx, tx, escrow, false)
		return err
	}
	escrow := rows[0]
	for _, r := range rows[1:] {
		if r.Id < escrow.Id {
			escrow = r
		}
	}
	applyEscrowFields(escrow, f)
	_, err = p.repos.PropertyTransactionEscrowRepo.UpdateOneByIdWithTx(ctx, tx, escrow.Id, escrow)
	return err
}

func escrowOnUploadSig(docId int64) string {
	return model.FileParseParserEscrowInfo + "|" + model.FileParseModeOnUpload + "|" + strconv.FormatInt(docId, 10)
}
