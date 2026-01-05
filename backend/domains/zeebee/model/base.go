package model

import (
	"time"

	"gorm.io/plugin/soft_delete"
)

type Base struct {
	Id        int64                 `gorm:"primaryKey;autoIncrement;<-:create" json:"id"`
	CreatedAt time.Time             `gorm:"column:created_at;default:now();<-:create" json:"created_at"`
	UpdatedAt time.Time             `gorm:"column:updated_at;default:now()" json:"updated_at"`
	IsDel     soft_delete.DeletedAt `gorm:"softDelete:flag" json:"-"`
}

type Address struct {
	State    string `json:"state" copier:"State"`
	City     string `json:"city"`
	Address1 string `json:"address1"`
	Address2 string `json:"address2"`
	ZipCode  string `json:"zip_code"`
}
