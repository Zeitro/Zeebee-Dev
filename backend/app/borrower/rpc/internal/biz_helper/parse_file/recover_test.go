package parse_file

import (
	"context"
	"encoding/json"
	"fmt"
	"sync/atomic"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/glebarez/sqlite"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"

	"zeebee-backend/app/borrower/rpc/internal/biz_helper/parse_file/handler"
	"zeebee-backend/app/borrower/rpc/internal/biz_helper/parse_file/progress"
	"zeebee-backend/common/llm"
	"zeebee-backend/common/redisx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/repo"
)

// These tests cover the recovery coordinator in recover.go (recoverJobs +
// tryReclaim) over a real AllRepos backed by in-memory SQLite + miniredis.
// tryReclaim only reads the document (+ its files) and re-derives the rule, so
// only those two tables are needed. The cross-node SetNX uniqueness the
// recovery relies on is covered by progress.TestRecoveryScan + TestHoldSetNX;
// process kill/restart is inherently a manual integration check.

// countingHandler is a non-blocking handler.Handler that records how many times
// Parse ran, so a test can assert a job was (or was not) reclaimed.
type countingHandler struct {
	calls int32
}

func (c *countingHandler) Stages([]model.File) []progress.StageState {
	return []progress.StageState{{Name: progress.StageCollect, Status: progress.StageStatusPending}}
}

func (c *countingHandler) Parse(ctx context.Context, _, _ int64, _ []model.File, _ dto.DocumentTemplateFileParseRule, rep handler.Reporter) error {
	atomic.AddInt32(&c.calls, 1)
	rep.MarkStage(ctx, progress.StageCollect, progress.StageStatusDone)
	return nil
}

func setupRecoverDB(t *testing.T) (*gorm.DB, *repo.AllRepos) {
	t.Helper()
	dsn := fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("sqlite: %v", err)
	}
	ddl := []string{
		`CREATE TABLE documents (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			is_del INTEGER DEFAULT 0,
			parent_id INTEGER DEFAULT 0,
			entity_document_id INTEGER DEFAULT 0,
			type TEXT DEFAULT '',
			name TEXT DEFAULT '',
			level INTEGER DEFAULT 0,
			sort_order INTEGER DEFAULT 9,
			last_uploaded_at DATETIME,
			metadata TEXT DEFAULT '{}'
		)`,
		`CREATE TABLE files (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			is_del INTEGER DEFAULT 0,
			entity_type TEXT DEFAULT 'LoanDocument',
			entity_id INTEGER DEFAULT 0,
			file_url TEXT DEFAULT '',
			file_name TEXT DEFAULT '',
			storage_type TEXT DEFAULT 'Aws S3',
			external_id TEXT DEFAULT '',
			upload_type TEXT DEFAULT 'Manual',
			drive_url TEXT DEFAULT '',
			drive_synced_at DATETIME
		)`,
	}
	for _, s := range ddl {
		if err := db.Exec(s).Error; err != nil {
			t.Fatalf("ddl: %v", err)
		}
	}
	return db, repo.NewAllRepos(db)
}

func setupRecoverParser(t *testing.T) (*gorm.DB, *Parser, *progress.ProgressManager) {
	db, repos := setupRecoverDB(t)
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("miniredis: %v", err)
	}
	t.Cleanup(mr.Close)
	rdb, err := redisx.GetRedis(redisx.Opt{Host: mr.Addr(), Type: "node"})
	if err != nil {
		t.Fatalf("redisx: %v", err)
	}
	pm := progress.NewProgressManager(rdb, logx.WithContext(context.Background()))
	return db, NewParser(repos, &llm.MockClient{}, pm, logx.WithContext(context.Background())), pm
}

// seedParseDoc inserts a document whose metadata resolves to the given parser,
// so recovery can re-derive its rule.
func seedParseDoc(t *testing.T, db *gorm.DB, id int64, parser string) {
	t.Helper()
	meta, _ := json.Marshal(map[string]any{
		"file_parse_config": map[string]any{
			"rules": []map[string]any{{"parser": parser, "source": "Directory"}},
		},
	})
	if err := db.Exec(
		`INSERT INTO documents (id, type, name, metadata) VALUES (?, 'Branch', 'Contract', ?)`,
		id, string(meta),
	).Error; err != nil {
		t.Fatalf("seed doc: %v", err)
	}
}

// TestRecoverReclaimsAbandonedJob verifies an in-flight job with no done record
// and no live hold is reclaimed on the recovery scan: the hold is acquired, the
// rule re-derived from the document, and the run completes → terminal Finished.
func TestRecoverReclaimsAbandonedJob(t *testing.T) {
	db, p, pm := setupRecoverParser(t)
	fake := &countingHandler{}
	p.Register(constants.DocumentTemplateFileParseParserRPA, fake)
	seedParseDoc(t, db, 100, string(constants.DocumentTemplateFileParseParserRPA))

	ctx := context.Background()
	job := &progress.Job{
		RequestId: "req-recover", TransactionId: 7, RuleSig: "RPA|100",
		Parser: string(constants.DocumentTemplateFileParseParserRPA), Mode: "ManualTree", DocumentId: 100,
	}
	if err := pm.WriteJob(ctx, job); err != nil {
		t.Fatalf("write job: %v", err)
	}

	p.recoverJobs(ctx)

	done := waitDone(t, pm, "req-recover")
	if done.Status != progress.StatusFinished {
		t.Fatalf("recovered done = %+v, want Finished", done)
	}
	if atomic.LoadInt32(&fake.calls) != 1 {
		t.Fatalf("handler calls = %d, want 1 (reclaim re-ran the job)", fake.calls)
	}
}

// TestRecoverSkipsTerminalJob verifies a job that already reached a terminal
// done record is NOT reclaimed — no second run, no handler call.
func TestRecoverSkipsTerminalJob(t *testing.T) {
	db, p, pm := setupRecoverParser(t)
	fake := &countingHandler{}
	p.Register(constants.DocumentTemplateFileParseParserRPA, fake)
	seedParseDoc(t, db, 200, string(constants.DocumentTemplateFileParseParserRPA))

	ctx := context.Background()
	if err := pm.WriteJob(ctx, &progress.Job{
		RequestId: "req-done", TransactionId: 7, RuleSig: "RPA|200",
		Parser: string(constants.DocumentTemplateFileParseParserRPA), Mode: "ManualTree", DocumentId: 200,
	}); err != nil {
		t.Fatalf("write job: %v", err)
	}
	if err := pm.WriteDone(ctx, "req-done", progress.StatusFinished, ""); err != nil {
		t.Fatalf("write done: %v", err)
	}

	p.recoverJobs(ctx)

	if atomic.LoadInt32(&fake.calls) != 0 {
		t.Fatalf("handler calls = %d, want 0 (terminal job must not be reclaimed)", fake.calls)
	}
}
