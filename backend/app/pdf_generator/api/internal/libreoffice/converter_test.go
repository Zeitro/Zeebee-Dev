package libreoffice

import (
	"archive/zip"
	"bytes"
	"context"
	"errors"
	"os/exec"
	"testing"
	"time"
)

// sofficeAvailable reports whether `soffice` is on PATH (CI/local without
// LibreOffice skip the conversion tests).
func sofficeAvailable(tb testing.TB) bool {
	tb.Helper()
	_, err := exec.LookPath("soffice")
	return err == nil
}

// buildDocx assembles a minimal valid OOXML .docx (the 3 required zip parts).
// Programmatic so the test is self-contained — no binary fixture in the repo.
func buildDocx(tb testing.TB) []byte {
	tb.Helper()
	parts := map[string]string{
		"[Content_Types].xml": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types"><Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/><Default Extension="xml" ContentType="application/xml"/><Override PartName="/word/document.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.document.main+xml"/></Types>`,
		"_rels/.rels": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"><Relationship Id="r1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="word/document.xml"/></Relationships>`,
		"word/document.xml": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main"><w:body><w:p><w:r><w:t>Escrow No. 202656-RL Mission Escrow Officer Renee Lo</w:t></w:r></w:p></w:body></w:document>`,
	}
	var buf bytes.Buffer
	w := zip.NewWriter(&buf)
	for name, body := range parts {
		f, err := w.Create(name)
		if err != nil {
			tb.Fatal(err)
		}
		if _, err := f.Write([]byte(body)); err != nil {
			tb.Fatal(err)
		}
	}
	if err := w.Close(); err != nil {
		tb.Fatal(err)
	}
	return buf.Bytes()
}

func TestConvert_DocxToPdf(t *testing.T) {
	if !sofficeAvailable(t) {
		t.Skip("soffice not installed; skipping conversion smoke test")
	}
	conv, err := NewConverter(2, 60*time.Second)
	if err != nil {
		t.Fatalf("NewConverter: %v", err)
	}
	defer conv.Close()

	pdf, err := conv.Convert(context.Background(), buildDocx(t), "docx")
	if err != nil {
		t.Fatalf("Convert: %v", err)
	}
	if !bytes.HasPrefix(pdf, []byte("%PDF-")) {
		t.Fatalf("output is not a pdf, got %q", pdf[:min(16, len(pdf))])
	}
}

func TestConvert_CorruptInput(t *testing.T) {
	if !sofficeAvailable(t) {
		t.Skip("soffice not installed")
	}
	conv, err := NewConverter(1, 30*time.Second)
	if err != nil {
		t.Fatalf("NewConverter: %v", err)
	}
	defer conv.Close()

	// Broken zip (PK magic + junk, no central directory): soffice fails to load
	// the package, exits 0 but writes no pdf — Convert surfaces it as an error.
	// Plain text named .docx is NOT corrupt — soffice imports it as text, hence
	// the structurally-broken zip is the reliable corrupt fixture.
	corrupt := append([]byte("PK\x03\x04\x00\x00"), []byte("garbage-not-a-valid-zip-central-dir")...)
	if _, err := conv.Convert(context.Background(), corrupt, "docx"); err == nil {
		t.Fatal("expected error for corrupt input, got nil")
	}
}

func TestNewConverter_MissingSoffice(t *testing.T) {
	// Empty PATH makes LookPath fail everywhere, simulating an image without LibreOffice.
	t.Setenv("PATH", "")
	if _, err := exec.LookPath("soffice"); err == nil {
		t.Skip("cannot simulate missing soffice on this host")
	}
	_, err := NewConverter(1, time.Second)
	if !errors.Is(err, ErrSofficeNotFound) {
		t.Fatalf("want ErrSofficeNotFound, got %v", err)
	}
}
