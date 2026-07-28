// Package handler holds per-parser workers (escrow, rpa) + shared
// extract/collect/write-back primitives. The coordinator owns the run lifecycle.
package handler

import (
	"context"
	"encoding/json"

	"zeebee-backend/app/borrower/rpc/internal/biz_helper/parse_file/progress"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
)

// Reporter is a handler's window onto run progress: stage granularity, the
// Extract files checkpoint, and persisted retry. Redis/lifecycle-free; the
// coordinator supplies the Redis-backed impl.
type Reporter interface {
	MarkStage(ctx context.Context, stage, status string)
	SetStageTotal(ctx context.Context, stage string, total int)
	SetStageProcessed(ctx context.Context, stage string, processed int)
	// SetStageState writes a stage-specific blob (Extract: {uploaded,total}).
	SetStageState(ctx context.Context, stage string, state json.RawMessage)
	// StageRetry/SaveStageRetry: persisted retry counter (T009b) — recovery continues counting.
	StageRetry(ctx context.Context, stage string) progress.RetryState
	SaveStageRetry(ctx context.Context, stage string, retry progress.RetryState)
	// SaveFile/Files: Extract resume checkpoint (fileId → provider file id).
	SaveFile(ctx context.Context, fileId int64, fc progress.FileCheckpoint)
	Files(ctx context.Context) map[int64]progress.FileCheckpoint
}

// Handler implements one parser's PURE extraction work — no run lifecycle (the
// coordinator acquires the lock, seeds, heartbeats, finalizes, writes terminal
// Done). Stages is the layout the coordinator seeds; Parse extracts + writes
// back and reports transitions via rep. Trigger is encoded in the rule (Source),
// not a method.
// Adding a parser = write a Handler + register it.
type Handler interface {
	// Stages is the per-parser layout the coordinator seeds. files = the
	// triggering batch (on_upload) or nil (manual) so an on_upload parser can
	// size its stage.
	Stages(files []model.File) []progress.StageState
	Parse(ctx context.Context, txId, docId int64, files []model.File, rule dto.DocumentTemplateFileParseRule, rep Reporter) error
}
