package parse_file

import (
	"context"
	"encoding/json"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"zeebee-backend/app/borrower/rpc/internal/biz_helper/parse_file/handler"
	"zeebee-backend/app/borrower/rpc/internal/biz_helper/parse_file/progress"
	"zeebee-backend/common/llm"
	"zeebee-backend/common/redisx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
)

// These tests cover the coordinator lifecycle in parser.go (Dispatch's
// start/run/seed/finalize + dedup): staged progress, terminal Done, key
// cleanup, and the on_upload-join vs manual-409 divergence. They drive the
// parser with a fake handler over miniredis — start/run never touch the repo
// layer (only handlers do), so no DB/S3/PDF fixture is needed. Per-parser
// extract/arbitrate behavior is covered in package handler.

// stubLLM removed: tests use llm.MockClient (implements the full Client iface).

// fakeHandler is a handler.Handler that reports canned stage progress and can
// block (to hold a run in flight for dedup tests) or return an error. It owns
// no lifecycle and no repos — exactly the contract the coordinator expects.
type fakeHandler struct {
	block chan struct{}
	err   error
}

func newFakeHandler() *fakeHandler { return &fakeHandler{block: make(chan struct{})} }

func (f *fakeHandler) Stages([]model.File) []progress.StageState {
	return []progress.StageState{
		{Name: progress.StageCollect, Status: progress.StageStatusPending},
		{Name: progress.StageExtract, Status: progress.StageStatusPending},
	}
}

func (f *fakeHandler) Parse(ctx context.Context, _, _ int64, _ []model.File, _ dto.DocumentTemplateFileParseRule, rep handler.Reporter) error {
	<-f.block // stay in flight until the test releases us
	if f.err != nil {
		return f.err
	}
	rep.MarkStage(ctx, progress.StageExtract, progress.StageStatusRunning)
	rep.SetStageTotal(ctx, progress.StageExtract, 2)
	rep.SetStageProcessed(ctx, progress.StageExtract, 2)
	rep.MarkStage(ctx, progress.StageExtract, progress.StageStatusDone)
	return nil
}

func newTestParser(t *testing.T) (*Parser, *progress.ProgressManager) {
	t.Helper()
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("miniredis: %v", err)
	}
	t.Cleanup(mr.Close)
	rdb, err := redisx.GetRedis(redisx.Opt{Host: mr.Addr(), Type: "node"})
	if err != nil {
		t.Fatalf("redisx: %v", err)
	}
	logger := logx.WithContext(context.Background())
	pm := progress.NewProgressManager(rdb, logger)
	return NewParser(nil, &llm.MockClient{}, pm, nil, logger), pm
}

func rpaDirectoryRule() dto.DocumentTemplateFileParseRule {
	return dto.DocumentTemplateFileParseRule{
		Parser: constants.DocumentTemplateFileParseParserRPA,
		Source: constants.DocumentTemplateFileParseSourceDirectory,
	}
}

// waitDone polls the terminal record; miniredis is in-process so runs resolve
// within milliseconds once the handler unblocks.
func waitDone(t *testing.T, pm *progress.ProgressManager, requestId string) *progress.DoneState {
	t.Helper()
	ctx := context.Background()
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		if d, _ := pm.GetDone(ctx, requestId); d != nil {
			return d
		}
		time.Sleep(2 * time.Millisecond)
	}
	t.Fatalf("done record for %s never written", requestId)
	return nil
}

// TestStartSeedsRunningProgress verifies the synchronous seed: by the time
// start returns, a live Running progress record exists for the requestId.
func TestStartSeedsRunningProgress(t *testing.T) {
	p, pm := newTestParser(t)
	fake := newFakeHandler()
	defer close(fake.block)

	requestId, err := p.start(context.Background(), 1, 1, rpaDirectoryRule(), fake, nil, false)
	if err != nil {
		t.Fatalf("start: %v", err)
	}
	prog, _ := pm.GetProgress(context.Background(), requestId)
	if prog == nil || prog.Status != progress.StatusRunning {
		t.Fatalf("progress = %+v, want live Running right after start", prog)
	}
}

// TestRunWritesFinishedAndCleansKeys drives a full success run and asserts the
// terminal Done is Finished and the run keys (progress/hold) are reclaimed.
func TestRunWritesFinishedAndCleansKeys(t *testing.T) {
	p, pm := newTestParser(t)
	fake := newFakeHandler()
	ctx := context.Background()

	requestId, err := p.start(ctx, 1, 1, rpaDirectoryRule(), fake, nil, false)
	if err != nil {
		t.Fatalf("start: %v", err)
	}
	close(fake.block)

	done := waitDone(t, pm, requestId)
	if done.Status != progress.StatusFinished {
		t.Fatalf("done = %+v, want Finished", done)
	}
	if prog, _ := pm.GetProgress(ctx, requestId); prog != nil {
		t.Fatalf("progress key must be cleared after WriteDone, got %+v", prog)
	}
	// finalize releases the hold — a fresh SetNX must succeed.
	if ok, _ := pm.AcquireHold(ctx, requestId, "after"); !ok {
		t.Fatal("hold must be released after run so a reclaim SetNX succeeds")
	}
}

