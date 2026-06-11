package model

import "zeebee-backend/domains/zeebee/constants"

type DocumentTemplate struct {
	Base
	Name      string                       `gorm:"column:name;type:varchar(255);not null;default:''" json:"name"`
	Type      string                       `gorm:"column:type;type:varchar(20);not null;default:''" json:"type"`
	Structure DocumentTemplateStructureMap `gorm:"column:structure;type:json;not null;default:'{}';serializer:json" json:"structure"`
	Status    constants.CommonStatus       `gorm:"column:status;type:varchar(20);not null;default:''" json:"status"`
}

func (m DocumentTemplate) TableName() string {
	return "document_templates"
}

type DocumentTemplateStructureMap map[string][]DocumentTemplateStructure

type DocumentTemplateStructure struct {
	Name       string                      `json:"name"`
	SortOrder  int32                       `json:"sort_order"`
	IsRequired bool                        `json:"is_required"`
	Children   []DocumentTemplateStructure `json:"children"`
}
