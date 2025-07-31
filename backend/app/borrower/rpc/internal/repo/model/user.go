package model

import (
	"zeebee-backend/common/constants"
)

type User struct {
	Base
	Type         constants.UserType   `json:"type"`
	Email        string               `json:"email"`
	PasswordHash string               `json:"-"`
	Status       constants.UserStatus `json:"status"`

	Borrower *Borrower `json:"borrower" gorm:"foreignKey:UserId;references:Id"`
}

func (m *User) TableName() string {
	return "users"
}

func (m User) FieldType() string {
	return "type"
}

func (m User) FieldEmail() string {
	return "email"
}

func (m User) FieldPasswordHash() string {
	return "password_hash"
}

func (m User) FieldStatus() string {
	return "status"
}
