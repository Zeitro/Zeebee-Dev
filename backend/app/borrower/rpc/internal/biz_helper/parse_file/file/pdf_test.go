package file

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// These tests cover the non-pdfium parts of pdf.go: the file-type gate and the
// download size cap / status guard. readFirstNPages / RenderPagePNG need a real
// pdfium instance + a PDF fixture (the go-pdfium interface is too large to mock
// cheaply), so they are exercised in integration rather than here.

func TestIsPDF(t *testing.T) {
	cases := map[string]bool{
		"wire.pdf":     true,
		"WIRE.PDF":     true,
		"contract.Pdf": true,
		"notes.docx":   false,
		"noext":        false,
		"pdf":          false,
		"scan.pdf.bak": false,
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
	// pdfium seeks the xref itself — the returned reader must be seekable.
	if _, err := rs.Seek(0, io.SeekStart); err != nil {
		t.Fatalf("seek: %v", err)
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
