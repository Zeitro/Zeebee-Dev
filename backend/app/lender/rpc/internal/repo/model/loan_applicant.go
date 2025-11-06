package model

import (
	"database/sql"

	"zeebee-backend/common/constants"
	"zeebee-backend/common/types"
)

type LoanApplicant struct {
	Base
	BorrowerId          int64                               `gorm:"column:borrower_id;not null;default:0" json:"borrower_id"`
	BorrowerApplicantId int64                               `gorm:"column:borrower_applicant_id;not null;default:0" json:"borrower_applicant_id"`
	LoanId              int64                               `gorm:"column:loan_id;not null;default:0" json:"loan_id"`
	FullName            string                              `gorm:"column:full_name;size:30;not null;default:''" json:"full_name"`
	FirstName           string                              `gorm:"column:first_name;size:30;not null;default:''" json:"first_name"`
	MidName             string                              `gorm:"column:mid_name;size:30;not null;default:''" json:"mid_name"`
	LastName            string                              `gorm:"column:last_name;size:30;not null;default:''" json:"last_name"`
	IsPrimary           bool                                `gorm:"column:is_primary;not null;default:false" json:"is_primary"`
	LoanPropertyId      int64                               `gorm:"column:loan_property_id;not null;default:0" json:"loan_property_id"`
	Type                constants.ApplicantType             `gorm:"column:type;size:10;not null;default:'Individual';commit:'Individual,Corporation'" json:"type"`
	InvitationStatus    constants.ApplicantInvitationStatus `gorm:"column:invitation_status;size:10;not null;default:'Pending';commit:'Accepted,Pending,Rejected'" json:"invitation_status"`
	Ssn                 string                              `gorm:"column:ssn;size:20;not null;default:''" json:"ssn"`
	Dob                 *sql.NullTime                       `gorm:"column:dob;size:10;not null;default:''" json:"dob"`
	MaritalStatus       constants.MaritalStatus             `gorm:"column:marital_status;size:10;not null;default:''" json:"marital_status"`
	Email               string                              `gorm:"column:email;size:100;not null;default:''" json:"email"`
	EstCreditScore      int32                               `gorm:"column:est_credit_score;not null;default:0" json:"est_credit_score"`
	DriverLicenseNum    string                              `gorm:"column:driver_license_num;size:20;not null;default:''" json:"driver_license_num"`
	Address             types.Address                       `gorm:"column:address;type:jsonb;serializer:json" json:"address"`
	LoanBorrowerId      int64                               `gorm:"column:loan_borrower_id;not null;default:0" json:"loan_borrower_id"`
	Phone               string                              `gorm:"column:phone;size:20;not null;default:''" json:"phone"`

	// relationships
	LoanProperty *LoanProperty `json:"loan_property" gorm:"foreignKey:LoanPropertyId;references:Id"`
	Borrower     *Borrower     `json:"borrower" gorm:"foreignKey:BorrowerId;references:Id"`
	Loan         *Loan         `json:"loan" gorm:"foreignKey:LoanId;references:Id"`
}

func (m LoanApplicant) TableName() string {
	return "applicants"
}

func (m *LoanApplicant) RemoveRelationships() {
	m.LoanProperty = nil
	m.Borrower = nil
	m.Loan = nil
}
