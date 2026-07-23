package parse_file

import (
	"bytes"
	"context"
	"fmt"
	"sync/atomic"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/rpc/internal/biz_helper/parse_file/handler"
	"zeebee-backend/app/borrower/rpc/internal/biz_helper/parse_file/progress"
	"zeebee-backend/common/redisx"
)

// These benchmarks back the two performance success criteria:
//   - SC-002: an on_upload parse must not slow the upload response. The extract
//     runs off the upload's critical path (GoSafe), so the parse-induced upload
//     overhead is just start(): seed Running progress + acquire the dedup lock +
//     spawn the async run. BenchmarkStartOnUpload measures that.
//   - SC-003: a manual rpa parse over 5–10 PDFs finishes in minutes. The wall
//     clock is dominated by per-file LLM latency (not measurable headless), so
//     BenchmarkRPAExtractFile measures the deterministic per-file floor the
//     system controls (pdfium open + leading-pages text + one ChatJSON round
//     trip with a no-op LLM); multiply by N and add real LLM latency per file.

// benchmarkParser builds a coordinator over miniredis for the start() benchmark.
func benchmarkParser(b *testing.B) (*Parser, *progress.ProgressManager) {
	b.Helper()
	mr, err := miniredis.Run()
	if err != nil {
		b.Fatalf("miniredis: %v", err)
	}
	b.Cleanup(mr.Close)
	rdb, err := redisx.GetRedis(redisx.Opt{Host: mr.Addr(), Type: "node"})
	if err != nil {
		b.Fatalf("redisx: %v", err)
	}
	pm := progress.NewProgressManager(rdb, logx.WithContext(context.Background()))
	return NewParser(nil, stubLLM{}, pm, logx.WithContext(context.Background())), pm
}

// BenchmarkStartOnUpload is the SC-002 upload-overhead benchmark: the
// synchronous work an on_upload dispatch adds to the upload response. Each
// iteration uses a fresh documentId so runs do not join. The spawned runs use an
// instant handler and self-clean (WriteDone + key reclaim), so the timer captures
// only the upload-critical-path (seed + lock + goroutine spawn), not the extract.
func BenchmarkStartOnUpload(b *testing.B) {
	p, _ := benchmarkParser(b)
	fake := &countingHandler{}
	rule := rpaDirectoryRule()
	ctx := context.Background()
	var docID atomic.Int64
	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		if _, err := p.start(ctx, 1, docID.Add(1), rule, fake, nil, true); err != nil {
			b.Fatalf("start: %v", err)
		}
	}
}

// minimalPerfPDF is a valid blank 1-page PDF (correct xref) so pdfium opens +
// reads it. The mock LLM ignores the (empty) text, so this isolates the pdfium +
// glue floor from LLM cost. Mirrors handler/escrow_test.minimalPDF; kept here to
// keep the perf benchmark self-contained without exporting a test helper.
func minimalPerfPDF() []byte {
	var buf bytes.Buffer
	buf.WriteString("%PDF-1.4\n")
	o1 := buf.Len()
	buf.WriteString("1 0 obj\n<< /Type /Catalog /Pages 2 0 R >>\nendobj\n")
	o2 := buf.Len()
	buf.WriteString("2 0 obj\n<< /Type /Pages /Kids [3 0 R] /Count 1 >>\nendobj\n")
	o3 := buf.Len()
	buf.WriteString("3 0 obj\n<< /Type /Page /Parent 2 0 R /MediaBox [0 0 612 792] >>\nendobj\n")
	xref := buf.Len()
	buf.WriteString("xref\n0 4\n0000000000 65535 f \n")
	fmt.Fprintf(&buf, "%010d 00000 n \n", o1)
	fmt.Fprintf(&buf, "%010d 00000 n \n", o2)
	fmt.Fprintf(&buf, "%010d 00000 n \n", o3)
	fmt.Fprintf(&buf, "trailer\n<< /Size 4 /Root 1 0 R >>\nstartxref\n%d\n%%%%EOF\n", xref)
	return buf.Bytes()
}

// BenchmarkRPAExtractFile is the SC-003 per-file floor: download-free pdfium open
// + leading-pages text read + one ChatJSON (no-op stub LLM) over a real PDF.
// ns/op × N (5–10) is the non-LLM portion of a manual rpa parse; real provider
// latency per file is additive and dominates wall clock.
func BenchmarkRPAExtractFile(b *testing.B) {
	pdf := minimalPerfPDF()
	ctx := context.Background()
	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		if _, err := handler.ExtractRPAFields(ctx, stubLLM{}, 1, bytes.NewReader(pdf)); err != nil {
			b.Fatalf("extract: %v", err)
		}
	}
}
