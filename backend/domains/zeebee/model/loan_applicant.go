package model

import (
	"database/sql"

	"zeebee-backend/domains/zeebee/constants"
)

type LoanApplicant struct {
	Base
	BorrowerOrganizationId int64                               `gorm:"column:borrower_organization_id;not null;default:0" json:"borrower_organization_id"`
	BorrowerApplicantId    int64                               `gorm:"column:borrower_applicant_id;not null;default:0" json:"borrower_applicant_id"`
	LoanId                 int64                               `gorm:"column:loan_id;not null;default:0" json:"loan_id"`
	FullName               string                              `gorm:"column:full_name;size:30;not null;default:''" json:"full_name"`
	FirstName              string                              `gorm:"column:first_name;size:30;not null;default:''" json:"first_name"`
	MidName                string                              `gorm:"column:mid_name;size:30;not null;default:''" json:"mid_name"`
	LastName               string                              `gorm:"column:last_name;size:30;not null;default:''" json:"last_name"`
	IsPrimary              bool                                `gorm:"column:is_primary;not null;default:false" json:"is_primary"`
	LoanPropertyId         int64                               `gorm:"column:loan_property_id;not null;default:0" json:"loan_property_id"`
	InvitationStatus       constants.ApplicantInvitationStatus `gorm:"column:invitation_status;size:10;not null;default:'Pending';commit:'Accepted,Pending,Rejected'" json:"invitation_status"`
	Ssn                    string                              `gorm:"column:ssn;size:20;not null;default:''" json:"ssn"`
	Dob                    *sql.NullTime                       `gorm:"column:dob;size:10;not null;default:''" json:"dob"`
	MaritalStatus          constants.MaritalStatus             `gorm:"column:marital_status;size:10;not null;default:''" json:"marital_status"`
	Email                  string                              `gorm:"column:email;size:100;not null;default:''" json:"email"`
	EstCreditScore         int32                               `gorm:"column:est_credit_score;not null;default:0" json:"est_credit_score"`
	DriverLicenseNum       string                              `gorm:"column:driver_license_num;size:20;not null;default:''" json:"driver_license_num"`
	Address                Address                             `gorm:"column:address;type:jsonb;serializer:json" json:"address"`
	LoanBorrowerId         int64                               `gorm:"column:loan_borrower_id;not null;default:0" json:"loan_borrower_id"`
	Phone                  string                              `gorm:"column:phone;size:20;not null;default:''" json:"phone"`
	OwnPercent             int32                               `gorm:"column:own_percent;not null;default:0" json:"own_percent"`

	// relationships
	LoanProperty *LoanProperty `json:"loan_property" gorm:"foreignKey:LoanPropertyId;references:Id"`
	// User              *User              `json:"user" gorm:"foreignKey:UserId;references:Id"`
	BorrowerOrganization *Organization      `json:"borrower_organization" gorm:"foreignKey:BorrowerOrganizationId;references:Id"`
	Loan                 *Loan              `json:"loan" gorm:"foreignKey:LoanId;references:Id"`
	BorrowerApplicant    *BorrowerApplicant `json:"borrower_applicant" gorm:"foreignKey:BorrowerApplicantId;references:Id"`
}

func (m LoanApplicant) TableName() string {
	return "loan_applicants"
}

func (m *LoanApplicant) RemoveRelationships() {
	m.LoanProperty = nil
	m.BorrowerOrganization = nil
	// m.User = nil
	m.Loan = nil
	m.BorrowerApplicant = nil
}
