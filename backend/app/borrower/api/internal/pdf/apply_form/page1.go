package apply_form

import (
	"github.com/signintech/gopdf"

	"zeebee-backend/app/borrower/api/internal/pdf"
)

func (form *ApplyForm) page1() (err error) {
	form.pdf.AddPage()
	form.pdf.SetRegularFont()

	x := form.pdf.GetX()
	y := form.pdf.GetY()
	// gray bg
	err = form.pdf.ImageFrom(pdf.GrayBg, x, y, &gopdf.Rect{W: gopdf.PageSizeA4.W, H: 80})
	if err != nil {
		return err
	}
	// blue bg
	err = form.pdf.ImageFrom(pdf.BlueBg, x, y+9, &gopdf.Rect{W: 395, H: 63})
	if err != nil {
		return err
	}

	err = form.pdf.ImageFrom(pdf.BookIco, 45, 27, &gopdf.Rect{W: 29, H: 29})
	if err != nil {
		return err
	}

	// white bg
	err = form.pdf.ImageFrom(pdf.WhiteBg, x+400, y+8, &gopdf.Rect{W: 200, H: 63})
	if err != nil {
		return err
	}

	// logo
	err = form.pdf.ImageFrom(pdf.Logo, x+430, y+20, &gopdf.Rect{W: 120, H: 40})
	if err != nil {
		return err
	}
	// application text
	form.pdf.SetXY(82, 35)
	form.pdf.SetBoldFontWithSize(20)
	form.pdf.SetWhiteFont()
	_ = form.pdf.Cell(nil, "APPLICATION ")
	form.pdf.SetBlackFont()
	_ = form.pdf.Cell(nil, "PACKET")

	form.pdf.SetXY(20, 91)
	form.pdf.SetMarginLeft(20)
	form.pdf.SetMarginRight(20)

	form.pdf.SetRegularFont()

	_ = form.pdf.LongText(15, "To initiate the loan process, please complete the application form and submit with applicable documents. Reach out to your loan officer or email pl@ratednagroup.com with any questions.")
	// content
	x = form.pdf.GetX()
	y = form.pdf.GetY()
	_ = form.pdf.ImageFrom(pdf.GrayBg, x, y, &gopdf.Rect{W: form.pdf.GetFullWidth(), H: 27})
	_ = form.pdf.ImageFrom(pdf.BlueBg, x+10, y, &gopdf.Rect{W: form.pdf.GetFullWidth() - 10, H: 27})
	form.pdf.SetWhiteBoldFont()

	form.pdf.SetXY(x+30, y+10)
	_ = form.pdf.Cell(nil, "CONTENT")
	form.pdf.SetXY(520, y+10)
	_ = form.pdf.Cell(nil, "PAGE")
	form.pdf.Br(30)

	var firstCellWidth float64 = 220
	var secondCellWidth float64 = 280
	var thirdCellWidth float64 = 100
	var rowHeight float64 = 30
	var applicantCompleteText = "All applicants must complete"
	form.pdf.SetBlackRegularFont()
	_ = form.pdf.Cell(&gopdf.Rect{W: firstCellWidth, H: rowHeight}, "Borrower Information")
	_ = form.pdf.Cell(&gopdf.Rect{W: secondCellWidth, H: rowHeight}, applicantCompleteText)
	_ = form.pdf.Cell(&gopdf.Rect{W: thirdCellWidth, H: rowHeight}, "2 - 3")
	form.pdf.Br(20)
	form.pdf.Line(form.pdf.GetPageStartX(), form.pdf.GetY(), form.pdf.GetPageEndX(), form.pdf.GetY())

	form.pdf.Br(20)
	x = form.pdf.GetX()
	y = form.pdf.GetY()
	form.pdf.SetXY(x, y)
	_ = form.pdf.Cell(&gopdf.Rect{W: firstCellWidth, H: rowHeight}, "Property Information/")
	form.pdf.SetXY(x, y+15)
	_ = form.pdf.Cell(&gopdf.Rect{W: firstCellWidth, H: rowHeight}, "Loan Information")
	form.pdf.SetXY(x+firstCellWidth, y)
	_ = form.pdf.Cell(&gopdf.Rect{W: secondCellWidth, H: rowHeight}, "Please complete for each applicable property")
	_ = form.pdf.Cell(&gopdf.Rect{W: thirdCellWidth, H: rowHeight}, "4 - 5")
	form.pdf.Br(40)
	form.pdf.Line(form.pdf.GetPageStartX(), form.pdf.GetY(), form.pdf.GetPageEndX(), form.pdf.GetY())

	form.pdf.Br(20)
	_ = form.pdf.Cell(&gopdf.Rect{W: firstCellWidth, H: rowHeight}, "Experience")
	x = form.pdf.GetX()
	y = form.pdf.GetY()
	_ = form.pdf.Cell(&gopdf.Rect{W: secondCellWidth, H: rowHeight}, "Please list 1-10 properties or provide REO")
	form.pdf.SetXY(x, y+15)
	_ = form.pdf.Cell(&gopdf.Rect{W: secondCellWidth, H: rowHeight}, "schedule of completed projects")
	form.pdf.SetXY(x+secondCellWidth, y)
	_ = form.pdf.Cell(&gopdf.Rect{W: thirdCellWidth, H: rowHeight}, "6")
	form.pdf.Br(40)
	form.pdf.Line(form.pdf.GetPageStartX(), form.pdf.GetY(), form.pdf.GetPageEndX(), form.pdf.GetY())

	form.pdf.Br(20)
	_ = form.pdf.Cell(&gopdf.Rect{W: firstCellWidth, H: rowHeight}, "Declarations/ Demographic Information")
	_ = form.pdf.Cell(&gopdf.Rect{W: secondCellWidth, H: rowHeight}, applicantCompleteText)
	_ = form.pdf.Cell(&gopdf.Rect{W: thirdCellWidth, H: rowHeight}, "7 - 9")
	form.pdf.Br(20)
	form.pdf.Line(form.pdf.GetPageStartX(), form.pdf.GetY(), form.pdf.GetPageEndX(), form.pdf.GetY())

	form.pdf.Br(20)
	_ = form.pdf.Cell(&gopdf.Rect{W: firstCellWidth, H: rowHeight}, "Disclosures and Consents")
	_ = form.pdf.Cell(&gopdf.Rect{W: secondCellWidth, H: rowHeight}, applicantCompleteText)
	_ = form.pdf.Cell(&gopdf.Rect{W: thirdCellWidth, H: rowHeight}, "10 - 11")
	form.pdf.Br(20)
	form.pdf.Line(form.pdf.GetPageStartX(), form.pdf.GetY(), form.pdf.GetPageEndX(), form.pdf.GetY())

	form.pdf.Br(20)
	_ = form.pdf.Cell(&gopdf.Rect{W: firstCellWidth, H: rowHeight}, "Construction Reserve/Fund Control")
	_ = form.pdf.Cell(&gopdf.Rect{W: secondCellWidth, H: rowHeight}, applicantCompleteText)
	_ = form.pdf.Cell(&gopdf.Rect{W: thirdCellWidth, H: rowHeight}, "12")
	form.pdf.Br(20)
	form.pdf.Line(form.pdf.GetPageStartX(), form.pdf.GetY(), form.pdf.GetPageEndX(), form.pdf.GetY())

	form.pdf.Br(20)
	_ = form.pdf.Cell(&gopdf.Rect{W: firstCellWidth, H: rowHeight}, "Borrower Authorization")
	_ = form.pdf.Cell(&gopdf.Rect{W: secondCellWidth, H: rowHeight}, applicantCompleteText)
	_ = form.pdf.Cell(&gopdf.Rect{W: thirdCellWidth, H: rowHeight}, "13")
	form.pdf.Br(20)
	form.pdf.Line(form.pdf.GetPageStartX(), form.pdf.GetY(), form.pdf.GetPageEndX(), form.pdf.GetY())

	form.pdf.BrLine()
	_ = form.pdf.LongText(15, "For loans being held in an Entity, all owners with 20% or greater ownership are subject to Credit and Background Review.  Applicant 1 will be designated the Personal Guarantor, unless a waiver is granted by Yeelee Capital.")

	form.pdf.BrLine()
	x = form.pdf.GetX()
	y = form.pdf.GetY()
	_ = form.pdf.ImageFrom(pdf.GrayBg, x, y, &gopdf.Rect{W: form.pdf.GetFullWidth(), H: 27})
	_ = form.pdf.ImageFrom(pdf.BlueBg, x+10, y, &gopdf.Rect{W: form.pdf.GetFullWidth() - 10, H: 27})
	form.pdf.SetWhiteBoldFont()
	form.pdf.SetXY(x+30, y+10)
	_ = form.pdf.Cell(nil, "DOCUMENTS")
	form.pdf.SetBlackRegularFont()
	form.pdf.Br(17)

	firstCellWidth = 150
	secondCellWidth = 350

	docs := []Page1Doc{
		{"All Loans", []string{
			"·  2 months recent bank statements (dated within 30 days of settlement date)",
			"·  Entity Documents",
			"  °  LLC: Articles of Organization and Operating Agreement",
			"  °  Corporations: Articles of Incorporation and Bylaws",
			"·  Valid, Legible Photo ID for all signatories",
		}},
		{"Refinance Loans", []string{
			"·  Current Mortgage Statement",
			"·  Itemized Scope of Work (Budget)",
			"·  Costs Spent to Date",
		}},
		{"Purchase Loans", []string{
			"·  Executed Purchase Contract with all addendums",
			"·  EItemized Scope of work (Budget)",
		}},
	}
	for _, doc := range docs {
		form.pdf.Br(20)
		_ = form.pdf.Cell(&gopdf.Rect{W: firstCellWidth, H: rowHeight}, doc.title)
		x = form.pdf.GetX()
		y = form.pdf.GetY()
		for idx, item := range doc.contents {
			form.pdf.SetXY(x, y+float64(15*idx))
			_ = form.pdf.Cell(&gopdf.Rect{W: secondCellWidth, H: rowHeight}, item)
		}
		form.pdf.Br(20)
		form.pdf.Line(form.pdf.MarginLeft(), form.pdf.GetY(), gopdf.PageSizeA4.W-form.pdf.MarginLeft()-form.pdf.MarginRight(), form.pdf.GetY())
	}

	return err
}

type Page1Doc struct {
	title    string
	contents []string
}
