// Package file wraps pdfium (via wazero) for PDF download + text render and
// gates parsing by file type. No CGO.
package file

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/klippa-app/go-pdfium"
	"github.com/klippa-app/go-pdfium/references"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/webassembly"
)

const (
	// maxPDFBytes caps a downloaded PDF at 10 MB.
	maxPDFBytes = 10 * 1024 * 1024
	// pdfPoolSize caps concurrent pdfium parses; extra callers block on GetInstance.
	// ponytail: raise if upload throughput demands it.
	pdfPoolSize = 5
	// MaxRPAPages caps how many leading pages of each contract PDF are fed to
	// the rpa text extractor (key info sits up front; whole-file is wasteful).
	MaxRPAPages = 10
)

var (
	pdfPoolOnce sync.Once
	pdfPool     pdfium.Pool
	pdfPoolErr  error
)

func getPdfPool() (pdfium.Pool, error) {
	pdfPoolOnce.Do(func() {
		// ponytail: wazero pure-Go runtime, no CGO, no pdfium binary install.
		pdfPool, pdfPoolErr = webassembly.Init(webassembly.Config{
			MinIdle:  0,
			MaxIdle:  2,
			MaxTotal: pdfPoolSize,
		})
	})
	return pdfPool, pdfPoolErr
}

// WithPdfInstance runs fn on a pooled pdfium instance. Exported so the handler
// package (which owns the LLM step) can drive pdfium without re-importing the
// pool plumbing.
func WithPdfInstance(fn func(instance pdfium.Pdfium) error) error {
	pool, err := getPdfPool()
	if err != nil {
		return err
	}
	inst, err := pool.GetInstance(60 * time.Second)
	if err != nil {
		return err
	}
	defer func() { _ = inst.Close() }()
	return fn(inst)
}

// DownloadPDF fetches a presigned file URL (≤10 MB) and returns a seekable
// reader for pdfium. The body is buffered into memory because a streaming http
// body is not seekable and pdfium seeks the xref ranges itself.
func DownloadPDF(ctx context.Context, url string) (io.ReadSeeker, error) {
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
	if res.ContentLength > maxPDFBytes {
		return nil, fmt.Errorf("pdf too large: %d bytes", res.ContentLength)
	}
	// LimitReader +1 so an undersized Content-Length still trips the cap.
	b, err := io.ReadAll(io.LimitReader(res.Body, maxPDFBytes+1))
	if err != nil {
		return nil, err
	}
	if len(b) > maxPDFBytes {
		return nil, fmt.Errorf("pdf exceeds 10MB")
	}
	return bytes.NewReader(b), nil
}

// OpenPDF opens a seekable reader as a pdfium document.
func OpenPDF(instance pdfium.Pdfium, r io.ReadSeeker) (references.FPDF_DOCUMENT, func(), error) {
	size, err := r.Seek(0, io.SeekEnd)
	if err != nil {
		return "", nil, err
	}
	if _, err = r.Seek(0, io.SeekStart); err != nil {
		return "", nil, err
	}
	doc, err := instance.OpenDocument(&requests.OpenDocument{
		FileReader:     r,
		FileReaderSize: size,
	})
	if err != nil {
		return "", nil, err
	}
	return doc.Document, func() {
		_, _ = instance.FPDF_CloseDocument(&requests.FPDF_CloseDocument{Document: doc.Document})
	}, nil
}

func readPageText(instance pdfium.Pdfium, doc references.FPDF_DOCUMENT, page int) (string, error) {
	tp, err := instance.FPDFText_LoadPage(&requests.FPDFText_LoadPage{
		Page: requests.Page{ByIndex: &requests.PageByIndex{Document: doc, Index: page}},
	})
	if err != nil {
		return "", err
	}
	defer func() { _, _ = instance.FPDFText_ClosePage(&requests.FPDFText_ClosePage{TextPage: tp.TextPage}) }()

	count, err := instance.FPDFText_CountChars(&requests.FPDFText_CountChars{TextPage: tp.TextPage})
	if err != nil {
		return "", err
	}
	if count.Count <= 0 {
		return "", nil
	}
	res, err := instance.FPDFText_GetText(&requests.FPDFText_GetText{
		TextPage:   tp.TextPage,
		StartIndex: 0,
		Count:      count.Count,
	})
	if err != nil {
		return "", err
	}
	return res.Text, nil
}

// ReadFirstNPages concatenates text from the first n pages (or fewer if the
// document is shorter) for the rpa text-extraction path.
func ReadFirstNPages(instance pdfium.Pdfium, doc references.FPDF_DOCUMENT, n int) (string, error) {
	var b strings.Builder
	for i := 0; i < n; i++ {
		text, err := readPageText(instance, doc, i)
		if err != nil {
			// A missing page beyond the document's length is not fatal.
			if isPageOutOfRange(err) {
				break
			}
			return "", err
		}
		if text == "" {
			break
		}
		b.WriteString(text)
		b.WriteString("\n")
	}
	return b.String(), nil
}

// isPageOutOfRange returns true for pdfium's "page index too large" error so
// readFirstNPages can stop cleanly at the document end.
func isPageOutOfRange(err error) bool {
	return err != nil && strings.Contains(err.Error(), "page")
}

// RenderPagePNG renders one page to PNG bytes for the vision-extraction path.
func RenderPagePNG(instance pdfium.Pdfium, doc references.FPDF_DOCUMENT, page int) ([]byte, error) {
	res, err := instance.RenderToFile(&requests.RenderToFile{
		RenderPageInDPI: &requests.RenderPageInDPI{
			Page:       requests.Page{ByIndex: &requests.PageByIndex{Document: doc, Index: page}},
			DPI:        150,
			RenderForm: true,
		},
		OutputFormat: requests.RenderToFileOutputFormatPNG,
		OutputTarget: requests.RenderToFileOutputTargetBytes,
	})
	if err != nil {
		return nil, err
	}
	if res.ImageBytes == nil {
		return nil, fmt.Errorf("render returned no image bytes")
	}
	return *res.ImageBytes, nil
}

// IsPDF gates parsing to PDF files only; non-PDF uploads are skipped.
func IsPDF(name string) bool {
	return strings.HasSuffix(strings.ToLower(name), ".pdf")
}
