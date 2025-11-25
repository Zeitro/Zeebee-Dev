package model

import (
	"zeebee-backend/common/constants"
)

type BorrowerProperty struct {
	Base
	UserId         int64                  `gorm:"column:user_id;default:0" json:"user_id"`
	OrganizationId int64                  `gorm:"column:organization_id;default:0" json:"organization_id"`
	Unit           int32                  `gorm:"column:unit;not null;default:0" json:"unit"`
	PropertyType   constants.PropertyType `gorm:"column:property_type;size:30;not null;default:''" json:"property_type"`
	LivingAreaSqft float32                `gorm:"column:living_area_sqft;type:numeric(8,2);not null;default:0" json:"living_area_sqft"`
	Bedrooms       int32                  `gorm:"column:bedrooms;not null;default:0" json:"bedrooms"`
	Bathrooms      int32                  `gorm:"column:bathrooms;not null;default:0" json:"bathrooms"`
	Address        Address                `gorm:"column:address;type:jsonb;serializer:json" json:"address"`
	TotalAmount    float64                `gorm:"column:total_amount;not null;default:0" json:"total_amount"`
	LoanCount      int32                  `gorm:"column:loan_count;not null;default:0" json:"loan_count"`

	// Relationships
	User           *User            `json:"user" gorm:"foreignKey:UserId;references:Id"`
	Organization   *Organization    `json:"organization" gorm:"foreignKey:OrganizationId;references:Id"`
	Loans          []*Loan          `json:"loans" gorm:"-"`
	Applicants     []*LoanApplicant `json:"applicants" gorm:"-"`
	LoanProperties []*LoanProperty  `json:"loan_properties" gorm:"foreignKey:BorrowerPropertyId;references:Id"`
}

func (m BorrowerProperty) TableName() string {
	return "borrower_properties"
}

func (m *BorrowerProperty) RemoveRelationships() {
	m.User = nil
	m.Organization = nil
	m.Loans = nil
	m.Applicants = nil
	m.LoanProperties = nil
}
