package model

import (
	"database/sql"

	"github.com/lib/pq"

	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
)

type BorrowerPropertyTransaction struct {
	Base
	BorrowerPropertyId     int64                                       `gorm:"column:borrower_property_id;not null;default:0" json:"borrower_property_id"`
	OrganizationId         int64                                       `gorm:"column:organization_id;not null;default:0" json:"organization_id"`
	Code                   string                                      `gorm:"column:code;type:varchar(20);not null;default:''" json:"code"`
	Type                   constants.BorrowerPropertyTransactionType   `gorm:"column:type;type:varchar(20);not null;default:'Purchase'" json:"type"`
	Status                 constants.BorrowerPropertyTransactionStatus `gorm:"column:status;type:varchar(20);not null;default:'Pending'" json:"status"`
	CoeDate                *sql.NullTime                               `gorm:"column:coe_date;type:date" json:"coe_date"`
	PurchaseAmount         float64                                     `gorm:"column:purchase_amount;type:decimal(10,2);not null;default:0" json:"purchase_amount"`
	EmdAmount              float64                                     `gorm:"column:emd_amount;type:decimal(10,2);not null;default:0" json:"emd_amount"`
	PropertyType           constants.PropertyType                      `gorm:"column:property_type;type:varchar(50);not null;default:''" json:"property_type"`
	Pm                     string                                      `gorm:"column:pm;type:varchar(255);not null;default:''" json:"pm"`
	BusinessStrategy       pq.StringArray                              `gorm:"column:business_strategy;type:text[]" json:"business_strategy"`
	Note                   string                                      `gorm:"column:note;type:varchar(255);not null;default:''" json:"note"`
	TitleHolder            string                                      `gorm:"column:title_holder;type:varchar(100);not null;default:''" json:"title_holder"`
	ContractDate           *sql.NullTime                               `gorm:"column:contract_date;type:date" json:"contract_date"`
	EmdDueDate             *sql.NullTime                               `gorm:"column:emd_due_date;type:date" json:"emd_due_date"`
	ContingencyRemovalDate *sql.NullTime                               `gorm:"column:contingency_removal_date;type:date" json:"contingency_removal_date"`
	Tc                     string                                      `gorm:"column:tc;type:varchar(50);not null;default:''" json:"tc"`
	CostAllocation         string                                      `gorm:"column:cost_allocation;type:varchar(255);not null;default:''" json:"cost_allocation"`
	Checklist              dto.ChecklistRoot                           `gorm:"column:checklist;type:jsonb;not null;default:'{}';serializer:json" json:"checklist"`

	// Relationships
	BorrowerProperty *BorrowerProperty            `json:"borrower_property" gorm:"foreignKey:BorrowerPropertyId;references:Id"`
	Escrows          []*PropertyTransactionEscrow `json:"escrows" gorm:"foreignKey:BorrowerPropertyTransactionId;references:Id"`
	Property         *PropertyTransactionProperty `json:"property" gorm:"foreignKey:BorrowerPropertyTransactionId;references:Id"`
}

func (m BorrowerPropertyTransaction) TableName() string {
	return "borrower_property_transactions"
}

func (m *BorrowerPropertyTransaction) RemoveRelationships() {
	m.BorrowerProperty = nil
	m.Escrows = nil
	m.Property = nil
}
