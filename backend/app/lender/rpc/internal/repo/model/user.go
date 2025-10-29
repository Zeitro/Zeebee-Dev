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

	Borrower       *Borrower       `json:"borrower" gorm:"foreignKey:UserId;references:Id"`
	LenderEmployee *LenderEmployee `json:"lender_employee" gorm:"foreignKey:UserId;references:Id"`
}

func (m *User) TableName() string {
	return "users"
}
