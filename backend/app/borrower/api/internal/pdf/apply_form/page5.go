package apply_form

import (
	"github.com/signintech/gopdf"

	"zeebee-backend/app/borrower/api/internal/pdf"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/common/stringx"
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
	budget := form.info.LoanProperty.Budget
	if budget == nil {
		budget = &types.LoanBudgetEntry{}
	}
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
	_ = form.pdf.Cell(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: fontHeight}, "Original Budget:  "+stringx.FormatUsDollarCurrency(budget.OriginalBudget))
	form.pdf.Br(lineHeight)
	_ = form.pdf.Cell(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: fontHeight}, "Remaining Budget:  "+stringx.FormatUsDollarCurrency(budget.RemainingBudget))
	form.pdf.Br(lineHeight)
	_ = form.pdf.Cell(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: fontHeight}, "Cost Spent to Date:  "+stringx.FormatUsDollarCurrency(budget.CostSpentToDateAmount))

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
	err = form.buildPage5ThirdServiceTable(position{x: x, y: y})
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
	property := form.info.LoanProperty
	table.AddColumn("", 100, "center")
	table.AddColumn("Current Status", 150, "center")
	table.AddColumn("Future Status", 150, "center")

	table.AddRow([]string{"Property Type", property.PropertyType, property.FuturePropertyType})
	table.AddRow([]string{"Living Area SQFT", stringx.TransNumberToString(property.LivingAreaSqft, true), stringx.TransNumberToString(property.FutureLivingAreaSqft, true)})
	table.AddRow([]string{"Number of Units", stringx.TransNumberToString(property.Unit, true), stringx.TransNumberToString(property.FutureUnit, true)})
	table.AddRow([]string{"Number of Bedrooms", stringx.TransNumberToString(property.Bedrooms, true), stringx.TransNumberToString(property.FutureBedrooms, true)})
	table.AddRow([]string{"Number of Bathrooms", stringx.TransNumberToString(property.Bathrooms, true), stringx.TransNumberToString(property.FutureBathrooms, true)})
	err = table.DrawTable()

	return err
}

func (form *ApplyForm) buildPage5ThirdServiceTable(tablePosition position) (err error) {
	lineWidth := 0.5
	thirdServices := form.info.LoanThirdServices
	maxRows := len(thirdServices)
	if maxRows <= 5 {
		maxRows = 5
	}
	table := form.pdf.NewTableLayout(tablePosition.x, tablePosition.y, 30, maxRows)

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
	table.AddColumn("Service Type", 109, "center")
	table.AddColumn("Name", 109, "center")
	table.AddColumn("Company", 109, "center")
	table.AddColumn("Phone", 109, "center")
	table.AddColumn("Email", 109, "center")

	for _, thirdService := range thirdServices {
		table.AddRow([]string{thirdService.Type, thirdService.Name, thirdService.Company, stringx.FormatUsPhoneNumber(thirdService.Phone), thirdService.Email})
	}
	err = table.DrawTable()

	return err
}
