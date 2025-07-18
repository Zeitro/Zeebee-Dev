package model

import (
	"time"

	"gorm.io/plugin/soft_delete"

	"zeebee-backend/common/constants"
)

type User struct {
	ID           int64                 `json:"id"`
	Type         constants.UserType    `json:"type"`
	Email        string                `json:"email"`
	PasswordHash string                `json:"-"`
	Status       constants.UserStatus  `json:"status"`
	CreatedAt    time.Time             `json:"created_at"`
	UpdatedAt    time.Time             `json:"updated_at"`
	IsDel        soft_delete.DeletedAt `gorm:"softDelete:flag" json:"-"`

	Borrower *Borrower `json:"borrower" gorm:"foreignKey:UserID;references:ID"`
}

func (u *User) TableName() string {
	return "users"
}
