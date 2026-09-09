package handler

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"

	"zeebee-backend/common/llm"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/repo"
)

// These tests cover base.resolveFiles — the DB wiring (loadTxDocs →
// GetByEntity + GetAllByEntityDocumentId, then gatherFiles → GetFiles) — for
// the Directory (manual_tree) and TagFiltered (manual_select) scopes. The pure
// CollectSubtree/CollectBySelector selection is covered in selector_test.go;
// here we prove the repos feed those selectors correctly over in-memory SQLite.

func setupResolveDB(t *testing.T) (*gorm.DB, *repo.AllRepos) {
	t.Helper()
	dsn := fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("sqlite: %v", err)
	}
	for _, s := range []string{
		`CREATE TABLE entity_documents (
			id INTEGER PRIMARY KEY AUTOINCREMENT, created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP, is_del INTEGER DEFAULT 0,
			organization_id INTEGER DEFAULT 0, entity_type TEXT DEFAULT '', entity_id INTEGER DEFAULT 0,
			root_document_id INTEGER DEFAULT 0
		)`,
		`CREATE TABLE documents (
			id INTEGER PRIMARY KEY AUTOINCREMENT, created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP, is_del INTEGER DEFAULT 0,
			parent_id INTEGER DEFAULT 0, entity_document_id INTEGER DEFAULT 0, type TEXT DEFAULT '',
			name TEXT DEFAULT '', level INTEGER DEFAULT 0, sort_order INTEGER DEFAULT 9,
			last_uploaded_at DATETIME, metadata TEXT DEFAULT '{}'
		)`,
		`CREATE TABLE files (
			id INTEGER PRIMARY KEY AUTOINCREMENT, created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP, is_del INTEGER DEFAULT 0,
			entity_type TEXT DEFAULT 'LoanDocument', entity_id INTEGER DEFAULT 0, file_url TEXT DEFAULT '',
			file_name TEXT DEFAULT '', storage_type TEXT DEFAULT 'Aws S3', external_id TEXT DEFAULT '',
			upload_type TEXT DEFAULT 'Manual', drive_url TEXT DEFAULT '', drive_synced_at DATETIME
		)`,
	} {
		if err := db.Exec(s).Error; err != nil {
			t.Fatalf("ddl: %v", err)
		}
	}
	// entity_documents(id=1) → BorrowerPropertyTransaction 7, root document 10.
	if err := db.Exec(`INSERT INTO entity_documents (id, entity_type, entity_id, root_document_id) VALUES (1, 'BorrowerPropertyTransaction', 7, 10)`).Error; err != nil {
		t.Fatalf("seed entdoc: %v", err)
	}
	return db, repo.NewAllRepos(db)
}

// seedDocTree seeds a branch (10) with two leaves under entity_document 1:
// leaf 20 tagged "emd", leaf 21 tagged "signed", each carrying one file.
func seedDocTree(t *testing.T, db *gorm.DB) {
	t.Helper()
	docs := []struct {
		id     int64
		parent int64
		typ    string
		meta   string
	}{
		{10, 0, "Branch", "{}"},
		{20, 10, "Leaf", `{"tags":["emd"]}`},
		{21, 10, "Leaf", `{"tags":["signed"]}`},
	}
	for _, d := range docs {
		if err := db.Exec(`INSERT INTO documents (id, entity_document_id, parent_id, type, name, metadata) VALUES (?, 1, ?, ?, ?, ?)`,
			d.id, d.parent, d.typ, fmt.Sprintf("doc-%d", d.id), d.meta).Error; err != nil {
			t.Fatalf("seed doc %d: %v", d.id, err)
		}
	}
	for _, fid := range []int64{20, 21} {
		if err := db.Exec(`INSERT INTO files (id, entity_type, entity_id, file_url, file_name) VALUES (?, 'Document', ?, ?, ?)`,
			fid, fid, fmt.Sprintf("s3://%d.pdf", fid), fmt.Sprintf("%d.pdf", fid)).Error; err != nil {
			t.Fatalf("seed file %d: %v", fid, err)
		}
	}
}

