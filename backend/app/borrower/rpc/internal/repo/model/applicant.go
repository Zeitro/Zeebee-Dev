package model

import "zeebee-backend/common/constants"

type Applicant struct {
	Base
	BorrowerId       int64                               `gorm:"column:borrower_id;not null;default:0" json:"borrower_id"`
	LoanId           int64                               `gorm:"column:loan_id;not null;default:0" json:"loan_id"`
	Name             string                              `gorm:"column:name;size:30;not null;default:''" json:"name"`
	IsPrimary        bool                                `gorm:"column:is_primary;not null;default:false" json:"is_primary"`
	LoanPropertyId   int64                               `gorm:"column:loan_property_id;not null;default:0" json:"loan_property_id"`
	Type             constants.ApplicantType             `gorm:"column:type;size:10;not null;default:'Individual';commit:'Individual,Corporation'" json:"type"`
	InvitationStatus constants.ApplicantInvitationStatus `gorm:"column:invitation_status;size:10;not null;default:'Pending';commit:'Accepted,Pending,Rejected'" json:"invitation_status"`

	LoanProperty *LoanProperty `json:"loan_property" gorm:"foreignKey:LoanPropertyId;references:Id"`
	Borrower     *Borrower     `json:"borrower" gorm:"foreignKey:BorrowerId;references:Id"`
	Loan         *Loan         `json:"loan" gorm:"foreignKey:LoanId;references:Id"`
}

func (m Applicant) TableName() string {
	return "applicants"
}

func (m Applicant) FieldBorrowerId() string {
	return "borrower_id"
}
func (m Applicant) FieldLoanId() string {
	return "loan_id"
}
