package apply_form

import (
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

	checkBoxOpt := &gopdf.Rect{W: lineHeight * 0.9, H: lineHeight * 0.9}
	samllCheckBoxOpt := &gopdf.Rect{W: lineHeight * 0.8, H: lineHeight * 0.8}
	cellTextOpt := gopdf.CellOption{
		Align: gopdf.Left | gopdf.Middle,
	}
	cellTextWithBottomBorderOpt := gopdf.CellOption{
		Align:  gopdf.Left | gopdf.Middle,
		Border: gopdf.Bottom,
	}
	// Ethnicity
	form.pdf.SetGrayStroke()
	_ = form.pdf.CellWithOption(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: lineHeight}, "Ethnicity (Check One or More)", cellTextOpt)
	form.pdf.SetRegularFont()

	form.pdf.Br(lineHeight)
	x = form.pdf.GetX()
	y = form.pdf.GetY()
	_ = form.pdf.ImageFrom(pdf.EmptyCheckBox, x, y, checkBoxOpt)
	form.pdf.SetX(x + lineHeight*1.5)
	_ = form.pdf.CellWithOption(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: lineHeight}, "Hispanic or Latino", cellTextOpt)
	form.pdf.Br(lineHeight)
	x = form.pdf.GetX()
	y = form.pdf.GetY()
	_ = form.pdf.ImageFrom(pdf.EmptyCheckBox, x, y, checkBoxOpt)
	form.pdf.SetX(x + lineHeight*1.5)
	_ = form.pdf.CellWithOption(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: lineHeight}, "Mexican", cellTextOpt)
	form.pdf.Br(lineHeight)
	x = form.pdf.GetX()
	y = form.pdf.GetY()
	_ = form.pdf.ImageFrom(pdf.EmptyCheckBox, x, y, checkBoxOpt)
	form.pdf.SetX(x + lineHeight*1.5)
	_ = form.pdf.CellWithOption(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: lineHeight}, "Cuban", cellTextOpt)
	form.pdf.Br(lineHeight)
	x = form.pdf.GetX()
	y = form.pdf.GetY()
	_ = form.pdf.ImageFrom(pdf.EmptyCheckBox, x, y, checkBoxOpt)
	form.pdf.SetX(x + lineHeight*1.5)
	textLen, _ := form.pdf.MeasureTextWidth("Other Hispanic or Latino – Print Origin")
	_ = form.pdf.CellWithOption(&gopdf.Rect{W: textLen, H: lineHeight}, "Other Hispanic or Latino – Print Origin", cellTextOpt)
	_ = form.pdf.CellWithOption(&gopdf.Rect{W: 100, H: lineHeight}, "", cellTextWithBottomBorderOpt)
	form.pdf.Br(lineHeight)
	_ = form.pdf.Cell(&gopdf.Rect{W: lineHeight * 1.5, H: lineHeight}, "")
	_ = form.pdf.CellWithOption(&gopdf.Rect{W: form.pdf.GetFullWidth() - lineHeight*1.5, H: lineHeight}, "(For example: Colombian, Dominican, Nicaraguan, Salvadoran, Spaniard, and so on) ", cellTextOpt)
	form.pdf.Br(lineHeight)
	x = form.pdf.GetX()
	y = form.pdf.GetY()
	_ = form.pdf.ImageFrom(pdf.EmptyCheckBox, x, y, checkBoxOpt)
	form.pdf.SetX(x + lineHeight*1.5)
	_ = form.pdf.CellWithOption(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: lineHeight}, "Not Hispanic or Latino", cellTextOpt)
	form.pdf.Br(lineHeight)
	x = form.pdf.GetX()
	y = form.pdf.GetY()
	_ = form.pdf.ImageFrom(pdf.EmptyCheckBox, x, y, checkBoxOpt)
	form.pdf.SetX(x + lineHeight*1.5)
	_ = form.pdf.CellWithOption(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: lineHeight}, "I do not wish to provide this information", cellTextOpt)
	form.pdf.Br(lineHeight * 2.0)
	// Race
	form.pdf.SetBoldFont()
	_ = form.pdf.CellWithOption(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: lineHeight}, "Race (Check One or More)", cellTextOpt)
	form.pdf.SetRegularFont()
	form.pdf.Br(lineHeight)
	x = form.pdf.GetX()
	y = form.pdf.GetY()
	_ = form.pdf.ImageFrom(pdf.EmptyCheckBox, x, y, checkBoxOpt)
	form.pdf.SetX(x + lineHeight*1.5)
	textLen, _ = form.pdf.MeasureTextWidth("American Indian or Alaska Native – Print name of enrolled or principal tribe")
	_ = form.pdf.CellWithOption(&gopdf.Rect{W: textLen, H: lineHeight}, "American Indian or Alaska Native – Print name of enrolled or principal tribe", cellTextOpt)
	_ = form.pdf.CellWithOption(&gopdf.Rect{W: 100, H: lineHeight}, "", cellTextWithBottomBorderOpt)
	form.pdf.Br(lineHeight)
	x = form.pdf.GetX()
	y = form.pdf.GetY()
	_ = form.pdf.ImageFrom(pdf.EmptyCheckBox, x, y, checkBoxOpt)
	form.pdf.SetX(x + lineHeight*1.5)
	_ = form.pdf.CellWithOption(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: lineHeight}, "Asian", cellTextOpt)
	form.pdf.Br(lineHeight)
	x = form.pdf.GetX()
	y = form.pdf.GetY()
	_ = form.pdf.ImageFrom(pdf.EmptyCheckBox, x+lineHeight, y, samllCheckBoxOpt)
	form.pdf.SetX(x + lineHeight*2.5)
	_ = form.pdf.CellWithOption(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: lineHeight}, "Asian Indian", cellTextOpt)
	form.pdf.Br(lineHeight)
	x = form.pdf.GetX()
	y = form.pdf.GetY()
	_ = form.pdf.ImageFrom(pdf.EmptyCheckBox, x+lineHeight, y, samllCheckBoxOpt)
	form.pdf.SetX(x + lineHeight*2.5)
	_ = form.pdf.CellWithOption(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: lineHeight}, "Chinese", cellTextOpt)
	form.pdf.Br(lineHeight)
	x = form.pdf.GetX()
	y = form.pdf.GetY()
	_ = form.pdf.ImageFrom(pdf.EmptyCheckBox, x+lineHeight, y, samllCheckBoxOpt)
	form.pdf.SetX(x + lineHeight*2.5)
	_ = form.pdf.CellWithOption(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: lineHeight}, "Filipino", cellTextOpt)
	form.pdf.Br(lineHeight)
	x = form.pdf.GetX()
	y = form.pdf.GetY()
	_ = form.pdf.ImageFrom(pdf.EmptyCheckBox, x+lineHeight, y, samllCheckBoxOpt)
	form.pdf.SetX(x + lineHeight*2.5)
	_ = form.pdf.CellWithOption(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: lineHeight}, "Japanese", cellTextOpt)
	form.pdf.Br(lineHeight)
	x = form.pdf.GetX()
	y = form.pdf.GetY()
	_ = form.pdf.ImageFrom(pdf.EmptyCheckBox, x+lineHeight, y, samllCheckBoxOpt)
	form.pdf.SetX(x + lineHeight*2.5)
	_ = form.pdf.CellWithOption(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: lineHeight}, "Korean", cellTextOpt)
	form.pdf.Br(lineHeight)
	x = form.pdf.GetX()
	y = form.pdf.GetY()
	_ = form.pdf.ImageFrom(pdf.EmptyCheckBox, x+lineHeight, y, samllCheckBoxOpt)
	form.pdf.SetX(x + lineHeight*2.5)
	_ = form.pdf.CellWithOption(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: lineHeight}, "Vietnamese", cellTextOpt)
	form.pdf.Br(lineHeight)
	x = form.pdf.GetX()
	y = form.pdf.GetY()
	_ = form.pdf.ImageFrom(pdf.EmptyCheckBox, x+lineHeight, y, samllCheckBoxOpt)
	form.pdf.SetX(x + lineHeight*2.5)
	textLen, _ = form.pdf.MeasureTextWidth("Other Asian – Print Race")
	_ = form.pdf.CellWithOption(&gopdf.Rect{W: textLen, H: lineHeight}, "Other Asian – Print Race", cellTextOpt)
	_ = form.pdf.CellWithOption(&gopdf.Rect{W: 100, H: lineHeight}, "", cellTextWithBottomBorderOpt)
	form.pdf.Br(lineHeight)
	_ = form.pdf.Cell(&gopdf.Rect{W: lineHeight * 2.5, H: lineHeight}, "")
	_ = form.pdf.CellWithOption(&gopdf.Rect{W: form.pdf.GetFullWidth() - lineHeight*1.5, H: lineHeight}, " (For example: Hmong, Lootian, Thai, Pakistani, Cambodian, and so on.)", cellTextOpt)
	form.pdf.Br(lineHeight)
	x = form.pdf.GetX()
	y = form.pdf.GetY()
	_ = form.pdf.ImageFrom(pdf.EmptyCheckBox, x, y, checkBoxOpt)
	form.pdf.SetX(x + lineHeight*1.5)
	_ = form.pdf.CellWithOption(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: lineHeight}, "Black or African American", cellTextOpt)
	form.pdf.Br(lineHeight)
	x = form.pdf.GetX()
	y = form.pdf.GetY()
	_ = form.pdf.ImageFrom(pdf.EmptyCheckBox, x, y, checkBoxOpt)
	form.pdf.SetX(x + lineHeight*1.5)
	_ = form.pdf.CellWithOption(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: lineHeight}, "Native Hawaiian or Other Pacific Islander ", cellTextOpt)
	form.pdf.Br(lineHeight)
	x = form.pdf.GetX()
	y = form.pdf.GetY()
	_ = form.pdf.ImageFrom(pdf.EmptyCheckBox, x, y, checkBoxOpt)
	form.pdf.SetX(x + lineHeight*1.5)
	_ = form.pdf.CellWithOption(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: lineHeight}, "Native Hawaiian Guamanian or Chamorro Samoan", cellTextOpt)
	form.pdf.Br(lineHeight)
	x = form.pdf.GetX()
	y = form.pdf.GetY()
	_ = form.pdf.ImageFrom(pdf.EmptyCheckBox, x, y, checkBoxOpt)
	form.pdf.SetX(x + lineHeight*1.5)
	textLen, _ = form.pdf.MeasureTextWidth("Other Asian – Print Race")
	_ = form.pdf.CellWithOption(&gopdf.Rect{W: textLen, H: lineHeight}, "Other Asian – Print Race", cellTextOpt)
	_ = form.pdf.CellWithOption(&gopdf.Rect{W: 100, H: lineHeight}, "", cellTextWithBottomBorderOpt)
	form.pdf.Br(lineHeight)
	_ = form.pdf.Cell(&gopdf.Rect{W: lineHeight * 1.5, H: lineHeight}, "")
	_ = form.pdf.CellWithOption(&gopdf.Rect{W: form.pdf.GetFullWidth() - lineHeight*1.5, H: lineHeight}, " (For example: Hmong, Lootian, Thai, Pakistani, Cambodian, and so on.)", cellTextOpt)
	form.pdf.Br(lineHeight)
	x = form.pdf.GetX()
	y = form.pdf.GetY()
	_ = form.pdf.ImageFrom(pdf.EmptyCheckBox, x, y, checkBoxOpt)
	form.pdf.SetX(x + lineHeight*1.5)
	_ = form.pdf.CellWithOption(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: lineHeight}, "White", cellTextOpt)
	form.pdf.Br(lineHeight)
	x = form.pdf.GetX()
	y = form.pdf.GetY()
	_ = form.pdf.ImageFrom(pdf.EmptyCheckBox, x, y, checkBoxOpt)
	form.pdf.SetX(x + lineHeight*1.5)
	_ = form.pdf.CellWithOption(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: lineHeight}, "I do not wish to provide this information", cellTextOpt)

	return err
}
