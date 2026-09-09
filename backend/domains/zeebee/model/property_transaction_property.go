package model

import "zeebee-backend/domains/zeebee/constants"

type PropertyTransactionProperty struct {
	Base
	BorrowerPropertyTransactionId int64                  `gorm:"column:borrower_property_transaction_id;not null;default:0" json:"borrower_property_transaction_id"`
	BorrowerPropertyId            int64                  `gorm:"column:borrower_property_id;not null;default:0" json:"borrower_property_id"`
	Unit                          int32                  `gorm:"column:unit;not null;default:0" json:"unit"`
	PropertyType                  constants.PropertyType `gorm:"column:property_type;size:30;not null;default:''" json:"property_type"`
	LivingAreaSqft                float32                `gorm:"column:living_area_sqft;type:numeric(8,2);not null;default:0" json:"living_area_sqft"`
	LotAreaSqft                   float32                `gorm:"column:lot_area_sqft;type:numeric(8,2);not null;default:0" json:"lot_area_sqft"`
	Bedrooms                      int32                  `gorm:"column:bedrooms;not null;default:0" json:"bedrooms"`
	Bathrooms                     int32                  `gorm:"column:bathrooms;not null;default:0" json:"bathrooms"`
	Address                       Address                `gorm:"column:address;type:jsonb;serializer:json" json:"address"`

	// Relationships
	BorrowerPropertyTransaction *BorrowerPropertyTransaction `json:"borrower_property_transaction" gorm:"foreignKey:BorrowerPropertyTransactionId;references:Id"`
	BorrowerProperty            *BorrowerProperty            `json:"borrower_property" gorm:"foreignKey:BorrowerPropertyId;references:Id"`
}

func (m PropertyTransactionProperty) TableName() string {
	return "property_transaction_properties"
}

func (m *PropertyTransactionProperty) RemoveRelationships() {
	m.BorrowerPropertyTransaction = nil
	m.BorrowerProperty = nil
}
