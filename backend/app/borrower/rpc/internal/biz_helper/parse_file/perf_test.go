package parse_file

import (
	"context"
	"math"
	"slices"
	"testing"
	"time"

	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/model"
)

// TestSC002OnUploadNonBlocking: on_upload dispatch returns immediately — parse
// runs in a goroutine, never on the upload path (SC-002). A handler that parks
// indefinitely isolates pure dispatch cost; if dispatch waited on parse, the
// loop would deadlock.
func TestSC002OnUploadNonBlocking(t *testing.T) {
	db, p, pm := setupRecoverParser(t)
	fake := newFakeHandler() // Parse blocks on block chan until released
	p.Register(constants.DocumentTemplateFileParseParserEscrowInfo, fake)

	const n = 40
	lats := make([]time.Duration, 0, n)
	reqIds := make([]string, 0, n)
	ctx := context.Background()
	for i := range n {
		docId := int64(300 + i)
		seedOnUploadDoc(t, db, docId)
		start := time.Now()
		reqId, err := p.Dispatch(ctx, DispatchParam{
			TransactionId: 7,
			DocumentId:    docId,
			Trigger:       TriggerOnUpload,
			Files:         []model.File{{Base: model.Base{Id: int64(i + 1)}, FileName: "wire.pdf"}},
		})
		lats = append(lats, time.Since(start))
		if err != nil || reqId == "" {
			t.Fatalf("dispatch %d: err=%v reqId=%q", i, err, reqId)
		}
		reqIds = append(reqIds, reqId)
	}

	p95 := percentile(lats, 95)
	if p95 >= 100*time.Millisecond {
		t.Fatalf("dispatch p95=%s, want <100ms (parse must not block upload)", p95)
	}
	close(fake.block) // release parked parses, then await their terminal Done
	for _, rid := range reqIds {
		waitDone(t, pm, rid)
	}
	t.Logf("SC-002 dispatch over %d on_upload: p50=%s p95=%s max=%s (parse async)",
		n, percentile(lats, 50), p95, percentile(lats, 100))
}

// percentile returns the pct-th percentile of xs (pct in 1..100).
func percentile(xs []time.Duration, pct int) time.Duration {
	if len(xs) == 0 {
		return 0
	}
	s := make([]time.Duration, len(xs))
	copy(s, xs)
	slices.Sort(s)
	idx := max(int(math.Ceil(float64(pct)/100*float64(len(s))))-1, 0)
	return s[idx]
}
