package model

import (
	"time"

	"gorm.io/plugin/soft_delete"
)

type Base struct {
	Id        int64                 `gorm:"primaryKey;<-:false" json:"id"`
	CreatedAt time.Time             `gorm:"column:created_at;default:now();<-:create" json:"created_at"`
	UpdatedAt time.Time             `gorm:"column:updated_at;default:now()" json:"updated_at"`
	IsDel     soft_delete.DeletedAt `gorm:"softDelete:flag" json:"-"`
}

func (Base) FieldId() string {
	return "id"
}
func (Base) FieldCreatedAt() string {
	return "created_at"
}
func (Base) FieldUpdatedAt() string {
	return "updated_at"
}

func (Base) FieldIsDel() string {
	return "is_del"
}
