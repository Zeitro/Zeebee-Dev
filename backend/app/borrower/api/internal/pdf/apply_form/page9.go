package apply_form

import (
	"github.com/signintech/gopdf"

	"zeebee-backend/app/borrower/api/internal/pdf"
)

func (form *ApplyForm) page9() (err error) {
	form.pdf.AddPage()
	err = form.setPageTitle(pdf.NoteIco, "DECLARATIONS AND", "DEMOGRAGPHICS")
	if err != nil {
		return err
	}
	form.pdf.Br(form.pdf.LineHeight)
	// Race
	checkBoxOpt := &gopdf.Rect{W: form.pdf.LineHeight, H: form.pdf.LineHeight}
	cellTextOpt := gopdf.CellOption{
		Align: gopdf.Left | gopdf.Middle,
	}
	form.pdf.SetBoldFont()
	_ = form.pdf.CellWithOption(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: form.pdf.LineHeight}, "Sex", cellTextOpt)
	form.pdf.SetRegularFont()
	form.pdf.Br(form.pdf.LineHeight)
	x := form.pdf.GetX()
	y := form.pdf.GetY()
	_ = form.pdf.ImageFrom(pdf.EmptyCheckBox, x, y, checkBoxOpt)
	form.pdf.SetX(x + form.pdf.LineHeight*1.5)
	_ = form.pdf.CellWithOption(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: form.pdf.LineHeight}, "Female", cellTextOpt)
	form.pdf.Br(form.pdf.LineHeight * 1.3)
	x = form.pdf.GetX()
	y = form.pdf.GetY()
	_ = form.pdf.ImageFrom(pdf.EmptyCheckBox, x, y, checkBoxOpt)
	form.pdf.SetX(x + form.pdf.LineHeight*1.5)
	_ = form.pdf.CellWithOption(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: form.pdf.LineHeight}, "Male", cellTextOpt)
	form.pdf.Br(form.pdf.LineHeight * 1.3)
	x = form.pdf.GetX()
	y = form.pdf.GetY()
	_ = form.pdf.ImageFrom(pdf.EmptyCheckBox, x, y, checkBoxOpt)
	form.pdf.SetX(x + form.pdf.LineHeight*1.5)
	_ = form.pdf.CellWithOption(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: form.pdf.LineHeight}, " I do not wish to provide this information", cellTextOpt)
	form.pdf.Br(form.pdf.LineHeight * 2)

	col1Width := 405.0
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
	x = form.pdf.GetX()
	y = form.pdf.GetY()
	_ = form.buildTitleWithBg("FOR OFFICE USE ONLY, DO NOT FILL OUT:", position{w: form.pdf.GetFullWidth(), h: 27, x: x, y: y}, position{w: form.pdf.GetFullWidth() - 10, h: 27, x: x + 10, y: y})
	form.pdf.Br(form.pdf.LineHeight * 3)
	x = form.pdf.GetX()
	y = form.pdf.GetY()
	titleHeight := 27.0
	_ = form.pdf.ImageFrom(pdf.GrayBg, x+form.pdf.MarginLeft()*2, y, &gopdf.Rect{W: form.pdf.GetFullWidth() - form.pdf.MarginLeft()*2, H: titleHeight})
	_ = form.pdf.ImageFrom(pdf.BlueBg, x+form.pdf.MarginLeft()*2+10, y, &gopdf.Rect{W: form.pdf.GetFullWidth() - form.pdf.MarginLeft()*2 - 10, H: titleHeight})
	form.pdf.SetBoldFont()
	form.pdf.SetWhiteFont()
	form.pdf.SetGrayStroke()
	form.pdf.SetXY(x+form.pdf.MarginLeft()*2+10, y)
	_ = form.pdf.CellWithOption(&gopdf.Rect{W: col1Width, H: titleHeight}, " Demographic Info", col1Opt)
	_ = form.pdf.CellWithOption(&gopdf.Rect{W: col2Width, H: titleHeight}, "YES", col2Opt)
	_ = form.pdf.CellWithOption(&gopdf.Rect{W: col3Width, H: titleHeight}, "NO", col3Opt)
	form.pdf.SetBlackFont()
	form.pdf.SetRegularFont()
	form.pdf.Br(titleHeight)
	x = form.pdf.GetX()
	y = form.pdf.GetY()
	text := " Was the ethnicity of the Borrower collected on the basis of visual observation or surname?"
	_ = form.buildPage7TableLine(text, x+form.pdf.MarginLeft()*2, y, form.pdf.LineHeight, col1Width+10, col2Width, col3Width, false, "")
	x = form.pdf.GetX()
	y = form.pdf.GetY()
	text = " Was the sex of the Borrower collected on the basis of visual observation or surname?"
	_ = form.buildPage7TableLine(text, x+form.pdf.MarginLeft()*2, y, form.pdf.LineHeight, col1Width+10, col2Width, col3Width, false, "")
	x = form.pdf.GetX()
	y = form.pdf.GetY()
	text = " Was the race of the Borrower collected on the basis of visual observation or surname?"
	_ = form.buildPage7TableLine(text, x+form.pdf.MarginLeft()*2, y, form.pdf.LineHeight, col1Width+10, col2Width, col3Width, false, "")

	form.pdf.Br(form.pdf.LineHeight)
	x = form.pdf.GetX()
	y = form.pdf.GetY()
	_ = form.buildTitleWithBg(" The Demographic Information Was Provided Through:", position{
		w: form.pdf.GetFullWidth() - form.pdf.MarginLeft() - form.pdf.MarginLeft(),
		h: 27,
		x: x + form.pdf.MarginLeft()*2,
		y: y,
	},
		position{
			w: col1Width,
			h: 27,
			x: x + form.pdf.MarginLeft()*2 + 10,
			y: y,
		})
	form.pdf.Br(form.pdf.LineHeight * 2.5)
	x = form.pdf.GetX()
	y = form.pdf.GetY()
	startX := x + form.pdf.MarginLeft()*2.5
	_ = form.pdf.ImageFrom(pdf.EmptyCheckBox, startX, y, checkBoxOpt)
	form.pdf.SetX(startX + form.pdf.LineHeight*1.5)
	_ = form.pdf.CellWithOption(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: form.pdf.LineHeight}, "Face to Face (Video Conf)", cellTextOpt)
	form.pdf.Br(form.pdf.LineHeight * 1.3)
	x = form.pdf.GetX()
	y = form.pdf.GetY()
	_ = form.pdf.ImageFrom(pdf.EmptyCheckBox, startX, y, checkBoxOpt)
	form.pdf.SetX(startX + form.pdf.LineHeight*1.5)
	_ = form.pdf.CellWithOption(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: form.pdf.LineHeight}, "Telephone Interview", cellTextOpt)
	form.pdf.Br(form.pdf.LineHeight * 1.3)
	x = form.pdf.GetX()
	y = form.pdf.GetY()
	_ = form.pdf.ImageFrom(pdf.EmptyCheckBox, startX, y, checkBoxOpt)
	form.pdf.SetX(startX + form.pdf.LineHeight*1.5)
	_ = form.pdf.CellWithOption(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: form.pdf.LineHeight}, "Fax or Mail", cellTextOpt)
	form.pdf.Br(form.pdf.LineHeight * 1.3)
	x = form.pdf.GetX()
	y = form.pdf.GetY()
	_ = form.pdf.ImageFrom(pdf.EmptyCheckBox, startX, y, checkBoxOpt)
	form.pdf.SetX(startX + form.pdf.LineHeight*1.5)
	_ = form.pdf.CellWithOption(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: form.pdf.LineHeight}, "Email or Internet", cellTextOpt)
	form.pdf.Br(form.pdf.LineHeight * 1.3)

	return err
}
