package parse_file

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/common/redisx"
)

func newTestProgressManager(t *testing.T) (*ProgressManager, *miniredis.Miniredis) {
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
	return NewProgressManager(rdb, logx.WithContext(context.Background())), mr
}

func TestDedupLockSetNX(t *testing.T) {
	pm, _ := newTestProgressManager(t)
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

func TestQueueSerial(t *testing.T) {
	pm, _ := newTestProgressManager(t)
	ctx := context.Background()
	for _, fid := range []string{"10", "11", "12"} {
		if err := pm.EnqueueFile(ctx, 1, 5, fid); err != nil {
			t.Fatalf("enqueue %s: %v", fid, err)
		}
	}
	for _, want := range []string{"10", "11", "12"} {
		got, ok, err := pm.DequeueFile(ctx, 1, 5)
		if err != nil || !ok || got != want {
			t.Fatalf("dequeue: got=%q ok=%v err=%v, want %q", got, ok, err, want)
		}
	}
	got, ok, err := pm.DequeueFile(ctx, 1, 5)
	if err != nil || ok || got != "" {
		t.Fatalf("empty dequeue: got=%q ok=%v err=%v, want empty", got, ok, err)
	}
}

func TestProgressWriteReadDone(t *testing.T) {
	pm, _ := newTestProgressManager(t)
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
	pm, _ := newTestProgressManager(t)
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
	pm, _ := newTestProgressManager(t)
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
	pm, _ := newTestProgressManager(t)
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
