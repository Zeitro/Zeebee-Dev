package model

import (
	"database/sql"

	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
)

type BorrowerPropertyTransaction struct {
	Base
	BorrowerPropertyId int64                                       `gorm:"column:borrower_property_id;not null;default:0" json:"borrower_property_id"`
	OrganizationId     int64                                       `gorm:"column:organization_id;not null;default:0" json:"organization_id"`
	Code               string                                      `gorm:"column:code;type:varchar(20);not null;default:''" json:"code"`
	Type               constants.BorrowerPropertyTransactionType   `gorm:"column:type;type:varchar(20);not null;default:'Purchase'" json:"type"`
	Status             constants.BorrowerPropertyTransactionStatus `gorm:"column:status;type:varchar(20);not null;default:'Pending'" json:"status"`
	CoeDate            *sql.NullTime                               `gorm:"column:coe_date;type:date" json:"coe_date"`
	PurchaseAmount     float64                                     `gorm:"column:purchase_amount;type:decimal(10,2);not null;default:0" json:"purchase_amount"`
	EmdAmount          float64                                     `gorm:"column:emd_amount;type:decimal(10,2);not null;default:0" json:"emd_amount"`
	PropertyType       constants.PropertyType                      `gorm:"column:property_type;type:varchar(50);not null;default:''" json:"property_type"`
	Pm                 string                                      `gorm:"column:pm;type:varchar(255);not null;default:''" json:"pm"`
	Strategy           string                                      `gorm:"column:strategy;type:varchar(100);not null;default:''" json:"strategy"`
	Note               string                                      `gorm:"column:note;type:varchar(255);not null;default:''" json:"note"`
	Checklist          dto.ChecklistRoot                           `gorm:"column:checklist;type:jsonb;not null;default:'{}';serializer:json" json:"checklist"`

	// Relationships
	BorrowerProperty *BorrowerProperty `json:"borrower_property" gorm:"foreignKey:BorrowerPropertyId;references:Id"`
}

func (m BorrowerPropertyTransaction) TableName() string {
	return "borrower_property_transactions"
}
