package biz_helper

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/datatypes"

	"zeebee-backend/domains/zeebee/model"
)

func docWithMeta(json string) *model.Document {
	return &model.Document{Metadata: datatypes.JSON([]byte(json))}
}

func TestReadParserName(t *testing.T) {
	cases := []struct {
		name string
		doc  *model.Document
		want string
	}{
		{"rpa", docWithMeta(`{"file_parse_config":{"parser":"rpa"}}`), "rpa"},
		{"escrow", docWithMeta(`{"file_parse_config":{"parser":"escrow_info"},"is_required":true}`), "escrow_info"},
		{"no config", docWithMeta(`{"is_required":true}`), ""},
		{"empty metadata", &model.Document{}, ""},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			assert.Equal(t, c.want, readParserName(c.doc))
		})
	}
}

func TestEscrowFieldsComplete(t *testing.T) {
	full := escrowFields{EscrowCompany: "a", OfficerName: "b", OfficerEmail: "c", OfficerPhone: "d", EscrowNumber: "e"}
	assert.True(t, full.complete())
	// missing the escrow number -> not complete (drives early-exit)
	assert.False(t, escrowFields{EscrowCompany: "a", OfficerName: "b", OfficerEmail: "c", OfficerPhone: "d"}.complete())
	assert.False(t, escrowFields{}.complete())
}

func TestApplyEscrowFields_SkipEmpty(t *testing.T) {
	esc := &model.PropertyTransactionEscrow{
		No: "OLD", EscrowCompany: "OLDCO", OfficerName: "OLDNAME",
		OfficerPhone: "OLDPHONE", OfficerEmail: "OLDEMAIL",
	}
	// only EscrowNumber + OfficerName extracted; empty fields must NOT overwrite.
	applyEscrowFields(esc, escrowFields{EscrowNumber: "NEW", OfficerName: "NEWNAME"})

	assert.Equal(t, "NEW", esc.No)
	assert.Equal(t, "NEWNAME", esc.OfficerName)
	assert.Equal(t, "OLDCO", esc.EscrowCompany, "empty extraction must preserve existing")
	assert.Equal(t, "OLDPHONE", esc.OfficerPhone, "empty extraction must preserve existing")
	assert.Equal(t, "OLDEMAIL", esc.OfficerEmail, "empty extraction must preserve existing")
}

func TestIsPDF(t *testing.T) {
	assert.True(t, isPDF("offer.pdf"))
	assert.True(t, isPDF("OFFER.PDF"))
	assert.False(t, isPDF("offer.docx"))
	assert.False(t, isPDF("scan.png"))
	assert.False(t, isPDF(""))
}

func TestParsePurchasePrice(t *testing.T) {
	price, err := parsePurchasePrice(`{"purchase_price": 1750000}`)
	assert.NoError(t, err)
	assert.Equal(t, 1750000.0, price)

	// absent -> 0 (drives the "keep existing" branch in parseRPA)
	price, err = parsePurchasePrice(`{"purchase_price": 0}`)
	assert.NoError(t, err)
	assert.Equal(t, 0.0, price)

	_, err = parsePurchasePrice(`not json`)
	assert.Error(t, err)
}

func TestParseEscrowFields(t *testing.T) {
	raw := `{"escrow_company":"Mission Escrow, Inc.","officer_name":"Renee Lo","officer_email":"r@x.com","officer_phone":"4085136005","escrow_number":"202656-RL"}`
	f, err := parseEscrowFields(raw)
	assert.NoError(t, err)
	assert.Equal(t, "Mission Escrow, Inc.", f.EscrowCompany)
	assert.Equal(t, "202656-RL", f.EscrowNumber)
	assert.True(t, f.complete())

	// all-null extraction -> empty fields, not complete
	f, err = parseEscrowFields(`{"escrow_company":null,"officer_name":null,"officer_email":null,"officer_phone":null,"escrow_number":null}`)
	assert.NoError(t, err)
	assert.False(t, f.complete())

	_, err = parseEscrowFields(`not json`)
	assert.Error(t, err)
}
