package model

import (
	"time"

	"gorm.io/datatypes"

	"zeebee-backend/domains/zeebee/constants"
)

type Document struct {
	Base
	ParentId       int64                  `gorm:"column:parent_id;not null;default:0" json:"parent_id"`
	Type           constants.DocumentType `gorm:"column:type;type:varchar(20);not null;default:''" json:"type"`
	Name           string                 `gorm:"column:name;type:varchar(255);not null" json:"name"`
	Level          int32                  `gorm:"column:level;not null;default:0" json:"level"`
	SortOrder      int32                  `gorm:"column:sort_order;not null;default:9" json:"sort_order"`
	LastUploadedAt *time.Time             `gorm:"column:last_uploaded_at" json:"last_uploaded_at"`
	Metadata       datatypes.JSON         `gorm:"column:metadata;type:jsonb;not null;default:'{}'" json:"metadata"`

	// Relationships
	Children []*Document `json:"children" gorm:"foreignKey:ParentId;references:Id"`
	Files    []*File     `json:"files" gorm:"-"`
}

func (m Document) TableName() string {
	return "documents"
}
