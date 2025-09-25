package model

import "zeebee-backend/common/constants"

type File struct {
	Base

	FileType constants.FileType    `gorm:"column:file_type;size:10;not null;default:'file'" json:"file_type"`
	BizType  constants.FileBizType `gorm:"column:biz_type;size:50;not null;default:'Loan'" json:"biz_type"`
	RefId    int64                 `gorm:"column:ref_id;not null;default:0" json:"ref_id"`
	File     string                `gorm:"column:file;size:255;not null;default:''" json:"file"`
}

func (m *File) TableName() string {
	return "files"
}
