package apply_form

import "github.com/signintech/gopdf"

func (form *ApplyForm) page6() (err error) {
	form.pdf.AddPage()
	err = form.setPageTitle(nil, "EXPERIENCE |", "TRACK RECORD")
	if err != nil {
		return err
	}
	form.pdf.Br(form.pdf.LineHeight)
	x := form.pdf.GetX()
	y := form.pdf.GetY()
	_ = form.buildGraySmallTitle("Transaction History (Flips, Ground Up, Value Add Properties)  Note: Please exclude all homes that were a primary residence", position{
		x: x,
		y: y,
		w: form.pdf.GetFullWidth(),
		h: form.pdf.LineHeight,
	})
	form.pdf.Br(form.pdf.LineHeight)
	x = form.pdf.GetX()
	y = form.pdf.GetY()
	err = form.buildPage6TransactionTable(position{x: x, y: y})
	if err != nil {
		return err
	}
	form.pdf.Br(form.pdf.LineHeight * 2)
	x = form.pdf.GetX()
	y = form.pdf.GetY()
	_ = form.buildGraySmallTitle("Real Estate Owned (Properties held as Rentals)", position{
		x: x,
		y: y,
		w: form.pdf.GetFullWidth(),
		h: form.pdf.LineHeight,
	})
	form.pdf.Br(form.pdf.LineHeight)
	x = form.pdf.GetX()
	y = form.pdf.GetY()
	err = form.buildPage6RealEstateTable(position{x: x, y: y})
	if err != nil {
		return err
	}
	return err
}

func (form *ApplyForm) buildPage6TransactionTable(tablePosition position) (err error) {
	lineWidth := 0.5
	tableFontSize := 8.5
	table := form.pdf.NewTableLayout(tablePosition.x, tablePosition.y, 30, 10)

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
		FontSize:  tableFontSize,
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
		FontSize:  tableFontSize,
	})

	// Add columns to the table
	table.AddColumn("Address", 70, "center")
	table.AddColumn("Ownership on Title", 55, "center")
	table.AddColumn("Property Type", 64, "center")
	table.AddColumn("Purchase Price", 60, "center")
	table.AddColumn("Purchase Date", 60, "center")
	table.AddColumn("Rehab Budget/ Construction Cost", 65, "center")
	table.AddColumn("Sale Price $", 60, "center")
	table.AddColumn("Date Sold", 60, "center")
	table.AddColumn("Profit / Loss $", 60, "center")

	table.AddRow([]string{"2625 Augustine Drive\n Santa Clara CA", "Me", "SFR Attached", "234555", " 04/03/2025", "1234", "234556", " 04/17/2025", " 12345"})
	table.AddRow([]string{"2625 Augustine Drive Santa Clara CA", "Me", "SFR Attached", "234555", " 04/03/2025", "1234", "234556", " 04/17/2025", " 12345"})
	table.AddRow([]string{"2625 Augustine Drive Santa Clara CA", "Me", "SFR Attached", "234555", " 04/03/2025", "1234", "234556", " 04/17/2025", " 12345"})
	table.AddRow([]string{"2625 Augustine Drive Santa Clara CA", "Me", "SFR Attached", "234555", " 04/03/2025", "1234", "234556", " 04/17/2025", " 12345"})
	err = table.DrawTable()

	return err
}

func (form *ApplyForm) buildPage6RealEstateTable(tablePosition position) (err error) {
	lineWidth := 0.5
	tableFontSize := 8.5
	table := form.pdf.NewTableLayout(tablePosition.x, tablePosition.y, 30, 8)

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
		FontSize:  tableFontSize,
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
		FontSize:  tableFontSize,
	})

	// Add columns to the table
	table.AddColumn("Address", 70, "center")
	table.AddColumn("Ownership on Title", 55, "center")
	table.AddColumn("Property Type", 64, "center")
	table.AddColumn("Purchase Price", 60, "center")
	table.AddColumn("Purchase Date", 60, "center")
	table.AddColumn("Tenant Occupied\n Y/N", 65, "center")
	table.AddColumn("Sale Price $", 60, "center")
	table.AddColumn("Date Sold", 60, "center")
	table.AddColumn("Profit / Loss $", 60, "center")

	table.AddRow([]string{"2625 Augustine Drive Santa Clara CA", "Me", "SFR Attached", "234555", " 04/03/2025", "Y", "234556", " 04/17/2025", " 12345"})
	table.AddRow([]string{"2625 Augustine Drive Santa Clara CA", "Me", "SFR Attached", "234555", " 04/03/2025", "Y", "234556", " 04/17/2025", " 12345"})
	table.AddRow([]string{"2625 Augustine Drive Santa Clara CA", "Me", "SFR Attached", "234555", " 04/03/2025", "N", "234556", " 04/17/2025", " 12345"})
	table.AddRow([]string{"2625 Augustine Drive Santa Clara CA", "Me", "SFR Attached", "234555", " 04/03/2025", "N", "234556", " 04/17/2025", " 12345"})
	err = table.DrawTable()

	return err
}
