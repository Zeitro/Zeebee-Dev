package apply_form

import (
	"github.com/duke-git/lancet/v2/slice"
	"github.com/signintech/gopdf"

	"zeebee-backend/app/borrower/api/internal/pdf"
)

func (form *ApplyForm) page8() (err error) {
	form.pdf.AddPage()
	lineHeight := 18.0
	err = form.setPageTitle(pdf.NoteIco, "DECLARATIONS AND", "DEMOGRAGPHICS")
	if err != nil {
		return err
	}
	form.pdf.Br(lineHeight)
	x := form.pdf.GetX()
	y := form.pdf.GetY()
	_ = form.buildTitleWithBg("Borrower Demographic Information:", position{w: form.pdf.GetFullWidth(), h: 27, x: x, y: y}, position{w: form.pdf.GetFullWidth() - 10, h: 27, x: x + 10, y: y})
	form.pdf.Br(30)

	x = form.pdf.GetX()
	y = form.pdf.GetY()
	contents := []colorText{
		{text: "The purpose of collecting this information is to help ensure that all applicants are treated fairly and that the housing needs of communities and neighborhoods are being fulfilled. For residential mortgage lending, Federal law requires that we ask applicants for their demographic information (ethnicity, sex, and race) in order to monitor our compliance with equal credit opportunity, fair housing, and home mortgage disclosure laws. You are not required to provide this information, but are encouraged to do so. You may select one or more designations for “ethnicity” and one or more designations for “race.”",
			color: struct{ r, g, b uint8 }{0, 0, 0},
		},
		{text: " The law provides that we may not discriminate",
			color: struct{ r, g, b uint8 }{71, 117, 255},
		},
		{text: " on the basis of this information, or on where you choose to provide it. However, if you choose not to provide the information and you have made this application in person, Federal regulations require us to note your ethnicity, sex and race on the basis of visual observation or surname. The law also provides that we may not discriminate on the basis of age or marital status information provide in this application. If you do not wish to provide some or all of this information, please indicate so below: ",
			color: struct{ r, g, b uint8 }{0, 0, 0},
		},
	}
	err = form.buildLongTextWithColor(contents, x, lineHeight)
	if err != nil {
		return err
	}
	form.pdf.Br(lineHeight / 2.0)
	form.pdf.SetGrayStroke()
	x = form.pdf.GetX()
	y = form.pdf.GetY()
	form.pdf.Line(form.pdf.MarginLeft(), y, form.pdf.GetFullWidth()+form.pdf.MarginLeft(), y)
	form.pdf.Br(lineHeight / 2.0)
	form.pdf.SetBoldFont()
	cellTextOpt := gopdf.CellOption{
		Align: gopdf.Left | gopdf.Middle,
	}

	// Ethnicity
	ethnicity := form.info.LoanDemographic.Ethnicity
	form.pdf.SetGrayStroke()
	_ = form.pdf.CellWithOption(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: lineHeight}, "Ethnicity (Check One or More)", cellTextOpt)
	form.pdf.SetRegularFont()

	// checkboxImg := pdf.EmptyCheckBox

	form.pdf.Br(lineHeight)
	x = form.pdf.GetX()
	y = form.pdf.GetY()
	pos := position{
		x: x,
		y: y,
		h: lineHeight,
	}
	opts := []page9CheckBoxOpt{
		{
			label:      "Hispanic or Latino",
			values:     ethnicity.Values,
			targetText: "Hispanic or Latino",
			extraText:  "",
			inputText:  "",
		},
		{
			label:      "Mexican",
			values:     ethnicity.Values,
			targetText: "Mexican",
			extraText:  "",
			inputText:  "",
		},
		{
			label:      "Cuban",
			values:     ethnicity.Values,
			targetText: "Cuban",
			extraText:  "",
			inputText:  "",
		},
		{
			label:      "Other Hispanic or Latino – Print Origin",
			values:     ethnicity.Values,
			targetText: "Other Hispanic or Latino",
			extraText:  "(For example: Colombian, Dominican, Nicaraguan, Salvadoran, Spaniard, and so on) ",
			inputText:  ethnicity.OtherHispanicOrLatinoName,
		},
		{
			label:      "Not Hispanic or Latino",
			values:     ethnicity.Values,
			targetText: "Not Hispanic or Latino",
			extraText:  "",
			inputText:  "",
		},
		{
			label:      "I do not wish to provide this information",
			values:     ethnicity.Values,
			targetText: "Hidden",
			extraText:  "",
			inputText:  "",
		},
	}

	err = form.buildPage9CheckBoxes(pos, lineHeight, opts)
	if err != nil {
		return err
	}
	form.pdf.Br(lineHeight)

	// Race
	race := form.info.LoanDemographic.Race
	form.pdf.SetBoldFont()
	_ = form.pdf.CellWithOption(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: lineHeight}, "Race (Check One or More)", cellTextOpt)
	form.pdf.SetRegularFont()
	form.pdf.Br(lineHeight)
	x = form.pdf.GetX()
	y = form.pdf.GetY()
	pos = position{
		x: x,
		y: y,
		h: lineHeight,
	}
	opts = []page9CheckBoxOpt{
		{
			label:      "American Indian or Alaska Native – Print name of enrolled or principal tribe",
			values:     race.Values,
			targetText: "American Indian or Alaska Native",
			extraText:  "",
			inputText:  race.AmericanIndianOrAlaskaNativeName,
		},
		{
			label:      "Asian",
			values:     race.Values,
			targetText: "Asian",
			extraText:  "",
			inputText:  "",
		},
	}
	err = form.buildPage9CheckBoxes(pos, lineHeight, opts)
	if err != nil {
		return err
	}

	x = form.pdf.GetX()
	y = form.pdf.GetY()
	pos = position{
		x: x + lineHeight*2.5,
		y: y,
		h: lineHeight,
	}
	opts = []page9CheckBoxOpt{
		{
			label:      "Asian Indian",
			values:     race.AsianValues,
			targetText: "Asian Indian",
			extraText:  "",
			inputText:  "",
		},
		{
			label:      "Chinese",
			values:     race.AsianValues,
			targetText: "Chinese",
			extraText:  "",
			inputText:  "",
		},
		{
			label:      "Filipino",
			values:     race.AsianValues,
			targetText: "Filipino",
			extraText:  "",
			inputText:  "",
		},
		{
			label:      "Japanese",
			values:     race.AsianValues,
			targetText: "Japanese",
			extraText:  "",
			inputText:  "",
		},
		{
			label:      "Korean",
			values:     race.AsianValues,
			targetText: "Korean",
			extraText:  "",
			inputText:  "",
		},
		{
			label:      "Vietnamese",
			values:     race.AsianValues,
			targetText: "Vietnamese",
			extraText:  "",
			inputText:  "",
		},
		{
			label:      "Other Asian – Print Race",
			values:     race.AsianValues,
			targetText: "Other Asian",
			extraText:  "(For example: Hmong, Lootian, Thai, Pakistani, Cambodian, and so on.)",
			inputText:  race.AsianOtherName,
		},
	}

	err = form.buildPage9CheckBoxes(pos, lineHeight*0.8, opts)
	if err != nil {
		return err
	}

	x = form.pdf.GetX()
	y = form.pdf.GetY()
	pos = position{
		x: x,
		y: y,
		h: lineHeight,
	}
	opts = []page9CheckBoxOpt{
		{
			label:      "Black or African American",
			values:     race.Values,
			targetText: "Black or African American",
			extraText:  "",
			inputText:  "",
		},
		{
			label:      "Native Hawaiian or Other Pacific Islander",
			values:     race.Values,
			targetText: "Native Hawaiian or Other Pacific Islander",
			extraText:  "",
			inputText:  "",
		},
		{
			label:      "Native Hawaiian Guamanian or Chamorro Samoan",
			values:     race.Values,
			targetText: "Native Hawaiian Guamanian or Chamorro Samoan",
			extraText:  "",
			inputText:  "",
		},
		{
			label:      "Other Pacific Islander – Print Race",
			values:     race.Values,
			targetText: "Other Pacific Islander",
			extraText:  "(For example: Fijian, Tongan, and so on)",
			inputText:  race.OtherPacificIslanderName,
		},
		{
			label:      "White",
			values:     race.Values,
			targetText: "White",
			extraText:  "",
			inputText:  "",
		},
		{
			label:      "I do not wish to provide this information",
			values:     race.Values,
			targetText: "Hidden",
			extraText:  "",
			inputText:  "",
		},
	}
	err = form.buildPage9CheckBoxes(pos, lineHeight*0.8, opts)
	if err != nil {
		return err
	}

	return err
}

