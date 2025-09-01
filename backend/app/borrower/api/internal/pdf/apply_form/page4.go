package apply_form

import (
	"strings"

	"github.com/duke-git/lancet/v2/condition"
	"github.com/signintech/gopdf"

	"zeebee-backend/app/borrower/api/internal/pdf"
	"zeebee-backend/common/constants"
	"zeebee-backend/common/stringx"
)

func (form *ApplyForm) page4() (err error) {
	lineHeight := form.pdf.LineHeight
	fontHeight := form.pdf.LineHeight
	form.pdf.AddPage()
	err = form.setPageTitle(pdf.GlassIco, "PROPERTY |", "LOAN INFO")
	if err != nil {
		return err
	}
	form.pdf.Br(lineHeight)
	form.pdf.SetBoldFont()
	loan := form.info
	property := loan.LoanProperty
	address := property.PropertyBase.Address
	_ = form.pdf.MultiCell(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: 12}, "Subject Property Address:  "+address.Address1)
	form.pdf.Br(lineHeight / 2.0)
	_ = form.pdf.Cell(&gopdf.Rect{W: 120, H: 12}, "City:  "+address.City)
	_ = form.pdf.Cell(&gopdf.Rect{W: 120, H: 12}, "State:  "+address.State)
	_ = form.pdf.Cell(&gopdf.Rect{W: 120, H: 12}, "Zip: "+address.ZipCode)
	form.pdf.Br(lineHeight)

	_ = form.pdf.MultiCell(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: 12}, "Loan Type:")
	form.pdf.Br(lineHeight / 2.0)
	x := form.pdf.GetX()
	y := form.pdf.GetY()
	form.pdf.SetRegularFont()
	err = form.buildRadioCheckbox(x, y, constants.LoanTypeName, loan.LoanType)
	if err != nil {
		return err
	}

	// _ = form.buildCheckbox(x, y, []string{"Bridge", "Rehab", "Rental", "Ground Up Construction"})

	form.pdf.Br(lineHeight / 2.0)
	form.pdf.SetBoldFont()
	_ = form.pdf.MultiCell(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: 12}, "Loan Purpose:")
	form.pdf.Br(lineHeight / 2.0)
	x = form.pdf.GetX()
	y = form.pdf.GetY()
	form.pdf.SetRegularFont()
	err = form.buildRadioCheckbox(x, y, constants.LoanPurposeName, loan.LoanPurpose)
	if err != nil {
		return err
	}
	// _ = form.buildCheckbox(x, y, []string{"Purchase", "Refinance", "Cash Out Refinance"})

	form.pdf.Br(lineHeight / 2.0)
	form.pdf.SetBoldFont()
	_ = form.pdf.MultiCell(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: 12}, "Property Type:")
	form.pdf.Br(lineHeight / 2.0)
	x = form.pdf.GetX()
	y = form.pdf.GetY()
	form.pdf.SetRegularFont()
	err = form.buildRadioCheckbox(x, y, constants.PropertyTypeName, property.PropertyType)
	if err != nil {
		return err
	}

	form.pdf.Br(lineHeight / 2.0)

	x = form.pdf.GetX()
	y = form.pdf.GetY()
	_ = form.buildTitleWithBg("PURCHASE LOAN DETAIL:", position{w: form.pdf.GetFullWidth(), h: 27, x: x, y: y}, position{w: form.pdf.GetFullWidth() - 100, h: 27, x: x + 10, y: y})
	form.pdf.Br(30)
	if constants.LoanPurpose(loan.LoanPurpose) == constants.LoanPurposePurchase {
		_ = form.pdf.Cell(&gopdf.Rect{W: form.pdf.GetFullWidth() / 2.0, H: fontHeight}, "Purchase Price:  "+stringx.TransNumberToString[float64](loan.PurchaseAmount, true))
		_ = form.pdf.Cell(&gopdf.Rect{W: form.pdf.GetFullWidth() / 2.0, H: fontHeight}, "Assignment Fee Y/N:  "+condition.Ternary(loan.HasAssignmentFee, "Y", "N"))
		form.pdf.Br(lineHeight)
		_ = form.pdf.Cell(&gopdf.Rect{W: form.pdf.GetFullWidth() / 2.0, H: fontHeight}, "Scheduled Close of Escrow Date:  "+loan.ScheduledCloseEscrowDate)
		_ = form.pdf.Cell(&gopdf.Rect{W: form.pdf.GetFullWidth() / 2.0, H: fontHeight}, "Assignment Fee:  "+stringx.TransNumberToString[float64](loan.AssignmentFeeAmount, true))
		form.pdf.Br(lineHeight)
		_ = form.pdf.Cell(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: fontHeight}, "Initial Advance Requested:  "+stringx.TransNumberToString[float64](loan.InitialAmount, true))
		form.pdf.Br(lineHeight)
		_ = form.pdf.Cell(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: fontHeight}, "Rehab/Construction Loan Amount Requested:  "+stringx.TransNumberToString[float64](loan.HoldbackAmount, true))
		form.pdf.Br(lineHeight)
		_ = form.pdf.Cell(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: fontHeight}, "Total Loan Amount Requested:  "+stringx.TransNumberToString[float64](loan.TotalAmount, true))
		form.pdf.Br(lineHeight)
		_ = form.pdf.Cell(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: fontHeight}, "Loan Term Requested:  "+stringx.TransNumberToString[int32](loan.LoanTerm, true))
		form.pdf.Br(lineHeight)
		_ = form.pdf.Cell(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: fontHeight}, "Estimated After Repair Value (ARV):  "+stringx.TransNumberToString[float64](loan.Arv, true))
		form.pdf.Br(lineHeight)
		_ = form.pdf.Cell(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: fontHeight}, "Liquidity:  "+stringx.TransNumberToString[float64](loan.LiquidityAmount, true))
		form.pdf.Br(lineHeight)
		_ = form.pdf.Cell(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: fontHeight}, "Exit Strategy:  "+strings.Join(loan.BusinessStrategy, ","))
		form.pdf.Br(lineHeight)

		x = form.pdf.GetX()
		y = form.pdf.GetY()
		_ = form.buildTitleWithBg("REFINANCE LOAN DETAIL:", position{w: form.pdf.GetFullWidth(), h: 27, x: x, y: y}, position{w: form.pdf.GetFullWidth() - 100, h: 27, x: x + 10, y: y})
		form.pdf.Br(30)
		_ = form.pdf.Cell(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: fontHeight}, "Original Purchase Price:  ")
		form.pdf.Br(lineHeight)
		_ = form.pdf.Cell(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: fontHeight}, "Date of Purchase:  ")
		form.pdf.Br(lineHeight)
		_ = form.pdf.Cell(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: fontHeight}, "Initial Advance Requested:  ")
		form.pdf.Br(lineHeight)
		_ = form.pdf.Cell(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: fontHeight}, "Rehab/Construction Loan Amount Requested:  ")
		form.pdf.Br(lineHeight)
		_ = form.pdf.Cell(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: fontHeight}, "Total Loan Amount Requested:  ")
		form.pdf.Br(lineHeight)
		_ = form.pdf.Cell(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: fontHeight}, "Loan Term Requested:  ")
		form.pdf.Br(lineHeight)
		_ = form.pdf.Cell(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: fontHeight}, "Estimated Current Value (LTV):  ")
		form.pdf.Br(lineHeight)
		_ = form.pdf.Cell(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: fontHeight}, "Estimated After Repair Value (ARV):  ")
		form.pdf.Br(lineHeight)
		_ = form.pdf.Cell(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: fontHeight}, "Existing Loan Amount (Payoff):  ")
		form.pdf.Br(lineHeight)
		_ = form.pdf.Cell(&gopdf.Rect{W: form.pdf.GetFullWidth() / 2.0, H: fontHeight}, "Existing Lender:  ")
		_ = form.pdf.Cell(&gopdf.Rect{W: form.pdf.GetFullWidth() / 2.0, H: fontHeight}, "Existing Loan Rate:  ")
		form.pdf.Br(lineHeight)
		_ = form.pdf.Cell(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: fontHeight}, "Existing Loan Maturity Date:  ")
		form.pdf.Br(lineHeight)
		_ = form.pdf.Cell(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: fontHeight}, "Subordinate Debt (Y/N-$):  ")
		form.pdf.Br(lineHeight)
		_ = form.pdf.Cell(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: fontHeight}, "Liquidity:  ")
		form.pdf.Br(lineHeight)
		_ = form.pdf.Cell(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: fontHeight}, "Exit Strategy:  ")
		form.pdf.Br(lineHeight)
		_ = form.pdf.Cell(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: fontHeight}, "Mid Construction (Y/N):  ")
	} else {
		_ = form.pdf.Cell(&gopdf.Rect{W: form.pdf.GetFullWidth() / 2.0, H: fontHeight}, "Purchase Price:  ")
		_ = form.pdf.Cell(&gopdf.Rect{W: form.pdf.GetFullWidth() / 2.0, H: fontHeight}, "Assignment Fee Y/N:  ")
		form.pdf.Br(lineHeight)
		_ = form.pdf.Cell(&gopdf.Rect{W: form.pdf.GetFullWidth() / 2.0, H: fontHeight}, "Scheduled Close of Escrow Date:  ")
		_ = form.pdf.Cell(&gopdf.Rect{W: form.pdf.GetFullWidth() / 2.0, H: fontHeight}, "Assignment Fee:  ")
		form.pdf.Br(lineHeight)
		_ = form.pdf.Cell(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: fontHeight}, "Initial Advance Requested:  ")
		form.pdf.Br(lineHeight)
		_ = form.pdf.Cell(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: fontHeight}, "Rehab/Construction Loan Amount Requested:  ")
		form.pdf.Br(lineHeight)
		_ = form.pdf.Cell(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: fontHeight}, "Total Loan Amount Requested:  ")
		form.pdf.Br(lineHeight)
		_ = form.pdf.Cell(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: fontHeight}, "Loan Term Requested:  ")
		form.pdf.Br(lineHeight)
		_ = form.pdf.Cell(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: fontHeight}, "Estimated After Repair Value (ARV):  ")
		form.pdf.Br(lineHeight)
		_ = form.pdf.Cell(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: fontHeight}, "Liquidity:  ")
		form.pdf.Br(lineHeight)
		_ = form.pdf.Cell(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: fontHeight}, "Exit Strategy:  ")
		form.pdf.Br(lineHeight)

		x = form.pdf.GetX()
		y = form.pdf.GetY()
		_ = form.buildTitleWithBg("REFINANCE LOAN DETAIL:", position{w: form.pdf.GetFullWidth(), h: 27, x: x, y: y}, position{w: form.pdf.GetFullWidth() - 100, h: 27, x: x + 10, y: y})
		form.pdf.Br(30)
		_ = form.pdf.Cell(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: fontHeight}, "Original Purchase Price:  "+stringx.TransNumberToString[float64](loan.PurchaseAmount, true))
		form.pdf.Br(lineHeight)
		_ = form.pdf.Cell(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: fontHeight}, "Date of Purchase:  "+loan.PurchaseDate)
		form.pdf.Br(lineHeight)
		_ = form.pdf.Cell(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: fontHeight}, "Initial Advance Requested:  "+stringx.TransNumberToString[float64](loan.InitialAmount, true))
		form.pdf.Br(lineHeight)
		_ = form.pdf.Cell(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: fontHeight}, "Rehab/Construction Loan Amount Requested:  "+stringx.TransNumberToString[float64](loan.HoldbackAmount, true))
		form.pdf.Br(lineHeight)
		_ = form.pdf.Cell(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: fontHeight}, "Total Loan Amount Requested:  "+stringx.TransNumberToString[float64](loan.TotalAmount, true))
		form.pdf.Br(lineHeight)
		_ = form.pdf.Cell(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: fontHeight}, "Loan Term Requested:  "+stringx.TransNumberToString[int32](loan.LoanTerm, true))
		form.pdf.Br(lineHeight)
		_ = form.pdf.Cell(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: fontHeight}, "Estimated Current Value (LTV):  "+stringx.TransNumberToString[float64](loan.Ltv, true))
		form.pdf.Br(lineHeight)
		_ = form.pdf.Cell(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: fontHeight}, "Estimated After Repair Value (ARV):  "+stringx.TransNumberToString[float64](loan.Arv, true))
		form.pdf.Br(lineHeight)
		_ = form.pdf.Cell(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: fontHeight}, "Existing Loan Amount (Payoff):  "+stringx.TransNumberToString[float64](loan.PrincipalAmount, true))
		form.pdf.Br(lineHeight)
		_ = form.pdf.Cell(&gopdf.Rect{W: form.pdf.GetFullWidth() / 2.0, H: fontHeight}, "Existing Lender:  "+loan.ExistingLenderCompany)
		_ = form.pdf.Cell(&gopdf.Rect{W: form.pdf.GetFullWidth() / 2.0, H: fontHeight}, "Existing Loan Rate:  "+stringx.TransNumberToString[float32](loan.ExistingRate, true))
		form.pdf.Br(lineHeight)
		_ = form.pdf.Cell(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: fontHeight}, "Existing Loan Maturity Date:  "+loan.ExistingMaturityDate)
		form.pdf.Br(lineHeight)
		_ = form.pdf.Cell(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: fontHeight}, "Subordinate Debt (Y/N-$):  "+condition.Ternary(loan.HasSubordinateDebt, "Y", "N"))
		form.pdf.Br(lineHeight)
		_ = form.pdf.Cell(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: fontHeight}, "Liquidity:  "+stringx.TransNumberToString[float64](loan.LiquidityAmount, true))
		form.pdf.Br(lineHeight)
		_ = form.pdf.Cell(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: fontHeight}, "Exit Strategy:  "+strings.Join(loan.BusinessStrategy, ","))
		form.pdf.Br(lineHeight)
		_ = form.pdf.Cell(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: fontHeight}, "Mid Construction (Y/N):  "+condition.Ternary(loan.MidConstruction, "Y", "N"))
	}

	return err
}
