package model

import "zeebee-backend/common/constants"

// LoanThirdService 三方服务公司
type LoanThirdService struct {
	Base
	LoanId  int64                          `gorm:"column:loan_id;not null;default:0" json:"loan_id"`
	Type    constants.LoanThirdServiceType `gorm:"column:type;size:20;not null;default:''" json:"type"`
	Name    string                         `gorm:"column:name;size:100;not null;default:''" json:"name"`
	Email   string                         `gorm:"column:email;size:100;not null;default:''" json:"email"`
	Phone   string                         `gorm:"column:phone;size:20;not null;default:''" json:"phone"`
	Company string                         `gorm:"column:company;size:100;not null;default:''" json:"company"`
	// relationships
	Loan *Loan `json:"loan" gorm:"foreignKey:LoanId;references:Id"`
}

func (m LoanThirdService) TableName() string {
	return "loan_third_services"
}

func (m *LoanThirdService) RemoveRelationships() {
	m.Loan = nil
}
