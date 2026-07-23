package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/rpc/internal/biz_helper/parse_file/progress"
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

// TestApplyRPAFieldsSkipEmpty verifies the authoritative-replace write-back
// (FR-014): non-zero amounts and parseable dates overwrite; zero/empty/unparseable
// values preserve any existing real or manual value on the row.
func TestApplyRPAFieldsSkipEmpty(t *testing.T) {
	tx := &model.BorrowerPropertyTransaction{
		PurchaseAmount: 1000000,
		EmdAmount:      2000,
	}
	applyRPAFields(tx, rpaFields{
		PurchaseAmount:         1750000, // overwrite
		CoeDate:                "09/15/2026",
		EmdAmount:              0,            // keep 2000
		EmdDueDate:             "",           // keep nil
		ContingencyRemovalDate: "not-a-date", // unparseable → skip
	})
	if tx.PurchaseAmount != 1750000 {
		t.Errorf("PurchaseAmount = %v, want 1750000", tx.PurchaseAmount)
	}
	if tx.EmdAmount != 2000 {
		t.Errorf("EmdAmount = %v, want 2000 (skip-zero)", tx.EmdAmount)
	}
	if tx.CoeDate == nil || !tx.CoeDate.Valid {
		t.Error("CoeDate must be set from parseable date")
	}
	if tx.EmdDueDate != nil {
		t.Error("EmdDueDate must stay nil (empty extracted)")
	}
	if tx.ContingencyRemovalDate != nil {
		t.Error("ContingencyRemovalDate must stay nil (unparseable extracted)")
	}
}

func TestParseDate(t *testing.T) {
	if t1, ok := parseDate("09/15/2026"); !ok || !t1.Valid {
		t.Fatalf("valid date not parsed: %+v", t1)
	}
	if _, ok := parseDate(""); ok {
		t.Fatal("empty string must not parse")
	}
	if _, ok := parseDate("garbage"); ok {
		t.Fatal("garbage must not parse")
	}
	_ = sql.NullTime{}
}

