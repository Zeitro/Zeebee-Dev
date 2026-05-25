package model

import "zeebee-backend/domains/zeebee/constants"

type EntityDocument struct {
	Base
	OrganizationId int64                        `gorm:"column:organization_id;not null;default:0" json:"organization_id"`
	EntityType     constants.EntityDocumentType `gorm:"column:entity_type;type:varchar(50);not null;default:''" json:"entity_type"`
	EntityId       int64                        `gorm:"column:entity_id;not null;default:0" json:"entity_id"`
	RootDocumentId int64                        `gorm:"column:root_document_id;not null;default:0" json:"root_document_id"`

	// Relationships
	RootDocument *Document `json:"root_document" gorm:"foreignKey:RootDocumentId;references:Id"`
}

func (m EntityDocument) TableName() string {
	return "entity_documents"
}