// TestResolveFilesDirectory verifies manual_tree scope: the anchor branch's whole
// subtree of leaves is gathered (both tagged and untagged) → both files.
func TestResolveFilesDirectory(t *testing.T) {
	db, repos := setupResolveDB(t)
	seedDocTree(t, db)
	b := &base{repos: repos, logger: logx.WithContext(context.Background())}

	got, err := b.resolveFiles(context.Background(), 7, 10, nil, dto.DocumentTemplateFileParseRule{
		Source: constants.DocumentTemplateFileParseSourceDirectory,
	})
	if err != nil {
		t.Fatalf("resolveFiles: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("got %d files, want 2 (whole subtree): %+v", len(got), got)
	}
}

// TestResolveFilesTagFiltered verifies manual_select scope: only leaves whose
// tags ⊇ selector.tags are gathered → just the emd file, the signed one excluded.
func TestResolveFilesTagFiltered(t *testing.T) {
	db, repos := setupResolveDB(t)
	seedDocTree(t, db)
	b := &base{repos: repos, logger: logx.WithContext(context.Background())}

	got, err := b.resolveFiles(context.Background(), 7, 10, nil, dto.DocumentTemplateFileParseRule{
		Source:   constants.DocumentTemplateFileParseSourceTagFiltered,
		Selector: &dto.DocumentTemplateTagSelector{Tags: []string{"emd"}},
	})
	if err != nil {
		t.Fatalf("resolveFiles: %v", err)
	}
	if len(got) != 1 || got[0].FileId != 20 {
		t.Fatalf("got %+v, want only the emd-tagged file (id 20)", got)
	}
}

// fakePdfConverter is a stand-in pdfConverter that returns canned pdf bytes and
// records the DocToPdf args (no live pdf_generator needed).
type fakePdfConverter struct {
	pdf            []byte
	err            error
	gotTyp, gotExt string
	calls          int
}

func (f *fakePdfConverter) DocToPdf(_ context.Context, typ, _ string, _ []byte, ext string) ([]byte, error) {
	f.calls++
	f.gotTyp, f.gotExt = typ, ext
	return f.pdf, f.err
}

// TestFileStream_WordBranch: a .docx converts to application/pdf + <base>.pdf,
// with the reader carrying the converted bytes.
func TestFileStream_WordBranch(t *testing.T) {
	pdf := []byte("%PDF-1.4 fake body")
	fc := &fakePdfConverter{pdf: pdf}
	b := &base{pdfClient: fc, logger: logx.WithContext(context.Background())}

	res, err := b.fileStream(context.Background(), collectedFile{Name: "wire.docx", Url: "https://s3/wire.docx"})
	if err != nil {
		t.Fatalf("fileStream: %v", err)
	}
	if res.ContentType != "application/pdf" {
		t.Errorf("ContentType = %q, want application/pdf", res.ContentType)
	}
	if res.DisplayName != "wire.pdf" {
		t.Errorf("DisplayName = %q, want wire.pdf", res.DisplayName)
	}
	if fc.gotExt != "docx" || fc.gotTyp != "Url" {
		t.Errorf("DocToPdf args: typ=%q ext=%q", fc.gotTyp, fc.gotExt)
	}
	got, _ := io.ReadAll(res.Reader)
	if !bytes.Equal(got, pdf) {
		t.Errorf("reader bytes mismatch: %q", got)
	}
}

// TestFileStream_WordNoClient: absent pdf client → error (caller fail-soft skips).
func TestFileStream_WordNoClient(t *testing.T) {
	b := &base{logger: logx.WithContext(context.Background())}
	if _, err := b.fileStream(context.Background(), collectedFile{Name: "wire.docx"}); err == nil {
		t.Fatal("expected error when pdf client is nil")
	}
}

// TestMaxPagesForFile_Word: Word files get the cap (so cache key + slicing apply).
func TestMaxPagesForFile_Word(t *testing.T) {
	if maxPagesForFile("wire.docx", 1) != 1 {
		t.Error("docx should map to cap")
	}
	if maxPagesForFile("wire.doc", 5) != 5 {
		t.Error("doc should map to cap")
	}
	if maxPagesForFile("photo.png", 5) != 0 {
		t.Error("image should map to 0")
	}
}

// TestUploadOne_WordRename: uploading a Word file passes the .pdf display name
// to the provider (R6 — name/format must match the converted bytes).
func TestUploadOne_WordRename(t *testing.T) {
	fc := &fakePdfConverter{pdf: []byte("%PDF-1.4 body")}
	mc := &llm.MockClient{}
	b := &base{pdfClient: fc, llm: mc, logger: logx.WithContext(context.Background())}
	rep := newMemReporter()

	ckpt, err := b.uploadOne(context.Background(), rep, collectedFile{FileId: 1, Name: "wire.docx", Url: "s3://wire.docx"}, 0)
	if err != nil {
		t.Fatalf("uploadOne: %v", err)
	}
	if len(mc.UploadedNames) != 1 || mc.UploadedNames[0] != "wire.pdf" {
		t.Errorf("uploaded names = %v, want [wire.pdf]", mc.UploadedNames)
	}
	if ckpt.MaxPages != 0 { // maxPages 0 → whole file
		t.Errorf("checkpoint MaxPages = %d, want 0", ckpt.MaxPages)
	}
}
