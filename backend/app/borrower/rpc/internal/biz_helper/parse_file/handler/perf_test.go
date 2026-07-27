package handler

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/common/llm"
)

// TestSC003RPAConsolidatedSingleCall: 5–10 PDFs → ONE consolidated Chat over all
// uploaded ids (in-call arbitration), not one call per file (SC-003). Reports the
// single-call latency at each scale.
func TestSC003RPAConsolidatedSingleCall(t *testing.T) {
	srv := byteServer(t, buildPDF(t, 2))
	for _, n := range []int{5, 10} {
		t.Run(fmt.Sprintf("files=%d", n), func(t *testing.T) {
			mock := &llm.MockClient{ChatContent: `{"purchase_amount":1750000,"coe_date":"09/15/2026"}`}
			h := &RPAHandler{base{logger: logx.WithContext(context.Background()), llm: mock}}
			rep := newMemReporter()

			files := make([]collectedFile, n)
			for i := range files {
				files[i] = collectedFile{
					FileId: int64(i + 1), Name: fmt.Sprintf("doc-%d.pdf", i), Url: srv.URL,
				}
			}

			start := time.Now()
			fields, err := h.extractStage(context.Background(), rep, files)
			lat := time.Since(start)
			if err != nil {
				t.Fatalf("extractStage: %v", err)
			}
			if !fields.hasAny() {
				t.Fatal("expected extracted fields")
			}
			if calls := mock.ChatCalls(); calls != 1 {
				t.Fatalf("n=%d: Chat calls=%d, want 1 (single consolidated)", n, calls)
			}
			if len(mock.UploadedFiles) != n {
				t.Fatalf("n=%d: uploads=%d, want %d", n, len(mock.UploadedFiles), n)
			}
			if atts := mock.Attachments(); len(atts) != n {
				t.Fatalf("n=%d: attachments=%d, want %d in the one call", n, len(atts), n)
			}
			t.Logf("SC-003 n=%d: 1 consolidated call over %d uploads, latency=%s", n, n, lat)
		})
	}
}
