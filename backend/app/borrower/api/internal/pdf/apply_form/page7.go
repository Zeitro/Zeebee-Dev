package apply_form

import (
	"github.com/signintech/gopdf"

	"zeebee-backend/app/borrower/api/internal/pdf"
)

func (form *ApplyForm) page7() (err error) {
	form.pdf.AddPage()
	col1Width := 445.0
	col2Width := 50.0
	col3Width := 50.0
	col1Opt := gopdf.CellOption{
		Align: gopdf.Left | gopdf.Middle,
	}
	col2Opt := gopdf.CellOption{
		Align:  gopdf.Center | gopdf.Middle,
		Border: gopdf.Left | gopdf.Right,
	}
	col3Opt := gopdf.CellOption{
		Align: gopdf.Center | gopdf.Middle,
	}

	lineHeight := 18.0
	fontHeight := 16.0
	err = form.setPageTitle(pdf.NoteIco, "DECLARATIONS AND", "DEMOGRAGPHICS")
	if err != nil {
		return err
	}
	form.pdf.Br(lineHeight)
	x := form.pdf.GetX()
	y := form.pdf.GetY()
	titleHeight := 27.0
	_ = form.pdf.ImageFrom(pdf.GrayBg, form.pdf.MarginLeft(), y, &gopdf.Rect{W: form.pdf.GetFullWidth(), H: titleHeight})
	_ = form.pdf.ImageFrom(pdf.BlueBg, x+10, y, &gopdf.Rect{W: form.pdf.GetFullWidth() - 10, H: titleHeight})
	form.pdf.SetBoldFont()
	form.pdf.SetWhiteFont()
	form.pdf.SetGrayStroke()
	form.pdf.SetXY(x+10, y)
	_ = form.pdf.CellWithOption(&gopdf.Rect{W: col1Width, H: titleHeight}, "  About this Property and Your Money for this Loan", col1Opt)
	// pdf.SetGrayStroke(0.6)
	_ = form.pdf.CellWithOption(&gopdf.Rect{W: col2Width, H: titleHeight}, "YES", col2Opt)
	_ = form.pdf.CellWithOption(&gopdf.Rect{W: col3Width, H: titleHeight}, "NO", col3Opt)
	form.pdf.SetBlackFont()
	form.pdf.SetRegularFont()
	form.pdf.Br(titleHeight)
	x = form.pdf.GetX()
	y = form.pdf.GetY()
	text := "Will you occupy the property as your primary residence?"
	_ = form.buildPage7TableLine(text, x+10, y, fontHeight, col1Width, col2Width, col3Width, false, "")

	x = form.pdf.GetX()
	y = form.pdf.GetY()
	text = "If this is a Purchase Transaction. Do you have a family relationship or business affiliation with the seller of the property?"
	_ = form.buildPage7TableLine(text, x+10, y, fontHeight, col1Width, col2Width, col3Width, false, "")

	x = form.pdf.GetX()
	y = form.pdf.GetY()
	text = "Are you borrowing any money for this real estate transaction (e.g. money for your closing costs or down payment) or obtaining any money from another party such as the seller or realtor, that you have not disclosed on this loan application? If Yes, what is the amount of this money?"
	_ = form.buildPage7TableLine(text, x+10, y, fontHeight, col1Width, col2Width, col3Width, true, "")

	x = form.pdf.GetX()
	y = form.pdf.GetY()
	text = "Have you or will you be applying for a mortgage loan on another property (not the property securing this loan) on or before closing this transaction that is not disclosed on this loan application?"
	_ = form.buildPage7TableLine(text, x+10, y, fontHeight, col1Width, col2Width, col3Width, false, "")

	x = form.pdf.GetX()
	y = form.pdf.GetY()
	text = " Will this property be subject to a lien that could take priority over the first mortgage lien, such as a clean energy lien paid through your property taxes (e.g., the Property Assessed Clean Energy Program)?"
	_ = form.buildPage7TableLine(text, x+10, y, fontHeight, col1Width, col2Width, col3Width, false, "")
	form.pdf.Br(titleHeight)

	x = form.pdf.GetX()
	y = form.pdf.GetY()
	_ = form.pdf.ImageFrom(pdf.GrayBg, form.pdf.MarginLeft(), y, &gopdf.Rect{W: form.pdf.GetFullWidth(), H: titleHeight})
	_ = form.pdf.ImageFrom(pdf.BlueBg, x+10, y, &gopdf.Rect{W: form.pdf.GetFullWidth() - 10, H: titleHeight})
	form.pdf.SetBoldFont()
	form.pdf.SetWhiteFont()
	form.pdf.SetXY(x+10, y)
	_ = form.pdf.CellWithOption(&gopdf.Rect{W: col1Width, H: titleHeight}, "   About this Property and Your Money for this Loan", col1Opt)
	// pdf.SetGrayStroke(0.6)
	_ = form.pdf.CellWithOption(&gopdf.Rect{W: col2Width, H: titleHeight}, "YES", col2Opt)
	_ = form.pdf.CellWithOption(&gopdf.Rect{W: col3Width, H: titleHeight}, "NO", col3Opt)
	form.pdf.SetBlackFont()
	form.pdf.SetRegularFont()
	form.pdf.Br(titleHeight)

	x = form.pdf.GetX()
	y = form.pdf.GetY()
	text = "Are there any outstanding judgments against you?"
	_ = form.buildPage7TableLine(text, x+10, y, fontHeight, col1Width, col2Width, col3Width, false, "")
	x = form.pdf.GetX()
	y = form.pdf.GetY()
	text = "Are you currently delinquent or in default on a Federal Debt?"
	_ = form.buildPage7TableLine(text, x+10, y, fontHeight, col1Width, col2Width, col3Width, false, "")
	x = form.pdf.GetX()
	y = form.pdf.GetY()
	text = "Are you a party to a lawsuit in which you potentially have any personal financial liability?"
	_ = form.buildPage7TableLine(text, x+10, y, fontHeight, col1Width, col2Width, col3Width, false, "")
	x = form.pdf.GetX()
	y = form.pdf.GetY()
	text = "Have you conveyed title to any property in lieu of foreclosure in the past 7 years?"
	_ = form.buildPage7TableLine(text, x+10, y, fontHeight, col1Width, col2Width, col3Width, false, "")
	x = form.pdf.GetX()
	y = form.pdf.GetY()
	text = "Within the past 7 years, have you completed a pre-foreclosure sale or short sale, whereby the property was sold to a third party and the Lender agreed to accept less than the outstanding mortgage balance due?"
	_ = form.buildPage7TableLine(text, x+10, y, fontHeight, col1Width, col2Width, col3Width, false, "")
	x = form.pdf.GetX()
	y = form.pdf.GetY()
	text = "Have you had property foreclosed upon in the last 7 years?"
	_ = form.buildPage7TableLine(text, x+10, y, fontHeight, col1Width, col2Width, col3Width, false, "")

	// last
	x = form.pdf.GetX()
	y = form.pdf.GetY()
	startX := x + 10
	startY := y + fontHeight/2.0
	colHeight := fontHeight
	form.pdf.SetXY(startX, startY)
	text = "Have you declared bankruptcy within the past 7 years?"
	_ = form.pdf.CellWithOption(&gopdf.Rect{W: col1Width, H: fontHeight}, text, gopdf.CellOption{
		Align: gopdf.Left | gopdf.Middle,
	})
	startY += fontHeight
	form.pdf.SetXY(startX, startY)
	text = "If YES, identify the type(s) of bankruptcy:"
	_ = form.pdf.CellWithOption(&gopdf.Rect{W: col1Width, H: fontHeight}, text, gopdf.CellOption{
		Align: gopdf.Left | gopdf.Middle,
	})
	startY += fontHeight
	colHeight += fontHeight
	_ = form.buildCheckbox(startX+50, startY+5, []string{"Chapter 7", "Chapter 11", "Chapter 12", "Chapter 13"})
	colHeight += 35

	form.pdf.SetXY(startX+col1Width, y)
	_ = form.pdf.CellWithOption(&gopdf.Rect{W: col2Width, H: colHeight}, "", gopdf.CellOption{
		Align:  gopdf.Center | gopdf.Middle,
		Border: gopdf.Left | gopdf.Right,
	})
	_ = form.pdf.CellWithOption(&gopdf.Rect{W: col3Width, H: colHeight}, "", gopdf.CellOption{
		Align: gopdf.Center | gopdf.Middle,
	})
	form.pdf.Line(form.pdf.MarginLeft(), y+colHeight, form.pdf.GetFullWidth()+form.pdf.MarginLeft(), y+colHeight)
	form.pdf.Br(colHeight)

	return err
}

