package apply_form

import (
	"github.com/signintech/gopdf"

	"zeebee-backend/app/borrower/api/internal/pdf"
)

func (form *ApplyForm) page5() (err error) {
	lineHeight := form.pdf.LineHeight
	fontHeight := form.pdf.LineHeight
	form.pdf.AddPage()
	err = form.setPageTitle(pdf.GlassIco, "PROPERTY |", "LOAN INFO")
	if err != nil {
		return err
	}
	// budget
	form.pdf.Br(lineHeight)
	x := form.pdf.GetX()
	y := form.pdf.GetY()
	_ = form.buildTitleWithBg("BUDGET DETAILS:", position{
		x: x,
		y: y,
		w: form.pdf.GetFullWidth(),
		h: 27,
	}, position{
		x: x + 10,
		y: y,
		w: form.pdf.GetFullWidth() - 100,
		h: 27,
	})
	form.pdf.Br(30)
	_ = form.pdf.Cell(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: fontHeight}, "Original Budget:")
	form.pdf.Br(lineHeight)
	_ = form.pdf.Cell(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: fontHeight}, "Remaining Budget:")
	form.pdf.Br(lineHeight)
	_ = form.pdf.Cell(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: fontHeight}, "Cost Spent to Date:")

	form.pdf.Br(lineHeight)
	x = form.pdf.GetX()
	y = form.pdf.GetY()
	_ = form.buildTitleWithBg("PROPERTY DETAILS:", position{
		x: x,
		y: y,
		w: form.pdf.GetFullWidth(),
		h: 27,
	}, position{
		x: x + 10,
		y: y,
		w: form.pdf.GetFullWidth() - 100,
		h: 27,
	})
	form.pdf.Br(lineHeight * 3)
	x = form.pdf.GetX()
	y = form.pdf.GetY()
	err = form.buildPage5PropertyStatusTable(position{x: x, y: y})
	if err != nil {
		return err
	}
	form.pdf.Br(lineHeight * 3)
	x = form.pdf.GetX()
	y = form.pdf.GetY()
	err = form.buildPage5ThirdService1Table(position{x: x, y: y})
	if err != nil {
		return err
	}

	form.pdf.Br(lineHeight * 3)
	x = form.pdf.GetX()
	y = form.pdf.GetY()
	err = form.buildPage5ThirdService2Table(position{x: x, y: y})
	if err != nil {
		return err
	}

	return err
}

func (form *ApplyForm) buildPage5PropertyStatusTable(tablePosition position) (err error) {
	lineWidth := 0.5
	table := form.pdf.NewTableLayout(tablePosition.x, tablePosition.y, 30, 5)

	table.SetHeaderStyle(gopdf.CellStyle{
		BorderStyle: gopdf.BorderStyle{
			Top:      true,
			Left:     true,
			Bottom:   true,
			Right:    true,
			Width:    lineWidth,
			RGBColor: gopdf.RGBColor{R: 0, G: 0, B: 0},
		},
		FillColor: gopdf.RGBColor{R: 71, G: 117, B: 255},
		TextColor: gopdf.RGBColor{R: 255, G: 255, B: 255},
		Font:      form.pdf.BoldFontName,
		FontSize:  14,
	})

	table.SetCellStyle(gopdf.CellStyle{
		BorderStyle: gopdf.BorderStyle{
			Right:    true,
			Bottom:   true,
			Left:     true,
			Top:      true,
			Width:    lineWidth,
			RGBColor: gopdf.RGBColor{R: 0, G: 0, B: 0},
		},
		FillColor: gopdf.RGBColor{R: 255, G: 255, B: 255},
		TextColor: gopdf.RGBColor{R: 0, G: 0, B: 0},
		Font:      form.pdf.RegularFontName,
		FontSize:  form.pdf.RegularFontSize,
	})

	// Add columns to the table
	table.AddColumn("", 100, "center")
	table.AddColumn("Current Status", 150, "center")
	table.AddColumn("Future Status", 150, "center")

	table.AddRow([]string{"Property Type", "", ""})
	table.AddRow([]string{"Living Area SQFT", "", ""})
	table.AddRow([]string{"Number of Units", "", ""})
	table.AddRow([]string{"Number of Bedrooms", "", ""})
	table.AddRow([]string{"Number of Bathrooms", "", ""})
	err = table.DrawTable()

	return err
}

