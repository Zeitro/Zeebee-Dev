// Package file downloads document files and gates parsing by type. The provider
// (claude) reads PDF/image bytes directly via file upload — no pdfium, no text
// render, no CGO.
package file

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	// maxPDFBytes caps a downloaded file at 10 MB.
	maxPDFBytes = 10 * 1024 * 1024
)

// DownloadPDF fetches a presigned file URL (≤10 MB) and returns its bytes for
// provider upload. Buffered in memory because the body is consumed once, inline.
func DownloadPDF(ctx context.Context, url string) (io.Reader, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, fmt.Errorf("download http %d", res.StatusCode)
	}
	b, err := readCapped(res.Body, maxPDFBytes)
	if err != nil {
		return nil, err
	}
	return strings.NewReader(string(b)), nil
}

// IsPDF gates parsing to PDF files; non-PDF uploads are skipped.
func IsPDF(name string) bool {
	return strings.HasSuffix(strings.ToLower(name), ".pdf")
}

// readCapped reads up to max+1 bytes; the +1 trips the cap on undersized Content-Length.
func readCapped(r io.Reader, max int) ([]byte, error) {
	b, err := io.ReadAll(io.LimitReader(r, int64(max)+1))
	if err != nil {
		return nil, err
	}
	if len(b) > max {
		return nil, fmt.Errorf("file exceeds %d bytes", max)
	}
	return b, nil
}