type page9CheckBoxOpt struct {
	label      string
	values     []string
	targetText string
	extraText  string
	inputText  string
}

func (form *ApplyForm) buildPage9CheckBoxes(pos position, checkBoxWidth float64, opts []page9CheckBoxOpt) (err error) {
	currentX := pos.x
	currentY := pos.y
	checkBoxOpt := &gopdf.Rect{W: checkBoxWidth, H: checkBoxWidth}
	cellTextOpt := gopdf.CellOption{
		Align: gopdf.Left | gopdf.Middle,
	}
	cellTextWithBottomBorderOpt := gopdf.CellOption{
		Align:  gopdf.Left | gopdf.Middle,
		Border: gopdf.Bottom,
	}

	for _, opt := range opts {
		form.pdf.SetXY(currentX, currentY)
		checkboxImg := pdf.EmptyCheckBox

		if slice.Contain(opt.values, opt.targetText) {
			checkboxImg = pdf.SelectedCheckBox
		}
		// checkbox image
		_ = form.pdf.ImageFrom(checkboxImg, currentX, currentY, checkBoxOpt)
		form.pdf.SetX(currentX + checkBoxWidth*1.5)
		// _ = form.pdf.CellWithOption(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: pos.h}, opt.label, cellTextOpt)
		// form.pdf.Br(pos.h)

		// label
		textLen, _ := form.pdf.MeasureTextWidth(opt.label)
		_ = form.pdf.CellWithOption(&gopdf.Rect{W: textLen, H: pos.h}, opt.label, cellTextOpt)
		if len(opt.inputText) > 0 {
			_ = form.pdf.CellWithOption(&gopdf.Rect{W: 100, H: pos.h}, " "+opt.inputText, cellTextWithBottomBorderOpt)
		}
		currentY += checkBoxWidth * 1.1
		if len(opt.extraText) > 0 {
			form.pdf.Br(pos.h * 0.9)
			form.pdf.SetX(pos.x)
			_ = form.pdf.Cell(&gopdf.Rect{W: checkBoxWidth * 1.5, H: pos.h}, "")
			_ = form.pdf.CellWithOption(&gopdf.Rect{W: form.pdf.GetFullWidth() - checkBoxWidth*1.5, H: pos.h}, opt.extraText, cellTextOpt)
			currentY += pos.h * 0.8
		}
		form.pdf.Br(pos.h)
	}

	return err
}
