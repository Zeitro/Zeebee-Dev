package handler

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/signintech/gopdf"
	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/rpc/internal/biz_helper/parse_file/progress"
	"zeebee-backend/common/llm"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
)

func TestParseRPAFields(t *testing.T) {
	got, err := parseRPAFields(`{"purchase_amount":1750000,"coe_date":"2026-09-15","emd_amount":5000,"emd_due_date":null,"contingency_removal_date":"2026-08-01"}`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if got.PurchaseAmount != 1750000 || got.CoeDate != "2026-09-15" || got.EmdAmount != 5000 || got.EmdDueDate != "" || got.ContingencyRemovalDate != "2026-08-01" {
		t.Fatalf("fields = %+v", got)
	}
	if !got.hasAny() {
		t.Fatal("non-zero fields must report hasAny")
	}
}

func TestParseRPAFieldsBadJSON(t *testing.T) {
	if _, err := parseRPAFields("not json"); err == nil {
		t.Fatal("expected error on bad json")
	}
}

func TestRPAFieldsHasAny(t *testing.T) {
	if (rpaFields{}).hasAny() {
		t.Fatal("all-zero fields must not report hasAny")
	}
	if !(rpaFields{EmdDueDate: "2026-09-15"}).hasAny() {
		t.Fatal("a single non-empty date must report hasAny")
	}
}

// TestApplyRPAFieldsSkipEmpty: non-zero amounts + parseable dates overwrite; zero/empty skip.
func TestApplyRPAFieldsSkipEmpty(t *testing.T) {
	tx := &model.BorrowerPropertyTransaction{PurchaseAmount: 1000000, EmdAmount: 2000}
	applyRPAFields(tx, rpaFields{
		PurchaseAmount:         1750000,
		CoeDate:                "09/15/2026",
		EmdAmount:              0,
		EmdDueDate:             "",
		ContingencyRemovalDate: "not-a-date",
	})
	if tx.PurchaseAmount != 1750000 {
		t.Errorf("PurchaseAmount = %v, want 1750000", tx.PurchaseAmount)
	}
	if tx.EmdAmount != 2000 {
		t.Errorf("EmdAmount = %v, want 2000 (skip-zero)", tx.EmdAmount)
	}
	if tx.CoeDate == nil || !tx.CoeDate.Valid {
		t.Error("CoeDate must be set")
	}
	if tx.EmdDueDate != nil {
		t.Error("EmdDueDate must stay nil (empty extracted)")
	}
	if tx.ContingencyRemovalDate != nil {
		t.Error("ContingencyRemovalDate must stay nil (unparseable)")
	}
}

func TestParseDate(t *testing.T) {
	if t1, ok := parseDate("09/15/2026"); !ok || !t1.Valid {
		t.Fatalf("valid date not parsed: %+v", t1)
	}
	if _, ok := parseDate(""); ok {
		t.Fatal("empty must not parse")
	}
	if _, ok := parseDate("garbage"); ok {
		t.Fatal("garbage must not parse")
	}
	_ = sql.NullTime{}
}

// memReporter is an in-memory handler.Reporter for handler tests.
type memStage struct {
	status    string
	total     int
	processed int
	state     json.RawMessage
	retry     progress.RetryState
	files     map[int64]progress.FileCheckpoint
}

type memReporter struct {
	stages map[string]*memStage
}

func newMemReporter() *memReporter { return &memReporter{stages: map[string]*memStage{}} }

func (m *memReporter) stage(name string) *memStage {
	s, ok := m.stages[name]
	if !ok {
		s = &memStage{}
		m.stages[name] = s
	}
	return s
}

func (m *memReporter) MarkStage(_ context.Context, stage, status string) {
	m.stage(stage).status = status
}
func (m *memReporter) SetStageTotal(_ context.Context, stage string, total int) {
	m.stage(stage).total = total
}
func (m *memReporter) SetStageProcessed(_ context.Context, stage string, processed int) {
	m.stage(stage).processed = processed
}
func (m *memReporter) SetStageState(_ context.Context, stage string, state json.RawMessage) {
	m.stage(stage).state = state
}
func (m *memReporter) StageRetry(_ context.Context, stage string) progress.RetryState {
	return m.stage(stage).retry
}
func (m *memReporter) SaveStageRetry(_ context.Context, stage string, retry progress.RetryState) {
	m.stage(stage).retry = retry
}
func (m *memReporter) SaveFile(_ context.Context, fileId int64, fc progress.FileCheckpoint) {
	s := m.stage(progress.StageExtract)
	if s.files == nil {
		s.files = map[int64]progress.FileCheckpoint{}
	}
	s.files[fileId] = fc
}
func (m *memReporter) Files(_ context.Context) map[int64]progress.FileCheckpoint {
	return m.stage(progress.StageExtract).files
}

// buildPDF renders an N empty-page PDF for parse tests (pdfium needs real input).
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

// byteServer serves fixed bytes over http for download-driven tests.
func byteServer(t *testing.T, b []byte) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write(b)
	}))
	t.Cleanup(srv.Close)
	return srv
}

