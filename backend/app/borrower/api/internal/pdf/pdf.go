package pdf

import (
	"bytes"
	"embed"
	"image"

	"github.com/signintech/gopdf"
	"golang.org/x/sync/errgroup"
)

//go:embed statics
var StaticsFs embed.FS

//go:embed statics/font/LiberationSerif-Regular.ttf
var RegularFont []byte

//go:embed statics/font/LiberationSerif-Bold-1.ttf
var BoldFont []byte

//go:embed statics/images/gray_bg.png
var grayBgFile []byte

//go:embed statics/images/white_bg.png
var whiteBgFile []byte

//go:embed statics/images/blue_bg.png
var blueBgFile []byte

//go:embed statics/images/logo.png
var LogoFile []byte

//go:embed statics/images/light_blue_bg.png
var lightBlueBgFile []byte

//go:embed statics/images/book_ico.png
var bookIcoFile []byte

//go:embed statics/images/note_ico.png
var noteIcoFile []byte

//go:embed statics/images/list_ico.png
var listIcoFile []byte

//go:embed statics/images/glass_ico.png
var glassIcoFile []byte

//go:embed statics/images/chapter_ico.png
var chapterIcoFile []byte

//go:embed statics/images/empty_checkbox.png
var emptyCheckBoxFile []byte

//go:embed statics/images/selected_checkbox.png
var selectedCheckBoxFile []byte

var RegularFontName = "Regular"
var BoldFontName = "Bold"

var GrayBg image.Image
var BlueBg image.Image
var WhiteBg image.Image
var Logo image.Image
var LightBlueBg image.Image
var BookIco image.Image
var NoteIco image.Image
var ListIco image.Image
var GlassIco image.Image
var ChapterIco image.Image
var EmptyCheckBox image.Image
var SelectedCheckBox image.Image

var defaultRegularFontSize = 11.0
var defaultBoldFontSize = 11.0
var defaultLineHeight = 14.0

func init() {
	GrayBg, _, _ = image.Decode(bytes.NewReader(grayBgFile))
	BlueBg, _, _ = image.Decode(bytes.NewReader(blueBgFile))
	WhiteBg, _, _ = image.Decode(bytes.NewReader(whiteBgFile))
	LightBlueBg, _, _ = image.Decode(bytes.NewReader(lightBlueBgFile))
	BookIco, _, _ = image.Decode(bytes.NewReader(bookIcoFile))
	NoteIco, _, _ = image.Decode(bytes.NewReader(noteIcoFile))
	ListIco, _, _ = image.Decode(bytes.NewReader(listIcoFile))
	GlassIco, _, _ = image.Decode(bytes.NewReader(glassIcoFile))
	ChapterIco, _, _ = image.Decode(bytes.NewReader(chapterIcoFile))
	EmptyCheckBox, _, _ = image.Decode(bytes.NewReader(emptyCheckBoxFile))
	SelectedCheckBox, _, _ = image.Decode(bytes.NewReader(selectedCheckBoxFile))
	Logo, _, _ = image.Decode(bytes.NewReader(LogoFile))
}

type Pdf struct {
	*gopdf.GoPdf

	pdfConfig       gopdf.Config
	RegularFontSize float64
	BoldFontSize    float64
	RegularFontName string
	BoldFontName    string
	LineHeight      float64
}

type OptFunc func(*Pdf)

func WithRegularFontSize(size float64) OptFunc {
	return func(pdf *Pdf) {
		pdf.RegularFontSize = size
	}
}
func WithBoldFontSize(size float64) OptFunc {
	return func(pdf *Pdf) {
		pdf.BoldFontSize = size
	}
}
func WithLineHeight(height float64) OptFunc {
	return func(pdf *Pdf) {
		pdf.LineHeight = height
	}
}

func NewPdf(config gopdf.Config, opts ...OptFunc) (pdf *Pdf, err error) {
	p := &gopdf.GoPdf{}
	p.Start(config)

	pdf = &Pdf{
		GoPdf:           p,
		pdfConfig:       config,
		RegularFontSize: defaultRegularFontSize,
		BoldFontSize:    defaultBoldFontSize,
		LineHeight:      defaultLineHeight,
		RegularFontName: RegularFontName,
		BoldFontName:    BoldFontName,
	}
	for _, opt := range opts {
		opt(pdf)
	}
	err = pdf.init()

	return pdf, err
}

func (pdf *Pdf) init() (err error) {
	eg := errgroup.Group{}
	eg.Go(func() error {
		return pdf.AddTTFFontData(RegularFontName, RegularFont)
	})
	eg.Go(func() error {
		return pdf.AddTTFFontData(BoldFontName, BoldFont)
	})
	err = eg.Wait()
	return err
}

