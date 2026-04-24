package apply_form

import (
	"zeebee-backend/app/borrower/api/internal/pdf"
)

func (form *ApplyForm) page12() (err error) {
	form.pdf.AddPage()
	err = form.setPageTitle(pdf.ListIco, "CONSTRUCTIVE RESERVE /", "FUND CONTROL")
	if err != nil {
		return err
	}
	form.pdf.BrLine()
	form.pdf.SetBlackRegularFont()
	text := []string{
		"Construction reserve funds, if any, will be withheld from the initial funding of the loan. The Borrower will be required to provide a detailed scope of work for Lender’s review prior to origination of the loan. This scope of work will be included in the loan documents. Borrower understands, acknowledges and agrees that Lender will not disburse funds in excess of the construction reserve amount. If excess funds beyond those in the construction reserve are required to complete work items included in the scope of work, Borrower shall be solely responsible for supplying such additional funds. Borrower is also solely responsible for supplying funds for unexpected expenses or costs incurred that are not included in the scope of work. ",
	}
	_ = form.pdf.MultiLongText(form.pdf.LineHeight, text)

	x := form.pdf.GetX()
	y := form.pdf.GetY()
	contents := []colorText{
		{
			text:  "Work item(s) must be completed prior to disbursal of corresponding draw amount. ",
			color: struct{ r, g, b uint8 }{0, 0, 0},
		},
		{
			text:  "For avoidance of doubt, work items shall consist only of actual costs toward improvements to the property and shall not include any overhead of Borrower, Lender’s interest, property taxes, property insurance, HOA dues and other costs which do not constitute improvements. ",
			color: struct{ r, g, b uint8 }{71, 117, 255},
		},
		{
			text:  "Upon each draw request, an inspection will be ordered with a 3rd party inspection company (Trinity or CFSI) and an inspector will contact you to confirm and schedule the site visit. The inspector will go the property with your budget and verify completed work and take photos. We will then review the inspection report and release funds based on the completed work and overall percentage of completion less the inspection fee. Each inspection has a fee of $200 or $225 which will be deducted from the approved funds (Larger projects consisting of multiple units may have a higher cost & will be confirmed prior to the first inspection.) Typical turn-around time is 5-7 days from the date of request to the date of wire. The minimum amount of each draw request shall be equal to the greater of $5,000 or 10% of the total construction reserve.",
			color: struct{ r, g, b uint8 }{0, 0, 0},
		},
	}

	_ = form.buildLongTextWithColor(contents, x, form.pdf.LineHeight)

	form.pdf.BrLine()

	text = []string{
		"Borrower understands, acknowledges and agrees that each draw disbursed by Lender shall be a percent, equal to the Lender- approved loan-to-cost ratio (construction reserve funds/Lender-approved budget for the scope of work), of the total amount of funds actually spent by Borrower on the work items and evidenced in the draw request, and approved by Lender in its sole and absolute discretion.",
		"Borrower understands, acknowledges and agrees that Borrower must have the necessary funds from Borrower’s own resources to fund the difference in actual costs in excess of loan funds advanced, if any, by Lender toward the work items. This includes paying for work items in advance of Borrower’s receipt of loan funds advanced, if any, by Lender toward the work items. ",
	}
	_ = form.pdf.MultiLongText(form.pdf.LineHeight, text)

	x = form.pdf.GetX()
	y = form.pdf.GetY()
	_ = form.buildSmallTitle("Need to know", position{
		x: x,
		y: y,
		w: form.pdf.GetFullWidth(),
		h: form.pdf.LineHeight,
	})
	form.pdf.Br(form.pdf.LineHeight * 1.5)
	x = form.pdf.GetX()
	y = form.pdf.GetY()

	items := []string{
		"Prior to the 1st draw, the Deed of Trust must be recorded",
		"Soft costs require proof of payment for reimbursement",
		"Funds throughout the project are reimbursed based on physical completed work verified by an inspector with corresponding invoices",
		"Draws in any amount over 20k require a lien release from the contractor (if applicable) or invoices/receipts when there is no contractor involved",
		"Costs that are over budget or not included in the budget will not be reimbursed.",
		"Materials are not reimbursed until installed & verified by the inspection company",
		"No draws will be disbursed if the loan is delinquent",
		"Final inspections are required even if the property is listed.",
		"We do not allow subordinated financing behind our loan without authorization",
	}
	startPos := position{
		x: x,
		y: y,
		w: form.pdf.GetFullWidth(),
	}
	_ = form.buildLists(items, startPos, form.pdf.LineHeight)

	return err
}
