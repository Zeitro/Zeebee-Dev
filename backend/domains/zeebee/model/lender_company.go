package model

// import (
// 	"zeebee-backend/common/constants"
// 	"zeebee-backend/common/types"
// )
//
// type LenderCompany struct {
// 	Base
//
// 	Name         string                           `gorm:"column:name;size:100;not null;default:''" json:"name"`
// 	Address      types.Address                    `gorm:"column:address;type:jsonb;serializer:json" json:"address"`
// 	Status       constants.LenderCompanyStatus    `gorm:"column:status;size:10;not null;default:'active'" json:"status"`
// 	ContactEmail string                           `gorm:"column:contact_email;size:100;not null;default:''" json:"contact_email"`
// 	OwnerType    constants.LenderCompanyOwnerType `gorm:"column:owner_type;size:30;not null;default:'ThirdParty'" json:"owner_type"`
//
// 	Loans []Loan `json:"loans" gorm:"foreignKey:LenderCompanyId;references:Id"`
// }
//
// func (m *LenderCompany) TableName() string {
// 	return "lender_companies"
// }