// memReporter is an in-memory handler.Reporter for extract-stage tests.
type memStage struct {
	status    string
	total     int
	processed int
	output    json.RawMessage
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
func (m *memReporter) StageOutput(_ context.Context, stage string) json.RawMessage {
	return m.stage(stage).output
}
func (m *memReporter) SaveStageOutput(_ context.Context, stage string, output json.RawMessage) {
	m.stage(stage).output = output
}

// TestExtractStageResumesFromCheckpoint verifies a restarted extract stage
// recalls the persisted envelope map and skips already-extracted files (no
// re-extraction / no LLM call) while returning their recalled content — so rpa
// synthesis across all files survives an interrupt. Files here are all
// pre-checkpointed, so extractOneFile (pdfium + LLM) is never reached and the
// handler's nil deps are safe.
func TestExtractStageResumesFromCheckpoint(t *testing.T) {
	rep := newMemReporter()
	pre := map[int64]rpaEnvelope{
		10: {FileId: 10, Fields: rpaFields{PurchaseAmount: 1750000}},
		11: {FileId: 11, Fields: rpaFields{CoeDate: "09/15/2026"}},
	}
	raw, _ := json.Marshal(pre)
	rep.SaveStageOutput(context.Background(), progress.StageExtract, raw)

	h := &RPAHandler{}
	got := h.extractStage(context.Background(), rep, []collectedFile{
		{FileId: 10, Name: "a.pdf"}, {FileId: 11, Name: "b.pdf"},
	})
	if len(got) != 2 {
		t.Fatalf("got %d envelopes, want 2 (recalled, not re-extracted)", len(got))
	}
	if got[0].FileId != 10 || got[1].FileId != 11 {
		t.Fatalf("order = [%d,%d], want fileId-sorted [10,11]", got[0].FileId, got[1].FileId)
	}
	if got[0].Fields.PurchaseAmount != 1750000 || got[1].Fields.CoeDate != "09/15/2026" {
		t.Fatalf("recalled content wrong: %+v", got)
	}
}

// seqChatLLM is an llm.Client whose ChatJSON returns one canned rpa JSON per
// call (in order), so a multi-file rpa run synthesizes across files without a
// real provider.
type seqChatLLM struct {
	responses []string
	n         int
}

func (s *seqChatLLM) ChatJSON(context.Context, string, string) (string, error) {
	if s.n >= len(s.responses) {
		return "{}", nil
	}
	r := s.responses[s.n]
	s.n++
	return r, nil
}

func (seqChatLLM) VisionJSON(context.Context, string, string, []byte) (string, error) {
	return "{}", nil
}

// TestRPAParseSynthesizesAndWritebacks drives the full manual_tree rpa pipeline
// (quickstart scenario B): collect subtree → extract each PDF (real pdfium text
// path, mock LLM) → skip-empty synthesis into the transaction row. Two files
// agree on purchase_amount and conflict on coe_date; write-back's
// first-non-empty-wins in fileId order takes the lower fileId's date, and an
// existing manual purchase_amount is authoritatively replaced.
func TestRPAParseSynthesizesAndWritebacks(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write(minimalPDF())
	}))
	t.Cleanup(srv.Close)

	db, repos := setupResolveDB(t)
	if err := db.Exec(`CREATE TABLE borrower_property_transactions (
		id INTEGER PRIMARY KEY AUTOINCREMENT, created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP, is_del INTEGER DEFAULT 0,
		purchase_amount REAL DEFAULT 0, coe_date DATETIME, emd_amount REAL DEFAULT 0,
		emd_due_date DATETIME, contingency_removal_date DATETIME
	)`).Error; err != nil {
		t.Fatalf("ddl tx: %v", err)
	}
	// contract branch 10 (root) + two PDF leaves under it.
	leaves := []int64{20, 21}
	if err := db.Exec(`INSERT INTO documents (id, entity_document_id, parent_id, type, name) VALUES (10, 1, 0, 'Branch', 'Contract')`).Error; err != nil {
		t.Fatalf("seed branch: %v", err)
	}
	for _, lid := range leaves {
		if err := db.Exec(`INSERT INTO documents (id, entity_document_id, parent_id, type, name) VALUES (?, 1, 10, 'Leaf', ?)`,
			lid, fmt.Sprintf("leaf-%d", lid)).Error; err != nil {
			t.Fatalf("seed leaf %d: %v", lid, err)
		}
		if err := db.Exec(`INSERT INTO files (id, entity_type, entity_id, file_url, file_name) VALUES (?, 'Document', ?, ?, ?)`,
			lid, lid, srv.URL, fmt.Sprintf("%d.pdf", lid)).Error; err != nil {
			t.Fatalf("seed file %d: %v", lid, err)
		}
	}
	// existing manual purchase_amount → must be authoritatively replaced.
	if err := db.Exec(`INSERT INTO borrower_property_transactions (id, purchase_amount) VALUES (7, 1000000)`).Error; err != nil {
		t.Fatalf("seed tx: %v", err)
	}

	llm := &seqChatLLM{responses: []string{
		`{"purchase_amount":1750000,"coe_date":"09/15/2026","emd_amount":5000}`, // leaf 20 (lower fileId)
		`{"purchase_amount":1750000,"coe_date":"10/01/2026","emd_amount":0}`,    // leaf 21 (conflict date)
	}}
	h := &RPAHandler{base{repos: repos, logger: logx.WithContext(context.Background()), llm: llm}}
	rule := dto.DocumentTemplateFileParseRule{Source: constants.DocumentTemplateFileParseSourceDirectory}

	if err := h.Parse(context.Background(), 7, 10, nil, rule, newMemReporter()); err != nil {
		t.Fatalf("parse: %v", err)
	}

	var tx model.BorrowerPropertyTransaction
	if err := db.First(&tx, 7).Error; err != nil {
		t.Fatalf("load tx: %v", err)
	}
	if tx.PurchaseAmount != 1750000 {
		t.Errorf("PurchaseAmount = %v, want 1750000 (agreed across files, replaced manual)", tx.PurchaseAmount)
	}
	if tx.EmdAmount != 5000 {
		t.Errorf("EmdAmount = %v, want 5000 (only leaf 20 carried it)", tx.EmdAmount)
	}
	// coe_date conflict: envelopes apply in ascending fileId order and each
	// non-empty value overwrites, so the HIGHEST fileId's value is final — the
	// deterministic fallback now that LLM arbitration is removed.
	if tx.CoeDate == nil || !tx.CoeDate.Valid {
		t.Fatalf("CoeDate must be set; got %+v", tx.CoeDate)
	}
	wantCOE, _ := time.Parse("01/02/2006", "10/01/2026")
	if !tx.CoeDate.Time.Equal(wantCOE) {
		t.Errorf("CoeDate = %v, want 10/01/2026 (highest fileId wins conflict)", tx.CoeDate.Time)
	}
}
