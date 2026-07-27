package file

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/signintech/gopdf"
)

// These tests cover pdf.go + image.go: the file-type gates and the download
// size cap / status guard. The provider reads bytes directly — no pdfium render.

func TestIsPDF(t *testing.T) {
	cases := map[string]bool{
		"wire.pdf": true, "WIRE.PDF": true, "contract.Pdf": true,
		"notes.docx": false, "noext": false, "pdf": false, "scan.pdf.bak": false,
	}
	for name, want := range cases {
		if got := IsPDF(name); got != want {
			t.Errorf("IsPDF(%q) = %v, want %v", name, got, want)
		}
	}
}

func TestDownloadPDFSuccess(t *testing.T) {
	body := []byte("%PDF-1.4 fake but seekable")
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write(body)
	}))
	defer srv.Close()

	rs, err := DownloadPDF(context.Background(), srv.URL)
	if err != nil {
		t.Fatalf("download: %v", err)
	}
	got, err := io.ReadAll(rs)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if !bytes.Equal(got, body) {
		t.Fatalf("body = %q, want %q", got, body)
	}
}

// TestDownloadPDFRejectsOversize verifies the 10 MB cap rejects a body that
// exceeds it — the robust path (LimitReader) that catches a server lying about
// Content-Length.
func TestDownloadPDFRejectsOversize(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// maxPDFBytes+1 bytes; one over the cap.
		chunk := make([]byte, 1024)
		_, _ = w.Write(bytes.Repeat(chunk, maxPDFBytes/1024+1))
	}))
	defer srv.Close()

	if _, err := DownloadPDF(context.Background(), srv.URL); err == nil {
		t.Fatal("oversized download must be rejected")
	}
}

func TestDownloadPDFBadStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	if _, err := DownloadPDF(context.Background(), srv.URL); err == nil {
		t.Fatal("non-2xx download must error")
	}
}

// buildPDF renders an N empty-page PDF (no font needed) for slicing tests.
func buildPDF(t *testing.T, pages int) []byte {
	t.Helper()
	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	for i := 0; i < pages; i++ {
		pdf.AddPage()
	}
	var buf bytes.Buffer
	if _, err := pdf.WriteTo(&buf); err != nil {
		t.Fatalf("gopdf write: %v", err)
	}
	return buf.Bytes()
}

func TestFirstPagesPDFShortReturnsNil(t *testing.T) {
	got, err := FirstPagesPDF(bytes.NewReader(buildPDF(t, MaxRPAPages-1)), MaxRPAPages)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if got != nil {
		t.Errorf("≤cap must return nil (whole file), got %d bytes", len(got))
	}
}

func TestFirstPagesPDFSlicesToCap(t *testing.T) {
	full := buildPDF(t, 11)
	sliced, err := FirstPagesPDF(bytes.NewReader(full), MaxRPAPages)
	if err != nil {
		t.Fatalf("slice: %v", err)
	}
	if sliced == nil {
		t.Fatal("11-page PDF must be sliced")
	}
	if !bytes.HasPrefix(sliced, []byte("%PDF")) {
		t.Fatalf("sliced output not a PDF: %q", sliced[:8])
	}
	if len(sliced) >= len(full) {
		t.Fatalf("slice must shrink: %d >= %d", len(sliced), len(full))
	}
}

func TestFirstPagesPDFBadInput(t *testing.T) {
	if _, err := FirstPagesPDF(strings.NewReader("garbage"), MaxRPAPages); err == nil {
		t.Fatal("garbage must error")
	}
}
