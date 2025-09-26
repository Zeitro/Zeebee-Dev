package model

import "zeebee-backend/common/constants"

type LoanProperty struct {
	Base
	LoanId               int64                  `gorm:"column:loan_id;not null;default:0" json:"loan_id"`
	PropertyBaseId       int64                  `gorm:"column:property_base_id;not null;default:0" json:"property_base_id"`
	Unit                 int32                  `gorm:"column:unit;not null;default:1" json:"unit"`
	PropertyType         constants.PropertyType `gorm:"column:property_type;size:30;not null;default:''" json:"property_type"`
	LivingAreaSqft       float32                `gorm:"column:living_area_sqft;type:numeric(8,2);not null;default:0" json:"living_area_sqft"`
	Bedrooms             int32                  `gorm:"column:bedrooms;not null;default:0" json:"bedrooms"`
	Bathrooms            int32                  `gorm:"column:bathrooms;not null;default:0" json:"bathrooms"`
	FuturePropertyType   constants.PropertyType `gorm:"column:future_property_type;size:30;not null;default:''" json:"future_property_type"`
	FutureUnit           int32                  `gorm:"column:future_unit;not null;default:0" json:"future_unit"`
	FutureLivingAreaSqft float32                `gorm:"column:future_living_area_sqft;type:numeric(8,2);not null;default:0" json:"future_living_area_sqft"`
	FutureBedrooms       int32                  `gorm:"column:future_bedrooms;not null;default:0" json:"future_bedrooms"`
	FutureBathrooms      int32                  `gorm:"column:future_bathrooms;not null;default:0" json:"future_bathrooms"`
	Owner                string                 `gorm:"column:owner;size:255;not null;default:''" json:"owner"`
	// relationships
	PropertyBase *PropertyBase `json:"property_base" gorm:"foreignKey:PropertyBaseId;references:Id"`
	Loan         *Loan         `json:"loan" gorm:"foreignKey:LoanId;references:Id"`
	Applicants   []*Applicant  `json:"applicants" gorm:"foreignKey:LoanPropertyId;references:Id"`
}

func (m LoanProperty) TableName() string {
	return "loan_properties"
}

func (m *LoanProperty) RemoveRelationships() {
	m.PropertyBase = nil
	m.Loan = nil
	m.Applicants = nil
}
