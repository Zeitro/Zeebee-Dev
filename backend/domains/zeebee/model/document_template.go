package model

import (
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
)

type DocumentTemplate struct {
	Base
	Name      string                           `gorm:"column:name;type:varchar(255);not null;default:''" json:"name"`
	Type      string                           `gorm:"column:type;type:varchar(20);not null;default:''" json:"type"`
	Structure dto.DocumentTemplateStructureMap `gorm:"column:structure;type:json;not null;default:'{}';serializer:json" json:"structure"`
	Status    constants.CommonStatus           `gorm:"column:status;type:varchar(20);not null;default:''" json:"status"`
}

func (m DocumentTemplate) TableName() string {
	return "document_templates"
}
