package handler

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"

	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/repo"
)

func TestParseEscrowFields(t *testing.T) {
	got, err := parseEscrowFields(`{"escrow_company":"Mission Escrow","officer_name":"Renee Lo","officer_email":null,"officer_phone":"4085136005","escrow_number":"202656-RL"}`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if got.EscrowCompany != "Mission Escrow" || got.OfficerName != "Renee Lo" || got.OfficerEmail != "" || got.OfficerPhone != "4085136005" || got.EscrowNumber != "202656-RL" {
		t.Fatalf("fields = %+v", got)
	}
	if got.OfficerEmail != "" {
		t.Fatal("null email must deserialize to empty")
	}
}

func TestParseEscrowFieldsBadJSON(t *testing.T) {
	if _, err := parseEscrowFields("not json"); err == nil {
		t.Fatal("expected error on bad json")
	}
}

// TestApplyEscrowFieldsSkipEmpty verifies the skip-empty merge rule (FR-015):
// non-empty extracted values overwrite; empty extracted values preserve the
// existing real value on the row.
func TestApplyEscrowFieldsSkipEmpty(t *testing.T) {
	esc := &model.PropertyTransactionEscrow{
		EscrowCompany: "Old Co",
		OfficerName:   "Old Officer",
		OfficerEmail:  "old@x.com",
		OfficerPhone:  "111",
		No:            "OLD-NO",
	}
	ApplyEscrowFields(esc, escrowFields{
		EscrowCompany: "New Co", // overwrite
		OfficerName:   "",       // keep Old Officer
		OfficerEmail:  "new@x.com",
		OfficerPhone:  "",
		EscrowNumber:  "NEW-NO",
	})
	if esc.EscrowCompany != "New Co" {
		t.Errorf("EscrowCompany = %q, want New Co", esc.EscrowCompany)
	}
	if esc.OfficerName != "Old Officer" {
		t.Errorf("OfficerName = %q, want Old Officer (skip-empty)", esc.OfficerName)
	}
	if esc.OfficerEmail != "new@x.com" {
		t.Errorf("OfficerEmail = %q, want new@x.com", esc.OfficerEmail)
	}
	if esc.OfficerPhone != "111" {
		t.Errorf("OfficerPhone = %q, want 111 (skip-empty)", esc.OfficerPhone)
	}
	if esc.No != "NEW-NO" {
		t.Errorf("No = %q, want NEW-NO", esc.No)
	}
}

// minimalPDF builds a valid minimal 1-page (blank) PDF with a correct xref so
// pdfium opens + renders it. Used only to drive the vision-extract glue
// end-to-end; the mock LLM ignores the rendered image.
func minimalPDF() []byte {
	var b bytes.Buffer
	b.WriteString("%PDF-1.4\n")
	o1 := b.Len()
	b.WriteString("1 0 obj\n<< /Type /Catalog /Pages 2 0 R >>\nendobj\n")
	o2 := b.Len()
	b.WriteString("2 0 obj\n<< /Type /Pages /Kids [3 0 R] /Count 1 >>\nendobj\n")
	o3 := b.Len()
	b.WriteString("3 0 obj\n<< /Type /Page /Parent 2 0 R /MediaBox [0 0 612 792] >>\nendobj\n")
	xref := b.Len()
	b.WriteString("xref\n0 4\n0000000000 65535 f \n")
	fmt.Fprintf(&b, "%010d 00000 n \n", o1)
	fmt.Fprintf(&b, "%010d 00000 n \n", o2)
	fmt.Fprintf(&b, "%010d 00000 n \n", o3)
	fmt.Fprintf(&b, "trailer\n<< /Size 4 /Root 1 0 R >>\nstartxref\n%d\n%%%%EOF\n", xref)
	return b.Bytes()
}

// visionStubLLM is an llm.Client whose VisionJSON returns a canned escrow JSON,
// so the extract glue runs end-to-end without a real provider.
type visionStubLLM struct{ resp string }

func (visionStubLLM) ChatJSON(context.Context, string, string) (string, error) {
	return "{}", nil
}
func (l visionStubLLM) VisionJSON(context.Context, string, string, []byte) (string, error) {
	return l.resp, nil
}

// setupEscrowDB builds an in-memory SQLite DB with the escrows table.
func setupEscrowDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())), &gorm.Config{})
	if err != nil {
		t.Fatalf("sqlite: %v", err)
	}
	if err := db.Exec(`CREATE TABLE property_transaction_escrows (
		id INTEGER PRIMARY KEY AUTOINCREMENT, created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP, is_del INTEGER DEFAULT 0,
		borrower_property_transaction_id INTEGER DEFAULT 0, no TEXT DEFAULT '',
		escrow_company TEXT DEFAULT '', officer_name TEXT DEFAULT '',
		officer_phone TEXT DEFAULT '', officer_email TEXT DEFAULT '',
		checklist TEXT DEFAULT '{}'
	)`).Error; err != nil {
		t.Fatalf("ddl: %v", err)
	}
	return db
}

// TestEscrowParseEndToEnd drives the full on_upload escrow path over a real PDF
// (pdfium/wazero) + httptest download + mock vision LLM + SQLite upsert:
// resolveFiles → download → render first page → VisionJSON → parse → skip-empty
// upsert into the primary escrow row. Proves the escrow glue the unit tests
// above cover piece-by-piece actually wires together.
func TestEscrowParseEndToEnd(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write(minimalPDF())
	}))
	t.Cleanup(srv.Close)

	db := setupEscrowDB(t)
	// create_logic seeds an empty primary escrow row; mirror that.
	if err := db.Exec(`INSERT INTO property_transaction_escrows (id, borrower_property_transaction_id) VALUES (1, 7)`).Error; err != nil {
		t.Fatalf("seed: %v", err)
	}
	h := &EscrowHandler{base{
		repos:  repo.NewAllRepos(db),
		logger: logx.WithContext(context.Background()),
		llm: visionStubLLM{resp: `{"escrow_company":"Mission Escrow, Inc.",` +
			`"officer_name":"Renee Lo","officer_email":"renee@mission.com",` +
			`"officer_phone":"4085136005","escrow_number":"202656-RL"}`},
	}}
	files := []model.File{{Base: model.Base{Id: 1}, FileUrl: srv.URL, FileName: "wire.pdf"}}
	rule := dto.DocumentTemplateFileParseRule{Source: constants.DocumentTemplateFileParseSourceSingleFile}

	if err := h.Parse(context.Background(), 7, 100, files, rule, newMemReporter()); err != nil {
		t.Fatalf("parse: %v", err)
	}

	var esc model.PropertyTransactionEscrow
	if err := db.First(&esc, 1).Error; err != nil {
		t.Fatalf("load escrow: %v", err)
	}
	if esc.EscrowCompany != "Mission Escrow, Inc." {
		t.Fatalf("EscrowCompany = %q", esc.EscrowCompany)
	}
	if esc.OfficerName != "Renee Lo" {
		t.Fatalf("OfficerName = %q", esc.OfficerName)
	}
	if esc.OfficerEmail != "renee@mission.com" {
		t.Fatalf("OfficerEmail = %q", esc.OfficerEmail)
	}
	if esc.OfficerPhone != "4085136005" {
		t.Fatalf("OfficerPhone = %q", esc.OfficerPhone)
	}
	if esc.No != "202656-RL" {
		t.Fatalf("No = %q", esc.No)
	}
}
