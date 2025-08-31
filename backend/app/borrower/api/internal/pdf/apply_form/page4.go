package apply_form

import (
	"github.com/signintech/gopdf"

	"zeebee-backend/app/borrower/api/internal/pdf"
)

func (form *ApplyForm) page4() (err error) {
	// lineHeight := 18.0
	// fontHeight := 12.0
	lineHeight := form.pdf.LineHeight
	fontHeight := form.pdf.LineHeight
	form.pdf.AddPage()
	err = form.setPageTitle(pdf.GlassIco, "PROPERTY |", "LOAN INFO")
	if err != nil {
		return err
	}
	form.pdf.Br(lineHeight)
	form.pdf.SetBoldFont()
	_ = form.pdf.MultiCell(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: 12}, "Subject Property Address:")
	form.pdf.Br(lineHeight / 2.0)
	_ = form.pdf.Cell(&gopdf.Rect{W: 120, H: 12}, "City:")
	_ = form.pdf.Cell(&gopdf.Rect{W: 120, H: 12}, "State: ")
	_ = form.pdf.Cell(&gopdf.Rect{W: 120, H: 12}, "Zip:")
	form.pdf.Br(lineHeight)

	_ = form.pdf.MultiCell(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: 12}, "Loan Type:")
	form.pdf.Br(lineHeight / 2.0)
	x := form.pdf.GetX()
	y := form.pdf.GetY()
	form.pdf.SetRegularFont()
	_ = form.buildCheckbox(x, y, []string{"Bridge", "Rehab", "Rental", "Ground Up Construction"})

	form.pdf.Br(lineHeight / 2.0)
	form.pdf.SetBoldFont()
	_ = form.pdf.MultiCell(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: 12}, "Loan Purpose:")
	form.pdf.Br(lineHeight / 2.0)
	x = form.pdf.GetX()
	y = form.pdf.GetY()
	form.pdf.SetRegularFont()
	_ = form.buildCheckbox(x, y, []string{"Purchase", "Refinance", "Cash Out Refinance"})

	form.pdf.Br(lineHeight / 2.0)
	form.pdf.SetBoldFont()
	_ = form.pdf.MultiCell(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: 12}, "Property Type:")
	form.pdf.Br(lineHeight / 2.0)
	x = form.pdf.GetX()
	y = form.pdf.GetY()
	form.pdf.SetRegularFont()
	_ = form.buildCheckbox(x, y, []string{"Single Family Residence", "2 unit", "3 unit", "4 unit", "4+ Multi"})
	form.pdf.Br(lineHeight / 2.0)
	x = form.pdf.GetX()
	y = form.pdf.GetY()
	_ = form.buildCheckbox(x, y, []string{"Townhouse", "PUD/Condominium", "Mixed Use", "Commercial"})
	form.pdf.Br(lineHeight / 2.0)

	x = form.pdf.GetX()
	y = form.pdf.GetY()
	_ = form.buildTitleWithBg("PURCHASE LOAN DETAIL:", position{w: form.pdf.GetFullWidth(), h: 27, x: x, y: y}, position{w: form.pdf.GetFullWidth() - 100, h: 27, x: x + 10, y: y})
	form.pdf.Br(30)

	_ = form.pdf.Cell(&gopdf.Rect{W: form.pdf.GetFullWidth() / 2.0, H: fontHeight}, "Purchase Price:")
	_ = form.pdf.Cell(&gopdf.Rect{W: form.pdf.GetFullWidth() / 2.0, H: fontHeight}, "Assignment Fee Y/N:")
	form.pdf.Br(lineHeight)
	_ = form.pdf.Cell(&gopdf.Rect{W: form.pdf.GetFullWidth() / 2.0, H: fontHeight}, "Scheduled Close of Escrow Date: ")
	_ = form.pdf.Cell(&gopdf.Rect{W: form.pdf.GetFullWidth() / 2.0, H: fontHeight}, "Assignment Fee:")
	form.pdf.Br(lineHeight)
	_ = form.pdf.Cell(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: fontHeight}, "Initial Advance Requested:")
	form.pdf.Br(lineHeight)
	_ = form.pdf.Cell(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: fontHeight}, " Rehab/Construction Loan Amount Requested:")
	form.pdf.Br(lineHeight)
	_ = form.pdf.Cell(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: fontHeight}, "Total Loan Amount Requested:")
	form.pdf.Br(lineHeight)
	_ = form.pdf.Cell(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: fontHeight}, "Loan Term Requested:")
	form.pdf.Br(lineHeight)
	_ = form.pdf.Cell(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: fontHeight}, "Estimated After Repair Value (ARV):")
	form.pdf.Br(lineHeight)
	_ = form.pdf.Cell(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: fontHeight}, "Liquidity:")
	form.pdf.Br(lineHeight)
	_ = form.pdf.Cell(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: fontHeight}, "Exit Strategy:")
	form.pdf.Br(lineHeight)

	x = form.pdf.GetX()
	y = form.pdf.GetY()
	_ = form.buildTitleWithBg("REFINANCE LOAN DETAIL:", position{w: form.pdf.GetFullWidth(), h: 27, x: x, y: y}, position{w: form.pdf.GetFullWidth() - 100, h: 27, x: x + 10, y: y})
	form.pdf.Br(30)
	_ = form.pdf.Cell(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: fontHeight}, "Original Purchase Price:")
	form.pdf.Br(lineHeight)
	_ = form.pdf.Cell(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: fontHeight}, "Date of Purchase:")
	form.pdf.Br(lineHeight)
	_ = form.pdf.Cell(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: fontHeight}, "Initial Advance Requested:")
	form.pdf.Br(lineHeight)
	_ = form.pdf.Cell(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: fontHeight}, "Rehab/Construction Loan Amount Requested:")
	form.pdf.Br(lineHeight)
	_ = form.pdf.Cell(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: fontHeight}, "Total Loan Amount Requested:")
	form.pdf.Br(lineHeight)
	_ = form.pdf.Cell(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: fontHeight}, "Loan Term Requested:")
	form.pdf.Br(lineHeight)
	_ = form.pdf.Cell(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: fontHeight}, "Estimated Current Value (LTV):")
	form.pdf.Br(lineHeight)
	_ = form.pdf.Cell(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: fontHeight}, "Estimated After Repair Value (ARV):")
	form.pdf.Br(lineHeight)
	_ = form.pdf.Cell(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: fontHeight}, "Existing Loan Amount (Payoff):")
	form.pdf.Br(lineHeight)
	_ = form.pdf.Cell(&gopdf.Rect{W: form.pdf.GetFullWidth() / 2.0, H: fontHeight}, "Existing Lender: ")
	_ = form.pdf.Cell(&gopdf.Rect{W: form.pdf.GetFullWidth() / 2.0, H: fontHeight}, "Existing Loan Rate:")
	form.pdf.Br(lineHeight)
	_ = form.pdf.Cell(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: fontHeight}, "Existing Loan Maturity Date:")
	form.pdf.Br(lineHeight)
	_ = form.pdf.Cell(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: fontHeight}, "Subordinate Debt (Y/N-$):")
	form.pdf.Br(lineHeight)
	_ = form.pdf.Cell(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: fontHeight}, "Liquidity:")
	form.pdf.Br(lineHeight)
	_ = form.pdf.Cell(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: fontHeight}, "Exit Strategy:")
	form.pdf.Br(lineHeight)
	_ = form.pdf.Cell(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: fontHeight}, "Mid Construction (Y/N):")

	return err
}
