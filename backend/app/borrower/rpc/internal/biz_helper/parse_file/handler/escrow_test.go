package handler

import (
	"context"
	"fmt"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"

	"zeebee-backend/common/llm"
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
}

func TestParseEscrowFieldsBadJSON(t *testing.T) {
	if _, err := parseEscrowFields("not json"); err == nil {
		t.Fatal("expected error on bad json")
	}
}

// TestApplyEscrowFieldsSkipEmpty: non-empty overwrites; empty preserves existing (FR-015).
func TestApplyEscrowFieldsSkipEmpty(t *testing.T) {
	esc := &model.PropertyTransactionEscrow{
		EscrowCompany: "Old Co", OfficerName: "Old Officer", OfficerEmail: "old@x.com",
		OfficerPhone: "111", No: "OLD-NO",
	}
	ApplyEscrowFields(esc, escrowFields{
		EscrowCompany: "New Co", OfficerName: "", OfficerEmail: "new@x.com",
		OfficerPhone: "", EscrowNumber: "NEW-NO",
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

// TestEscrowParseEndToEnd (T016): on_upload path — download → UploadFile → single
// consolidated Chat([fileId]) → skip-empty upsert. Mock file-client; no provider.
func TestEscrowParseEndToEnd(t *testing.T) {
	srv := byteServer(t, buildPDF(t, 1))
	db := setupEscrowDB(t)
	if err := db.Exec(`INSERT INTO property_transaction_escrows (id, borrower_property_transaction_id) VALUES (1, 7)`).Error; err != nil {
		t.Fatalf("seed: %v", err)
	}
	mock := &llm.MockClient{ChatContent: `{"escrow_company":"Mission Escrow, Inc.",` +
		`"officer_name":"Renee Lo","officer_email":"renee@mission.com",` +
		`"officer_phone":"4085136005","escrow_number":"202656-RL"}`}
	h := &EscrowHandler{base{
		repos:  repo.NewAllRepos(db),
		logger: logx.WithContext(context.Background()),
		llm:    mock,
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
	want := map[string]string{
		esc.EscrowCompany: "Mission Escrow, Inc.", esc.OfficerName: "Renee Lo",
		esc.OfficerEmail: "renee@mission.com", esc.OfficerPhone: "4085136005", esc.No: "202656-RL",
	}
	for got, w := range want {
		if got != w {
			t.Fatalf("escrow not upserted: %+v (want %s)", esc, w)
		}
	}
	if len(mock.UploadedFiles) != 1 {
		t.Fatalf("want 1 upload, got %v", mock.UploadedFiles)
	}
	atts := mock.Attachments()
	if len(atts) != 1 || atts[0].FileID == "" {
		t.Fatalf("want 1 chat attachment with the uploaded file id, got %+v", atts)
	}
}