// TestRPAParseConsolidated drives the manual_tree path: collect subtree → upload
// both PDFs → ONE consolidated call (in-call arbitration) → authoritative write-back.
func TestRPAParseConsolidated(t *testing.T) {
	srv := byteServer(t, buildPDF(t, 2))
	db, repos := setupResolveDB(t)
	if err := db.Exec(`CREATE TABLE borrower_property_transactions (
		id INTEGER PRIMARY KEY AUTOINCREMENT, created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP, is_del INTEGER DEFAULT 0,
		purchase_amount REAL DEFAULT 0, coe_date DATETIME, emd_amount REAL DEFAULT 0,
		emd_due_date DATETIME, contingency_removal_date DATETIME
	)`).Error; err != nil {
		t.Fatalf("ddl tx: %v", err)
	}
	if err := db.Exec(`INSERT INTO documents (id, entity_document_id, parent_id, type, name) VALUES (10, 1, 0, 'Branch', 'Contract')`).Error; err != nil {
		t.Fatalf("seed branch: %v", err)
	}
	for _, lid := range []int64{20, 21} {
		if err := db.Exec(`INSERT INTO documents (id, entity_document_id, parent_id, type, name) VALUES (?, 1, 10, 'Leaf', ?)`, lid, fmt.Sprintf("leaf-%d", lid)).Error; err != nil {
			t.Fatalf("seed leaf %d: %v", lid, err)
		}
		if err := db.Exec(`INSERT INTO files (id, entity_type, entity_id, file_url, file_name) VALUES (?, 'Document', ?, ?, ?)`, lid, lid, srv.URL, fmt.Sprintf("%d.pdf", lid)).Error; err != nil {
			t.Fatalf("seed file %d: %v", lid, err)
		}
	}
	if err := db.Exec(`INSERT INTO borrower_property_transactions (id, purchase_amount) VALUES (7, 1000000)`).Error; err != nil {
		t.Fatalf("seed tx: %v", err)
	}

	mock := &llm.MockClient{ChatContent: `{"purchase_amount":1750000,"coe_date":"09/15/2026","emd_amount":5000}`}
	h := &RPAHandler{base{repos: repos, logger: logx.WithContext(context.Background()), llm: mock}}
	rule := dto.DocumentTemplateFileParseRule{Source: constants.DocumentTemplateFileParseSourceDirectory}

	if err := h.Parse(context.Background(), 7, 10, nil, rule, newMemReporter()); err != nil {
		t.Fatalf("parse: %v", err)
	}

	var tx model.BorrowerPropertyTransaction
	if err := db.First(&tx, 7).Error; err != nil {
		t.Fatalf("load tx: %v", err)
	}
	if tx.PurchaseAmount != 1750000 {
		t.Errorf("PurchaseAmount = %v, want 1750000 (authoritative replace)", tx.PurchaseAmount)
	}
	if tx.EmdAmount != 5000 {
		t.Errorf("EmdAmount = %v, want 5000", tx.EmdAmount)
	}
	if tx.CoeDate == nil || !tx.CoeDate.Valid {
		t.Fatalf("CoeDate must be set")
	}
	// one consolidated call over both uploaded files
	atts := mock.Attachments()
	if len(atts) != 2 {
		t.Errorf("want 2 attachments in the single consolidated call, got %d", len(atts))
	}
	if len(mock.UploadedFiles) != 2 {
		t.Errorf("want 2 uploads, got %d", len(mock.UploadedFiles))
	}
}

// TestUploadOneReusesCachedFileID: a pre-checkpointed matching-provider id is
// reused (no re-upload) — the resume reuse path.
func TestUploadOneReusesCachedFileID(t *testing.T) {
	mock := &llm.MockClient{ChatContent: "{}"}
	b := &base{logger: logx.WithContext(context.Background()), llm: mock}
	rep := newMemReporter()
	rep.SaveFile(context.Background(), 42, progress.FileCheckpoint{Provider: "mock", ProviderFileID: "existing-1"})

	fc, err := b.uploadOne(context.Background(), rep, collectedFile{FileId: 42, Name: "x.pdf"}, 0)
	if err != nil {
		t.Fatalf("uploadOne: %v", err)
	}
	if fc.ProviderFileID != "existing-1" {
		t.Fatalf("expected reuse of existing-1, got %s", fc.ProviderFileID)
	}
	if len(mock.UploadedFiles) != 0 {
		t.Fatalf("expected no re-upload, got %v", mock.UploadedFiles)
	}
}