func (form *ApplyForm) buildPage7TableLine(txt string, startX float64, startY float64, fontHeight float64, col1Width, col2Width, col3Width float64, hasExtraTxt bool, extraTxt string) (err error) {
	breakOpt := &gopdf.BreakOption{
		Mode:           gopdf.BreakModeIndicatorSensitive,
		BreakIndicator: ' ', // 在空格处分割
	}
	strs, _ := form.pdf.SplitTextWithOption(txt, col1Width, breakOpt)
	strLen := len(strs)
	colHeight := float64(len(strs))*fontHeight + fontHeight
	tmpY := startY + fontHeight/2.0
	for index, str := range strs {
		form.pdf.SetXY(startX, tmpY)
		if index == strLen-1 && hasExtraTxt {
			textLen, _ := form.pdf.MeasureTextWidth(str)
			extraLen, _ := form.pdf.MeasureTextWidth(extraTxt)
			_ = form.pdf.CellWithOption(&gopdf.Rect{W: textLen, H: fontHeight}, str, gopdf.CellOption{
				Align: gopdf.Left | gopdf.Middle,
			})

			if extraLen == 0 {
				extraLen = 50
			}
			_ = form.pdf.CellWithOption(&gopdf.Rect{W: extraLen + 10, H: fontHeight}, extraTxt, gopdf.CellOption{
				Align:  gopdf.Left | gopdf.Middle,
				Border: gopdf.Bottom,
			})
		} else {
			_ = form.pdf.CellWithOption(&gopdf.Rect{W: col1Width, H: fontHeight}, str, gopdf.CellOption{
				Align: gopdf.Left | gopdf.Middle,
			})
			tmpY += fontHeight
		}

	}
	form.pdf.SetXY(startX+col1Width, startY)
	_ = form.pdf.CellWithOption(&gopdf.Rect{W: col2Width, H: colHeight}, "", gopdf.CellOption{
		Align:  gopdf.Center | gopdf.Middle,
		Border: gopdf.Left | gopdf.Right,
	})
	_ = form.pdf.CellWithOption(&gopdf.Rect{W: col3Width, H: colHeight}, "", gopdf.CellOption{
		Align: gopdf.Center | gopdf.Middle,
	})
	form.pdf.Line(startX, startY+colHeight, form.pdf.GetFullWidth()+form.pdf.MarginLeft(), startY+colHeight)
	form.pdf.Br(colHeight)

	return err
}
