package model

import (
	"zeebee-backend/common/constants"
	"zeebee-backend/common/types"
)

type Applicant struct {
	Base
	BorrowerId       int64                               `gorm:"column:borrower_id;not null;default:0" json:"borrower_id"`
	LoanId           int64                               `gorm:"column:loan_id;not null;default:0" json:"loan_id"`
	Name             string                              `gorm:"column:name;size:30;not null;default:''" json:"name"`
	IsPrimary        bool                                `gorm:"column:is_primary;not null;default:false" json:"is_primary"`
	LoanPropertyId   int64                               `gorm:"column:loan_property_id;not null;default:0" json:"loan_property_id"`
	Type             constants.ApplicantType             `gorm:"column:type;size:10;not null;default:'Individual';commit:'Individual,Corporation'" json:"type"`
	InvitationStatus constants.ApplicantInvitationStatus `gorm:"column:invitation_status;size:10;not null;default:'Pending';commit:'Accepted,Pending,Rejected'" json:"invitation_status"`
	Ssn              string                              `gorm:"column:ssn;size:20;not null;default:''" json:"ssn"`
	Dob              string                              `gorm:"column:dob;size:10;not null;default:''" json:"dob"`
	MaritalStatus    constants.MaritalStatus             `gorm:"column:marital_status;size:10;not null;default:''" json:"marital_status"`
	Email            string                              `gorm:"column:email;size:100;not null;default:''" json:"email"`
	EstCreditScore   int32                               `gorm:"column:est_credit_score;not null;default:0" json:"est_credit_score"`
	DriveLicenseNum  string                              `gorm:"column:drive_license_num;size:20;not null;default:''" json:"drive_license_num"`
	Address          types.Address                       `gorm:"column:address;type:jsonb;serializer:json" json:"address"`
	LoanBorrowerId   int64                               `gorm:"column:loan_borrower_id;not null;default:0" json:"loan_borrower_id"`
	// relationships
	LoanProperty *LoanProperty `json:"loan_property" gorm:"foreignKey:LoanPropertyId;references:Id"`
	Borrower     *Borrower     `json:"borrower" gorm:"foreignKey:BorrowerId;references:Id"`
	Loan         *Loan         `json:"loan" gorm:"foreignKey:LoanId;references:Id"`
	Applicants   []*Applicant  `json:"applicants" gorm:"foreignKey:LoanBorrowerId;references:Id"`
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
