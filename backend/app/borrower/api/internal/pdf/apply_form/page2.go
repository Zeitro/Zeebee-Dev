package apply_form

import (
	"fmt"

	"github.com/signintech/gopdf"

	"zeebee-backend/app/borrower/api/internal/pdf"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/common/constants"
	"zeebee-backend/common/stringx"
)

func (form *ApplyForm) page2() (err error) {
	form.pdf.AddPage()
	err = form.setPageTitle(pdf.NoteIco, "BORROWER", "INFORMATION")
	if err != nil {
		return err
	}
	form.pdf.SetBoldFont()
	_ = form.pdf.Cell(nil, "Borrower Type:")
	// form.pdf.Br(20)
	form.pdf.BrLine()
	x := form.pdf.GetX()
	y := form.pdf.GetY()
	startX := x + 50.0
	startY := y + 5
	borrowerInfo := form.info.LoanBorrower

	err = form.buildRadioCheckbox(startX, startY, constants.ApplicantTypeName, borrowerInfo.Type)
	if err != nil {
		return err
	}
	form.pdf.BrLine()
	// 绘制灰色线条
	form.pdf.SetGrayStroke()
	form.pdf.Line(form.pdf.MarginLeft(), form.pdf.GetY(), gopdf.PageSizeA4.W-form.pdf.MarginRight(), form.pdf.GetY())
	form.pdf.SetRegularFont()
	rowHeight := 12.0
	brHeight := 6.0
	form.pdf.Br(brHeight)
	_ = form.pdf.Cell(&gopdf.Rect{W: 250, H: rowHeight}, "Entity Name/Individual:  "+borrowerInfo.Name)
	_ = form.pdf.Cell(&gopdf.Rect{W: 150, H: rowHeight}, "TIN:  "+borrowerInfo.Tin)
	_ = form.pdf.Cell(&gopdf.Rect{W: 100, H: rowHeight}, "% Own:  "+fmt.Sprintf("%d", borrowerInfo.OwnPercent))
	form.pdf.Br(brHeight * 3.0)
	_ = form.pdf.MultiCell(&gopdf.Rect{W: gopdf.PageSizeA4.W - form.pdf.MarginLeft() - form.pdf.MarginRight(), H: rowHeight}, "Mailing Address:  "+borrowerInfo.MailingAddress.Address1)
	form.pdf.Br(brHeight)
	_ = form.pdf.Cell(&gopdf.Rect{W: 250, H: rowHeight}, "City:  "+borrowerInfo.MailingAddress.City)
	_ = form.pdf.Cell(&gopdf.Rect{W: 150, H: rowHeight}, "State:  "+borrowerInfo.MailingAddress.State)
	_ = form.pdf.Cell(&gopdf.Rect{W: 100, H: rowHeight}, "Zip:  "+borrowerInfo.MailingAddress.ZipCode)
	form.pdf.Br(brHeight * 3.0)
	_ = form.pdf.MultiCell(&gopdf.Rect{W: gopdf.PageSizeA4.W - form.pdf.MarginLeft() - form.pdf.MarginRight(), H: rowHeight}, "Phone:  "+borrowerInfo.ContactPhone)
	form.pdf.Br(brHeight)
	_ = form.pdf.MultiCell(&gopdf.Rect{W: gopdf.PageSizeA4.W - form.pdf.MarginLeft() - form.pdf.MarginRight(), H: rowHeight}, "Email:  "+borrowerInfo.ContactEmail)
	form.pdf.Br(brHeight)
	_ = form.pdf.MultiCell(&gopdf.Rect{W: gopdf.PageSizeA4.W - form.pdf.MarginLeft() - form.pdf.MarginRight(), H: rowHeight}, "Contact Name:  "+borrowerInfo.ContactName)
	form.pdf.Br(brHeight)
	form.pdf.Line(form.pdf.MarginLeft(), form.pdf.GetY(), gopdf.PageSizeA4.W-form.pdf.MarginRight(), form.pdf.GetY())

	applicantsInfo := form.info.Applicants
	for i := 0; i < 2; i++ {
		var applicant *types.LoanApplyApplicantEntry
		if i > len(applicantsInfo)-1 {
			applicant = &types.LoanApplyApplicantEntry{
				Address: &types.AddressEntry{},
			}
		} else {
			applicant = applicantsInfo[i]
		}
		err = form.applicantInfo(i+1, applicant)
		if err != nil {
			return err
		}

	}

	return err
}

