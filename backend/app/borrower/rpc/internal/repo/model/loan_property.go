package model

import "zeebee-backend/common/constants"

type LoanProperty struct {
	Base
	LoanId         int64                  `gorm:"column:loan_id;not null;default:0" json:"loan_id"`
	PropertyBaseId int64                  `gorm:"column:property_base_id;not null;default:0" json:"property_base_id"`
	PropertyType   constants.PropertyType `gorm:"column:property_type;size:30;not null;default:''" json:"property_type"`
	LivingAreaSqft float32                `gorm:"column:living_area_sqft;type:numeric(8,2);not null;default:0" json:"living_area_sqft"`
	Bedrooms       int32                  `gorm:"column:bedrooms;not null;default:0" json:"bedrooms"`
	Bathrooms      int32                  `gorm:"column:bathrooms;not null;default:0" json:"bathrooms"`
	// relationships
	PropertyBase *PropertyBase `json:"property_base" gorm:"foreignKey:PropertyBaseId;references:Id"`
	Loan         *Loan         `json:"loan" gorm:"foreignKey:LoanId;references:Id"`
	Applicants   []*Applicant  `json:"applicants" gorm:"foreignKey:LoanPropertyId;references:Id"`
}

func (m LoanProperty) TableName() string {
	return "loan_properties"
}
