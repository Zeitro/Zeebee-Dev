// Package handler holds the per-parser workers (escrow, rpa) plus the shared
// extract/collect/write-back primitives. The coordinator in package parse_file
// owns the run lifecycle; handlers are pure workers reporting stage progress.
package handler

import (
	"context"
	"encoding/json"

	"zeebee-backend/app/borrower/rpc/internal/biz_helper/parse_file/progress"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
)

// Reporter is a handler's window onto run progress: stage granularity + the
// Extract files checkpoint + persisted retry. Free of Redis/lifecycle knowledge;
// the coordinator supplies the Redis-backed impl.
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

// Handler implements one parser's PURE file-extraction work. It owns no run
// lifecycle — no dedup lock, no goroutine, no Redis writes: the coordinator
// decides entry (routing by parser/source/trigger/rule), acquires the lock,
// seeds job/hold/progress, heartbeats, finalizes, and writes the terminal Done
// state. Stages returns the stage layout the coordinator seeds; Parse extracts
// + writes back and reports stage transitions through rep. Whether it is
// upload- or manual-triggered is encoded in the rule (Source), not a method.
//
// Adding a parser = "write a Handler + register it".
type Handler interface {
	// Stages is the per-parser stage layout the coordinator seeds into progress.
	// files is the triggering batch (on_upload) or nil (manual) so an on_upload
	// parser can size its single stage.
	Stages(files []model.File) []progress.StageState
	Parse(ctx context.Context, txId, docId int64, files []model.File, rule dto.DocumentTemplateFileParseRule, rep Reporter) error
}
