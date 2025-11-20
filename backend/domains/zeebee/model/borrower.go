package model

type Borrower struct {
	Base

	UserId    int64  `gorm:"column:user_id;not null;default:0" json:"user_id"`
	FirstName string `gorm:"column:first_name;size:100;not null;default:''" json:"first_name"`
	LastName  string `gorm:"column:last_name;size:100;not null;default:''" json:"last_name"`
	FullName  string `gorm:"column:full_name;size:100;not null;default:''" json:"full_name"`

	User *User `json:"user" gorm:"foreignKey:UserId;references:Id"`
}

func (m *Borrower) TableName() string {
	return "borrowers"
}
