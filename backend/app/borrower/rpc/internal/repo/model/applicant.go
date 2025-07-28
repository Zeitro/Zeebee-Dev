package model

import (
	"time"

	"gorm.io/plugin/soft_delete"
)

type Applicant struct {
	ID         int64                 `gorm:"primaryKey" json:"id"`
	BorrowerID int64                 `gorm:"column:borrower_id;not null;default:0" json:"borrower_id"`
	LoanID     int64                 `gorm:"column:loan_id;not null;default:0" json:"loan_id"`
	FirstName  string                `gorm:"column:first_name;size:30;not null;default:''" json:"first_name"`
	LastName   string                `gorm:"column:last_name;size:30;not null;default:''" json:"last_name"`
	MidName    string                `gorm:"column:mid_name;size:30;not null;default:''" json:"mid_name"`
	FullName   string                `gorm:"column:mid_name;size:30;not null;default:''" json:"full_name"`
	IsPrimary  bool                  `gorm:"column:is_primary;not null;default:false" json:"is_primary"`
	CreatedAt  time.Time             `gorm:"column:created_at;default:now()" json:"created_at"`
	UpdatedAt  time.Time             `gorm:"column:updated_at;default:now()" json:"updated_at"`
	IsDel      soft_delete.DeletedAt `gorm:"softDelete:flag" json:"-"`

	Borrower *Borrower `json:"borrower" gorm:"foreignKey:BorrowerID;references:ID"`
	Loan     *Loan     `json:"loan" gorm:"foreignKey:LoanID;references:ID"`
}

func (a *Applicant) TableName() string {
	return "applicants"
}
