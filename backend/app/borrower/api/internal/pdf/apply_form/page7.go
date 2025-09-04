package apply_form

import (
	"github.com/signintech/gopdf"

	"zeebee-backend/app/borrower/api/internal/pdf"
	"zeebee-backend/common/constants"
	"zeebee-backend/common/stringx"
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

	property := form.info.LoanDemographic.Property
	x = form.pdf.GetX()
	y = form.pdf.GetY()

	startPos := page7TableLinePos{
		x:          x + 10,
		y:          y,
		col1Width:  col1Width,
		col2Width:  col2Width,
		col3Width:  col3Width,
		fontHeight: form.pdf.LineHeight,
	}
	vals := []page7TableRowVal{
		{
			text:        "Will you occupy the property as your primary residence?",
			isYes:       property.AsPrimaryResident,
			hasExtraTxt: false,
			extraTxt:    "",
		},
		{
			text:        "If this is a Purchase Transaction. Do you have a family relationship or business affiliation with the seller of the property?",
			isYes:       property.HasTransactionRelationship,
			hasExtraTxt: false,
			extraTxt:    "",
		},
		{
			text:        "Are you borrowing any money for this real estate transaction (e.g. money for your closing costs or down payment) or obtaining any money from another party such as the seller or realtor, that you have not disclosed on this loan application? If Yes, what is the amount of this money?",
			isYes:       property.HasBorrowedMoney,
			hasExtraTxt: true,
			extraTxt:    stringx.FormatUsDollarCurrency(property.BorrowedMoneyAmount),
		},
		{
			text:        "Have you or will you be applying for a mortgage loan on another property (not the property securing this loan) on or before closing this transaction that is not disclosed on this loan application?",
			isYes:       property.HasAnotherMortgageLoan,
			hasExtraTxt: false,
			extraTxt:    "",
		},
		{
			text:        "Will this property be subject to a lien that could take priority over the first mortgage lien, such as a clean energy lien paid through your property taxes (e.g., the Property Assessed Clean Energy Program)?",
			isYes:       property.HasPriorityOver1st,
			hasExtraTxt: false,
			extraTxt:    "",
		},
	}

	err = form.buildPage7TableRows(startPos, vals)
	if err != nil {
		return err
	}
	form.pdf.Br(titleHeight)
	x = form.pdf.GetX()
	y = form.pdf.GetY()
	_ = form.pdf.ImageFrom(pdf.GrayBg, form.pdf.MarginLeft(), y, &gopdf.Rect{W: form.pdf.GetFullWidth(), H: titleHeight})
	_ = form.pdf.ImageFrom(pdf.BlueBg, x+10, y, &gopdf.Rect{W: form.pdf.GetFullWidth() - 10, H: titleHeight})
	form.pdf.SetBoldFont()
	form.pdf.SetWhiteFont()
	form.pdf.SetXY(x+10, y)
	_ = form.pdf.CellWithOption(&gopdf.Rect{W: col1Width, H: titleHeight}, " About this Property and Your Money for this Loan", col1Opt)
	// pdf.SetGrayStroke(0.6)
	_ = form.pdf.CellWithOption(&gopdf.Rect{W: col2Width, H: titleHeight}, "YES", col2Opt)
	_ = form.pdf.CellWithOption(&gopdf.Rect{W: col3Width, H: titleHeight}, "NO", col3Opt)
	form.pdf.SetBlackFont()
	form.pdf.SetRegularFont()
	form.pdf.Br(titleHeight)
	x = form.pdf.GetX()
	y = form.pdf.GetY()

	startPos = page7TableLinePos{
		x:          x + 10,
		y:          y,
		col1Width:  col1Width,
		col2Width:  col2Width,
		col3Width:  col3Width,
		fontHeight: form.pdf.LineHeight,
	}
	vals = []page7TableRowVal{
		{
			text:        "Are there any outstanding judgments against you?",
			isYes:       property.HasOutstandingJudgments,
			hasExtraTxt: false,
			extraTxt:    "",
		},
		{
			text:        "Are you currently delinquent or in default on a Federal Debt?",
			isYes:       property.HasFederalDebt,
			hasExtraTxt: false,
			extraTxt:    "",
		},
		{
			text:        "Are you a party to a lawsuit in which you potentially have any personal financial liability?",
			isYes:       property.HasALawsuit,
			hasExtraTxt: false,
			extraTxt:    "",
		},
		{
			text:        "Have you conveyed title to any property in lieu of foreclosure in the past 7 years?",
			isYes:       property.HasConveyedTitle,
			hasExtraTxt: false,
			extraTxt:    "",
		},
		{
			text:        "Within the past 7 years, have you completed a pre-foreclosure sale or short sale, whereby the property was sold to a third party and the Lender agreed to accept less than the outstanding mortgage balance due?",
			isYes:       property.HasAPreForeclosureSale,
			hasExtraTxt: false,
			extraTxt:    "",
		},
		{
			text:        "Have you had property foreclosed upon in the last 7 years?",
			isYes:       property.HasForeclosedUpon,
			hasExtraTxt: false,
			extraTxt:    "",
		},
	}

	err = form.buildPage7TableRows(startPos, vals)
	if err != nil {
		return err
	}
	// last
	x = form.pdf.GetX()
	y = form.pdf.GetY()
	startX := x + 10
	startY := y + fontHeight/2.0
	colHeight := fontHeight
	form.pdf.SetXY(startX, startY)
	text := "Have you declared bankruptcy within the past 7 years?"
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
	err = form.buildRadioCheckbox(startX+50, startY+5, constants.BankruptcyChapterName, property.BankruptcyChapter)
	if err != nil {
		return err
	}
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

type page7TableLinePos struct {
	x          float64
	y          float64
	col1Width  float64
	col2Width  float64
	col3Width  float64
	fontHeight float64
}

type page7TableLineValue struct {
	yes         string
	no          string
	hasExtraTxt bool
	extraTxt    string
}

type page7TableRowVal struct {
	text        string
	isYes       bool
	hasExtraTxt bool
	extraTxt    string

	// col1Width  float64
	// col2Width  float64
	// col3Width  float64
	// fontHeight float64
}

func (form *ApplyForm) buildPage7TableRows(startPos page7TableLinePos, vals []page7TableRowVal) (err error) {
	for _, val := range vals {
		text := val.text
		// x := form.pdf.GetX()
		y := form.pdf.GetY()
		pos := page7TableLinePos{
			x:          startPos.x,
			y:          y,
			col1Width:  startPos.col1Width,
			col2Width:  startPos.col2Width,
			col3Width:  startPos.col3Width,
			fontHeight: startPos.fontHeight,
		}
		v := page7TableLineValue{

			hasExtraTxt: val.hasExtraTxt,
			extraTxt:    val.extraTxt,
		}
		if val.isYes {
			v.yes = "Y"
		} else {
			v.no = "N"
		}
		err = form.buildPage7TableRow(text, pos, v)
		if err != nil {
			return err
		}

	}
	return err
}

func (form *ApplyForm) buildPage7TableRow(txt string, pos page7TableLinePos, val page7TableLineValue) (err error) {
	breakOpt := &gopdf.BreakOption{
		Mode:           gopdf.BreakModeIndicatorSensitive,
		BreakIndicator: ' ', // 在空格处分割
	}
	strs, _ := form.pdf.SplitTextWithOption(txt, pos.col1Width, breakOpt)
	strLen := len(strs)
	colHeight := float64(len(strs))*pos.fontHeight + pos.fontHeight
	tmpY := pos.y + pos.fontHeight/2.0
	for index, str := range strs {
		form.pdf.SetXY(pos.x, tmpY)
		if index == strLen-1 && val.hasExtraTxt {
			textLen, _ := form.pdf.MeasureTextWidth(str)
			extraLen, _ := form.pdf.MeasureTextWidth(val.extraTxt)
			_ = form.pdf.CellWithOption(&gopdf.Rect{W: textLen, H: pos.fontHeight}, str, gopdf.CellOption{
				Align: gopdf.Left | gopdf.Middle,
			})

			if extraLen == 0 {
				extraLen = 50
			}
			_ = form.pdf.CellWithOption(&gopdf.Rect{W: extraLen + 10, H: pos.fontHeight}, val.extraTxt, gopdf.CellOption{
				Align:  gopdf.Left | gopdf.Middle,
				Border: gopdf.Bottom,
			})
		} else {
			_ = form.pdf.CellWithOption(&gopdf.Rect{W: pos.col1Width, H: pos.fontHeight}, str, gopdf.CellOption{
				Align: gopdf.Left | gopdf.Middle,
			})
			tmpY += pos.fontHeight
		}

	}
	form.pdf.SetXY(pos.x+pos.col1Width, pos.y)
	_ = form.pdf.CellWithOption(&gopdf.Rect{W: pos.col2Width, H: colHeight}, val.yes, gopdf.CellOption{
		Align:  gopdf.Center | gopdf.Middle,
		Border: gopdf.Left | gopdf.Right,
	})
	_ = form.pdf.CellWithOption(&gopdf.Rect{W: pos.col3Width, H: colHeight}, val.no, gopdf.CellOption{
		Align: gopdf.Center | gopdf.Middle,
	})
	form.pdf.Line(pos.x, pos.y+colHeight, form.pdf.GetFullWidth()+form.pdf.MarginLeft(), pos.y+colHeight)
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
