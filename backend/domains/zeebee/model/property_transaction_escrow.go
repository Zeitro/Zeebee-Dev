package model

import (
	"zeebee-backend/domains/zeebee/data/dto"
)

type PropertyTransactionEscrow struct {
	Base
	BorrowerPropertyTransactionId int64                                  `gorm:"column:borrower_property_transaction_id;not null;default:0" json:"borrower_property_transaction_id"`
	No                            string                                 `gorm:"column:no;type:varchar(50);not null;default:''" json:"no"`
	EscrowCompany                 string                                 `gorm:"column:escrow_company;type:varchar(100);not null;default:''" json:"escrow_company"`
	OfficerName                   string                                 `gorm:"column:officer_name;type:varchar(30);not null;default:''" json:"officer_name"`
	OfficerPhone                  string                                 `gorm:"column:officer_phone;type:varchar(20);not null;default:''" json:"officer_phone"`
	OfficerEmail                  string                                 `gorm:"column:officer_email;type:varchar(50);not null;default:''" json:"officer_email"`
	Checklist                     dto.PropertyTransactionEscrowChecklist `gorm:"column:checklist;type:jsonb;not null;default:'{}';serializer:json" json:"checklist"`

	// Relationships
	BorrowerPropertyTransaction *BorrowerPropertyTransaction `json:"borrower_property_transaction" gorm:"foreignKey:BorrowerPropertyTransactionId;references:Id"`
}

func (m PropertyTransactionEscrow) TableName() string {
	return "property_transaction_escrows"
}
