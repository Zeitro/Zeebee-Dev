package apply_form

import (
	"fmt"
	"image"
	"io"

	"github.com/signintech/gopdf"

	"zeebee-backend/app/borrower/api/internal/pdf"
	"zeebee-backend/common/errorx"
)

var checkBoxWidth = 17.0

type position struct {
	x, y float64
	w, h float64
}

type ApplyForm struct {
	pdf *pdf.Pdf
}

func NewApplyForm() (form *ApplyForm, err error) {
	p, err := pdf.NewPdf(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	if err != nil {
		return nil, err
	}
	return &ApplyForm{
		pdf: p,
	}, nil
}

// Generate writes the pdf to the given writer
func (form *ApplyForm) Generate(w io.Writer) (err error) {
	form.pdf.SetFooter(func() {
		currentX := form.pdf.GetX()
		currentY := form.pdf.GetY()
		pageHeight := gopdf.PageSizeA4.H
		form.pdf.SetXY(0, pageHeight-20)
		form.pdf.SetRegularFont()
		footerText := fmt.Sprintf("- %d -", form.pdf.GetNumberOfPages())
		textWidth, _ := form.pdf.MeasureTextWidth(footerText)
		pageWidth := gopdf.PageSizeA4.W
		form.pdf.SetX((pageWidth - textWidth) / 2)
		_ = form.pdf.Cell(nil, footerText)
		form.pdf.SetXY(currentX, currentY)
	})

	eg := errorx.NewErrorGroup()
	eg.Do(
		func() error {
			return form.page1()
		},
		func() error {
			return form.page2()
		},
		func() error {
			return form.page3()
		},

		// write the pdf
		func() error {
			_, err := form.pdf.WriteTo(w)
			return err
		},
	)
	err = eg.Error()
	return err
}

// setPageTitle sets the page title
func (form *ApplyForm) setPageTitle(icoImage image.Image, title1 string, title2 string) (err error) {
	var icoWith float64 = 20
	titleLen1, err := form.pdf.MeasureTextWidth(title1)
	if err != nil {
		return err
	}

	titleLen2, err := form.pdf.MeasureTextWidth(title2)
	if err != nil {
		return err
	}

	totalLen := titleLen1 + titleLen2
	y := form.pdf.GetY()
	startX := form.pdf.GetFullWidth()/2.0 - totalLen/2.0 - form.pdf.MarginLeft()
	if icoImage != nil {
		startX = form.pdf.GetFullWidth()/2.0 - (totalLen+icoWith+10)/2.0 - form.pdf.MarginLeft() - 20
		err = form.pdf.ImageFrom(icoImage, startX, y, &gopdf.Rect{W: icoWith, H: icoWith})
		if err != nil {
			return err
		}
		startX += icoWith
	}

	form.pdf.SetXY(startX, y)
	form.pdf.SetBoldFontWithSize(18)
	form.pdf.SetBlackFont()

	_ = form.pdf.Cell(nil, title1)
	form.setBlueFont()
	_ = form.pdf.Cell(nil, " "+title2)

	form.pdf.SetY(40)
	form.pdf.Br(1)
	form.pdf.SetBlackRegularFont()

	return err
}

func (form *ApplyForm) setBlueFont() {
	form.pdf.SetFontColor(71, 117, 255)
}

func (form *ApplyForm) buildCheckbox(startX, startY float64, items []string) (err error) {

	for _, item := range items {
		itemLen, _ := form.pdf.MeasureTextWidth(item)
		err = form.pdf.ImageFrom(pdf.EmptyCheckBox, startX, startY, &gopdf.Rect{W: checkBoxWidth, H: checkBoxWidth})
		if err != nil {
			return err
		}
		form.pdf.SetXY(startX+checkBoxWidth+5, startY+5)
		_ = form.pdf.Cell(nil, item)
		startX += itemLen + 30
	}
	form.pdf.SetY(startY + 17)
	form.pdf.Br(1)

	return err
}

func (form *ApplyForm) buildTitleWithBg(title string, grayBgPos position, blueBgPos position) (err error) {
	_ = form.pdf.ImageFrom(pdf.GrayBg, grayBgPos.x, grayBgPos.y, &gopdf.Rect{W: grayBgPos.w, H: grayBgPos.h})
	_ = form.pdf.ImageFrom(pdf.BlueBg, blueBgPos.x, blueBgPos.y, &gopdf.Rect{W: blueBgPos.w, H: blueBgPos.h})
	form.pdf.SetWhiteBoldFont()
	form.pdf.SetXY(blueBgPos.x+10, blueBgPos.y)
	_ = form.pdf.CellWithOption(&gopdf.Rect{W: blueBgPos.w, H: blueBgPos.h}, title, gopdf.CellOption{
		Align: gopdf.Left | gopdf.Middle,
	})
	form.pdf.SetBlackRegularFont()
	return err
}

func (form *ApplyForm) buildSmallTitle(title string, position position) (err error) {
	_ = form.pdf.ImageFrom(pdf.BlueBg, position.x, position.y, &gopdf.Rect{W: position.w, H: position.h})
	form.pdf.SetWhiteBoldFont()
	form.pdf.SetXY(position.x, position.y)
	_ = form.pdf.CellWithOption(&gopdf.Rect{W: position.w, H: position.h}, title, gopdf.CellOption{
		Align: gopdf.Center | gopdf.Middle,
	})
	form.pdf.SetBlackRegularFont()
	return err
}

func (form *ApplyForm) buildGraySmallTitle(title string, position position) (err error) {
	_ = form.pdf.ImageFrom(pdf.GrayBg, position.x, position.y, &gopdf.Rect{W: position.w, H: position.h})

	form.pdf.SetWhiteBoldFont()
	_ = form.pdf.SetFontSize(10)
	form.pdf.SetXY(position.x, position.y)
	_ = form.pdf.CellWithOption(&gopdf.Rect{W: position.w, H: position.h}, title, gopdf.CellOption{
		Align: gopdf.Left | gopdf.Middle,
	})
	form.pdf.SetBlackRegularFont()
	return err
}

func (form *ApplyForm) buildSignature() (err error) {
	form.pdf.SetBoldFont()
	_ = form.pdf.MultiCell(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: 12}, "Agreed & Accepted")
	form.pdf.Br(20)
	_ = form.pdf.Cell(&gopdf.Rect{W: 140, H: 12}, "Applicant 1 Signature:")
	_ = form.pdf.CellWithOption(&gopdf.Rect{W: 200, H: 12}, "", gopdf.CellOption{
		Align:  gopdf.Left,
		Border: gopdf.Bottom,
	})
	_ = form.pdf.Cell(&gopdf.Rect{W: 20, H: 12}, "")
	_ = form.pdf.Cell(&gopdf.Rect{W: 40, H: 12}, "Date:")
	_ = form.pdf.CellWithOption(&gopdf.Rect{W: 80, H: 12}, "", gopdf.CellOption{
		Align:  gopdf.Left,
		Border: gopdf.Bottom,
	})
	form.pdf.Br(1)
	form.pdf.Br(40)
	_ = form.pdf.Cell(&gopdf.Rect{W: 140, H: 12}, "Applicant 2 Signature:")
	_ = form.pdf.CellWithOption(&gopdf.Rect{W: 200, H: 12}, "", gopdf.CellOption{
		Align:  gopdf.Left,
		Border: gopdf.Bottom,
	})
	_ = form.pdf.Cell(&gopdf.Rect{W: 20, H: 12}, "")
	_ = form.pdf.Cell(&gopdf.Rect{W: 40, H: 12}, "Date:")
	_ = form.pdf.CellWithOption(&gopdf.Rect{W: 80, H: 12}, "", gopdf.CellOption{
		Align:  gopdf.Left,
		Border: gopdf.Bottom,
	})
	form.pdf.Br(1)
	form.pdf.SetRegularFont()
	return err
}

func (form *ApplyForm) OneFullLine(text string) (err error) {
	form.pdf.Br(0)
	err = form.pdf.CellWithOption(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: form.pdf.LineHeight}, text, gopdf.CellOption{
		Align: gopdf.Left | gopdf.Middle,
	})
	form.pdf.Br(form.pdf.LineHeight)
	return
}

// func (form *ApplyForm) MultiFullLine(text string) (err error) {
// 	form.pdf.Br(0)
// 	height, err := form.pdf.MeasureCellHeightByText(text)
// 	if err != nil {
// 		return err
// 	}
// 	err = form.pdf.MultiCellWithOption(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: height}, text, gopdf.CellOption{
// 		Align: gopdf.Left | gopdf.Middle,
// 	})
// 	form.pdf.Br(form.pdf.LineHeight)
// 	return
// }
