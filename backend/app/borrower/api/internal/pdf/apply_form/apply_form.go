package apply_form

import (
	"encoding/base64"
	"fmt"
	"image"
	"io"
	"strings"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/pkg/errors"
	"github.com/signintech/gopdf"

	"zeebee-backend/app/borrower/api/internal/pdf"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/common/enums"
	"zeebee-backend/common/errorx"
)

var checkBoxWidth = 17.0

type position struct {
	x, y float64
	w, h float64
}

type ApplyForm struct {
	pdf  *pdf.Pdf
	info *types.ApplyFormPdfRequest
}

func NewApplyForm(info *types.ApplyFormPdfRequest) (form *ApplyForm, err error) {
	if info == nil {
		return nil, errors.New("pdf info is empty")
	}
	p, err := pdf.NewPdf(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	if err != nil {
		return nil, err
	}
	return &ApplyForm{
		pdf:  p,
		info: info,
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
		func() error {
			return form.page4()
		},
		func() error {
			return form.page5()
		},
		func() error {
			return form.page6()
		},
		func() error {
			return form.page7()
		},
		func() error {
			return form.page8()
		},
		func() error {
			return form.page9()
		},
		func() error {
			return form.page10()
		},
		func() error {
			return form.page11()
		},
		func() error {
			return form.page12()
		},
		func() error {
			return form.page13()
		},
		func() error {
			return form.page14()
		},
		func() error {
			return form.page15()
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

	title2 = strings.Trim(title2, " ")
	title2 = " " + title2

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
		startX += icoWith * 1.3
	}

	form.pdf.SetXY(startX, y+3)
	form.pdf.SetBoldFontWithSize(18)
	form.pdf.SetBlackFont()

	_ = form.pdf.Cell(nil, title1)
	form.setBlueFont()
	_ = form.pdf.Cell(nil, title2)

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

func (form *ApplyForm) buildRadioCheckbox(startX, startY float64, optName string, targetValue string) (err error) {

	opts := enums.Get(optName)
	currentX := startX
	currentY := startY
	for _, opt := range opts {
		itemLen, _ := form.pdf.MeasureTextWidth(opt.Label)
		checkBox := pdf.EmptyCheckBox
		if opt.Value == targetValue {
			checkBox = pdf.SelectedCheckBox
		}

		if currentX+checkBoxWidth+5+itemLen > form.pdf.GetFullWidth()+form.pdf.MarginLeft() {
			currentX = startX
			currentY += checkBoxWidth * 1.2
		}

		err = form.pdf.ImageFrom(checkBox, currentX, currentY, &gopdf.Rect{W: checkBoxWidth, H: checkBoxWidth})
		if err != nil {
			return err
		}
		form.pdf.SetXY(currentX+checkBoxWidth+5, currentY+5)
		_ = form.pdf.Cell(nil, opt.Label)
		currentX += itemLen + 30
	}
	form.pdf.SetY(currentY + checkBoxWidth)
	form.pdf.Br(1)
	return err
}

func (form *ApplyForm) buildMultiCheckbox(startX, startY float64, optName string, targetValues []string) (err error) {
	opts := enums.Get(optName)
	currentX := startX
	currentY := startY
	for _, opt := range opts {
		itemLen, _ := form.pdf.MeasureTextWidth(opt.Label)
		checkBox := pdf.EmptyCheckBox

		if slice.Contain(targetValues, opt.Value) {
			checkBox = pdf.SelectedCheckBox
		}
		if currentX+checkBoxWidth+5+itemLen > form.pdf.GetFullWidth()+form.pdf.MarginLeft() {
			currentX = startX
			currentY += checkBoxWidth * 1.2
		}

		err = form.pdf.ImageFrom(checkBox, currentX, currentY, &gopdf.Rect{W: checkBoxWidth, H: checkBoxWidth})
		if err != nil {
			return err
		}
		form.pdf.SetXY(currentX+checkBoxWidth+5, currentY+5)
		_ = form.pdf.Cell(nil, opt.Label)
		currentX += itemLen + 30
	}
	form.pdf.SetY(currentY + checkBoxWidth)
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
	form.pdf.SetTextColor(71, 117, 255)
	form.pdf.SetXY(position.x, position.y)
	_ = form.pdf.CellWithOption(&gopdf.Rect{W: position.w, H: position.h}, title, gopdf.CellOption{
		Align: gopdf.Left | gopdf.Middle,
	})
	form.pdf.SetBlackRegularFont()
	return err
}

func (form *ApplyForm) buildSignature(typ string) (err error) {
	signatures := form.info.Signatures
	targetSignature := signatures[typ]

	form.pdf.SetBoldFont()
	_ = form.pdf.MultiCell(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: 12}, "Agreed & Accepted")
	form.pdf.Br(20)

	for i := 0; i < 2; i++ {
		var signature *types.SignatureEntry
		if i > len(targetSignature)-1 {
			signature = &types.SignatureEntry{}
		} else {
			signature = targetSignature[i]
		}
		_ = form.pdf.Cell(&gopdf.Rect{W: 140, H: 12}, fmt.Sprintf("Applicant %d Signature:", i+1))
		_ = form.pdf.CellWithOption(&gopdf.Rect{W: 200, H: 12}, "", gopdf.CellOption{
			Align:  gopdf.Left,
			Border: gopdf.Bottom,
		})
		signatureImg, err := base64ToImage(signature.SignatureImage)
		if err != nil {
			return err
		}
		if signatureImg != nil {
			err = form.pdf.ImageFrom(signatureImg, 200, form.pdf.GetY()-20, &gopdf.Rect{W: 100, H: 30})
			if err != nil {
				return err
			}
		}

		_ = form.pdf.Cell(&gopdf.Rect{W: 20, H: 12}, "")
		_ = form.pdf.Cell(&gopdf.Rect{W: 40, H: 12}, "Date:")
		_ = form.pdf.CellWithOption(&gopdf.Rect{W: 80, H: 12}, signature.SignatureDate, gopdf.CellOption{
			Align:  gopdf.Left,
			Border: gopdf.Bottom,
		})
		form.pdf.Br(1)
		form.pdf.Br(40)

	}
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

type colorText struct {
	text  string
	color struct {
		r, g, b uint8
	}
}

func (form *ApplyForm) buildLongTextWithColor(contents []colorText, startX float64, fontHeight float64) (err error) {
	breakOption := &gopdf.BreakOption{
		Mode:           gopdf.BreakModeIndicatorSensitive,
		BreakIndicator: ' ', // 在空格处分割
	}
	tmpX := startX

OutLoop:
	for _, content := range contents {
		form.pdf.SetTextColor(content.color.r, content.color.g, content.color.b)
		text := content.text
		// start from prev position
		if tmpX > startX {
			leftLen := form.pdf.GetFullWidth() - (tmpX - form.pdf.MarginLeft())
			form.pdf.SetX(tmpX)
			lines, _ := form.pdf.SplitTextWithOption(text, leftLen, breakOption)
			_ = form.pdf.CellWithOption(&gopdf.Rect{W: leftLen, H: fontHeight}, lines[0], gopdf.CellOption{
				Align: gopdf.Left | gopdf.Middle,
			})
			form.pdf.Br(fontHeight)
			text = text[len(lines[0]):]
			tmpX = startX
		}

		if len(text) == 0 {
			continue OutLoop
		}

		// 将文本分割成多行，每行宽度不超过指定值
		lines, _ := form.pdf.SplitTextWithOption(text, form.pdf.GetFullWidth(), breakOption)
		lineCount := len(lines)
	InnerLoop:
		for i, line := range lines {
			if i < lineCount-1 {
				_ = form.pdf.CellWithOption(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: fontHeight}, line, gopdf.CellOption{
					Align: gopdf.Left | gopdf.Middle,
				})
				form.pdf.Br(fontHeight)
				continue InnerLoop
			}

			lineWidth, _ := form.pdf.MeasureTextWidth(line)
			tmpX += lineWidth
			_ = form.pdf.CellWithOption(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: fontHeight}, line, gopdf.CellOption{
				Align: gopdf.Left | gopdf.Middle,
			})

		}
	}
	form.pdf.Br(fontHeight)

	return err

}

func (form *ApplyForm) buildLists(lists []string, position position, fontHeight float64) (err error) {
	point := "·"
	col1Width := 10.0
	col2Width := position.w - col1Width - 10.0
	tempY := position.y
	for _, list := range lists {
		lines, err := form.pdf.SplitText(list, col2Width)
		if err != nil {
			return err
		}
		form.pdf.SetXY(position.x, tempY)
		_ = form.pdf.CellWithOption(&gopdf.Rect{W: col1Width, H: fontHeight}, point, gopdf.CellOption{
			Align: gopdf.Left,
		})

		for _, line := range lines {
			form.pdf.SetXY(position.x+col1Width, tempY)
			_ = form.pdf.CellWithOption(&gopdf.Rect{W: col2Width, H: fontHeight}, line, gopdf.CellOption{
				Align: gopdf.Left | gopdf.Middle,
			})
			tempY += fontHeight
		}
	}
	return err
}

func formateAddress(addr *types.AddressEntry) string {
	if addr == nil {
		return ""
	}
	return fmt.Sprintf("%s %s %s %s", addr.Address1, addr.City, addr.State, addr.ZipCode)

}

func base64ToImage(content string) (img image.Image, err error) {
	if content == "" {
		return nil, nil
	}
	// remove base64 header
	content = strings.TrimPrefix(content, "data:image/png;base64,")
	content = strings.TrimPrefix(content, "data:image/jpeg;base64,")
	content = strings.TrimPrefix(content, "data:image/jpg;base64,")
	img, _, err = image.Decode(base64.NewDecoder(base64.StdEncoding, strings.NewReader(content)))

	return img, err
}
