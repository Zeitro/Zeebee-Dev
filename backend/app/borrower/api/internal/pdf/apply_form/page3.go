package apply_form

import (
	"zeebee-backend/app/borrower/api/internal/pdf"
)

func (form *ApplyForm) page3() (err error) {
	form.pdf.AddPage()
	err = form.setPageTitle(pdf.NoteIco, "BORROWER", "INFORMATION")
	if err != nil {
		return err
	}
	form.pdf.Br(10)
	x := form.pdf.GetX()
	y := form.pdf.GetY()
	_ = form.buildTitleWithBg("DISCLOSURE TO APPLICANT REGARDING PROCUREMENT OF A BACKGROUND REPORT", position{
		x: x,
		y: y,
		w: form.pdf.GetFullWidth(),
		h: 27,
	}, position{
		x: x + 10,
		y: y,
		w: form.pdf.GetFullWidth() - 10,
		h: 27,
	})

	form.pdf.BrLine()
	form.pdf.BrLine()
	form.pdf.BrLine()

	text := []string{
		"In connection with your application for a loan, we may procure a background investigative report and/or credit report as part of the process of considering your candidacy as a borrower.",
		"With your signature below, you hereby authorize us, our affiliates, investors, other lending partners, and their agents and representatives to obtain a background and credit check to determine your financial responsibility, background and general credit worthiness. The Information requested in this application shall be used to underwrite a loan transaction. The background report may include, but not be limited to, criminal history, verifications of employment, credit report, education, and driving records.",
	}
	_ = form.pdf.MultiLongText(form.pdf.LineHeight, text)
	form.pdf.BrLine()
	form.pdf.BrLine()
	form.pdf.BrLine()
	err = form.buildSignature(form.info.Signatures.BorrowerInfo)
	if err != nil {
		return err
	}
	form.pdf.SetRegularFont()

	return err
}
