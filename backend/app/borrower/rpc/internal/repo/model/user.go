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

func (u *User) TableName() string {
	return "users"
}
