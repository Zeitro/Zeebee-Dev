package apply_form

import (
	"github.com/signintech/gopdf"

	"zeebee-backend/app/borrower/api/internal/pdf"
)

func (form *ApplyForm) page14() (err error) {
	form.pdf.AddPage()
	err = form.setPageTitle(nil, "OVERFLOW", "TEXT PAGE")
	if err != nil {
		return err
	}
	form.pdf.Br(form.pdf.LineHeight)
	x := form.pdf.GetX()
	y := form.pdf.GetY()
	_ = form.buildTitleWithBg("Additional Information:", position{w: form.pdf.GetFullWidth(), h: 27, x: x, y: y}, position{w: form.pdf.GetFullWidth() - 10, h: 27, x: x + 10, y: y})
	form.pdf.Br(form.pdf.LineHeight * 3)

	_ = form.pdf.ImageFrom(pdf.LightBlueBg, x, y, &gopdf.Rect{W: form.pdf.GetFullWidth(), H: gopdf.PageSizeA4.H - form.pdf.LineHeight*6})
	return err
}
