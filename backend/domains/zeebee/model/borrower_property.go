package model

import (
	"zeebee-backend/domains/zeebee/constants"
)

type BorrowerProperty struct {
	Base
	CreatorId      int64                  `gorm:"column:creator_id;default:0" json:"creator_id"`
	CreatorName    string                 `gorm:"column:creator_name;type:varchar(100);default:'';not null" json:"creator_name"`
	OrganizationId int64                  `gorm:"column:organization_id;default:0" json:"organization_id"`
	Unit           int32                  `gorm:"column:unit;not null;default:0" json:"unit"`
	PropertyType   constants.PropertyType `gorm:"column:property_type;size:30;not null;default:''" json:"property_type"`
	LivingAreaSqft float32                `gorm:"column:living_area_sqft;type:numeric(8,2);not null;default:0" json:"living_area_sqft"`
	LotAreaSqft    float32                `gorm:"column:lot_area_sqft;type:numeric(8,2);not null;default:0" json:"lot_area_sqft"`
	Bedrooms       int32                  `gorm:"column:bedrooms;not null;default:0" json:"bedrooms"`
	Bathrooms      int32                  `gorm:"column:bathrooms;not null;default:0" json:"bathrooms"`
	Address        Address                `gorm:"column:address;type:jsonb;serializer:json" json:"address"`
	TotalAmount    float64                `gorm:"column:total_amount;not null;default:0" json:"total_amount"`
	LoanCount      int32                  `gorm:"column:loan_count;not null;default:0" json:"loan_count"`

	// Relationships
	Creator        *User            `json:"creator" gorm:"foreignKey:CreatorId;references:Id"`
	Organization   *Organization    `json:"organization" gorm:"foreignKey:OrganizationId;references:Id"`
	Loans          []*Loan          `json:"loans" gorm:"-"`
	Applicants     []*LoanApplicant `json:"applicants" gorm:"-"`
	LoanProperties []*LoanProperty  `json:"loan_properties" gorm:"foreignKey:BorrowerPropertyId;references:Id"`
}

func (m BorrowerProperty) TableName() string {
	return "borrower_properties"
}

func (m *BorrowerProperty) RemoveRelationships() {
	m.Creator = nil
	m.Organization = nil
	m.Loans = nil
	m.Applicants = nil
	m.LoanProperties = nil
}
