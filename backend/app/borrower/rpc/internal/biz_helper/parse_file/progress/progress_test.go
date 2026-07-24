package progress

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/common/redisx"
)

func newTestProgressManager(t *testing.T) *ProgressManager {
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
	return NewProgressManager(rdb, logx.WithContext(context.Background()))
}

func TestDedupLockSetNX(t *testing.T) {
	pm := newTestProgressManager(t)
	ctx := context.Background()
	ok, err := pm.AcquireDedupLock(ctx, 1, "sig", "req1")
	if err != nil || !ok {
		t.Fatalf("first acquire: ok=%v err=%v", ok, err)
	}
	ok2, err := pm.AcquireDedupLock(ctx, 1, "sig", "req2")
	if err != nil {
		t.Fatalf("second acquire err: %v", err)
	}
	if ok2 {
		t.Fatal("second acquire must fail (already held)")
	}
	if err := pm.ReleaseDedupLock(ctx, 1, "sig"); err != nil {
		t.Fatalf("release: %v", err)
	}
	ok3, _ := pm.AcquireDedupLock(ctx, 1, "sig", "req3")
	if !ok3 {
		t.Fatal("re-acquire after release must succeed")
	}
}

func TestProgressWriteReadDone(t *testing.T) {
	pm := newTestProgressManager(t)
	ctx := context.Background()
	rid := "req-1"
	stages := []StageState{
		{Name: StageCollect, Status: StageStatusDone, Total: 1, Processed: 1},
		{Name: StageExtract, Status: StageStatusRunning, Total: 8, Processed: 3},
	}
	if err := pm.InitProgress(ctx, rid, stages); err != nil {
		t.Fatalf("init: %v", err)
	}
	got, err := pm.GetProgress(ctx, rid)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got == nil || got.Stage != StageExtract || len(got.Stages) != 2 {
		t.Fatalf("progress = %+v", got)
	}
	if err := pm.WriteDone(ctx, rid, StatusFinished, ""); err != nil {
		t.Fatalf("write done: %v", err)
	}
	if p, _ := pm.GetProgress(ctx, rid); p != nil {
		t.Fatal("progress key must be cleared after WriteDone")
	}
	d, err := pm.GetDone(ctx, rid)
	if err != nil || d == nil || d.Status != StatusFinished {
		t.Fatalf("done = %+v err=%v", d, err)
	}
}

func TestReqMap(t *testing.T) {
	pm := newTestProgressManager(t)
	ctx := context.Background()
	if err := pm.SetReqMap(ctx, 1, "sig", "req-9"); err != nil {
		t.Fatalf("set: %v", err)
	}
	got, err := pm.GetReqMap(ctx, 1, "sig")
	if err != nil || got != "req-9" {
		t.Fatalf("get: got=%q err=%v", got, err)
	}
}

func TestHoldSetNX(t *testing.T) {
	pm := newTestProgressManager(t)
	ctx := context.Background()
	ok, _ := pm.AcquireHold(ctx, "req", "nodeA")
	if !ok {
		t.Fatal("first hold must succeed")
	}
	ok2, _ := pm.AcquireHold(ctx, "req", "nodeB")
	if ok2 {
		t.Fatal("second hold must fail")
	}
}

func TestProgressJSONRoundTrip(t *testing.T) {
	pm := newTestProgressManager(t)
	ctx := context.Background()
	p := &Progress{Status: StatusRunning, Stage: StageExtract,
		Stages: []StageState{{Name: StageExtract, Status: StageStatusRunning, Processed: 1, Total: 2}}}
	if err := pm.SaveProgress(ctx, "req", p); err != nil {
		t.Fatalf("save: %v", err)
	}
	got, _ := pm.GetProgress(ctx, "req")
	if got == nil || got.Stages[0].Total != 2 {
		t.Fatalf("roundtrip: %+v", got)
	}
	// progress JSON must NOT carry an overall percent field.
	raw, _ := json.Marshal(p)
	if strings.Contains(string(raw), `"percent"`) {
		t.Fatalf("progress JSON leaks percent: %s", raw)
	}
}

// TestRecoveryScan verifies the Redis contract the recovery coordinator leans
// on: a job persists and is scannable; an in-flight job has no done record
// (reclaimable); the hold guarantees a single owner across nodes; once done is
// written the job is terminal. (The tryReclaim dispatch above this is glue —
// faking AllRepos would need an interface refactor across the repo layer.)
func TestRecoveryScan(t *testing.T) {
	pm := newTestProgressManager(t)
	ctx := context.Background()
	job := &Job{RequestId: "req-recover", TransactionId: 7, RuleSig: "rpa|9", Parser: "rpa", Mode: "Directory", DocumentId: 9}
	if err := pm.WriteJob(ctx, job); err != nil {
		t.Fatalf("write job: %v", err)
	}
	jobs, err := pm.ListJobs(ctx)
	if err != nil {
		t.Fatalf("list jobs: %v", err)
	}
	if len(jobs) != 1 || jobs[0].RequestId != "req-recover" || jobs[0].DocumentId != 9 {
		t.Fatalf("ListJobs = %+v", jobs)
	}
	if done, _ := pm.GetDone(ctx, "req-recover"); done != nil {
		t.Fatal("in-flight job must have no done record (reclaimable)")
	}
	if ok, _ := pm.AcquireHold(ctx, "req-recover", "nodeA"); !ok {
		t.Fatal("first hold acquire must succeed")
	}
	if ok, _ := pm.AcquireHold(ctx, "req-recover", "nodeB"); ok {
		t.Fatal("second node must not steal a live hold")
	}
	if err := pm.WriteDone(ctx, "req-recover", StatusFinished, ""); err != nil {
		t.Fatalf("write done: %v", err)
	}
	if done, _ := pm.GetDone(ctx, "req-recover"); done == nil || done.Status != StatusFinished {
		t.Fatalf("terminal done = %+v", done)
	}
}

// TestProgressFilesCheckpointRoundTrip verifies the Extract resume checkpoint
// (fileId → provider file id) survives in its own job-TTL key — rpa resume leans
// on it to reuse uploaded ids instead of re-uploading.
func TestProgressFilesCheckpointRoundTrip(t *testing.T) {
	pm := newTestProgressManager(t)
	ctx := context.Background()
	if got, _ := pm.GetFilesCheckpoint(ctx, "req"); got != nil {
		t.Fatalf("fresh checkpoint = %+v, want nil", got)
	}
	if err := pm.SaveFileCheckpoint(ctx, "req", 10, FileCheckpoint{Provider: "claude", ProviderFileID: "file-10"}); err != nil {
		t.Fatalf("save 10: %v", err)
	}
	if err := pm.SaveFileCheckpoint(ctx, "req", 11, FileCheckpoint{Provider: "claude", ProviderFileID: "file-11"}); err != nil {
		t.Fatalf("save 11: %v", err)
	}
	got, _ := pm.GetFilesCheckpoint(ctx, "req")
	if len(got) != 2 || got[10].ProviderFileID != "file-10" || got[11].Provider != "claude" {
		t.Fatalf("checkpoint = %+v, want 2 claude entries", got)
	}
	if err := pm.DeleteFilesCheckpoint(ctx, "req"); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if got, _ := pm.GetFilesCheckpoint(ctx, "req"); got != nil {
		t.Fatalf("after delete = %+v, want nil", got)
	}
}
