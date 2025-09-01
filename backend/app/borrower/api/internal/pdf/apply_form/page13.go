package apply_form

import (
	"github.com/signintech/gopdf"

	"zeebee-backend/app/borrower/api/internal/pdf"
)

func (form *ApplyForm) page13() (err error) {
	form.pdf.AddPage()
	err = form.setPageTitle(pdf.ListIco, "BORROWER SIGNATURE", "AUTHORIZATION")
	if err != nil {
		return err
	}
	form.pdf.Br(form.pdf.LineHeight)
	x := form.pdf.GetX()
	y := form.pdf.GetY()
	_ = form.buildSmallTitle("Privacy Act Notice", position{
		x: x,
		y: y,
		w: form.pdf.GetFullWidth(),
		h: form.pdf.LineHeight,
	})
	form.pdf.Br(form.pdf.LineHeight * 1.5)

	x = form.pdf.GetX()
	y = form.pdf.GetY()
	contents := []colorText{
		{
			text:  "Privacy Act Notice:",
			color: struct{ r, g, b uint8 }{71, 117, 255},
		},
		{
			text:  " This information is to be used by the agency collecting it or its assignees in determining whether you qualify as a prospective mortgagor under its program. It will not be disclosed outside the agency except as required and permitted by law. You do not have to provide this information, but if you do not your application for approval as a prospective mortgagor or borrower maybe delayed or rejected. The information requested in this form is authorized by Title 38, USC, Chapter 37 (if VA); by 12 USC, Section 1701et. seq. (if HUD/FHA); by 42 USC, Section 1452b (if HUD/CPD); and Title 42 USC, 1471 et. seq., or 7 USC, 1921 et. seq. (if USDA/FmHA). ",
			color: struct{ r, g, b uint8 }{0, 0, 0},
		},
	}
	_ = form.buildLongTextWithColor(contents, x, form.pdf.LineHeight)
	form.pdf.Br(form.pdf.LineHeight)
	x = form.pdf.GetX()
	y = form.pdf.GetY()
	_ = form.buildTitleWithBg("Part 1 - General Information", position{w: form.pdf.GetFullWidth(), h: 27, x: x, y: y}, position{w: form.pdf.GetFullWidth() - 10, h: 27, x: x + 10, y: y})
	form.pdf.Br(form.pdf.LineHeight * 2)
	// borrower
	x = form.pdf.GetX()
	y = form.pdf.GetY()
	form.pdf.SetBlackStroke()
	form.pdf.Line(x, y, x+form.pdf.GetFullWidth(), y)
	rowHeight := form.pdf.LineHeight * 5.0
	form.pdf.SetBoldFont()
	_ = form.pdf.CellWithOption(&gopdf.Rect{W: 100.0, H: rowHeight}, "1. Borrower", gopdf.CellOption{
		Align: gopdf.Center | gopdf.Middle,
	})
	_ = form.pdf.ImageFrom(pdf.LightBlueBg, x+100, y+form.pdf.LineHeight*1.5, &gopdf.Rect{W: 120.0, H: form.pdf.LineHeight * 2})
	form.pdf.SetRegularFont()
	_ = form.pdf.CellWithOption(&gopdf.Rect{W: 120.0, H: rowHeight}, "", gopdf.CellOption{
		Align: gopdf.Center | gopdf.Middle,
	})
	_ = form.pdf.CellWithOption(&gopdf.Rect{W: 10.0, H: rowHeight}, "", gopdf.CellOption{
		Align:  gopdf.Center | gopdf.Middle,
		Border: gopdf.Right,
	})

	// lender
	form.pdf.SetXY(x+230, y+form.pdf.LineHeight/2.0)
	form.pdf.SetBoldFont()
	x = form.pdf.GetX()
	y = form.pdf.GetY()
	_ = form.pdf.CellWithOption(&gopdf.Rect{W: 230.0, H: form.pdf.LineHeight}, "2. Name and address of Lender/ Broker", gopdf.CellOption{
		Align: gopdf.Center | gopdf.Middle,
	})

	_ = form.pdf.ImageFrom(pdf.LightBlueBg, x+230, y, &gopdf.Rect{W: 95.0, H: form.pdf.LineHeight * 4.0})

	form.pdf.SetXY(x+20, y+form.pdf.LineHeight)
	form.pdf.SetRegularFont()
	_ = form.pdf.CellWithOption(&gopdf.Rect{W: 210.0, H: form.pdf.LineHeight}, "Yeelee Capital", gopdf.CellOption{
		Align: gopdf.Left | gopdf.Middle,
	})
	form.pdf.SetXY(x+20, y+form.pdf.LineHeight*2.0)
	_ = form.pdf.CellWithOption(&gopdf.Rect{W: 210.0, H: form.pdf.LineHeight}, "340 E Middlefield Rd, STE 100.", gopdf.CellOption{
		Align: gopdf.Left | gopdf.Middle,
	})
	form.pdf.SetXY(x+20, y+form.pdf.LineHeight*3.0)
	_ = form.pdf.CellWithOption(&gopdf.Rect{W: 210.0, H: form.pdf.LineHeight}, "Mountain View, CA 94043", gopdf.CellOption{
		Align: gopdf.Left | gopdf.Middle,
	})
	form.pdf.Br(form.pdf.LineHeight * 1.5)
	x = form.pdf.GetX()
	y = form.pdf.GetY()
	form.pdf.Line(x, y, x+form.pdf.GetFullWidth(), y)

	// date
	form.pdf.SetBoldFont()
	_ = form.pdf.CellWithOption(&gopdf.Rect{W: 100.0, H: rowHeight}, "date", gopdf.CellOption{
		Align: gopdf.Center | gopdf.Middle,
	})
	_ = form.pdf.ImageFrom(pdf.LightBlueBg, x+100, y+form.pdf.LineHeight*1.5, &gopdf.Rect{W: 120.0, H: form.pdf.LineHeight * 2})
	form.pdf.SetRegularFont()
	_ = form.pdf.CellWithOption(&gopdf.Rect{W: 120.0, H: rowHeight}, "", gopdf.CellOption{
		Align: gopdf.Center | gopdf.Middle,
	})
	_ = form.pdf.CellWithOption(&gopdf.Rect{W: 10.0, H: rowHeight}, "", gopdf.CellOption{
		Align:  gopdf.Center | gopdf.Middle,
		Border: gopdf.Right,
	})
	// Loan Number
	form.pdf.SetXY(x+230, y)
	form.pdf.SetBoldFont()
	x = form.pdf.GetX()
	y = form.pdf.GetY()
	form.pdf.SetBoldFont()
	_ = form.pdf.CellWithOption(&gopdf.Rect{W: 100.0, H: rowHeight}, "4. Loan Number", gopdf.CellOption{
		Align: gopdf.Center | gopdf.Middle,
	})
	_ = form.pdf.ImageFrom(pdf.LightBlueBg, x+100, y+form.pdf.LineHeight*1.5, &gopdf.Rect{W: 225.0, H: form.pdf.LineHeight * 2})
	form.pdf.SetRegularFont()
	_ = form.pdf.CellWithOption(&gopdf.Rect{W: 225.0, H: rowHeight}, "", gopdf.CellOption{
		Align: gopdf.Center | gopdf.Middle,
	})
	form.pdf.Br(form.pdf.LineHeight * 5)
	x = form.pdf.GetX()
	y = form.pdf.GetY()
	form.pdf.Line(x, y, x+form.pdf.GetFullWidth(), y)

	form.pdf.Br(form.pdf.LineHeight)
	_ = form.buildSmallTitle("Part II – Borrower Authorization", position{
		x: form.pdf.GetX(),
		y: form.pdf.GetY(),
		w: form.pdf.GetFullWidth(),
		h: form.pdf.LineHeight,
	})
	form.pdf.Br(form.pdf.LineHeight * 2)

	text := []string{
		"I hereby authorize the Lender/ Broker to verify my past and present employment earnings records, bank accounts, stock holdings, insurance information, and any other asset balances that are needed to process my mortgage loan application. I further authorize the Lender/ Broker to order a consumer credit report and verify other credit information, including past and present mortgage and landlord references. It is understood that a copy of this form will also serve as authorization.",
		"The information the Lender/ Broker obtains is only to be used in the processing of my application for a mortgage loan.",
	}
	err = form.pdf.MultiLongText(form.pdf.LineHeight, text)
	if err != nil {
		return err
	}
	_ = form.buildSignature("")

	return err
}
