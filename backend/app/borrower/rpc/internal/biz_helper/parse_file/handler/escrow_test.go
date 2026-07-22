package handler

import (
	"testing"

	"zeebee-backend/domains/zeebee/model"
)

func TestParseEscrowFields(t *testing.T) {
	got, err := parseEscrowFields(`{"escrow_company":"Mission Escrow","officer_name":"Renee Lo","officer_email":null,"officer_phone":"4085136005","escrow_number":"202656-RL"}`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if got.EscrowCompany != "Mission Escrow" || got.OfficerName != "Renee Lo" ||
		got.OfficerEmail != "" || got.OfficerPhone != "4085136005" || got.EscrowNumber != "202656-RL" {
		t.Fatalf("fields = %+v", got)
	}
	if got.complete() {
		t.Fatal("must not be complete with empty email")
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

func TestEscrowFieldsComplete(t *testing.T) {
	full := escrowFields{EscrowCompany: "c", OfficerName: "n", OfficerEmail: "e", OfficerPhone: "p", EscrowNumber: "x"}
	if !full.complete() {
		t.Fatal("full fields must be complete")
	}
	empty := escrowFields{}
	if empty.complete() {
		t.Fatal("empty fields must not be complete")
	}
}