func (form *ApplyForm) applicantInfo(index int, applicant *types.LoanApplyApplicantEntry) (err error) {
	brHeight := 6.0
	x := form.pdf.GetX()
	y := form.pdf.GetY()
	_ = form.pdf.ImageFrom(pdf.GrayBg, x, y, &gopdf.Rect{W: form.pdf.GetFullWidth(), H: 27})
	_ = form.pdf.ImageFrom(pdf.BlueBg, x+10, y, &gopdf.Rect{W: form.pdf.GetFullWidth() - 10, H: 27})

	form.pdf.SetWhiteBoldFont()
	form.pdf.SetXY(x+20, y+8)
	_ = form.pdf.Cell(nil, fmt.Sprintf("Applicant %d:", index))
	form.pdf.SetBlackRegularFont()
	form.pdf.Br(30)
	rowHeight := 12.0
	_ = form.pdf.MultiCell(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: rowHeight}, "First Name:  "+applicant.FirstName)
	form.pdf.Br(brHeight)
	_ = form.pdf.MultiCell(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: rowHeight}, "Middle Name:  "+applicant.MidName)
	form.pdf.Br(brHeight)
	_ = form.pdf.MultiCell(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: rowHeight}, "Last Name:  "+applicant.LastName)
	form.pdf.Br(brHeight)
	_ = form.pdf.MultiCell(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: rowHeight}, "SSN:  "+applicant.Ssn)
	form.pdf.Br(brHeight)
	_ = form.pdf.MultiCell(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: rowHeight}, "DOB:  "+applicant.Dob)
	form.pdf.Br(brHeight)
	_ = form.pdf.MultiCell(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: rowHeight}, "Phone:  "+applicant.Phone)
	form.pdf.Br(brHeight)
	_ = form.pdf.MultiCell(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: rowHeight}, "Marital Status:")
	x = form.pdf.GetX()
	y = form.pdf.GetY()
	startX := x + 50.0
	startY := y + 5

	err = form.buildRadioCheckbox(startX, startY, constants.MaritalStatusName, applicant.MaritalStatus)
	if err != nil {
		return err
	}

	form.pdf.Br(brHeight)
	_ = form.pdf.MultiCell(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: rowHeight}, "Email:  "+applicant.Email)
	form.pdf.Br(brHeight)
	_ = form.pdf.MultiCell(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: rowHeight}, "Est. Credit Score:  "+stringx.TransNumberToString[int32](applicant.EstCreditScore, true))
	form.pdf.Br(brHeight)
	_ = form.pdf.MultiCell(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: rowHeight}, "Primary Address:  "+applicant.Address.Address1)
	form.pdf.Br(brHeight)
	_ = form.pdf.Cell(&gopdf.Rect{W: 100, H: rowHeight}, "City:  "+applicant.Address.City)
	_ = form.pdf.Cell(&gopdf.Rect{W: (form.pdf.GetFullWidth()) / 2.0, H: rowHeight}, "State:  "+applicant.Address.State)
	form.pdf.Br(brHeight)
	form.pdf.Br(15)
	_ = form.pdf.MultiCell(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: rowHeight}, "Zip:  "+applicant.Address.ZipCode)
	form.pdf.Br(brHeight)
	_ = form.pdf.MultiCell(&gopdf.Rect{W: form.pdf.GetFullWidth(), H: rowHeight}, "Drivers License #:  "+applicant.DriverLicenseNum)
	form.pdf.Br(brHeight)
	return err
}
