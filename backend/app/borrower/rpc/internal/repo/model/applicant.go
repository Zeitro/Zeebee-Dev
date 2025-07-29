package model

type Applicant struct {
	Base
	BorrowerId int64  `gorm:"column:borrower_id;not null;default:0" json:"borrower_id"`
	LoanId     int64  `gorm:"column:loan_id;not null;default:0" json:"loan_id"`
	FirstName  string `gorm:"column:first_name;size:30;not null;default:''" json:"first_name"`
	LastName   string `gorm:"column:last_name;size:30;not null;default:''" json:"last_name"`
	MidName    string `gorm:"column:mid_name;size:30;not null;default:''" json:"mid_name"`
	FullName   string `gorm:"column:full_name;size:30;not null;default:''" json:"full_name"`
	IsPrimary  bool   `gorm:"column:is_primary;not null;default:false" json:"is_primary"`

	Borrower *Borrower `json:"borrower" gorm:"foreignKey:BorrowerId;references:Id"`
	Loan     *Loan     `json:"loan" gorm:"foreignKey:LoanId;references:Id"`
}

func (a *Applicant) TableName() string {
	return "applicants"
}
