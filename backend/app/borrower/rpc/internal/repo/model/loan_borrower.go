package model

import (
	"zeebee-backend/common/constants"
	"zeebee-backend/common/types"
)

type LoanBorrower struct {
	Base
	LoanId         int64                   `gorm:"column:loan_id;not null;default:0" json:"loan_id"`
	Type           constants.ApplicantType `gorm:"column:type;size:10;not null;default:''" json:"type"`
	Name           string                  `gorm:"column:name;size:255;not null;default:''" json:"name"`
	TinType        string                  `gorm:"column:tin_type;size:10;not null;default:'SSN'" json:"tin_type"`
	Tin            string                  `gorm:"column:tin;size:20;not null;default:''" json:"tin"`
	MailingAddress types.Address           `gorm:"column:mailing_address;type:jsonb;serializer:json" json:"mailing_address"`
	ContactEmail   string                  `gorm:"column:contact_email;size:100;not null;default:''" json:"contact_email"`
	ContactPhone   string                  `gorm:"column:contact_phone;size:20;not null;default:''" json:"contact_phone"`
	ContactName    string                  `gorm:"column:contact_name;size:100;not null;default:''" json:"contact_name"`
	OwnPercent     int32                   `gorm:"column:own_percent;not null;default:0" json:"own_percent"`

	Loan       *Loan        `json:"loan" gorm:"foreignKey:LoanId;references:Id"`
	Applicants []*Applicant `json:"applicants" gorm:"foreignKey:LoanBorrowerId;references:Id"`
}

func (m LoanBorrower) TableName() string {
	return "loan_borrowers"
}

func (m *LoanBorrower) RemoveRelationships() {
	m.Loan = nil
	m.Applicants = nil
}
