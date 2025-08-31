package apply_form

import (
	"github.com/signintech/gopdf"

	"zeebee-backend/app/borrower/api/internal/pdf"
)

func (form *ApplyForm) page11() (err error) {
	form.pdf.AddPage()
	err = form.setPageTitle(pdf.NoteIco, "DISCLOSURES AND", "CONSENTS")
	if err != nil {
		return err
	}
	form.pdf.Br(form.pdf.LineHeight)
	x := form.pdf.GetX()
	y := form.pdf.GetY()
	_ = form.buildSmallTitle("ECOA – Equal Credit Opportunity Act Disclosures", position{
		x: x,
		y: y,
		w: form.pdf.GetFullWidth(),
		h: form.pdf.LineHeight,
	})
	form.pdf.Br(form.pdf.LineHeight * 1.5)
	form.setBlueFont()
	form.pdf.SetBoldFont()
	_ = form.pdf.Cell(nil, "ECOA Appraisal Report Disclosure ")
	form.pdf.Br(form.pdf.LineHeight)
	form.pdf.SetBlackRegularFont()

	text := []string{
		"We may order an appraisal to determine the property’s value and charge you for this appraisal. We will promptly give you a copy of the appraisal, even if your loan does not close. You can pay for an additional appraisal for your own use, at your own cost.",
	}
	_ = form.pdf.MultiLongText(form.pdf.LineHeight, text)
	// form.pdf.Br(form.pdf.LineHeight)

	x = form.pdf.GetX()
	y = form.pdf.GetY()
	contents := []colorText{
		{
			text:  "Equal Credit Opportunity Act Disclosure of Right to Request Specific Reasons for Credit Denial",
			color: struct{ r, g, b uint8 }{71, 117, 255},
		},
		{
			text:  " If your application for business credit is denied, you have the right to a written statement of the specific reasons for the denial. To obtain the statement, please contact: Credit Department at 949-244-1090, within 60 days from the date you are notified of our decision. We will send you a written statement of reasons for the denial within 30 days of receiving your request for the statement. Notice: The Federal Equal Credit Opportunity Act prohibits creditors from discriminating against credit applicants on the basis of race, color, religion, national origin, sex, marital status, age (provided the applicant has the capacity to enter into a binding contact); because all or part of the applicant’s income derives from any public assistance program; or because the applicant has a good faith exercised any right under the Consumer Credit Protection Act. The Federal agency that administers compliance with this law concerning this creditor is:",
			color: struct{ r, g, b uint8 }{0, 0, 0},
		},
	}
	_ = form.buildLongTextWithColor(contents, x, form.pdf.LineHeight)
	form.pdf.Br(form.pdf.LineHeight)

	form.setBlueFont()
	form.pdf.SetBoldFont()
	_ = form.pdf.Cell(nil, "Equal Credit Opportunity Act Information")
	form.pdf.Br(form.pdf.LineHeight)
	form.pdf.SetBlackRegularFont()
	x = form.pdf.GetX()
	y = form.pdf.GetY()
	contents = []colorText{
		{
			text:  "Fair Lending Notice",
			color: struct{ r, g, b uint8 }{71, 117, 255},
		},
		{
			text:  " The Federal Equal Credit Opportunity Act prohibits creditors from discriminating against credit applicants on the basis of race, color, religion, national origin, sex, marital status, age (provided the applicant has the capacity to enter into a binding contact); because all or part of the applicant’s income derives from any public assistance program; or because the applicant has a good faith exercised any right under the Consumer Credit Protection Act. The Federal agency that administers compliance with this law concerning this creditor is: ",
			color: struct{ r, g, b uint8 }{0, 0, 0},
		},
	}
	_ = form.buildLongTextWithColor(contents, x, form.pdf.LineHeight)
	form.pdf.Br(form.pdf.LineHeight)
	cellOpt := gopdf.CellOption{
		Align: gopdf.Left | gopdf.Middle,
	}
	_ = form.pdf.MultiCellWithOption(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: form.pdf.LineHeight * 2}, "Federal Trade Commission", cellOpt)
	form.pdf.Br(form.pdf.LineHeight * 0.5)
	_ = form.pdf.MultiCellWithOption(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: form.pdf.LineHeight * 2}, "Equal Credit Opportunity", cellOpt)
	form.pdf.Br(form.pdf.LineHeight * 0.5)
	_ = form.pdf.MultiCellWithOption(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: form.pdf.LineHeight * 2}, "Washington, DC 20580", cellOpt)

	form.pdf.Br(form.pdf.LineHeight)
	x = form.pdf.GetX()
	y = form.pdf.GetY()
	_ = form.buildSmallTitle("Patriot Act.", position{
		x: x,
		y: y,
		w: form.pdf.GetFullWidth(),
		h: form.pdf.LineHeight,
	})
	form.pdf.Br(form.pdf.LineHeight * 1.5)
	form.setBlueFont()
	form.pdf.SetBoldFont()
	cellRect := &gopdf.Rect{W: form.pdf.GetFullWidth(), H: form.pdf.LineHeight}
	cellOpt = gopdf.CellOption{
		Align: gopdf.Center | gopdf.Middle,
	}
	_ = form.pdf.CellWithOption(cellRect, "Patriot Act Information Disclosure Customer Identification", cellOpt)
	form.pdf.Br(form.pdf.LineHeight)
	_ = form.pdf.CellWithOption(cellRect, "Program Important Information About Application Procedures", cellOpt)
	form.pdf.Br(form.pdf.LineHeight)

	form.pdf.SetBlackRegularFont()
	text = []string{
		"To help the US Government fight the funding of terrorism and money laundering activities, Federal law requires all financial institutions to obtain, verify and record information that identified every person who opens an account. What this means to you:",
		"When you apply for a loan, we will ask for your name, address, date of birth, and other information that will allow us to identify you. We may also ask to see your driver’s license or other identifying documents.",
	}

	err = form.pdf.MultiLongText(form.pdf.LineHeight, text)

	// text = "To help the US Government fight the funding of terrorism and money laundering activities, Federal law requires all financial institutions to obtain, verify and record information that identified every person who opens an account. What this means to you: "
	// _ = form.pdf.LongText(form.pdf.LineHeight, text)
	// form.pdf.Br(form.pdf.LineHeight)
	// text = "When you apply for a loan, we will ask for your name, address, date of birth, and other information that will allow us to identify you. We may also ask to see your driver’s license or other identifying documents.  "
	// _ = form.pdf.LongText(form.pdf.LineHeight, text)
	// form.pdf.Br(form.pdf.LineHeight)

	_ = form.buildSignature()

	return err
}
