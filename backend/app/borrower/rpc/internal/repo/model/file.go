package model

import "zeebee-backend/common/constants"

type File struct {
	Base

	Type  constants.FileType `gorm:"column:type;size:10;not null;default:'file'" json:"type"`
	RefId int64              `gorm:"column:ref_id;not null;default:0" json:"ref_id"`
	File  string             `gorm:"column:file;size:255;not null;default:''" json:"file"`
}

func (m *File) TableName() string {
	return "files"
}
