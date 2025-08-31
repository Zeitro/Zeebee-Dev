package apply_form

import (
	"zeebee-backend/app/borrower/api/internal/pdf"
)

func (form *ApplyForm) page10() (err error) {
	form.pdf.AddPage()
	err = form.setPageTitle(pdf.NoteIco, "DISCLOSURES AND", "CONSENTS")
	if err != nil {
		return err
	}
	form.pdf.Br(form.pdf.LineHeight)
	contents := []string{
		"The information contained in this application is provided for the purpose of obtaining or maintaining credit with Lender, on behalf of the undersigned, or persons, firms or corporations on whose behalf the undersigned may either severally or jointly with others execute a guaranty in Lender’s favor. Each undersigned understands that Lender is relying on the information provided herein (including the designation made as to the ownership of property) in deciding to grant or continue credit. Each undersigned represents and warrants that the information provided is true and complete and that Lender may consider this statement as continuing to be true and correct until a written notice of a change is given to Lender by the undersigned.",
		"Furthermore, each of the undersigned specifically represents to Lender and to Lender’s actual or potential agents, brokers, processors, attorneys, insurers, servicers, successors and assigns and agrees and acknowledges that: (1) the information provided in this application is true and correct as of the date set forth opposite my/our signature(s) and that any intentional or negligent misrepresentation of this information contained in this application may result in civil liability, including monetary damages, to any person who may suffer any loss due to reliance upon any misrepresentation that I/we have made on this application, and/or in criminal penalties including, but not limited to, fine or imprisonment or both under the provisions of Title 18, United States Code, Sec. 1001, et seq.; (2) the loan requested pursuant to this application (the “Loan”) will be secured by a mortgage or deed of trust on the property described in this application; (3) the property will not be used for any illegal or prohibited purpose or use; (4) all statements made in this application are made for the purpose of obtaining a loan on a residential property for business purposes only; (5) the property will be NOT be occupied by borrower, sponsors or their families as has been indicated in this application; (6) the Lender, its servicers, successors or assigns may retain the original and/or an electronic record of this application, whether or not the Loan is approved; (7) the Lender and its agents, brokers, insurers, servicers, successors, and assigns may continuously rely on the information contained in the application, and I/we are obligated to amend and/or supplement the information provided in this application if any of the material facts that I/we have represented herein should change prior to closing of the Loan; (8) in the event that my/our payments on the Loan become delinquent, the Lender, its servicers, successors or assigns may, in addition to any other rights and remedies that it may have relating to such delinquency, report my/our name and account information to one or more reporting agencies; (9) ownership of the Loan and/or administration of the Loan account may be transferred with such notice as may be required by law; and (10) neither Lender nor its agents, brokers, insurers, servicers, successors or assigns has made any representation or warranty, express or implied, to me regarding the property or the condition or value of the property.",
	}
	err = form.pdf.MultiLongText(form.pdf.LineHeight, contents)
	if err != nil {
		return err
	}
	x := form.pdf.GetX()
	y := form.pdf.GetY()
	_ = form.buildSmallTitle(" Business Purpose & Appraisal Acknowledgment", position{
		x: x,
		y: y,
		w: form.pdf.GetFullWidth(),
		h: form.pdf.LineHeight,
	})
	form.pdf.Br(form.pdf.LineHeight * 1.5)
	contents = []string{
		"Acknowledgment: Each of the undersigned hereby acknowledges that any owner of the Loan, its servicers, successors and assigns, may verify or re-verify any information contained in this application or obtain any information or data relating to the Loan, for any legitimate business purpose through any source, including a source named in this application or a credit reporting agency.",
		"You are entitled to receive a copy of appraisals and other written valuations developed in connection with your application for credit secured by a first lien on a dwelling promptly upon completion, or three days prior to closing, whichever is earlier. You may waive the timing requirement and agree to receive any copy at or before closing.",
		"By signing below, you hereby acknowledge timely receipt of this notice, waive the timing requirement, and agree to receive such appraisal/valuation at closing.",
		"Please confirm that you understand this is a commercial loan for business purposes only and you do not intend to occupy any property financed using proceeds from this loan as your primary residence.",
	}
	err = form.pdf.MultiLongText(form.pdf.LineHeight, contents)
	if err != nil {
		return err
	}

	_ = form.buildSignature()

	return err
}
