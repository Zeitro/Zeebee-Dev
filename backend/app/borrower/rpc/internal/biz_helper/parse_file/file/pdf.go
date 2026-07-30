// Package file downloads document files, gates parsing by type, and slices
// leading PDF pages (wasm pdfium, no CGO) to cap provider upload cost.
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
	// maxPDFBytes caps a downloaded file at 10 MB.
	maxPDFBytes = 10 * 1024 * 1024
)

// wazero WebAssembly: no CGO, no libpdfium. CGO_ENABLED=0 alpine build unchanged.
var (
	pdfPoolOnce sync.Once
	pdfPool     pdfium.Pool
	pdfPoolErr  error
)

func getPdfPool() (pdfium.Pool, error) {
	pdfPoolOnce.Do(func() {
		pdfPool, pdfPoolErr = webassembly.Init(webassembly.Config{
			MinIdle:  0,
			MaxIdle:  2,
			MaxTotal: 5, // caps concurrent pdfium parses; ponytail: raise if throughput demands
		})
	})
	return pdfPool, pdfPoolErr
}

func withPdfInstance(fn func(instance pdfium.Pdfium) error) error {
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

// openPDF opens r as a doc; defer the returned closer.
func openPDF(instance pdfium.Pdfium, r io.ReadSeeker) (references.FPDF_DOCUMENT, func(), error) {
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

// FirstPagesPDF returns a new PDF containing the first n pages of src.
// nil bytes + nil error = src has ≤n pages (caller uploads whole file).
// error = pdfium failed (caller logs + skips the file).
func FirstPagesPDF(src io.ReadSeeker, n int) ([]byte, error) {
	var out []byte
	err := withPdfInstance(func(inst pdfium.Pdfium) error {
		doc, closeSrc, err := openPDF(inst, src)
		if err != nil {
			return err
		}
		defer closeSrc()
		count, err := inst.FPDF_GetPageCount(&requests.FPDF_GetPageCount{Document: doc})
		if err != nil {
			return err
		}
		if count.PageCount <= n {
			return nil // ≤n pages: no slice needed
		}
		newDoc, err := inst.FPDF_CreateNewDocument(&requests.FPDF_CreateNewDocument{})
		if err != nil {
			return err
		}
		defer func() {
			_, _ = inst.FPDF_CloseDocument(&requests.FPDF_CloseDocument{Document: newDoc.Document})
		}()
		pageRange := fmt.Sprintf("1-%d", n)
		if _, err := inst.FPDF_ImportPages(&requests.FPDF_ImportPages{
			Source:      doc,
			Destination: newDoc.Document,
			PageRange:   &pageRange,
			Index:       0,
		}); err != nil {
			return err
		}
		var buf bytes.Buffer
		if _, err := inst.FPDF_SaveAsCopy(&requests.FPDF_SaveAsCopy{
			Document:   newDoc.Document,
			FileWriter: &buf,
		}); err != nil {
			return err
		}
		out = buf.Bytes()
		return nil
	})
	return out, err
}