// TestExtractStageIncludesWord (US2): a mixed batch (docx + png) — the docx is
// converted to pdf and uploaded as <base>.pdf; both reach the consolidated call.
func TestExtractStageIncludesWord(t *testing.T) {
	fc := &fakePdfConverter{pdf: buildPDF(t, 1)}
	mock := &llm.MockClient{ChatContent: "{}"}
	h := &RPAHandler{base{logger: logx.WithContext(context.Background()), llm: mock, pdfClient: fc}}
	rep := newMemReporter()
	files := []collectedFile{
		{FileId: 1, Name: "notes.docx", Url: "https://s3/notes.docx"},
		{FileId: 2, Name: "scan.png", Url: byteServer(t, []byte("png")).URL},
	}

	fields, err := h.extractStage(context.Background(), rep, files)
	if err != nil {
		t.Fatalf("extractStage: %v", err)
	}
	if fields.hasAny() {
		t.Fatal("empty mock response → no fields")
	}
	if len(mock.UploadedFiles) != 2 {
		t.Fatalf("want 2 uploads (docx→pdf + png), got %v", mock.UploadedFiles)
	}
	if len(mock.UploadedNames) != 2 || mock.UploadedNames[0] != "notes.pdf" {
		t.Errorf("uploaded names = %v, want notes.pdf first", mock.UploadedNames)
	}
	atts := mock.Attachments()
	if len(atts) != 2 {
		t.Fatalf("want 2 consolidated attachments, got %d", len(atts))
	}
	if atts[0].Filename != "notes.pdf" || atts[0].ContentType != "application/pdf" {
		t.Errorf("docx attachment = %+v, want notes.pdf/application/pdf", atts[0])
	}
}

// TestUploadOneReuploadsOnProviderMismatch: a cached id under a different provider
// is NOT reused (provider switch invalidates old ids) — the file is re-uploaded.
func TestUploadOneReuploadsOnProviderMismatch(t *testing.T) {
	srv := byteServer(t, []byte("pdf bytes"))
	mock := &llm.MockClient{} // Provider() == "mock"
	b := &base{logger: logx.WithContext(context.Background()), llm: mock}
	rep := newMemReporter()
	rep.SaveFile(context.Background(), 42, progress.FileCheckpoint{Provider: "claude", ProviderFileID: "old-claude-id"})

	fc, err := b.uploadOne(context.Background(), rep, collectedFile{FileId: 42, Name: "x.pdf", Url: srv.URL}, 0)
	if err != nil {
		t.Fatalf("uploadOne: %v", err)
	}
	if fc.ProviderFileID == "old-claude-id" {
		t.Fatal("must not reuse id from a different provider")
	}
	if len(mock.UploadedFiles) != 1 {
		t.Fatalf("want 1 re-upload on provider mismatch, got %v", mock.UploadedFiles)
	}
}

// TestUploadOneOpenAIProviderSlices: under an openai-named provider, a >cap PDF
// is still sliced before upload and the checkpoint records provider + cap. Proves
// the slicing path is provider-agnostic through the changed uploadOne.
func TestUploadOneOpenAIProviderSlices(t *testing.T) {
	srv := byteServer(t, buildPDF(t, 11))
	mock := &llm.MockClient{UploadProvider: "openai", ChatContent: "{}"}
	b := &base{logger: logx.WithContext(context.Background()), llm: mock}
	rep := newMemReporter()

	fc, err := b.uploadOne(context.Background(), rep, collectedFile{FileId: 1, Name: "big.pdf", Url: srv.URL}, maxRPAPages)
	if err != nil {
		t.Fatalf("uploadOne: %v", err)
	}
	if fc.Provider != "openai" {
		t.Errorf("Provider = %s, want openai", fc.Provider)
	}
	if fc.MaxPages != maxRPAPages {
		t.Errorf("MaxPages = %d, want %d", fc.MaxPages, maxRPAPages)
	}
	if len(mock.UploadedFiles) != 1 {
		t.Fatalf("want 1 sliced upload, got %v", mock.UploadedFiles)
	}
	// cached now → second call reuses (same provider + cap), no re-upload.
	if _, err := b.uploadOne(context.Background(), rep, collectedFile{FileId: 1, Name: "big.pdf"}, maxRPAPages); err != nil {
		t.Fatalf("reuse: %v", err)
	}
	if len(mock.UploadedFiles) != 1 {
		t.Fatalf("want reuse (still 1 upload), got %v", mock.UploadedFiles)
	}
}
