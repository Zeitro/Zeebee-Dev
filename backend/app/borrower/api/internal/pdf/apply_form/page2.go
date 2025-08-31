package apply_form

import (
	"fmt"

	"github.com/signintech/gopdf"

	"zeebee-backend/app/borrower/api/internal/pdf"
)

func (form *ApplyForm) page2() (err error) {
	form.pdf.AddPage()
	err = form.setPageTitle(pdf.NoteIco, "BORROWER", "INFORMATION")
	if err != nil {
		return err
	}
	form.pdf.SetBoldFont()
	_ = form.pdf.Cell(nil, "Borrower Type:")
	// form.pdf.Br(20)
	form.pdf.BrLine()
	x := form.pdf.GetX()
	y := form.pdf.GetY()
	startX := x + 50.0
	startY := y + 5
	boxWidth := 17.0
	items := []string{"Individual", " LLC", "Corporation", " Partnership", "Trust"}
	for _, item := range items {
		itemLen, _ := form.pdf.MeasureTextWidth(item)
		err = form.pdf.ImageFrom(pdf.EmptyCheckBox, startX, y, &gopdf.Rect{W: boxWidth, H: boxWidth})
		if err != nil {
			return err
		}
		form.pdf.SetXY(startX+boxWidth+5, startY)
		_ = form.pdf.Cell(nil, item)
		startX += itemLen + 30
	}
	// form.pdf.SetXY(x, y)
	// form.pdf.Br(15)
	form.pdf.BrLine()
	// 绘制灰色线条
	form.pdf.SetGrayStroke()
	form.pdf.Line(form.pdf.MarginLeft(), form.pdf.GetY(), gopdf.PageSizeA4.W-form.pdf.MarginRight(), form.pdf.GetY())
	form.pdf.SetRegularFont()
	rowHeight := 12.0
	brHeight := 6.0
	form.pdf.Br(brHeight)
	_ = form.pdf.Cell(&gopdf.Rect{W: 300, H: rowHeight}, "Entity Name/Individual: ")
	_ = form.pdf.Cell(&gopdf.Rect{W: 100, H: rowHeight}, "TIN: ")
	_ = form.pdf.Cell(&gopdf.Rect{W: 100, H: rowHeight}, "% Own:")
	form.pdf.Br(brHeight * 3.0)
	_ = form.pdf.MultiCell(&gopdf.Rect{W: gopdf.PageSizeA4.W - form.pdf.MarginLeft() - form.pdf.MarginRight(), H: rowHeight}, "Mailing Address: ")
	form.pdf.Br(brHeight)
	_ = form.pdf.Cell(&gopdf.Rect{W: 300, H: rowHeight}, "City: ")
	_ = form.pdf.Cell(&gopdf.Rect{W: 100, H: rowHeight}, "State: ")
	_ = form.pdf.Cell(&gopdf.Rect{W: 100, H: rowHeight}, "Zip:")
	form.pdf.Br(brHeight * 3.0)
	_ = form.pdf.MultiCell(&gopdf.Rect{W: gopdf.PageSizeA4.W - form.pdf.MarginLeft() - form.pdf.MarginRight(), H: rowHeight}, "Phone: ")
	form.pdf.Br(brHeight)
	_ = form.pdf.MultiCell(&gopdf.Rect{W: gopdf.PageSizeA4.W - form.pdf.MarginLeft() - form.pdf.MarginRight(), H: rowHeight}, "Email: ")
	form.pdf.Br(brHeight)
	_ = form.pdf.MultiCell(&gopdf.Rect{W: gopdf.PageSizeA4.W - form.pdf.MarginLeft() - form.pdf.MarginRight(), H: rowHeight}, "Contact Name: ")
	form.pdf.Br(brHeight)
	form.pdf.Line(form.pdf.MarginLeft(), form.pdf.GetY(), gopdf.PageSizeA4.W-form.pdf.MarginRight(), form.pdf.GetY())

	err = form.applicantInfo(1)
	if err != nil {
		return err
	}
	err = form.applicantInfo(2)
	if err != nil {
		return err
	}

	return err
}

func (form *ApplyForm) applicantInfo(index int) (err error) {
	brHeight := 6.0
	x := form.pdf.GetX()
	y := form.pdf.GetY()
	_ = form.pdf.ImageFrom(pdf.GrayBg, x, y, &gopdf.Rect{W: form.pdf.GetFullWidth(), H: 27})
	_ = form.pdf.ImageFrom(pdf.BlueBg, x+10, y, &gopdf.Rect{W: form.pdf.GetFullWidth() - 10, H: 27})

	form.pdf.SetWhiteBoldFont()
	form.pdf.SetXY(x+20, y+8)
	_ = form.pdf.Cell(nil, fmt.Sprintf("Applicant %d:", index))
	form.pdf.SetBlackRegularFont()
	form.pdf.Br(30)
	rowHeight := 12.0
	_ = form.pdf.MultiCell(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: rowHeight}, "First Name:")
	form.pdf.Br(brHeight)
	_ = form.pdf.MultiCell(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: rowHeight}, "Middle Name:")
	form.pdf.Br(brHeight)
	_ = form.pdf.MultiCell(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: rowHeight}, "Last Name:")
	form.pdf.Br(brHeight)
	_ = form.pdf.MultiCell(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: rowHeight}, "SSN:")
	form.pdf.Br(brHeight)
	_ = form.pdf.MultiCell(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: rowHeight}, "DOB:")
	form.pdf.Br(brHeight)
	_ = form.pdf.MultiCell(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: rowHeight}, "Phone:")
	form.pdf.Br(brHeight)
	_ = form.pdf.MultiCell(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: rowHeight}, "Marital Status:")
	x = form.pdf.GetX()
	y = form.pdf.GetY()
	startX := x + 50.0
	startY := y + 5
	// boxWidth := 17.0
	items := []string{"Married", "Unmarried", "Separated"}
	err = form.buildCheckbox(startX, startY, items)
	if err != nil {
		return err
	}
	form.pdf.Br(brHeight)
	_ = form.pdf.MultiCell(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: rowHeight}, "Email:")
	form.pdf.Br(brHeight)
	_ = form.pdf.MultiCell(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: rowHeight}, "Est. Credit Score:")
	form.pdf.Br(brHeight)
	_ = form.pdf.MultiCell(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: rowHeight}, "Primary Address:")
	form.pdf.Br(brHeight)
	_ = form.pdf.Cell(&gopdf.Rect{W: 100, H: rowHeight}, "City: ")
	_ = form.pdf.Cell(&gopdf.Rect{W: (form.pdf.GetFullWidth()) / 2.0, H: rowHeight}, "State: ")
	form.pdf.Br(brHeight)
	form.pdf.Br(15)
	_ = form.pdf.MultiCell(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: rowHeight}, "Zip:")
	form.pdf.Br(brHeight)
	_ = form.pdf.MultiCell(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: rowHeight}, "Drivers License #:")
	form.pdf.Br(brHeight)
	return err
}