// TestRunFailureWritesFailed verifies a handler error surfaces as a terminal
// Failed done record carrying the error message.
func TestRunFailureWritesFailed(t *testing.T) {
	p, pm := newTestParser(t)
	fake := newFakeHandler()
	fake.err = errors.New("boom")
	ctx := context.Background()

	requestId, err := p.start(ctx, 1, 1, rpaDirectoryRule(), fake, nil, false)
	if err != nil {
		t.Fatalf("start: %v", err)
	}
	close(fake.block)

	done := waitDone(t, pm, requestId)
	if done.Status != progress.StatusFailed || done.Error == "" {
		t.Fatalf("done = %+v, want Failed with error", done)
	}
}

// TestStartManualDedupConflict verifies a duplicate manual trigger for the same
// (parser, document) while one is in flight is rejected as AlreadyExists (the
// gRPC shape the API layer maps to HTTP 409).
func TestStartManualDedupConflict(t *testing.T) {
	p, _ := newTestParser(t)
	fake := newFakeHandler()
	defer close(fake.block)
	rule := rpaDirectoryRule()
	ctx := context.Background()

	if _, err := p.start(ctx, 1, 1, rule, fake, nil, false); err != nil {
		t.Fatalf("first start: %v", err)
	}
	_, err := p.start(ctx, 1, 1, rule, fake, nil, false)
	if status.Code(err) != codes.AlreadyExists {
		t.Fatalf("second manual start err = %v, want AlreadyExists", err)
	}
}

// TestStartOnUploadJoinsInFlight verifies on_upload does NOT 409 on a duplicate:
// a concurrent dispatch joins the in-flight run and returns its requestId.
func TestStartOnUploadJoinsInFlight(t *testing.T) {
	p, _ := newTestParser(t)
	fake := newFakeHandler()
	defer close(fake.block)
	rule := rpaDirectoryRule()
	ctx := context.Background()

	req1, err := p.start(ctx, 1, 1, rule, fake, nil, true)
	if err != nil {
		t.Fatalf("first on_upload start: %v", err)
	}
	req2, err := p.start(ctx, 1, 1, rule, fake, nil, true)
	if err != nil {
		t.Fatalf("second on_upload start: %v, want join (nil err)", err)
	}
	if req2 == "" || req2 != req1 {
		t.Fatalf("on_upload join: got %q, want the in-flight requestId %q", req2, req1)
	}
}

// seedOnUploadDoc inserts a leaf document carrying an on_upload escrow rule, so
// Dispatch routes an upload batch to the escrow handler.
func seedOnUploadDoc(t *testing.T, db *gorm.DB, id int64) {
	t.Helper()
	meta, _ := json.Marshal(map[string]any{
		"file_parse_config": map[string]any{
			"rules": []map[string]any{{
				"parser":       string(constants.DocumentTemplateFileParseParserEscrowInfo),
				"source":       string(constants.DocumentTemplateFileParseSourceSingleFile),
				"triggered_by": string(constants.FileParseTriggeredTypeUpload),
			}},
		},
	})
	if err := db.Exec(
		`INSERT INTO documents (id, type, name, metadata) VALUES (?, 'Leaf', 'Wire Instructions', ?)`,
		id, string(meta),
	).Error; err != nil {
		t.Fatalf("seed doc: %v", err)
	}
}

// TestDispatchOnUploadRunsAndFinishes drives the on_upload dispatch path over a
// real AllRepos (in-memory SQLite): the uploaded doc's on_upload rule routes to
// the escrow handler, the run fires async, and reaches terminal Finished. The
// escrow extract itself (S3 download + pdfium + vision LLM) is handler-tested
// separately; this covers the dispatch routing + lifecycle wiring over the DB.
func TestDispatchOnUploadRunsAndFinishes(t *testing.T) {
	db, p, pm := setupRecoverParser(t)
	fake := &countingHandler{}
	p.Register(constants.DocumentTemplateFileParseParserEscrowInfo, fake)
	seedOnUploadDoc(t, db, 300)

	reqId, err := p.Dispatch(context.Background(), DispatchParam{
		TransactionId: 7,
		DocumentId:    300,
		Trigger:       TriggerOnUpload,
		Files:         []model.File{{Base: model.Base{Id: 1}, FileUrl: "s3://wire.pdf", FileName: "wire.pdf"}},
	})
	if err != nil || reqId == "" {
		t.Fatalf("dispatch: err=%v reqId=%q", err, reqId)
	}
	done := waitDone(t, pm, reqId)
	if done.Status != progress.StatusFinished {
		t.Fatalf("done = %+v, want Finished", done)
	}
	if atomic.LoadInt32(&fake.calls) != 1 {
		t.Fatalf("handler calls = %d, want 1", fake.calls)
	}
}

// TestReporterFilesCheckpointRoundTrip verifies the Redis-backed reporter
// persists the Extract files checkpoint (fileId → provider file id) so a
// recovered run reuses uploaded ids.
func TestReporterFilesCheckpointRoundTrip(t *testing.T) {
	_, pm := newTestParser(t)
	ctx := context.Background()
	rep := redisReporter{progMgr: pm, requestId: "req"}

	if files := rep.Files(ctx); files != nil {
		t.Fatalf("fresh files = %+v, want nil", files)
	}
	rep.SaveFile(ctx, 10, progress.FileCheckpoint{Provider: "claude", ProviderFileID: "file-10"})
	files := rep.Files(ctx)
	if len(files) != 1 || files[10].ProviderFileID != "file-10" {
		t.Fatalf("round-trip = %+v, want file-10", files)
	}
}
