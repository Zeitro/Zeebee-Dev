package handler

import (
	"database/sql"
	"testing"

	"zeebee-backend/domains/zeebee/model"
)

func TestParseRPAFields(t *testing.T) {
	got, err := parseRPAFields(`{"purchase_amount":1750000,"coe_date":"2026-09-15","emd_amount":5000,"emd_due_date":null,"contingency_removal_date":"2026-08-01"}`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if got.PurchaseAmount != 1750000 || got.CoeDate != "2026-09-15" ||
		got.EmdAmount != 5000 || got.EmdDueDate != "" || got.ContingencyRemovalDate != "2026-08-01" {
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
		CoeDate:                "2026-09-15",
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

func TestParseISODate(t *testing.T) {
	if t1, ok := parseISODate("2026-09-15"); !ok || !t1.Valid {
		t.Fatalf("valid date not parsed: %+v", t1)
	}
	if _, ok := parseISODate(""); ok {
		t.Fatal("empty string must not parse")
	}
	if _, ok := parseISODate("garbage"); ok {
		t.Fatal("garbage must not parse")
	}
	_ = sql.NullTime{}
}