// SetFooter sets the footer
func (pdf *Pdf) SetFooter(f func()) {
	pdf.AddFooter(f)
}

// set font

// SetRegularFont sets the font to regular
func (pdf *Pdf) SetRegularFont() {
	_ = pdf.SetFont(pdf.RegularFontName, "", pdf.RegularFontSize)
}

// SetBoldFont sets the font to bold
func (pdf *Pdf) SetBoldFont() {
	_ = pdf.SetFont(pdf.BoldFontName, "", pdf.BoldFontSize)
}

// SetBoldFontWithSize sets the font to bold with a specific size
func (pdf *Pdf) SetBoldFontWithSize(size float64) {
	_ = pdf.SetFont(pdf.BoldFontName, "", size)
}

// SetFontColor sets the font color
func (pdf *Pdf) SetFontColor(r uint8, g uint8, b uint8) {
	pdf.SetTextColor(r, g, b)
}

// SetWhiteFont sets the font color to white
func (pdf *Pdf) SetWhiteFont() {
	pdf.SetTextColor(255, 255, 255)
}

// SetBlackFont sets the font color to black
func (pdf *Pdf) SetBlackFont() {
	pdf.SetTextColor(0, 0, 0)
}

// SetBlackBoldFont sets the font color to black and bold
func (pdf *Pdf) SetBlackBoldFont() {
	pdf.SetBlackFont()
	pdf.SetBoldFont()
}

// SetWhiteBoldFont sets the font color to white and bold
func (pdf *Pdf) SetWhiteBoldFont() {
	pdf.SetWhiteFont()
	pdf.SetBoldFont()
}

// SetBlackRegularFont sets the font color to black and regular
func (pdf *Pdf) SetBlackRegularFont() {
	pdf.SetBlackFont()
	pdf.SetRegularFont()
}

// SetWhiteRegularFont sets the font color to white and regular
func (pdf *Pdf) SetWhiteRegularFont() {
	pdf.SetWhiteFont()
	pdf.SetRegularFont()
}

// LongText adds a long text with line breaks
func (pdf *Pdf) LongText(lineHeight float64, text string) (err error) {
	breakOption := &gopdf.BreakOption{
		Mode:           gopdf.BreakModeIndicatorSensitive,
		BreakIndicator: ' ',
	}

	lines, err := pdf.SplitTextWithOption(text, gopdf.PageSizeA4.W-pdf.MarginLeft()-pdf.MarginRight(), breakOption)
	if err != nil {
		return err
	}

	for _, line := range lines {
		_ = pdf.Cell(&gopdf.Rect{W: pdf.GetFullWidth(), H: lineHeight}, line)
		pdf.Br(lineHeight)
	}
	return
}

// MultiLongText adds multiple long texts with line breaks
func (pdf *Pdf) MultiLongText(lineHeight float64, text []string) (err error) {
	for _, s := range text {
		err = pdf.LongText(lineHeight, s)
		if err != nil {
			return err
		}
		pdf.Br(lineHeight)
	}
	return
}

func (pdf *Pdf) OneFullLine(text string) (err error) {
	pdf.Br(0)
	err = pdf.CellWithOption(&gopdf.Rect{W: pdf.GetFullWidth(), H: pdf.LineHeight}, text, gopdf.CellOption{
		Align: gopdf.Left | gopdf.Middle,
	})
	pdf.Br(pdf.LineHeight)
	return
}

// GetFullPageWidth returns the width of the full page
func (pdf *Pdf) GetFullPageWidth() float64 {
	return pdf.pdfConfig.PageSize.W
}

// GetFullWidth returns the width of the full page minus the left and right margins
func (pdf *Pdf) GetFullWidth() float64 {
	return pdf.pdfConfig.PageSize.W - pdf.MarginLeft() - pdf.MarginRight()
}

// position

// GetPageStartX returns the x coordinate of the start of the page
func (pdf *Pdf) GetPageStartX() float64 {
	return pdf.MarginLeft()
}

// GetPageEndX returns the x coordinate of the end of the page
func (pdf *Pdf) GetPageEndX() float64 {
	return pdf.pdfConfig.PageSize.W - pdf.MarginLeft()
}

// BrLine adds a line break with default line height
func (pdf *Pdf) BrLine() {
	pdf.Br(pdf.LineHeight)
}

func (pdf *Pdf) SetGrayStroke() {
	pdf.GoPdf.SetGrayStroke(0.6)
}

func (pdf *Pdf) SetWhiteStroke() {
	pdf.GoPdf.SetGrayStroke(1)
}

func (pdf *Pdf) SetBlackStroke() {
	pdf.GoPdf.SetGrayStroke(0)
}