func (form *ApplyForm) buildPage5ThirdService1Table(tablePosition position) (err error) {
	lineWidth := 0.5
	table := form.pdf.NewTableLayout(tablePosition.x, tablePosition.y, 30, 4)

	table.SetHeaderStyle(gopdf.CellStyle{
		BorderStyle: gopdf.BorderStyle{
			Top:      true,
			Left:     true,
			Bottom:   true,
			Right:    true,
			Width:    lineWidth,
			RGBColor: gopdf.RGBColor{R: 0, G: 0, B: 0},
		},
		FillColor: gopdf.RGBColor{R: 71, G: 117, B: 255},
		TextColor: gopdf.RGBColor{R: 255, G: 255, B: 255},
		Font:      form.pdf.BoldFontName,
		FontSize:  14,
	})

	table.SetCellStyle(gopdf.CellStyle{
		BorderStyle: gopdf.BorderStyle{
			Right:    true,
			Bottom:   true,
			Left:     true,
			Top:      true,
			Width:    lineWidth,
			RGBColor: gopdf.RGBColor{R: 0, G: 0, B: 0},
		},
		FillColor: gopdf.RGBColor{R: 255, G: 255, B: 255},
		TextColor: gopdf.RGBColor{R: 0, G: 0, B: 0},
		Font:      form.pdf.RegularFontName,
		FontSize:  form.pdf.RegularFontSize,
	})

	// Add columns to the table
	table.AddColumn("Contact Info", 100, "center")
	table.AddColumn("Title Co.", 150, "center")
	table.AddColumn("Escrow Co.", 150, "center")
	table.AddColumn("Insurance", 150, "center")

	table.AddRow([]string{"Name", "", "", ""})
	table.AddRow([]string{"Company", "", "", ""})
	table.AddRow([]string{"Phone", "", "", ""})
	table.AddRow([]string{"Email", "", "", ""})
	err = table.DrawTable()

	return err
}

func (form *ApplyForm) buildPage5ThirdService2Table(tablePosition position) (err error) {
	lineWidth := 0.5
	table := form.pdf.NewTableLayout(tablePosition.x, tablePosition.y, 30, 3)

	table.SetHeaderStyle(gopdf.CellStyle{
		BorderStyle: gopdf.BorderStyle{
			Top:      true,
			Left:     true,
			Bottom:   true,
			Right:    true,
			Width:    lineWidth,
			RGBColor: gopdf.RGBColor{R: 0, G: 0, B: 0},
		},
		FillColor: gopdf.RGBColor{R: 71, G: 117, B: 255},
		TextColor: gopdf.RGBColor{R: 255, G: 255, B: 255},
		Font:      form.pdf.BoldFontName,
		FontSize:  14,
	})

	table.SetCellStyle(gopdf.CellStyle{
		BorderStyle: gopdf.BorderStyle{
			Right:    true,
			Bottom:   true,
			Left:     true,
			Top:      true,
			Width:    lineWidth,
			RGBColor: gopdf.RGBColor{R: 0, G: 0, B: 0},
		},
		FillColor: gopdf.RGBColor{R: 255, G: 255, B: 255},
		TextColor: gopdf.RGBColor{R: 0, G: 0, B: 0},
		Font:      form.pdf.RegularFontName,
		FontSize:  form.pdf.RegularFontSize,
	})

	// Add columns to the table
	table.AddColumn("Contact Info", 100, "center")
	table.AddColumn("Appraisal Scheduling/ Access Contact", 250, "center")
	table.AddColumn("Appraisal Payment Contact", 200, "center")

	table.AddRow([]string{"Name", "", ""})
	table.AddRow([]string{"Phone", "", ""})
	table.AddRow([]string{"Email", "", ""})
	err = table.DrawTable()

	return err
}
