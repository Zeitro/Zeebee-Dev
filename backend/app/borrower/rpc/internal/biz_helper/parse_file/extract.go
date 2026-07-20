package parse_file

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/klippa-app/go-pdfium"

	"zeebee-backend/common/llm"
	"zeebee-backend/domains/zeebee/model"
)

// escrowFields is the vision-extracted escrow row. Empty fields mean "absent"
// and are skipped on write-back (skip-empty merge).
type escrowFields struct {
	EscrowCompany string `json:"escrow_company"`
	OfficerName   string `json:"officer_name"`
	OfficerEmail  string `json:"officer_email"`
	OfficerPhone  string `json:"officer_phone"`
	EscrowNumber  string `json:"escrow_number"`
}

func (e escrowFields) complete() bool {
	return e.EscrowCompany != "" && e.OfficerName != "" && e.OfficerEmail != "" &&
		e.OfficerPhone != "" && e.EscrowNumber != ""
}

func parseEscrowFields(content string) (escrowFields, error) {
	var f escrowFields
	if err := json.Unmarshal([]byte(content), &f); err != nil {
		return f, fmt.Errorf("escrow response not json: %s", content)
	}
	return f, nil
}

// extractEscrowFields vision-reads escrow fields off the rendered first page.
func extractEscrowFields(ctx context.Context, client llm.Client, r io.ReadSeeker) (fields escrowFields, err error) {
	err = withPdfInstance(func(inst pdfium.Pdfium) error {
		doc, cleanup, e := openPDF(inst, r)
		if e != nil {
			return e
		}
		defer cleanup()
		png, e := renderPagePNG(inst, doc, 0)
		if e != nil {
			return e
		}
		content, e := client.VisionJSON(ctx, extractSystemPrompt, escrowPrompt(), png)
		if e != nil {
			return e
		}
		fields, e = parseEscrowFields(content)
		return e
	})
	return fields, err
}

// applyEscrowFields copies only non-empty extracted fields (skip-empty merge),
// preserving any real value already on the row.
func applyEscrowFields(esc *model.PropertyTransactionEscrow, f escrowFields) {
	if f.EscrowCompany != "" {
		esc.EscrowCompany = f.EscrowCompany
	}
	if f.OfficerName != "" {
		esc.OfficerName = f.OfficerName
	}
	if f.OfficerEmail != "" {
		esc.OfficerEmail = f.OfficerEmail
	}
	if f.OfficerPhone != "" {
		esc.OfficerPhone = f.OfficerPhone
	}
	if f.EscrowNumber != "" {
		esc.No = f.EscrowNumber
	}
}

// ---- rpa (v1 single-field; US2 extends to 5-field rpaFields/extractRPAFields) ----

func parsePurchasePrice(content string) (float64, error) {
	var probe struct {
		PurchasePrice float64 `json:"purchase_price"`
	}
	if err := json.Unmarshal([]byte(content), &probe); err != nil {
		return 0, fmt.Errorf("rpa response not json: %s", content)
	}
	return probe.PurchasePrice, nil
}

// extractPurchasePrice pulls purchase_price (float) from page-1 text. 0 = absent.
func extractPurchasePrice(ctx context.Context, client llm.Client, r io.ReadSeeker) (price float64, err error) {
	err = withPdfInstance(func(inst pdfium.Pdfium) error {
		doc, cleanup, e := openPDF(inst, r)
		if e != nil {
			return e
		}
		defer cleanup()
		text, e := readPageText(inst, doc, 0)
		if e != nil {
			return e
		}
		content, e := client.ChatJSON(ctx, extractSystemPrompt, rpaPrompt(text))
		if e != nil {
			return e
		}
		price, e = parsePurchasePrice(content)
		return e
	})
	return price, err
}
