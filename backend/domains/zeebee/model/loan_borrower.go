package model

import (
	"zeebee-backend/domains/zeebee/constants"
)

type LoanBorrower struct {
	Base
	LoanId           int64                   `gorm:"column:loan_id;not null;default:0" json:"loan_id"`
	BorrowerEntityId int64                   `gorm:"column:borrower_entity_id;not null;default:0" json:"borrower_entity_id"`
	Type             constants.ApplicantType `gorm:"column:type;size:10;not null;default:''" json:"type"`
	Name             string                  `gorm:"column:name;size:255;not null;default:''" json:"name"`
	TinType          string                  `gorm:"column:tin_type;size:10;not null;default:'SSN'" json:"tin_type"`
	Tin              string                  `gorm:"column:tin;size:20;not null;default:''" json:"tin"`
	MailingAddress   Address                 `gorm:"column:mailing_address;type:jsonb;serializer:json" json:"mailing_address"`
	ContactEmail     string                  `gorm:"column:contact_email;size:100;not null;default:''" json:"contact_email"`
	ContactPhone     string                  `gorm:"column:contact_phone;size:20;not null;default:''" json:"contact_phone"`
	ContactName      string                  `gorm:"column:contact_name;size:100;not null;default:''" json:"contact_name"`

	Loan           *Loan            `json:"loan" gorm:"foreignKey:LoanId;references:Id"`
	Applicants     []*LoanApplicant `json:"applicants" gorm:"foreignKey:LoanBorrowerId;references:Id"`
	BorrowerEntity *BorrowerEntity  `json:"borrower_entity" gorm:"foreignKey:BorrowerEntityId;references:Id"`
}

func (m LoanBorrower) TableName() string {
	return "loan_borrowers"
}

func (m *LoanBorrower) RemoveRelationships() {
	m.Loan = nil
	m.Applicants = nil
	m.BorrowerEntity = nil
}
