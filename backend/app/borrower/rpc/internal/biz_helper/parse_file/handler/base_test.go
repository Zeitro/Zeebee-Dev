package handler

import (
	"context"
	"fmt"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"

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
		id       int64
		parent   int64
		typ      string
		meta string
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
