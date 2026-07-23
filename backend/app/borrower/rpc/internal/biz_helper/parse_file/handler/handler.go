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

// Reporter is a handler's narrow window onto run progress: it reports stage
// granularity while staying free of any Redis/lifecycle knowledge. The
// coordinator (parent parse_file package) supplies the Redis-backed impl.
type Reporter interface {
	MarkStage(ctx context.Context, stage, status string)
	SetStageTotal(ctx context.Context, stage string, total int)
	SetStageProcessed(ctx context.Context, stage string, processed int)
	// StageOutput reads the stage's resume checkpoint — the marshaled
	// map[int64] of already-extracted per-file results. A recovered run recalls
	// it so it skips finished files without losing their results.
	StageOutput(ctx context.Context, stage string) json.RawMessage
	// SaveStageOutput writes the checkpoint after a file is processed.
	SaveStageOutput(ctx context.Context, stage string, output json.RawMessage)
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
